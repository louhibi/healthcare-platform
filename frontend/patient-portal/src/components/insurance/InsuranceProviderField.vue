<template>
  <div>
    <!-- Insurance Provider Select -->
    <div>
      <label :for="field.name" class="block text-sm font-medium text-gray-700 mb-1">
        {{ field.display_name || field.label }}
        <span v-if="field.is_required" class="text-red-500">*</span>
      </label>
      <select
        :id="field.name"
        :name="field.name"
        v-model="localValue"
        :disabled="field.is_disabled || isLoading || !insuranceTypeSelected"
        :required="field.is_required && insuranceTypeSelected"
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        :class="{
          'bg-gray-100 cursor-not-allowed': field.is_disabled || !insuranceTypeSelected,
          'border-red-300': error,
          'opacity-50': isLoading
        }"
        @change="handleProviderChange"
      >
        <option v-for="provider in insuranceProviders" :key="provider.id" :value="provider.id">{{ provider.name }}</option>
      </select>
  <p v-if="field.description" class="text-xs text-gray-500 mt-1">{{ field.description }}</p>
      <div v-if="isLoading" class="mt-1 text-sm text-blue-600 flex items-center">
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        Loading insurance providers...
      </div>
  <div v-if="!insuranceTypeSelected" class="mt-1 text-sm text-gray-500">
        Please select an insurance type first to see available providers.
      </div>
      <div v-if="error" class="mt-1 text-sm text-red-600">{{ error }}</div>
    </div>
  </div>
</template>
<script setup>
import { locationsApi } from '@/api/locations'
import logger from '@/utils/logger'
import { computed, inject, ref, watch } from 'vue'

const props = defineProps({
  field: { type: Object, required: true },
  modelValue: { type: [String, Number], default: null },
  // Unified single error prop
  error: { type: String, default: '' },
  insuranceTypeId: { type: [String, Number], default: null }
})

logger.debug('[InsuranceProviderField] initial props', {
  fieldName: props.field?.name,
  modelValue: props.modelValue,
  error: props.error,
  insuranceTypeId: props.insuranceTypeId
})

const emit = defineEmits(['update:modelValue'])

// Simple local state (no caching)
const insuranceProviders = ref([])
const loadingProviders = ref(false)

const formData = inject('formData', null)

const localValue = ref(props.modelValue)

const isLoading = computed(() => loadingProviders.value)
const insuranceTypeSelected = computed(() => props.insuranceTypeId && props.insuranceTypeId !== '')

const handleProviderChange = () => {
  logger.debug('[InsuranceProviderField] handleProviderChange', { value: localValue.value })
  emit('update:modelValue', localValue.value)
}


watch(() => props.modelValue, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    logger.debug('[InsuranceProviderField] modelValue prop changed', { oldVal, newVal })
  }
  localValue.value = newVal
}, { immediate: true })

const addAndSortProviders = (providers, typeId) => {
  const hasOther = providers.some(p => p.code === 'other' || p.code.includes('other'))
  if (!hasOther) {
    providers.push({
      id: -1,
      insurance_type_id: typeId,
      code: 'other_provider',
      name: 'Other',
      name_en: 'Other',
      name_fr: 'Autre',
      name_ar: 'أخرى',
      is_default: false,
      sort_order: 999
    })
  }
  return providers.sort((a, b) => {
    if (a.sort_order !== b.sort_order) return a.sort_order - b.sort_order
    return a.name.localeCompare(b.name)
  })
}

const loadInsuranceProvidersForType = async (typeId) => {
  if (!typeId || typeId === -1) {
    insuranceProviders.value = []
    return
  }
  loadingProviders.value = true
  try {
    const providers = await locationsApi.getInsuranceProvidersByType(typeId)
  insuranceProviders.value = addAndSortProviders(Array.isArray(providers) ? providers : [], typeId)
  } catch (e) {
    logger.error('[InsuranceProviderField] Failed to load providers', e)
    insuranceProviders.value = []
  } finally {
    loadingProviders.value = false
  }
}

watch(() => props.insuranceTypeId, async (newTypeId, oldTypeId) => {
  logger.debug('[InsuranceProviderField] insuranceTypeId watcher triggered', { newTypeId, oldTypeId })
  localValue.value = null
  emit('update:modelValue', null)
  if (newTypeId) {
    logger.debug('[InsuranceProviderField] Loading providers for type', { typeId: newTypeId })
    await loadInsuranceProvidersForType(newTypeId)
    logger.debug('[InsuranceProviderField] Providers loaded', { count: insuranceProviders.value.length })
    const realProviders = insuranceProviders.value.filter(p => p.id !== -1)
    if (realProviders.length === 1) {
      localValue.value = realProviders[0].id
      logger.debug('[InsuranceProviderField] Auto-selected sole provider after type change', { providerId: localValue.value })
      emit('update:modelValue', localValue.value)
    }
  } else {
    insuranceProviders.value = []
  }
}, { immediate: true })

watch(insuranceProviders, (list, old) => {
  logger.debug('[InsuranceProviderField] insuranceProviders changed', {
    previousCount: old ? old.length : 0,
    newCount: list.length
  })
  const realProviders = list.filter(p => p.id !== -1)
  if (!localValue.value && realProviders.length === 1) {
    localValue.value = realProviders[0].id
    logger.debug('[InsuranceProviderField] Auto-selected sole provider on provider list change', { providerId: localValue.value })
    emit('update:modelValue', localValue.value)
  }
})

// removed onMounted (not needed beyond watcher logic)
</script>
