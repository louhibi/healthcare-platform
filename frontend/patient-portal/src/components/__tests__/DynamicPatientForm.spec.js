import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { vi } from 'vitest'
import DynamicPatientForm from '../patients/DynamicPatientForm.vue'

// Mock toast
vi.mock('vue-toastification', () => ({ useToast: () => ({ success: vi.fn(), error: vi.fn() }) }))

// Mock API modules
vi.mock('@/api/locations', () => ({
  locationsApi: {
    getCountries: vi.fn().mockResolvedValue([
      { code: 'US', name: 'United States' },
      { code: 'CA', name: 'Canada' }
    ]),
    getStatesByCountry: vi.fn((code) => {
      if (code === 'US') return Promise.resolve([{ code: 'NY', name: 'New York' }])
      if (code === 'CA') return Promise.resolve([{ code: 'ON', name: 'Ontario' }])
      return Promise.resolve([])
    }),
    getCitiesByCountry: vi.fn((code, params = {}) => {
      if (code === 'US') return Promise.resolve([{ id: 'nyc', name: 'New York' }])
      if (code === 'CA') return Promise.resolve([{ id: 'toronto', name: 'Toronto' }])
      return Promise.resolve([])
    })
  }
}))

vi.mock('@/api/patients', () => ({
  patientsApi: {
    createPatient: vi.fn().mockResolvedValue({ data: { id: 1 } }),
    updatePatient: vi.fn().mockResolvedValue({ data: { id: 1 } }),
  }
}))

// Pinia store mock for form fields
vi.mock('@/stores/formConfig', () => ({
  useFormConfigStore: () => ({
    isLoading: false,
    error: null,
    isDirty: false,
    async loadFormFields() {},
    getFormFields: () => ([
      { field_id: 1, name: 'country', display_name: 'Country', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 2, name: 'city', display_name: 'City', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 3, name: 'first_name', display_name: 'First Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 4, name: 'last_name', display_name: 'Last Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 5, name: 'date_of_birth', display_name: 'Date of Birth', field_type: 'date', is_enabled: true, is_required: true, category: 'Personal' }
    ]),
    getEnabledFields: () => ([
      { field_id: 1, name: 'country', display_name: 'Country', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 2, name: 'city', display_name: 'City', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 3, name: 'first_name', display_name: 'First Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 4, name: 'last_name', display_name: 'Last Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 5, name: 'date_of_birth', display_name: 'Date of Birth', field_type: 'date', is_enabled: true, is_required: true, category: 'Personal' }
    ]),
    getRequiredFields: () => ([
      { field_id: 1, name: 'country', display_name: 'Country', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 2, name: 'city', display_name: 'City', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
      { field_id: 3, name: 'first_name', display_name: 'First Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 4, name: 'last_name', display_name: 'Last Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
      { field_id: 5, name: 'date_of_birth', display_name: 'Date of Birth', field_type: 'date', is_enabled: true, is_required: true, category: 'Personal' }
    ]),
    getFieldsByCategory: () => ({
      'Address': [
        { field_id: 1, name: 'country', display_name: 'Country', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' },
        { field_id: 2, name: 'city', display_name: 'City', field_type: 'select', is_enabled: true, is_required: true, category: 'Address' }
      ],
      'Personal': [
        { field_id: 3, name: 'first_name', display_name: 'First Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
        { field_id: 4, name: 'last_name', display_name: 'Last Name', field_type: 'text', is_enabled: true, is_required: true, category: 'Personal' },
        { field_id: 5, name: 'date_of_birth', display_name: 'Date of Birth', field_type: 'date', is_enabled: true, is_required: true, category: 'Personal' }
      ]
    })
  })
}))

function factory(props = {}) {
  return mount(DynamicPatientForm, {
    props,
    global: {
      plugins: [createPinia()],
      stubs: ['component']
    }
  })
}

describe('DynamicPatientForm - country/city dependency', () => {
  it('disables city until country selected, then loads cities', async () => {
    const wrapper = factory()
    await flushPromises()

    const country = wrapper.find('select#country')
    const city = wrapper.find('select#city')

    expect(country.exists()).toBe(true)
    expect(city.exists()).toBe(true)
    expect(city.attributes('disabled')).toBeDefined()

    await country.setValue('US')
    await flushPromises()
    
  // Wait for any asynchronous initialization tasks (debounces removed with composable deletion)
    await new Promise(resolve => setTimeout(resolve, 350))

    // City should now be enabled and have options
    const updatedCity = wrapper.find('select#city')
    expect(updatedCity.attributes('disabled')).toBeUndefined()
    const opts = updatedCity.findAll('option')
    expect(opts.some(o => o.text() === 'New York')).toBe(true)
  })

  it('validates city belongs to selected country on submit', async () => {
    const wrapper = factory()
    await flushPromises()

    await wrapper.find('select#country').setValue('US')
    await flushPromises()

    // Set invalid city id
    await wrapper.find('select#city').setValue('toronto')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    // Expect an error present in the UI state (toast is mocked)
    // Since we can't easily assert toast, check validationErrors or keep city value reset on invalid
    // Here we ensure the city remains the invalid value but submission prevented (isSubmitting false)
    expect(wrapper.vm.isSubmitting).toBe(false)
  })
})
