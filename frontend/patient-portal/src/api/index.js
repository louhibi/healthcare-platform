import axios from 'axios'

// Create axios instance with base configuration  
const base = import.meta.env.VITE_API_URL
const api = axios.create({
  baseURL: base,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token and user headers
api.interceptors.request.use(
  (config) => {
    
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
      
      // Parse JWT token to get user info for microservices
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        if (payload.user_id) {
          config.headers['X-User-ID'] = payload.user_id.toString()
        }
        if (payload.email) {
          config.headers['X-User-Email'] = payload.email
        }
        if (payload.role) {
          config.headers['X-User-Role'] = payload.role
        }
        if (payload.healthcare_entity_id) {
          config.headers['X-Healthcare-Entity-ID'] = payload.healthcare_entity_id.toString()
        }
      } catch (e) {
        console.warn('Failed to parse JWT token:', e)
      }
    }
    return config
  },
  (error) => {
    console.error('Request interceptor error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors globally
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    console.error('API response error:', {
      message: error.message,
      response: error.response
    })
    
    if (error.response) {
      // Server responded with error status
      const { status, data } = error.response
      
      // Handle 401 Unauthorized - token expired or invalid
      if (status === 401) {
        localStorage.removeItem('auth_token')
        // Don't automatically redirect, let the auth store handle it
        return Promise.reject(new Error('Session expired. Please login again.'))
      }
      
      // Handle other HTTP errors
      const errorMessage = data?.message || data?.error || `HTTP ${status} Error`
      return Promise.reject(new Error(errorMessage))
    } else if (error.request) {
      // Network error
      console.error('Network error details:', error.request)
      return Promise.reject(new Error('Network error. Please check your connection.'))
    } else {
      // Something else happened
      return Promise.reject(new Error('An unexpected error occurred.'))
    }
  }
)

export default api