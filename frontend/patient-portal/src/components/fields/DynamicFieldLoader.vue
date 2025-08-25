<template>
  <!-- Loading State -->
  <div v-if="isLoading" class="animate-pulse">
    <div class="h-4 bg-gray-200 rounded mb-2"></div>
    <div class="h-8 bg-gray-200 rounded"></div>
  </div>
  
  <!-- Error State -->
  <div v-else-if="loadError" class="text-red-600 text-sm">
    Failed to load field component: {{ fieldType }}
  </div>
  
  <!-- Dynamic Component -->
  <component
    v-else-if="dynamicComponent"
    :is="dynamicComponent"
    v-bind="componentProps"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    @input="$emit('input', $event)"
    @change="$emit('change', $event)"
  >
    <template v-for="(_, name) in $slots" #[name]="slotProps">
      <slot :name="name" v-bind="slotProps" />
    </template>
  </component>
  
  <!-- Fallback to basic HTML input -->
  <input
    v-else
    v-bind="fallbackProps"
    :value="modelValue"
    :class="fieldClasses"
    @input="handleFallbackInput"
    @change="handleFallbackChange"
  />
</template>

<script setup>
import { buildFieldProps } from '@/utils/fieldComponentMap'
import { loadFieldComponent } from '@/utils/fieldLoader'
import { computed, ref, watch } from 'vue'

// Props
const props = defineProps({
  fieldType: {
    type: String,
    required: true
  },
  field: {
    type: Object,
    required: true
  },
  modelValue: {
    type: [String, Number, Boolean, Array, Object],
    default: null
  },
  options: {
    type: Array,
    default: () => []
  },
  hasError: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => true,
  'input': (value) => true,
  'change': (value) => true
})

// State
const dynamicComponent = ref(null)
const isLoading = ref(false)
const loadError = ref(false)

// Field component mapping with async loading
const FIELD_COMPONENT_MAP = {
  // Complex components that benefit from dynamic loading
  date: () => import('./DateField.vue'),
  datetime: () => import('./DateField.vue'),
  time: () => import('./DateField.vue'),
  textarea: () => import('./TextareaField.vue'),
  select: () => import('./SelectField.vue'),
  multiselect: () => import('./SelectField.vue'),
  number: () => import('./NumberField.vue'),
  checkbox: () => import('./CheckboxField.vue'),
  boolean: () => import('./CheckboxField.vue'),
  file: () => import('./FileField.vue'),
  
  // ID-based select field for location and reference fields
  id_select: () => import('./IdSelectField.vue'),
  
  // Insurance field components (moved to insurance/ folder)
  insurance_type: () => import('../insurance/InsuranceTypeField.vue'),
  insurance_provider: () => import('../insurance/InsuranceProviderField.vue'),
  
  // Simple components that can use basic HTML elements (no async loading needed)
  text: null,
  email: null,
  tel: null,
  phone: null,
  url: null,
  password: null,
  radio: null,
  hidden: null
}

// Computed
const componentProps = computed(() => {
  const baseProps = buildFieldProps(props.field)
  return {
    ...baseProps,
    field: props.field,
    fieldName: props.field.name, // Add field name for debugging
    options: props.options,
    hasError: props.hasError,
    disabled: props.disabled || props.loading,
    loading: props.loading,
    placeholder: props.field.placeholder_text || baseProps.placeholder,
    required: props.field.is_required // Ensure required status is passed
  }
})

const fallbackProps = computed(() => {
  const baseProps = buildFieldProps(props.field)
  return {
    ...baseProps,
    disabled: props.disabled || props.loading,
    class: fieldClasses.value
  }
})

const fieldClasses = computed(() => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
  const border = props.hasError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  const disabledClass = (props.disabled || props.loading) ? 'bg-gray-50 text-gray-500' : ''
  return `${base} ${border} ${disabledClass}`
})

// Helper function to check if field is ID-based
const isIdBasedField = (field) => {
  const fieldName = field.name?.toLowerCase() || ''
  const idFields = [
    'country_id', 'state_id', 'city_id', 'nationality_id',
    'doctor_id', 'patient_id', 'user_id', 'entity_id'
  ]
  
  // Direct match
  if (idFields.includes(fieldName)) {
    return true
  }
  
  // Pattern match for _id suffix
  if (fieldName.endsWith('_id')) {
    return true
  }
  
  // Pattern match for id_ prefix  
  if (fieldName.startsWith('id_')) {
    return true
  }
  
  return false
}

// Methods
const loadComponent = async (fieldType) => {
  let actualFieldType = fieldType
  
  // Auto-detect ID-based select fields
  if (fieldType === 'select' && isIdBasedField(props.field)) {
    actualFieldType = 'id_select'
  }
  
  const componentLoader = FIELD_COMPONENT_MAP[actualFieldType]
  
  // If no custom component is defined, use fallback
  if (!componentLoader) {
    dynamicComponent.value = null
    return
  }
  
  isLoading.value = true
  loadError.value = false
  
  try {
    // Load the component using the enhanced field loader
    const component = await loadFieldComponent(actualFieldType, componentLoader)
    dynamicComponent.value = component
  } catch (error) {
    console.error(`Failed to load component for field type: ${actualFieldType}`, error)
    loadError.value = true
    dynamicComponent.value = null
  } finally {
    isLoading.value = false
  }
}

const handleFallbackInput = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  emit('input', value)
}

const handleFallbackChange = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  emit('change', value)
}

// Watch for field type changes
watch(() => props.fieldType, (newType) => {
  if (newType) {
    loadComponent(newType)
  }
}, { immediate: true })
</script>