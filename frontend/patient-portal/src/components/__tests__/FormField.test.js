import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import FormField from '../FormField.vue'
import DynamicFieldLoader from '../fields/DynamicFieldLoader.vue'

// Mock DynamicFieldLoader
vi.mock('../fields/DynamicFieldLoader.vue', () => ({
  default: {
    name: 'DynamicFieldLoader',
    template: '<div class="mock-dynamic-field-loader">{{ fieldType }} field</div>',
    props: ['fieldType', 'field', 'modelValue', 'options', 'hasError', 'disabled', 'loading'],
    emits: ['update:modelValue', 'input', 'change']
  }
}))

describe('FormField', () => {
  const mockField = {
    field_id: 1,
    name: 'test_field',
    field_type: 'text',
    display_name: 'Test Field',
    is_enabled: true,
    is_required: true,
    description: 'Test field description'
  }

  let wrapper

  beforeEach(() => {
    wrapper = mount(FormField, {
      props: {
        field: mockField,
        modelValue: 'test value',
        options: [],
        errors: [],
        loading: false,
        disabled: false
      }
    })
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
      wrapper = null
    }
  })

  describe('rendering', () => {
    it('should render field label with display name', () => {
      const label = wrapper.find('label')
      
      expect(label.exists()).toBe(true)
      expect(label.text()).toContain('Test Field')
      expect(label.attributes('for')).toBe('test_field')
    })

    it('should show required asterisk for required fields', () => {
      const requiredSpan = wrapper.find('.text-red-500')
      
      expect(requiredSpan.exists()).toBe(true)
      expect(requiredSpan.text()).toBe('*')
    })

    it('should not show required asterisk for optional fields', async () => {
      const optionalField = { ...mockField, is_required: false }
      await wrapper.setProps({ field: optionalField })
      
      const requiredSpan = wrapper.find('.text-red-500')
      expect(requiredSpan.exists()).toBe(false)
    })

    it('should render DynamicFieldLoader component', () => {
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      
      expect(dynamicLoader.exists()).toBe(true)
      expect(dynamicLoader.props('fieldType')).toBe('text')
      expect(dynamicLoader.props('field')).toEqual(mockField)
      expect(dynamicLoader.props('modelValue')).toBe('test value')
    })

    it('should render field description when provided', () => {
      const description = wrapper.find('.text-xs.text-gray-500')
      
      expect(description.exists()).toBe(true)
      expect(description.text()).toBe('Test field description')
    })

    it('should not render description when not provided', async () => {
      const fieldWithoutDescription = { ...mockField }
      delete fieldWithoutDescription.description
      
      await wrapper.setProps({ field: fieldWithoutDescription })
      
      const description = wrapper.find('.text-xs.text-gray-500')
      expect(description.exists()).toBe(false)
    })
  })

  describe('error handling', () => {
    it('should display field errors', async () => {
      const errors = ['This field is required', 'Must be at least 5 characters']
      await wrapper.setProps({ errors })
      
      const errorContainer = wrapper.find('.mt-1')
      const errorMessages = wrapper.findAll('.text-sm.text-red-600')
      
      expect(errorContainer.exists()).toBe(true)
      expect(errorMessages).toHaveLength(2)
      expect(errorMessages[0].text()).toBe('This field is required')
      expect(errorMessages[1].text()).toBe('Must be at least 5 characters')
    })

    it('should not display error container when no errors', () => {
      const errorMessages = wrapper.findAll('.text-sm.text-red-600')
      expect(errorMessages).toHaveLength(0)
    })

    it('should pass hasError prop to DynamicFieldLoader', async () => {
      await wrapper.setProps({ errors: ['Error message'] })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('hasError')).toBe(true)
    })
  })

  describe('props validation', () => {
    it('should validate field prop with required properties', () => {
      const validField = {
        name: 'test',
        field_type: 'text',
        display_name: 'Test',
        is_enabled: true
      }

      expect(() => {
        mount(FormField, {
          props: { field: validField }
        })
      }).not.toThrow()
    })

    it('should handle invalid field prop gracefully', () => {
      // This would trigger validation warning in development
      const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
      
      try {
        mount(FormField, {
          props: { field: { invalid: 'field' } }
        })
      } catch (error) {
        // Expected to potentially throw due to validation
      }
      
      consoleSpy.mockRestore()
    })

    it('should validate options prop', async () => {
      const validOptions = [
        { value: 'option1', label: 'Option 1' },
        { value: 'option2', label: 'Option 2' },
        'string_option'
      ]

      await wrapper.setProps({ options: validOptions })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('options')).toEqual(validOptions)
    })

    it('should validate errors prop', async () => {
      const validErrors = ['Error 1', 'Error 2']

      await wrapper.setProps({ errors: validErrors })
      
      const errorMessages = wrapper.findAll('.text-sm.text-red-600')
      expect(errorMessages).toHaveLength(2)
    })
  })

  describe('event handling', () => {
    it('should handle model value updates', async () => {
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      
      await dynamicLoader.vm.$emit('update:modelValue', 'new value')
      
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual(['new value'])
    })

    it('should handle input events', async () => {
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      
      await dynamicLoader.vm.$emit('input', 'input value')
      
      expect(wrapper.emitted('input')).toBeTruthy()
      expect(wrapper.emitted('input')[0]).toEqual(['input value'])
    })

    it('should handle change events', async () => {
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      
      await dynamicLoader.vm.$emit('change', 'changed value')
      
      expect(wrapper.emitted('change')).toBeTruthy()
      expect(wrapper.emitted('change')[0]).toEqual(['changed value'])
    })
  })

  describe('field grid classes', () => {
    it('should apply wide class for textarea fields', async () => {
      const textareaField = { ...mockField, field_type: 'textarea' }
      await wrapper.setProps({ field: textareaField })
      
      expect(wrapper.classes()).toContain('md:col-span-2')
    })

    it('should apply wide class for specific field names', async () => {
      const addressField = { ...mockField, name: 'address' }
      await wrapper.setProps({ field: addressField })
      
      expect(wrapper.classes()).toContain('md:col-span-2')
    })

    it('should not apply wide class for regular fields', () => {
      expect(wrapper.classes()).not.toContain('md:col-span-2')
    })
  })

  describe('loading and disabled states', () => {
    it('should pass loading state to DynamicFieldLoader', async () => {
      await wrapper.setProps({ loading: true })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('loading')).toBe(true)
    })

    it('should pass disabled state to DynamicFieldLoader', async () => {
      await wrapper.setProps({ disabled: true })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('disabled')).toBe(true)
    })
  })

  describe('accessibility', () => {
    it('should have proper label-field association', () => {
      const label = wrapper.find('label')
      const fieldId = mockField.name
      
      expect(label.attributes('for')).toBe(fieldId)
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('field').name).toBe(fieldId)
    })

    it('should provide accessible error messages', async () => {
      await wrapper.setProps({ errors: ['Test error'] })
      
      const errorMessage = wrapper.find('.text-sm.text-red-600')
      expect(errorMessage.exists()).toBe(true)
    })
  })

  describe('field options handling', () => {
    it('should compute field options correctly', async () => {
      const fieldOptions = [
        { value: 'option1', label: 'Option 1' },
        { value: 'option2', label: 'Option 2' }
      ]
      
      await wrapper.setProps({ options: fieldOptions })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('options')).toEqual(fieldOptions)
    })

    it('should pass options prop to DynamicFieldLoader', async () => {
      const options = [{ value: 'field_option', label: 'Field Option' }]
      
      await wrapper.setProps({ options })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('options')).toEqual(options)
    })
  })

  describe('reactivity', () => {
    it('should update when field prop changes', async () => {
      const newField = {
        ...mockField,
        display_name: 'Updated Field Name'
      }
      
      await wrapper.setProps({ field: newField })
      
      const label = wrapper.find('label')
      expect(label.text()).toContain('Updated Field Name')
    })

    it('should update when modelValue changes', async () => {
      await wrapper.setProps({ modelValue: 'updated value' })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('modelValue')).toBe('updated value')
    })
  })

  describe('edge cases', () => {
    it('should handle field without display_name', async () => {
      const fieldWithoutDisplayName = { ...mockField }
      delete fieldWithoutDisplayName.display_name
      
      await wrapper.setProps({ field: fieldWithoutDisplayName })
      
      const label = wrapper.find('label')
      expect(label.exists()).toBe(true)
    })

    it('should handle null modelValue', async () => {
      await wrapper.setProps({ modelValue: null })
      
      const dynamicLoader = wrapper.findComponent(DynamicFieldLoader)
      expect(dynamicLoader.props('modelValue')).toBeNull()
    })

    it('should handle undefined errors', async () => {
      await wrapper.setProps({ errors: undefined })
      
      const errorMessages = wrapper.findAll('.text-sm.text-red-600')
      expect(errorMessages).toHaveLength(0)
    })
  })
})