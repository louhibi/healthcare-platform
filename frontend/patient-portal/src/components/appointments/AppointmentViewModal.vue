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
                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
                  <EyeIcon class="h-6 w-6 text-blue-600" />
                </div>
                <div class="mt-3 text-center sm:mt-5">
                  <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                    Appointment Details
                  </DialogTitle>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500">
                      Complete information for appointment #{{ appointment?.id }}
                    </p>
                  </div>
                </div>
              </div>

              <div v-if="appointment" class="mt-6">
                <!-- Status Badge -->
                <div class="mb-6 flex justify-center">
                  <span class="inline-flex items-center rounded-full px-3 py-1 text-sm font-medium"
                        :class="getStatusClass(appointment.status)">
                    {{ appointment.status }}
                  </span>
                </div>

                <!-- Appointment Information -->
                <div class="border-t border-gray-200 pt-6">
                  <dl class="divide-y divide-gray-200">
                    <!-- Patient Information -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Patient</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <div class="flex items-center">
                          <div class="h-10 w-10 flex-shrink-0">
                            <div class="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
                              <UserIcon class="h-6 w-6 text-gray-400" />
                            </div>
                          </div>
                          <div class="ml-4">
                            <div class="font-medium">{{ appointment.patient_name || 'Patient' }}</div>
                            <div class="text-gray-500">ID: {{ appointment.patient_id }}</div>
                          </div>
                        </div>
                      </dd>
                    </div>

                    <!-- Doctor Information -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Doctor</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <div class="flex items-center">
                          <div class="h-10 w-10 flex-shrink-0">
                            <div class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                              <UserIcon class="h-6 w-6 text-blue-600" />
                            </div>
                          </div>
                          <div class="ml-4">
                            <div class="font-medium">{{ appointment.doctor_name || 'Doctor' }}</div>
                            <div class="text-gray-500">ID: {{ appointment.doctor_id }}</div>
                          </div>
                        </div>
                      </dd>
                    </div>

                    <!-- Date & Time -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Date & Time</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <div class="flex items-center space-x-2">
                          <CalendarDaysIcon class="h-5 w-5 text-gray-400" />
                          <span class="font-medium">{{ formatDate(appointment.date_time) }}</span>
                          <span class="text-gray-500">at</span>
                          <span class="font-medium">{{ formatTime(appointment.date_time) }}</span>
                        </div>
                        <div class="mt-1 text-gray-500 text-xs">
                          Duration: {{ appointment.duration }} minutes
                          <span v-if="appointment.end_time">
                            (ends at {{ formatTime(appointment.end_time) }})
                          </span>
                        </div>
                      </dd>
                    </div>

                    <!-- Type -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Appointment Type</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                              :class="getTypeClass(appointment.type)">
                          {{ appointment.type }}
                        </span>
                      </dd>
                    </div>

                    <!-- Priority -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Priority</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                              :class="getPriorityClass(appointment.priority)">
                          {{ appointment.priority || 'normal' }}
                        </span>
                      </dd>
                    </div>

                    <!-- Room -->
                    <div v-if="appointment.room_number" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Room</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        <div class="flex items-center">
                          <BuildingOfficeIcon class="h-5 w-5 text-gray-400 mr-2" />
                          {{ appointment.room_number || 'Not assigned' }}
                        </div>
                      </dd>
                    </div>

                    <!-- Reason -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Reason for Visit</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        {{ appointment.reason }}
                      </dd>
                    </div>

                    <!-- Notes -->
                    <div v-if="appointment.notes" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Notes</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        {{ appointment.notes }}
                      </dd>
                    </div>

                    <!-- Timestamps -->
                    <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Created</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        {{ formatDateTime(appointment.created_at) }}
                      </dd>
                    </div>

                    <div v-if="appointment.updated_at !== appointment.created_at" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:py-5">
                      <dt class="text-sm font-medium text-gray-500">Last Updated</dt>
                      <dd class="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                        {{ formatDateTime(appointment.updated_at) }}
                      </dd>
                    </div>
                  </dl>
                </div>
              </div>

              <div class="mt-6 sm:flex sm:flex-row-reverse">
                <button
                  @click="close"
                  type="button"
                  class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:ml-3 sm:w-auto"
                >
                  Close
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { formatEntityDateTime, formatEntityTime, formatEntityDate } from '@/utils/timezoneUtils'
import { useAuthStore } from '@/stores/auth'
import { useEntityStore } from '@/stores/entity'
import {
  EyeIcon,
  UserIcon,
  CalendarDaysIcon,
  BuildingOfficeIcon
} from '@heroicons/vue/24/outline'

// Props & Emits
defineProps({
  show: {
    type: Boolean,
    default: true
  },
  appointment: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

// Auth store
const authStore = useAuthStore()
const entityStore = useEntityStore()

// Methods
const close = () => {
  emit('close')
}

// Utility functions with timezone awareness
const formatDate = (utcDateTimeString) => {
  if (!utcDateTimeString) return ''
  const entityTimezone = entityStore.entityTimezone
  if (!entityTimezone) {
    return new Date(utcDateTimeString).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
  
  try {
    return formatEntityDate(utcDateTimeString, entityTimezone, {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch (error) {
    console.warn('Failed to format date with timezone:', error)
    return new Date(utcDateTimeString).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
}

const formatTime = (utcDateTimeString) => {
  if (!utcDateTimeString) return ''
  const entityTimezone = entityStore.entityTimezone
  if (!entityTimezone) return new Date(utcDateTimeString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  
  try {
    return formatEntityTime(utcDateTimeString, entityTimezone, true)
  } catch (error) {
    console.warn('Failed to format time with timezone:', error)
    return new Date(utcDateTimeString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
}

const formatDateTime = (utcDateTimeString) => {
  if (!utcDateTimeString) return ''
  const entityTimezone = entityStore.entityTimezone
  if (!entityTimezone) return new Date(utcDateTimeString).toLocaleString()
  
  try {
    return formatEntityDateTime(utcDateTimeString, entityTimezone)
  } catch (error) {
    console.warn('Failed to format datetime with timezone:', error)
    return new Date(utcDateTimeString).toLocaleString()
  }
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

const getTypeClass = (type) => {
  const classes = {
    consultation: 'bg-blue-100 text-blue-800',
    'follow-up': 'bg-green-100 text-green-800',
    procedure: 'bg-purple-100 text-purple-800',
    emergency: 'bg-red-100 text-red-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getPriorityClass = (priority) => {
  const classes = {
    low: 'bg-gray-100 text-gray-800',
    normal: 'bg-blue-100 text-blue-800',
    high: 'bg-orange-100 text-orange-800',
    urgent: 'bg-red-100 text-red-800'
  }
  return classes[priority] || classes.normal
}
</script>