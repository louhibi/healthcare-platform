<template>
  <div class="w-full">
    <label v-if="label" :for="idComputed" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <select
      :id="idComputed"
      class="mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
      :class="selectBorderClass"
      :disabled="disabled || loading || !stateId"
      v-model.number="internalValue"
      @change="emitChange"
    >
      <option :value="''" :disabled="true">
        {{ !stateId ? t('select.stateFirst', 'Select state first') : (loading ? t('loading.cities', 'Loading cities…') : placeholder) }}
      </option>
      <option v-for="c in sortedCities" :key="c.id" :value="c.id">
        {{ localizedName(c) }}
      </option>
    </select>
    <p v-if="help" class="text-xs text-gray-500 mt-1">{{ help }}</p>
    <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
  </div>
</template>

<script setup>
import { locationsApi } from '@/api/locations'
import { useI18n } from '@/composables/useI18n'
import logger from '@/utils/logger'
import { computed, onMounted, ref, watch } from 'vue'

// Props mirror StateSelect with city-specific naming
const props = defineProps({
  modelValue: { type: [Number, String, null], default: null },
  cityId: { type: [Number, String, null], default: null }, // alias for initial selection
  stateId: { type: [Number, String, null], default: null },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Select City' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  help: { type: String, default: '' },
  error: { type: String, default: '' },
  id: { type: String, default: null },
  refetchOnLocaleChange: { type: Boolean, default: true },
  nameResolver: { type: Function, default: null },
  clearOnStateChange: { type: Boolean, default: true }
})

// Initial props snapshot
logger.debug('[CitySelect] initial props', {
  modelValue: props.modelValue,
  cityId: props.cityId,
  stateId: props.stateId,
  label: props.label,
  placeholder: props.placeholder,
  required: props.required,
  disabled: props.disabled,
  help: props.help,
  error: props.error,
  id: props.id,
  refetchOnLocaleChange: props.refetchOnLocaleChange,
  hasNameResolver: !!props.nameResolver,
  clearOnStateChange: props.clearOnStateChange
})

const emit = defineEmits(['update:modelValue', 'city-change', 'loaded'])
const { t, locale } = useI18n()
const cities = ref([])
const loading = ref(false)
const internalValue = ref(null)
const lastLoadedState = ref(null)
const idComputed = computed(() => props.id || 'city-select')

const numeric = (val) => {
  if (val === null || val === undefined || val === '') return null
  const n = Number(val)
  return isNaN(n) ? null : n
}

const loadCities = async () => {
  logger.debug('[CitySelect] loadCities called', { stateId: props.stateId })
  if (!props.stateId) {
    cities.value = []
    internalValue.value = null
    logger.debug('[CitySelect] No stateId provided, cleared cities')
    return
  }
  // Avoid refetch if same state and we already have data unless locale changed
  if (lastLoadedState.value === props.stateId && cities.value.length > 0) {
    logger.debug('[CitySelect] Skipping refetch, cities already loaded for state', props.stateId)
    return
  }
  loading.value = true
  try {
    logger.debug('[CitySelect] Fetching cities for state', props.stateId)
    const data = await locationsApi.getCitiesByStateId(props.stateId)
    cities.value = Array.isArray(data) ? data : []
    lastLoadedState.value = props.stateId
    emit('loaded', cities.value)
    logger.debug('[CitySelect] Cities loaded', { count: cities.value.length })
    selectInitial()
  } catch (e) {
    logger.error('[CitySelect] Failed to load cities', e)
    cities.value = []
  } finally {
    logger.debug('[CitySelect] Finished loadCities', { loading: false })
    loading.value = false
  }
}

const selectInitial = () => {
  const initial = props.modelValue ?? props.cityId
  logger.debug('[CitySelect] selectInitial invoked', { initial, citiesCount: cities.value.length })
  if (initial && cities.value.some(c => c.id === Number(initial))) {
    internalValue.value = Number(initial)
    logger.debug('[CitySelect] Initial city selected', internalValue.value)
  } else {
    internalValue.value = null
    logger.debug('[CitySelect] Initial city not found, cleared selection')
  }
}

const localizedName = (city) => {
  if (!city) return ''
  if (props.nameResolver) {
    try { return props.nameResolver(city, locale.value) } catch { /* noop */ }
  }
  const keyCandidates = []
  if (city.code) keyCandidates.push(`cities.${city.code}`)
  if (city.id) keyCandidates.push(`cities.${city.id}`)
  for (const key of keyCandidates) {
    const translated = t(key, city.name || key)
    if (translated !== key) return translated
  }
  return city.name
}

const sortedCities = computed(() => [...cities.value].sort((a, b) => localizedName(a).localeCompare(localizedName(b))) )
const selectBorderClass = computed(() => props.error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300')
const emitChange = () => {
  logger.debug('[CitySelect] emitChange', { value: internalValue.value })
  emit('update:modelValue', internalValue.value)
  emit('city-change', internalValue.value)
}

// Watchers
watch(() => props.modelValue, (val) => {
  const num = numeric(val)
  if (num !== internalValue.value) {
  logger.debug('[CitySelect] modelValue watcher updating internalValue', { from: internalValue.value, to: num })
    internalValue.value = num
  }
})

watch(() => props.stateId, async (newState, oldState) => {
  logger.debug('[CitySelect] stateId changed', { oldState, newState })
  if (newState !== oldState) {
    if (props.clearOnStateChange) {
      logger.debug('[CitySelect] Clearing selection due to state change')
      internalValue.value = null
    }
    await loadCities()
  } else if (newState && cities.value.length === 0) {
    logger.debug('[CitySelect] State unchanged but cities empty, loading')
    await loadCities()
  }
})

watch(locale, async (newLoc, oldLoc) => {
  logger.debug('[CitySelect] locale changed', { oldLoc, newLoc })
  if (props.refetchOnLocaleChange && props.stateId) {
    lastLoadedState.value = null
    await loadCities()
  }
})

onMounted(() => {
  logger.debug('[CitySelect] mounted', { stateId: props.stateId, modelValue: props.modelValue, cityId: props.cityId })
  if (props.stateId) loadCities()
})
</script>

<style scoped>
</style>