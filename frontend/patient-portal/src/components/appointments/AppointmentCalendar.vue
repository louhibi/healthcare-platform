<template>
  <div class="bg-white shadow rounded-lg">
    <!-- Calendar Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <h3 class="text-lg font-medium text-gray-900">
            {{ currentMonthYear }}
          </h3>
          <div class="flex items-center space-x-1">
            <button
              @click="previousMonth"
              class="p-1.5 text-gray-400 hover:text-gray-600 rounded-md hover:bg-gray-100"
            >
              <ChevronLeftIcon class="h-5 w-5" />
            </button>
            <button
              @click="nextMonth"
              class="p-1.5 text-gray-400 hover:text-gray-600 rounded-md hover:bg-gray-100"
            >
              <ChevronRightIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
        <button
          @click="goToToday"
          class="px-3 py-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700 border border-indigo-300 rounded-md hover:bg-indigo-50"
        >
          Today
        </button>
      </div>
    </div>

    <!-- Calendar Grid -->
    <div class="p-6">
      <!-- Weekday Headers -->
      <div class="grid grid-cols-7 gap-px mb-2">
        <div
          v-for="day in weekdays"
          :key="day"
          class="py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wide"
        >
          {{ day }}
        </div>
      </div>

      <!-- Calendar Days -->
      <div class="grid grid-cols-7 gap-px bg-gray-200 rounded-lg overflow-hidden">
        <div
          v-for="day in calendarDays"
          :key="`${day.date}-${day.isCurrentMonth}`"
          :class="[
            'min-h-[120px] bg-white p-2 hover:bg-gray-50',
            !day.isCurrentMonth && 'bg-gray-50 text-gray-400',
            day.isToday && 'bg-indigo-50'
          ]"
        >
          <!-- Day Number -->
          <div class="flex items-center justify-between mb-2">
            <span
              :class="[
                'text-sm font-medium',
                day.isToday
                  ? 'bg-indigo-600 text-white rounded-full w-6 h-6 flex items-center justify-center'
                  : day.isCurrentMonth
                    ? 'text-gray-900'
                    : 'text-gray-400'
              ]"
            >
              {{ day.dayNumber }}
            </span>
            <span
              v-if="day.appointments.length > 0"
              class="text-xs text-gray-500"
            >
              {{ day.appointments.length }}
            </span>
          </div>

          <!-- Appointments -->
          <div class="space-y-1">
            <div
              v-for="(appointment, index) in day.appointments.slice(0, 3)"
              :key="appointment.id"
              @click="$emit('appointment-click', appointment)"
              :class="[
                'text-xs px-2 py-1 rounded cursor-pointer truncate',
                getAppointmentColor(appointment.status)
              ]"
              :title="`${appointment.patient_name || 'Patient'} - ${formatTime(appointment.date_time)}`"
            >
              <div class="font-medium truncate">
                {{ formatTime(appointment.date_time) }}
              </div>
              <div class="truncate opacity-90">
                {{ appointment.patient_name || 'Patient' }}
              </div>
            </div>
            
            <!-- Show more indicator -->
            <div
              v-if="day.appointments.length > 3"
              class="text-xs text-gray-500 px-2 py-1 cursor-pointer hover:text-gray-700"
              @click="$emit('show-more', day)"
            >
              +{{ day.appointments.length - 3 }} more
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="px-6 py-4 border-t border-gray-200">
      <div class="flex items-center space-x-6 text-xs">
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-yellow-200 rounded"></div>
          <span class="text-gray-600">Scheduled</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-blue-200 rounded"></div>
          <span class="text-gray-600">Confirmed</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-green-200 rounded"></div>
          <span class="text-gray-600">In Progress</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-gray-200 rounded"></div>
          <span class="text-gray-600">Completed</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import { formatEntityTime, getCurrentEntityDate } from '@/utils/timezoneUtils'
import { useAuthStore } from '@/stores/auth'

// Props
const props = defineProps({
  appointments: {
    type: Array,
    default: () => []
  },
  entityTimezone: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits(['appointment-click', 'show-more', 'month-changed'])

// Auth store
const authStore = useAuthStore()

// Reactive state
const currentDate = ref(new Date())

// Computed
const weekdays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

const currentMonthYear = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', {
    month: 'long',
    year: 'numeric'
  })
})

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  // Get first day of month and how many days in month
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const daysInMonth = lastDay.getDate()
  const startingDayOfWeek = firstDay.getDay()
  
  // Get previous month's last days
  const prevMonth = new Date(year, month, 0)
  const daysInPrevMonth = prevMonth.getDate()
  
  const days = []
  
  // Previous month's trailing days
  for (let day = daysInPrevMonth - startingDayOfWeek + 1; day <= daysInPrevMonth; day++) {
    const date = new Date(year, month - 1, day)
    days.push({
      date: date.toISOString().split('T')[0],
      dayNumber: day,
      isCurrentMonth: false,
      isToday: false,
      appointments: []
    })
  }
  
  // Current month days
  for (let day = 1; day <= daysInMonth; day++) {
    const date = new Date(year, month, day)
    const dateString = date.toISOString().split('T')[0]
    const isToday = dateString === new Date().toISOString().split('T')[0]
    
    // Filter appointments for this day
    const dayAppointments = props.appointments.filter(appointment => {
      if (!appointment.date_time) return false
      const appointmentDate = new Date(appointment.date_time).toISOString().split('T')[0]
      return appointmentDate === dateString
    })
    
    days.push({
      date: dateString,
      dayNumber: day,
      isCurrentMonth: true,
      isToday,
      appointments: dayAppointments
    })
  }
  
  // Next month's leading days
  const remainingCells = 42 - days.length // 6 rows * 7 days
  for (let day = 1; day <= remainingCells; day++) {
    const date = new Date(year, month + 1, day)
    days.push({
      date: date.toISOString().split('T')[0],
      dayNumber: day,
      isCurrentMonth: false,
      isToday: false,
      appointments: []
    })
  }
  
  return days
})

// Methods
const previousMonth = () => {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
  emit('month-changed', currentDate.value)
}

const nextMonth = () => {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
  emit('month-changed', currentDate.value)
}

const goToToday = () => {
  currentDate.value = new Date()
  emit('month-changed', currentDate.value)
}

const formatTime = (utcDateTimeString) => {
  if (!utcDateTimeString || !props.entityTimezone) return ''
  
  try {
    return formatEntityTime(utcDateTimeString, props.entityTimezone, true)
  } catch (error) {
    console.warn('Failed to format time:', error)
    return new Date(utcDateTimeString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
}

const getAppointmentColor = (status) => {
  const colors = {
    scheduled: 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200',
    confirmed: 'bg-blue-100 text-blue-800 hover:bg-blue-200',
    'in-progress': 'bg-green-100 text-green-800 hover:bg-green-200',
    completed: 'bg-gray-100 text-gray-800 hover:bg-gray-200',
    cancelled: 'bg-red-100 text-red-800 hover:bg-red-200',
    'no-show': 'bg-orange-100 text-orange-800 hover:bg-orange-200'
  }
  return colors[status] || colors.scheduled
}

// Watch for month changes to potentially load more data
watch(currentDate, (newDate) => {
  emit('month-changed', newDate)
})
</script>