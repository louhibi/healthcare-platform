<template>
  <div class="max-w-4xl mx-auto">

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Configuration Error</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Dynamic Form -->
    <form v-else @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Location Fields Section -->
      <FormSection v-if="locationFields.length > 0" title="Location Information">
        <LocationFieldGroup
          :fields="locationFields"
          :form-data="formData"
          :has-field-error="hasFieldError"
          :get-field-error="getFieldError"
          @location-change="onAggregatedLocationChange"
        />
      </FormSection>

      <!-- Insurance Fields Section -->
      <FormSection v-if="insuranceFields.length > 0" title="Insurance Information">
        <InsuranceFieldGroup
          :fields="insuranceFields"
          :form-data="formData"
          :has-field-error="hasFieldError"
          :get-field-error="getFieldError"
          :loading="isSubmitting"
          :country-id="formData.country_id"
          @insurance-change="handleAggregatedInsuranceChange"
        />
      </FormSection>

      <!-- Other Form Sections with Virtual Scrolling -->
      <VirtualFormSection
        v-for="section in nonLocationFieldsByCategory"
        :key="section.category"
        :title="section.category"
        :fields="orderedFields(section.fields)"
        :form-data="formData"
        :get-field-options="getSelectOptions"
        :get-field-error="getFieldError"
        :loading="isSubmitting"
        :show-field-count="isDevelopment && section.fields.length > 10"
        :min-virtual-fields="15"
        :virtual-height="section.fields.length > 30 ? '600px' : '400px'"
        @field-update="(fieldName, value) => handleFieldInput(fieldName, value)"
        @field-input="(fieldName, value) => handleFieldInput(fieldName, value)"
        @field-change="(fieldName, value) => handleFieldChange(fieldName, value)"
      />

      <!-- Form Actions -->
      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex justify-between items-center">
          <div class="text-sm text-gray-500">
            <span class="font-medium">{{ enabledFields.length }}</span> fields enabled,
            <span class="font-medium">{{ requiredFields.length }}</span> required
          </div>
          <div class="flex space-x-3">
            <button
              type="button"
              @click="() => resetForm(enabledFields)"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Reset
            </button>
            <button
              type="submit"
              :disabled="!isFormValid || isSubmitting"
              class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmitting ? 'Submitting...' : 'Submit Patient' }}
            </button>
          </div>
        </div>
      </div>
    </form>


  </div>
</template>

<script setup>
import { locationsApi } from '@/api/locations'
import { patientsApi } from '@/api/patients'
import { useFormConfig } from '@/composables/useFormConfig'
import { useFormState } from '@/composables/useFormState'
import { useValidation } from '@/composables/useValidation'
import { useEntityStore } from '@/stores/entity'
import { buildFieldProps, getFieldComponentName } from '@/utils/fieldComponentMap'
import { preloadFieldComponents } from '@/utils/fieldLoader'
import { memoizedArrayFilter } from '@/utils/memoization'
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { useToast } from 'vue-toastification'
import FormSection from '../FormSection.vue'
import InsuranceFieldGroup from '../insurance/InsuranceFieldGroup.vue'
import LocationFieldGroup from '../location/LocationFieldGroup.vue'
import VirtualFormSection from '../VirtualFormSection.vue'

// Props with comprehensive validation
const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({}),
    validator: (value) => {
      // Should be a plain object (not null, not array)
      return value !== null && typeof value === 'object' && !Array.isArray(value)
    }
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

// Emits with validation
const emit = defineEmits({
  submit: (result) => {
    // Submit event should include result data
    return result && typeof result === 'object'
  },
  cancel: () => {
    // Cancel event has no payload
    return true
  }
})

// Composables
const toast = useToast()
const validation = useValidation()
const entityStore = useEntityStore()

const {
  isInitialized,
  fields,
  enabledFields,
  requiredFields,
  fieldsByCategory,
  isLoading,
  error,
  initialize,
  getFieldConfig
} = useFormConfig('patient')

// Form state management
const {
  formData,
  validationErrors,
  isSubmitting,
  isDirty,
  hasErrors,
  isFormValid,
  initializeFormData,
  updateField,
  updateMultipleFields,
  validateField,
  validateForm,
  getFieldError,
  hasFieldError,
  clearFieldError,
  clearAllErrors,
  resetForm,
  submitForm
} = useFormState(patientsApi, {
  createMethod: 'createPatient',
  updateMethod: 'updatePatient',
  successMessage: {
    create: 'Patient created successfully!',
    update: 'Patient updated successfully!'
  },
  errorMessage: {
    create: 'Failed to create patient',
    update: 'Failed to update patient'
  }
})

// Local state
const isDevelopment = computed(() => import.meta.env.DEV)
const isInitializing = ref(false)

// Nationalities (simple direct load, no composable)
const nationalityOptions = ref([])
const loadNationalitiesImmediate = async () => {
  try {
    const data = await locationsApi.getNationalities()
    nationalityOptions.value = Array.isArray(data) ? data.map(n => ({ value: n.id, label: n.name, country_id: n.country_id, is_primary: n.is_primary })) : []
    nationalityOptions.value.sort((a,b)=>a.label.localeCompare(b.label))
  } catch (e) {
    console.error('[DynamicPatientForm] Failed to load nationalities', e)
    nationalityOptions.value = []
  }
}

// Removed insurance composable; insurance logic handled within field components directly

// Country/state/city loading removed from parent – child component manages cascading logic

// Recognizers for location fields - ID-based (support common variants)
const isCountryField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'country_id' || name === 'country' || name.endsWith('_country_id') || name.endsWith('_country')
}

const isStateField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'state_id' || name === 'state' || name === 'province_id' || name === 'province' || 
         name.endsWith('_state_id') || name.endsWith('_state') || name.endsWith('_province_id') || name.endsWith('_province') || 
         name.startsWith('state_') || name.startsWith('province_')
}

const isCityField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'city_id' || name === 'city' || name.endsWith('_city_id') || name.endsWith('_city') || name.startsWith('city_')
}

// Insurance field recognizers - precise field matching
const isInsuranceTypeField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'insurance' || name === 'insurance_type' || name === 'insurance_type_id'
}

const isInsuranceProviderField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'insurance_provider' || name === 'insurance_provider_id'
}

const isInsuranceRelatedField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'policy_number' || name === 'insurance_custom' || name === 'insurance_provider_custom'
}

const isInsuranceField = (field) => {
  return isInsuranceTypeField(field) || isInsuranceProviderField(field) || isInsuranceRelatedField(field)
}

// Removed location update methods (handled internally by child components)

// Convert value for ID-based fields
const convertIdFieldValue = (fieldName, value) => {
  // Check if this is an ID-based field
  if (fieldName.endsWith('_id')) {
    // Handle empty/null cases
    if (value === '' || value === null || value === undefined) {
      return null
    }
    
    // Convert to integer if it's a valid number
    const numValue = parseInt(value, 10)
    if (isNaN(numValue)) {
      return null
    }
    
    return numValue
  }
  
  return value
}

// Form field change handlers
const handleFieldInput = async (fieldName, value) => {
  const convertedValue = convertIdFieldValue(fieldName, value)
  updateField(fieldName, convertedValue)
  await nextTick()
  const fieldConfig = getFieldConfig(fieldName)
  validateField(fieldName, convertedValue, fieldConfig)
}

const handleFieldChange = async (fieldName, value) => {
  const convertedValue = convertIdFieldValue(fieldName, value)
  updateField(fieldName, convertedValue)
  await nextTick()
  const fieldConfig = getFieldConfig(fieldName)
  validateField(fieldName, convertedValue, fieldConfig)
}

// Removed explicit country/state change handlers (no longer needed in parent)

// Aggregated location change (from LocationFieldGroup) – update + validate
const onAggregatedLocationChange = ({ countryId, stateId, cityId }) => {
  if (countryId !== undefined && formData.value.country_id !== countryId) {
    updateField('country_id', countryId)
    const cfg = getFieldConfig('country_id')
    validateField('country_id', countryId, cfg)
  }
  if (stateId !== undefined && formData.value.state_id !== stateId) {
    updateField('state_id', stateId)
    const cfg = getFieldConfig('state_id')
    validateField('state_id', stateId, cfg)
  }
  if (cityId !== undefined && formData.value.city_id !== cityId) {
    updateField('city_id', cityId)
    const cfg = getFieldConfig('city_id')
    validateField('city_id', cityId, cfg)
  }
}

// Aggregated insurance change handler (only source of insurance updates now)
const handleAggregatedInsuranceChange = async ({ insuranceTypeId, insuranceProviderId, policyNumber }) => {
  const updates = {}
  if (insuranceTypeId !== undefined) updates.insurance_type_id = convertIdFieldValue('insurance_type_id', insuranceTypeId)
  if (insuranceProviderId !== undefined) updates.insurance_provider_id = convertIdFieldValue('insurance_provider_id', insuranceProviderId)
  if (policyNumber !== undefined) updates.policy_number = policyNumber
  Object.entries(updates).forEach(([k,v]) => updateField(k, v))
  await nextTick()
  Object.keys(updates).forEach(name => {
    const cfg = getFieldConfig(name)
    validateField(name, updates[name], cfg)
  })
}

// Utility methods
const getFieldComponent = (field) => {
  return getFieldComponentName(field.field_type)
}

// Input props per field
const getFieldProps = (field) => {
  return buildFieldProps(field)
}

// Classes per input with error state
const getFieldClasses = (field) => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
  const border = hasFieldError(field.name) ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  return `${base} ${border}`
}

// Memoized select options function
const selectOptionsCache = new Map()

const getSelectOptions = (field) => {
  // Location fields (country, state, city) are handled separately
  if (isCountryField(field) || isStateField(field) || isCityField(field)) return []
  
  // Handle nationality field with options from location API
  if (field.name === 'nationality_id') {
    return nationalityOptions.value || []
  }
  
  const cacheKey = `${field.name}-${field.field_type}-${field.options?.length || 0}`
  
  if (selectOptionsCache.has(cacheKey)) {
    return selectOptionsCache.get(cacheKey)
  }
  
  // Prefer field.options; fallback to defaults for other fields
  const options = field.options && field.options.length ? field.options : getDefaultOptions(field)
  
  // Cache result (keep only last 20 entries)
  if (selectOptionsCache.size >= 20) {
    const firstKey = selectOptionsCache.keys().next().value
    selectOptionsCache.delete(firstKey)
  }
  selectOptionsCache.set(cacheKey, options)
  
  return options
}

const getDefaultOptions = (field) => {
  // Provide default options for common select fields (excluding country now handled via API)
  const defaultOptions = {
    gender: [
      { value: 'male', label: 'Male' },
      { value: 'female', label: 'Female' },
      { value: 'other', label: 'Other' }
    ],
    blood_type: [
      { value: 'A+', label: 'A+' },
      { value: 'A-', label: 'A-' },
      { value: 'B+', label: 'B+' },
      { value: 'B-', label: 'B-' },
      { value: 'AB+', label: 'AB+' },
      { value: 'AB-', label: 'AB-' },
      { value: 'O+', label: 'O+' },
      { value: 'O-', label: 'O-' }
    ],
    marital_status: [
      { value: 'single', label: 'Single' },
      { value: 'married', label: 'Married' },
      { value: 'divorced', label: 'Divorced' },
      { value: 'widowed', label: 'Widowed' }
    ],
    preferred_language: [
      { value: 'en-MA', label: 'English (Morocco)' },
      { value: 'fr-MA', label: 'Français (Maroc)' }
    ]
  }
  
  return defaultOptions[field.name] || []
}

// Pure ID-based system - no name mapping needed

// Helper function to get initial data with entity defaults - Pure ID-based
const getInitialDataWithEntityDefaults = (initialData = {}) => {
  
  // If initialData is provided (editing mode), use it as-is
  if (initialData && Object.keys(initialData).length > 0) {
    return initialData
  }
  
  // For new patients, use entity's location IDs directly (no name mapping)
  const defaults = {}
  
  // Get IDs directly from entity store if available
  if (entityStore.entity) {
    if (entityStore.entity.country_id) {
      // Set both country_id and country for compatibility
      defaults.country_id = entityStore.entity.country_id
      defaults.country = entityStore.entity.country_id
    }
    if (entityStore.entity.state_id) {
      // Set both state_id and state for compatibility  
      defaults.state_id = entityStore.entity.state_id
      defaults.state = entityStore.entity.state_id
    }
    if (entityStore.entity.city_id) {
      // Set both city_id and city for compatibility
      defaults.city_id = entityStore.entity.city_id
      defaults.city = entityStore.entity.city_id
    }
  }
  
  return defaults
}

// Form submission handler with location validation
const handleSubmit = async () => {
  
  // Custom validation logic for location fields
  const additionalValidation = () => {
    
    if (!validateForm(enabledFields.value)) {
      toast.error('Please fix the validation errors before submitting')
      return false
    }

    // Additional healthcare-specific validation
    const patientErrors = validation.validatePatientData(formData.value)
    Object.entries(patientErrors).forEach(([field, errors]) => {
      validation.setFieldError(field, errors)
    })

    if (Object.keys(patientErrors).length > 0) {
      toast.error('Please correct the patient data errors')
      return false
    }

    return true
  }

  try {
    const result = await submitForm(
      props.isEdit,
      props.initialData?.id,
      additionalValidation
    )
    
    if (result) {
      emit('submit', result)
    }
  } catch (err) {
    console.error('Form submission error:', err)
  }
}

// Form initialization
const initializeForm = async () => {
  try {
    isInitializing.value = true
    
    await initialize()
    
  // Load nationalities only (countries/states/cities handled in child components)
  await loadNationalitiesImmediate()
  await entityStore.ensureNationalitiesLoaded()
    
    // Wait for next tick to ensure computed properties are updated
    await nextTick()
    
    
    // Prepare initial data with entity defaults when no initial data is provided
    const initialDataWithDefaults = getInitialDataWithEntityDefaults(props.initialData)
    
    // Use the unified useFormState solution with entity store for defaults
    initializeFormData(initialDataWithDefaults, enabledFields.value, entityStore)
    
    // Small delay to ensure all reactive updates are complete
    await nextTick()
    
    // Preload commonly used field components for better performance
    const commonFieldTypes = ['select', 'date', 'textarea', 'checkbox', 'number']
    const enabledFieldTypes = enabledFields.value.map(f => f.field_type)
    const preloadTypes = [...new Set([...commonFieldTypes, ...enabledFieldTypes])]
    
    // Field component loaders map
    const fieldLoaders = {
      date: () => import('../fields/DateField.vue'),
      datetime: () => import('../fields/DateField.vue'),
      time: () => import('../fields/DateField.vue'),
      textarea: () => import('../fields/TextareaField.vue'),
      select: () => import('../fields/SelectField.vue'),
      multiselect: () => import('../fields/SelectField.vue'),
      number: () => import('../fields/NumberField.vue'),
      checkbox: () => import('../fields/CheckboxField.vue'),
      boolean: () => import('../fields/CheckboxField.vue'),
      file: () => import('../fields/FileField.vue')
    }
    
    // Start preloading in background (don't await)
    preloadFieldComponents(preloadTypes, fieldLoaders).catch(console.warn)
    
  // No country/state/city preloading here anymore
  } catch (err) {
    console.error('💥 Form initialization error:', err)
    toast.error('Failed to load form configuration')
  } finally {
    isInitializing.value = false
  }
}

// Removed watcher for country selection – cascading handled in child component

watch(() => props.initialData, async (newData, oldData) => {
  // Reinitialize whenever we get new meaningful data (even during initial mount)
  if (newData && Object.keys(newData).length > 0) {
    // Wait for form config to be ready if it's not initialized yet
    if (!isInitialized.value) {
      let attempts = 0
      while (!isInitialized.value && attempts < 20) {
        await new Promise(resolve => setTimeout(resolve, 100))
        attempts++
      }
    }
    
    if (isInitialized.value) {
      isInitializing.value = true
      
      // Prepare data with entity defaults - ID-based
      const dataWithDefaults = getInitialDataWithEntityDefaults(newData)
      initializeFormData(dataWithDefaults, enabledFields.value, entityStore)
      
      // Force reactive update
      await nextTick()
  // No immediate loading of states/cities here; delegated to child
      isInitializing.value = false
    }
  }
}, { deep: true, immediate: true })

// Initialize on mount
initializeForm()

// Layout helper for wide fields
function getFieldGridClass(field) {
  if (field.field_type === 'textarea' ||
      ['address', 'medicalHistory', 'allergies', 'medications'].includes(field.name)) {
    return 'md:col-span-2'
  }
  return ''
}

// Memoized field ordering function
const fieldOrderCache = new Map()

function orderedFields(list) {
  if (!Array.isArray(list)) return []
  
  // Create cache key based on field names and sort orders
  const cacheKey = list.map(f => `${f.name}-${f.sort_order}`).join('|')
  
  if (fieldOrderCache.has(cacheKey)) {
    return fieldOrderCache.get(cacheKey)
  }
  
  const arr = [...list]
  arr.sort((a, b) => {
    const aCountry = isCountryField(a)
    const bCountry = isCountryField(b)
    const aState = isStateField(a)
    const bState = isStateField(b)
    const aCity = isCityField(a)
    const bCity = isCityField(b)
    
    // Location field hierarchy: country < state < city
    if (aCountry && (bState || bCity)) return -1
    if ((aState || aCity) && bCountry) return 1
    if (aState && bCity) return -1
    if (aCity && bState) return 1
    
    // fallback to provided sort_order if present
    const ao = a.sort_order ?? 0
    const bo = b.sort_order ?? 0
    if (ao !== bo) return ao - bo
    return 0
  })
  
  // Cache result (keep only last 10 entries)
  if (fieldOrderCache.size >= 10) {
    const firstKey = fieldOrderCache.keys().next().value
    fieldOrderCache.delete(firstKey)
  }
  fieldOrderCache.set(cacheKey, arr)
  
  return arr
}

// Separate location and non-location fields (memoized)
const locationFields = memoizedArrayFilter(
  enabledFields,
  (field) => isCountryField(field) || isStateField(field) || isCityField(field),
  []
)

// Separate insurance fields (memoized)
const insuranceFields = memoizedArrayFilter(
  enabledFields,
  (field) => isInsuranceField(field),
  []
)

// Memoized non-location, non-insurance fields by category with sorting
const nonLocationFieldsByCategory = computed(() => {
  const entries = Object.entries(fieldsByCategory.value || {})
  const filteredEntries = entries.map(([category, fields]) => {
    const enabled = (fields || []).filter(f => 
      f.is_enabled && !isCountryField(f) && !isStateField(f) && !isCityField(f) && !isInsuranceField(f)
    )
    
    return { category, fields: enabled }
  }).filter(section => section.fields.length > 0)
  
  // Sort sections by minimum field order
  filteredEntries.sort((a, b) => {
    const aMinOrder = a.fields.length ? Math.min(...a.fields.map(f => f.sort_order ?? 0)) : 0
    const bMinOrder = b.fields.length ? Math.min(...b.fields.map(f => f.sort_order ?? 0)) : 0
    if (aMinOrder !== bMinOrder) return aMinOrder - bMinOrder
    return a.category.localeCompare(b.category)
  })
  
  return filteredEntries
})

// Sort categories for ergonomic ordering (location fields in proper order)
const categoriesSorted = computed(() => {
  const entries = Object.entries(fieldsByCategory.value || {})
  const sections = entries.map(([category, fields]) => {
    const enabled = (fields || []).filter(f => f.is_enabled)
    const minOrder = enabled.length ? Math.min(...enabled.map(f => f.sort_order ?? 0)) : 0
    const hasCountry = enabled.some(isCountryField)
    const hasState = enabled.some(isStateField)
    const hasCity = enabled.some(isCityField)
    return { category, fields, minOrder, hasCountry, hasState, hasCity }
  })
  sections.sort((a, b) => {
    // Location sections should be ordered: country, state, city
    if (a.hasCountry && (b.hasState || b.hasCity)) return -1
    if ((a.hasState || a.hasCity) && b.hasCountry) return 1
    if (a.hasState && b.hasCity) return -1
    if (a.hasCity && b.hasState) return 1
    
    // Fallback to sort order and name
    if (a.minOrder !== b.minOrder) return a.minOrder - b.minOrder
    return a.category.localeCompare(b.category)
  })
  return sections
})

// Cleanup on unmount
onUnmounted(() => {
  fieldOrderCache.clear()
  selectOptionsCache.clear()
})
</script>
