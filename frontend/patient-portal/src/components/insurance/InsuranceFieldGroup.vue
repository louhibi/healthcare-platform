<template>
  <div class="space-y-4">
    <!-- Insurance Type Field -->
    <div v-if="insuranceTypeField">
      <InsuranceTypeField
        :field="insuranceTypeField"
        v-model="formData[insuranceTypeField.name]"
        :error="hasFieldError && getFieldError && hasFieldError(insuranceTypeField.name) ? (getFieldError(insuranceTypeField.name)[0] || '') : ''"
        :country-id="countryId"
      />
    </div>

    <!-- Insurance Provider Field -->
    <div v-if="insuranceProviderField">
      <InsuranceProviderField
        :field="insuranceProviderField"
        v-model="formData[insuranceProviderField.name]"
        :error="hasFieldError && getFieldError && hasFieldError(insuranceProviderField.name) ? (getFieldError(insuranceProviderField.name)[0] || '') : ''"
        :insurance-type-id="formData[insuranceTypeField?.name]"
      />
    </div>

    <!-- Other Insurance Fields (like policy number) -->
    <div v-for="field in otherInsuranceFields" :key="field.name" class="space-y-2">
      <FormField
        :field="field"
        :model-value="formData[field.name]"
        :has-error="hasFieldError ? hasFieldError(field.name) : false"
        :error-message="getFieldError ? getFieldError(field.name) : null"
        :disabled="loading || field.is_disabled"
        :loading="loading"
        @field-update="val => updateValue(field.name, val)"
        @field-input="val => updateValue(field.name, val)"
        @field-change="val => updateValue(field.name, val)"
      />
    </div>
  </div>
</template>

<script setup>
import logger from '@/utils/logger'
import { computed, provide, watch } from 'vue'
import FormField from '../FormField.vue'
import InsuranceProviderField from './InsuranceProviderField.vue'
import InsuranceTypeField from './InsuranceTypeField.vue'

const props = defineProps({
  fields: { type: Array, required: true },
  formData: { type: Object, required: true },
  hasFieldError: { type: Function, default: null },
  getFieldError: { type: Function, default: null },
  loading: { type: Boolean, default: false },
  countryId: { type: [Number, String], default: null }
})

const emit = defineEmits(['insurance-change'])

logger.debug('[InsuranceFieldGroup] initial props', {
  fieldNames: props.fields.map(f => f.name),
  hasType: props.fields.some(f => ['insurance','insurance_type','insurance_type_id'].includes(f.name)),
  hasProvider: props.fields.some(f => ['insurance_provider','insurance_provider_id'].includes(f.name)),
  countryId: props.countryId
})

provide('formData', props.formData)

const insuranceTypeField = computed(() => {
  return props.fields.find(field =>
    field.name === 'insurance' ||
    field.name === 'insurance_type' ||
    field.name === 'insurance_type_id'
  )
})

const insuranceProviderField = computed(() => {
  return props.fields.find(field =>
    field.name === 'insurance_provider' ||
    field.name === 'insurance_provider_id'
  )
})

const otherInsuranceFields = computed(() => {
  return props.fields.filter(field =>
    field.name !== 'insurance' &&
    field.name !== 'insurance_type' &&
    field.name !== 'insurance_type_id' &&
    field.name !== 'insurance_provider' &&
    field.name !== 'insurance_provider_id'
  )
})

const updateValue = (name, val) => {
  logger.debug('[InsuranceFieldGroup] updateValue', { name, val })
  if (props.formData) {
    props.formData[name] = val
  }
  emitAggregatedInsurance()
}

const emitAggregatedInsurance = () => {
  const insuranceTypeId = insuranceTypeField.value ? (props.formData[insuranceTypeField.value.name] ?? null) : null
  const insuranceProviderId = insuranceProviderField.value ? (props.formData[insuranceProviderField.value.name] ?? null) : null
  const policyNumber = props.formData['policy_number'] ?? null
  const payload = { insuranceTypeId, insuranceProviderId, policyNumber }
  logger.debug('[InsuranceFieldGroup] insurance-change emit', payload)
  emit('insurance-change', payload)
}

watch(() => [
  insuranceTypeField.value ? props.formData[insuranceTypeField.value.name] : null,
  insuranceProviderField.value ? props.formData[insuranceProviderField.value.name] : null,
  props.formData['policy_number']
], () => {
  emitAggregatedInsurance()
})
</script>
