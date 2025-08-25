<template>
  <select
    v-bind="$attrs"
    :value="modelValue"
    :class="fieldClasses"
    @input="handleInput"
    @change="handleChange"
  >
    <!-- Placeholder option -->
    <option value="">
      {{ placeholder }}
    </option>
    
    <!-- Options -->
    <option
      v-for="option in normalizedOptions"
      :key="option.key"
      :value="option.value"
      :disabled="option.disabled"
    >
      {{ option.label }}
    </option>
  </select>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    default: () => [],
    validator: (value) => {
      return value.every(opt => 
        typeof opt === 'string' || 
        typeof opt === 'number' ||
        (typeof opt === 'object' && opt !== null)
      )
    }
  },
  placeholder: {
    type: String,
    default: 'Select an option'
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

// Computed
const fieldClasses = computed(() => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
  const border = props.hasError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  const disabledClass = (props.disabled || props.loading) ? 'bg-gray-50 text-gray-500' : ''
  return `${base} ${border} ${disabledClass}`
})

const normalizedOptions = computed(() => {
  return props.options.map(opt => {
    // Handle different option formats
    if (typeof opt === 'string' || typeof opt === 'number') {
      return {
        key: opt,
        value: opt,
        label: String(opt),
        disabled: false
      }
    }
    
    // Handle object options
    const value = opt.value ?? opt.code ?? opt.id ?? opt
    const label = opt.label ?? opt.name ?? opt.display_name ?? String(value)
    const key = `${value}-${label}`
    
    return {
      key,
      value,
      label,
      disabled: !!opt.disabled
    }
  })
})

// Methods
const handleInput = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  emit('input', value)
}

const handleChange = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  emit('change', value)
}
</script>