<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-75 overflow-y-auto h-full w-full z-50" v-if="show">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Change Password Required</h3>
            <p class="text-sm text-gray-600 mt-1">
              You are using a temporary password. Please change it to continue.
            </p>
          </div>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Current Password -->
          <div>
            <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">
              Current Password
            </label>
            <input
              id="currentPassword"
              v-model="form.currentPassword"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-300': errors.currentPassword }"
            />
            <p v-if="errors.currentPassword" class="text-red-600 text-xs mt-1">
              {{ errors.currentPassword }}
            </p>
          </div>

          <!-- New Password -->
          <div>
            <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">
              New Password
            </label>
            <input
              id="newPassword"
              v-model="form.newPassword"
              type="password"
              required
              autocomplete="new-password"
              minlength="8"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-300': errors.newPassword }"
            />
            <p v-if="errors.newPassword" class="text-red-600 text-xs mt-1">
              {{ errors.newPassword }}
            </p>
            <p class="text-gray-500 text-xs mt-1">
              Minimum 8 characters required
            </p>
          </div>

          <!-- Confirm New Password -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
              Confirm New Password
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              autocomplete="new-password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-300': errors.confirmPassword }"
            />
            <p v-if="errors.confirmPassword" class="text-red-600 text-xs mt-1">
              {{ errors.confirmPassword }}
            </p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
            <p class="text-sm">{{ error }}</p>
          </div>

          <!-- Success Message -->
          <div v-if="success" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md">
            <p class="text-sm">{{ success }}</p>
          </div>

          <!-- Submit Button -->
          <div class="flex justify-end pt-4">
            <button
              type="submit"
              :disabled="loading || !isFormValid"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <span v-if="loading">Changing Password...</span>
              <span v-else>Change Password</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'

// Props
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['close', 'success'])

// Store
const authStore = useAuthStore()

// State
const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const errors = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Computed
const isFormValid = computed(() => {
  return form.currentPassword.length > 0 &&
         form.newPassword.length >= 8 &&
         form.confirmPassword.length > 0 &&
         form.newPassword === form.confirmPassword
})

// Methods
const validateForm = () => {
  // Reset errors
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })

  let isValid = true

  // Validate current password
  if (!form.currentPassword) {
    errors.currentPassword = 'Current password is required'
    isValid = false
  }

  // Validate new password
  if (!form.newPassword) {
    errors.newPassword = 'New password is required'
    isValid = false
  } else if (form.newPassword.length < 8) {
    errors.newPassword = 'Password must be at least 8 characters long'
    isValid = false
  }

  // Validate confirm password
  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your new password'
    isValid = false
  } else if (form.newPassword !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    isValid = false
  }

  // Check if new password is different from current
  if (form.currentPassword === form.newPassword) {
    errors.newPassword = 'New password must be different from current password'
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  error.value = ''
  success.value = ''
  
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    await authStore.changePassword(form.currentPassword, form.newPassword)
    
    success.value = 'Password changed successfully!'
    
    // Reset form
    form.currentPassword = ''
    form.newPassword = ''
    form.confirmPassword = ''
    
    // Close modal after short delay
    setTimeout(() => {
      emit('success')
      emit('close')
    }, 1500)
    
  } catch (err) {
    console.error('Password change error:', err)
    error.value = err.message || 'Failed to change password. Please try again.'
  } finally {
    loading.value = false
  }
}

// Reset form when modal is closed
const resetForm = () => {
  form.currentPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
  error.value = ''
  success.value = ''
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })
}

// Watch for modal close to reset form
import { watch } from 'vue'
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