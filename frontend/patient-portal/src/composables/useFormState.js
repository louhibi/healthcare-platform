import { ref, computed, nextTick } from 'vue'
import { useToast } from 'vue-toastification'
import { useValidation } from './useValidation'

/**
 * Composable for managing form state, validation, and submission
 * Handles form data, validation errors, and submission states
 */
export function useFormState(apiService, options = {}) {
  const toast = useToast()
  const validation = useValidation()
  
  // Options with defaults
  const {
    createMethod = 'createPatient',
    updateMethod = 'updatePatient',
    successMessage = {
      create: 'Record created successfully!',
      update: 'Record updated successfully!'
    },
    errorMessage = {
      create: 'Failed to create record',
      update: 'Failed to update record',
      general: 'An error occurred'
    }
  } = options

  // Form state
  const formData = ref({})
  const isSubmitting = ref(false)
  const isDirty = ref(false)

  // Computed properties (delegated to validation composable)
  const hasErrors = computed(() => validation.hasErrors.value)
  const isFormValid = computed(() => validation.isValid.value)
  const validationErrors = computed(() => validation.validationErrors.value)

  // Methods
  const initializeFormData = (initialData = {}, enabledFields = [], entityStore = null) => {
    
    // Clear existing data first
    Object.keys(formData.value).forEach(key => {
      delete formData.value[key]
    })
    
    // Populate with new data, maintaining reactivity
    enabledFields.forEach(field => {
      const fieldName = field.name
      const initialValue = initialData[fieldName]
      const defaultValue = getDefaultValue(field, entityStore)
      const finalValue = initialValue !== undefined ? initialValue : defaultValue
      
      formData.value[fieldName] = finalValue
    })
    
    // Reset dirty and error states
    isDirty.value = false
    validation.clearAllErrors()
    
    return formData.value
  }

  const getDefaultValue = (field, entityStore = null) => {
    // Special handling for nationality field - use entity default if available
    if (field.name === 'nationality_id' && entityStore?.entityDefaultNationalityId) {
      return entityStore.entityDefaultNationalityId
    }
    
    switch (field.field_type) {
      case 'boolean':
        return false
      case 'number':
        return 0
      case 'select':
        return field.options && field.options.length > 0 ? field.options[0].value : ''
      case 'multiselect':
        return []
      case 'date':
        return ''
      case 'email':
      case 'tel':
      case 'url':
      case 'text':
      case 'textarea':
      default:
        return ''
    }
  }

  const updateField = (fieldName, value) => {
    if (formData.value[fieldName] !== value) {
      formData.value[fieldName] = value
      isDirty.value = true
    }
  }

  const validateField = (fieldName, value, fieldConfig) => {
    if (!fieldConfig || !fieldConfig.is_enabled) {
      validation.clearFieldError(fieldName)
      return true
    }

    // Build validation rules from field configuration
    const rules = validation.buildRulesFromConfig(fieldConfig)
    return validation.validateField(fieldName, value, rules)
  }

  const validateForm = (fieldsConfig = []) => {
    validation.clearAllErrors()
    
    fieldsConfig.forEach(field => {
      const fieldName = field.name
      const value = formData.value[fieldName]
      validateField(fieldName, value, field)
    })

    return isFormValid.value
  }

  // Delegate error management to validation composable
  const getFieldError = validation.getFieldError
  const hasFieldError = validation.hasFieldError
  const clearFieldError = validation.clearFieldError
  const clearAllErrors = validation.clearAllErrors

  const resetForm = (enabledFields = []) => {
    // Reset form data to default values
    enabledFields.forEach(field => {
      formData.value[field.name] = getDefaultValue(field)
    })
    
    // Clear validation state
    clearAllErrors()
    isDirty.value = false
  }

  const submitForm = async (isEdit = false, recordId = null, additionalValidation = null) => {
    if (isSubmitting.value) return null

    // Run additional validation if provided
    if (additionalValidation && !additionalValidation()) {
      return null
    }

    isSubmitting.value = true
    
    try {
      let result
      
      if (isEdit && recordId) {
        // Update existing record
        result = await apiService[updateMethod](recordId, formData.value)
        toast.success(successMessage.update)
      } else {
        // Create new record
        result = await apiService[createMethod](formData.value)
        toast.success(successMessage.create)
      }
      
      // Mark as clean after successful submission
      isDirty.value = false
      
      return result
    } catch (err) {
      const message = isEdit ? errorMessage.update : errorMessage.create
      toast.error(err.message || message)
      console.error('Form submission error:', err)
      throw err
    } finally {
      isSubmitting.value = false
    }
  }

  // Bulk field updates
  const updateMultipleFields = async (updates = {}) => {
    Object.entries(updates).forEach(([fieldName, value]) => {
      updateField(fieldName, value)
    })
    await nextTick()
  }

  return {
    // State
    formData,
    validationErrors,
    isSubmitting,
    isDirty,
    
    // Computed
    hasErrors,
    isFormValid,
    
    // Methods
    initializeFormData,
    updateField,
    updateMultipleFields,
    validateField,
    validateForm,
    getFieldError,
    hasFieldError,
    clearFieldError,
    clearAllErrors,
    resetForm,
    submitForm
  }
}