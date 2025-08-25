<template>
  <tr class="hover:bg-gray-50">
    <!-- Patient Info -->
    <td class="px-6 py-4 whitespace-nowrap">
      <div class="flex items-center">
        <div class="flex-shrink-0 h-10 w-10">
          <div class="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
            <span class="text-sm font-medium text-gray-700">
              {{ patient.first_name.charAt(0) }}{{ patient.last_name.charAt(0) }}
            </span>
          </div>
        </div>
        <div class="ml-4">
          <div class="text-sm font-medium text-gray-900">
            {{ patient.first_name }} {{ patient.last_name }}
          </div>
          <div class="text-sm text-gray-500">
            {{ $t('patients.patientId') }}: {{ patient.patient_id || patient.id }}
          </div>
        </div>
      </div>
    </td>

    <!-- Contact Info -->
    <td class="px-6 py-4 whitespace-nowrap">
      <div class="text-sm text-gray-900">{{ patient.email }}</div>
      <div class="text-sm text-gray-500">{{ patient.phone }}</div>
    </td>

    <!-- Demographics -->
    <td class="px-6 py-4 whitespace-nowrap">
      <div class="text-sm text-gray-900">
        {{ formatGender(patient.gender) }}, {{ patient.age }} {{ $t('common.yearsOld') }}
      </div>
      <div class="text-sm text-gray-500">
        {{ patient.city }}, {{ getCountryName(patient.country) }}
      </div>
      <div v-if="patient.blood_type" class="text-sm text-gray-500">
        {{ $t('forms.patient.fields.bloodType.label') }}: {{ patient.blood_type }}
      </div>
    </td>

    <!-- Insurance -->
    <td class="px-6 py-4 whitespace-nowrap">
      <!-- Show provider first, then type, then policy number -->
      <div v-if="patient.insurance_provider || patient.insurance_type" class="space-y-1">
        <div v-if="patient.insurance_provider" class="text-sm text-gray-900">
          {{ patient.insurance_provider }}
        </div>
        <div v-else-if="patient.insurance_type" class="text-sm text-gray-900">
          {{ patient.insurance_type }}
        </div>
        <div v-if="patient.insurance_type && patient.insurance_provider" class="text-xs text-gray-500">
          {{ patient.insurance_type }}
        </div>
        <div v-if="patient.policy_number" class="text-xs text-gray-500">
          {{ patient.policy_number }}
        </div>
      </div>
      <div v-else class="text-sm text-gray-500">
        {{ $t('patients.noInsurance') }}
      </div>
    </td>

    <!-- Actions -->
    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
      <PatientActions
        :patient="patient"
        @view="$emit('view', patient)"
        @edit="$emit('edit', patient)"
        @delete="$emit('delete', patient)"
      />
    </td>
  </tr>
</template>

<script setup>
import PatientActions from './PatientActions.vue'

// Props with validation
const props = defineProps({
  patient: {
    type: Object,
    required: true,
    validator: (value) => {
      const required = ['id', 'first_name', 'last_name', 'email', 'phone', 'gender', 'age', 'country']
      return required.every(field => field in value)
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
  'view': (patient) => patient && typeof patient === 'object',
  'edit': (patient) => patient && typeof patient === 'object',
  'delete': (patient) => patient && typeof patient === 'object'
})
</script>