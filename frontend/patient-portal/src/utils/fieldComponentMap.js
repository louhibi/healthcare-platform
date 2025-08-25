/**
 * Field component mapping system for dynamic form fields
 * Centralized configuration for field types and their rendering components
 */

// Base field component configurations
const FIELD_COMPONENTS = {
  // Basic text inputs
  text: {
    component: 'input',
    type: 'text',
    props: {}
  },
  email: {
    component: 'input',
    type: 'email',
    props: {}
  },
  tel: {
    component: 'input',
    type: 'tel',
    props: {}
  },
  phone: {
    component: 'input',
    type: 'tel',
    props: {}
  },
  url: {
    component: 'input',
    type: 'url',
    props: {}
  },
  password: {
    component: 'input',
    type: 'password',
    props: {}
  },
  
  // Number inputs
  number: {
    component: 'input',
    type: 'number',
    props: {}
  },
  
  // Date/time inputs
  date: {
    component: 'input',
    type: 'date',
    props: {}
  },
  datetime: {
    component: 'input',
    type: 'datetime-local',
    props: {}
  },
  time: {
    component: 'input',
    type: 'time',
    props: {}
  },
  
  // Text areas
  textarea: {
    component: 'textarea',
    props: {
      rows: 3
    }
  },
  
  // Select inputs
  select: {
    component: 'select',
    props: {}
  },
  multiselect: {
    component: 'select',
    props: {
      multiple: true
    }
  },
  
  // Checkboxes and radios
  checkbox: {
    component: 'input',
    type: 'checkbox',
    props: {}
  },
  boolean: {
    component: 'input',
    type: 'checkbox',
    props: {}
  },
  radio: {
    component: 'input',
    type: 'radio',
    props: {}
  },
  
  // File inputs
  file: {
    component: 'input',
    type: 'file',
    props: {}
  },
  
  // Hidden inputs
  hidden: {
    component: 'input',
    type: 'hidden',
    props: {}
  }
}

/**
 * Get the component configuration for a field type
 * @param {string} fieldType - The field type
 * @returns {Object} Component configuration object
 */
export function getFieldComponent(fieldType) {
  const config = FIELD_COMPONENTS[fieldType] || FIELD_COMPONENTS.text
  return {
    component: config.component,
    type: config.type,
    props: { ...config.props }
  }
}

/**
 * Get the HTML component tag name for a field type
 * @param {string} fieldType - The field type
 * @returns {string} HTML component tag name
 */
export function getFieldComponentName(fieldType) {
  const config = getFieldComponent(fieldType)
  return config.component
}

/**
 * Get the input type attribute for a field type
 * @param {string} fieldType - The field type
 * @returns {string|undefined} Input type attribute value
 */
export function getFieldType(fieldType) {
  const config = getFieldComponent(fieldType)
  return config.type
}

/**
 * Get the default props for a field type
 * @param {string} fieldType - The field type
 * @returns {Object} Default props object
 */
export function getFieldProps(fieldType) {
  const config = getFieldComponent(fieldType)
  return { ...config.props }
}

/**
 * Check if a field type requires options (select, radio, etc.)
 * @param {string} fieldType - The field type
 * @returns {boolean} True if the field type requires options
 */
export function fieldRequiresOptions(fieldType) {
  return ['select', 'multiselect', 'radio'].includes(fieldType)
}

/**
 * Check if a field type is a text-based input
 * @param {string} fieldType - The field type
 * @returns {boolean} True if the field type is text-based
 */
export function isTextBasedField(fieldType) {
  return ['text', 'email', 'tel', 'phone', 'url', 'password', 'textarea'].includes(fieldType)
}

/**
 * Check if a field type is a numeric input
 * @param {string} fieldType - The field type
 * @returns {boolean} True if the field type is numeric
 */
export function isNumericField(fieldType) {
  return ['number'].includes(fieldType)
}

/**
 * Check if a field type is a date/time input
 * @param {string} fieldType - The field type
 * @returns {boolean} True if the field type is date/time related
 */
export function isDateTimeField(fieldType) {
  return ['date', 'datetime', 'time'].includes(fieldType)
}

/**
 * Check if a field type is a boolean input
 * @param {string} fieldType - The field type
 * @returns {boolean} True if the field type is boolean
 */
export function isBooleanField(fieldType) {
  return ['checkbox', 'boolean'].includes(fieldType)
}

/**
 * Get all supported field types
 * @returns {Array<string>} Array of all supported field types
 */
export function getAllFieldTypes() {
  return Object.keys(FIELD_COMPONENTS)
}

/**
 * Build complete field props including validation rules and field-specific attributes
 * @param {Object} field - The field configuration object
 * @returns {Object} Complete props object for the field
 */
export function buildFieldProps(field) {
  const baseConfig = getFieldComponent(field.field_type)
  const props = {
    id: field.name,
    name: field.name,
    placeholder: field.placeholder_text || '',
    required: !!field.is_required,
    ...baseConfig.props
  }

  // Add input type if it exists
  if (baseConfig.type) {
    props.type = baseConfig.type
  }

  // Add validation attributes from rules
  if (field.validation_rules) {
    const rules = field.validation_rules
    if (rules.min_length) props.minlength = rules.min_length
    if (rules.max_length) props.maxlength = rules.max_length
    if (rules.pattern) props.pattern = rules.pattern
    if (rules.min !== undefined) props.min = rules.min
    if (rules.max !== undefined) props.max = rules.max
    if (rules.step) props.step = rules.step
  }

  return props
}

export default {
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
}