<template>
  <div>
    <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Medical Information</h4>
    <div class="space-y-3">
      <div v-if="patient.blood_type">
        <label class="text-sm font-medium text-gray-600">Blood Type</label>
        <p class="text-gray-900">{{ patient.blood_type }}</p>
      </div>
      <div v-if="patient.allergies">
        <label class="text-sm font-medium text-gray-600">Allergies</label>
        <p class="text-gray-900">{{ patient.allergies }}</p>
      </div>
      <div v-if="patient.medications">
        <label class="text-sm font-medium text-gray-600">Current Medications</label>
        <p class="text-gray-900">{{ patient.medications }}</p>
      </div>
      <div v-if="patient.medical_history">
        <label class="text-sm font-medium text-gray-600">Medical History</label>
        <p class="text-gray-900 whitespace-pre-wrap">{{ patient.medical_history }}</p>
      </div>
      <div v-if="!hasAnyMedicalInfo" class="text-gray-500">
        No medical information on file
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props with validation
const props = defineProps({
  patient: {
    type: Object,
    required: true
  }
})

// Computed
const hasAnyMedicalInfo = computed(() => {
  return !!(props.patient.blood_type || 
           props.patient.allergies || 
           props.patient.medications || 
           props.patient.medical_history)
})
</script>