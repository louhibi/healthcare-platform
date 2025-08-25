<template>
  <div 
    v-if="show && patient" 
    class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" 
    @click="$emit('close')"
  >
    <div 
      class="relative top-10 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white" 
      @click.stop
    >
      <div class="mt-3">
        <!-- Modal Header -->
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold text-gray-900">Patient Details</h3>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Patient Details Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Personal Information -->
          <PatientPersonalInfo 
            :patient="patient" 
            :format-gender="formatGender" 
          />

          <!-- Contact Information -->
          <PatientContactInfo 
            :patient="patient" 
            :get-country-name="getCountryName" 
          />

          <!-- Medical Information -->
          <PatientMedicalInfo :patient="patient" />

          <!-- Insurance & Emergency Contact -->
          <PatientInsuranceInfo 
            :patient="patient" 
            :format-gender="formatGender" 
          />
        </div>

        <!-- Action Buttons -->
        <div class="mt-8 pt-4 border-t border-gray-200 flex justify-end space-x-3">
          <button
            @click="$emit('edit', patient)"
            class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
          >
            Edit Patient
          </button>
          <button
            @click="$emit('close')"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import PatientPersonalInfo from './PatientPersonalInfo.vue'
import PatientContactInfo from './PatientContactInfo.vue'
import PatientMedicalInfo from './PatientMedicalInfo.vue'
import PatientInsuranceInfo from './PatientInsuranceInfo.vue'

// Props with validation
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  patient: {
    type: Object,
    default: null,
    validator: (value) => {
      if (!value) return true
      return typeof value === 'object' && 'id' in value
    }
  },
  formatGender: {
    type: Function,
    required: true
  },
  getCountryName: {
    type: Function,
    required: true
  }
})

// Emits with validation
defineEmits({
  'close': () => true,
  'edit': (patient) => patient && typeof patient === 'object'
})
</script>