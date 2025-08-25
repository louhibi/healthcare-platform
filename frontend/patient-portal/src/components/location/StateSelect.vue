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
      :disabled="disabled || loading || !countryId"
      v-model.number="internalValue"
      @change="emitChange"
    >
      <option :value="''" :disabled="true">
        {{ !countryId ? t('select.countryFirst', 'Select country first') : (loading ? t('loading.states', 'Loading states…') : placeholder) }}
      </option>
      <option v-for="s in sortedStates" :key="s.id" :value="s.id">
        {{ localizedName(s) }}
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

const props = defineProps({
  modelValue: { type: [Number, String, null], default: null },
  stateId: { type: [Number, String, null], default: null }, // alias for initial selection
  countryId: { type: [Number, String, null], default: null },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Select State/Province' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  help: { type: String, default: '' },
  error: { type: String, default: '' },
  id: { type: String, default: null },
  refetchOnLocaleChange: { type: Boolean, default: true },
  nameResolver: { type: Function, default: null },
  clearOnCountryChange: { type: Boolean, default: true }
})

// Log initial props snapshot
logger.debug('[StateSelect] initial props', {
  modelValue: props.modelValue,
  stateId: props.stateId,
  countryId: props.countryId,
  label: props.label,
  placeholder: props.placeholder,
  required: props.required,
  disabled: props.disabled,
  help: props.help,
  error: props.error,
  id: props.id,
  refetchOnLocaleChange: props.refetchOnLocaleChange,
  hasNameResolver: !!props.nameResolver,
  clearOnCountryChange: props.clearOnCountryChange
})

const emit = defineEmits(['update:modelValue', 'state-change', 'loaded'])
const { t, locale } = useI18n()
const states = ref([])
const loading = ref(false)
const internalValue = ref(null)
const lastLoadedCountry = ref(null)
const idComputed = computed(() => props.id || 'state-select')

const numeric = (val) => {
  if (val === null || val === undefined || val === '') return null
  const n = Number(val)
  return isNaN(n) ? null : n
}

const loadStates = async () => {
  logger.debug('[StateSelect] loadStates called', { countryId: props.countryId })
  if (!props.countryId) {
    states.value = []
    internalValue.value = null
  logger.debug('[StateSelect] No countryId provided, cleared states')
    return
  }
  // Avoid refetch if same country and we already have data unless locale changed (handled separately)
  if (lastLoadedCountry.value === props.countryId && states.value.length > 0) {
  logger.debug('[StateSelect] Skipping refetch, states already loaded for country', props.countryId)
    return
  }
  loading.value = true
  try {
  logger.debug('[StateSelect] Fetching states for country', props.countryId)
    const data = await locationsApi.getStatesByCountryId(props.countryId)
    states.value = Array.isArray(data) ? data : []
    lastLoadedCountry.value = props.countryId
    emit('loaded', states.value)
  logger.debug('[StateSelect] States loaded', { count: states.value.length })
    selectInitial()
  } catch (e) {
  logger.error('[StateSelect] Failed to load states', e)
    states.value = []
  } finally {
  logger.debug('[StateSelect] Finished loadStates', { loading: false })
    loading.value = false
  }
}

const selectInitial = () => {
  const initial = props.modelValue ?? props.stateId
  logger.debug('[StateSelect] selectInitial invoked', { initial, statesCount: states.value.length })
  if (initial && states.value.some(s => s.id === Number(initial))) {
    internalValue.value = Number(initial)
  logger.debug('[StateSelect] Initial state selected', internalValue.value)
  } else {
    internalValue.value = null
  logger.debug('[StateSelect] Initial state not found, cleared selection')
  }
}

const localizedName = (state) => {
  if (!state) return ''
  if (props.nameResolver) {
    try { return props.nameResolver(state, locale.value) } catch { /* noop */ }
  }
  const keyCandidates = []
  if (state.code) keyCandidates.push(`states.${state.code}`)
  if (state.id) keyCandidates.push(`states.${state.id}`)
  for (const key of keyCandidates) {
    const translated = t(key, state.name || key)
    if (translated !== key) return translated
  }
  return state.name
}

const sortedStates = computed(() => [...states.value].sort((a, b) => localizedName(a).localeCompare(localizedName(b))) )
const selectBorderClass = computed(() => props.error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300')
const emitChange = () => {
  logger.debug('[StateSelect] emitChange', { value: internalValue.value })
  emit('update:modelValue', internalValue.value)
  emit('state-change', internalValue.value)
}

// Watchers
watch(() => props.modelValue, (val) => {
  const num = numeric(val)
  if (num !== internalValue.value) {
  logger.debug('[StateSelect] modelValue watcher updating internalValue', { from: internalValue.value, to: num })
    internalValue.value = num
  }
})
watch(() => props.countryId, async (newCountry, oldCountry) => {
  logger.debug('[StateSelect] countryId changed', { oldCountry, newCountry })
  if (newCountry !== oldCountry) {
    if (props.clearOnCountryChange) {
  logger.debug('[StateSelect] Clearing selection due to country change')
      internalValue.value = null
    }
    await loadStates()
  } else if (newCountry && states.value.length === 0) {
  logger.debug('[StateSelect] Country unchanged but states empty, loading')
    await loadStates()
  }
})
watch(locale, async (newLoc, oldLoc) => {
  logger.debug('[StateSelect] locale changed', { oldLoc, newLoc })
  if (props.refetchOnLocaleChange && props.countryId) {
    lastLoadedCountry.value = null
    await loadStates()
  }
})

onMounted(() => {
  logger.debug('[StateSelect] mounted', { countryId: props.countryId, modelValue: props.modelValue, stateId: props.stateId })
  if (props.countryId) loadStates()
})
</script>

<style scoped>
</style>
