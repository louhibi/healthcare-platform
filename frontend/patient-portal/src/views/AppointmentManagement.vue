<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="sm:flex sm:items-center">
        <div class="sm:flex-auto">
          <h1 class="text-2xl font-semibold text-gray-900">Appointment Management</h1>
          <p class="mt-2 text-sm text-gray-700">
            Manage patient appointments, view schedules, and book new appointments
          </p>
        </div>
        <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none sm:flex sm:items-center sm:space-x-3">
          <!-- View Toggle -->
          <div class="flex items-center rounded-lg bg-gray-100 p-1">
            <button
              @click="currentView = 'table'"
              :class="[
                'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
                currentView === 'table'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              ]"
            >
              <ViewColumnsIcon class="h-4 w-4 mr-1.5 inline" />
              Table
            </button>
            <button
              @click="currentView = 'calendar'"
              :class="[
                'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
                currentView === 'calendar'
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              ]"
            >
              <CalendarDaysIcon class="h-4 w-4 mr-1.5 inline" />
              Calendar
            </button>
          </div>
          
          <button
            @click="showBookingModal = true"
            type="button"
            class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            <PlusIcon class="-ml-0.5 mr-1.5 h-5 w-5" />
            Book Appointment
          </button>
        </div>
      </div>

      <!-- Stats -->
      <div class="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">Today's Appointments</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
            {{ stats.today_appointments || 0 }}
          </dd>
        </div>
        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">This Week</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
            {{ stats.week_appointments || 0 }}
          </dd>
        </div>
        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">Pending</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
            {{ stats.by_status?.scheduled || 0 }}
          </dd>
        </div>
        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">Total Active</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
            {{ stats.total_appointments || 0 }}
          </dd>
        </div>
      </div>

      <!-- Filters -->
      <div class="mt-8 flex flex-col sm:flex-row gap-4 bg-white p-4 rounded-lg shadow">
        <div class="flex-1">
          <label for="date-filter" class="block text-sm font-medium text-gray-700">Date Range</label>
          <div class="flex gap-2 mt-1">
            <input
              id="date-from"
              v-model="filters.dateFrom"
              type="date"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
            <input
              id="date-to"
              v-model="filters.dateTo"
              type="date"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>
        </div>
        <div class="flex-1">
          <label for="status-filter" class="block text-sm font-medium text-gray-700">Status</label>
          <select
            id="status-filter"
            v-model="filters.status"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Statuses</option>
            <option value="scheduled">Scheduled</option>
            <option value="confirmed">Confirmed</option>
            <option value="in-progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
            <option value="no-show">No Show</option>
          </select>
        </div>
        <div class="flex-1">
          <label for="doctor-filter" class="block text-sm font-medium text-gray-700">Doctor</label>
          <select
            id="doctor-filter"
            v-model="filters.doctorId"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Doctors</option>
            <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
              Dr. {{ doctor.first_name }} {{ doctor.last_name }}
            </option>
          </select>
        </div>
        <div class="flex items-end">
          <button
            @click="loadAppointments"
            type="button"
            class="inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
          >
            <MagnifyingGlassIcon class="-ml-0.5 mr-1.5 h-5 w-5 text-gray-400" />
            Search
          </button>
        </div>
      </div>

      <!-- Table View -->
      <div v-if="currentView === 'table'" class="mt-8 flow-root">
        <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
              <table class="min-w-full divide-y divide-gray-300">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Patient
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Doctor
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Date & Time
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Room
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Type
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Status
                    </th>
                    <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wide text-gray-500">
                      Priority
                    </th>
                    <th class="relative px-6 py-3">
                      <span class="sr-only">Actions</span>
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 bg-white">
                  <tr
                    v-for="appointment in appointments"
                    :key="appointment.id"
                    class="hover:bg-gray-50"
                  >
                    <td class="whitespace-nowrap px-6 py-4 text-sm">
                      <div class="flex items-center">
                        <div class="h-10 w-10 flex-shrink-0">
                          <div class="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
                            <UserIcon class="h-6 w-6 text-gray-400" />
                          </div>
                        </div>
                        <div class="ml-4">
                          <div class="font-medium text-gray-900">{{ appointment.patient_name || 'Patient' }}</div>
                          <div class="text-gray-500">ID: {{ appointment.patient_id }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm">
                      <div class="text-gray-900">{{ appointment.doctor_name || 'Doctor' }}</div>
                      <div class="text-gray-500">ID: {{ appointment.doctor_id }}</div>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                      <div>{{ formatDate(appointment.date_time) }}</div>
                      <div class="text-gray-500">{{ formatTime(appointment.date_time) }}</div>
                      <div class="text-xs text-gray-400">{{ appointment.duration }}min</div>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                      <div v-if="appointment.room_number">
                        <div class="font-medium">Room {{ appointment.room_number }}</div>
                        <div v-if="appointment.room_name" class="text-xs text-gray-500">{{ appointment.room_name }}</div>
                      </div>
                      <div v-else class="text-gray-400 text-sm">No room assigned</div>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                      <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                            :class="getTypeClass(appointment.type)">
                        {{ appointment.type }}
                      </span>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm">
                      <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                            :class="getStatusClass(appointment.status)">
                        {{ appointment.status }}
                      </span>
                    </td>
                    <td class="whitespace-nowrap px-6 py-4 text-sm">
                      <span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                            :class="getPriorityClass(appointment.priority)">
                        {{ appointment.priority || 'normal' }}
                      </span>
                    </td>
                    <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                      <div class="flex items-center justify-end space-x-2">
                        <button
                          @click="viewAppointment(appointment)"
                          class="text-indigo-600 hover:text-indigo-900"
                        >
                          <EyeIcon class="h-5 w-5" />
                        </button>
                        <button
                          @click="editAppointment(appointment)"
                          class="text-green-600 hover:text-green-900"
                        >
                          <PencilIcon class="h-5 w-5" />
                        </button>
                        <button
                          @click="updateStatus(appointment)"
                          class="text-blue-600 hover:text-blue-900"
                        >
                          <CheckIcon class="h-5 w-5" />
                        </button>
                        <button
                          @click="deleteAppointment(appointment)"
                          class="text-red-600 hover:text-red-900"
                        >
                          <TrashIcon class="h-5 w-5" />
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
              
              <!-- Empty State -->
              <div v-if="appointments.length === 0 && !loading" class="text-center py-12">
                <CalendarDaysIcon class="mx-auto h-12 w-12 text-gray-400" />
                <h3 class="mt-2 text-sm font-medium text-gray-900">No appointments found</h3>
                <p class="mt-1 text-sm text-gray-500">Get started by booking a new appointment.</p>
                <div class="mt-6">
                  <button
                    @click="showBookingModal = true"
                    type="button"
                    class="inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
                  >
                    <PlusIcon class="-ml-0.5 mr-1.5 h-5 w-5" />
                    Book Appointment
                  </button>
                </div>
              </div>

              <!-- Loading State -->
              <div v-if="loading" class="text-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
                <p class="mt-2 text-sm text-gray-500">Loading appointments...</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination (Table View Only) -->
      <div v-if="currentView === 'table' && appointments.length > 0" class="mt-6 flex items-center justify-between">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="previousPage"
            :disabled="pagination.offset === 0"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="!pagination.hasMore"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing
              <span class="font-medium">{{ pagination.offset + 1 }}</span>
              to
              <span class="font-medium">{{ Math.min(pagination.offset + pagination.limit, pagination.totalCount) }}</span>
              of
              <span class="font-medium">{{ pagination.totalCount }}</span>
              results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
              <button
                @click="previousPage"
                :disabled="pagination.offset === 0"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
              >
                <ChevronLeftIcon class="h-5 w-5" />
              </button>
              <button
                @click="nextPage"
                :disabled="!pagination.hasMore"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
              >
                <ChevronRightIcon class="h-5 w-5" />
              </button>
            </nav>
          </div>
        </div>
      </div>

      <!-- Calendar View -->
      <div v-if="currentView === 'calendar'" class="mt-8">
        <AppointmentCalendar
          :appointments="appointments"
          :entity-timezone="entityStore.entityTimezone"
          @appointment-click="viewAppointment"
          @show-more="handleShowMore"
          @month-changed="handleMonthChanged"
        />
      </div>
    </div>

    <!-- Booking Modal -->
    <AppointmentBookingModal
      v-if="showBookingModal"
      @close="showBookingModal = false"
      @appointment-booked="handleAppointmentBooked"
    />

    <!-- View Modal -->
    <AppointmentViewModal
      v-if="showViewModal"
      :appointment="selectedAppointment"
      @close="showViewModal = false"
    />

    <!-- Edit Modal -->
    <AppointmentEditModal
      v-if="showEditModal"
      :appointment="selectedAppointment"
      @close="showEditModal = false"
      @appointment-updated="handleAppointmentUpdated"
    />

    <!-- Status Update Modal -->
    <AppointmentStatusModal
      v-if="showStatusModal"
      :appointment="selectedAppointment"
      @close="showStatusModal = false"
      @status-updated="handleStatusUpdated"
    />
  </div>
</template>

<script setup>
import { appointmentsApi } from '@/api/appointments'
import { useAuthStore } from '@/stores/auth'
import { useEntityStore } from '@/stores/entity'
import { formatEntityDate, formatEntityTime, getCurrentEntityDate } from '@/utils/timezoneUtils'
import {
  CalendarDaysIcon,
  CheckIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  EyeIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
  UserIcon,
  ViewColumnsIcon
} from '@heroicons/vue/24/outline'
import { computed, onMounted, reactive, ref, watch } from 'vue'

// Import modals and components
import AppointmentBookingModal from '@/components/appointments/AppointmentBookingModal.vue'
import AppointmentCalendar from '@/components/appointments/AppointmentCalendar.vue'
import AppointmentEditModal from '@/components/appointments/AppointmentEditModal.vue'
import AppointmentStatusModal from '@/components/appointments/AppointmentStatusModal.vue'
import AppointmentViewModal from '@/components/appointments/AppointmentViewModal.vue'

// Auth store
const authStore = useAuthStore()
const entityStore = useEntityStore()

// Reactive state
const appointments = ref([])
const doctors = ref([])
const stats = ref({})
const loading = ref(false)
const error = ref('')
const currentView = ref('table') // 'table' or 'calendar'

// Modal states
const showBookingModal = ref(false)
const showViewModal = ref(false)
const showEditModal = ref(false)
const showStatusModal = ref(false)
const selectedAppointment = ref(null)

// Filters
const filters = reactive({
  dateFrom: '', // Will be set to today's date after auth store is available
  dateTo: '',
  status: '',
  doctorId: ''
})

// Pagination
const pagination = reactive({
  limit: 20,
  offset: 0,
  totalCount: 0,
  hasMore: false
})

// Computed
const appointmentParams = computed(() => ({
  limit: pagination.limit,
  offset: pagination.offset,
  ...(filters.dateFrom && { date_from: filters.dateFrom }),
  ...(filters.dateTo && { date_to: filters.dateTo }),
  ...(filters.status && { status: filters.status }),
  ...(filters.doctorId && { doctor_id: filters.doctorId })
}))

// Methods
const loadAppointments = async () => {
  try {
    loading.value = true
    error.value = ''
    
    const response = await appointmentsApi.getAppointments(appointmentParams.value)
    appointments.value = response.data.appointments || []
    pagination.totalCount = response.data.total_count || 0
    pagination.hasMore = response.data.has_more || false
  } catch (err) {
    error.value = 'Failed to load appointments'
    console.error('Load appointments error:', err)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    // This would be a new API endpoint for statistics
    // For now, we'll calculate basic stats from appointments
    const today = new Date().toISOString().split('T')[0]
    const response = await appointmentsApi.getAppointments({
      date_from: today,
      date_to: today,
      limit: 100
    })
    
    stats.value = {
      today_appointments: response.data.appointments?.length || 0,
      total_appointments: appointments.value.length,
      by_status: appointments.value.reduce((acc, app) => {
        acc[app.status] = (acc[app.status] || 0) + 1
        return acc
      }, {})
    }
  } catch (err) {
    console.error('Load stats error:', err)
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

const previousPage = () => {
  if (pagination.offset > 0) {
    pagination.offset = Math.max(0, pagination.offset - pagination.limit)
    loadAppointments()
  }
}

const nextPage = () => {
  if (pagination.hasMore) {
    pagination.offset += pagination.limit
    loadAppointments()
  }
}

// Appointment actions
const viewAppointment = (appointment) => {
  selectedAppointment.value = appointment
  showViewModal.value = true
}

const editAppointment = (appointment) => {
  selectedAppointment.value = appointment
  showEditModal.value = true
}

const updateStatus = (appointment) => {
  selectedAppointment.value = appointment
  showStatusModal.value = true
}

const deleteAppointment = async (appointment) => {
  if (!confirm('Are you sure you want to delete this appointment?')) {
    return
  }

  try {
    await appointmentsApi.deleteAppointment(appointment.id)
    await loadAppointments()
    await loadStats()
  } catch (err) {
    error.value = 'Failed to delete appointment'
    console.error('Delete appointment error:', err)
  }
}

// Event handlers
const handleAppointmentBooked = () => {
  showBookingModal.value = false
  loadAppointments()
  loadStats()
}

const handleAppointmentUpdated = () => {
  showEditModal.value = false
  loadAppointments()
  loadStats()
}

const handleStatusUpdated = () => {
  showStatusModal.value = false
  loadAppointments()
  loadStats()
}

// Calendar event handlers
const handleShowMore = (day) => {
  // Could open a modal showing all appointments for that day
}

const handleMonthChanged = (newDate) => {
  // Could load appointments for the new month if needed
}

// Utility functions with timezone awareness
const formatDate = (utcDateTimeString) => {
  if (!utcDateTimeString) return ''
  const entityTimezone = entityStore.entityTimezone
  if (!entityTimezone) return new Date(utcDateTimeString).toLocaleDateString()
  
  try {
    return formatEntityDate(utcDateTimeString, entityTimezone)
  } catch (error) {
    console.warn('Failed to format date with timezone:', error)
    return new Date(utcDateTimeString).toLocaleDateString()
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

// Watchers
watch(filters, () => {
  pagination.offset = 0
  loadAppointments()
}, { deep: true })

// Lifecycle
onMounted(async () => {
  // Set dateFrom to today in entity timezone
  const entityTimezone = entityStore.entityTimezone
  if (entityTimezone) {
    filters.dateFrom = getCurrentEntityDate(entityTimezone)
  } else {
    console.warn('Entity timezone not available - cannot set default date filter')
  }
  
  await Promise.all([
    loadAppointments(),
    loadDoctors(),
    loadStats()
  ])
})
</script>