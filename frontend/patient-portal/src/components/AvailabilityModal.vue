<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="$emit('close')">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-2/3 lg:w-1/2 shadow-lg rounded-md bg-white" @click.stop>
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold text-gray-900">
            {{ availability?.id ? 'Edit Availability' : 'Create Availability' }}
          </h3>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Doctor and Date -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Doctor *</label>
              <select
                v-model="form.doctor_id"
                required
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Select Doctor</option>
                <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                  Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                  <span v-if="doctor.specialization">({{ doctor.specialization }})</span>
                </option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Date *</label>
              <input
                v-model="form.date"
                type="date"
                required
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>

          <!-- Status -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Status *</label>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
              <label
                v-for="status in statusOptions"
                :key="status.value"
                :class="[
                  'relative flex cursor-pointer rounded-lg border p-4 focus:outline-none',
                  form.status === status.value
                    ? 'border-blue-600 ring-2 ring-blue-600'
                    : 'border-gray-300'
                ]"
              >
                <input
                  v-model="form.status"
                  type="radio"
                  :value="status.value"
                  class="sr-only"
                />
                <div class="flex items-center justify-between w-full">
                  <div class="flex items-center">
                    <div class="text-2xl mr-2">{{ status.emoji }}</div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">{{ status.label }}</div>
                      <div class="text-xs text-gray-500">{{ status.description }}</div>
                    </div>
                  </div>
                  <div v-if="form.status === status.value" class="text-blue-600">
                    <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <!-- Working Hours (only show if status is available) -->
          <div v-if="form.status === 'available'" class="space-y-4">
            <h4 class="text-lg font-medium text-gray-900 border-b pb-2">Working Hours</h4>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Start Time</label>
                <input
                  v-model="form.start_time"
                  type="time"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">End Time</label>
                <input
                  v-model="form.end_time"
                  type="time"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Break Start</label>
                <input
                  v-model="form.break_start"
                  type="time"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Break End</label>
                <input
                  v-model="form.break_end"
                  type="time"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <!-- Working Hours Preview -->
            <div v-if="form.start_time && form.end_time" class="bg-blue-50 border border-blue-200 rounded-md p-4">
              <h5 class="text-sm font-medium text-blue-900 mb-2">Schedule Preview</h5>
              <div class="text-sm text-blue-800">
                <div>📅 Working: {{ form.start_time }} - {{ form.end_time }}</div>
                <div v-if="form.break_start && form.break_end">
                  ☕ Break: {{ form.break_start }} - {{ form.break_end }}
                </div>
                <div class="mt-2 text-xs">
                  Total hours: {{ calculateWorkingHours() }}
                </div>
              </div>
            </div>
          </div>

          <!-- Notes -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Notes</label>
            <textarea
              v-model="form.notes"
              rows="3"
              placeholder="Add any additional notes about this availability..."
              class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
            ></textarea>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="!isFormValid"
              class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ availability?.id ? 'Update Availability' : 'Create Availability' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useEntityStore } from '../stores/entity'
import { createUTCDateTime, convertUTCToEntityDate, formatEntityTime } from '../utils/timezoneUtils'

export default {
  name: 'AvailabilityModal',
  props: {
    availability: {
      type: Object,
      default: null
    },
    doctors: {
      type: Array,
      required: true
    }
  },
  emits: ['close', 'save'],
  setup(props, { emit }) {
    const authStore = useAuthStore()
    const entityStore = useEntityStore()
    
    const form = reactive({
      doctor_id: '',
      date: '',
      status: 'available',
      start_time: '09:00',
      end_time: '17:00',
      break_start: '12:00',
      break_end: '13:00',
      notes: ''
    })

    const statusOptions = [
      {
        value: 'available',
        label: 'Available',
        emoji: '✅',
        description: 'Ready for appointments'
      },
      {
        value: 'unavailable',
        label: 'Unavailable',
        emoji: '❌',
        description: 'Not working today'
      },
      {
        value: 'vacation',
        label: 'Vacation',
        emoji: '🏖️',
        description: 'On vacation'
      },
      {
        value: 'training',
        label: 'Training',
        emoji: '📚',
        description: 'In training/education'
      },
      {
        value: 'sick_leave',
        label: 'Sick Leave',
        emoji: '🤒',
        description: 'Out sick'
      },
      {
        value: 'meeting',
        label: 'Meeting',
        emoji: '👥',
        description: 'In meetings/admin'
      }
    ]

    const isFormValid = computed(() => {
      return form.doctor_id && form.date && form.status
    })

    const calculateWorkingHours = () => {
      if (!form.start_time || !form.end_time) return '0h'
      
      const start = new Date(`1970-01-01T${form.start_time}:00`)
      const end = new Date(`1970-01-01T${form.end_time}:00`)
      let totalMinutes = (end - start) / (1000 * 60)
      
      if (form.break_start && form.break_end) {
        const breakStart = new Date(`1970-01-01T${form.break_start}:00`)
        const breakEnd = new Date(`1970-01-01T${form.break_end}:00`)
        const breakMinutes = (breakEnd - breakStart) / (1000 * 60)
        totalMinutes -= breakMinutes
      }
      
      const hours = Math.floor(totalMinutes / 60)
      const minutes = totalMinutes % 60
      
      return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
    }

    const handleSubmit = () => {
      if (!isFormValid.value) return
      
      // Get the healthcare entity timezone
      const entityTimezone = entityStore.entityTimezone
      if (!entityTimezone) {
        console.error('Healthcare entity timezone not available')
        return
      }
      
      // Convert form data to UTC datetime format expected by the backend
      const submitData = {
        doctor_id: form.doctor_id,
        status: form.status,
        notes: form.notes
      }
      
      // Convert time fields from entity timezone to UTC datetime strings
      if (form.status === 'available' && form.start_time && form.end_time) {
        const date = form.date // YYYY-MM-DD format
        
        // Convert entity timezone times to UTC
        submitData.start_datetime = createUTCDateTime(date, form.start_time, entityTimezone)
        submitData.end_datetime = createUTCDateTime(date, form.end_time, entityTimezone)
        
        // Handle break times if provided
        if (form.break_start && form.break_end) {
          submitData.break_start_datetime = createUTCDateTime(date, form.break_start, entityTimezone)
          submitData.break_end_datetime = createUTCDateTime(date, form.break_end, entityTimezone)
        }
      }
      
      emit('save', submitData)
    }

    // Initialize form with existing data
    onMounted(() => {
      if (props.availability) {
        // Set basic fields
        form.doctor_id = props.availability.doctor_id || ''
        form.status = props.availability.status || 'available'
        form.notes = props.availability.notes || ''
        
        // Set date from the existing date field or derive from start_datetime
        if (props.availability.date) {
          form.date = props.availability.date
        } else if (props.availability.start_datetime) {
          // Extract date from UTC datetime
          form.date = props.availability.start_datetime.split('T')[0]
        }
        
        // Convert UTC datetime fields back to entity timezone for UI editing
        const entityTimezone = entityStore.entityTimezone
        
        if (entityTimezone && props.availability.start_datetime) {
          // Convert UTC to entity timezone and extract time
          const entityTime = formatEntityTime(props.availability.start_datetime, entityTimezone, false) // 24-hour format
          form.start_time = entityTime.replace(' ', '').toLowerCase() // Remove AM/PM for 24h format
          
          // For proper 24h format, let's extract from converted date
          const entityDate = new Date(props.availability.start_datetime.replace('Z', ''))
          const localString = entityDate.toLocaleString('en-CA', { 
            timeZone: entityTimezone,
            hour12: false,
            hour: '2-digit',
            minute: '2-digit'
          })
          form.start_time = localString.split(' ')[1] || localString
        } else if (props.availability.start_time) {
          form.start_time = props.availability.start_time
        }
        
        if (entityTimezone && props.availability.end_datetime) {
          const entityDate = new Date(props.availability.end_datetime.replace('Z', ''))
          const localString = entityDate.toLocaleString('en-CA', { 
            timeZone: entityTimezone,
            hour12: false,
            hour: '2-digit',
            minute: '2-digit'
          })
          form.end_time = localString.split(' ')[1] || localString
        } else if (props.availability.end_time) {
          form.end_time = props.availability.end_time
        }
        
        if (entityTimezone && props.availability.break_start_datetime) {
          const entityDate = new Date(props.availability.break_start_datetime.replace('Z', ''))
          const localString = entityDate.toLocaleString('en-CA', { 
            timeZone: entityTimezone,
            hour12: false,
            hour: '2-digit',
            minute: '2-digit'
          })
          form.break_start = localString.split(' ')[1] || localString
        } else if (props.availability.break_start) {
          form.break_start = props.availability.break_start
        }
        
        if (entityTimezone && props.availability.break_end_datetime) {
          const entityDate = new Date(props.availability.break_end_datetime.replace('Z', ''))
          const localString = entityDate.toLocaleString('en-CA', { 
            timeZone: entityTimezone,
            hour12: false,
            hour: '2-digit',
            minute: '2-digit'
          })
          form.break_end = localString.split(' ')[1] || localString
        } else if (props.availability.break_end) {
          form.break_end = props.availability.break_end
        }
      }
    })

    // Clear working hours when status changes to non-available
    watch(() => form.status, (newStatus) => {
      if (newStatus !== 'available') {
        form.start_time = ''
        form.end_time = ''
        form.break_start = ''
        form.break_end = ''
      } else if (!form.start_time) {
        // Set defaults when switching to available
        form.start_time = '09:00'
        form.end_time = '17:00'
        form.break_start = '12:00'
        form.break_end = '13:00'
      }
    })

    return {
      form,
      statusOptions,
      isFormValid,
      calculateWorkingHours,
      handleSubmit
    }
  }
}
</script>

<style scoped>
/* Custom radio button styling */
input[type="radio"]:checked + div {
  @apply ring-2 ring-blue-600 border-blue-600;
}

/* Smooth transitions */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out, border-color 0.2s ease-in-out;
}
</style>