import { ref, computed } from 'vue'

/**
 * Composable for form validation with custom rules and error handling
 * Provides standardized validation patterns for healthcare forms
 */
export function useValidation() {
  // Validation state
  const validationErrors = ref({})
  const validationRules = ref({})

  // Computed properties
  const hasErrors = computed(() => {
    return Object.keys(validationErrors.value).length > 0
  })

  const isValid = computed(() => {
    return !hasErrors.value
  })

  const errorCount = computed(() => {
    return Object.keys(validationErrors.value).length
  })

  // Built-in validation rules
  const builtInRules = {
    required: (value, message = 'This field is required') => {
      if (value === null || value === undefined || value === '' || 
          (Array.isArray(value) && value.length === 0)) {
        return message
      }
      return null
    },

    email: (value, message = 'Please enter a valid email address') => {
      if (!value) return null
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return emailRegex.test(value) ? null : message
    },

    phone: (value, message = 'Please enter a valid phone number') => {
      if (!value) return null
      const phoneRegex = /^[\+]?[\d\-\(\)\s]+$/
      return phoneRegex.test(value) ? null : message
    },

    url: (value, message = 'Please enter a valid URL') => {
      if (!value) return null
      try {
        new URL(value)
        return null
      } catch {
        return message
      }
    },

    number: (value, message = 'Please enter a valid number') => {
      if (!value) return null
      return isNaN(value) ? message : null
    },

    date: (value, message = 'Please enter a valid date (YYYY-MM-DD)') => {
      if (!value) return null
      const dateRegex = /^\d{4}-\d{2}-\d{2}$/
      if (!dateRegex.test(value)) return message
      
      const date = new Date(value)
      return isNaN(date.getTime()) ? message : null
    },

    minLength: (minLen) => (value, message = `Minimum length is ${minLen} characters`) => {
      if (!value) return null
      return value.length >= minLen ? null : message
    },

    maxLength: (maxLen) => (value, message = `Maximum length is ${maxLen} characters`) => {
      if (!value) return null
      return value.length <= maxLen ? null : message
    },

    pattern: (regex, message = 'Please enter a valid format') => (value) => {
      if (!value) return null
      const pattern = new RegExp(regex)
      return pattern.test(value) ? null : message
    },

    // Healthcare-specific validations
    nationalId: (value, message = 'Please enter a valid national ID') => {
      if (!value) return null
      // Basic validation - could be enhanced per country
      return /^[A-Za-z0-9\-]+$/.test(value) ? null : message
    },

    postalCode: (value, message = 'Please enter a valid postal code') => {
      if (!value) return null
      // Basic validation - could be enhanced per country
      return /^[A-Za-z0-9\s\-]+$/.test(value) ? null : message
    },

    bloodType: (value, message = 'Please select a valid blood type') => {
      if (!value) return null
      const validTypes = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']
      return validTypes.includes(value) ? null : message
    },

    age: (value, message = 'Please enter a valid age') => {
      if (!value) return null
      const age = parseInt(value)
      return (age >= 0 && age <= 150) ? null : message
    },

    // Conditional validations
    requiredIf: (condition) => (value, message = 'This field is required') => {
      if (condition && builtInRules.required(value)) {
        return message
      }
      return null
    },

    // Cross-field validations
    matchField: (otherValue, message = 'Fields do not match') => (value) => {
      return value === otherValue ? null : message
    }
  }

  // Methods
  const setRules = (fieldName, rules) => {
    validationRules.value[fieldName] = Array.isArray(rules) ? rules : [rules]
  }

  const addRule = (fieldName, rule) => {
    if (!validationRules.value[fieldName]) {
      validationRules.value[fieldName] = []
    }
    validationRules.value[fieldName].push(rule)
  }

  const removeRule = (fieldName, ruleIndex = null) => {
    if (ruleIndex === null) {
      delete validationRules.value[fieldName]
    } else if (validationRules.value[fieldName]) {
      validationRules.value[fieldName].splice(ruleIndex, 1)
    }
  }

  const validateField = (fieldName, value, customRules = null) => {
    const rules = customRules || validationRules.value[fieldName] || []
    const errors = []

    for (const rule of rules) {
      let error = null

      if (typeof rule === 'function') {
        // Custom function rule
        error = rule(value)
      } else if (typeof rule === 'object' && rule.rule) {
        // Rule object with parameters
        if (typeof rule.rule === 'string' && builtInRules[rule.rule]) {
          error = builtInRules[rule.rule](value, rule.message)
        } else if (typeof rule.rule === 'function') {
          error = rule.rule(value, rule.message)
        }
      } else if (typeof rule === 'string' && builtInRules[rule]) {
        // Built-in rule name
        error = builtInRules[rule](value)
      }

      if (error) {
        errors.push(error)
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

  const validateFields = (fieldsData, fieldConfigs = null) => {
    const results = {}

    Object.entries(fieldsData).forEach(([fieldName, value]) => {
      let rules = validationRules.value[fieldName]

      // Use field config rules if provided
      if (fieldConfigs && fieldConfigs[fieldName]) {
        const config = fieldConfigs[fieldName]
        rules = buildRulesFromConfig(config)
      }

      results[fieldName] = validateField(fieldName, value, rules)
    })

    return results
  }

  const validateAll = (formData, fieldConfigs = null) => {
    const results = validateFields(formData, fieldConfigs)
    return Object.values(results).every(isValid => isValid)
  }

  const buildRulesFromConfig = (fieldConfig) => {
    const rules = []

    // Required validation
    if (fieldConfig.is_required) {
      rules.push(builtInRules.required)
    }

    // Type-specific validations
    switch (fieldConfig.field_type) {
      case 'email':
        rules.push(builtInRules.email)
        break
      case 'tel':
      case 'phone':
        rules.push(builtInRules.phone)
        break
      case 'url':
        rules.push(builtInRules.url)
        break
      case 'number':
        rules.push(builtInRules.number)
        break
      case 'date':
        rules.push(builtInRules.date)
        break
    }

    // Validation rules from config
    if (fieldConfig.validation_rules) {
      const vRules = fieldConfig.validation_rules

      if (vRules.min_length) {
        rules.push(builtInRules.minLength(vRules.min_length))
      }

      if (vRules.max_length) {
        rules.push(builtInRules.maxLength(vRules.max_length))
      }

      if (vRules.pattern) {
        rules.push(builtInRules.pattern(vRules.pattern))
      }
    }

    return rules
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

  const setFieldError = (fieldName, errors) => {
    if (Array.isArray(errors) && errors.length > 0) {
      validationErrors.value[fieldName] = errors
    } else if (typeof errors === 'string') {
      validationErrors.value[fieldName] = [errors]
    } else {
      delete validationErrors.value[fieldName]
    }
  }

  // Healthcare-specific validation helpers
  const validateLocation = (country, state, city, stateOptions, cityOptions, isLoading = {}) => {
    const errors = {}

    // Skip validation if location data is still loading
    if (isLoading.states || isLoading.cities) {
      return errors
    }

    if (state && country) {
      const allowedStates = stateOptions.map(s => s.id)
      // State should now be a number from v-model.number
      if (allowedStates.length > 0 && !allowedStates.includes(state)) {
        errors.state = ['Selected state/province is not valid for the chosen country']
      }
    }

    if (city && country) {
      const allowedCities = cityOptions.map(c => c.id)
      // City should now be a number from v-model.number
      if (allowedCities.length > 0 && !allowedCities.includes(city)) {
        errors.city = ['Selected city is not valid for the chosen country']
      }
    }

    return errors
  }

  const validatePatientData = (patientData) => {
    const errors = {}

    // Age validation based on date of birth
    if (patientData.date_of_birth) {
      const birthDate = new Date(patientData.date_of_birth)
      const today = new Date()
      const age = Math.floor((today - birthDate) / (365.25 * 24 * 60 * 60 * 1000))
      
      if (age < 0 || age > 150) {
        errors.date_of_birth = ['Please enter a valid birth date']
      }
    }

    // Emergency contact validation
    if (patientData.emergency_contact_name && !patientData.emergency_contact_phone) {
      errors.emergency_contact_phone = ['Emergency contact phone is required when contact name is provided']
    }

    return errors
  }

  return {
    // State
    validationErrors: computed(() => validationErrors.value),
    validationRules: computed(() => validationRules.value),
    
    // Computed
    hasErrors,
    isValid,
    errorCount,
    
    // Rule management
    setRules,
    addRule,
    removeRule,
    builtInRules,
    
    // Validation methods
    validateField,
    validateFields,
    validateAll,
    buildRulesFromConfig,
    
    // Error management
    getFieldError,
    hasFieldError,
    clearFieldError,
    clearAllErrors,
    setFieldError,
    
    // Healthcare-specific
    validateLocation,
    validatePatientData
  }
}