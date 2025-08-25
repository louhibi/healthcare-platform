import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { ref, computed } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import { useFormConfig } from '../useFormConfig'
import { useFormConfigStore } from '@/stores/formConfig'

// Mock the form config store that useFormConfig uses
vi.mock('@/stores/formConfig', () => ({
  useFormConfigStore: vi.fn()
}))

describe('useFormConfig', () => {
  let formConfig
  let mockStore

  const mockFields = [
    {
      field_id: 1,
      name: 'first_name',
      field_type: 'text',
      display_name: 'First Name',
      is_enabled: true,
      is_required: true,
      is_core: false,
      sort_order: 1,
      category: 'Personal Information',
      placeholder_text: 'Enter first name',
      description: 'Patient first name',
      validation_rules: { min_length: 2, max_length: 50 }
    },
    {
      field_id: 2,
      name: 'email',
      field_type: 'email',
      display_name: 'Email Address',
      is_enabled: true,
      is_required: false,
      is_core: false,
      sort_order: 2,
      category: 'Contact Information',
      placeholder_text: 'Enter email address'
    },
    {
      field_id: 3,
      name: 'notes',
      field_type: 'textarea',
      display_name: 'Notes',
      is_enabled: false,
      is_required: false,
      is_core: false,
      sort_order: 3,
      category: 'Additional Information'
    },
    {
      field_id: 4,
      name: 'gender',
      field_type: 'select',
      display_name: 'Gender',
      is_enabled: true,
      is_required: true,
      is_core: false,
      sort_order: 4,
      category: 'Personal Information',
      options: [
        { value: 'male', label: 'Male' },
        { value: 'female', label: 'Female' },
        { value: 'other', label: 'Other' }
      ]
    }
  ]

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()

    // Create mock store
    mockStore = {
      getFormFields: vi.fn(() => mockFields),
      getEnabledFields: vi.fn(() => mockFields.filter(f => f.is_enabled)),
      getRequiredFields: vi.fn(() => mockFields.filter(f => f.is_enabled && f.is_required)),
      getFieldsByCategory: vi.fn(() => {
        const categories = {}
        mockFields.filter(f => f.is_enabled).forEach(field => {
          const category = field.category || 'Other'
          if (!categories[category]) categories[category] = []
          categories[category].push(field)
        })
        return categories
      }),
      getCoreFields: vi.fn(() => mockFields.filter(f => f.is_core)),
      getConfigurableFields: vi.fn(() => mockFields.filter(f => !f.is_core)),
      isLoading: ref(false),
      error: ref(null),
      isDirty: ref(false),
      loadFormFields: vi.fn().mockResolvedValue({ fields: mockFields })
    }

    useFormConfigStore.mockReturnValue(mockStore)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  beforeEach(() => {
    formConfig = useFormConfig('patient')
  })

  describe('initialization', () => {
    it('should initialize automatically due to watcher', async () => {
      // The watcher runs immediately and calls initialize
      await vi.waitFor(() => {
        expect(formConfig.isInitialized.value).toBe(true)
      })
      expect(formConfig.fields.value).toEqual(mockFields)
      expect(formConfig.isLoading.value.value).toBe(false)
      expect(formConfig.error.value.value).toBe(null)
    })

    it('should initialize with form type', () => {
      const appointmentConfig = useFormConfig('appointment')
      expect(appointmentConfig).toBeDefined()
    })
  })

  describe('configuration loading', () => {
    it('should load form configuration successfully', async () => {
      await formConfig.initialize()
      
      expect(formConfig.isInitialized.value).toBe(true)
      expect(formConfig.fields.value).toEqual(mockFields)
      expect(formConfig.isLoading.value.value).toBe(false)
      expect(formConfig.error.value.value).toBe(null)
      expect(mockStore.loadFormFields).toHaveBeenCalledWith('patient')
    })


    it('should not reload if already initialized', async () => {
      await formConfig.initialize()
      await formConfig.initialize() // Second call
      
      expect(mockStore.loadFormFields).toHaveBeenCalledTimes(1)
    })

    it('should allow force reload by reinitializing', async () => {
      await formConfig.initialize()
      
      // Reset initialization to allow reload
      formConfig.isInitialized.value = false
      await formConfig.initialize()
      
      expect(mockStore.loadFormFields).toHaveBeenCalledTimes(2)
    })
  })

  describe('field filtering', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should filter enabled fields', () => {
      const enabledFields = formConfig.enabledFields.value
      
      expect(enabledFields).toHaveLength(3) // first_name, email, gender
      expect(enabledFields.every(field => field.is_enabled)).toBe(true)
      expect(enabledFields.find(field => field.name === 'notes')).toBeUndefined()
    })

    it('should filter required fields', () => {
      const requiredFields = formConfig.requiredFields.value
      
      expect(requiredFields).toHaveLength(2) // first_name, gender
      expect(requiredFields.every(field => field.is_required)).toBe(true)
      expect(requiredFields.find(field => field.name === 'email')).toBeUndefined()
    })
  })

  describe('field categorization', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should group fields by category', () => {
      const fieldsByCategory = formConfig.fieldsByCategory.value
      
      expect(fieldsByCategory).toHaveProperty('Personal Information')
      expect(fieldsByCategory).toHaveProperty('Contact Information')
      expect(fieldsByCategory['Personal Information']).toHaveLength(2) // first_name, gender
      expect(fieldsByCategory['Contact Information']).toHaveLength(1) // email
    })

    it('should only include enabled fields in categories', () => {
      const fieldsByCategory = formConfig.fieldsByCategory.value
      
      // notes field is disabled, so Additional Information category should not exist
      expect(fieldsByCategory).not.toHaveProperty('Additional Information')
    })

    it('should sort fields within categories by sort_order', () => {
      const personalFields = formConfig.fieldsByCategory.value['Personal Information']
      
      expect(personalFields[0].sort_order).toBeLessThanOrEqual(personalFields[1].sort_order)
    })
  })

  describe('field lookup', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should get field configuration by name', () => {
      const fieldConfig = formConfig.getFieldConfig('first_name')
      
      expect(fieldConfig).toBeDefined()
      expect(fieldConfig.name).toBe('first_name')
      expect(fieldConfig.field_type).toBe('text')
      expect(fieldConfig.is_required).toBe(true)
    })

    it('should return undefined for non-existent field', () => {
      const fieldConfig = formConfig.getFieldConfig('non_existent_field')
      
      expect(fieldConfig).toBeUndefined()
    })

    it('should get field display name', () => {
      expect(formConfig.getFieldDisplayName('first_name')).toBe('First Name')
      expect(formConfig.getFieldDisplayName('email')).toBe('Email Address')
    })
  })

  describe('field validation', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should check if field is required', () => {
      expect(formConfig.isFieldRequired('first_name')).toBe(true)
      expect(formConfig.isFieldRequired('email')).toBe(false)
      expect(formConfig.isFieldRequired('non_existent')).toBe(false)
    })

    it('should check if field is enabled', () => {
      expect(formConfig.isFieldEnabled('first_name')).toBe(true)
      expect(formConfig.isFieldEnabled('notes')).toBe(false)
      expect(formConfig.isFieldEnabled('non_existent')).toBe(false)
    })

    it('should get field options', () => {
      const genderOptions = formConfig.getFieldOptions('gender')
      
      expect(genderOptions).toEqual([
        { value: 'male', label: 'Male' },
        { value: 'female', label: 'Female' },
        { value: 'other', label: 'Other' }
      ])
    })

    it('should return empty array for field without options', () => {
      const textOptions = formConfig.getFieldOptions('first_name')
      
      expect(textOptions).toEqual([])
    })
  })

  describe('field type utilities', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should identify field types', () => {
      expect(formConfig.getFieldType('first_name')).toBe('text')
      expect(formConfig.getFieldType('email')).toBe('email')
      expect(formConfig.getFieldType('gender')).toBe('select')
    })

    it('should check if field is core', () => {
      expect(formConfig.isFieldCore('first_name')).toBe(false)
      expect(formConfig.isFieldCore('email')).toBe(false)
    })

    it('should get field category', () => {
      expect(formConfig.getFieldCategory('first_name')).toBe('Personal Information')
      expect(formConfig.getFieldCategory('email')).toBe('Contact Information')
    })

    it('should get field description', () => {
      expect(formConfig.getFieldDescription('first_name')).toBe('Patient first name')
      expect(formConfig.getFieldDescription('non_existent')).toBe('')
    })
  })

  describe('form data management', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should initialize form data with enabled fields', () => {
      const initialData = { first_name: 'John', email: 'john@example.com' }
      const formData = formConfig.initializeFormData(initialData)
      
      expect(formData.first_name).toBe('John')
      expect(formData.email).toBe('john@example.com')
      expect(formData.gender).toBe('male') // Default from select options
    })

    it('should update form data', () => {
      formConfig.initializeFormData({})
      formConfig.updateFormData('first_name', 'Jane')
      
      expect(formConfig.formData.value.first_name).toBe('Jane')
    })

    it('should validate individual fields', () => {
      formConfig.initializeFormData({})
      const isValid = formConfig.validateField('first_name', '')
      
      expect(isValid).toBe(false)
      expect(formConfig.hasFieldError('first_name')).toBe(true)
    })

    it('should validate the entire form', () => {
      const testData = { first_name: '', email: 'valid@email.com', gender: 'male' }
      const isValid = formConfig.validateForm(testData)
      
      expect(isValid).toBe(false) // first_name is required but empty
    })
  })

  describe('field rendering', () => {
    beforeEach(async () => {
      await formConfig.initialize()
    })

    it('should render text field configuration', () => {
      const field = mockFields[0] // first_name
      const renderConfig = formConfig.renderField(field, 'John', vi.fn())
      
      expect(renderConfig.component).toBe('input')
      expect(renderConfig.props.type).toBe('text')
      expect(renderConfig.props.name).toBe('first_name')
    })

    it('should render select field configuration', () => {
      const field = mockFields[3] // gender
      const renderConfig = formConfig.renderField(field, 'male', vi.fn())
      
      expect(renderConfig.component).toBe('select')
      expect(renderConfig.options).toEqual(field.options)
    })

    it('should get form fields for rendering', () => {
      formConfig.initializeFormData({})
      const renderFields = formConfig.getFormFieldsForRendering()
      
      expect(renderFields).toHaveLength(3) // Only enabled fields
      expect(renderFields[0]).toHaveProperty('renderConfig')
    })
  })
})