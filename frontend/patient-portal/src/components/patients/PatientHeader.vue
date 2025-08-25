<template>
  <div class="bg-white shadow rounded-lg p-6 mb-6">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">{{ $t('patients.title') }}</h1>
        <p class="mt-1 text-sm text-gray-500">{{ $t('patients.subtitle') }}</p>
      </div>
      <div class="mt-4 sm:mt-0">
        <button
          @click="$emit('add-patient')"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {{ $t('patients.addPatient') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// Emits with validation
defineEmits({
  'add-patient': () => true
})
</script>