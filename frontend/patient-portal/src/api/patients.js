import api from './index'

export const patientsApi = {
  // Get patients with filtering and pagination
  async getPatients(params = {}) {
    const response = await api.get('/api/patients/', { params })
    return response.data
  },

  // Get single patient by ID
  async getPatient(id) {
    const response = await api.get(`/api/patients/${id}`)
    return response.data
  },

  // Create new patient
  async createPatient(patientData) {
    const response = await api.post('/api/patients/', patientData)
    return response.data
  },

  // Update patient
  async updatePatient(id, patientData) {
    const response = await api.put(`/api/patients/${id}`, patientData)
    return response.data
  },

  // Delete patient
  async deletePatient(id) {
    const response = await api.delete(`/api/patients/${id}`)
    return response.data
  },

  // Search patients
  async searchPatients(query, params = {}) {
    const response = await api.get('/api/patients/search', { 
      params: { q: query, ...params } 
    })
    return response.data
  },

  // Get patient statistics
  async getPatientStats() {
    const response = await api.get('/api/patients/stats')
    return response.data
  },

  // Export patients data
  async exportPatients(format = 'csv', params = {}) {
    const response = await api.get('/api/patients/export', { 
      params: { format, ...params },
      responseType: 'blob'
    })
    return response.data
  }
}