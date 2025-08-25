<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- Country field -->
    <div v-if="countryField" :class="getFieldGridClass(countryField)">
      <CountrySelect
        :id="countryField.name"
        :label="countryField.display_name"
        v-model="formData[countryField.name]"
        :required="countryField.is_required"
        :help="countryField.description"
        :error="hasFieldError(countryField.name) ? getFieldError(countryField.name)[0] : ''"
        @country-change="val => handleCountryChange(countryField.name, val)"
      />
    </div>

    <!-- State field -->
    <div v-if="stateField" :class="getFieldGridClass(stateField)">
      <StateSelect
        :id="stateField.name"
        :label="stateField.display_name"
        v-model="formData[stateField.name]"
        :required="stateField.is_required"
        :help="stateField.description"
        :error="hasFieldError(stateField.name) ? getFieldError(stateField.name)[0] : ''"
        :country-id="selectedCountryId"
        @state-change="val => handleStateChange(stateField.name, val)"
      />
      <div v-if="hasFieldError(stateField.name)" class="mt-1">
        <p
          v-for="error in getFieldError(stateField.name)"
          :key="error"
          class="text-sm text-red-600"
        >
          {{ error }}
        </p>
      </div>
    </div>

    <!-- City field -->
    <div v-if="cityField" :class="getFieldGridClass(cityField)">
      <CitySelect
        :id="cityField.name"
        :label="cityField.display_name"
        v-model="formData[cityField.name]"
        :required="cityField.is_required"
        :help="cityField.description"
        :error="hasFieldError(cityField.name) ? getFieldError(cityField.name)[0] : ''"
        :state-id="formData[stateField?.name] || null"
        @city-change="val => handleCityChange(cityField.name, val)"
      />
      <div v-if="hasFieldError(cityField.name)" class="mt-1">
        <p
          v-for="error in getFieldError(cityField.name)"
          :key="error"
          class="text-sm text-red-600"
        >
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import logger from '@/utils/logger'
import { computed, watch } from 'vue'
import CitySelect from './CitySelect.vue'
import CountrySelect from './CountrySelect.vue'
import StateSelect from './StateSelect.vue'

// Props (data lists removed - child selects load their own)
const props = defineProps({
  fields: {
    type: Array,
    required: true,
    validator: (value) => Array.isArray(value) && value.every(f => f && typeof f === 'object' && 'name' in f && 'field_type' in f)
  },
  formData: {
    type: Object,
    required: true,
    validator: (value) => value !== null && typeof value === 'object'
  },
  hasFieldError: { type: Function, required: true },
  getFieldError: { type: Function, required: true }
})

logger.debug('[LocationFieldGroup] initial props', {
  fields: props.fields?.map(f => f.name),
  formDataKeys: Object.keys(props.formData || {})
})

// Emits
const emit = defineEmits({
  'country-change': (fieldName, value) => typeof fieldName === 'string' && value !== undefined,
  'state-change': (fieldName, value) => typeof fieldName === 'string' && value !== undefined,
  'city-change': (fieldName, value) => typeof fieldName === 'string' && value !== undefined,
  'location-change': (payload) => payload && typeof payload === 'object'
})

// Computed properties
const countryField = computed(() => {
  const f = props.fields.find(field => isCountryField(field))
  if (f) logger.debug('[LocationFieldGroup] countryField detected', f.name)
  return f
})

const stateField = computed(() => {
  const f = props.fields.find(field => isStateField(field))
  if (f) logger.debug('[LocationFieldGroup] stateField detected', f.name)
  return f
})

const cityField = computed(() => {
  const f = props.fields.find(field => isCityField(field))
  if (f) logger.debug('[LocationFieldGroup] cityField detected', f.name)
  return f
})

const selectedCountryId = computed(() => {
  const raw = props.formData?.country_id ?? props.formData?.country ?? null
  if (raw == null || raw === '') return null
  const num = Number(raw)
  const resolved = isNaN(num) ? null : num
  logger.debug('[LocationFieldGroup] selectedCountryId computed', { raw, resolved })
  return resolved
})

// Helper factory and matchers for field recognition (exact name matches)
const createExactNameMatcher = (...names) => {
  const lowered = names.map(n => n.toLowerCase())
  return (field) => lowered.includes((field?.name || '').toLowerCase())
}
const isCountryField = createExactNameMatcher('country_id')
const isStateField   = createExactNameMatcher('state_id')
const isCityField    = createExactNameMatcher('city_id')

// Styling functions
const getFieldGridClass = (field) => {
  if (field.field_type === 'textarea' ||
      ['address', 'medicalHistory', 'allergies', 'medications'].includes(field.name)) {
    return 'md:col-span-2'
  }
  return ''
}

const getFieldClasses = (field) => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
  const border = props.hasFieldError(field.name) ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  return `${base} ${border}`
}

// Event handlers
const emitAggregatedLocation = () => {
  const countryId = selectedCountryId.value
  const stateId = stateField.value ? (props.formData[stateField.value.name] ?? null) : null
  const cityId = cityField.value ? (props.formData[cityField.value.name] ?? null) : null
  emit('location-change', { countryId, stateId, cityId })
  logger.debug('[LocationFieldGroup] location-change emitted', { countryId, stateId, cityId })
}

const handleCountryChange = (fieldName, value) => {
  logger.debug('[LocationFieldGroup] handleCountryChange', { fieldName, value })
  emit('country-change', fieldName, value)
  emitAggregatedLocation()
}

const handleStateChange = (fieldName, value) => {
  logger.debug('[LocationFieldGroup] handleStateChange', { fieldName, value })
  emit('state-change', fieldName, value)
  emitAggregatedLocation()
}

const handleCityChange = (fieldName, value) => {
  logger.debug('[LocationFieldGroup] handleCityChange', { fieldName, value })
  emit('city-change', fieldName, value)
  emitAggregatedLocation()
}

// --- Child component prop debug logging ---
// Build prop snapshots similar to what template binds, so we can log them when they change
const countrySelectProps = computed(() => {
  if (!countryField.value) return null
  const field = countryField.value
  return {
    id: field.name,
    label: field.display_name,
    modelValue: props.formData[field.name],
    required: field.is_required,
    help: field.description,
    error: props.hasFieldError(field.name) ? props.getFieldError(field.name)[0] : ''
  }
})

const stateSelectProps = computed(() => {
  if (!stateField.value) return null
  const field = stateField.value
  return {
    id: field.name,
    label: field.display_name,
    modelValue: props.formData[field.name],
    required: field.is_required,
    help: field.description,
    error: props.hasFieldError(field.name) ? props.getFieldError(field.name)[0] : '',
    countryId: selectedCountryId.value
  }
})

const citySelectProps = computed(() => {
  if (!cityField.value) return null
  const field = cityField.value
  const stateName = stateField.value?.name
  return {
    id: field.name,
    label: field.display_name,
    modelValue: props.formData[field.name],
    required: field.is_required,
    help: field.description,
    error: props.hasFieldError(field.name) ? props.getFieldError(field.name)[0] : '',
    stateId: stateName ? props.formData[stateName] : null
  }
})

let prevCountrySnapshot = null
let prevStateSnapshot = null
let prevCitySnapshot = null

watch(countrySelectProps, (val) => {
  if (!val) return
  const changed = JSON.stringify(val) !== JSON.stringify(prevCountrySnapshot)
  if (changed) {
    logger.debug('[LocationFieldGroup] CountrySelect props', val)
    prevCountrySnapshot = { ...val }
  }
}, { deep: true })

watch(stateSelectProps, (val) => {
  if (!val) return
  const changed = JSON.stringify(val) !== JSON.stringify(prevStateSnapshot)
  if (changed) {
    logger.debug('[LocationFieldGroup] StateSelect props', val)
    prevStateSnapshot = { ...val }
  }
}, { deep: true })

watch(citySelectProps, (val) => {
  if (!val) return
  const changed = JSON.stringify(val) !== JSON.stringify(prevCitySnapshot)
  if (changed) {
    logger.debug('[LocationFieldGroup] CitySelect props', val)
    prevCitySnapshot = { ...val }
  }
}, { deep: true })

// Watch explicit stateField presence/changes (including removal)
let prevStateFieldName = null
watch(() => stateField.value ? stateField.value.name : null, (newName, oldName) => {
  if (newName !== oldName) {
    logger.debug('[LocationFieldGroup] stateField change', { from: oldName, to: newName })
  }
  prevStateFieldName = newName
})
</script>
