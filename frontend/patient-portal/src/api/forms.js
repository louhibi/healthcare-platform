import api from './index.js'

export const formsApi = {
  /**
   * Get all available form types
   * @returns {Promise} Response with form types
   */
  async getFormTypes() {
    const response = await api.get('/api/forms/types')
    return response.data
  },

  /**
   * Get complete form configuration metadata
   * @param {string} formType - The form type (e.g., 'patient', 'appointment')
   * @returns {Promise} Response with complete form metadata
   */
  async getFormMetadata(formType) {
    const response = await api.get(`/api/forms/${formType}/metadata`)
    return response.data
  },

  /**
   * Get field configurations for a specific form
   * @param {string} formType - The form type (e.g., 'patient', 'appointment')
   * @returns {Promise} Response with field configurations
   */
  async getFormFields(formType) {
    const response = await api.get(`/api/forms/${formType}/fields`)
    return response.data
  },

  /**
   * Update a single field configuration
   * @param {string} formType - The form type
   * @param {number} fieldId - The field ID to update
   * @param {Object} fieldConfig - The field configuration updates
   * @returns {Promise} Response with updated field
   */
  async updateFormField(formType, fieldId, fieldConfig) {
    const response = await api.put(`/api/forms/${formType}/fields/${fieldId}`, fieldConfig)
    return response.data
  },

  /**
   * Update multiple field configurations
   * @param {string} formType - The form type
   * @param {Array} fieldsConfig - Array of field configurations to update
   * @returns {Promise} Response with updated fields
   */
  async updateFormFields(formType, fieldsConfig) {
    const response = await api.put(`/api/forms/${formType}/fields`, {
      fields: fieldsConfig
    })
    return response.data
  },

  /**
   * Reset form configuration to defaults
   * @param {string} formType - The form type to reset
   * @returns {Promise} Response confirming reset
   */
  async resetFormToDefaults(formType) {
    const response = await api.post(`/api/forms/${formType}/reset`)
    return response.data
  },

  /**
   * Update field sort order
   * @param {string} formType - The form type
   * @param {Array} fieldOrders - Array of {field_id, sort_order} objects
   * @returns {Promise} Response with updated field orders
   */
  async updateFieldOrder(formType, fieldOrders) {
    const response = await api.put(`/api/forms/${formType}/fields/order`, {
      field_orders: fieldOrders
    })
    return response.data
  }
}