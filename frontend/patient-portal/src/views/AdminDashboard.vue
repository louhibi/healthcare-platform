<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <CogIcon class="h-8 w-8 text-indigo-600" />
            <h1 class="ml-3 text-2xl font-bold text-gray-900">Admin Dashboard</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-500">Healthcare Entity: {{ currentEntity?.name }}</span>
            <button
              @click="$router.push('/')"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              <ArrowLeftIcon class="h-4 w-4 mr-2" />
              Back to Portal
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Stats Overview -->
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <ClockIcon class="h-6 w-6 text-gray-400" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Duration Settings</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.durationSettings }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <BuildingOfficeIcon class="h-6 w-6 text-gray-400" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Active Rooms</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.activeRooms }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <CalendarDaysIcon class="h-6 w-6 text-gray-400" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Today's Appointments</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.todayAppointments }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <UserGroupIcon class="h-6 w-6 text-gray-400" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Active Patients</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.activePatients }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-white shadow rounded-lg mb-8">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Quick Actions</h3>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
            <button
              @click="activeTab = 'duration-settings'"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <ClockIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">Manage Duration Settings</span>
            </button>

            <button
              @click="activeTab = 'room-management'"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <BuildingOfficeIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">Manage Rooms</span>
            </button>

            <button
              @click="activeTab = 'doctor-management'"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <UserIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">Manage Doctors</span>
            </button>

            <button
              @click="$router.push('/appointments')"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <CalendarDaysIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">View Appointments</span>
            </button>

            <button
              @click="activeTab = 'form-configuration'"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <ClipboardDocumentListIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">Configure Forms</span>
            </button>

            <button
              @click="$router.push('/form-demo')"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <span class="mt-2 block text-sm font-medium text-gray-900">Form Demo</span>
            </button>

            <button
              @click="$router.push('/date-test')"
              class="relative block w-full p-6 border-2 border-dashed border-gray-300 rounded-lg text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <CalendarIcon class="mx-auto h-8 w-8 text-gray-400" />
              <span class="mt-2 block text-sm font-medium text-gray-900">Date Localization Test</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content Tabs -->
      <div class="bg-white shadow rounded-lg">
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex" aria-label="Tabs">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                activeTab === tab.id
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                'flex-1 py-2 px-1 text-center border-b-2 font-medium text-sm'
              ]"
            >
              <component :is="tab.icon" class="h-5 w-5 mx-auto mb-1" />
              {{ tab.name }}
            </button>
          </nav>
        </div>

        <div class="p-6">
          <!-- Duration Settings Tab -->
          <DurationSettings v-if="activeTab === 'duration-settings'" @stats-updated="loadStats" />
          
          <!-- Room Management Tab -->
          <RoomManagement v-if="activeTab === 'room-management'" @stats-updated="loadStats" />

          <!-- Doctor Management Tab -->
          <DoctorManagement v-if="activeTab === 'doctor-management'" @stats-updated="loadStats" />

          <!-- Form Configuration Tab -->
          <FormConfigManager v-if="activeTab === 'form-configuration'" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import {
  CogIcon,
  ClockIcon,
  BuildingOfficeIcon,
  CalendarDaysIcon,
  CalendarIcon,
  UserGroupIcon,
  UserIcon,
  ArrowLeftIcon,
  ClipboardDocumentListIcon
} from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'
import { patientsApi } from '@/api/patients'
import DurationSettings from '@/components/admin/DurationSettings.vue'
import RoomManagement from '@/components/admin/RoomManagement.vue'
import DoctorManagement from '@/components/admin/DoctorManagement.vue'
import FormConfigManager from '@/components/admin/FormConfigManager.vue'

// Stores
const authStore = useAuthStore()

// Reactive state
const activeTab = ref('duration-settings')
const stats = reactive({
  durationSettings: 0,
  activeRooms: 0,
  todayAppointments: 0,
  activePatients: 0
})

// Computed
const currentEntity = computed(() => authStore.currentEntity)

const tabs = [
  { id: 'duration-settings', name: 'Duration Settings', icon: ClockIcon },
  { id: 'room-management', name: 'Room Management', icon: BuildingOfficeIcon },
  { id: 'doctor-management', name: 'Doctor Management', icon: UserIcon },
  { id: 'form-configuration', name: 'Form Configuration', icon: ClipboardDocumentListIcon }
]

// Methods
const loadStats = async () => {
  try {
    // Load duration settings count
    const durationResponse = await appointmentsApi.getDurationSettings()
    stats.durationSettings = durationResponse.data?.length || 0

    // Load rooms count
    const roomsResponse = await appointmentsApi.getRooms()
    stats.activeRooms = roomsResponse.data?.filter(room => room.is_active)?.length || 0

    // Load today's appointments count
    const today = new Date().toISOString().split('T')[0]
    const appointmentsResponse = await appointmentsApi.getAppointments({ date: today })
    stats.todayAppointments = appointmentsResponse.data?.length || 0

    // Load active patients count
    const patientsResponse = await patientsApi.getPatients({ is_active: true, limit: 1000 })
    stats.activePatients = patientsResponse.total_count || 0

  } catch (error) {
    console.error('Error loading admin stats:', error)
  }
}

// Lifecycle
onMounted(() => {
  loadStats()
})
</script>