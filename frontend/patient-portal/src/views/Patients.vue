<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <PatientHeader @add-patient="showAddPatientModal = true" />

      <!-- Search -->
      <PatientSearch 
        v-model="searchQuery" 
        @search="debouncedSearch" 
      />

      <!-- Patients Table -->
      <PatientsTable
        :patients="patients"
        :loading="loading"
        :error="error"
        :total-count="pagination.totalCount"
        :offset="pagination.offset"
        :limit="pagination.limit"
        :has-more="pagination.hasMore"
        :format-gender="formatGender"
        :get-country-name="getCountryName"
        @retry="loadPatients"
        @view-patient="viewPatient"
        @edit-patient="editPatient"
        @delete-patient="deletePatient"
        @prev-page="prevPage"
        @next-page="nextPage"
      />
    </div>

    <!-- Add Patient Modal -->
    <div v-if="showAddPatientModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeAddPatientModal">
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-gray-900">Add New Patient</h3>
            <button @click="closeAddPatientModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Use DynamicPatientForm Component -->
          <DynamicPatientForm 
            key="add-patient"
            :initial-data="{}"
            :is-edit="false"
            @submit="handleFormSubmit"
            @cancel="closeAddPatientModal"
          />
        </div>
      </div>
    </div>

    <!-- Edit Patient Modal -->
    <div v-if="editingPatient && Object.keys(editingPatient).length > 0" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeEditModal">
      <div class="relative top-4 mx-auto p-0 w-11/12 md:w-4/5 lg:w-3/4 xl:w-2/3 max-w-6xl" @click.stop>
        <div class="bg-white rounded-lg shadow-xl">
          <div class="flex items-center justify-between p-6 border-b">
            <h3 class="text-xl font-semibold text-gray-900">Edit Patient</h3>
            <button @click="closeEditModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <!-- Dynamic Patient Form for Editing -->
          <div class="p-6">
            <DynamicPatientForm 
              :key="`edit-${editingPatient?.id}-${Date.now()}`"
              :initial-data="editingPatient"
              :is-edit="true"
              @submit="handleEditFormSubmit"
              @cancel="closeEditModal"
            />
          </div>
        </div>
      </div>
    </div>


    <!-- Patient View Modal -->
    <PatientViewModal
      :show="showViewModal"
      :patient="viewingPatient"
      :format-gender="formatGender"
      :get-country-name="getCountryName"
      @close="closeViewModal"
      @edit="handleEditFromView"
    />
  </div>
</template>

<script>
import { debounce } from 'lodash-es'
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toastification'
import { patientsApi } from '../api/patients'
import DynamicPatientForm from '../components/patients/DynamicPatientForm.vue'
import PatientHeader from '../components/patients/PatientHeader.vue'
import PatientSearch from '../components/patients/PatientSearch.vue'
import PatientsTable from '../components/patients/PatientsTable.vue'
import PatientViewModal from '../components/patients/PatientViewModal.vue'
import { useFormConfig } from '../composables/useFormConfig'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Patients',
  components: {
    DynamicPatientForm,
    PatientHeader,
    PatientSearch,
    PatientsTable,
    PatientViewModal
  },
  setup() {
    const { t } = useI18n()
    const toast = useToast()
    const authStore = useAuthStore()
    
    // Form configuration
    const patientFormConfig = useFormConfig('patient')

    // Reactive data
    const patients = ref([])
    const loading = ref(false)
    const error = ref('')
    const searchQuery = ref('')
    const showAddPatientModal = ref(false)
    const editingPatient = ref(null)
    const viewingPatient = ref(null)
    const showViewModal = ref(false)
    const addingPatient = ref(false)
    const addPatientError = ref('')


    // Generate year range (from 1920 to current year)
    const currentYear = new Date().getFullYear()
    const yearRange = computed(() => {
      const years = []
      for (let year = currentYear; year >= 1920; year--) {
        years.push(year)
      }
      return years
    })

    // Add Patient modal no longer needs form data - DynamicPatientForm handles it

    // Watch for changes in edit form birth date components and combine them

    const pagination = reactive({
      limit: 20,
      offset: 0,
      totalCount: 0,
      hasMore: false
    })

    // Initialize form data based on configuration - NO FALLBACKS
    // Retry form configuration loading  
    const retryFormConfiguration = async () => {
      try {
        await patientFormConfig.initialize()
        // DynamicPatientForm will handle its own initialization
      } catch (err) {
        console.error('Failed to retry form configuration:', err)
      }
    }

    // Methods
    const loadPatients = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const params = {
          limit: pagination.limit,
          offset: pagination.offset
        }
        
        if (searchQuery.value) {
          params.q = searchQuery.value
        }

        const response = await patientsApi.getPatients(params)
        patients.value = response.patients || []
        pagination.totalCount = response.total_count || 0
        pagination.hasMore = response.has_more || false
      } catch (err) {
        error.value = err.message || 'Error loading patients'
        console.error('Error loading patients:', err)
      } finally {
        loading.value = false
      }
    }

    const debouncedSearch = debounce(() => {
      pagination.offset = 0
      loadPatients()
    }, 300)

    const nextPage = () => {
      if (pagination.hasMore) {
        pagination.offset += pagination.limit
        loadPatients()
      }
    }

    const prevPage = () => {
      if (pagination.offset > 0) {
        pagination.offset = Math.max(0, pagination.offset - pagination.limit)
        loadPatients()
      }
    }

    const viewPatient = (patient) => {
      viewingPatient.value = patient
      showViewModal.value = true
    }

    const closeViewModal = () => {
      showViewModal.value = false
      viewingPatient.value = null
    }

    const editPatient = (patient) => {
      // Create a deep copy for editing
      editingPatient.value = {
        ...patient,
        date_of_birth: patient.date_of_birth ? patient.date_of_birth.split('T')[0] : '',
        emergency_contact: {
          name: patient.emergency_contact?.name || '',
          phone: patient.emergency_contact?.phone || '',
          relationship: patient.emergency_contact?.relationship || ''
        }
      }
      
      viewingPatient.value = null
    }

    const closeEditModal = () => {
      editingPatient.value = null
    }


    const handleEditFormSubmit = async (patientData) => {
      try {
        // The DynamicPatientForm component handles validation and API call
        closeEditModal()
        await loadPatients() // Refresh the patient list
        toast.success(t('patients.patientUpdated'))
      } catch (err) {
        console.error('Edit form submission error:', err)
        toast.error(err.message || t('errors.general'))
      }
    }

    const deletePatient = async (patient) => {
      if (!confirm(t('patients.confirmDelete'))) {
        return
      }

      try {
        await patientsApi.deletePatient(patient.id)
        toast.success(t('patients.patientDeleted'))
        loadPatients()
      } catch (err) {
        toast.error(err.message || t('errors.general'))
        console.error('Error deleting patient:', err)
      }
    }

    const formatGender = (gender) => {
      return gender.charAt(0).toUpperCase() + gender.slice(1)
    }

    const getCountryName = (countryCode) => {
      const countries = {
        'Canada': 'Canada',
        'USA': 'United States',
        'Morocco': 'Morocco',
        'France': 'France'
      }
      return countries[countryCode] || countryCode
    }

    // Add Patient Methods
    const closeAddPatientModal = () => {
      showAddPatientModal.value = false
      addPatientError.value = ''
      // DynamicPatientForm handles its own cleanup
    }

    // Handle successful form submission from DynamicPatientForm
    const handleFormSubmit = async (result) => {
      // DynamicPatientForm already created the patient and shows success toast
      // Just handle the UI updates
      closeAddPatientModal()
      loadPatients() // Refresh the patient list
    }

    const handleEditFromView = (patient) => {
      closeViewModal()
      editPatient(patient)
    }


    // Lifecycle
    onMounted(async () => {
      // DynamicPatientForm handles its own initialization
      // Just load the patient list
      loadPatients()
    })

    return {
      // Data
      patients,
      loading,
      error,
      searchQuery,
      pagination,
      showAddPatientModal,
      editingPatient,
      viewingPatient,
      showViewModal,
      addingPatient,
      addPatientError,
      yearRange,

      // Methods
      loadPatients,
      debouncedSearch,
      nextPage,
      prevPage,
      viewPatient,
      closeViewModal,
      editPatient,
      closeEditModal,
      handleEditFormSubmit,
      deletePatient,
      formatGender,
      getCountryName,
      closeAddPatientModal,
      handleFormSubmit,
      handleEditFromView,
      t,
      
      // Form configuration methods
      retryFormConfiguration,
      patientFormConfig
    }
  }
}
</script>