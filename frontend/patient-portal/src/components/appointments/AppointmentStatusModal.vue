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
            <DialogPanel class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
              <div>
                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
                  <CheckIcon class="h-6 w-6 text-blue-600" />
                </div>
                <div class="mt-3 text-center sm:mt-5">
                  <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                    Update Appointment Status
                  </DialogTitle>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500">
                      Change the status for {{ appointment?.patient_name || 'Patient' }}'s appointment
                    </p>
                  </div>
                </div>
              </div>

              <form @submit.prevent="updateStatus" class="mt-6">
                <!-- Current Status Display -->
                <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                  <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-gray-700">Current Status:</span>
                    <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                          :class="getStatusClass(appointment?.status)">
                      {{ appointment?.status }}
                    </span>
                  </div>
                  <div class="mt-2 text-xs text-gray-500">
                    {{ formatDateTime(appointment?.updated_at) }}
                  </div>
                </div>

                <!-- New Status Selection -->
                <div class="mb-4">
                  <label for="new-status" class="block text-sm font-medium text-gray-700">New Status *</label>
                  <select
                    id="new-status"
                    v-model="form.status"
                    required
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  >
                    <option value="">Select Status</option>
                    <option 
                      v-for="status in availableStatuses" 
                      :key="status.value" 
                      :value="status.value"
                      :disabled="status.value === appointment?.status"
                    >
                      {{ status.label }}
                      <span v-if="status.value === appointment?.status" class="text-gray-400">(Current)</span>
                    </option>
                  </select>
                </div>

                <!-- Status-specific Messages -->
                <div v-if="form.status" class="mb-4 p-3 rounded-md" :class="getStatusMessageClass(form.status)">
                  <div class="flex">
                    <component :is="getStatusIcon(form.status)" class="h-5 w-5" :class="getStatusIconClass(form.status)" />
                    <div class="ml-3">
                      <p class="text-sm" :class="getStatusTextClass(form.status)">
                        {{ getStatusMessage(form.status) }}
                      </p>
                    </div>
                  </div>
                </div>

                <!-- Notes -->
                <div class="mb-6">
                  <label for="status-notes" class="block text-sm font-medium text-gray-700">
                    Status Update Notes
                    <span v-if="requiresNotes(form.status)" class="text-red-500">*</span>
                  </label>
                  <textarea
                    id="status-notes"
                    v-model="form.notes"
                    :required="requiresNotes(form.status)"
                    rows="3"
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    :placeholder="getNotesPlaceholder(form.status)"
                  />
                </div>

                <!-- Quick Actions -->
                <div v-if="form.status" class="mb-4">
                  <h4 class="text-sm font-medium text-gray-700 mb-2">Quick Actions</h4>
                  <div class="space-y-2">
                    <button
                      v-for="action in getQuickActions(form.status)"
                      :key="action.text"
                      @click="form.notes = action.text"
                      type="button"
                      class="w-full text-left px-3 py-2 text-sm text-gray-700 bg-gray-50 hover:bg-gray-100 rounded-md border border-gray-200"
                    >
                      {{ action.text }}
                    </button>
                  </div>
                </div>

                <!-- Error Message -->
                <div v-if="error" class="mb-4 rounded-md bg-red-50 p-4">
                  <div class="flex">
                    <ExclamationCircleIcon class="h-5 w-5 text-red-400" />
                    <div class="ml-3">
                      <h3 class="text-sm font-medium text-red-800">
                        Error updating status
                      </h3>
                      <div class="mt-2 text-sm text-red-700">
                        <p>{{ error }}</p>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    :disabled="loading || !form.status"
                    class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:ml-3 sm:w-auto disabled:opacity-50"
                  >
                    <span v-if="loading" class="inline-flex items-center">
                      <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Updating...
                    </span>
                    <span v-else>Update Status</span>
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
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
  CheckIcon,
  ExclamationCircleIcon,
  CheckCircleIcon,
  ClockIcon,
  XMarkIcon,
  PlayIcon,
  StopIcon
} from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'

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

const emit = defineEmits(['close', 'status-updated'])

// Reactive state
const loading = ref(false)
const error = ref('')

const form = reactive({
  status: '',
  notes: ''
})

// Computed
const availableStatuses = computed(() => [
  { value: 'scheduled', label: 'Scheduled' },
  { value: 'confirmed', label: 'Confirmed' },
  { value: 'in-progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
  { value: 'cancelled', label: 'Cancelled' },
  { value: 'no-show', label: 'No Show' }
])

// Methods
const close = () => {
  emit('close')
}

const updateStatus = async () => {
  try {
    loading.value = true
    error.value = ''
    
    await appointmentsApi.updateAppointmentStatus(
      props.appointment.id, 
      form.status, 
      form.notes
    )
    
    emit('status-updated')
  } catch (err) {
    error.value = err.message || 'Failed to update appointment status'
    console.error('Update status error:', err)
  } finally {
    loading.value = false
  }
}

// Utility functions
const formatDateTime = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString()
}

const getStatusClass = (status) => {
  const classes = {
    scheduled: 'bg-yellow-100 text-yellow-800',
    confirmed: 'bg-blue-100 text-blue-800',
    'in-progress': 'bg-green-100 text-green-800',
    completed: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-red-100 text-red-800',
    'no-show': 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusIcon = (status) => {
  const icons = {
    scheduled: ClockIcon,
    confirmed: CheckCircleIcon,
    'in-progress': PlayIcon,
    completed: CheckCircleIcon,
    cancelled: XMarkIcon,
    'no-show': StopIcon
  }
  return icons[status] || CheckIcon
}

const getStatusIconClass = (status) => {
  const classes = {
    scheduled: 'text-yellow-400',
    confirmed: 'text-blue-400',
    'in-progress': 'text-green-400',
    completed: 'text-gray-400',
    cancelled: 'text-red-400',
    'no-show': 'text-orange-400'
  }
  return classes[status] || 'text-gray-400'
}

const getStatusMessageClass = (status) => {
  const classes = {
    scheduled: 'bg-yellow-50',
    confirmed: 'bg-blue-50',
    'in-progress': 'bg-green-50',
    completed: 'bg-gray-50',
    cancelled: 'bg-red-50',
    'no-show': 'bg-orange-50'
  }
  return classes[status] || 'bg-gray-50'
}

const getStatusTextClass = (status) => {
  const classes = {
    scheduled: 'text-yellow-700',
    confirmed: 'text-blue-700',
    'in-progress': 'text-green-700',
    completed: 'text-gray-700',
    cancelled: 'text-red-700',
    'no-show': 'text-orange-700'
  }
  return classes[status] || 'text-gray-700'
}

const getStatusMessage = (status) => {
  const messages = {
    scheduled: 'Appointment is scheduled and awaiting confirmation.',
    confirmed: 'Appointment is confirmed and ready to proceed.',
    'in-progress': 'Appointment is currently in progress.',
    completed: 'Appointment has been completed successfully.',
    cancelled: 'Appointment has been cancelled. Please provide a reason.',
    'no-show': 'Patient did not show up for the appointment.'
  }
  return messages[status] || ''
}

const requiresNotes = (status) => {
  return ['cancelled', 'no-show'].includes(status)
}

const getNotesPlaceholder = (status) => {
  const placeholders = {
    scheduled: 'Optional notes about scheduling...',
    confirmed: 'Confirmation details or instructions...',
    'in-progress': 'Notes about the ongoing appointment...',
    completed: 'Summary of the completed appointment...',
    cancelled: 'Required: Reason for cancellation...',
    'no-show': 'Details about the no-show...'
  }
  return placeholders[status] || 'Additional notes...'
}

const getQuickActions = (status) => {
  const actions = {
    confirmed: [
      { text: 'Patient confirmed via phone' },
      { text: 'Reminder sent to patient' },
      { text: 'Insurance verified' }
    ],
    'in-progress': [
      { text: 'Patient checked in' },
      { text: 'Vitals taken' },
      { text: 'Moved to examination room' }
    ],
    completed: [
      { text: 'Appointment completed successfully' },
      { text: 'Follow-up required' },
      { text: 'Prescription provided' },
      { text: 'Tests ordered' }
    ],
    cancelled: [
      { text: 'Patient requested cancellation' },
      { text: 'Doctor unavailable' },
      { text: 'Emergency situation' },
      { text: 'Equipment malfunction' }
    ],
    'no-show': [
      { text: 'Patient did not arrive' },
      { text: 'No response to calls' },
      { text: 'Waited 15 minutes past appointment time' }
    ]
  }
  return actions[status] || []
}
</script>