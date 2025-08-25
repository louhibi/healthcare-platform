import api from './index.js'

// I18n API endpoints
const i18nAPI = {
  // Get supported locales
  async getSupportedLocales() {
    try {
      const response = await api.get('/api/i18n/locales')
      return response.data
    } catch (error) {
      console.error('Error fetching supported locales:', error)
      throw error
    }
  },

  // Update user's preferred locale
  async updateUserLocale(locale) {
    try {
      const response = await api.put('/api/i18n/user/locale', {
        preferred_locale: locale
      })
      return response.data
    } catch (error) {
      console.error('Error updating user locale:', error)
      throw error
    }
  },

  // Get translations for a specific locale
  async getTranslations(locale, context = null) {
    try {
      const params = {}
      if (context) {
        params.context = context
      }
      
      const response = await api.get(`/api/i18n/translations/${locale}`, { params })
      return response.data
    } catch (error) {
      console.error('Error fetching translations:', error)
      throw error
    }
  },

  // Create or update a translation (admin only)
  async createOrUpdateTranslation(translationData) {
    try {
      const response = await api.post('/api/i18n/translations', translationData)
      return response.data
    } catch (error) {
      console.error('Error creating/updating translation:', error)
      throw error
    }
  },

  // Create or update field translation (admin only)
  async createOrUpdateFieldTranslation(fieldTranslationData) {
    try {
      const response = await api.post('/api/i18n/field-translations', fieldTranslationData)
      return response.data
    } catch (error) {
      console.error('Error creating/updating field translation:', error)
      throw error
    }
  },

  // Get localized form metadata
  async getLocalizedFormMetadata(formType, locale = null) {
    try {
      const params = {}
      if (locale) {
        params.locale = locale
      }
      
      const response = await api.get(`/api/i18n/forms/${formType}/metadata`, { params })
      return response.data
    } catch (error) {
      console.error('Error fetching localized form metadata:', error)
      throw error
    }
  },

  // Bulk create/update translations
  async bulkUpdateTranslations(locale, translations) {
    try {
      const promises = Object.entries(translations).map(([key, content]) => 
        this.createOrUpdateTranslation({
          translation_key: key,
          locale,
          content,
          context: key.split('.')[0] // Use first part of key as context
        })
      )
      
      const results = await Promise.allSettled(promises)
      const failed = results.filter(result => result.status === 'rejected')
      
      if (failed.length > 0) {
        console.warn(`${failed.length} translations failed to update:`, failed)
      }
      
      return {
        total: promises.length,
        succeeded: results.length - failed.length,
        failed: failed.length,
        results
      }
    } catch (error) {
      console.error('Error in bulk translation update:', error)
      throw error
    }
  },

  // Get form field translations for a specific locale
  async getFormFieldTranslations(formType, locale) {
    try {
      const response = await this.getLocalizedFormMetadata(formType, locale)
      return response.data
    } catch (error) {
      console.error('Error fetching form field translations:', error)
      throw error
    }
  },

  // Update form field translations (admin only)
  async updateFormFieldTranslations(formType, locale, fieldTranslations) {
    try {
      // Get form metadata to get field definition IDs
      const formMetadata = await this.getLocalizedFormMetadata(formType, locale)
      
      if (!formMetadata.data?.fields) {
        throw new Error('No form fields found')
      }
      
      const promises = formMetadata.data.fields.map(field => {
        const translation = fieldTranslations[field.name]
        if (!translation) return null
        
        return this.createOrUpdateFieldTranslation({
          field_definition_id: field.field_id,
          locale,
          display_name: translation.label || field.display_name,
          description: translation.description || field.description,
          placeholder_text: translation.placeholder || field.placeholder_text
        })
      }).filter(Boolean)
      
      const results = await Promise.allSettled(promises)
      const failed = results.filter(result => result.status === 'rejected')
      
      if (failed.length > 0) {
        console.warn(`${failed.length} field translations failed to update:`, failed)
      }
      
      return {
        total: promises.length,
        succeeded: results.length - failed.length,
        failed: failed.length,
        results
      }
    } catch (error) {
      console.error('Error updating form field translations:', error)
      throw error
    }
  },

  // Helper method to get user's preferred locale from auth store (no API calls)
  getUserPreferredLocale(user) {
    return user?.preferred_locale || 'en-US'
  }
}

export default i18nAPI