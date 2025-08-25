<template>
  <textarea
    v-bind="$attrs"
    :value="modelValue"
    :class="fieldClasses"
    :rows="rows"
    @input="handleInput"
    @change="handleChange"
  ></textarea>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: String,
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
  rows: {
    type: Number,
    default: 3
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => typeof value === 'string',
  'input': (value) => typeof value === 'string',
  'change': (value) => typeof value === 'string'
})

// Computed
const fieldClasses = computed(() => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm resize-y'
  const border = props.hasError ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  const disabledClass = (props.disabled || props.loading) ? 'bg-gray-50 text-gray-500' : ''
  return `${base} ${border} ${disabledClass}`
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