<template>
  <div>
    <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Contact Information</h4>
    <div class="space-y-3">
      <div>
        <label class="text-sm font-medium text-gray-600">Email</label>
        <p class="text-gray-900">{{ patient.email }}</p>
      </div>
      <div>
        <label class="text-sm font-medium text-gray-600">Phone</label>
        <p class="text-gray-900">{{ patient.phone }}</p>
      </div>
      <div v-if="patient.address">
        <label class="text-sm font-medium text-gray-600">Address</label>
        <p class="text-gray-900">{{ patient.address }}</p>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div v-if="patient.city">
          <label class="text-sm font-medium text-gray-600">City</label>
          <p class="text-gray-900">{{ patient.city }}</p>
        </div>
        <div v-if="patient.state">
          <label class="text-sm font-medium text-gray-600">State/Province</label>
          <p class="text-gray-900">{{ patient.state }}</p>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div v-if="patient.postal_code">
          <label class="text-sm font-medium text-gray-600">Postal Code</label>
          <p class="text-gray-900">{{ patient.postal_code }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Country</label>
          <p class="text-gray-900">{{ getCountryName(patient.country) }}</p>
        </div>
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
      const required = ['email', 'phone', 'country']
      return required.every(field => field in value)
    }
  },
  getCountryName: {
    type: Function,
    required: true
  }
})
</script>