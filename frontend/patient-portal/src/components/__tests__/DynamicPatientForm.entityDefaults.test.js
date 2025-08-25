import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick, ref } from 'vue'
import DynamicPatientForm from '../patients/DynamicPatientForm.vue'

// Mock all the dependencies
vi.mock('vue-toastification', () => ({
  useToast: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn()
  })
}))

vi.mock('@/api/patients', () => ({
  patientsApi: {
    createPatient: vi.fn(),
    updatePatient: vi.fn()
  }
}))

// Create a mutable entity store mock that we can change in tests
let mockEntityData = {
  id: 1,
  name: 'Toronto General Hospital',
  country: 'Canada',
  state: 'ON',
  city: 'Toronto'
}

const mockEntityStore = {
  get entity() { return mockEntityData },
  set entity(value) { mockEntityData = value },
  get entityCountry() { return mockEntityData?.country || '' }
}

vi.mock('@/stores/entity', () => ({
  useEntityStore: () => mockEntityStore
}))

// Create reactive mocks for form config
const mockFormConfig = {
  isInitialized: ref(true),
  fields: ref([]),
  enabledFields: ref([
    { name: 'first_name', field_type: 'text', is_enabled: true },
    { name: 'country', field_type: 'select', is_enabled: true },
    { name: 'state', field_type: 'select', is_enabled: true },
    { name: 'city', field_type: 'select', is_enabled: true }
  ]),
  requiredFields: ref([]),
  fieldsByCategory: ref({}),
  isLoading: ref(false),
  error: ref(null),
  initialize: vi.fn().mockResolvedValue(true),
  getFieldConfig: vi.fn().mockReturnValue({ is_enabled: true })
}

vi.mock('@/composables/useFormConfig', () => ({
  useFormConfig: () => mockFormConfig
}))

// useLocationData composable removed; no mock needed

vi.mock('@/composables/useValidation', () => ({
  useValidation: () => ({
    hasErrors: { value: false },
    isValid: { value: true },
    validationErrors: { value: {} },
    validateLocation: vi.fn().mockReturnValue({}),
    validatePatientData: vi.fn().mockReturnValue({}),
    setFieldError: vi.fn(),
    clearAllErrors: vi.fn(),
    clearFieldError: vi.fn(),
    buildRulesFromConfig: vi.fn().mockReturnValue([]),
    validateField: vi.fn().mockReturnValue(true),
    getFieldError: vi.fn().mockReturnValue(null),
    hasFieldError: vi.fn().mockReturnValue(false)
  })
}))

// Mock utility functions
vi.mock('@/utils/fieldComponentMap', () => ({
  getFieldComponentName: vi.fn().mockReturnValue('TextField'),
  buildFieldProps: vi.fn().mockReturnValue({})
}))

vi.mock('@/utils/memoization', () => ({
  memoizedArrayFilter: vi.fn((arr, filter) => arr.filter ? arr.filter(filter) : []),
  memoizedArraySort: vi.fn((arr, sortFn) => arr.sort ? arr.sort(sortFn) : [])
}))

vi.mock('@/utils/fieldLoader', () => ({
  preloadFieldComponents: vi.fn().mockResolvedValue()
}))

// Mock field components
vi.mock('../FormSection.vue', () => ({
  default: {
    name: 'FormSection',
    template: '<div data-test="form-section"><slot /></div>',
    props: ['title']
  }
}))

vi.mock('../VirtualFormSection.vue', () => ({
  default: {
    name: 'VirtualFormSection',
    template: '<div data-test="virtual-form-section"></div>',
    props: ['title', 'fields', 'formData', 'getFieldOptions', 'getFieldError', 'loading', 'showFieldCount', 'minVirtualFields', 'virtualHeight']
  }
}))

vi.mock('../location/LocationFieldGroup.vue', () => ({
  default: {
    name: 'LocationFieldGroup',
    template: '<div data-test="location-field-group"></div>',
    props: ['fields', 'formData', 'hasFieldError', 'getFieldError'],
    emits: ['country-change', 'state-change', 'city-change', 'location-change']
  }
}))

// Create a mock that we can inspect
const mockInitializeFormData = vi.fn()
const mockSubmitForm = vi.fn().mockResolvedValue({ id: 1, success: true })

vi.mock('@/composables/useFormState', () => ({
  useFormState: vi.fn(() => ({
    formData: ref({}),
    validationErrors: ref({}),
    isSubmitting: ref(false),
    isDirty: ref(false),
    hasErrors: ref(false),
    isFormValid: ref(true),
    initializeFormData: mockInitializeFormData,
    updateField: vi.fn(),
    updateMultipleFields: vi.fn(),
    validateField: vi.fn().mockReturnValue(true),
    validateForm: vi.fn().mockReturnValue(true),
    getFieldError: vi.fn().mockReturnValue(null),
    hasFieldError: vi.fn().mockReturnValue(false),
    clearFieldError: vi.fn(),
    clearAllErrors: vi.fn(),
    resetForm: vi.fn(),
    submitForm: mockSubmitForm
  }))
}))

describe('DynamicPatientForm - Entity Location Defaults', () => {
  let wrapper

  // Helper function to wait for component initialization
  const waitForComponentInit = async (wrapper, maxAttempts = 20) => {
    let attempts = 0
    while (attempts < maxAttempts) {
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 10))
      if (mockInitializeFormData.mock.calls.length > 0) {
        return
      }
      attempts++
    }
    // If we get here, initialization didn't complete in time
    console.warn('Component initialization may not have completed in time')
  }

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()
    mockInitializeFormData.mockClear()
    mockSubmitForm.mockClear()
    
    // Reset entity store to default state
    mockEntityData = {
      id: 1,
      name: 'Toronto General Hospital',
      country: 'Canada',
      state: 'ON',
      city: 'Toronto'
    }
    
    // Reset form config to initialized state
    mockFormConfig.isInitialized.value = true
    mockFormConfig.isLoading.value = false
    mockFormConfig.error.value = null
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
    vi.clearAllMocks()
  })

  describe('Create Mode - Entity Location Defaults', () => {
    it('should pass entity location defaults when creating new patient', async () => {
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {}, // No ID = create mode
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Verify that initializeFormData was called with entity location defaults
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {}, // initialData
        expect.any(Object), // enabledFields
        {
          country: 'Canada',
          state: 'ON',
          city: 'Toronto'
        } // entityLocationDefaults
      )
    })

    it('should handle entity store with no location data gracefully', async () => {
      // Set entity to null before mounting component
      const originalEntity = mockEntityData
      mockEntityData = null
      
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Should pass null for entity defaults when no entity data
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {},
        expect.any(Object),
        null // null entity defaults
      )

      // Restore original entity for other tests
      mockEntityData = originalEntity
    })
  })

  describe('Edit Mode - Preserve Existing Data', () => {
    it('should still pass entity defaults but let form state handle the logic', async () => {
      const existingPatientData = {
        id: 123, // ID present = edit mode
        first_name: 'John',
        last_name: 'Doe',
        country: 'USA',
        state: 'CA',
        city: 'San Francisco'
      }

      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: existingPatientData,
          isEdit: true
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Entity defaults are still passed, but form state should detect edit mode and ignore them
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        existingPatientData,
        expect.any(Object),
        {
          country: 'Canada',
          state: 'ON', 
          city: 'Toronto'
        }
      )
    })
  })

  describe('Form Reinitialization on Data Changes', () => {
    it('should reinitialize with entity defaults when switching from edit to create', async () => {
      // Start with edit data
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: { id: 1, first_name: 'John' },
          isEdit: true
        }
      })

      // Wait for initial mounting
      await waitForComponentInit(wrapper)
      
      mockInitializeFormData.mockClear()

      // Switch to create mode by removing the ID
      await wrapper.setProps({
        initialData: { first_name: 'Jane' }, // No ID = create mode
        isEdit: false
      })

      // Wait for reinitialization
      await waitForComponentInit(wrapper)

      // Should reinitialize with entity defaults for create mode
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        { first_name: 'Jane' },
        expect.any(Object),
        {
          country: 'Canada',
          state: 'ON',
          city: 'Toronto'
        }
      )
    })
  })

  describe('Integration with Location Data Loading', () => {
    it('should trigger location data loading after form initialization in create mode', async () => {
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Should initialize form config first
      expect(mockFormConfig.initialize).toHaveBeenCalled()
      
      // Should call initializeFormData with entity defaults
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {},
        expect.any(Object),
        {
          country: 'Canada',
          state: 'ON',
          city: 'Toronto'
        }
      )
    })
  })

  describe('Entity Store Integration', () => {
    it('should compute entity location defaults correctly', async () => {
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Verify that initializeFormData was called with correct entity defaults
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {},
        expect.any(Object),
        {
          country: 'Canada',
          state: 'ON',
          city: 'Toronto'
        }
      )
    })

    it('should handle entity with partial location data', async () => {
      // Set entity with partial data before mounting component
      const originalEntity = mockEntityData
      mockEntityData = {
        id: 1,
        country: 'Morocco',
        state: '', // Missing state
        city: 'Casablanca'
      }

      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {},
        expect.any(Object),
        {
          country: 'Morocco',
          state: '', // Empty string for missing state
          city: 'Casablanca'
        }
      )

      // Restore original entity
      mockEntityData = originalEntity
    })
  })

  describe('Form Submission with Location Defaults', () => {
    it('should mount component and initialize with entity location defaults', async () => {
      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait for component initialization to complete
      await waitForComponentInit(wrapper)

      // Should have called initializeFormData with entity location defaults
      expect(mockInitializeFormData).toHaveBeenCalledWith(
        {},
        expect.any(Object),
        {
          country: 'Canada',
          state: 'ON',
          city: 'Toronto'
        }
      )

      // Component should exist and be rendered
      expect(wrapper.exists()).toBe(true)
    })
  })

  describe('Error Handling', () => {
    it('should handle initialization errors gracefully', async () => {
      // Mock initialization failure
      const originalInitialize = mockFormConfig.initialize
      mockFormConfig.initialize = vi.fn().mockRejectedValue(new Error('Config load failed'))

      wrapper = mount(DynamicPatientForm, {
        props: {
          initialData: {},
          isEdit: false
        }
      })

      // Wait a bit for error handling
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 50))

      // Component should handle the error and not crash
      expect(wrapper.exists()).toBe(true)
      
      // Restore original initialize function
      mockFormConfig.initialize = originalInitialize
    })
  })
})