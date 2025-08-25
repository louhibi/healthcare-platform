<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-lg font-medium text-gray-900">Room Management</h2>
        <p class="text-sm text-gray-500">Manage rooms and their availability for your healthcare entity</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700"
      >
        <PlusIcon class="h-4 w-4 mr-2" />
        Add Room
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white shadow rounded-lg p-4 mb-6">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div>
          <label for="room-type-filter" class="block text-sm font-medium text-gray-700">Room Type</label>
          <select
            id="room-type-filter"
            v-model="filters.roomType"
            @change="applyFilters"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Types</option>
            <option value="consultation">Consultation</option>
            <option value="examination">Examination</option>
            <option value="procedure">Procedure</option>
            <option value="operating">Operating</option>
            <option value="emergency">Emergency</option>
          </select>
        </div>
        <div>
          <label for="floor-filter" class="block text-sm font-medium text-gray-700">Floor</label>
          <select
            id="floor-filter"
            v-model="filters.floor"
            @change="applyFilters"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Floors</option>
            <option v-for="floor in availableFloors" :key="floor" :value="floor">
              Floor {{ floor }}
            </option>
          </select>
        </div>
        <div>
          <label for="status-filter" class="block text-sm font-medium text-gray-700">Status</label>
          <select
            id="status-filter"
            v-model="filters.status"
            @change="applyFilters"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-2 text-sm text-gray-500">Loading rooms...</p>
    </div>

    <!-- Rooms Grid -->
    <div v-else-if="filteredRooms.length > 0" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="room in filteredRooms"
        :key="room.id"
        class="bg-white overflow-hidden shadow rounded-lg"
      >
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="h-10 w-10 rounded-lg bg-indigo-100 flex items-center justify-center">
                  <span class="text-lg">{{ getRoomTypeIcon(room.room_type) }}</span>
                </div>
              </div>
              <div class="ml-3">
                <h3 class="text-lg font-medium text-gray-900">{{ room.room_number }}</h3>
                <p class="text-sm text-gray-500">{{ room.room_name || room.room_type }}</p>
              </div>
            </div>
            <div class="flex space-x-2">
              <button
                @click="editRoom(room)"
                class="text-indigo-600 hover:text-indigo-900"
              >
                <PencilIcon class="h-4 w-4" />
              </button>
              <button
                @click="toggleRoomStatus(room)"
                class="text-gray-600 hover:text-gray-900"
              >
                <PowerIcon class="h-4 w-4" />
              </button>
              <button
                @click="deleteRoom(room)"
                class="text-red-600 hover:text-red-900"
              >
                <TrashIcon class="h-4 w-4" />
              </button>
            </div>
          </div>

          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Type:</span>
              <span class="font-medium capitalize">{{ room.room_type.replace('-', ' ') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Floor:</span>
              <span class="font-medium">{{ room.floor }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Department:</span>
              <span class="font-medium">{{ room.department || 'General' }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Capacity:</span>
              <span class="font-medium">{{ room.capacity }} {{ room.capacity === 1 ? 'person' : 'people' }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Status:</span>
              <span :class="[
                'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
                room.is_active 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              ]">
                {{ room.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>
          </div>

          <div v-if="room.equipment" class="mt-4">
            <h4 class="text-sm font-medium text-gray-900 mb-2">Equipment:</h4>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="equipment in getEquipmentList(room.equipment)"
                :key="equipment.trim()"
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-100 text-gray-800"
              >
                {{ equipment.trim() }}
              </span>
            </div>
          </div>

          <div v-if="room.notes" class="mt-4">
            <h4 class="text-sm font-medium text-gray-900 mb-1">Notes:</h4>
            <p class="text-sm text-gray-600">{{ room.notes }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <BuildingOfficeIcon class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">No rooms found</h3>
      <p class="mt-1 text-sm text-gray-500">
        {{ filters.roomType || filters.floor || filters.status ? 'Try adjusting your filters or create a new room.' : 'Get started by creating your first room.' }}
      </p>
      <div class="mt-6">
        <button
          @click="showCreateModal = true"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
          <PlusIcon class="h-4 w-4 mr-2" />
          Add Room
        </button>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <TransitionRoot as="template" :show="showCreateModal || showEditModal">
      <Dialog as="div" class="relative z-50" @close="closeModal">
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
                  <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-indigo-100">
                    <BuildingOfficeIcon class="h-6 w-6 text-indigo-600" />
                  </div>
                  <div class="mt-3 text-center sm:mt-5">
                    <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                      {{ isEditing ? 'Edit Room' : 'Add New Room' }}
                    </DialogTitle>
                  </div>
                </div>

                <form @submit.prevent="saveRoom" class="mt-6 space-y-6">
                  <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                    <div>
                      <label for="room-number" class="block text-sm font-medium text-gray-700">
                        Room Number *
                      </label>
                      <input
                        id="room-number"
                        v-model="form.room_number"
                        type="text"
                        required
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>

                    <div>
                      <label for="room-name" class="block text-sm font-medium text-gray-700">
                        Room Name
                      </label>
                      <input
                        id="room-name"
                        v-model="form.room_name"
                        type="text"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>

                    <div>
                      <label for="room-type" class="block text-sm font-medium text-gray-700">
                        Room Type *
                      </label>
                      <select
                        id="room-type"
                        v-model="form.room_type"
                        required
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      >
                        <option value="">Select type...</option>
                        <option value="consultation">Consultation</option>
                        <option value="examination">Examination</option>
                        <option value="procedure">Procedure</option>
                        <option value="operating">Operating</option>
                        <option value="emergency">Emergency</option>
                      </select>
                    </div>

                    <div>
                      <label for="floor" class="block text-sm font-medium text-gray-700">
                        Floor
                      </label>
                      <input
                        id="floor"
                        v-model.number="form.floor"
                        type="number"
                        min="1"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>

                    <div>
                      <label for="department" class="block text-sm font-medium text-gray-700">
                        Department
                      </label>
                      <input
                        id="department"
                        v-model="form.department"
                        type="text"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>

                    <div>
                      <label for="capacity" class="block text-sm font-medium text-gray-700">
                        Capacity
                      </label>
                      <input
                        id="capacity"
                        v-model.number="form.capacity"
                        type="number"
                        min="1"
                        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>
                  </div>

                  <div>
                    <label for="equipment" class="block text-sm font-medium text-gray-700">
                      Equipment
                    </label>
                    <textarea
                      id="equipment"
                      v-model="form.equipment"
                      rows="3"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      placeholder="Comma-separated list of equipment..."
                    />
                  </div>

                  <div>
                    <label for="notes" class="block text-sm font-medium text-gray-700">
                      Notes
                    </label>
                    <textarea
                      id="notes"
                      v-model="form.notes"
                      rows="3"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    />
                  </div>

                  <div class="flex items-center">
                    <input
                      id="is-active"
                      v-model="form.is_active"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label for="is-active" class="ml-2 block text-sm text-gray-900">
                      Active
                    </label>
                  </div>

                  <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                    <button
                      type="submit"
                      :disabled="saving"
                      class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2 disabled:opacity-50"
                    >
                      <span v-if="!saving">{{ isEditing ? 'Update' : 'Create' }}</span>
                      <span v-else>{{ isEditing ? 'Updating...' : 'Creating...' }}</span>
                    </button>
                    <button
                      type="button"
                      @click="closeModal"
                      class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
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

    <!-- Delete Confirmation Modal -->
    <TransitionRoot as="template" :show="showDeleteModal">
      <Dialog as="div" class="relative z-50" @close="showDeleteModal = false">
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
                <div class="sm:flex sm:items-start">
                  <div class="mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                    <ExclamationTriangleIcon class="h-6 w-6 text-red-600" />
                  </div>
                  <div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                    <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                      Delete Room
                    </DialogTitle>
                    <div class="mt-2">
                      <p class="text-sm text-gray-500">
                        Are you sure you want to delete room {{ roomToDelete?.room_number }}? This action cannot be undone.
                      </p>
                    </div>
                  </div>
                </div>
                <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
                  <button
                    @click="confirmDelete"
                    :disabled="deleting"
                    class="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 sm:ml-3 sm:w-auto disabled:opacity-50"
                  >
                    {{ deleting ? 'Deleting...' : 'Delete' }}
                  </button>
                  <button
                    @click="showDeleteModal = false"
                    class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                  >
                    Cancel
                  </button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </TransitionRoot>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
  BuildingOfficeIcon,
  PlusIcon,
  PencilIcon,
  PowerIcon,
  TrashIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'

// Emits
const emit = defineEmits(['stats-updated'])

// Reactive state
const rooms = ref([])
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const roomToDelete = ref(null)

const filters = reactive({
  roomType: '',
  floor: '',
  status: ''
})

const form = reactive({
  room_number: '',
  room_name: '',
  room_type: '',
  floor: 1,
  department: '',
  capacity: 1,
  equipment: '',
  notes: '',
  is_active: true
})

// Computed
const filteredRooms = computed(() => {
  return rooms.value.filter(room => {
    if (filters.roomType && room.room_type !== filters.roomType) return false
    if (filters.floor && room.floor !== parseInt(filters.floor)) return false
    if (filters.status === 'active' && !room.is_active) return false
    if (filters.status === 'inactive' && room.is_active) return false
    return true
  })
})

const availableFloors = computed(() => {
  const floors = [...new Set(rooms.value.map(room => room.floor))].sort((a, b) => a - b)
  return floors
})

// Methods
const loadRooms = async () => {
  try {
    loading.value = true
    const response = await appointmentsApi.getRooms()
    rooms.value = response.data || []
    emit('stats-updated')
  } catch (error) {
    console.error('Error loading rooms:', error)
  } finally {
    loading.value = false
  }
}

const saveRoom = async () => {
  try {
    saving.value = true
    
    if (isEditing.value) {
      await appointmentsApi.updateRoom(editingId.value, form)
    } else {
      await appointmentsApi.createRoom(form)
    }
    
    await loadRooms()
    closeModal()
  } catch (error) {
    console.error('Error saving room:', error)
  } finally {
    saving.value = false
  }
}

const editRoom = (room) => {
  // Convert equipment to string format for editing
  let equipmentStr = ''
  if (room.equipment) {
    if (Array.isArray(room.equipment)) {
      equipmentStr = room.equipment.join(', ')
    } else if (typeof room.equipment === 'string') {
      equipmentStr = room.equipment
    }
  }

  Object.assign(form, {
    room_number: room.room_number,
    room_name: room.room_name || '',
    room_type: room.room_type,
    floor: room.floor,
    department: room.department || '',
    capacity: room.capacity,
    equipment: equipmentStr,
    notes: room.notes || '',
    is_active: room.is_active
  })
  
  isEditing.value = true
  editingId.value = room.id
  showEditModal.value = true
}

const toggleRoomStatus = async (room) => {
  try {
    await appointmentsApi.updateRoom(room.id, {
      ...room,
      is_active: !room.is_active
    })
    await loadRooms()
  } catch (error) {
    console.error('Error toggling room status:', error)
  }
}

const deleteRoom = (room) => {
  roomToDelete.value = room
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  try {
    deleting.value = true
    await appointmentsApi.deleteRoom(roomToDelete.value.id)
    await loadRooms()
    showDeleteModal.value = false
    roomToDelete.value = null
  } catch (error) {
    console.error('Error deleting room:', error)
  } finally {
    deleting.value = false
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  isEditing.value = false
  editingId.value = null
  
  // Reset form
  Object.assign(form, {
    room_number: '',
    room_name: '',
    room_type: '',
    floor: 1,
    department: '',
    capacity: 1,
    equipment: '',
    notes: '',
    is_active: true
  })
}

const applyFilters = () => {
  // Filters are applied via computed property
}

const getRoomTypeIcon = (type) => {
  const icons = {
    'consultation': '💬',
    'examination': '🔍',
    'procedure': '⚕️',
    'operating': '🏥',
    'emergency': '🚨'
  }
  return icons[type] || '🏠'
}

const getEquipmentList = (equipment) => {
  if (!equipment) return []
  if (Array.isArray(equipment)) return equipment
  if (typeof equipment === 'string') return equipment.split(',')
  return []
}

// Lifecycle
onMounted(() => {
  loadRooms()
})
</script>