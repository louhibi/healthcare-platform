<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="$emit('close')">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-2/3 lg:w-1/2 shadow-lg rounded-md bg-white" @click.stop>
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Availability Details</h3>
            <p class="text-sm text-gray-600 mt-1">{{ availability.doctor_name || 'Doctor' }} - {{ formatDate(availability.date) }}</p>
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

        <!-- Status Badge -->
        <div class="mb-6">
          <div
            :class="[
              'inline-flex items-center px-4 py-2 rounded-full text-sm font-medium',
              getStatusClass(availability.status)
            ]"
          >
            <span class="text-lg mr-2">{{ getStatusEmoji(availability.status) }}</span>
            {{ getStatusLabel(availability.status) }}
          </div>
        </div>

        <!-- Content Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Basic Information -->
          <div>
            <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Basic Information</h4>
            <div class="space-y-3">
              <div>
                <label class="text-sm font-medium text-gray-600">Doctor</label>
                <p class="text-gray-900 mt-1">{{ availability.doctor_name || 'Doctor' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Date</label>
                <p class="text-gray-900 mt-1">
                  {{ formatDate(availability.date) }}
                  <span class="text-gray-500 text-sm ml-2">({{ getDayOfWeek(availability.date) }})</span>
                </p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Status</label>
                <p class="text-gray-900 mt-1">
                  <span :class="getStatusClass(availability.status)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                    {{ getStatusEmoji(availability.status) }} {{ getStatusLabel(availability.status) }}
                  </span>
                </p>
              </div>
            </div>
          </div>

          <!-- Working Hours (only show if available) -->
          <div v-if="availability.status === 'available' && (availability.start_time || availability.end_time)">
            <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Working Hours</h4>
            <div class="space-y-3">
              <div v-if="availability.start_time && availability.end_time">
                <label class="text-sm font-medium text-gray-600">Work Schedule</label>
                <div class="text-gray-900 mt-1 flex items-center">
                  <span class="text-lg mr-2">🕐</span>
                  {{ availability.start_time }} - {{ availability.end_time }}
                  <span class="text-gray-500 text-sm ml-2">({{ calculateWorkingHours() }})</span>
                </div>
              </div>
              
              <div v-if="availability.break_start && availability.break_end">
                <label class="text-sm font-medium text-gray-600">Break Time</label>
                <div class="text-gray-900 mt-1 flex items-center">
                  <span class="text-lg mr-2">☕</span>
                  {{ availability.break_start }} - {{ availability.break_end }}
                  <span class="text-gray-500 text-sm ml-2">({{ calculateBreakHours() }})</span>
                </div>
              </div>

              <!-- Visual Schedule -->
              <div class="bg-gray-50 border border-gray-200 rounded-md p-3 mt-4">
                <h5 class="text-sm font-medium text-gray-700 mb-2">📅 Schedule Timeline</h5>
                <div class="space-y-1">
                  <div class="flex items-center text-sm text-gray-600">
                    <div class="w-16 text-right mr-3">{{ availability.start_time }}</div>
                    <div class="flex-1 bg-green-200 h-2 rounded-l"></div>
                    <div class="ml-2 text-xs text-green-700">Start Work</div>
                  </div>
                  
                  <div v-if="availability.break_start" class="flex items-center text-sm text-gray-600">
                    <div class="w-16 text-right mr-3">{{ availability.break_start }}</div>
                    <div class="flex-1 bg-yellow-200 h-2"></div>
                    <div class="ml-2 text-xs text-yellow-700">Break</div>
                  </div>
                  
                  <div v-if="availability.break_end" class="flex items-center text-sm text-gray-600">
                    <div class="w-16 text-right mr-3">{{ availability.break_end }}</div>
                    <div class="flex-1 bg-green-200 h-2"></div>
                    <div class="ml-2 text-xs text-green-700">Resume Work</div>
                  </div>
                  
                  <div class="flex items-center text-sm text-gray-600">
                    <div class="w-16 text-right mr-3">{{ availability.end_time }}</div>
                    <div class="flex-1 bg-red-200 h-2 rounded-r"></div>
                    <div class="ml-2 text-xs text-red-700">End Work</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Non-Available Status Info -->
          <div v-else-if="availability.status !== 'available'">
            <h4 class="text-lg font-medium text-gray-900 mb-4 border-b pb-2">Status Information</h4>
            <div class="space-y-3">
              <div class="flex items-start space-x-3">
                <div class="text-2xl">{{ getStatusEmoji(availability.status) }}</div>
                <div>
                  <h5 class="font-medium text-gray-900">{{ getStatusLabel(availability.status) }}</h5>
                  <p class="text-sm text-gray-600 mt-1">{{ getStatusDescription(availability.status) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Notes -->
        <div v-if="availability.notes" class="mt-6">
          <h4 class="text-lg font-medium text-gray-900 mb-3 border-b pb-2">Notes</h4>
          <div class="bg-yellow-50 border border-yellow-200 rounded-md p-4">
            <div class="flex items-start">
              <div class="text-yellow-500 mr-3">📝</div>
              <div class="text-gray-900 text-sm whitespace-pre-wrap">{{ availability.notes }}</div>
            </div>
          </div>
        </div>

        <!-- Metadata -->
        <div class="mt-6 pt-6 border-t border-gray-200">
          <h4 class="text-lg font-medium text-gray-900 mb-3">Record Information</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600">
            <div>
              <label class="font-medium text-gray-700">Created</label>
              <p>{{ formatDateTime(availability.created_at) }}</p>
            </div>
            <div>
              <label class="font-medium text-gray-700">Last Updated</label>
              <p>{{ formatDateTime(availability.updated_at) }}</p>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200 mt-6">
          <button
            @click="$emit('edit', availability)"
            class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors"
          >
            ✏️ Edit Availability
          </button>
          <button
            @click="$emit('close')"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ViewAvailabilityModal',
  props: {
    availability: {
      type: Object,
      required: true
    }
  },
  emits: ['close', 'edit'],
  setup(props) {
    // Utility functions
    const getStatusClass = (status) => {
      const classes = {
        available: 'bg-green-100 text-green-800',
        unavailable: 'bg-gray-100 text-gray-800',
        vacation: 'bg-blue-100 text-blue-800',
        training: 'bg-yellow-100 text-yellow-800',
        sick_leave: 'bg-red-100 text-red-800',
        meeting: 'bg-purple-100 text-purple-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    }

    const getStatusEmoji = (status) => {
      const emojis = {
        available: '✅',
        unavailable: '❌',
        vacation: '🏖️',
        training: '📚',
        sick_leave: '🤒',
        meeting: '👥'
      }
      return emojis[status] || '❓'
    }

    const getStatusLabel = (status) => {
      const labels = {
        available: 'Available',
        unavailable: 'Unavailable',
        vacation: 'On Vacation',
        training: 'In Training',
        sick_leave: 'Sick Leave',
        meeting: 'In Meeting'
      }
      return labels[status] || status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
    }

    const getStatusDescription = (status) => {
      const descriptions = {
        available: 'Doctor is available for appointments during scheduled hours',
        unavailable: 'Doctor is not available for appointments today',
        vacation: 'Doctor is on vacation and not available for appointments',
        training: 'Doctor is attending training or educational activities',
        sick_leave: 'Doctor is on sick leave and not available',
        meeting: 'Doctor is in meetings or administrative duties'
      }
      return descriptions[status] || 'No additional information available'
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const formatDateTime = (dateTimeString) => {
      return new Date(dateTimeString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const getDayOfWeek = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        weekday: 'long'
      })
    }

    const calculateWorkingHours = () => {
      const { start_time, end_time, break_start, break_end } = props.availability
      
      if (!start_time || !end_time) return '0h'
      
      const start = new Date(`1970-01-01T${start_time}:00`)
      const end = new Date(`1970-01-01T${end_time}:00`)
      let totalMinutes = (end - start) / (1000 * 60)
      
      if (break_start && break_end) {
        const breakStart = new Date(`1970-01-01T${break_start}:00`)
        const breakEnd = new Date(`1970-01-01T${break_end}:00`)
        const breakMinutes = (breakEnd - breakStart) / (1000 * 60)
        totalMinutes -= breakMinutes
      }
      
      const hours = Math.floor(totalMinutes / 60)
      const minutes = totalMinutes % 60
      
      return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
    }

    const calculateBreakHours = () => {
      const { break_start, break_end } = props.availability
      
      if (!break_start || !break_end) return '0h'
      
      const breakStart = new Date(`1970-01-01T${break_start}:00`)
      const breakEnd = new Date(`1970-01-01T${break_end}:00`)
      const breakMinutes = (breakEnd - breakStart) / (1000 * 60)
      
      const hours = Math.floor(breakMinutes / 60)
      const minutes = breakMinutes % 60
      
      return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`
    }

    return {
      getStatusClass,
      getStatusEmoji,
      getStatusLabel,
      getStatusDescription,
      formatDate,
      formatDateTime,
      getDayOfWeek,
      calculateWorkingHours,
      calculateBreakHours
    }
  }
}
</script>

<style scoped>
/* Schedule timeline styling */
.bg-green-200 { background-color: #bbf7d0; }
.bg-yellow-200 { background-color: #fef3c7; }
.bg-red-200 { background-color: #fecaca; }

.text-green-700 { color: #15803d; }
.text-yellow-700 { color: #a16207; }
.text-red-700 { color: #b91c1c; }

/* Smooth transitions */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out, border-color 0.2s ease-in-out;
}
</style>