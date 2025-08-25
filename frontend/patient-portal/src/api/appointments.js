import api from './index.js'

export const appointmentsApi = {
  // Appointments CRUD
  async getAppointments(params = {}) {
    const response = await api.get('/api/appointments/', { params })
    return response.data
  },

  async getAppointment(id) {
    const response = await api.get(`/api/appointments/${id}`)
    return response.data
  },

  async createAppointment(appointmentData) {
    const response = await api.post('/api/appointments/', appointmentData)
    return response.data
  },

  async updateAppointment(id, appointmentData) {
    const response = await api.put(`/api/appointments/${id}`, appointmentData)
    return response.data
  },

  async updateAppointmentStatus(id, status, notes = '') {
    const response = await api.patch(`/api/appointments/${id}/status`, { status, notes })
    return response.data
  },

  async deleteAppointment(id) {
    const response = await api.delete(`/api/appointments/${id}`)
    return response.data
  },

  // Doctor Availability Management
  async getDoctorAvailability(params = {}) {
    const response = await api.get('/api/availability/', { params })
    return response.data
  },

  async getDoctorAvailabilityById(id) {
    const response = await api.get(`/api/availability/${id}`)
    return response.data
  },

  async createDoctorAvailability(availabilityData) {
    const response = await api.post('/api/availability/', availabilityData)
    return response.data
  },

  async updateDoctorAvailability(id, availabilityData) {
    const response = await api.put(`/api/availability/${id}`, availabilityData)
    return response.data
  },

  async deleteDoctorAvailability(id) {
    const response = await api.delete(`/api/availability/${id}`)
    return response.data
  },

  // Calendar and Bulk Operationse/${en
  async getAvailabilityCalendar(doctorId, month) {
    const response = await api.get('/api/availability/calendar', {
      params: { doctor_id: doctorId, month }
    })
    return response.data
  },

  async createBulkAvailability(bulkData) {
    const response = await api.post('/api/availability/bulk', bulkData)
    return response.data
  },

  // Doctor Management
  async getDoctorsByEntity() {
    const response = await api.get('/api/doctors/')
    return response.data
  },

  // Doctor Schedules (for appointment slot availability)
  async getDoctorSchedule(doctorId, date) {
    const response = await api.get('/api/schedules/', {
      params: { doctor_id: doctorId, date }
    })
    return response.data
  },

  // Smart Booking & Time Slots
  async bookAppointment(bookingData) {
    const response = await api.post('/api/appointments/book', bookingData)
    return response.data
  },

  async getAvailableTimeSlots(doctorId, date, duration = 30) {
    const response = await api.get('/api/appointments/slots', {
      params: { doctor_id: doctorId, date, duration }
    })
    return response.data
  },

  // Admin Management API
  
  // Duration Settings
  async getDurationSettings() {
    const response = await api.get('/api/admin/duration-settings')
    return response.data
  },

  async createDurationSetting(durationData) {
    const response = await api.post('/api/admin/duration-settings', durationData)
    return response.data
  },

  async updateDurationSetting(settingId, durationData) {
    const response = await api.put(`/api/admin/duration-settings/${settingId}`, durationData)
    return response.data
  },

  // Room Management
  async getRooms(roomType = '', floor = '') {
    const response = await api.get('/api/appointments/rooms', {
      params: { room_type: roomType, floor }
    })
    return response.data
  },

  async createRoom(roomData) {
    const response = await api.post('/api/admin/rooms', roomData)
    return response.data
  },

  async updateRoom(roomId, roomData) {
    const response = await api.put(`/api/admin/rooms/${roomId}`, roomData)
    return response.data
  },

  async deleteRoom(roomId) {
    const response = await api.delete(`/api/admin/rooms/${roomId}`)
    return response.data
  },

  async getAvailableRooms(dateTime, duration = 30, roomType = '') {
    const response = await api.get('/api/appointments/available-rooms', {
      params: { date_time: dateTime, duration, room_type: roomType }
    })
    return response.data
  },

  // Duration Options Management
  async getDurationOptions(appointmentType = '') {
    const response = await api.get('/api/appointments/duration-options', {
      params: { appointment_type: appointmentType }
    })
    return response.data
  },

  async createDurationOption(optionData) {
    const response = await api.post('/api/admin/duration-options', optionData)
    return response.data
  },

  async updateDurationOption(optionId, optionData) {
    const response = await api.put(`/api/admin/duration-options/${optionId}`, optionData)
    return response.data
  },

  async deleteDurationOption(optionId) {
    const response = await api.delete(`/api/admin/duration-options/${optionId}`)
    return response.data
  },

  // Healthcare Entity Settings
}