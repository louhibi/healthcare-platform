/**
 * Form Configuration API
 * Handles dynamic form configuration loading
 */

import apiClient from './index'

export const formConfigApi = {
  /**
   * Get form configuration for a specific form type
   * @param {string} formType - Type of form (e.g., 'patient', 'appointment')
   * @returns {Promise<Object>} Form configuration object
   */
  async getFormConfiguration(formType) {
    try {
      const response = await apiClient.get(`/form-config/${formType}`)
      return response.data
    } catch (error) {
      throw new Error(`Failed to load form configuration: ${error.message}`)
    }
  }
}