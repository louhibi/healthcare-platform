<template>
  <div>
    <!-- Header with Create Button -->
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-lg font-medium text-gray-900">Doctor Management</h2>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        Add Doctor
      </button>
    </div>

    <!-- Doctors List -->
    <div class="bg-white shadow overflow-hidden sm:rounded-md">
      <div v-if="loading" class="px-4 py-12 text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-sm text-gray-500">Loading doctors...</p>
      </div>

      <div v-else-if="error" class="px-4 py-12 text-center">
        <div class="text-red-600">
          <ExclamationTriangleIcon class="h-8 w-8 mx-auto mb-2" />
          <p class="text-sm">{{ error }}</p>
          <button 
            @click="loadDoctors" 
            class="mt-2 text-indigo-600 hover:text-indigo-500 text-sm font-medium"
          >
            Try again
          </button>
        </div>
      </div>

      <ul v-else-if="doctors.length > 0" class="divide-y divide-gray-200">
        <li v-for="doctor in doctors" :key="doctor.id" class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="flex-shrink-0 h-12 w-12">
                <div class="h-12 w-12 bg-gray-300 rounded-full flex items-center justify-center">
                  <UserIcon class="h-6 w-6 text-gray-600" />
                </div>
              </div>
              <div class="ml-4">
                <div class="flex items-center">
                  <p class="text-sm font-medium text-gray-900">
                    Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                  </p>
                  <span v-if="doctor.is_temp_password" 
                        class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    <ClockIcon class="h-3 w-3 mr-1" />
                    Temporary Password
                  </span>
                </div>
                <p class="text-sm text-gray-500">{{ doctor.email }}</p>
                <div class="flex items-center text-xs text-gray-500 mt-1">
                  <span>{{ doctor.specialization }}</span>
                  <span class="mx-2">•</span>
                  <span>License: {{ doctor.license_number }}</span>
                  <span class="mx-2">•</span>
                  <span>Language: {{ doctor.preferred_language?.toUpperCase() || 'EN' }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-3">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                Active
              </span>
              <div class="flex space-x-2">
                <button 
                  @click="viewDoctor(doctor)"
                  class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
                >
                  View
                </button>
                <button 
                  @click="editDoctor(doctor)"
                  class="text-gray-600 hover:text-gray-900 text-sm font-medium"
                >
                  Edit
                </button>
              </div>
            </div>
          </div>
        </li>
      </ul>

      <div v-else class="px-4 py-12 text-center">
        <UserIcon class="h-12 w-12 mx-auto text-gray-400" />
        <p class="mt-2 text-sm text-gray-500">No doctors found</p>
        <p class="text-xs text-gray-400">Add your first doctor to get started</p>
      </div>
    </div>

    <!-- Create Doctor Modal -->
    <CreateDoctorModal 
      :show="showCreateModal"
      @close="showCreateModal = false"
      @success="handleDoctorCreated"
    />

    <!-- Success Modal for Showing Temp Password -->
    <div v-if="showSuccessModal" class="fixed inset-0 bg-gray-600 bg-opacity-75 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <!-- Success Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100">
            <CheckCircleIcon class="h-6 w-6 text-green-600" />
          </div>
          
          <!-- Content -->
          <div class="mt-5 text-center">
            <h3 class="text-lg font-medium text-gray-900">Doctor Created Successfully!</h3>
            <div class="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
              <p class="text-sm text-gray-700 mb-2">
                <strong>{{ newDoctorInfo.doctor.first_name }} {{ newDoctorInfo.doctor.last_name }}</strong> has been added to your team.
              </p>
              <div class="text-left space-y-2">
                <div class="text-sm">
                  <span class="font-medium">Email:</span> {{ newDoctorInfo.doctor.email }}
                </div>
                <div class="text-sm">
                  <span class="font-medium">Temporary Password:</span>
                  <code class="bg-gray-100 px-2 py-1 rounded text-red-600 font-mono">{{ newDoctorInfo.temp_password }}</code>
                </div>
                <div class="text-sm text-gray-600">
                  <span class="font-medium">Expires:</span> {{ formatDate(newDoctorInfo.password_expires) }}
                </div>
              </div>
              <div class="mt-3 p-2 bg-blue-50 border border-blue-200 rounded text-xs text-blue-700">
                <strong>Important:</strong> Please provide these credentials to the doctor. They will be required to change the password on first login.
              </div>
            </div>
            
            <!-- Actions -->
            <div class="mt-6 flex justify-end space-x-3">
              <button
                @click="copyCredentials"
                class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <ClipboardIcon class="h-4 w-4 mr-2" />
                Copy Credentials
              </button>
              <button
                @click="closeSuccessModal"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Done
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import {
  PlusIcon,
  UserIcon,
  ClockIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  ClipboardIcon
} from '@heroicons/vue/24/outline'
import CreateDoctorModal from './CreateDoctorModal.vue'
import { adminApi } from '@/api/admin'

// Props
const emit = defineEmits(['stats-updated'])

// State
const doctors = ref([])
const loading = ref(false)
const error = ref('')
const showCreateModal = ref(false)
const showSuccessModal = ref(false)
const newDoctorInfo = ref(null)
const toast = useToast()

// Methods
const loadDoctors = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await adminApi.getDoctors()
    doctors.value = response.doctors || response || []
  } catch (err) {
    console.error('Error loading doctors:', err)
    error.value = 'Failed to load doctors. Please try again.'
    doctors.value = []
  } finally {
    loading.value = false
  }
}

const handleDoctorCreated = (doctorData) => {
  showCreateModal.value = false
  newDoctorInfo.value = doctorData
  showSuccessModal.value = true
  
  // Add the new doctor to the list
  doctors.value.unshift(doctorData.doctor)
  
  // Emit stats update
  emit('stats-updated')
  
  toast.success('Doctor created successfully!')
}

const closeSuccessModal = () => {
  showSuccessModal.value = false
  newDoctorInfo.value = null
}

const copyCredentials = () => {
  if (!newDoctorInfo.value) return
  
  const credentials = `Healthcare Portal - Doctor Account Created
  
Name: Dr. ${newDoctorInfo.value.doctor.first_name} ${newDoctorInfo.value.doctor.last_name}
Email: ${newDoctorInfo.value.doctor.email}
Temporary Password: ${newDoctorInfo.value.temp_password}
Password Expires: ${formatDate(newDoctorInfo.value.password_expires)}

Please log in at the Healthcare Portal and change your password on first login.`

  navigator.clipboard.writeText(credentials).then(() => {
    toast.success('Credentials copied to clipboard!')
  }).catch(() => {
    toast.error('Failed to copy credentials')
  })
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const viewDoctor = (doctor) => {
  // TODO: Implement view doctor details
  toast.info('View doctor details - Coming soon!')
}

const editDoctor = (doctor) => {
  // TODO: Implement edit doctor
  toast.info('Edit doctor - Coming soon!')
}

// Lifecycle
onMounted(() => {
  loadDoctors()
})
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>