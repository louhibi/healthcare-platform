<template>
  <div class="min-h-screen bg-gray-50 p-6">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Doctor Availability Management</h1>
      <p class="mt-2 text-gray-600">Manage doctor schedules, availability, and working hours</p>
    </div>

    <!-- Quick Actions Bar -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div class="flex flex-wrap items-center gap-4">
          <!-- Doctor Selection -->
          <div class="flex items-center space-x-2">
            <label class="text-sm font-medium text-gray-700">Doctor:</label>
            <select
              v-model="selectedDoctorId"
              @change="handleDoctorChange"
              class="border border-gray-300 rounded-md px-3 py-1.5 text-sm focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Doctors</option>
              <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                Dr. {{ doctor.first_name }} {{ doctor.last_name }}
                <span v-if="doctor.specialization" class="text-gray-500">({{ doctor.specialization }})</span>
              </option>
            </select>
          </div>

          <!-- Date Range Selection -->
          <div class="flex items-center space-x-2">
            <label class="text-sm font-medium text-gray-700">Date Range:</label>
            <input
              v-model="dateFrom"
              type="date"
              @change="loadAvailability"
              class="border border-gray-300 rounded-md px-3 py-1.5 text-sm focus:ring-blue-500 focus:border-blue-500"
            />
            <span class="text-gray-500">to</span>
            <input
              v-model="dateTo"
              type="date"
              @change="loadAvailability"
              class="border border-gray-300 rounded-md px-3 py-1.5 text-sm focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <!-- Status Filter -->
          <div class="flex items-center space-x-2">
            <label class="text-sm font-medium text-gray-700">Status:</label>
            <select
              v-model="selectedStatus"
              @change="loadAvailability"
              class="border border-gray-300 rounded-md px-3 py-1.5 text-sm focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">All Statuses</option>
              <option value="available">Available</option>
              <option value="unavailable">Unavailable</option>
              <option value="vacation">Vacation</option>
              <option value="training">Training</option>
              <option value="sick_leave">Sick Leave</option>
              <option value="meeting">Meeting</option>
            </select>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center space-x-3">
          <button
            @click="showCalendarView = !showCalendarView"
            :class="showCalendarView ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 border border-gray-300'"
            class="px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-50 transition-colors"
          >
            {{ showCalendarView ? '📋 List View' : '📅 Calendar View' }}
          </button>
          <button
            @click="showAddModal = true"
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
          >
            ➕ Add Availability
          </button>
          <button
            @click="showBulkModal = true"
            class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
          >
            📅 Bulk Schedule
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-2 text-gray-600">Loading availability...</span>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
      <div class="flex">
        <div class="text-red-400">⚠️</div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <p class="mt-1 text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Calendar View -->
    <div v-if="showCalendarView && selectedDoctorId">
      <AvailabilityCalendar
        :doctor-id="selectedDoctorId"
        :availability-data="availabilityList"
        @edit-availability="editAvailability"
        @create-availability="createAvailabilityForDate"
      />
    </div>

    <!-- List View -->
    <div v-else-if="!showCalendarView" class="space-y-6">
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4 mb-6">
        <div
          v-for="(count, status) in statusSummary"
          :key="status"
          class="bg-white rounded-lg p-4 border border-gray-200 shadow-sm"
        >
          <div class="flex items-center">
            <div :class="getStatusIcon(status)" class="text-2xl mr-3">
              {{ getStatusEmoji(status) }}
            </div>
            <div>
              <p class="text-2xl font-bold text-gray-900">{{ count }}</p>
              <p class="text-sm text-gray-600 capitalize">{{ status.replace('_', ' ') }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Availability List -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">
            Doctor Availability
            <span class="text-sm text-gray-500 font-normal">({{ availabilityList.length }} records)</span>
          </h2>
        </div>

        <div v-if="availabilityList.length === 0" class="p-12 text-center">
          <div class="text-6xl text-gray-300 mb-4">📅</div>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No availability records found</h3>
          <p class="text-gray-600 mb-6">Create availability schedules to manage doctor working hours</p>
          <button
            @click="showAddModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
          >
            Create First Schedule
          </button>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Doctor
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Working Hours
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Break Time
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Notes
                </th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="availability in availabilityList"
                :key="availability.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3">
                      <span class="text-sm font-medium text-blue-600">
                        {{ availability.doctor_name?.split(' ')[1]?.[0] || 'D' }}
                      </span>
                    </div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">{{ availability.doctor_name || 'Doctor' }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-900">{{ formatDate(availability.start_datetime) }}</div>
                  <div class="text-xs text-gray-500">{{ getDayOfWeek(availability.start_datetime) }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getStatusClass(availability.status)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                    {{ getStatusEmoji(availability.status) }} {{ getStatusLabel(availability.status) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <div v-if="availability.start_datetime && availability.end_datetime">
                    {{ formatTime(availability.start_datetime) }} - {{ formatTime(availability.end_datetime) }}
                  </div>
                  <div v-else class="text-gray-400">Not set</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <div v-if="availability.break_start_datetime && availability.break_end_datetime">
                    {{ formatTime(availability.break_start_datetime) }} - {{ formatTime(availability.break_end_datetime) }}
                  </div>
                  <div v-else class="text-gray-400">No break</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-900 max-w-xs truncate" :title="availability.notes">
                    {{ availability.notes || 'No notes' }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div class="flex items-center justify-end space-x-2">
                    <button
                      @click="viewAvailability(availability)"
                      class="text-blue-600 hover:text-blue-900 transition-colors"
                      title="View Details"
                    >
                      👁️
                    </button>
                    <button
                      @click="editAvailability(availability)"
                      class="text-green-600 hover:text-green-900 transition-colors"
                      title="Edit"
                    >
                      ✏️
                    </button>
                    <button
                      @click="deleteAvailability(availability)"
                      class="text-red-600 hover:text-red-900 transition-colors"
                      title="Delete"
                    >
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Prompt to Select Doctor for Calendar View -->
    <div v-else class="bg-yellow-50 border border-yellow-200 rounded-md p-6 text-center">
      <div class="text-4xl text-yellow-500 mb-3">👩‍⚕️</div>
      <h3 class="text-lg font-medium text-yellow-800 mb-2">Select a Doctor</h3>
      <p class="text-yellow-700">Please select a specific doctor above to view their availability calendar.</p>
    </div>

    <!-- Add/Edit Availability Modal -->
    <AvailabilityModal
      v-if="showAddModal || editingAvailability"
      :availability="editingAvailability"
      :doctors="doctors"
      @close="closeModal"
      @save="handleSaveAvailability"
    />

    <!-- Bulk Schedule Modal -->
    <BulkScheduleModal
      v-if="showBulkModal"
      :doctors="doctors"
      @close="showBulkModal = false"
      @save="handleBulkSchedule"
    />

    <!-- View Availability Modal -->
    <ViewAvailabilityModal
      v-if="viewingAvailability"
      :availability="viewingAvailability"
      @close="viewingAvailability = null"
      @edit="editAvailability"
    />
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { appointmentsApi } from '../api/appointments'
import AvailabilityCalendar from '../components/AvailabilityCalendar.vue'
import AvailabilityModal from '../components/AvailabilityModal.vue'
import BulkScheduleModal from '../components/BulkScheduleModal.vue'
import ViewAvailabilityModal from '../components/ViewAvailabilityModal.vue'
import { formatDate, getDayOfWeek, getTodayString, formatTimeFromDateTime } from '../utils/dateUtils'

export default {
  name: 'DoctorAvailability',
  components: {
    AvailabilityCalendar,
    AvailabilityModal,
    BulkScheduleModal,
    ViewAvailabilityModal
  },
  setup() {
    const toast = useToast()

    // Reactive data
    const loading = ref(false)
    const error = ref('')
    const doctors = ref([])
    const availabilityList = ref([])
    const selectedDoctorId = ref('')
    const selectedStatus = ref('')
    const showCalendarView = ref(false)
    const showAddModal = ref(false)
    const showBulkModal = ref(false)
    const editingAvailability = ref(null)
    const viewingAvailability = ref(null)

    // Date range
    const todayString = getTodayString()
    const dateFrom = ref(todayString)
    const dateTo = ref(new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0])

    // Computed properties
    const statusSummary = computed(() => {
      const summary = {
        available: 0,
        unavailable: 0,
        vacation: 0,
        training: 0,
        sick_leave: 0,
        meeting: 0
      }
      
      availabilityList.value.forEach(availability => {
        if (summary.hasOwnProperty(availability.status)) {
          summary[availability.status]++
        }
      })
      
      return summary
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
      return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
    }

    const getStatusIcon = (status) => {
      const icons = {
        available: 'text-green-600',
        unavailable: 'text-gray-600',
        vacation: 'text-blue-600',
        training: 'text-yellow-600',
        sick_leave: 'text-red-600',
        meeting: 'text-purple-600'
      }
      return icons[status] || 'text-gray-600'
    }

    // Date formatting functions are now imported from utilities
    const formatTime = (datetime) => {
      return formatTimeFromDateTime(datetime)
    }

    // API calls
    const loadDoctors = async () => {
      try {
        const response = await appointmentsApi.getDoctorsByEntity()
        doctors.value = response.data || []
      } catch (err) {
        console.error('Error loading doctors:', err)
        toast.error('Failed to load doctors')
      }
    }

    const loadAvailability = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const params = {
          date_from: dateFrom.value,
          date_to: dateTo.value,
          limit: 100
        }
        
        if (selectedDoctorId.value) {
          params.doctor_id = selectedDoctorId.value
        }
        
        if (selectedStatus.value) {
          params.status = selectedStatus.value
        }

        const response = await appointmentsApi.getDoctorAvailability(params)
        availabilityList.value = response.data || []
      } catch (err) {
        error.value = err.message || 'Failed to load availability'
        console.error('Error loading availability:', err)
      } finally {
        loading.value = false
      }
    }

    // Event handlers
    const handleDoctorChange = () => {
      loadAvailability()
    }

    const viewAvailability = (availability) => {
      viewingAvailability.value = availability
    }

    const editAvailability = (availability) => {
      editingAvailability.value = { ...availability }
      showAddModal.value = true
    }

    const createAvailabilityForDate = (date) => {
      editingAvailability.value = {
        doctor_id: selectedDoctorId.value,
        start_datetime: `${date}T09:00:00Z`,
        end_datetime: `${date}T17:00:00Z`,
        break_start_datetime: `${date}T12:00:00Z`,
        break_end_datetime: `${date}T13:00:00Z`,
        status: 'available',
        notes: ''
      }
      showAddModal.value = true
    }

    const closeModal = () => {
      showAddModal.value = false
      editingAvailability.value = null
    }

    const handleSaveAvailability = async (availabilityData) => {
      try {
        if (editingAvailability.value && editingAvailability.value.id) {
          // Update existing
          await appointmentsApi.updateDoctorAvailability(editingAvailability.value.id, availabilityData)
          toast.success('Availability updated successfully')
        } else {
          // Create new
          await appointmentsApi.createDoctorAvailability(availabilityData)
          toast.success('Availability created successfully')
        }
        
        closeModal()
        loadAvailability()
      } catch (err) {
        toast.error(err.message || 'Failed to save availability')
        console.error('Error saving availability:', err)
      }
    }

    const handleBulkSchedule = async (bulkData) => {
      try {
        await appointmentsApi.createBulkAvailability(bulkData)
        toast.success('Bulk schedule created successfully')
        showBulkModal.value = false
        loadAvailability()
      } catch (err) {
        toast.error(err.message || 'Failed to create bulk schedule')
        console.error('Error creating bulk schedule:', err)
      }
    }

    const deleteAvailability = async (availability) => {
      if (!confirm(`Are you sure you want to delete availability for ${availability.doctor_name || 'Doctor'} on ${formatDate(availability.start_datetime)}?`)) {
        return
      }

      try {
        await appointmentsApi.deleteDoctorAvailability(availability.id)
        toast.success('Availability deleted successfully')
        loadAvailability()
      } catch (err) {
        toast.error(err.message || 'Failed to delete availability')
        console.error('Error deleting availability:', err)
      }
    }

    // Lifecycle
    onMounted(async () => {
      await loadDoctors()
      await loadAvailability()
    })

    return {
      // Data
      loading,
      error,
      doctors,
      availabilityList,
      selectedDoctorId,
      selectedStatus,
      dateFrom,
      dateTo,
      showCalendarView,
      showAddModal,
      showBulkModal,
      editingAvailability,
      viewingAvailability,
      
      // Computed
      statusSummary,
      
      // Methods
      getStatusClass,
      getStatusEmoji,
      getStatusLabel,
      getStatusIcon,
      formatDate,
      getDayOfWeek,
      formatTime,
      loadAvailability,
      handleDoctorChange,
      viewAvailability,
      editAvailability,
      createAvailabilityForDate,
      closeModal,
      handleSaveAvailability,
      handleBulkSchedule,
      deleteAvailability
    }
  }
}
</script>

<style scoped>
/* Custom animations and transitions */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Hover effects for interactive elements */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
}

/* Status indicator animations */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .7;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>