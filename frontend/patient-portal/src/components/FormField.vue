<template>
  <div :class="getFieldGridClass">
    <label :for="field.name" class="block text-sm font-medium text-gray-700 mb-1">
      {{ field.display_name }}
      <span v-if="field.is_required" class="text-red-500">*</span>
    </label>

    <!-- Dynamic Field Loader -->
    <DynamicFieldLoader
      :field-type="field.field_type"
      :field="field"
      :model-value="modelValue"
      :options="fieldOptions"
      :has-error="hasError"
      :disabled="disabled"
      :loading="loading"
      @update:model-value="handleModelUpdate"
      @input="handleInput"
      @change="handleChange"
    />
    
    <!-- Field description -->
    <p v-if="field.description" class="text-xs text-gray-500 mt-1">
      {{ field.description }}
    </p>
    
    <!-- Field errors -->
    <div v-if="hasError" class="mt-1">
      <p
        v-for="error in errors"
        :key="error"
        class="text-sm text-red-600"
      >
        {{ error }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import DynamicFieldLoader from './fields/DynamicFieldLoader.vue'

// Props with comprehensive validation
const props = defineProps({
  field: {
    type: Object,
    required: true,
    validator: (value) => {
      // Field must have required properties
      const requiredProps = ['name', 'field_type', 'display_name', 'is_enabled']
      return requiredProps.every(prop => prop in value)
    }
  },
  modelValue: {
    type: [String, Number, Boolean, Array],
    default: ''
  },
  options: {
    type: Array,
    default: () => [],
    validator: (value) => {
      // Each option should have at least a value or be a string
      return value.every(opt => 
        typeof opt === 'string' || 
        typeof opt === 'number' ||
        (typeof opt === 'object' && opt !== null && ('value' in opt || 'code' in opt || 'id' in opt))
      )
    }
  },
  errors: {
    type: Array,
    default: () => [],
    validator: (value) => {
      // Errors should be strings
      return value.every(error => typeof error === 'string')
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

// Emits with validation
const emit = defineEmits({
  'update:modelValue': (value) => {
    // Should accept any value for v-model
    return true
  },
  input: (value) => {
    // Input event should include the new value
    return value !== undefined
  },
  change: (value) => {
    // Change event should include the new value  
    return value !== undefined
  }
})

const fieldOptions = computed(() => {
  // Use provided options or fall back to field options
  return props.options.length > 0 ? props.options : (props.field.options || [])
})

const hasError = computed(() => {
  return props.errors && props.errors.length > 0
})

const getFieldGridClass = computed(() => {
  // Wide fields that should span full width
  if (props.field.field_type === 'textarea' ||
      ['address', 'medicalHistory', 'allergies', 'medications'].includes(props.field.name)) {
    return 'md:col-span-2'
  }
  return ''
})

// Methods
const handleModelUpdate = (value) => {
  emit('update:modelValue', value)
}

const handleInput = (value) => {
  emit('input', value)
}

const handleChange = (value) => {
  emit('change', value)
}
</script>