<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
    <!-- Calendar Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-semibold text-gray-900">
          {{ doctorName }} - Availability Calendar
        </h2>
        <p class="text-sm text-gray-600">{{ currentMonthYear }}</p>
      </div>
      
      <div class="flex items-center space-x-2">
        <!-- Month Navigation -->
        <button
          @click="previousMonth"
          class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
          title="Previous Month"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
        </button>
        
        <div class="text-sm font-medium text-gray-900 min-w-[120px] text-center">
          {{ currentMonthYear }}
        </div>
        
        <button
          @click="nextMonth"
          class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full transition-colors"
          title="Next Month"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Calendar Legend -->
    <div class="flex flex-wrap items-center gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
      <div class="text-sm font-medium text-gray-700">Legend:</div>
      <div
        v-for="status in statusLegend"
        :key="status.value"
        class="flex items-center space-x-2"
      >
        <div :class="status.color" class="w-3 h-3 rounded-full"></div>
        <span class="text-sm text-gray-600">{{ status.emoji }} {{ status.label }}</span>
      </div>
    </div>

    <!-- Calendar Grid -->
    <div class="grid grid-cols-7 gap-1">
      <!-- Day Headers -->
      <div
        v-for="day in dayHeaders"
        :key="day"
        class="p-2 text-center text-sm font-medium text-gray-500 bg-gray-50"
      >
        {{ day }}
      </div>

      <!-- Calendar Days -->
      <div
        v-for="day in calendarDays"
        :key="`${day.date}-${day.isCurrentMonth}`"
        :class="[
          'min-h-[120px] p-2 border border-gray-100 relative cursor-pointer transition-all duration-200',
          day.isCurrentMonth ? 'bg-white hover:bg-gray-50' : 'bg-gray-50 text-gray-400',
          day.isToday ? 'ring-2 ring-blue-500' : '',
          day.isWeekend ? 'bg-gray-25' : ''
        ]"
        @click="handleDayClick(day)"
      >
        <!-- Date Number -->
        <div class="flex justify-between items-start mb-1">
          <span
            :class="[
              'text-sm font-medium',
              day.isToday ? 'text-blue-600' : day.isCurrentMonth ? 'text-gray-900' : 'text-gray-400'
            ]"
          >
            {{ day.dayNumber }}
          </span>
          
          <!-- Add Availability Button -->
          <button
            v-if="day.isCurrentMonth && !day.availability"
            @click.stop="$emit('create-availability', day.date)"
            class="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-blue-600 transition-opacity"
            title="Add Availability"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </button>
        </div>

        <!-- Availability Info -->
        <div v-if="day.availability" class="space-y-1">
          <!-- Status Badge -->
          <div
            :class="[
              'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
              getStatusClass(day.availability.status)
            ]"
          >
            {{ getStatusEmoji(day.availability.status) }}
            {{ getStatusLabel(day.availability.status) }}
          </div>

          <!-- Working Hours -->
          <div v-if="day.availability.status === 'available' && day.availability.start_time" class="text-xs text-gray-600">
            🕐 {{ day.availability.start_time }} - {{ day.availability.end_time }}
          </div>

          <!-- Break Time -->
          <div v-if="day.availability.break_start && day.availability.break_end" class="text-xs text-gray-500">
            ☕ {{ day.availability.break_start }} - {{ day.availability.break_end }}
          </div>

          <!-- Notes Preview -->
          <div v-if="day.availability.notes" class="text-xs text-gray-500 truncate" :title="day.availability.notes">
            📝 {{ day.availability.notes }}
          </div>

          <!-- Action Buttons -->
          <div class="flex space-x-1 mt-2">
            <button
              @click.stop="$emit('edit-availability', day.availability)"
              class="text-xs text-blue-600 hover:text-blue-800 font-medium"
            >
              Edit
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="day.isCurrentMonth" class="flex items-center justify-center h-16 text-gray-400 group">
          <div class="text-center opacity-0 group-hover:opacity-100 transition-opacity">
            <div class="text-2xl mb-1">📅</div>
            <div class="text-xs">Click to add</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Calendar Footer with Summary -->
    <div class="mt-6 pt-4 border-t border-gray-200">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center space-x-6 text-sm text-gray-600">
          <div>📊 This Month:</div>
          <div v-for="(count, status) in monthSummary" :key="status" class="flex items-center space-x-1">
            <span>{{ getStatusEmoji(status) }}</span>
            <span class="font-medium">{{ count }}</span>
          </div>
        </div>
        
        <div class="text-sm text-gray-500">
          Total working days: {{ workingDaysCount }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue'

export default {
  name: 'AvailabilityCalendar',
  props: {
    doctorId: {
      type: [Number, String],
      required: true
    },
    availabilityData: {
      type: Array,
      default: () => []
    }
  },
  emits: ['edit-availability', 'create-availability'],
  setup(props, { emit }) {
    const currentDate = ref(new Date())
    
    const dayHeaders = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
    
    const statusLegend = [
      { value: 'available', label: 'Available', emoji: '✅', color: 'bg-green-500' },
      { value: 'unavailable', label: 'Unavailable', emoji: '❌', color: 'bg-gray-500' },
      { value: 'vacation', label: 'Vacation', emoji: '🏖️', color: 'bg-blue-500' },
      { value: 'training', label: 'Training', emoji: '📚', color: 'bg-yellow-500' },
      { value: 'sick_leave', label: 'Sick Leave', emoji: '🤒', color: 'bg-red-500' },
      { value: 'meeting', label: 'Meeting', emoji: '👥', color: 'bg-purple-500' }
    ]

    const currentMonthYear = computed(() => {
      return currentDate.value.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long'
      })
    })

    const doctorName = computed(() => {
      // Would normally get this from the doctor data
      return `Doctor #${props.doctorId}`
    })

    const calendarDays = computed(() => {
      const year = currentDate.value.getFullYear()
      const month = currentDate.value.getMonth()
      
      // Get first day of month and how many days to show from previous month
      const firstDay = new Date(year, month, 1)
      const lastDay = new Date(year, month + 1, 0)
      const startDate = new Date(firstDay)
      startDate.setDate(startDate.getDate() - firstDay.getDay())
      
      // Create array of 42 days (6 weeks)
      const days = []
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      
      for (let i = 0; i < 42; i++) {
        const date = new Date(startDate)
        date.setDate(startDate.getDate() + i)
        
        const dateString = date.toISOString().split('T')[0]
        const availability = props.availabilityData.find(a => a.date === dateString)
        
        days.push({
          date: dateString,
          dayNumber: date.getDate(),
          isCurrentMonth: date.getMonth() === month,
          isToday: date.getTime() === today.getTime(),
          isWeekend: date.getDay() === 0 || date.getDay() === 6,
          availability
        })
      }
      
      return days
    })

    const monthSummary = computed(() => {
      const summary = {}
      const currentMonth = currentDate.value.getMonth()
      const currentYear = currentDate.value.getFullYear()
      
      calendarDays.value.forEach(day => {
        const date = new Date(day.date)
        if (date.getMonth() === currentMonth && date.getFullYear() === currentYear && day.availability) {
          const status = day.availability.status
          summary[status] = (summary[status] || 0) + 1
        }
      })
      
      return summary
    })

    const workingDaysCount = computed(() => {
      return calendarDays.value.filter(day => 
        day.isCurrentMonth && 
        day.availability && 
        day.availability.status === 'available'
      ).length
    })

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
        vacation: 'Vacation',
        training: 'Training',
        sick_leave: 'Sick Leave',
        meeting: 'Meeting'
      }
      return labels[status] || status
    }

    // Navigation
    const previousMonth = () => {
      currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
    }

    const nextMonth = () => {
      currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
    }

    const handleDayClick = (day) => {
      if (!day.isCurrentMonth) return
      
      if (day.availability) {
        emit('edit-availability', day.availability)
      } else {
        emit('create-availability', day.date)
      }
    }

    return {
      currentDate,
      dayHeaders,
      statusLegend,
      currentMonthYear,
      doctorName,
      calendarDays,
      monthSummary,
      workingDaysCount,
      getStatusClass,
      getStatusEmoji,
      getStatusLabel,
      previousMonth,
      nextMonth,
      handleDayClick
    }
  }
}
</script>

<style scoped>
/* Calendar grid styling */
.grid-cols-7 > div {
  @apply min-h-[120px];
}

/* Hover effects */
.group:hover .opacity-0 {
  @apply opacity-100;
}

/* Weekend styling */
.bg-gray-25 {
  background-color: #fafafa;
}

/* Today indicator */
.ring-2.ring-blue-500 {
  box-shadow: 0 0 0 2px #3b82f6;
}

/* Smooth transitions */
.transition-all {
  transition: all 0.2s ease-in-out;
}

.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

.transition-opacity {
  transition: opacity 0.2s ease-in-out;
}
</style>