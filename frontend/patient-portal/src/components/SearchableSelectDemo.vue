<template>
  <div class="max-w-4xl mx-auto p-6 space-y-8">
    <h1 class="text-2xl font-bold text-gray-900">SearchableSelect Component Demo</h1>
    <p class="text-gray-600">This demo shows how the SearchableSelect component can be used with different data types.</p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <!-- Patients Example -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-gray-800">Patients Selector</h2>
        <SearchableSelect
          v-model="selectedPatient"
          label="Select Patient"
          placeholder="Search patients by name, ID, phone..."
          :items="samplePatients"
          item-type="patients"
          :display-key="patient => `${patient.first_name} ${patient.last_name}`"
          :secondary-key="patient => `ID: ${patient.patient_id} • ${patient.phone}`"
          :search-keys="['first_name', 'last_name', 'patient_id', 'phone']"
          item-icon="UserIcon"
          @select="onPatientSelect"
        />
        <div v-if="selectedPatient" class="text-sm text-gray-600">
          Selected Patient ID: {{ selectedPatient }}
        </div>
      </div>

      <!-- Doctors Example -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-gray-800">Doctors Selector</h2>
        <SearchableSelect
          v-model="selectedDoctor"
          label="Select Doctor"
          placeholder="Search doctors by name or specialization..."
          :items="sampleDoctors"
          item-type="doctors"
          :display-key="doctor => `Dr. ${doctor.first_name} ${doctor.last_name}`"
          :secondary-key="doctor => `${doctor.specialization} • ${doctor.department}`"
          :search-keys="['first_name', 'last_name', 'specialization', 'department']"
          item-icon="UserIcon"
          icon-class="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center"
          icon-size="h-4 w-4 text-green-600"
          @select="onDoctorSelect"
        />
        <div v-if="selectedDoctor" class="text-sm text-gray-600">
          Selected Doctor ID: {{ selectedDoctor }}
        </div>
      </div>

      <!-- Medications Example -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-gray-800">Medications Selector</h2>
        <SearchableSelect
          v-model="selectedMedication"
          label="Select Medication"
          placeholder="Search medications by name or type..."
          :items="sampleMedications"
          item-type="medications"
          display-key="name"
          :secondary-key="med => `${med.dosage} • ${med.manufacturer}`"
          :search-keys="['name', 'generic_name', 'manufacturer']"
          item-icon="BeakerIcon"
          icon-class="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center"
          icon-size="h-4 w-4 text-blue-600"
          @select="onMedicationSelect"
        />
        <div v-if="selectedMedication" class="text-sm text-gray-600">
          Selected Medication ID: {{ selectedMedication }}
        </div>
      </div>

      <!-- Departments Example -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-gray-800">Departments Selector</h2>
        <SearchableSelect
          v-model="selectedDepartment"
          label="Select Department"
          placeholder="Search departments..."
          :items="sampleDepartments"
          item-type="departments"
          display-key="name"
          :secondary-key="dept => `Floor ${dept.floor} • ${dept.head_doctor}`"
          :search-keys="['name', 'head_doctor', 'description']"
          item-icon="BuildingOfficeIcon"
          icon-class="h-8 w-8 rounded-full bg-purple-100 flex items-center justify-center"
          icon-size="h-4 w-4 text-purple-600"
          @select="onDepartmentSelect"
        />
        <div v-if="selectedDepartment" class="text-sm text-gray-600">
          Selected Department ID: {{ selectedDepartment }}
        </div>
      </div>
    </div>

    <!-- Custom Template Example -->
    <div class="space-y-4">
      <h2 class="text-lg font-semibold text-gray-800">Custom Template Example</h2>
      <p class="text-sm text-gray-600">This example shows how to use custom slot templates for complex layouts.</p>
      
      <SearchableSelect
        v-model="selectedPatientCustom"
        label="Select Patient (Custom Template)"
        placeholder="Search patients..."
        :items="samplePatients"
        item-type="patients"
        :display-key="patient => `${patient.first_name} ${patient.last_name}`"
        :search-keys="['first_name', 'last_name', 'patient_id', 'phone']"
        :show-icon="false"
        @select="onPatientCustomSelect"
      >
        <template #option="{ item: patient, selected }">
          <div class="flex items-center w-full p-2">
            <!-- Custom Avatar -->
            <div class="flex-shrink-0">
              <div class="h-10 w-10 rounded-full bg-gradient-to-r from-indigo-400 to-purple-500 flex items-center justify-center">
                <span class="text-white font-semibold text-sm">
                  {{ patient.first_name[0] }}{{ patient.last_name[0] }}
                </span>
              </div>
            </div>

            <!-- Patient Info -->
            <div class="ml-4 flex-1">
              <div class="flex items-center justify-between">
                <div class="text-sm font-medium text-gray-900">
                  {{ patient.first_name }} {{ patient.last_name }}
                </div>
                <div v-if="selected" class="text-indigo-600">
                  <CheckIcon class="h-4 w-4" />
                </div>
              </div>
              <div class="text-xs text-gray-500 mt-1">
                <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800 mr-2">
                  ID: {{ patient.patient_id }}
                </span>
                <span class="text-gray-400">{{ patient.phone }}</span>
              </div>
            </div>
          </div>
        </template>
      </SearchableSelect>
      <div v-if="selectedPatientCustom" class="text-sm text-gray-600">
        Selected Patient (Custom) ID: {{ selectedPatientCustom }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { CheckIcon } from '@heroicons/vue/24/outline'
import SearchableSelect from '@/components/SearchableSelect.vue'

// Reactive state
const selectedPatient = ref(null)
const selectedDoctor = ref(null)
const selectedMedication = ref(null)
const selectedDepartment = ref(null)
const selectedPatientCustom = ref(null)

// Sample data
const samplePatients = ref([
  {
    id: 1,
    patient_id: 'PAT001',
    first_name: 'John',
    last_name: 'Smith',
    phone: '(555) 123-4567',
    email: 'john.smith@email.com',
    date_of_birth: '1985-06-15'
  },
  {
    id: 2,
    patient_id: 'PAT002',
    first_name: 'Sarah',
    last_name: 'Johnson',
    phone: '(555) 234-5678',
    email: 'sarah.johnson@email.com',
    date_of_birth: '1990-03-22'
  },
  {
    id: 3,
    patient_id: 'PAT003',
    first_name: 'Michael',
    last_name: 'Brown',
    phone: '(555) 345-6789',
    email: 'michael.brown@email.com',
    date_of_birth: '1978-11-08'
  },
  {
    id: 4,
    patient_id: 'PAT004',
    first_name: 'Emma',
    last_name: 'Davis',
    phone: '(555) 456-7890',
    email: 'emma.davis@email.com',
    date_of_birth: '1995-09-12'
  }
])

const sampleDoctors = ref([
  {
    id: 1,
    first_name: 'Robert',
    last_name: 'Wilson',
    specialization: 'Cardiology',
    department: 'Heart Center',
    license: 'MD123456'
  },
  {
    id: 2,
    first_name: 'Lisa',
    last_name: 'Anderson',
    specialization: 'Pediatrics',
    department: 'Children\'s Wing',
    license: 'MD234567'
  },
  {
    id: 3,
    first_name: 'David',
    last_name: 'Martinez',
    specialization: 'Orthopedics',
    department: 'Bone & Joint',
    license: 'MD345678'
  }
])

const sampleMedications = ref([
  {
    id: 1,
    name: 'Lisinopril',
    generic_name: 'Lisinopril',
    dosage: '10mg tablets',
    manufacturer: 'Generic Pharma'
  },
  {
    id: 2,
    name: 'Metformin',
    generic_name: 'Metformin HCl',
    dosage: '500mg tablets',
    manufacturer: 'Diabetes Care Co.'
  },
  {
    id: 3,
    name: 'Amoxicillin',
    generic_name: 'Amoxicillin',
    dosage: '250mg capsules',
    manufacturer: 'Antibiotic Labs'
  }
])

const sampleDepartments = ref([
  {
    id: 1,
    name: 'Emergency Department',
    floor: 1,
    head_doctor: 'Dr. Sarah Chen',
    description: '24/7 emergency care'
  },
  {
    id: 2,
    name: 'Intensive Care Unit',
    floor: 3,
    head_doctor: 'Dr. Michael Torres',
    description: 'Critical care services'
  },
  {
    id: 3,
    name: 'Maternity Ward',
    floor: 2,
    head_doctor: 'Dr. Jennifer Lee',
    description: 'Obstetrics and newborn care'
  }
])

// Event handlers
const onPatientSelect = (patient) => {
}

const onDoctorSelect = (doctor) => {
}

const onMedicationSelect = (medication) => {
}

const onDepartmentSelect = (department) => {
}

const onPatientCustomSelect = (patient) => {
}
</script>