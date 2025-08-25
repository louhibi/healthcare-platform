import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import api from '../api'; // Use shared API gateway client
import { authApi } from '../api/auth'
import { useEntityStore } from './entity'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('auth_token'))
  const isLoading = ref(false)
  const isTempPassword = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const userName = computed(() => {
    if (!user.value) return ''
    return `${user.value.first_name} ${user.value.last_name}`
  })
  
  const healthcareEntityId = computed(() => {
    return user.value?.healthcare_entity_id || null
  })

  // Get entity store for accessing healthcare entity data
  const getEntityStore = () => useEntityStore()

  // Actions
  const login = async (credentials, router = null) => {
    isLoading.value = true
    try {
      const response = await authApi.login(credentials)
      
      // Validate response structure - authApi returns the data directly, not wrapped
      if (!response) {
        console.error('No response received from API')
        throw new Error('Invalid response from server')
      }
      
      if (!response.access_token) {
        console.error('Missing access_token in response:', response)
        throw new Error('No access token received from server')
      }
      
      if (!response.user) {
        console.error('Missing user data in response:', response)
        throw new Error('No user data received from server')
      }
      
      // Store token and user data (authApi returns data directly)
      token.value = response.access_token
      user.value = response.user
      
      // Persist token in localStorage
      localStorage.setItem('auth_token', token.value)
      
      // Set default authorization header for future requests
      setAuthHeader(token.value)
      
      
      // Fetch complete healthcare entity information using entity store
      const entityStore = getEntityStore()
      await entityStore.fetchEntity(user.value.healthcare_entity_id)
      
      // Check if user has temporary password
      checkTempPassword()
      
      // Handle redirect after successful login
      if (router) {
        const intendedRoute = localStorage.getItem('intended_route')
        if (intendedRoute) {
          localStorage.removeItem('intended_route')
          router.push(intendedRoute)
        } else {
          router.push('/')
        }
      }
      
      return response
    } catch (error) {
      // Clear any existing auth data on login failure
      logout()
      // Always throw a generic error message for security
      throw new Error('Authentication failed')
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData) => {
    isLoading.value = true
    try {
      const response = await authApi.register(userData)
      
      // Validate response structure - authApi returns the data directly
      if (!response) {
        throw new Error('Invalid response from server')
      }
      
      if (!response.access_token) {
        throw new Error('No access token received from server')
      }
      
      if (!response.user) {
        throw new Error('No user data received from server')
      }
      
      // Store token and user data (authApi returns data directly)
      token.value = response.access_token
      user.value = response.user
      
      // Persist token in localStorage
      localStorage.setItem('auth_token', token.value)
      
      // Set default authorization header for future requests
      setAuthHeader(token.value)
      
      return response
    } catch (error) {
      // Always throw a generic error message for security
      throw new Error('Registration failed')
    } finally {
      isLoading.value = false
    }
  }

  const logout = () => {
    // Clear state
    user.value = null
    token.value = null
    
    // Clear entity store
    const entityStore = getEntityStore()
    entityStore.clearEntity()
    
    // Clear localStorage
    localStorage.removeItem('auth_token')
    
    // Clear authorization header
    clearAuthHeader()
  }

  const refreshToken = async () => {
    try {
      const response = await authApi.refreshToken()
      
      // Update token (authApi returns data directly)
      token.value = response.access_token
      
      // Persist new token
      localStorage.setItem('auth_token', token.value)
      
      // Update authorization header
      setAuthHeader(token.value)
      
      return response
    } catch (error) {
      // If refresh fails, logout user
      logout()
      throw error
    }
  }

  const fetchProfile = async () => {
    try {
      const response = await authApi.getProfile()
      user.value = response // authApi returns data directly, not wrapped
      return response
    } catch (error) {
      // If profile fetch fails, logout user
      console.error('Profile fetch failed:', error)
      logout()
      throw error
    }
  }

  const updateProfile = async (profileData) => {
    isLoading.value = true
    try {
      const response = await authApi.updateProfile(profileData)
      user.value = { ...user.value, ...response } // authApi returns data directly, not wrapped
      return response
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const setAuthHeader = (authToken) => {
    // This will be used by the API service to set default headers
    if (authToken) {
      localStorage.setItem('auth_token', authToken)
    }
  }

  const clearAuthHeader = () => {
    localStorage.removeItem('auth_token')
  }

  // Initialize store - check if user is already logged in
  const initialize = async () => {
    const storedToken = localStorage.getItem('auth_token')
    if (storedToken) {
      token.value = storedToken
      setAuthHeader(storedToken)
      
      try {
        // Try to fetch user profile to validate token
        await fetchProfile()
        // Fetch healthcare entity information for existing session using entity store
        const entityStore = getEntityStore()
        await entityStore.fetchEntity(user.value.healthcare_entity_id)
        // Check if user has temporary password
        checkTempPassword()
      } catch (error) {
        // Token is invalid, clear it
        logout()
      }
    }
  }

  // Permission helpers
  const hasRole = (requiredRole) => {
    return userRole.value === requiredRole
  }

  const hasAnyRole = (roles) => {
    return roles.includes(userRole.value)
  }

  const isAdmin = computed(() => userRole.value === 'admin')
  const isDoctor = computed(() => userRole.value === 'doctor')
  const isNurse = computed(() => userRole.value === 'nurse')
  const isStaff = computed(() => userRole.value === 'staff')


  // Helper to decode JWT token
  const decodeJWT = (token) => {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload
    } catch (error) {
      console.error('Failed to decode JWT token:', error)
      return null
    }
  }

  // Helper to check if user has temporary password
  const checkTempPassword = () => {
    if (!token.value) {
      isTempPassword.value = false
      return
    }
    
    const payload = decodeJWT(token.value)
    isTempPassword.value = payload?.is_temp_password === true
  }

  // Password change function
  const changePassword = async (currentPassword, newPassword) => {
    isLoading.value = true
    try {
      await authApi.changePassword({
        current_password: currentPassword,
        new_password: newPassword
      })
      
      // After successful password change, update temp password flag
      isTempPassword.value = false
      
      return true
    } catch (error) {
      console.error('Password change failed:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    user,
    token,
    isLoading,
    isTempPassword,
    
    // Getters
    isAuthenticated,
    userRole,
    userName,
    healthcareEntityId,
    isAdmin,
    isDoctor,
    isNurse,
    isStaff,
    getEntityStore,
    
    // Actions
    login,
    register,
    logout,
    refreshToken,
    fetchProfile,
    updateProfile,
    initialize,
    hasRole,
    hasAnyRole,
    changePassword,
    checkTempPassword
  }
})