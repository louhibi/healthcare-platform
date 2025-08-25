import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useValidation } from '../useValidation'

describe('useValidation', () => {
  let validation

  beforeEach(() => {
    validation = useValidation()
  })

  describe('initialization', () => {
    it('should initialize with empty errors', () => {
      expect(validation.validationErrors.value).toEqual({})
      expect(validation.hasErrors.value).toBe(false)
      expect(validation.isValid.value).toBe(true)
      expect(validation.errorCount.value).toBe(0)
    })

    it('should provide built-in validation rules', () => {
      expect(validation.builtInRules).toHaveProperty('required')
      expect(validation.builtInRules).toHaveProperty('email')
      expect(validation.builtInRules).toHaveProperty('phone')
      expect(validation.builtInRules).toHaveProperty('bloodType')
    })
  })

  describe('built-in rules', () => {
    describe('required rule', () => {
      it('should validate required fields', () => {
        const { required } = validation.builtInRules
        
        expect(required('')).toBe('This field is required')
        expect(required(null)).toBe('This field is required')
        expect(required(undefined)).toBe('This field is required')
        expect(required([])).toBe('This field is required')
        expect(required('valid')).toBe(null)
        expect(required('0')).toBe(null)
        expect(required(['item'])).toBe(null)
      })

      it('should use custom error message', () => {
        const { required } = validation.builtInRules
        const customMessage = 'Custom required message'
        
        expect(required('', customMessage)).toBe(customMessage)
      })
    })

    describe('email rule', () => {
      it('should validate email addresses', () => {
        const { email } = validation.builtInRules
        
        expect(email('test@example.com')).toBe(null)
        expect(email('user.name+tag@domain.co.uk')).toBe(null)
        expect(email('invalid-email')).toBe('Please enter a valid email address')
        expect(email('test@')).toBe('Please enter a valid email address')
        expect(email('@domain.com')).toBe('Please enter a valid email address')
        expect(email('')).toBe(null) // Empty should pass (handled by required)
      })
    })

    describe('phone rule', () => {
      it('should validate phone numbers', () => {
        const { phone } = validation.builtInRules
        
        expect(phone('+1-555-123-4567')).toBe(null)
        expect(phone('(555) 123-4567')).toBe(null)
        expect(phone('555-123-4567')).toBe(null)
        expect(phone('5551234567')).toBe(null)
        expect(phone('invalid-phone')).toBe('Please enter a valid phone number')
        expect(phone('')).toBe(null) // Empty should pass
      })
    })

    describe('bloodType rule', () => {
      it('should validate blood types', () => {
        const { bloodType } = validation.builtInRules
        
        expect(bloodType('A+')).toBe(null)
        expect(bloodType('A-')).toBe(null)
        expect(bloodType('B+')).toBe(null)
        expect(bloodType('B-')).toBe(null)
        expect(bloodType('AB+')).toBe(null)
        expect(bloodType('AB-')).toBe(null)
        expect(bloodType('O+')).toBe(null)
        expect(bloodType('O-')).toBe(null)
        expect(bloodType('X+')).toBe('Please select a valid blood type')
        expect(bloodType('')).toBe(null) // Empty should pass
      })
    })

    describe('date rule', () => {
      it('should validate dates', () => {
        const { date } = validation.builtInRules
        
        expect(date('2024-01-15')).toBe(null)
        expect(date('2024-12-31')).toBe(null)
        expect(date('invalid-date')).toBe('Please enter a valid date (YYYY-MM-DD)')
        expect(date('2024-13-01')).toBe('Please enter a valid date (YYYY-MM-DD)')
        expect(date('')).toBe(null) // Empty should pass
      })
    })

    describe('number rule', () => {
      it('should validate numbers', () => {
        const { number } = validation.builtInRules
        
        expect(number('123')).toBe(null)
        expect(number('123.45')).toBe(null)
        expect(number('-123')).toBe(null)
        expect(number('0')).toBe(null)
        expect(number('not-a-number')).toBe('Please enter a valid number')
        expect(number('')).toBe(null) // Empty should pass
      })
    })
  })

  describe('rule management', () => {
    it('should set rules for a field', () => {
      const rules = ['required', 'email']
      validation.setRules('email_field', rules)
      
      expect(validation.validationRules.value.email_field).toEqual(rules)
    })

    it('should add rule to existing field', () => {
      validation.setRules('test_field', ['required'])
      validation.addRule('test_field', 'email')
      
      expect(validation.validationRules.value.test_field).toEqual(['required', 'email'])
    })

    it('should remove all rules for field', () => {
      validation.setRules('test_field', ['required', 'email'])
      validation.removeRule('test_field')
      
      expect(validation.validationRules.value.test_field).toBeUndefined()
    })
  })

  describe('field validation', () => {
    it('should validate field with single rule', () => {
      const result = validation.validateField('test_field', '', ['required'])
      
      expect(result).toBe(false)
      expect(validation.validationErrors.value.test_field).toEqual(['This field is required'])
    })

    it('should validate field with multiple rules', () => {
      const rules = ['required', 'email']
      const result = validation.validateField('email_field', 'invalid-email', rules)
      
      expect(result).toBe(false)
      expect(validation.validationErrors.value.email_field).toEqual([
        'Please enter a valid email address'
      ])
    })

    it('should pass validation with valid data', () => {
      const rules = ['required', 'email']
      const result = validation.validateField('email_field', 'test@example.com', rules)
      
      expect(result).toBe(true)
      expect(validation.validationErrors.value.email_field).toBeUndefined()
    })

    it('should validate with function rules', () => {
      const customRule = vi.fn().mockReturnValue('Custom error')
      const result = validation.validateField('test_field', 'value', [customRule])
      
      expect(result).toBe(false)
      expect(validation.validationErrors.value.test_field).toEqual(['Custom error'])
      expect(customRule).toHaveBeenCalledWith('value')
    })

    it('should validate with rule objects', () => {
      const ruleObject = {
        rule: 'required',
        message: 'Custom required message'
      }
      const result = validation.validateField('test_field', '', [ruleObject])
      
      expect(result).toBe(false)
      expect(validation.validationErrors.value.test_field).toEqual(['Custom required message'])
    })
  })

  describe('buildRulesFromConfig', () => {
    it('should build rules from field configuration', () => {
      const fieldConfig = {
        is_required: true,
        field_type: 'email',
        validation_rules: {
          min_length: 5,
          max_length: 50
        }
      }
      
      const rules = validation.buildRulesFromConfig(fieldConfig)
      
      expect(rules).toHaveLength(4) // required, email, minLength, maxLength
      expect(rules[0]).toBe(validation.builtInRules.required)
      expect(rules[1]).toBe(validation.builtInRules.email)
    })

    it('should handle field config without validation rules', () => {
      const fieldConfig = {
        is_required: false,
        field_type: 'text'
      }
      
      const rules = validation.buildRulesFromConfig(fieldConfig)
      
      expect(rules).toHaveLength(0)
    })
  })

  describe('error management', () => {
    beforeEach(() => {
      validation.validateField('test_field', '', ['required'])
      validation.validateField('email_field', 'invalid', ['email'])
    })

    it('should get field errors', () => {
      const errors = validation.getFieldError('test_field')
      
      expect(errors).toEqual(['This field is required'])
    })

    it('should check if field has error', () => {
      expect(validation.hasFieldError('test_field')).toBe(true)
      expect(validation.hasFieldError('nonexistent_field')).toBe(false)
    })

    it('should clear field error', () => {
      validation.clearFieldError('test_field')
      
      expect(validation.hasFieldError('test_field')).toBe(false)
      expect(validation.validationErrors.value.test_field).toBeUndefined()
    })

    it('should clear all errors', () => {
      validation.clearAllErrors()
      
      expect(validation.validationErrors.value).toEqual({})
      expect(validation.hasErrors.value).toBe(false)
    })

    it('should set field errors', () => {
      const errors = ['Error 1', 'Error 2']
      validation.setFieldError('new_field', errors)
      
      expect(validation.getFieldError('new_field')).toEqual(errors)
    })
  })

  describe('healthcare-specific validations', () => {
    describe('validateLocation', () => {
      it('should validate location cascade', () => {
        const stateOptions = [{ code: 'CA' }, { code: 'NY' }]
        const cityOptions = [{ id: 'toronto' }, { id: 'vancouver' }]
        
        const errors = validation.validateLocation('Canada', 'TX', 'toronto', stateOptions, cityOptions)
        
        expect(errors.state).toEqual(['Selected state/province is not valid for the chosen country'])
      })

      it('should validate city selection', () => {
        const stateOptions = [{ code: 'CA' }]
        const cityOptions = [{ id: 'toronto' }]
        
        const errors = validation.validateLocation('Canada', 'CA', 'invalid_city', stateOptions, cityOptions)
        
        expect(errors.city).toEqual(['Selected city is not valid for the chosen country'])
      })

      it('should pass with valid location data', () => {
        const stateOptions = [{ code: 'CA' }]
        const cityOptions = [{ id: 'toronto' }]
        
        const errors = validation.validateLocation('Canada', 'CA', 'toronto', stateOptions, cityOptions)
        
        expect(errors).toEqual({})
      })
    })

    describe('validatePatientData', () => {
      it('should validate birth date age', () => {
        const futureDate = new Date()
        futureDate.setFullYear(futureDate.getFullYear() + 1)
        
        const patientData = {
          date_of_birth: futureDate.toISOString().split('T')[0]
        }
        
        const errors = validation.validatePatientData(patientData)
        
        expect(errors.date_of_birth).toEqual(['Please enter a valid birth date'])
      })

      it('should validate emergency contact consistency', () => {
        const patientData = {
          emergency_contact_name: 'John Doe'
          // missing emergency_contact_phone
        }
        
        const errors = validation.validatePatientData(patientData)
        
        expect(errors.emergency_contact_phone).toEqual(['Emergency contact phone is required when contact name is provided'])
      })

      it('should pass with valid patient data', () => {
        const birthDate = new Date()
        birthDate.setFullYear(birthDate.getFullYear() - 25)
        
        const patientData = {
          date_of_birth: birthDate.toISOString().split('T')[0],
          emergency_contact_name: 'John Doe',
          emergency_contact_phone: '555-123-4567'
        }
        
        const errors = validation.validatePatientData(patientData)
        
        expect(errors).toEqual({})
      })
    })
  })

  describe('reactive properties', () => {
    it('should update computed properties when errors change', () => {
      expect(validation.hasErrors.value).toBe(false)
      expect(validation.isValid.value).toBe(true)
      expect(validation.errorCount.value).toBe(0)
      
      validation.validateField('test_field', '', ['required'])
      
      expect(validation.hasErrors.value).toBe(true)
      expect(validation.isValid.value).toBe(false)
      expect(validation.errorCount.value).toBe(1)
    })
  })
})