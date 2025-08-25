<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="bg-white shadow rounded-lg p-6 mb-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Form Configuration</h1>
          <p class="mt-1 text-sm text-gray-500">Configure field visibility and requirements for forms</p>
        </div>
        <div class="mt-4 sm:mt-0 flex space-x-3">
          <button
            v-if="hasChanges"
            @click="saveAllChanges"
            :disabled="isLoading"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
          >
            <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ isLoading ? 'Saving...' : 'Save Changes' }}
          </button>
          <button
            @click="refreshConfiguration"
            :disabled="isLoading"
            class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Form Type Selection -->
    <div class="bg-white shadow rounded-lg mb-6">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Select Form Type</h3>
      </div>
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <button
            v-for="formType in formTypes"
            :key="formType.name"
            @click="selectFormType(formType.name)"
            :class="[
              'relative block w-full p-6 border-2 rounded-lg text-left hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors',
              selectedFormType === formType.name
                ? 'border-indigo-500 bg-indigo-50 text-indigo-900'
                : 'border-gray-300 text-gray-900'
            ]"
          >
            <div class="flex items-center">
              <component 
                :is="getFormTypeIcon(formType.name)" 
                class="h-8 w-8 mr-3"
                :class="selectedFormType === formType.name ? 'text-indigo-600' : 'text-gray-400'"
              />
              <div>
                <span class="block text-lg font-medium">{{ formType.display_name }}</span>
                <span class="block text-sm opacity-75">{{ formType.description }}</span>
              </div>
            </div>
            <div v-if="selectedFormType === formType.type" class="absolute top-2 right-2">
              <svg class="h-6 w-6 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Form Configuration -->
    <div v-if="selectedFormType" class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <div>
          <h3 class="text-lg font-medium text-gray-900">
            {{ getFormTypeDisplayName(selectedFormType) }} Form Configuration
          </h3>
          <p class="text-sm text-gray-500">Configure field visibility, requirements, and order</p>
        </div>
        <div class="flex space-x-2">
          <button
            @click="showPreview = !showPreview"
            class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            <svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            {{ showPreview ? 'Hide Preview' : 'Show Preview' }}
          </button>
          <button
            @click="confirmResetForm"
            :disabled="isLoading"
            class="inline-flex items-center px-3 py-2 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50 disabled:opacity-50"
          >
            <svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Reset to Defaults
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading && !fields.length" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="p-6 text-center text-red-600">
        <p>{{ error }}</p>
        <button @click="loadFormConfiguration" class="mt-2 text-indigo-600 hover:text-indigo-800">Retry</button>
      </div>

      <!-- Form Fields Configuration -->
      <div v-else-if="fields.length" class="p-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Field Configuration Panel -->
          <div>
            <h4 class="text-md font-medium text-gray-900 mb-4">Field Configuration</h4>
            
            <!-- Category Tabs -->
            <div class="border-b border-gray-200 mb-4">
              <nav class="-mb-px flex space-x-8">
                <button
                  v-for="category in categories"
                  :key="category"
                  @click="selectedCategory = category"
                  :class="[
                    'py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap',
                    selectedCategory === category
                      ? 'border-indigo-500 text-indigo-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  ]"
                >
                  {{ category }}
                  <span class="ml-2 bg-gray-100 text-gray-900 py-0.5 px-2 rounded-full text-xs">
                    {{ getCategoryFieldCount(category) }}
                  </span>
                </button>
              </nav>
            </div>

            <!-- Fields List -->
            <div class="space-y-3 max-h-96 overflow-y-auto">
              <div
                v-for="field in getCategoryFields(selectedCategory)"
                :key="field.field_id"
                :class="[
                  'border rounded-lg p-4 transition-colors',
                  field.is_enabled ? 'border-gray-200 bg-white' : 'border-gray-100 bg-gray-50'
                ]"
              >
                <div class="flex items-start justify-between">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center">
                      <h5 :class="[
                        'text-sm font-medium truncate',
                        field.is_enabled ? 'text-gray-900' : 'text-gray-500'
                      ]">
                        {{ field.display_name }}
                      </h5>
                      <div class="ml-2 flex space-x-1">
                        <span v-if="field.is_core" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                          Core
                        </span>
                        <span v-if="field.is_required" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
                          Required
                        </span>
                      </div>
                    </div>
                    <p :class="[
                      'text-xs mt-1',
                      field.is_enabled ? 'text-gray-500' : 'text-gray-400'
                    ]">
                      {{ field.description || field.name }}
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                      Type: {{ field.field_type }} | Order: {{ field.sort_order }}
                    </p>
                  </div>
                  <div class="ml-4 flex items-center space-x-2">
                    <!-- Enable/Disable Toggle -->
                    <button
                      v-if="!field.is_core"
                      @click="toggleFieldEnabled(field)"
                      :class="[
                        'relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500',
                        field.is_enabled ? 'bg-indigo-600' : 'bg-gray-200'
                      ]"
                    >
                      <span :class="[
                        'pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200',
                        field.is_enabled ? 'translate-x-5' : 'translate-x-0'
                      ]"></span>
                    </button>
                    
                    <!-- Required Toggle -->
                    <button
                      v-if="field.is_enabled && !field.is_core"
                      @click="toggleFieldRequired(field)"
                      :class="[
                        'px-2 py-1 text-xs font-medium rounded border',
                        field.is_required 
                          ? 'bg-red-100 text-red-800 border-red-200 hover:bg-red-200'
                          : 'bg-gray-100 text-gray-700 border-gray-200 hover:bg-gray-200'
                      ]"
                    >
                      {{ field.is_required ? 'Required' : 'Optional' }}
                    </button>

                    <!-- Move buttons -->
                    <div class="flex flex-col space-y-1">
                      <button
                        @click="moveFieldUp(field)"
                        :disabled="isFirstInCategory(field)"
                        class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-25"
                      >
                        <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                      <button
                        @click="moveFieldDown(field)"
                        :disabled="isLastInCategory(field)"
                        class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-25"
                      >
                        <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Preview Panel -->
          <div v-if="showPreview">
            <h4 class="text-md font-medium text-gray-900 mb-4">Form Preview</h4>
            <div class="border border-gray-200 rounded-lg p-4 bg-gray-50 max-h-96 overflow-y-auto">
              <div v-for="category in categories" :key="category" class="mb-6">
                <h5 class="text-sm font-medium text-gray-700 mb-3">{{ category }}</h5>
                <div class="space-y-3">
                  <div
                    v-for="field in getCategoryFields(category).filter(f => f.is_enabled)"
                    :key="field.field_id"
                    class="bg-white p-3 rounded border"
                  >
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                      {{ field.display_name }}
                      <span v-if="field.is_required" class="text-red-500">*</span>
                    </label>
                    <component
                      :is="getPreviewComponent(field)"
                      v-bind="getPreviewProps(field)"
                      disabled
                      class="w-full"
                    />
                    <p v-if="field.description" class="text-xs text-gray-500 mt-1">{{ field.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Reset Confirmation Modal -->
    <div v-if="showResetModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="showResetModal = false">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
            <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.956-.833-2.726 0L3.088 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <div class="mt-2 text-center">
            <h3 class="text-lg font-medium text-gray-900">Reset Form Configuration</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500">
                Are you sure you want to reset the {{ getFormTypeDisplayName(selectedFormType) }} form to default settings? 
                This will undo all your customizations and cannot be undone.
              </p>
            </div>
            <div class="items-center px-4 py-3 flex justify-center space-x-3">
              <button
                @click="showResetModal = false"
                class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                @click="resetForm"
                :disabled="isLoading"
                class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50"
              >
                {{ isLoading ? 'Resetting...' : 'Reset Form' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useFormConfigStore } from '@/stores/formConfig'
import {
  UserIcon,
  CalendarDaysIcon,
  DocumentTextIcon,
  ClipboardDocumentListIcon
} from '@heroicons/vue/24/outline'

// Composables
const toast = useToast()
const formConfigStore = useFormConfigStore()

// Local state
const selectedFormType = ref('')
const selectedCategory = ref('')
const showPreview = ref(false)
const showResetModal = ref(false)

// Computed properties
const formTypes = computed(() => formConfigStore.formTypes)
const fields = computed(() => formConfigStore.getFormFields(selectedFormType.value))
const fieldsByCategory = computed(() => formConfigStore.getFieldsByCategory(selectedFormType.value))
const isLoading = computed(() => formConfigStore.isLoading)
const error = computed(() => formConfigStore.error)
const hasChanges = computed(() => formConfigStore.isDirty)

const categories = computed(() => {
  return Object.keys(fieldsByCategory.value).sort()
})

// Methods
const getFormTypeIcon = (formType) => {
  const icons = {
    patient: UserIcon,
    appointment: CalendarDaysIcon,
    medical_record: DocumentTextIcon,
    default: ClipboardDocumentListIcon
  }
  return icons[formType] || icons.default
}

const getFormTypeDisplayName = (formType) => {
  const displayNames = {
    patient: 'Patient Registration',
    appointment: 'Appointment Booking',
    medical_record: 'Medical Record'
  }
  return displayNames[formType] || formType
}

const selectFormType = async (formType) => {
  if (selectedFormType.value === formType) return
  
  selectedFormType.value = formType
  selectedCategory.value = categories.value[0] || ''
  
  try {
    await formConfigStore.loadFormFields(formType)
  } catch (err) {
    toast.error(`Failed to load ${formType} form configuration`)
  }
}

const getCategoryFields = (category) => {
  return fieldsByCategory.value[category] || []
}

const getCategoryFieldCount = (category) => {
  const categoryFields = getCategoryFields(category)
  const enabled = categoryFields.filter(f => f.is_enabled).length
  const total = categoryFields.length
  return `${enabled}/${total}`
}

const toggleFieldEnabled = async (field) => {
  if (field.is_core) return
  
  try {
    await formConfigStore.toggleFieldEnabled(selectedFormType.value, field.field_id)
    toast.success(`${field.display_name} ${field.is_enabled ? 'disabled' : 'enabled'}`)
  } catch (err) {
    toast.error(`Failed to toggle ${field.display_name}`)
  }
}

const toggleFieldRequired = async (field) => {
  if (field.is_core || !field.is_enabled) return
  
  try {
    await formConfigStore.toggleFieldRequired(selectedFormType.value, field.field_id)
    toast.success(`${field.display_name} marked as ${field.is_required ? 'optional' : 'required'}`)
  } catch (err) {
    toast.error(`Failed to update ${field.display_name} requirement`)
  }
}

const moveFieldUp = async (field) => {
  const categoryFields = getCategoryFields(selectedCategory.value)
  const currentIndex = categoryFields.findIndex(f => f.field_id === field.field_id)
  
  if (currentIndex <= 0) return
  
  const previousField = categoryFields[currentIndex - 1]
  const fieldOrders = [
    { field_id: field.field_id, sort_order: previousField.sort_order },
    { field_id: previousField.field_id, sort_order: field.sort_order }
  ]
  
  try {
    await formConfigStore.updateFieldOrder(selectedFormType.value, fieldOrders)
    toast.success(`Moved ${field.display_name} up`)
  } catch (err) {
    toast.error(`Failed to move ${field.display_name}`)
  }
}

const moveFieldDown = async (field) => {
  const categoryFields = getCategoryFields(selectedCategory.value)
  const currentIndex = categoryFields.findIndex(f => f.field_id === field.field_id)
  
  if (currentIndex >= categoryFields.length - 1) return
  
  const nextField = categoryFields[currentIndex + 1]
  const fieldOrders = [
    { field_id: field.field_id, sort_order: nextField.sort_order },
    { field_id: nextField.field_id, sort_order: field.sort_order }
  ]
  
  try {
    await formConfigStore.updateFieldOrder(selectedFormType.value, fieldOrders)
    toast.success(`Moved ${field.display_name} down`)
  } catch (err) {
    toast.error(`Failed to move ${field.display_name}`)
  }
}

const isFirstInCategory = (field) => {
  const categoryFields = getCategoryFields(selectedCategory.value)
  return categoryFields[0]?.field_id === field.field_id
}

const isLastInCategory = (field) => {
  const categoryFields = getCategoryFields(selectedCategory.value)
  return categoryFields[categoryFields.length - 1]?.field_id === field.field_id
}

const getPreviewComponent = (field) => {
  switch (field.field_type) {
    case 'textarea':
      return 'textarea'
    case 'select':
    case 'multiselect':
      return 'select'
    default:
      return 'input'
  }
}

const getPreviewProps = (field) => {
  const baseProps = {
    placeholder: `Enter ${field.display_name.toLowerCase()}`,
    class: 'block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm'
  }

  switch (field.field_type) {
    case 'email':
      return { ...baseProps, type: 'email' }
    case 'tel':
      return { ...baseProps, type: 'tel' }
    case 'url':
      return { ...baseProps, type: 'url' }
    case 'number':
      return { ...baseProps, type: 'number' }
    case 'date':
      return { ...baseProps, type: 'date' }
    case 'textarea':
      return { ...baseProps, rows: 3 }
    case 'select':
    case 'multiselect':
      return {
        ...baseProps,
        multiple: field.field_type === 'multiselect'
      }
    default:
      return { ...baseProps, type: 'text' }
  }
}

const confirmResetForm = () => {
  showResetModal.value = true
}

const resetForm = async () => {
  if (!selectedFormType.value) return
  
  try {
    await formConfigStore.resetFormToDefaults(selectedFormType.value)
    toast.success(`${getFormTypeDisplayName(selectedFormType.value)} form reset to defaults`)
    showResetModal.value = false
  } catch (err) {
    toast.error(`Failed to reset ${getFormTypeDisplayName(selectedFormType.value)} form`)
  }
}

const saveAllChanges = async () => {
  // This would batch save any pending changes
  // For now, changes are saved immediately when toggled
  formConfigStore.markClean()
  toast.success('All changes have been saved')
}

const refreshConfiguration = async () => {
  if (!selectedFormType.value) return
  
  try {
    // Force reload from API, bypassing cache
    await formConfigStore.reloadFormFields(selectedFormType.value)
    toast.success('Configuration refreshed')
  } catch (err) {
    toast.error('Failed to refresh configuration')
  }
}

const loadFormConfiguration = () => {
  if (selectedFormType.value) {
    selectFormType(selectedFormType.value)
  }
}

// Watch for form type changes to update category selection
watch(categories, (newCategories) => {
  if (newCategories.length > 0 && !selectedCategory.value) {
    selectedCategory.value = newCategories[0]
  }
}, { immediate: true })

// Lifecycle
onMounted(async () => {
  try {
    await formConfigStore.initialize()
    
    // Auto-select first form type if available
    if (formTypes.value.length > 0) {
      await selectFormType(formTypes.value[0].name)
    }
  } catch (err) {
    toast.error('Failed to initialize form configuration')
  }
})
</script>