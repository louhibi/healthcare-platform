<template>
  <div class="relative flex items-start">
    <div class="flex items-center h-5">
      <input
        v-bind="$attrs"
        :checked="isChecked"
        type="checkbox"
        :class="fieldClasses"
        @input="handleInput"
        @change="handleChange"
      />
    </div>
    <div class="ml-3 text-sm">
      <label :for="$attrs.id" class="font-medium text-gray-700">
        <slot>{{ label }}</slot>
      </label>
      <p v-if="description" class="text-gray-500">
        {{ description }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: [Boolean, String, Number, Array],
    default: false
  },
  trueValue: {
    type: [Boolean, String, Number],
    default: true
  },
  falseValue: {
    type: [Boolean, String, Number],
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  description: {
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
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => true,
  'input': (value) => true,
  'change': (value) => true
})

// Computed
const isChecked = computed(() => {
  if (Array.isArray(props.modelValue)) {
    return props.modelValue.includes(props.trueValue)
  }
  return props.modelValue === props.trueValue
})

const fieldClasses = computed(() => {
  const base = 'focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300 rounded'
  const errorClass = props.hasError ? 'border-red-300 focus:ring-red-500' : ''
  const disabledClass = (props.disabled || props.loading) ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
  return `${base} ${errorClass} ${disabledClass}`
})

// Methods
const handleInput = (event) => {
  const isChecked = event.target.checked
  let newValue
  
  if (Array.isArray(props.modelValue)) {
    // Handle array of values (for checkbox groups)
    newValue = [...props.modelValue]
    if (isChecked) {
      if (!newValue.includes(props.trueValue)) {
        newValue.push(props.trueValue)
      }
    } else {
      const index = newValue.indexOf(props.trueValue)
      if (index > -1) {
        newValue.splice(index, 1)
      }
    }
  } else {
    // Handle single boolean value
    newValue = isChecked ? props.trueValue : props.falseValue
  }
  
  emit('update:modelValue', newValue)
  emit('input', newValue)
}

const handleChange = (event) => {
  handleInput(event)
  emit('change', props.modelValue)
}
</script>