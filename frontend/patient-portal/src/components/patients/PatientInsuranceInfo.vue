<template>
  <div>
    <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Insurance & Emergency Contact</h4>
    <div class="space-y-4">
      <!-- Insurance Information -->
      <div>
        <h5 class="text-md font-medium text-gray-800 mb-2">Insurance</h5>
        <div class="space-y-2">
          <div v-if="patient.insurance_provider">
            <label class="text-sm font-medium text-gray-600">Insurance Provider</label>
            <p class="text-gray-900">{{ patient.insurance_provider }}</p>
          </div>
          <div v-if="patient.insurance_type">
            <label class="text-sm font-medium text-gray-600">Insurance Type</label>
            <p class="text-gray-900">{{ patient.insurance_type }}</p>
          </div>
          <div v-if="patient.policy_number">
            <label class="text-sm font-medium text-gray-600">Policy Number</label>
            <p class="text-gray-900">{{ patient.policy_number }}</p>
          </div>
          <div v-if="!patient.insurance_provider && !patient.insurance_type">
            <p class="text-gray-500">No insurance information on file</p>
          </div>
        </div>
      </div>

      <!-- Emergency Contact -->
      <div v-if="hasEmergencyContact">
        <h5 class="text-md font-medium text-gray-800 mb-2">Emergency Contact</h5>
        <div class="space-y-2">
          <div v-if="patient.emergency_contact.name">
            <label class="text-sm font-medium text-gray-600">Name</label>
            <p class="text-gray-900">{{ patient.emergency_contact.name }}</p>
          </div>
          <div v-if="patient.emergency_contact.phone">
            <label class="text-sm font-medium text-gray-600">Phone</label>
            <p class="text-gray-900">{{ patient.emergency_contact.phone }}</p>
          </div>
          <div v-if="patient.emergency_contact.relationship">
            <label class="text-sm font-medium text-gray-600">Relationship</label>
            <p class="text-gray-900">{{ formatGender(patient.emergency_contact.relationship) }}</p>
          </div>
        </div>
      </div>
      <div v-else>
        <h5 class="text-md font-medium text-gray-800 mb-2">Emergency Contact</h5>
        <p class="text-gray-500">No emergency contact information on file</p>
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
  },
  formatGender: {
    type: Function,
    required: true
  }
})

// Computed
const hasEmergencyContact = computed(() => {
  const contact = props.patient.emergency_contact
  return !!(contact && (contact.name || contact.phone))
})
</script>