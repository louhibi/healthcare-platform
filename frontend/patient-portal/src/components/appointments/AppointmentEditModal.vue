<template>
  <TransitionRoot as="template" :show="show">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        as="template"
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
      </TransitionChild>

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <TransitionChild
            as="template"
            enter="ease-out duration-300"
            enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to="opacity-100 translate-y-0 sm:scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 translate-y-0 sm:scale-100"
            leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <DialogPanel class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-2xl sm:p-6">
              <div>
                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
                  <PencilIcon class="h-6 w-6 text-green-600" />
                </div>
                <div class="mt-3 text-center sm:mt-5">
                  <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                    Edit Appointment
                  </DialogTitle>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500">
                      Update appointment information for {{ appointment?.patient_name || 'Patient' }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- Configuration Loading State -->
              <div v-if="appointmentFormConfig.isLoading.value" class="text-center py-8">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
                <p class="mt-4 text-sm text-gray-500">Loading form configuration...</p>
              </div>

              <!-- Configuration Error State -->
              <div v-else-if="appointmentFormConfig.error.value" class="text-center py-8">
                <div class="rounded-md bg-red-50 p-4">
                  <div class="flex">
                    <div class="flex-shrink-0">
                      <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                      </svg>
                    </div>
                    <div class="ml-3">
                      <h3 class="text-sm font-medium text-red-800">Configuration Error</h3>
                      <div class="mt-2 text-sm text-red-700">
                        <p>{{ appointmentFormConfig.error.value }}</p>
                      </div>
                      <div class="mt-3">
                        <button @click="retryFormConfiguration" class="bg-red-100 px-2 py-1 text-xs font-medium text-red-800 rounded hover:bg-red-200">
                          Retry Loading Configuration
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Dynamic Form - Only shown when configuration is loaded -->
              <form v-else-if="appointmentFormConfig.enabledFields.value.length > 0" @submit.prevent="updateAppointment" class="mt-6">
                <!-- Appointment fields in flexible layout -->
                <div class="space-y-6">
                  <!-- Core appointment fields -->
                  <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                    <!-- Patient Selection - special handling -->
                    <div v-if="appointmentFormConfig.isFieldEnabled('patient_id')">
                      <label for="patient-select" class="block text-sm font-medium text-gray-700">
                        {{ appointmentFormConfig.getFieldDisplayName('patient_id') }}
                        <span v-if="appointmentFormConfig.isFieldRequired('patient_id')" class="text-red-500 ml-1">*</span>
                      </label>
                      <select
                        id="patient-select"
                        v-model="form.patient_id"
                        :required="appointmentFormConfig.isFieldRequired('patient_id')"
                        :class="[
                          'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                          appointmentFormConfig.hasFieldError('patient_id') ? 'border-red-500' : ''
                        ]"
                      >
                        <option value="">Select Patient</option>
                        <option v-for="patient in patients" :key="patient.id" :value="patient.id">
                          {{ patient.first_name }} {{ patient.last_name }}
                        </option>
                      </select>
                      <div v-if="appointmentFormConfig.hasFieldError('patient_id')" class="mt-1">
                        <p v-for="error in appointmentFormConfig.getFieldError('patient_id')" :key="error" class="text-xs text-red-600">
                          {{ error }}
                        </p>
                      </div>
                    </div>

                    <!-- Doctor Selection - special handling -->
                    <div v-if="appointmentFormConfig.isFieldEnabled('doctor_id')">
                      <label for="doctor-select" class="block text-sm font-medium text-gray-700">
                        {{ appointmentFormConfig.getFieldDisplayName('doctor_id') }}
                        <span v-if="appointmentFormConfig.isFieldRequired('doctor_id')" class="text-red-500 ml-1">*</span>
                      </label>
                      <select
                        id="doctor-select"
                        v-model="form.doctor_id"
                        :required="appointmentFormConfig.isFieldRequired('doctor_id')"
                        :class="[
                          'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                          appointmentFormConfig.hasFieldError('doctor_id') ? 'border-red-500' : ''
                        ]"
                      >
                        <option value="">Select Doctor</option>
                        <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                          Dr. {{ doctor.first_name }} {{ doctor.last_name }} - {{ doctor.specialization }}
                        </option>
                      </select>
                      <div v-if="appointmentFormConfig.hasFieldError('doctor_id')" class="mt-1">
                        <p v-for="error in appointmentFormConfig.getFieldError('doctor_id')" :key="error" class="text-xs text-red-600">
                          {{ error }}
                        </p>
                      </div>
                    </div>
                  </div>
                  
                  <!-- Dynamic form fields by category -->
                  <div v-for="(categoryFields, categoryName) in appointmentFormConfig.fieldsByCategory.value" :key="categoryName" class="space-y-4">
                    <div v-if="categoryFields.filter(f => f.is_enabled && !['patient_id', 'doctor_id'].includes(f.name)).length > 0" class="border-t pt-4 first:border-t-0 first:pt-0">
                      <h4 v-if="categoryName !== 'Other'" class="text-md font-medium text-gray-900 mb-3">{{ categoryName }}</h4>
                      
                      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                        <template v-for="field in categoryFields" :key="field.name">
                          <div v-if="field.is_enabled && !['patient_id', 'doctor_id'].includes(field.name)" :class="field.field_type === 'textarea' ? 'sm:col-span-2' : ''">
                            <label :for="'appointment-' + field.name" class="block text-sm font-medium text-gray-700">
                              {{ field.display_name }}
                              <span v-if="field.is_required" class="text-red-500 ml-1">*</span>
                            </label>
                            
                            <!-- Text input fields -->
                            <input
                              v-if="['text', 'email', 'tel', 'url'].includes(field.field_type)"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              :type="field.field_type"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                              :placeholder="field.description || ''"
                            />
                            
                            <!-- Special handling for date fields -->
                            <input
                              v-else-if="field.field_type === 'date'"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              type="date"
                              :min="field.name === 'date' ? today : ''"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                            />
                            
                            <!-- Special handling for time fields -->
                            <input
                              v-else-if="field.field_type === 'time'"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              type="time"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                            />
                            
                            <!-- Textarea fields -->
                            <textarea
                              v-else-if="field.field_type === 'textarea'"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              :required="field.is_required"
                              rows="3"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                              :placeholder="field.description || ''"
                            ></textarea>
                            
                            <!-- Select fields -->
                            <select
                              v-else-if="field.field_type === 'select'"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                            >
                              <option value="">Select {{ field.display_name }}</option>
                              <option v-for="option in field.options" :key="option" :value="option">
                                {{ option }}
                              </option>
                            </select>
                            
                            <!-- Number fields -->
                            <input
                              v-else-if="field.field_type === 'number'"
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              type="number"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                            />
                            
                            <!-- Boolean/checkbox fields -->
                            <div v-else-if="field.field_type === 'boolean'" class="mt-1">
                              <label class="inline-flex items-center">
                                <input
                                  :id="'appointment-' + field.name"
                                  v-model="form[field.name]"
                                  type="checkbox"
                                  class="form-checkbox h-4 w-4 text-indigo-600"
                                />
                                <span class="ml-2 text-sm text-gray-600">{{ field.description || field.display_name }}</span>
                              </label>
                            </div>
                            
                            <!-- Fallback for other field types -->
                            <input
                              v-else
                              :id="'appointment-' + field.name"
                              v-model="form[field.name]"
                              type="text"
                              :required="field.is_required"
                              :class="[
                                'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                              ]"
                              :placeholder="field.description || ''"
                            />
                            
                            <!-- Field errors -->
                            <div v-if="appointmentFormConfig.hasFieldError(field.name)" class="mt-1">
                              <p v-for="error in appointmentFormConfig.getFieldError(field.name)" :key="error" class="text-xs text-red-600">
                                {{ error }}
                              </p>
                            </div>
                            
                            <!-- Field description -->
                            <p v-if="field.description && field.field_type !== 'boolean'" class="mt-1 text-xs text-gray-500">
                              {{ field.description }}
                            </p>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="mt-6 sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    :disabled="loading"
                    class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:ml-3 sm:w-auto disabled:opacity-50"
                  >
                    <span v-if="loading" class="inline-flex items-center">
                      <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Updating...
                    </span>
                    <span v-else>Update Appointment</span>
                  </button>
                  
                  <button
                    @click="close"
                    type="button"
                    class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                  >
                    Cancel
                  </button>
                </div>
              </form>

              <!-- No Configuration Available State -->
              <div v-else class="text-center py-8">
                <div class="rounded-md bg-yellow-50 p-4">
                  <div class="flex">
                    <div class="flex-shrink-0">
                      <svg class="h-5 w-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                      </svg>
                    </div>
                    <div class="ml-3">
                      <h3 class="text-sm font-medium text-yellow-800">No Form Configuration</h3>
                      <div class="mt-2 text-sm text-yellow-700">
                        <p>Appointment form configuration is not available. Please contact your administrator.</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Error Message -->
              <div v-if="error" class="mt-6 rounded-md bg-red-50 p-4">
                <div class="flex">
                  <ExclamationCircleIcon class="h-5 w-5 text-red-400" />
                  <div class="ml-3">
                    <h3 class="text-sm font-medium text-red-800">
                      Error updating appointment
                    </h3>
                    <div class="mt-2 text-sm text-red-700">
                      <p>{{ error }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { PencilIcon, ExclamationCircleIcon } from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'
import { patientsApi } from '@/api/patients'
import { getTodayString, formatDateForInput, parseDate } from '@/utils/dateUtils'
import { createUTCDateTime } from '@/utils/timezoneUtils'
import { useAuthStore } from '@/stores/auth'
import { useEntityStore } from '@/stores/entity'
import { useFormConfig } from '@/composables/useFormConfig'

// Props & Emits
const props = defineProps({
  show: {
    type: Boolean,
    default: true
  },
  appointment: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'appointment-updated'])

// Auth store
const authStore = useAuthStore()
const entityStore = useEntityStore()

// Form configuration
const appointmentFormConfig = useFormConfig('appointment')

// Reactive state
const patients = ref([])
const doctors = ref([])
const loading = ref(false)
const error = ref('')

const form = reactive({})

// Initialize form data based on configuration - NO FALLBACKS
const initializeForm = (appointmentData = {}) => {
  // Clear existing form data
  Object.keys(form).forEach(key => {
    delete form[key]
  })
  
  // Only initialize if configuration is loaded
  if (appointmentFormConfig.enabledFields.value.length > 0) {
    const initialData = appointmentFormConfig.initializeFormData(appointmentData)
    Object.assign(form, initialData)
  }
  // NO FALLBACK - form will not work without configuration
}

// Retry form configuration loading
const retryFormConfiguration = async () => {
  try {
    await appointmentFormConfig.initialize()
    populateForm()
  } catch (err) {
    console.error('Failed to retry appointment form configuration:', err)
  }
}

// Computed
const today = computed(() => getTodayString())

// Methods
const close = () => {
  emit('close')
}

const updateAppointment = async () => {
  try {
    loading.value = true
    error.value = ''
    
    // Configuration MUST be available - NO FALLBACKS
    if (appointmentFormConfig.enabledFields.value.length === 0) {
      error.value = 'Form configuration is not available. Please try refreshing the page.'
      return
    }
    
    // Validate using dynamic form configuration
    const isValid = appointmentFormConfig.validateForm(form)
    if (!isValid) {
      const errors = Object.values(appointmentFormConfig.validationErrors.value).flat()
      error.value = `Please fix the following errors: ${errors.join(', ')}`
      return
    }
    
    // Get healthcare entity timezone
    const entityTimezone = entityStore.entityTimezone
    if (!entityTimezone) {
      throw new Error('Healthcare entity timezone not available')
    }
    
    // Properly convert entity time to UTC
    const dateTime = createUTCDateTime(form.date, form.time, entityTimezone)
    
    const updateData = {
      patient_id: parseInt(form.patient_id),
      doctor_id: parseInt(form.doctor_id),
      date_time: dateTime,
      duration: parseInt(form.duration),
      type: form.type,
      reason: form.reason,
      notes: form.notes,
      priority: form.priority,
      room_number: form.room_number
    }
    
    await appointmentsApi.updateAppointment(props.appointment.id, updateData)
    emit('appointment-updated')
  } catch (err) {
    error.value = err.message || 'Failed to update appointment'
    console.error('Update appointment error:', err)
  } finally {
    loading.value = false
  }
}

const loadPatients = async () => {
  try {
    const response = await patientsApi.getPatients({ limit: 100 })
    patients.value = response.data.patients || []
  } catch (err) {
    console.error('Load patients error:', err)
  }
}

const loadDoctors = async () => {
  try {
    const response = await appointmentsApi.getDoctorsByEntity()
    doctors.value = response.data || []
  } catch (err) {
    console.error('Load doctors error:', err)
  }
}

const populateForm = () => {
  if (props.appointment && appointmentFormConfig.enabledFields.value.length > 0) {
    // Parse the appointment date_time properly
    const appointmentDateTime = new Date(props.appointment.date_time)
    
    const appointmentData = {
      patient_id: props.appointment.patient_id,
      doctor_id: props.appointment.doctor_id,
      date: formatDateForInput(appointmentDateTime),
      time: appointmentDateTime.toTimeString().slice(0, 5),
      duration: props.appointment.duration,
      type: props.appointment.type,
      reason: props.appointment.reason,
      notes: props.appointment.notes || '',
      priority: props.appointment.priority || 'normal',
      room_number: props.appointment.room_number || ''
    }
    
    initializeForm(appointmentData)
    appointmentFormConfig.clearAllErrors()
  }
}

// Watchers
watch(() => props.appointment, populateForm, { immediate: true })

// Lifecycle
onMounted(async () => {
  try {
    // Initialize form configuration first - REQUIRED
    await appointmentFormConfig.initialize()
  } catch (err) {
    console.error('Failed to load appointment form configuration:', err)
    // DO NOT initialize form without configuration
  }
  
  await Promise.all([loadPatients(), loadDoctors()])
  populateForm()
})
</script>