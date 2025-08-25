import { describe, it, expect } from 'vitest'
import {
  getFieldComponent,
  getFieldComponentName,
  getFieldType,
  getFieldProps,
  fieldRequiresOptions,
  isTextBasedField,
  isNumericField,
  isDateTimeField,
  isBooleanField,
  getAllFieldTypes,
  buildFieldProps
} from '../fieldComponentMap'

describe('fieldComponentMap', () => {
  describe('getFieldComponent', () => {
    it('should return correct component config for known field types', () => {
      const textConfig = getFieldComponent('text')
      expect(textConfig.component).toBe('input')
      expect(textConfig.type).toBe('text')

      const emailConfig = getFieldComponent('email')
      expect(emailConfig.component).toBe('input')
      expect(emailConfig.type).toBe('email')

      const textareaConfig = getFieldComponent('textarea')
      expect(textareaConfig.component).toBe('textarea')
      expect(textareaConfig.props.rows).toBe(3)
    })

    it('should return text config for unknown field types', () => {
      const unknownConfig = getFieldComponent('unknown_type')
      expect(unknownConfig.component).toBe('input')
      expect(unknownConfig.type).toBe('text')
    })

    it('should return new props object (not reference)', () => {
      const config1 = getFieldComponent('textarea')
      const config2 = getFieldComponent('textarea')
      
      expect(config1.props).not.toBe(config2.props)
      expect(config1.props).toEqual(config2.props)
    })
  })

  describe('getFieldComponentName', () => {
    it('should return component tag names', () => {
      expect(getFieldComponentName('text')).toBe('input')
      expect(getFieldComponentName('email')).toBe('input')
      expect(getFieldComponentName('textarea')).toBe('textarea')
      expect(getFieldComponentName('select')).toBe('select')
    })
  })

  describe('getFieldType', () => {
    it('should return input type attributes', () => {
      expect(getFieldType('text')).toBe('text')
      expect(getFieldType('email')).toBe('email')
      expect(getFieldType('password')).toBe('password')
      expect(getFieldType('number')).toBe('number')
      expect(getFieldType('date')).toBe('date')
      expect(getFieldType('checkbox')).toBe('checkbox')
    })

    it('should return undefined for components without type', () => {
      expect(getFieldType('textarea')).toBeUndefined()
      expect(getFieldType('select')).toBeUndefined()
    })
  })

  describe('getFieldProps', () => {
    it('should return default props for field types', () => {
      const textProps = getFieldProps('text')
      expect(textProps).toEqual({})

      const textareaProps = getFieldProps('textarea')
      expect(textareaProps).toEqual({ rows: 3 })

      const multiselectProps = getFieldProps('multiselect')
      expect(multiselectProps).toEqual({ multiple: true })
    })
  })

  describe('field type checkers', () => {
    describe('fieldRequiresOptions', () => {
      it('should identify fields that need options', () => {
        expect(fieldRequiresOptions('select')).toBe(true)
        expect(fieldRequiresOptions('multiselect')).toBe(true)
        expect(fieldRequiresOptions('radio')).toBe(true)
        expect(fieldRequiresOptions('text')).toBe(false)
        expect(fieldRequiresOptions('email')).toBe(false)
      })
    })

    describe('isTextBasedField', () => {
      it('should identify text-based fields', () => {
        expect(isTextBasedField('text')).toBe(true)
        expect(isTextBasedField('email')).toBe(true)
        expect(isTextBasedField('tel')).toBe(true)
        expect(isTextBasedField('phone')).toBe(true)
        expect(isTextBasedField('url')).toBe(true)
        expect(isTextBasedField('password')).toBe(true)
        expect(isTextBasedField('textarea')).toBe(true)
        expect(isTextBasedField('number')).toBe(false)
        expect(isTextBasedField('select')).toBe(false)
      })
    })

    describe('isNumericField', () => {
      it('should identify numeric fields', () => {
        expect(isNumericField('number')).toBe(true)
        expect(isNumericField('text')).toBe(false)
        expect(isNumericField('email')).toBe(false)
      })
    })

    describe('isDateTimeField', () => {
      it('should identify date/time fields', () => {
        expect(isDateTimeField('date')).toBe(true)
        expect(isDateTimeField('datetime')).toBe(true)
        expect(isDateTimeField('time')).toBe(true)
        expect(isDateTimeField('text')).toBe(false)
        expect(isDateTimeField('number')).toBe(false)
      })
    })

    describe('isBooleanField', () => {
      it('should identify boolean fields', () => {
        expect(isBooleanField('checkbox')).toBe(true)
        expect(isBooleanField('boolean')).toBe(true)
        expect(isBooleanField('text')).toBe(false)
        expect(isBooleanField('radio')).toBe(false)
      })
    })
  })

  describe('getAllFieldTypes', () => {
    it('should return array of all supported field types', () => {
      const fieldTypes = getAllFieldTypes()
      
      expect(Array.isArray(fieldTypes)).toBe(true)
      expect(fieldTypes).toContain('text')
      expect(fieldTypes).toContain('email')
      expect(fieldTypes).toContain('textarea')
      expect(fieldTypes).toContain('select')
      expect(fieldTypes).toContain('number')
      expect(fieldTypes).toContain('date')
      expect(fieldTypes).toContain('checkbox')
      expect(fieldTypes.length).toBeGreaterThan(10)
    })
  })

  describe('buildFieldProps', () => {
    it('should build basic field props', () => {
      const field = {
        name: 'test_field',
        field_type: 'text',
        placeholder_text: 'Enter text',
        is_required: true
      }

      const props = buildFieldProps(field)

      expect(props).toEqual({
        id: 'test_field',
        name: 'test_field',
        placeholder: 'Enter text',
        required: true,
        type: 'text'
      })
    })

    it('should include validation rules in props', () => {
      const field = {
        name: 'test_field',
        field_type: 'text',
        validation_rules: {
          min_length: 5,
          max_length: 50,
          pattern: '[A-Za-z]+',
          min: 10,
          max: 100,
          step: 0.1
        }
      }

      const props = buildFieldProps(field)

      expect(props).toMatchObject({
        minlength: 5,
        maxlength: 50,
        pattern: '[A-Za-z]+',
        min: 10,
        max: 100,
        step: 0.1
      })
    })

    it('should handle field without validation rules', () => {
      const field = {
        name: 'simple_field',
        field_type: 'text',
        is_required: false
      }

      const props = buildFieldProps(field)

      expect(props).toEqual({
        id: 'simple_field',
        name: 'simple_field',
        placeholder: '',
        required: false,
        type: 'text'
      })
    })

    it('should handle different field types correctly', () => {
      const emailField = {
        name: 'email_field',
        field_type: 'email',
        placeholder_text: 'Enter email'
      }

      const emailProps = buildFieldProps(emailField)
      expect(emailProps.type).toBe('email')

      const numberField = {
        name: 'number_field',
        field_type: 'number',
        validation_rules: { min: 0, max: 100 }
      }

      const numberProps = buildFieldProps(numberField)
      expect(numberProps.type).toBe('number')
      expect(numberProps.min).toBe(0)
      expect(numberProps.max).toBe(100)
    })

    it('should handle textarea field type', () => {
      const textareaField = {
        name: 'notes',
        field_type: 'textarea',
        placeholder_text: 'Enter notes'
      }

      const props = buildFieldProps(textareaField)
      expect(props.type).toBeUndefined() // textarea doesn't have type attribute
      expect(props.rows).toBe(3)
    })

    it('should handle boolean field requirements', () => {
      const checkboxField = {
        name: 'agree',
        field_type: 'checkbox',
        is_required: true
      }

      const props = buildFieldProps(checkboxField)
      expect(props.type).toBe('checkbox')
      expect(props.required).toBe(true)
    })

    it('should handle select field with multiple', () => {
      const multiselectField = {
        name: 'choices',
        field_type: 'multiselect'
      }

      const props = buildFieldProps(multiselectField)
      expect(props.multiple).toBe(true)
    })
  })

  describe('edge cases', () => {
    it('should handle null or undefined field types', () => {
      expect(() => getFieldComponent(null)).not.toThrow()
      expect(() => getFieldComponent(undefined)).not.toThrow()
      expect(getFieldComponent(null).component).toBe('input')
    })

    it('should handle empty field objects', () => {
      expect(() => buildFieldProps({})).not.toThrow()
      const props = buildFieldProps({})
      expect(props.id).toBeUndefined()
      expect(props.name).toBeUndefined()
      expect(props.placeholder).toBe('')
    })

    it('should handle field with null validation_rules', () => {
      const field = {
        name: 'test',
        field_type: 'text',
        validation_rules: null
      }

      expect(() => buildFieldProps(field)).not.toThrow()
      const props = buildFieldProps(field)
      expect(props.minlength).toBeUndefined()
    })
  })
})