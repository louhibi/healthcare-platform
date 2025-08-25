<template>
  <div>
    <!-- Vue Datepicker with proper theming -->
    <VueDatePicker
      :model-value="isoValue"
      @update:model-value="handleDateChange"
      :placeholder="placeholder"
      :disabled="disabled || loading"
      :class="themeClasses"
      :format="vueDatePickerFormat"
      :locale="currentLocaleCode"
      :enable-time-picker="false"
      :clearable="!required"
      :required="required"
      :auto-apply="true"
    />
    
    <!-- Validation errors -->
    <div v-if="internalValidationErrors.length > 0 && !hasError" class="mt-1">
      <p
        v-for="error in internalValidationErrors"
        :key="error"
        class="text-sm text-red-600"
      >
        {{ error }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import VueDatePicker from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import { useDateLocalization } from '@/composables/useDateLocalization'

// Props
const props = defineProps({
  modelValue: {
    type: [String, Date],
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
  required: {
    type: Boolean,
    default: false
  },
  showFormatHint: {
    type: Boolean,
    default: true
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => true,
  'input': (value) => true,
  'change': (value) => true,
  'validation-error': (errors) => true
})

// Composition API
const { t, locale } = useI18n()
const {
  validateDateFormat,
  placeholder,
  example,
  displayFormat
} = useDateLocalization()

// Refs
const internalValidationErrors = ref([])

// Computed values
const currentLocaleCode = computed(() => {
  const localeMap = {
    'en-MA': 'en',
    'fr-MA': 'fr', 
    'ar-MA': 'ar',
    'en-CA': 'en-CA',
    'fr-CA': 'fr-CA',
    'en-US': 'en-US',
    'fr-FR': 'fr'
  }
  return localeMap[locale.value] || 'en'
})

const isoValue = computed(() => {
  if (!props.modelValue) return null
  
  if (typeof props.modelValue === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(props.modelValue)) {
    return new Date(props.modelValue + 'T00:00:00')
  }
  
  if (props.modelValue instanceof Date) {
    return props.modelValue
  }
  
  try {
    const date = new Date(props.modelValue)
    return isNaN(date) ? null : date
  } catch {
    return null
  }
})

const hasValidationError = computed(() => {
  return props.hasError || internalValidationErrors.value.length > 0
})

const themeClasses = computed(() => {
  return hasValidationError.value ? 'healthcare-datepicker healthcare-datepicker-error' : 'healthcare-datepicker'
})

const vueDatePickerFormat = computed(() => {
  const format = displayFormat.value
  if (format === 'DD/MM/YYYY') {
    return 'dd/MM/yyyy'
  } else if (format === 'MM/DD/YYYY') {
    return 'MM/dd/yyyy'
  }
  return 'dd/MM/yyyy'
})

// Methods
const validateInput = (value) => {
  internalValidationErrors.value = []
  
  if (!value && props.required) {
    internalValidationErrors.value = [t('forms.validation.required')]
    return false
  }
  
  return true
}

const formatDateToString = (date) => {
  if (!date) return ''
  if (typeof date === 'string') return date
  if (date instanceof Date) {
    return date.toISOString().split('T')[0]
  }
  return ''
}

const handleDateChange = (date) => {
  const isoString = formatDateToString(date)
  
  if (validateInput(date)) {
    emit('update:modelValue', isoString)
    emit('change', isoString)
    emit('input', isoString)
    internalValidationErrors.value = []
  }
}

// Watch for external validation errors
watch(() => props.hasError, (hasError) => {
  if (hasError) {
    internalValidationErrors.value = []
  }
})
</script>

<style>
/* Proper Vue3 DatePicker theming using official approach */
:root {
  --dp-font-family: Inter, system-ui, sans-serif;
  --dp-border-radius: 0px;
  --dp-input-padding: 0.5rem 0.75rem;
  --dp-font-size: 1rem;
}

.healthcare-datepicker {
  /* Healthcare portal theme */
  --dp-background-color: #ffffff;
  --dp-text-color: rgb(17, 24, 39);
  --dp-border-color: #6b7280;
  --dp-border-color-hover: #9ca3af;
  --dp-border-color-focus: #4f46e5;
  --dp-primary-color: #4f46e5;
  --dp-primary-text-color: #ffffff;
  --dp-secondary-color: #f3f4f6;
  --dp-success-color: #10b981;
  --dp-danger-color: #ef4444;
  --dp-warning-color: #f59e0b;
  
  /* Input specific styling to match your forms */
  --dp-input-border-color: #6b7280;
  --dp-input-border-color-hover: #9ca3af;
  --dp-input-border-color-focus: #4f46e5;
  --dp-input-background-color: #ffffff;
  --dp-input-text-color: rgb(17, 24, 39);
  
  /* Remove shadows and set border radius to 0 */
  --dp-input-box-shadow: 0 0 #0000;
  --dp-input-box-shadow-focus: 0 0 0 3px rgba(79, 70, 229, 0.1);
  --dp-border-radius: 0px;
  
  /* Menu styling */
  --dp-menu-border-color: #d1d5db;
  --dp-menu-background-color: #ffffff;
  --dp-menu-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

/* Error state theme */
.healthcare-datepicker-error {
  --dp-border-color: #ef4444;
  --dp-border-color-hover: #dc2626;
  --dp-border-color-focus: #ef4444;
  --dp-input-border-color: #ef4444;
  --dp-input-border-color-hover: #dc2626;
  --dp-input-border-color-focus: #ef4444;
  --dp-input-box-shadow-focus: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

/* Clean integration with existing forms */
:deep(.dp__input_wrap) {
  margin: 0;
  padding: 0;
  border: none !important; /* Remove any inherited border */
  background: transparent !important;
}

:deep(.dp__main) {
  margin: 0;
  padding: 0;
  border: none !important; /* Remove any inherited border */
  background: transparent !important;
}

:deep(.dp__input) {
  border-style: solid;
  border-width: 1px;
  appearance: none;
  box-sizing: border-box;
  width: 100%;
  font-family: Inter, system-ui, sans-serif;
  font-size: 1rem;
  line-height: 1.5rem;
}

/* Target the parent wrapper that has form field attributes */
div[id][name][type="date"] {
  border: none !important;
  padding: 0 !important;
  margin: 0 !important;
  background: transparent !important;
}

/* Menu positioning */
:deep(.dp__menu) {
  z-index: 9999;
}
</style>