<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-lg font-medium text-gray-900">Appointment Duration Options</h2>
        <p class="text-sm text-gray-500">Configure multiple duration choices for each appointment type</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700"
      >
        <PlusIcon class="h-4 w-4 mr-2" />
        Add Duration Option
      </button>
    </div>

    <!-- Appointment Type Filter -->
    <div class="mb-6">
      <div class="flex space-x-4">
        <div>
          <label for="type-filter" class="block text-sm font-medium text-gray-700">Filter by Type</label>
          <select
            id="type-filter"
            v-model="selectedType"
            @change="applyFilter"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          >
            <option value="">All Types</option>
            <option value="consultation">Consultation</option>
            <option value="follow-up">Follow-up</option>
            <option value="procedure">Procedure</option>
            <option value="emergency">Emergency</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-2 text-sm text-gray-500">Loading duration settings...</p>
    </div>

    <!-- Duration Options by Type -->
    <div v-else-if="groupedOptions && Object.keys(groupedOptions).length > 0" class="space-y-6">
      <div v-for="(options, type) in filteredGroupedOptions" :key="type" class="bg-white shadow overflow-hidden sm:rounded-lg">
        <div class="px-4 py-5 sm:px-6 bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center mr-3">
                <span class="text-indigo-600 text-sm font-medium">
                  {{ getTypeIcon(type) }}
                </span>
              </div>
              <div>
                <h3 class="text-lg leading-6 font-medium text-gray-900 capitalize">
                  {{ type.replace('-', ' ') }}
                </h3>
                <p class="mt-1 max-w-2xl text-sm text-gray-500">
                  {{ options.length }} duration option{{ options.length !== 1 ? 's' : '' }} configured
                </p>
              </div>
            </div>
            <button
              @click="addOptionForType(type)"
              class="inline-flex items-center px-3 py-2 border border-transparent text-xs font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200"
            >
              <PlusIcon class="h-3 w-3 mr-1" />
              Add Option
            </button>
          </div>
        </div>
        
        <div class="border-t border-gray-200">
          <dl>
            <div v-for="(option, index) in options" :key="option.id" :class="index % 2 === 0 ? 'bg-gray-50' : 'bg-white'">
              <div class="px-4 py-5 sm:grid sm:grid-cols-6 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 flex items-center">
                  <ClockIcon class="h-4 w-4 mr-2" />
                  Duration
                </dt>
                <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-1">
                  <span class="font-medium">{{ option.duration_minutes }} minutes</span>
                  <div class="text-xs text-gray-500">{{ formatDuration(option.duration_minutes) }}</div>
                </dd>
                <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-1">
                  <span v-if="option.is_default" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    <CheckCircleIcon class="h-3 w-3 mr-1" />
                    Default
                  </span>
                  <span v-else class="text-gray-400 text-xs">Standard</span>
                </dd>
                <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-1">
                  <span :class="[
                    'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
                    option.is_active 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-red-100 text-red-800'
                  ]">
                    {{ option.is_active ? 'Active' : 'Inactive' }}
                  </span>
                </dd>
                <dd class="mt-1 text-sm text-gray-500 sm:mt-0 sm:col-span-1">
                  Order: {{ option.display_order }}
                </dd>
                <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-1">
                  <div class="flex justify-end space-x-2">
                    <button
                      @click="editOption(option)"
                      class="text-indigo-600 hover:text-indigo-900"
                      title="Edit"
                    >
                      <PencilIcon class="h-4 w-4" />
                    </button>
                    <button
                      @click="toggleDefault(option)"
                      class="text-yellow-600 hover:text-yellow-900"
                      title="Set as default"
                    >
                      <StarIcon class="h-4 w-4" :class="option.is_default ? 'fill-current' : ''" />
                    </button>
                    <button
                      @click="toggleStatus(option)"
                      class="text-gray-600 hover:text-gray-900"
                      title="Toggle status"
                    >
                      <PowerIcon class="h-4 w-4" />
                    </button>
                    <button
                      @click="deleteOption(option)"
                      class="text-red-600 hover:text-red-900"
                      title="Delete"
                    >
                      <TrashIcon class="h-4 w-4" />
                    </button>
                  </div>
                </dd>
              </div>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <ClockIcon class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">No duration settings</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by creating your first duration setting.</p>
      <div class="mt-6">
        <button
          @click="showCreateModal = true"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
          <PlusIcon class="h-4 w-4 mr-2" />
          Add Duration Setting
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
              <DialogPanel class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
                <div>
                  <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-indigo-100">
                    <ClockIcon class="h-6 w-6 text-indigo-600" />
                  </div>
                  <div class="mt-3 text-center sm:mt-5">
                    <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                      {{ isEditing ? 'Edit Duration Setting' : 'Add Duration Setting' }}
                    </DialogTitle>
                  </div>
                </div>

                <form @submit.prevent="saveOption" class="mt-6 space-y-4">
                  <div>
                    <label for="appointment-type" class="block text-sm font-medium text-gray-700">
                      Appointment Type *
                    </label>
                    <select
                      id="appointment-type"
                      v-model="form.appointment_type"
                      required
                      :disabled="isEditing"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    >
                      <option value="">Select type...</option>
                      <option value="consultation">Consultation</option>
                      <option value="follow-up">Follow-up</option>
                      <option value="procedure">Procedure</option>
                      <option value="emergency">Emergency</option>
                    </select>
                  </div>

                  <div>
                    <label for="duration" class="block text-sm font-medium text-gray-700">
                      Duration (minutes) *
                    </label>
                    <input
                      id="duration"
                      v-model.number="form.duration_minutes"
                      type="number"
                      min="15"
                      max="480"
                      required
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    />
                    <p class="mt-1 text-sm text-gray-500">Must be between 15 and 480 minutes</p>
                  </div>

                  <div class="flex items-center">
                    <input
                      id="is-default"
                      v-model="form.is_default"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label for="is-default" class="ml-2 block text-sm text-gray-900">
                      Set as default for this appointment type
                    </label>
                  </div>

                  <div>
                    <label for="display-order" class="block text-sm font-medium text-gray-700">
                      Display Order
                    </label>
                    <input
                      id="display-order"
                      v-model.number="form.display_order"
                      type="number"
                      min="1"
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                    />
                    <p class="mt-1 text-sm text-gray-500">Lower numbers appear first in the list</p>
                  </div>

                  <div class="flex items-center">
                    <input
                      id="is-default"
                      v-model="form.is_default"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label for="is-default" class="ml-2 block text-sm text-gray-900">
                      Set as default duration for this appointment type
                    </label>
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
  ClockIcon,
  PlusIcon,
  PencilIcon,
  PowerIcon,
  CheckCircleIcon,
  StarIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'

// Emits
const emit = defineEmits(['stats-updated'])

// Reactive state
const options = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const selectedType = ref('')

const form = reactive({
  appointment_type: '',
  duration_minutes: 30,
  is_default: false,
  is_active: true,
  display_order: 1
})

// Computed
const groupedOptions = computed(() => {
  return options.value.reduce((groups, option) => {
    if (!groups[option.appointment_type]) {
      groups[option.appointment_type] = []
    }
    groups[option.appointment_type].push(option)
    return groups
  }, {})
})

const filteredGroupedOptions = computed(() => {
  if (!selectedType.value) {
    return groupedOptions.value
  }
  return {
    [selectedType.value]: groupedOptions.value[selectedType.value] || []
  }
})

// Methods
const loadOptions = async () => {
  try {
    loading.value = true
    const response = await appointmentsApi.getDurationOptions()
    options.value = response.data || []
    
    // Sort options by type and display order
    options.value.sort((a, b) => {
      if (a.appointment_type !== b.appointment_type) {
        return a.appointment_type.localeCompare(b.appointment_type)
      }
      return a.display_order - b.display_order
    })
    
    emit('stats-updated')
  } catch (error) {
    console.error('Error loading duration options:', error)
  } finally {
    loading.value = false
  }
}

const saveOption = async () => {
  try {
    saving.value = true
    
    if (isEditing.value) {
      await appointmentsApi.updateDurationOption(editingId.value, form)
    } else {
      await appointmentsApi.createDurationOption(form)
    }
    
    await loadOptions()
    closeModal()
  } catch (error) {
    console.error('Error saving duration option:', error)
  } finally {
    saving.value = false
  }
}

const editOption = (option) => {
  Object.assign(form, {
    appointment_type: option.appointment_type,
    duration_minutes: option.duration_minutes,
    is_default: option.is_default,
    is_active: option.is_active,
    display_order: option.display_order
  })
  
  isEditing.value = true
  editingId.value = option.id
  showEditModal.value = true
}

const addOptionForType = (type) => {
  Object.assign(form, {
    appointment_type: type,
    duration_minutes: 30,
    is_default: false,
    is_active: true,
    display_order: getNextDisplayOrder(type)
  })
  
  isEditing.value = false
  editingId.value = null
  showCreateModal.value = true
}

const getNextDisplayOrder = (type) => {
  const typeOptions = groupedOptions.value[type] || []
  return typeOptions.length > 0 ? Math.max(...typeOptions.map(o => o.display_order)) + 1 : 1
}

const toggleStatus = async (option) => {
  try {
    await appointmentsApi.updateDurationOption(option.id, {
      ...option,
      is_active: !option.is_active
    })
    await loadOptions()
  } catch (error) {
    console.error('Error toggling option status:', error)
  }
}

const toggleDefault = async (option) => {
  try {
    // First, remove default from all other options of the same type
    const typeOptions = groupedOptions.value[option.appointment_type] || []
    const currentDefault = typeOptions.find(o => o.is_default && o.id !== option.id)
    
    if (currentDefault) {
      await appointmentsApi.updateDurationOption(currentDefault.id, {
        ...currentDefault,
        is_default: false
      })
    }
    
    // Then set this option as default
    await appointmentsApi.updateDurationOption(option.id, {
      ...option,
      is_default: !option.is_default
    })
    
    await loadOptions()
  } catch (error) {
    console.error('Error toggling default option:', error)
  }
}

const deleteOption = async (option) => {
  if (confirm(`Are you sure you want to delete the ${option.duration_minutes}-minute option for ${option.appointment_type}?`)) {
    try {
      await appointmentsApi.deleteDurationOption(option.id)
      await loadOptions()
    } catch (error) {
      console.error('Error deleting option:', error)
    }
  }
}

const applyFilter = () => {
  // Filter is handled by computed property
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  isEditing.value = false
  editingId.value = null
  
  // Reset form
  Object.assign(form, {
    appointment_type: '',
    duration_minutes: 30,
    is_default: false,
    is_active: true,
    display_order: 1
  })
}

const getTypeIcon = (type) => {
  const icons = {
    'consultation': '💬',
    'follow-up': '🔄',
    'procedure': '⚕️',
    'emergency': '🚨'
  }
  return icons[type] || '📅'
}

const formatDuration = (minutes) => {
  if (minutes < 60) {
    return `${minutes}m`
  }
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadOptions()
})
</script>