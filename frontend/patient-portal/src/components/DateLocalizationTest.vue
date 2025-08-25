<template>
  <div class="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow">
    <h2 class="text-2xl font-bold mb-6 text-gray-900">Date Localization Test</h2>
    
    <!-- Locale Selection -->
    <div class="mb-6">
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Current Locale
      </label>
      <select
        v-model="selectedLocale"
        class="mt-1 block w-64 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
        @change="handleLocaleChange"
      >
        <option value="en-MA">English (Morocco)</option>
        <option value="fr-MA">Français (Maroc)</option>
        <option value="en-CA">English (Canada)</option>
        <option value="fr-CA">Français (Canada)</option>
        <option value="en-US">English (USA)</option>
        <option value="fr-FR">Français (France)</option>
      </select>
    </div>

    <!-- Date Input Test -->
    <div class="mb-8">
      <h3 class="text-lg font-semibold mb-4 text-gray-800">Date Input Test</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Enhanced Date Field -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Enhanced Date Field
          </label>
          <DateField
            v-model="testDate"
            :required="true"
            @validation-error="onValidationError"
          />
          <p class="text-xs text-gray-500 mt-1">
            Value: {{ testDate || 'empty' }}
          </p>
        </div>

        <!-- Standard HTML5 Date Input -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Standard HTML5 Date Input
          </label>
          <input
            v-model="standardDate"
            type="date"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          <p class="text-xs text-gray-500 mt-1">
            Value: {{ standardDate || 'empty' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Format Display Test -->
    <div class="mb-8">
      <h3 class="text-lg font-semibold mb-4 text-gray-800">Format Display Test</h3>
      
      <div class="bg-gray-50 p-4 rounded-md">
        <p class="text-sm mb-2"><strong>Test Date:</strong> 1985-06-15</p>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h4 class="font-medium text-gray-700 mb-2">Current Locale Formatting:</h4>
            <ul class="text-sm space-y-1">
              <li><strong>Display Format:</strong> {{ currentDisplayFormat }}</li>
              <li><strong>Formatted:</strong> {{ formatTestDate('1985-06-15') }}</li>
              <li><strong>Placeholder:</strong> {{ currentPlaceholder }}</li>
              <li><strong>Example:</strong> {{ currentExample }}</li>
            </ul>
          </div>
          
          <div>
            <h4 class="font-medium text-gray-700 mb-2">Healthcare Format:</h4>
            <ul class="text-sm space-y-1">
              <li><strong>Morocco:</strong> {{ formatForCountry('1985-06-15', 'MA') }}</li>
              <li><strong>Canada:</strong> {{ formatForCountry('1985-06-15', 'CA') }}</li>
              <li><strong>USA:</strong> {{ formatForCountry('1985-06-15', 'US') }}</li>
              <li><strong>France:</strong> {{ formatForCountry('1985-06-15', 'FR') }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Validation Test -->
    <div class="mb-8">
      <h3 class="text-lg font-semibold mb-4 text-gray-800">Validation Test</h3>
      
      <div class="space-y-4">
        <div v-for="testCase in validationTestCases" :key="testCase.input">
          <div class="flex items-center justify-between p-3 border rounded-md"
               :class="testCase.isValid ? 'border-green-300 bg-green-50' : 'border-red-300 bg-red-50'">
            <span class="font-mono text-sm">{{ testCase.input }}</span>
            <div class="flex items-center space-x-2">
              <span :class="testCase.isValid ? 'text-green-700' : 'text-red-700'">
                {{ testCase.isValid ? '✓' : '✗' }}
              </span>
              <span class="text-xs text-gray-600">
                {{ testCase.error || 'Valid' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Current State Debug -->
    <div class="border-t pt-6">
      <h3 class="text-lg font-semibold mb-4 text-gray-800">Debug Information</h3>
      <div class="bg-gray-100 p-4 rounded-md">
        <pre class="text-xs">{{ debugInfo }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/i18n'
import { useDateLocalization } from '@/composables/useDateLocalization'
import { formatDateForHealthcare } from '@/utils/dateUtils'
import DateField from './fields/DateField.vue'

// Composables
const { locale } = useI18n()
const {
  currentFormat,
  formatDateForDisplay,
  validateDateFormat,
  placeholder: currentPlaceholder,
  example: currentExample,
  displayFormat: currentDisplayFormat
} = useDateLocalization()

// Reactive state
const selectedLocale = ref(locale.value || 'en-MA')
const testDate = ref('1985-06-15')
const standardDate = ref('1985-06-15')
const validationErrors = ref([])

// Methods
const handleLocaleChange = () => {
  setLocale(selectedLocale.value)
}

const onValidationError = (errors) => {
  validationErrors.value = errors
}

const formatTestDate = (dateString) => {
  return formatDateForDisplay(dateString)
}

const formatForCountry = (dateString, country) => {
  return formatDateForHealthcare(dateString, country, selectedLocale.value)
}

// Test cases for validation
const validationTestCases = computed(() => {
  const testInputs = [
    '15/06/1985',
    '06/15/1985', 
    '1985-06-15',
    '31/02/1985', // Invalid date
    '15/13/1985', // Invalid month
    '32/06/1985', // Invalid day
    '15-06-1985', // Wrong separator
    '15.06.1985', // Wrong separator
    '2050/06/15', // Future date
    '1800/06/15', // Too old
    'not-a-date'  // Invalid format
  ]
  
  return testInputs.map(input => {
    const errors = validateDateFormat(input)
    return {
      input,
      isValid: errors.length === 0,
      error: errors.length > 0 ? errors[0] : null
    }
  })
})

// Debug information
const debugInfo = computed(() => {
  return JSON.stringify({
    selectedLocale: selectedLocale.value,
    currentFormat: currentFormat.value,
    testDate: testDate.value,
    standardDate: standardDate.value,
    validationErrors: validationErrors.value,
    placeholder: currentPlaceholder.value,
    example: currentExample.value
  }, null, 2)
})

// Watch locale changes
watch(selectedLocale, (newLocale) => {
  // Locale changed for testing purposes
})
</script>