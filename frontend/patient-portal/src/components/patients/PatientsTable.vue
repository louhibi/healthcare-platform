<template>
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <!-- Table Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <h3 class="text-lg font-medium text-gray-900">
        {{ $t('patients.patientList') }}
        <span class="text-sm font-normal text-gray-500">
          ({{ totalCount }} {{ $t('common.total') }})
        </span>
      </h3>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-6 text-center text-red-600">
      <p>{{ error }}</p>
      <button @click="$emit('retry')" class="mt-2 text-blue-600 hover:text-blue-800">
        {{ $t('common.retry') }}
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="patients.length === 0" class="p-6 text-center text-gray-500">
      <p>{{ $t('patients.noPatients') }}</p>
    </div>

    <!-- Patients Table -->
    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('patients.patient') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('patients.contactInfo') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.demographics') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('patients.insurance') }}
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              {{ $t('common.actions') }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <PatientTableRow
            v-for="patient in patients"
            :key="patient.id"
            :patient="patient"
            :format-gender="formatGender"
            :get-country-name="getCountryName"
            @view="$emit('view-patient', patient)"
            @edit="$emit('edit-patient', patient)"
            @delete="$emit('delete-patient', patient)"
          />
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div 
      v-if="totalCount > limit" 
      class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6"
    >
      <!-- Mobile Pagination -->
      <div class="flex-1 flex justify-between sm:hidden">
        <button
          @click="$emit('prev-page')"
          :disabled="offset === 0"
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          {{ $t('common.previous') }}
        </button>
        <button
          @click="$emit('next-page')"
          :disabled="!hasMore"
          class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          {{ $t('common.next') }}
        </button>
      </div>

      <!-- Desktop Pagination -->
      <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            {{ $t('tables.showing') }}
            <span class="font-medium">{{ offset + 1 }}</span>
            {{ $t('tables.to') }}
            <span class="font-medium">{{ Math.min(offset + limit, totalCount) }}</span>
            {{ $t('tables.of') }}
            <span class="font-medium">{{ totalCount }}</span>
            {{ $t('tables.results') }}
          </p>
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
            <button
              @click="$emit('prev-page')"
              :disabled="offset === 0"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              <span class="sr-only">Previous</span>
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
            <button
              @click="$emit('next-page')"
              :disabled="!hasMore"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              <span class="sr-only">Next</span>
              <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import PatientTableRow from './PatientTableRow.vue'

// Props with validation
const props = defineProps({
  patients: {
    type: Array,
    required: true,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  totalCount: {
    type: Number,
    required: true,
    default: 0
  },
  offset: {
    type: Number,
    required: true,
    default: 0
  },
  limit: {
    type: Number,
    required: true,
    default: 20
  },
  hasMore: {
    type: Boolean,
    required: true,
    default: false
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
  'retry': () => true,
  'view-patient': (patient) => patient && typeof patient === 'object',
  'edit-patient': (patient) => patient && typeof patient === 'object',
  'delete-patient': (patient) => patient && typeof patient === 'object',
  'prev-page': () => true,
  'next-page': () => true
})
</script>