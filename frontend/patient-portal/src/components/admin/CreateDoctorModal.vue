<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-75 overflow-y-auto h-full w-full z-50" v-if="show">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Add New Doctor</h3>
            <p class="text-sm text-gray-600 mt-1">
              Create a new doctor account with temporary password
            </p>
          </div>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600 focus:outline-none"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- First Name -->
          <div>
            <label for="firstName" class="block text-sm font-medium text-gray-700 mb-1">
              First Name *
            </label>
            <input
              id="firstName"
              v-model="form.first_name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{ 'border-red-300': errors.first_name }"
              placeholder="John"
            />
            <p v-if="errors.first_name" class="text-red-600 text-xs mt-1">
              {{ errors.first_name }}
            </p>
          </div>

          <!-- Last Name -->
          <div>
            <label for="lastName" class="block text-sm font-medium text-gray-700 mb-1">
              Last Name *
            </label>
            <input
              id="lastName"
              v-model="form.last_name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{ 'border-red-300': errors.last_name }"
              placeholder="Smith"
            />
            <p v-if="errors.last_name" class="text-red-600 text-xs mt-1">
              {{ errors.last_name }}
            </p>
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
              Email Address *
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{ 'border-red-300': errors.email }"
              placeholder="john.smith@hospital.com"
            />
            <p v-if="errors.email" class="text-red-600 text-xs mt-1">
              {{ errors.email }}
            </p>
          </div>

          <!-- License Number -->
          <div>
            <label for="licenseNumber" class="block text-sm font-medium text-gray-700 mb-1">
              License Number *
            </label>
            <input
              id="licenseNumber"
              v-model="form.license_number"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{ 'border-red-300': errors.license_number }"
              placeholder="MD123456"
            />
            <p v-if="errors.license_number" class="text-red-600 text-xs mt-1">
              {{ errors.license_number }}
            </p>
          </div>

          <!-- Specialization -->
          <div>
            <label for="specialization" class="block text-sm font-medium text-gray-700 mb-1">
              Specialization *
            </label>
            <select
              id="specialization"
              v-model="form.specialization"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{ 'border-red-300': errors.specialization }"
            >
              <option value="">Select specialization</option>
              <option value="General Practice">General Practice</option>
              <option value="Cardiology">Cardiology</option>
              <option value="Dermatology">Dermatology</option>
              <option value="Emergency Medicine">Emergency Medicine</option>
              <option value="Family Medicine">Family Medicine</option>
              <option value="Internal Medicine">Internal Medicine</option>
              <option value="Neurology">Neurology</option>
              <option value="Oncology">Oncology</option>
              <option value="Orthopedics">Orthopedics</option>
              <option value="Pediatrics">Pediatrics</option>
              <option value="Psychiatry">Psychiatry</option>
              <option value="Radiology">Radiology</option>
              <option value="Surgery">Surgery</option>
              <option value="Other">Other</option>
            </select>
            <p v-if="errors.specialization" class="text-red-600 text-xs mt-1">
              {{ errors.specialization }}
            </p>
          </div>

          <!-- Preferred Language -->
          <div>
            <label for="preferredLanguage" class="block text-sm font-medium text-gray-700 mb-1">
              Preferred Language
            </label>
            <select
              id="preferredLanguage"
              v-model="form.preferred_language"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="en">English</option>
              <option value="fr">French</option>
              <option value="ar">Arabic</option>
            </select>
          </div>

          <!-- Information Note -->
          <div class="bg-blue-50 border border-blue-200 rounded-md p-3">
            <div class="flex">
              <InformationCircleIcon class="h-5 w-5 text-blue-400 mt-0.5 flex-shrink-0" />
              <div class="ml-3">
                <p class="text-sm text-blue-700">
                  <strong>Note:</strong> A temporary password will be generated automatically. 
                  The doctor will be required to change it on first login.
                </p>
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
            <p class="text-sm">{{ error }}</p>
          </div>

          <!-- Submit Button -->
          <div class="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading || !isFormValid"
              class="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <span v-if="loading">Creating Doctor...</span>
              <span v-else>Create Doctor</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, watch } from 'vue'
import { XMarkIcon, InformationCircleIcon } from '@heroicons/vue/24/outline'
import { adminApi } from '@/api/admin'

// Props
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['close', 'success'])

// State
const loading = ref(false)
const error = ref('')

const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
  license_number: '',
  specialization: '',
  preferred_language: 'en'
})

const errors = reactive({
  first_name: '',
  last_name: '',
  email: '',
  license_number: '',
  specialization: ''
})

// Computed
const isFormValid = computed(() => {
  return form.first_name.trim().length > 0 &&
         form.last_name.trim().length > 0 &&
         form.email.trim().length > 0 &&
         form.license_number.trim().length > 0 &&
         form.specialization.trim().length > 0 &&
         isValidEmail(form.email)
})

// Methods
const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

const validateForm = () => {
  // Reset errors
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })

  let isValid = true

  // Validate first name
  if (!form.first_name.trim()) {
    errors.first_name = 'First name is required'
    isValid = false
  }

  // Validate last name
  if (!form.last_name.trim()) {
    errors.last_name = 'Last name is required'
    isValid = false
  }

  // Validate email
  if (!form.email.trim()) {
    errors.email = 'Email is required'
    isValid = false
  } else if (!isValidEmail(form.email)) {
    errors.email = 'Please enter a valid email address'
    isValid = false
  }

  // Validate license number
  if (!form.license_number.trim()) {
    errors.license_number = 'License number is required'
    isValid = false
  }

  // Validate specialization
  if (!form.specialization.trim()) {
    errors.specialization = 'Specialization is required'
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  error.value = ''
  
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    const response = await adminApi.createDoctor({
      email: form.email.trim(),
      first_name: form.first_name.trim(),
      last_name: form.last_name.trim(),
      license_number: form.license_number.trim(),
      specialization: form.specialization.trim(),
      preferred_language: form.preferred_language
    })
    
    // Reset form
    resetForm()
    
    // Emit success with the response data
    emit('success', response)
    
  } catch (err) {
    console.error('Create doctor error:', err)
    error.value = err.response?.data?.error || err.message || 'Failed to create doctor. Please try again.'
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.first_name = ''
  form.last_name = ''
  form.email = ''
  form.license_number = ''
  form.specialization = ''
  form.preferred_language = 'en'
  error.value = ''
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })
}

// Watch for modal close to reset form
watch(() => props.show, (newVal) => {
  if (!newVal) {
    resetForm()
  }
})
</script>

<style scoped>
/* Additional styling if needed */
.transition-colors {
  transition: background-color 0.2s ease-in-out, border-color 0.2s ease-in-out;
}
</style>