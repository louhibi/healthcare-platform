import api from './index'

export const adminApi = {
  // Create doctor with temporary password
  async createDoctor(doctorData) {
    const response = await api.post('/api/admin/doctors', doctorData)
    return response.data
  },

  // Get all doctors (future functionality)
  async getDoctors(params = {}) {
    const response = await api.get('/api/admin/doctors', { params })
    return response.data
  },

  // Update doctor (future functionality)
  async updateDoctor(doctorId, doctorData) {
    const response = await api.put(`/api/admin/doctors/${doctorId}`, doctorData)
    return response.data
  },

  // Delete/deactivate doctor (future functionality)
  async deleteDoctor(doctorId) {
    const response = await api.delete(`/api/admin/doctors/${doctorId}`)
    return response.data
  },

  // Get doctor details (future functionality)
  async getDoctor(doctorId) {
    const response = await api.get(`/api/admin/doctors/${doctorId}`)
    return response.data
  }
}