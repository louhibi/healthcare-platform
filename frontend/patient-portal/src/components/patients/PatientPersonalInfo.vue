<template>
  <div>
    <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Personal Information</h4>
    <div class="space-y-3">
      <div>
        <label class="text-sm font-medium text-gray-600">Full Name</label>
        <p class="text-gray-900">{{ patient.first_name }} {{ patient.last_name }}</p>
      </div>
      <div v-if="patient.patient_id">
        <label class="text-sm font-medium text-gray-600">Patient ID</label>
        <p class="text-gray-900">{{ patient.patient_id }}</p>
      </div>
      <div>
        <label class="text-sm font-medium text-gray-600">Date of Birth</label>
        <p class="text-gray-900">{{ formatDate(patient.date_of_birth) }}</p>
      </div>
      <div>
        <label class="text-sm font-medium text-gray-600">Gender</label>
        <p class="text-gray-900">{{ formatGender(patient.gender) }}</p>
      </div>
      <div v-if="patient.nationality">
        <label class="text-sm font-medium text-gray-600">Nationality</label>
        <p class="text-gray-900">{{ patient.nationality }}</p>
      </div>
      <div v-if="patient.preferred_language">
        <label class="text-sm font-medium text-gray-600">Preferred Language</label>
        <p class="text-gray-900">{{ patient.preferred_language.toUpperCase() }}</p>
      </div>
      <div v-if="patient.marital_status">
        <label class="text-sm font-medium text-gray-600">Marital Status</label>
        <p class="text-gray-900">{{ formatGender(patient.marital_status) }}</p>
      </div>
      <div v-if="patient.occupation">
        <label class="text-sm font-medium text-gray-600">Occupation</label>
        <p class="text-gray-900">{{ patient.occupation }}</p>
      </div>
      <div v-if="patient.national_id">
        <label class="text-sm font-medium text-gray-600">National ID</label>
        <p class="text-gray-900">{{ patient.national_id }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
// Props with validation
const props = defineProps({
  patient: {
    type: Object,
    required: true,
    validator: (value) => {
      const required = ['first_name', 'last_name', 'date_of_birth', 'gender']
      return required.every(field => field in value)
    }
  },
  formatGender: {
    type: Function,
    required: true
  }
})

// Methods
const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}
</script>