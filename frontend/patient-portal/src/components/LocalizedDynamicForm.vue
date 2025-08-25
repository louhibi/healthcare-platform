<template>
  <div class="max-w-4xl mx-auto">
    <!-- Form Header -->
    <div class="bg-white shadow rounded-lg p-6 mb-6">
      <div class="flex justify-between items-start">
        <div>
          <h2 class="text-xl font-bold text-gray-900 mb-2">
            {{ localizedMetadata?.display_name || t(`forms.${formType}.title`) }}
          </h2>
          <p class="text-sm text-gray-600">
            {{ localizedMetadata?.description || t(`forms.${formType}.description`) }}
          </p>
        </div>
        <LanguageSelector />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      <span class="ml-3 text-gray-600">{{ t('common.loading') }}</span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
      <div class="flex">
        <div class="flex-shrink-0">
          <ExclamationTriangleIcon class="h-5 w-5 text-red-400" />
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">{{ t('errors.general') }}</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Dynamic Form -->
    <form v-else @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Render form sections by category -->
      <div
        v-for="(categoryFields, category) in fieldsByCategory"
        :key="category"
        class="bg-white shadow rounded-lg p-6"
      >
        <h3 class="text-lg font-medium text-gray-900 mb-4 border-b border-gray-200 pb-2">
          {{ getCategoryDisplayName(category) }}
        </h3>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="field in categoryFields.filter(f => f.is_enabled)"
            :key="field.field_id"
            :class="getFieldGridClass(field)"
          >
            <!-- Field Label -->
            <label :for="field.name" class="block text-sm font-medium text-gray-700 mb-1">
              {{ field.display_name }}
              <span v-if="field.is_required" class="text-red-500">*</span>
              <span v-else class="text-gray-400 text-xs ml-1">({{ t('common.optional') }})</span>
            </label>

            <!-- Field Input -->
            <component
              :is="getFieldComponent(field)"
              :id="field.name"
              v-model="formData[field.name]"
              :placeholder="field.placeholder_text"
              :required="field.is_required"
              :disabled="isSubmitting"
              v-bind="getFieldProps(field)"
              @blur="validateField(field.name)"
              @input="handleFieldInput(field.name, $event)"
              :class="[
                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                {
                  'border-red-300 focus:border-red-500 focus:ring-red-500': hasFieldError(field.name),
                  'bg-gray-50': isSubmitting
                }
              ]"
            >
              <!-- Options for select fields -->
              <option v-if="field.field_type === 'select'" value="">
                {{ field.placeholder_text || t('common.select') }}
              </option>
              <option
                v-for="option in getFieldOptions(field)"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </component>

            <!-- Field Description -->
            <p v-if="field.description" class="mt-1 text-xs text-gray-500">
              {{ field.description }}
            </p>

            <!-- Field Error -->
            <p v-if="hasFieldError(field.name)" class="mt-1 text-sm text-red-600">
              {{ getFieldError(field.name) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex justify-end space-x-3">
          <button
            type="button"
            @click="handleCancel"
            :disabled="isSubmitting"
            class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            {{ t('common.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <div v-if="isSubmitting" class="flex items-center">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              {{ t('common.loading') }}
            </div>
            <span v-else>
              {{ isEdit ? t('common.save') : t('common.submit') }}
            </span>
          </button>
        </div>
      </div>

      <!-- Development Debug Panel -->
      <div v-if="isDevelopment" class="bg-gray-50 border border-gray-200 rounded-lg p-4 mt-6">
        <h4 class="text-sm font-medium text-gray-900 mb-2">Debug Info (Development Only)</h4>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div>
            <h5 class="font-medium text-gray-700">Current Locale:</h5>
            <p class="text-gray-600">{{ currentLocale }}</p>
          </div>
          <div>
            <h5 class="font-medium text-gray-700">Form Fields:</h5>
            <p class="text-gray-600">{{ enabledFields.length }} enabled, {{ requiredFields.length }} required</p>
          </div>
          <div>
            <h5 class="font-medium text-gray-700">Validation Errors:</h5>
            <p class="text-gray-600">{{ Object.keys(validationErrors).length }} errors</p>
          </div>
          <div>
            <h5 class="font-medium text-gray-700">Form Valid:</h5>
            <p :class="isFormValid ? 'text-green-600' : 'text-red-600'">
              {{ isFormValid ? 'Yes' : 'No' }}
            </p>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n as useVueI18n } from 'vue-i18n'
import { useToast } from 'vue-toastification'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useI18n } from '../composables/useI18n'
import LanguageSelector from './LanguageSelector.vue'

// Props
const props = defineProps({
  formType: {
    type: String,
    required: true,
    validator: (value) => ['patient', 'appointment'].includes(value)
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['submit', 'cancel'])

// Composables
const { t } = useVueI18n()
const toast = useToast()
const { 
  currentLocale, 
  getLocalizedFormMetadata,
  getLocalizedFieldName,
  isRTL 
} = useI18n()

// Local state
const isLoading = ref(true)
const isSubmitting = ref(false)
const error = ref(null)
const localizedMetadata = ref(null)
const formData = ref({})
const validationErrors = ref({})
const isDevelopment = computed(() => import.meta.env.DEV)

// Computed properties
const fields = computed(() => localizedMetadata.value?.fields || [])
const enabledFields = computed(() => fields.value.filter(f => f.is_enabled))
const requiredFields = computed(() => fields.value.filter(f => f.is_required && f.is_enabled))

const fieldsByCategory = computed(() => {
  const categories = {}
  enabledFields.value.forEach(field => {
    const category = field.category || 'General'
    if (!categories[category]) {
      categories[category] = []
    }
    categories[category].push(field)
  })
  
  // Sort fields within each category by sort_order
  Object.keys(categories).forEach(category => {
    categories[category].sort((a, b) => (a.sort_order || 0) - (b.sort_order || 0))
  })
  
  return categories
})

const isFormValid = computed(() => {
  // Check if all required fields have values and there are no validation errors
  const hasRequiredFieldValues = requiredFields.value.every(field => {
    const value = formData.value[field.name]
    return value !== null && value !== undefined && value !== ''
  })
  
  const hasNoValidationErrors = Object.keys(validationErrors.value).length === 0
  
  return hasRequiredFieldValues && hasNoValidationErrors
})

// Methods
async function initialize() {
  try {
    isLoading.value = true
    error.value = null
    
    // Get localized form metadata
    const metadata = await getLocalizedFormMetadata(props.formType)
    localizedMetadata.value = metadata
    
    // Initialize form data
    initializeFormData()
    
  } catch (err) {
    console.error('Failed to initialize localized form:', err)
    error.value = err.message || t('errors.general')
  } finally {
    isLoading.value = false
  }
}

function initializeFormData() {
  const data = { ...props.initialData }
  
  // Ensure all enabled fields have a value (null, empty string, or array)
  enabledFields.value.forEach(field => {
    if (!(field.name in data)) {
      if (field.field_type === 'multiselect') {
        data[field.name] = []
      } else {
        data[field.name] = ''
      }
    }
  })
  
  formData.value = data
}

function getFieldComponent(field) {
  switch (field.field_type) {
    case 'textarea':
      return 'textarea'
    case 'select':
    case 'multiselect':
      return 'select'
    default:
      return 'input'
  }
}

function getFieldProps(field) {
  const props = {
    type: field.field_type === 'email' ? 'email' : 
          field.field_type === 'tel' ? 'tel' :
          field.field_type === 'date' ? 'date' :
          field.field_type === 'number' ? 'number' : 'text'
  }
  
  // Add validation attributes
  if (field.validation_rules) {
    if (field.validation_rules.min_length) {
      props.minlength = field.validation_rules.min_length
    }
    if (field.validation_rules.max_length) {
      props.maxlength = field.validation_rules.max_length
    }
    if (field.validation_rules.pattern) {
      props.pattern = field.validation_rules.pattern
    }
  }
  
  return props
}

function getFieldOptions(field) {
  if (!field.options || !Array.isArray(field.options)) {
    return []
  }
  
  return field.options.map(option => {
    if (typeof option === 'string') {
      return { value: option, label: option }
    }
    return option
  })
}

function getFieldGridClass(field) {
  // Full width for textarea and certain field types
  if (field.field_type === 'textarea' || 
      field.name === 'address' || 
      field.name === 'medicalHistory' ||
      field.name === 'allergies' ||
      field.name === 'medications') {
    return 'md:col-span-2'
  }
  return ''
}

function getCategoryDisplayName(category) {
  // Try to get localized category name
  const key = `forms.${props.formType}.categories.${category.toLowerCase()}`
  const translated = t(key)
  return translated !== key ? translated : category
}

function handleFieldInput(fieldName, event) {
  // Clear validation error when user starts typing
  if (validationErrors.value[fieldName]) {
    delete validationErrors.value[fieldName]
  }
  
  // Update form data
  const value = event.target ? event.target.value : event
  formData.value[fieldName] = value
}

function validateField(fieldName) {
  const field = fields.value.find(f => f.name === fieldName)
  if (!field || !field.is_enabled) return
  
  const value = formData.value[fieldName]
  const errors = []
  
  // Required field validation
  if (field.is_required && (!value || value.toString().trim() === '')) {
    errors.push(t('validation.required'))
  }
  
  // Type-specific validation
  if (value && value.toString().trim() !== '') {
    switch (field.field_type) {
      case 'email':
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
          errors.push(t('validation.email'))
        }
        break
      case 'tel':
        // Basic phone validation
        if (!/^[\d\s\-\+\(\)]+$/.test(value)) {
          errors.push(t('validation.phone'))
        }
        break
      case 'number':
        if (isNaN(value)) {
          errors.push(t('validation.number'))
        }
        break
    }
    
    // Custom validation rules
    if (field.validation_rules) {
      if (field.validation_rules.min_length && value.length < field.validation_rules.min_length) {
        errors.push(t('validation.minLength', { min: field.validation_rules.min_length }))
      }
      if (field.validation_rules.max_length && value.length > field.validation_rules.max_length) {
        errors.push(t('validation.maxLength', { max: field.validation_rules.max_length }))
      }
      if (field.validation_rules.pattern) {
        const regex = new RegExp(field.validation_rules.pattern)
        if (!regex.test(value)) {
          errors.push(t('validation.pattern'))
        }
      }
    }
  }
  
  // Update validation errors
  if (errors.length > 0) {
    validationErrors.value[fieldName] = errors[0] // Show first error
  } else {
    delete validationErrors.value[fieldName]
  }
}

function validateForm() {
  // Validate all enabled fields
  enabledFields.value.forEach(field => {
    validateField(field.name)
  })
  
  return Object.keys(validationErrors.value).length === 0
}

function hasFieldError(fieldName) {
  return !!validationErrors.value[fieldName]
}

function getFieldError(fieldName) {
  return validationErrors.value[fieldName]
}

function clearAllErrors() {
  validationErrors.value = {}
}

async function handleSubmit() {
  if (!validateForm()) {
    toast.error(t('errors.validation'))
    return
  }
  
  try {
    isSubmitting.value = true
    clearAllErrors()
    
    // Submit the form data
    emit('submit', { ...formData.value })
    
  } catch (err) {
    console.error('Form submission error:', err)
    toast.error(err.message || t('errors.general'))
  } finally {
    isSubmitting.value = false
  }
}

function handleCancel() {
  emit('cancel')
}

// Watch for locale changes and reinitialize
watch(currentLocale, () => {
  initialize()
})

// Initialize on mount
onMounted(() => {
  initialize()
})
</script>

<style scoped>
/* RTL support */
html[dir="rtl"] .space-x-3 > * + * {
  margin-left: 0;
  margin-right: 0.75rem;
}

html[dir="rtl"] .ml-3 {
  margin-left: 0;
  margin-right: 0.75rem;
}

html[dir="rtl"] .mr-2 {
  margin-left: 0.5rem;
  margin-right: 0;
}

/* Form field focus states */
.form-field:focus-within {
  transform: translateY(-1px);
  transition: transform 0.2s ease;
}

/* Loading animation */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>