<template>
  <div>
    <!-- Insurance Type Select -->
    <div>
      <label :for="field.name" class="block text-sm font-medium text-gray-700 mb-1">
        {{ field.display_name || field.label }}
        <span v-if="field.is_required" class="text-red-500">*</span>
      </label>
      <select
        :id="field.name"
        :name="field.name"
        :label="field.display_name"
        v-model="localValue"
        :disabled="field.is_disabled || isLoading"
        :required="field.is_required"
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        :class="{
          'bg-gray-100 cursor-not-allowed': field.is_disabled,
          'border-red-300': error,
          'opacity-50': isLoading
        }"
        @change="handleInsuranceTypeChange"
      >
        <option v-for="type in insuranceTypes" :key="type.id" :value="type.id">{{ type.name }}</option>
      </select>
  <p v-if="field.description" class="text-xs text-gray-500 mt-1">{{ field.description }}</p>
      <div v-if="isLoading" class="mt-1 text-sm text-blue-600 flex items-center">
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        Loading insurance types...
      </div>
      <div v-if="error" class="mt-1 text-sm text-red-600">{{ error }}</div>
    </div>
  </div>
</template>
<script setup>
import { locationsApi } from '@/api/locations'
import logger from '@/utils/logger'
import { computed, inject, onMounted, ref, watch } from 'vue'

const props = defineProps({
  field: { type: Object, required: true },
  modelValue: { type: [String, Number], default: null },
  // Unified single error prop (string) aligning with Location components pattern
  error: { type: String, default: '' },
  countryId: { type: [String, Number], default: null }
})

logger.debug('[InsuranceTypeField] initial props', {
  fieldName: props.field?.name,
  modelValue: props.modelValue,
  error: props.error,
  countryId: props.countryId
})

const emit = defineEmits(['update:modelValue', 'insuranceTypeChange'])

// Local simple state (no caching)
const insuranceTypes = ref([])
const loadingTypes = ref(false)

const formData = inject('formData', null)

const localValue = ref(props.modelValue)

const isLoading = computed(() => loadingTypes.value)

const handleInsuranceTypeChange = async () => {
  logger.debug('[InsuranceTypeField] handleInsuranceTypeChange invoked', { newValue: localValue.value })
  emit('update:modelValue', localValue.value)
  // Cascade clearing (mimic removed composable logic)
  if (formData) {
    // Clear provider related fields when type changes
    formData.insurance_provider_id = null
    formData.insurance_provider_custom = null
    // If switching away from "Other" clear custom insurance type text
    if (!isOtherType(localValue.value)) {
      formData.insurance_custom = null
    }
  }
  emit('insuranceTypeChange', localValue.value)
}

const isOtherType = (typeId) => typeId === -1 || typeId === '-1'

const addAndSortTypes = (types) => {
  const hasOther = types.some(t => t.code === 'other')
  if (!hasOther) {
    types.push({
      id: -1,
      country_id: null,
      code: 'other',
      name: 'Other',
      name_en: 'Other',
      name_fr: 'Autre',
      name_ar: 'أخرى',
      is_default: false,
      sort_order: 999
    })
  }
  return types.sort((a, b) => {
    if (a.sort_order !== b.sort_order) return a.sort_order - b.sort_order
    return a.name.localeCompare(b.name)
  })
}

const loadInsuranceTypesForCountry = async (countryId) => {
  if (countryId == null || countryId === '') {
    insuranceTypes.value = []
    return
  }
  loadingTypes.value = true
  try {
    const types = await locationsApi.getInsuranceTypesByCountry(countryId)
    insuranceTypes.value = addAndSortTypes(Array.isArray(types) ? types : [])
  } catch (e) {
    logger.error('[InsuranceTypeField] Failed to load insurance types', e)
    insuranceTypes.value = []
  } finally {
    loadingTypes.value = false
  }
}


watch(() => props.modelValue, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    logger.debug('[InsuranceTypeField] modelValue prop changed', { oldVal, newVal })
  }
  localValue.value = newVal
}, { immediate: true })

onMounted(async () => {
  logger.debug('[InsuranceTypeField] mounted', {
    fieldName: props.field?.name,
    modelValue: props.modelValue,
    countryId: props.countryId,
    initialTypes: insuranceTypes.value.length
  })
  if (props.countryId) {
    await loadInsuranceTypesForCountry(props.countryId)
  } else {
    logger.debug('[InsuranceTypeField] mount skipped loading: no countryId')
  }
})

watch(() => props.countryId, async (newCountryId, oldCountryId) => {
  if (newCountryId !== oldCountryId) {
    logger.debug('[InsuranceTypeField] countryId watcher triggered', { oldCountryId, newCountryId })
  }
  if (newCountryId && newCountryId !== oldCountryId) {
    await loadInsuranceTypesForCountry(newCountryId)
  } else if (!newCountryId) {
    insuranceTypes.value = []
    localValue.value = null
  }
})
// Debug watcher (optional)
watch(insuranceTypes, (newList, oldList) => {
  logger.debug('[InsuranceTypeField] insuranceTypes changed', { previousCount: oldList ? oldList.length : 0, newCount: newList.length })
})
</script>
