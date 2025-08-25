import { locationsApi } from '@/api/locations'
import { computed, ref, watch } from 'vue'
import { useEntityStore } from '../stores/entity'
import { useFormConfigStore } from '../stores/formConfig'

/**
 * Composable for form configuration management and dynamic form rendering
 * @param {string} formType - The type of form (e.g., 'patient', 'appointment')
 */
export function useFormConfig(formType) {
  const formConfigStore = useFormConfigStore()
  const entityStore = useEntityStore()
  const nationalityOptions = ref([])
  const loadNationalitiesImmediate = async () => {
    if (nationalityOptions.value.length) return
    try {
      const data = await locationsApi.getNationalities()
      nationalityOptions.value = Array.isArray(data) ? data.map(n => ({ value: n.id, label: n.name, country_id: n.country_id, is_primary: n.is_primary })) : []
      nationalityOptions.value.sort((a,b)=>a.label.localeCompare(b.label))
    } catch (e) {
      console.error('[useFormConfig] Failed to load nationalities', e)
    }
  }
  
  // Local state
  const isInitialized = ref(false)
  const validationErrors = ref({})
  const formData = ref({})

  // Computed properties
  const fields = computed(() => formConfigStore.getFormFields(formType))
  const enabledFields = computed(() => formConfigStore.getEnabledFields(formType))
  const requiredFields = computed(() => formConfigStore.getRequiredFields(formType))
  const fieldsByCategory = computed(() => formConfigStore.getFieldsByCategory(formType))
  const coreFields = computed(() => formConfigStore.getCoreFields(formType))
  const configurableFields = computed(() => formConfigStore.getConfigurableFields(formType))
  
  const isLoading = computed(() => formConfigStore.isLoading)
  const error = computed(() => formConfigStore.error)

  // Form validation computed
  const isFormValid = computed(() => {
    return Object.keys(validationErrors.value).length === 0
  })

  const hasChanges = computed(() => formConfigStore.isDirty)

  // Methods
  const initialize = async () => {
    if (isInitialized.value) return

    try {
      await formConfigStore.loadFormFields(formType)
      
      // Load nationalities if this form contains nationality fields
      const allFields = formConfigStore.getFormFields(formType)
      const hasNationalityField = allFields.some(
        field => field.name === 'nationality_id'
      )
      
      if (hasNationalityField) {
        await loadNationalitiesImmediate()
      }
      
      isInitialized.value = true
    } catch (err) {
      console.error(`Failed to initialize form config for ${formType}:`, err)
      throw err
    }
  }

  const getFieldConfig = (fieldName) => {
    return fields.value.find(field => field.name === fieldName)
  }

  const isFieldEnabled = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.is_enabled : false
  }

  const isFieldRequired = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.is_enabled && field.is_required : false
  }

  const isFieldCore = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.is_core : false
  }

  const getFieldType = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.field_type : 'text'
  }

  const getFieldOptions = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.options || [] : []
  }

  const getFieldCategory = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.category : 'Other'
  }

  const getFieldDescription = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.description : ''
  }

  const getFieldDisplayName = (fieldName) => {
    const field = getFieldConfig(fieldName)
    return field ? field.display_name : fieldName
  }

  // Form data management
  const initializeFormData = (initialData = {}) => {
    
    // Clear existing data first
    Object.keys(formData.value).forEach(key => {
      delete formData.value[key]
    })
    
    // Populate with new data, maintaining reactivity
    enabledFields.value.forEach(field => {
      const fieldName = field.name
      const initialValue = initialData[fieldName]
      const defaultValue = getDefaultValue(field)
      const finalValue = initialValue !== undefined ? initialValue : defaultValue
      
      formData.value[fieldName] = finalValue
    })
    
    return formData.value
  }

  const getDefaultValue = (field) => {
    // Special handling for nationality field - use centralized entity default
    if (field.name === 'nationality_id') {
  const defaultNationalityId = entityStore.entityDefaultNationalityId
  return defaultNationalityId || ''
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

  const updateFormData = (fieldName, value) => {
    if (formData.value) {
      formData.value[fieldName] = value
      validateField(fieldName, value)
    }
  }

  // Validation
  const validateField = (fieldName, value) => {
    const field = getFieldConfig(fieldName)
    if (!field || !field.is_enabled) {
      delete validationErrors.value[fieldName]
      return true
    }

    const errors = []

    // Required field validation
    if (field.is_required && (value === null || value === undefined || value === '')) {
      errors.push(`${field.display_name} is required`)
    }

    // Cross-field dependency: city_id depends on country_id
    if (fieldName === 'city_id' || fieldName === 'city') {
      const countryVal = formData.value?.country_id || formData.value?.country || null
      if (!countryVal) {
        errors.push('Please select a country first')
      }
    }

    // Type-specific validation
    if (value && value !== '') {
      switch (field.field_type) {
        case 'email':
          if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
            errors.push('Please enter a valid email address')
          }
          break
        case 'tel':
          if (!/^[\+]?[\d\-\(\)\s]+$/.test(value)) {
            errors.push('Please enter a valid phone number')
          }
          break
        case 'url':
          try {
            new URL(value)
          } catch {
            errors.push('Please enter a valid URL')
          }
          break
        case 'number':
          if (isNaN(value)) {
            errors.push('Please enter a valid number')
          }
          break
        case 'date':
          if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) {
            errors.push('Please enter a valid date (YYYY-MM-DD)')
          }
          break
      }
    }

    if (errors.length > 0) {
      validationErrors.value[fieldName] = errors
      return false
    } else {
      delete validationErrors.value[fieldName]
      return true
    }
  }

  const validateForm = (data = formData.value) => {
    validationErrors.value = {}
    
    enabledFields.value.forEach(field => {
      const fieldName = field.name
      const value = data[fieldName]
      validateField(fieldName, value)
    })

    return isFormValid.value
  }

  const getFieldError = (fieldName) => {
    return validationErrors.value[fieldName] || []
  }

  const hasFieldError = (fieldName) => {
    return !!validationErrors.value[fieldName]
  }

  const clearFieldError = (fieldName) => {
    delete validationErrors.value[fieldName]
  }

  const clearAllErrors = () => {
    validationErrors.value = {}
  }


  // Dynamic form rendering helpers
  const renderField = (field, value, onUpdate) => {
    const baseProps = {
      id: field.name,
      name: field.name,
      required: field.is_required,
      disabled: !field.is_enabled,
      value: value,
      'onUpdate:modelValue': onUpdate,
      class: hasFieldError(field.name) ? 'border-red-500' : 'border-gray-300'
    }

    switch (field.field_type) {
      case 'text':
        return {
          component: 'input',
          props: { ...baseProps, type: 'text' }
        }
      case 'email':
        return {
          component: 'input',
          props: { ...baseProps, type: 'email' }
        }
      case 'tel':
        return {
          component: 'input',
          props: { ...baseProps, type: 'tel' }
        }
      case 'url':
        return {
          component: 'input',
          props: { ...baseProps, type: 'url' }
        }
      case 'number':
        return {
          component: 'input',
          props: { ...baseProps, type: 'number' }
        }
      case 'date':
        return {
          component: 'input',
          props: { ...baseProps, type: 'date' }
        }
      case 'textarea':
        return {
          component: 'textarea',
          props: { ...baseProps, rows: 3 }
        }
      case 'select':
        return {
          component: 'select',
          props: baseProps,
          options: field.options || []
        }
      case 'multiselect':
        return {
          component: 'select',
          props: { ...baseProps, multiple: true },
          options: field.options || []
        }
      case 'boolean':
        return {
          component: 'input',
          props: { ...baseProps, type: 'checkbox', checked: !!value }
        }
      default:
        return {
          component: 'input',
          props: { ...baseProps, type: 'text' }
        }
    }
  }

  const getFormFieldsForRendering = () => {
    return enabledFields.value.map(field => ({
      ...field,
      hasError: hasFieldError(field.name),
      errors: getFieldError(field.name),
      renderConfig: renderField(field, formData.value[field.name], (value) => updateFormData(field.name, value))
    }))
  }

  // Watch for form type changes
  watch(() => formType, async (newFormType) => {
    if (newFormType) {
      isInitialized.value = false
      await initialize()
    }
  }, { immediate: true })

  return {
    // State
    isInitialized,
    validationErrors,
    formData,

    // Computed
    fields,
    enabledFields,
    requiredFields,
    fieldsByCategory,
    coreFields,
    configurableFields,
    isLoading,
    error,
    isFormValid,
    hasChanges,

    // Methods
    initialize,
    getFieldConfig,
    isFieldEnabled,
    isFieldRequired,
    isFieldCore,
    getFieldType,
    getFieldOptions,
    getFieldCategory,
    getFieldDescription,
    getFieldDisplayName,

    // Form data
    initializeFormData,
    updateFormData,

    // Validation
    validateField,
    validateForm,
    getFieldError,
    hasFieldError,
    clearFieldError,
    clearAllErrors,

    // Rendering helpers
    renderField,
    getFormFieldsForRendering
  }
}