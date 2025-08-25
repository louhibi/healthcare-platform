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
      :disabled="disabled || loading"
      v-model.number="internalValue"
      @change="emitChange"
    >
      <option :value="''" :disabled="true">
        {{ loading ? t('loading.countries', 'Loading countries…') : placeholder }}
      </option>
      <option v-for="c in sortedCountries" :key="c.id" :value="c.id">
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
import { computed, onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: [Number, String, null], default: null },
  countryId: { type: [Number, String, null], default: null },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Select Country' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  help: { type: String, default: '' },
  error: { type: String, default: '' },
  id: { type: String, default: null },
  refetchOnLocaleChange: { type: Boolean, default: true },
  nameResolver: { type: Function, default: null }
})

const emit = defineEmits(['update:modelValue', 'country-change', 'loaded'])
const { t, locale } = useI18n()
const countries = ref([])
const loading = ref(false)
const internalValue = ref(null)
const idComputed = computed(() => props.id || 'country-select')

const loadCountries = async () => {
  loading.value = true
  try {
    const data = await locationsApi.getCountries()
    countries.value = Array.isArray(data) ? data : []
    emit('loaded', countries.value)
    selectInitial()
  } catch (e) {
    console.error('Failed to load countries', e)
    countries.value = []
  } finally {
    loading.value = false
  }
}

const selectInitial = () => {
  const initial = props.modelValue ?? props.countryId
  if (initial && countries.value.some(c => c.id === Number(initial))) {
    internalValue.value = Number(initial)
  } else {
    internalValue.value = null
  }
}

const localizedName = (country) => {
  if (!country) return ''
  if (props.nameResolver) {
    try { return props.nameResolver(country, locale.value) } catch { /* noop */ }
  }
  const keyCandidates = []
  if (country.code) keyCandidates.push(`countries.${country.code}`)
  if (country.id) keyCandidates.push(`countries.${country.id}`)
  for (const key of keyCandidates) {
    const translated = t(key, country.name || key)
    if (translated !== key) return translated
  }
  return country.name
}

const sortedCountries = computed(() => [...countries.value].sort((a, b) => localizedName(a).localeCompare(localizedName(b))) )
const selectBorderClass = computed(() => props.error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300')
const emitChange = () => { emit('update:modelValue', internalValue.value); emit('country-change', internalValue.value) }

watch(() => props.modelValue, (val) => { if (val !== internalValue.value) internalValue.value = val ? Number(val) : null })
watch(locale, async () => { if (props.refetchOnLocaleChange) await loadCountries() })
onMounted(loadCountries)
</script>

<style scoped>
</style>
