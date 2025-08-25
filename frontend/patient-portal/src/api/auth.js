import api from './index'

export const authApi = {
  // Login user
  async login(credentials) {
    const response = await api.post('/api/auth/login', credentials)
    return response.data
  },

  // Register new user
  async register(userData) {
    const response = await api.post('/api/auth/register', userData)
    return response.data
  },

  // Refresh JWT token
  async refreshToken() {
    const response = await api.post('/api/auth/refresh')
    return response.data
  },

  // Get current user profile
  async getProfile() {
    const response = await api.get('/api/users/profile')
    return response.data
  },

  // Update user profile
  async updateProfile(profileData) {
    const response = await api.put('/api/users/profile', profileData)
    return response.data
  },

  // Change password
  async changePassword(passwordData) {
    const response = await api.post('/api/auth/change-password', passwordData)
    return response.data
  },

  // Logout (if server-side logout is needed)
  async logout() {
    // For now, we'll handle logout client-side only
    // In the future, this could invalidate tokens on the server
    return Promise.resolve()
  }
}