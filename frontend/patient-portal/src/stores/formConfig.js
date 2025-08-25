import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { formsApi } from '../api/forms'

export const useFormConfigStore = defineStore('formConfig', () => {
  // State
  const formTypes = ref([])
  const formConfigurations = ref({})
  const isLoading = ref(false)
  const error = ref(null)
  const isDirty = ref(false)

  // Getters
  const getFormConfig = computed(() => (formType) => {
    return formConfigurations.value[formType] || null
  })

  const getFormFields = computed(() => (formType) => {
    const config = formConfigurations.value[formType]
    return config?.fields || []
  })

  const getEnabledFields = computed(() => (formType) => {
    const fields = getFormFields.value(formType)
    return fields.filter(field => field.is_enabled)
  })

  const getRequiredFields = computed(() => (formType) => {
    const fields = getFormFields.value(formType)
    return fields.filter(field => field.is_enabled && field.is_required)
  })

  const getFieldsByCategory = computed(() => (formType) => {
    const fields = getFormFields.value(formType)
    const categories = {}
    
    fields.forEach(field => {
      const category = field.category || 'Other'
      if (!categories[category]) {
        categories[category] = []
      }
      categories[category].push(field)
    })

    // Sort fields within each category by sort_order
    Object.keys(categories).forEach(category => {
      categories[category].sort((a, b) => a.sort_order - b.sort_order)
    })

    return categories
  })

  const getCoreFields = computed(() => (formType) => {
    const fields = getFormFields.value(formType)
    return fields.filter(field => field.is_core)
  })

  const getConfigurableFields = computed(() => (formType) => {
    const fields = getFormFields.value(formType)
    return fields.filter(field => !field.is_core)
  })

  // Actions
  const loadFormTypes = async () => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await formsApi.getFormTypes()
      formTypes.value = response.form_types || response.data || []
      return response
    } catch (err) {
      error.value = err.message || 'Failed to load form types'
      console.error('Error loading form types:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const loadFormConfiguration = async (formType) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.getFormMetadata(formType)
      formConfigurations.value[formType] = response
      return response
    } catch (err) {
      error.value = err.message || `Failed to load configuration for ${formType}`
      console.error(`Error loading form configuration for ${formType}:`, err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const loadFormFields = async (formType, forceReload = false) => {
    // Check if fields are already loaded and not forcing reload
    const existingConfig = formConfigurations.value[formType]
    if (!forceReload && existingConfig && existingConfig.fields && existingConfig.fields.length > 0) {
      return existingConfig
    }

    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.getFormFields(formType)
      
      // Update the fields in the existing configuration or create new one
      if (!formConfigurations.value[formType]) {
        formConfigurations.value[formType] = {
          form_type: formType,
          fields: []
        }
      }
      
      formConfigurations.value[formType].fields = response.fields || response.data || []
      return response
    } catch (err) {
      error.value = err.message || `Failed to load fields for ${formType}`
      console.error(`Error loading form fields for ${formType}:`, err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateField = async (formType, fieldId, fieldConfig) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.updateFormField(formType, fieldId, fieldConfig)
      
      // Update the field in local state
      const config = formConfigurations.value[formType]
      if (config && config.fields) {
        const fieldIndex = config.fields.findIndex(f => f.field_id === fieldId)
        if (fieldIndex !== -1) {
          config.fields[fieldIndex] = { ...config.fields[fieldIndex], ...fieldConfig }
        }
      }
      
      isDirty.value = true
      return response
    } catch (err) {
      error.value = err.message || 'Failed to update field'
      console.error('Error updating field:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateMultipleFields = async (formType, fieldsConfig) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.updateFormFields(formType, fieldsConfig)
      
      // Update fields in local state
      const config = formConfigurations.value[formType]
      if (config && config.fields) {
        fieldsConfig.forEach(fieldUpdate => {
          const fieldIndex = config.fields.findIndex(f => f.field_id === fieldUpdate.field_id)
          if (fieldIndex !== -1) {
            config.fields[fieldIndex] = { ...config.fields[fieldIndex], ...fieldUpdate }
          }
        })
      }
      
      isDirty.value = true
      return response
    } catch (err) {
      error.value = err.message || 'Failed to update fields'
      console.error('Error updating multiple fields:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateFieldOrder = async (formType, fieldOrders) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.updateFieldOrder(formType, fieldOrders)
      
      // Update sort orders in local state
      const config = formConfigurations.value[formType]
      if (config && config.fields) {
        fieldOrders.forEach(orderUpdate => {
          const fieldIndex = config.fields.findIndex(f => f.field_id === orderUpdate.field_id)
          if (fieldIndex !== -1) {
            config.fields[fieldIndex].sort_order = orderUpdate.sort_order
          }
        })
        
        // Re-sort the fields by sort_order
        config.fields.sort((a, b) => a.sort_order - b.sort_order)
      }
      
      isDirty.value = true
      return response
    } catch (err) {
      error.value = err.message || 'Failed to update field order'
      console.error('Error updating field order:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const resetFormToDefaults = async (formType) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await formsApi.resetFormToDefaults(formType)
      
      // Reload the form configuration after reset
      await loadFormConfiguration(formType)
      
      isDirty.value = false
      return response
    } catch (err) {
      error.value = err.message || `Failed to reset ${formType} form`
      console.error(`Error resetting form ${formType}:`, err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Clear cached form fields (useful for forcing reload)
  const clearFormFieldsCache = (formType) => {
    if (formType) {
      if (formConfigurations.value[formType]) {
        formConfigurations.value[formType].fields = []
      }
    } else {
      Object.keys(formConfigurations.value).forEach(type => {
        formConfigurations.value[type].fields = []
      })
    }
  }

  // Force reload form fields from API
  const reloadFormFields = async (formType) => {
    clearFormFieldsCache(formType)
    return await loadFormFields(formType, true)
  }

  const toggleFieldEnabled = async (formType, fieldId) => {
    const config = formConfigurations.value[formType]
    if (!config || !config.fields) return

    const field = config.fields.find(f => f.field_id === fieldId)
    if (!field || field.is_core) return // Can't toggle core fields

    const newEnabledState = !field.is_enabled
    await updateField(formType, fieldId, { is_enabled: newEnabledState })
  }

  const toggleFieldRequired = async (formType, fieldId) => {
    const config = formConfigurations.value[formType]
    if (!config || !config.fields) return

    const field = config.fields.find(f => f.field_id === fieldId)
    if (!field || !field.is_enabled) return // Can't make disabled field required

    const newRequiredState = !field.is_required
    await updateField(formType, fieldId, { is_required: newRequiredState })
  }

  const clearError = () => {
    error.value = null
  }

  const markClean = () => {
    isDirty.value = false
  }

  // Initialize store
  const initialize = async () => {
    try {
      await loadFormTypes()
    } catch (err) {
      console.error('Failed to initialize form config store:', err)
    }
  }

  return {
    // State
    formTypes,
    formConfigurations,
    isLoading,
    error,
    isDirty,

    // Getters
    getFormConfig,
    getFormFields,
    getEnabledFields,
    getRequiredFields,
    getFieldsByCategory,
    getCoreFields,
    getConfigurableFields,

    // Actions
    loadFormTypes,
    loadFormConfiguration,
    loadFormFields,
    updateField,
    updateMultipleFields,
    updateFieldOrder,
    resetFormToDefaults,
    toggleFieldEnabled,
    toggleFieldRequired,
    clearError,
    markClean,
    initialize,
    
    // Cache management
    clearFormFieldsCache,
    reloadFormFields
  }
})