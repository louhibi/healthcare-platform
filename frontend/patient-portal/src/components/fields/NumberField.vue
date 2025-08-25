<template>
  <input
    v-bind="$attrs"
    :value="modelValue"
    type="number"
    :class="fieldClasses"
    :step="step"
    :min="min"
    :max="max"
    @input="handleInput"
    @change="handleChange"
    @wheel="handleWheel"
  />
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
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
  },
  min: {
    type: [String, Number],
    default: undefined
  },
  max: {
    type: [String, Number],
    default: undefined
  },
  step: {
    type: [String, Number],
    default: undefined
  },
  precision: {
    type: Number,
    default: 0,
    validator: (value) => value >= 0 && value <= 10
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

// Methods
const handleInput = (event) => {
  let value = event.target.value
  
  // Convert to number if it's a valid number
  if (value !== '' && !isNaN(value)) {
    value = props.precision > 0 ? parseFloat(value) : parseInt(value, 10)
  }
  
  emit('update:modelValue', value)
  emit('input', value)
}

const handleChange = (event) => {
  let value = event.target.value
  
  // Convert to number if it's a valid number
  if (value !== '' && !isNaN(value)) {
    value = props.precision > 0 ? parseFloat(value) : parseInt(value, 10)
  }
  
  emit('update:modelValue', value)
  emit('change', value)
}

// Prevent accidental scrolling changes in number inputs
const handleWheel = (event) => {
  if (event.target === document.activeElement) {
    event.preventDefault()
  }
}
</script>