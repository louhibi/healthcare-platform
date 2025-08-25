<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="$emit('close')">
    <div class="relative top-10 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white" @click.stop>
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Bulk Schedule Creation</h3>
            <p class="text-sm text-gray-600 mt-1">Create availability for multiple dates at once</p>
          </div>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Doctor and Date Range -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
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
                </option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">From Date *</label>
              <input
                v-model="form.date_from"
                type="date"
                required
                :min="today"
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">To Date *</label>
              <input
                v-model="form.date_to"
                type="date"
                required
                :min="form.date_from || today"
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>

          <!-- Date Range Preview -->
          <div v-if="form.date_from && form.date_to" class="bg-blue-50 border border-blue-200 rounded-md p-4">
            <h5 class="text-sm font-medium text-blue-900 mb-2">📅 Date Range Preview</h5>
            <div class="text-sm text-blue-800">
              <div><strong>{{ dayCount }}</strong> days from {{ formatDate(form.date_from) }} to {{ formatDate(form.date_to) }}</div>
              <div class="mt-1">{{ weekdayCount }} weekdays, {{ weekendCount }} weekend days</div>
            </div>
          </div>

          <!-- Day Selection -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-3">Apply to Days *</label>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
              <label
                v-for="day in daysOfWeek"
                :key="day.value"
                :class="[
                  'relative flex cursor-pointer rounded-lg border p-3 focus:outline-none',
                  selectedDays.includes(day.value)
                    ? 'border-blue-600 ring-2 ring-blue-600 bg-blue-50'
                    : 'border-gray-300 bg-white hover:bg-gray-50'
                ]"
              >
                <input
                  v-model="selectedDays"
                  type="checkbox"
                  :value="day.value"
                  class="sr-only"
                />
                <div class="flex items-center justify-between w-full">
                  <div class="flex items-center">
                    <div class="text-lg mr-2">{{ day.emoji }}</div>
                    <div class="text-sm font-medium text-gray-900">{{ day.label }}</div>
                  </div>
                  <div v-if="selectedDays.includes(day.value)" class="text-blue-600">
                    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
              </label>
            </div>
            <div class="mt-2 text-xs text-gray-500">
              Select which days of the week to create availability for
            </div>
          </div>

          <!-- Quick Selection Buttons -->
          <div class="flex flex-wrap gap-2">
            <button
              type="button"
              @click="selectWeekdays"
              class="px-3 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 transition-colors"
            >
              Weekdays Only
            </button>
            <button
              type="button"
              @click="selectWeekends"
              class="px-3 py-1.5 text-xs font-medium text-purple-600 bg-purple-50 border border-purple-200 rounded-md hover:bg-purple-100 transition-colors"
            >
              Weekends Only
            </button>
            <button
              type="button"
              @click="selectAllDays"
              class="px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 transition-colors"
            >
              All Days
            </button>
            <button
              type="button"
              @click="clearDays"
              class="px-3 py-1.5 text-xs font-medium text-red-600 bg-red-50 border border-red-200 rounded-md hover:bg-red-100 transition-colors"
            >
              Clear All
            </button>
          </div>

          <!-- Template Configuration -->
          <div class="border-t pt-6">
            <h4 class="text-lg font-medium text-gray-900 mb-4">Availability Template</h4>
            
            <!-- Status -->
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">Status *</label>
              <select
                v-model="form.template.status"
                required
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="available">✅ Available</option>
                <option value="unavailable">❌ Unavailable</option>
                <option value="vacation">🏖️ Vacation</option>
                <option value="training">📚 Training</option>
                <option value="sick_leave">🤒 Sick Leave</option>
                <option value="meeting">👥 Meeting</option>
              </select>
            </div>

            <!-- Working Hours (only for available status) -->
            <div v-if="form.template.status === 'available'" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Start Time</label>
                  <input
                    v-model="form.template.start_time"
                    type="time"
                    class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">End Time</label>
                  <input
                    v-model="form.template.end_time"
                    type="time"
                    class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Break Start</label>
                  <input
                    v-model="form.template.break_start"
                    type="time"
                    class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Break End</label>
                  <input
                    v-model="form.template.break_end"
                    type="time"
                    class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
            </div>

            <!-- Notes -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Notes</label>
              <textarea
                v-model="form.template.notes"
                rows="3"
                placeholder="Add notes that will apply to all created availability records..."
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
              ></textarea>
            </div>
          </div>

          <!-- Summary -->
          <div v-if="isFormValid" class="bg-green-50 border border-green-200 rounded-md p-4">
            <h5 class="text-sm font-medium text-green-900 mb-2">📋 Creation Summary</h5>
            <div class="text-sm text-green-800 space-y-1">
              <div>✅ Will create <strong>{{ estimatedRecords }}</strong> availability records</div>
              <div>👩‍⚕️ For: {{ selectedDoctorName }}</div>
              <div>📅 Status: {{ form.template.status }}</div>
              <div v-if="form.template.status === 'available' && form.template.start_time">
                🕐 Hours: {{ form.template.start_time }} - {{ form.template.end_time }}
              </div>
              <div>📆 Days: {{ selectedDayLabels.join(', ') }}</div>
            </div>
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
              class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Create {{ estimatedRecords }} Records
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { createUTCDateTime } from '../utils/timezoneUtils'
import { useAuthStore } from '@/stores/auth'
import { useEntityStore } from '@/stores/entity'

export default {
  name: 'BulkScheduleModal',
  props: {
    doctors: {
      type: Array,
      required: true
    }
  },
  emits: ['close', 'save'],
  setup(props, { emit }) {
    const authStore = useAuthStore()
    const entityStore = useEntityStore()
    const today = new Date().toISOString().split('T')[0]
    
    const form = reactive({
      doctor_id: '',
      date_from: '',
      date_to: '',
      template: {
        status: 'available',
        start_time: '09:00',
        end_time: '17:00',
        break_start: '12:00',
        break_end: '13:00',
        notes: ''
      }
    })

    const selectedDays = ref([1, 2, 3, 4, 5]) // Default to weekdays (Mon-Fri)

    const daysOfWeek = [
      { value: 0, label: 'Sunday', emoji: '🟦' },
      { value: 1, label: 'Monday', emoji: '🟩' },
      { value: 2, label: 'Tuesday', emoji: '🟩' },
      { value: 3, label: 'Wednesday', emoji: '🟩' },
      { value: 4, label: 'Thursday', emoji: '🟩' },
      { value: 5, label: 'Friday', emoji: '🟩' },
      { value: 6, label: 'Saturday', emoji: '🟦' }
    ]

    // Computed properties
    const isFormValid = computed(() => {
      return form.doctor_id && 
             form.date_from && 
             form.date_to && 
             selectedDays.value.length > 0 &&
             new Date(form.date_from) <= new Date(form.date_to)
    })

    const dayCount = computed(() => {
      if (!form.date_from || !form.date_to) return 0
      const start = new Date(form.date_from)
      const end = new Date(form.date_to)
      return Math.ceil((end - start) / (1000 * 60 * 60 * 24)) + 1
    })

    const weekdayCount = computed(() => {
      if (!form.date_from || !form.date_to) return 0
      let count = 0
      let current = new Date(form.date_from)
      const end = new Date(form.date_to)
      
      while (current <= end) {
        const dayOfWeek = current.getDay()
        if (selectedDays.value.includes(dayOfWeek) && dayOfWeek >= 1 && dayOfWeek <= 5) {
          count++
        }
        current.setDate(current.getDate() + 1)
      }
      
      return count
    })

    const weekendCount = computed(() => {
      if (!form.date_from || !form.date_to) return 0
      let count = 0
      let current = new Date(form.date_from)
      const end = new Date(form.date_to)
      
      while (current <= end) {
        const dayOfWeek = current.getDay()
        if (selectedDays.value.includes(dayOfWeek) && (dayOfWeek === 0 || dayOfWeek === 6)) {
          count++
        }
        current.setDate(current.getDate() + 1)
      }
      
      return count
    })

    const estimatedRecords = computed(() => {
      if (!form.date_from || !form.date_to || selectedDays.value.length === 0) return 0
      
      let count = 0
      let current = new Date(form.date_from)
      const end = new Date(form.date_to)
      
      while (current <= end) {
        if (selectedDays.value.includes(current.getDay())) {
          count++
        }
        current.setDate(current.getDate() + 1)
      }
      
      return count
    })

    const selectedDoctorName = computed(() => {
      const doctor = props.doctors.find(d => d.id === form.doctor_id)
      return doctor ? `Dr. ${doctor.first_name} ${doctor.last_name}` : ''
    })

    const selectedDayLabels = computed(() => {
      return selectedDays.value
        .map(dayValue => daysOfWeek.find(day => day.value === dayValue)?.label)
        .filter(Boolean)
    })

    // Utility functions
    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    // Day selection helpers
    const selectWeekdays = () => {
      selectedDays.value = [1, 2, 3, 4, 5]
    }

    const selectWeekends = () => {
      selectedDays.value = [0, 6]
    }

    const selectAllDays = () => {
      selectedDays.value = [0, 1, 2, 3, 4, 5, 6]
    }

    const clearDays = () => {
      selectedDays.value = []
    }

    const handleSubmit = () => {
      if (!isFormValid.value) return

      const entityTimezone = entityStore.entityTimezone
      if (!entityTimezone) {
        console.error('Healthcare entity timezone not available')
        return
      }

      // Convert template times to UTC datetime format
      const template = {
        status: form.template.status,
        notes: form.template.notes
      }
      
      // Convert time fields to UTC datetime strings if available
      if (form.template.status === 'available' && form.template.start_time && form.template.end_time) {
        // Use a reference date for the template (will be adjusted for each actual date)
        const referenceDate = form.date_from // Use start date as reference
        
        try {
          template.start_datetime = createUTCDateTime(referenceDate, form.template.start_time, entityTimezone)
          template.end_datetime = createUTCDateTime(referenceDate, form.template.end_time, entityTimezone)
          
          // Handle break times if provided
          if (form.template.break_start && form.template.break_end) {
            template.break_start_datetime = createUTCDateTime(referenceDate, form.template.break_start, entityTimezone)
            template.break_end_datetime = createUTCDateTime(referenceDate, form.template.break_end, entityTimezone)
          }
        } catch (error) {
          console.error('Failed to convert template times to UTC:', error)
          return
        }
      }

      const submitData = {
        doctor_id: form.doctor_id,
        date_from: form.date_from,
        date_to: form.date_to,
        selected_days: selectedDays.value,
        template: template
      }

      emit('save', submitData)
    }

    return {
      today,
      form,
      selectedDays,
      daysOfWeek,
      isFormValid,
      dayCount,
      weekdayCount,
      weekendCount,
      estimatedRecords,
      selectedDoctorName,
      selectedDayLabels,
      formatDate,
      selectWeekdays,
      selectWeekends,
      selectAllDays,
      clearDays,
      handleSubmit
    }
  }
}
</script>

<style scoped>
/* Checkbox styling */
input[type="checkbox"]:checked + div {
  @apply bg-blue-50 border-blue-600 ring-2 ring-blue-600;
}

/* Smooth transitions */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out, border-color 0.2s ease-in-out;
}
</style>