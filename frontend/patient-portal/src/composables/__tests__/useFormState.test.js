import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useFormState } from '../useFormState'

// Mock useToast
vi.mock('vue-toastification', () => ({
  useToast: () => ({
    success: vi.fn(),
    error: vi.fn()
  })
}))

// Mock useValidation
vi.mock('../useValidation', () => ({
  useValidation: () => ({
    hasErrors: { value: false },
    isValid: { value: true },
    validationErrors: { value: {} },
    clearAllErrors: vi.fn(),
    clearFieldError: vi.fn(),
    buildRulesFromConfig: vi.fn().mockReturnValue([]),
    validateField: vi.fn().mockReturnValue(true)
  })
}))

describe('useFormState', () => {
  let mockApiService
  let formState

  beforeEach(() => {
    mockApiService = {
      createPatient: vi.fn(),
      updatePatient: vi.fn()
    }
    
    formState = useFormState(mockApiService)
  })

  describe('initialization', () => {
    it('should initialize with empty form data', () => {
      expect(formState.formData.value).toEqual({})
      expect(formState.isSubmitting.value).toBe(false)
      expect(formState.isDirty.value).toBe(false)
    })

    it('should accept custom options', () => {
      const options = {
        createMethod: 'createCustom',
        updateMethod: 'updateCustom',
        successMessage: {
          create: 'Custom create success',
          update: 'Custom update success'
        }
      }
      
      const customFormState = useFormState(mockApiService, options)
      expect(customFormState).toBeDefined()
    })
  })

  describe('form data initialization', () => {
    it('should initialize form data with field defaults', () => {
      const enabledFields = [
        { name: 'first_name', field_type: 'text' },
        { name: 'age', field_type: 'number' },
        { name: 'is_active', field_type: 'boolean' },
        { name: 'notes', field_type: 'textarea' }
      ]
      
      formState.initializeFormData({}, enabledFields)
      
      expect(formState.formData.value).toEqual({
        first_name: '',
        age: 0,
        is_active: false,
        notes: ''
      })
    })

    it('should initialize with initial data', () => {
      const initialData = {
        first_name: 'John',
        last_name: 'Doe',
        age: 25
      }
      
      const enabledFields = [
        { name: 'first_name', field_type: 'text' },
        { name: 'last_name', field_type: 'text' },
        { name: 'age', field_type: 'number' }
      ]
      
      formState.initializeFormData(initialData, enabledFields)
      
      expect(formState.formData.value).toEqual(initialData)
    })

    it('should clear existing form data when reinitializing', () => {
      // Set initial data
      formState.updateField('old_field', 'old_value')
      
      const enabledFields = [
        { name: 'new_field', field_type: 'text' }
      ]
      
      formState.initializeFormData({}, enabledFields)
      
      expect(formState.formData.value).toEqual({
        new_field: ''
      })
      expect(formState.formData.value.old_field).toBeUndefined()
    })
  })

  describe('default value generation', () => {
    it('should provide correct defaults for different field types', () => {
      const testCases = [
        { field: { field_type: 'text' }, expected: '' },
        { field: { field_type: 'email' }, expected: '' },
        { field: { field_type: 'number' }, expected: 0 },
        { field: { field_type: 'boolean' }, expected: false },
        { field: { field_type: 'textarea' }, expected: '' },
        { field: { field_type: 'date' }, expected: '' },
        { field: { field_type: 'multiselect' }, expected: [] },
        { field: { field_type: 'select', options: [{ value: 'option1' }] }, expected: 'option1' }
      ]
      
      testCases.forEach(({ field, expected }) => {
        formState.initializeFormData({}, [field])
        expect(formState.formData.value[field.name || 'test']).toEqual(expected)
      })
    })
  })

  describe('field updates', () => {
    it('should update field values', () => {
      formState.updateField('test_field', 'new_value')
      
      expect(formState.formData.value.test_field).toBe('new_value')
      expect(formState.isDirty.value).toBe(true)
    })

    it('should not mark as dirty if value is the same', () => {
      formState.formData.value.test_field = 'existing_value'
      formState.isDirty.value = false
      
      formState.updateField('test_field', 'existing_value')
      
      expect(formState.isDirty.value).toBe(false)
    })

    it('should update multiple fields', async () => {
      const updates = {
        field1: 'value1',
        field2: 'value2',
        field3: 'value3'
      }
      
      await formState.updateMultipleFields(updates)
      
      expect(formState.formData.value.field1).toBe('value1')
      expect(formState.formData.value.field2).toBe('value2')
      expect(formState.formData.value.field3).toBe('value3')
      expect(formState.isDirty.value).toBe(true)
    })
  })

  describe('form validation', () => {
    it('should validate form with field configs', () => {
      const fieldsConfig = [
        { name: 'field1', is_enabled: true },
        { name: 'field2', is_enabled: true }
      ]
      
      const result = formState.validateForm(fieldsConfig)
      
      expect(result).toBe(true)
    })

    it('should validate individual field', () => {
      const fieldConfig = { 
        name: 'test_field', 
        is_enabled: true,
        is_required: true
      }
      
      const result = formState.validateField('test_field', 'test_value', fieldConfig)
      
      expect(result).toBe(true)
    })

    it('should skip validation for disabled fields', () => {
      const fieldConfig = { 
        name: 'test_field', 
        is_enabled: false
      }
      
      const result = formState.validateField('test_field', '', fieldConfig)
      
      expect(result).toBe(true)
    })
  })

  describe('form reset', () => {
    it('should reset form to default values', () => {
      // Set some data
      formState.updateField('field1', 'value1')
      formState.updateField('field2', 'value2')
      
      const enabledFields = [
        { name: 'field1', field_type: 'text' },
        { name: 'field2', field_type: 'number' }
      ]
      
      formState.resetForm(enabledFields)
      
      expect(formState.formData.value).toEqual({
        field1: '',
        field2: 0
      })
      expect(formState.isDirty.value).toBe(false)
    })
  })

  describe('form submission', () => {
    beforeEach(() => {
      formState.formData.value = {
        first_name: 'John',
        last_name: 'Doe'
      }
    })

    it('should create new record when not in edit mode', async () => {
      const mockResponse = { id: 1, first_name: 'John', last_name: 'Doe' }
      mockApiService.createPatient.mockResolvedValue(mockResponse)
      
      const result = await formState.submitForm(false)
      
      expect(mockApiService.createPatient).toHaveBeenCalledWith({
        first_name: 'John',
        last_name: 'Doe'
      })
      expect(result).toEqual(mockResponse)
      expect(formState.isDirty.value).toBe(false)
      expect(formState.isSubmitting.value).toBe(false)
    })

    it('should update existing record when in edit mode', async () => {
      const mockResponse = { id: 1, first_name: 'John', last_name: 'Doe' }
      mockApiService.updatePatient.mockResolvedValue(mockResponse)
      
      const result = await formState.submitForm(true, 1)
      
      expect(mockApiService.updatePatient).toHaveBeenCalledWith(1, {
        first_name: 'John',
        last_name: 'Doe'
      })
      expect(result).toEqual(mockResponse)
      expect(formState.isDirty.value).toBe(false)
    })

    it('should handle submission errors', async () => {
      const error = new Error('Submission failed')
      mockApiService.createPatient.mockRejectedValue(error)
      
      await expect(formState.submitForm(false)).rejects.toThrow('Submission failed')
      expect(formState.isSubmitting.value).toBe(false)
    })

    it('should prevent multiple simultaneous submissions', async () => {
      mockApiService.createPatient.mockImplementation(() => 
        new Promise(resolve => setTimeout(resolve, 100))
      )
      
      const promise1 = formState.submitForm(false)
      const promise2 = formState.submitForm(false)
      
      expect(formState.isSubmitting.value).toBe(true)
      
      const result1 = await promise1
      const result2 = await promise2
      
      expect(result2).toBeNull() // Second submission should return null
    })

    it('should run additional validation when provided', async () => {
      const additionalValidation = vi.fn().mockReturnValue(false)
      
      const result = await formState.submitForm(false, null, additionalValidation)
      
      expect(additionalValidation).toHaveBeenCalled()
      expect(result).toBeNull()
      expect(mockApiService.createPatient).not.toHaveBeenCalled()
    })

    it('should handle missing record ID in edit mode', async () => {
      const result = await formState.submitForm(true) // No record ID provided
      
      expect(result).toBeNull()
      expect(mockApiService.updatePatient).not.toHaveBeenCalled()
    })
  })

  describe('computed properties', () => {
    it('should delegate computed properties to validation composable', () => {
      expect(formState.hasErrors).toBeDefined()
      expect(formState.isFormValid).toBeDefined()
      expect(formState.validationErrors).toBeDefined()
    })
  })

  describe('error management delegation', () => {
    it('should delegate error management methods', () => {
      expect(typeof formState.getFieldError).toBe('function')
      expect(typeof formState.hasFieldError).toBe('function')
      expect(typeof formState.clearFieldError).toBe('function')
      expect(typeof formState.clearAllErrors).toBe('function')
    })
  })

  describe('edge cases', () => {
    it('should handle empty enabled fields array', () => {
      formState.initializeFormData({ test: 'value' }, [])
      
      expect(formState.formData.value).toEqual({})
    })

    it('should handle null initial data', () => {
      const enabledFields = [
        { name: 'field1', field_type: 'text' }
      ]
      
      formState.initializeFormData(null, enabledFields)
      
      expect(formState.formData.value).toEqual({
        field1: ''
      })
    })

    it('should handle field without field_type', () => {
      const enabledFields = [
        { name: 'field1' } // missing field_type
      ]
      
      formState.initializeFormData({}, enabledFields)
      
      expect(formState.formData.value).toEqual({
        field1: ''
      })
    })
  })
})