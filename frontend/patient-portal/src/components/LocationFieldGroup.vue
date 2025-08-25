<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- Country field -->
    <div v-if="countryField" :class="getFieldGridClass(countryField)">
      <label :for="countryField.name" class="block text-sm font-medium text-gray-700 mb-1">
        {{ countryField.display_name }}
        <span v-if="countryField.is_required" class="text-red-500">*</span>
      </label>
      
      <select
        :id="countryField.name"
        :name="countryField.name"
        v-model.number="formData[countryField.name]"
        :class="getFieldClasses(countryField)"
        @change="() => handleCountryChange(countryField.name, formData[countryField.name])"
      >
        <option value="">Select Country</option>
        <option v-for="c in countries" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
      
      <!-- Debug info -->
      <div class="text-xs text-gray-400 mt-1">
        Debug: {{ countries.length }} countries loaded
      </div>
      
      <p v-if="countryField.description" class="text-xs text-gray-500 mt-1">
        {{ countryField.description }}
      </p>
      
      <div v-if="hasFieldError(countryField.name)" class="mt-1">
        <p
          v-for="error in getFieldError(countryField.name)"
          :key="error"
          class="text-sm text-red-600"
        >
          {{ error }}
        </p>
      </div>
    </div>

    <!-- State field -->
    <div v-if="stateField" :class="getFieldGridClass(stateField)">
      <label :for="stateField.name" class="block text-sm font-medium text-gray-700 mb-1">
        {{ stateField.display_name }}
        <span v-if="stateField.is_required" class="text-red-500">*</span>
      </label>
      
      <select
        :id="stateField.name"
        :name="stateField.name"
        v-model.number="formData[stateField.name]"
        :disabled="!selectedCountryId || statesLoading"
        :class="getFieldClasses(stateField)"
        @change="() => handleStateChange(stateField.name, formData[stateField.name])"
      >
        <option value="" :disabled="true">
          {{ !selectedCountryId ? 'Select country first' : (statesLoading ? 'Loading states…' : 'Select State/Province') }}
        </option>
        <option v-for="opt in stateOptions" :key="opt.id" :value="opt.id">{{ opt.name }}</option>
      </select>
      
      <p v-if="stateField.description" class="text-xs text-gray-500 mt-1">
        {{ stateField.description }}
      </p>
      
      <div v-if="hasFieldError(stateField.name)" class="mt-1">
        <p
          v-for="error in getFieldError(stateField.name)"
          :key="error"
          class="text-sm text-red-600"
        >
          {{ error }}
        </p>
      </div>
    </div>

    <!-- City field -->
    <div v-if="cityField" :class="getFieldGridClass(cityField)">
      <label :for="cityField.name" class="block text-sm font-medium text-gray-700 mb-1">
        {{ cityField.display_name }}
        <span v-if="cityField.is_required" class="text-red-500">*</span>
      </label>
      
      <select
        :id="cityField.name"
        :name="cityField.name"
        v-model.number="formData[cityField.name]"
        :disabled="!selectedCountryId || citiesLoading"
        :class="getFieldClasses(cityField)"
        @change="() => handleCityChange(cityField.name, formData[cityField.name])"
      >
        <option value="" :disabled="true">
          {{ !selectedCountryId ? 'Select country first' : (citiesLoading ? 'Loading cities…' : 'Select City') }}
        </option>
        <option v-for="opt in cityOptions" :key="opt.id" :value="opt.id">{{ opt.name }}</option>
      </select>
      
      <p v-if="cityField.description" class="text-xs text-gray-500 mt-1">
        {{ cityField.description }}
      </p>
      
      <div v-if="hasFieldError(cityField.name)" class="mt-1">
        <p
          v-for="error in getFieldError(cityField.name)"
          :key="error"
          class="text-sm text-red-600"
        >
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props with comprehensive validation
const props = defineProps({
  fields: {
    type: Array,
    required: true,
    validator: (value) => {
      // Fields should be an array of field objects
      return Array.isArray(value) && value.every(field => 
        field && typeof field === 'object' && 'name' in field && 'field_type' in field
      )
    }
  },
  formData: {
    type: Object,
    required: true,
    validator: (value) => {
      // FormData should be a non-null object
      return value !== null && typeof value === 'object'
    }
  },
  countries: {
    type: Array,
    default: () => [],
    validator: (value) => {
      // Countries should have id and name properties
      return value.every(country => 
        country && typeof country === 'object' && 'id' in country && 'name' in country
      )
    }
  },
  stateOptions: {
    type: Array,
    default: () => [],
    validator: (value) => {
      // States should have id and name properties
      return value.every(state => 
        state && typeof state === 'object' && 'id' in state && 'name' in state
      )
    }
  },
  cityOptions: {
    type: Array,
    default: () => [],
    validator: (value) => {
      // Cities should have id and name properties
      return value.every(city => 
        city && typeof city === 'object' && 'id' in city && 'name' in city
      )
    }
  },
  statesLoading: {
    type: Boolean,
    default: false
  },
  citiesLoading: {
    type: Boolean,
    default: false
  },
  resolveCountryId: {
    type: Function,
    required: true,
    validator: (value) => {
      // Should be a function
      return typeof value === 'function'
    }
  },
  hasFieldError: {
    type: Function,
    required: true,
    validator: (value) => {
      return typeof value === 'function'
    }
  },
  getFieldError: {
    type: Function,
    required: true,
    validator: (value) => {
      return typeof value === 'function'
    }
  }
})

// Emits with validation
const emit = defineEmits({
  'country-change': (fieldName, value) => {
    // Should include field name and value
    return typeof fieldName === 'string' && value !== undefined
  },
  'state-change': (fieldName, value) => {
    // Should include field name and value
    return typeof fieldName === 'string' && value !== undefined
  },
  'city-change': (fieldName, value) => {
    // Should include field name and value
    return typeof fieldName === 'string' && value !== undefined
  }
})

// Computed properties
const countryField = computed(() => {
  return props.fields.find(field => isCountryField(field))
})

const stateField = computed(() => {
  return props.fields.find(field => isStateField(field))
})

const cityField = computed(() => {
  return props.fields.find(field => isCityField(field))
})

const selectedCountryId = computed(() => {
  const raw = props.formData?.country_id || props.formData?.country || null
  return props.resolveCountryId(raw)
})

// Helper functions for field recognition - ID-based fields
const isCountryField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'country_id' || name === 'country' || name.endsWith('_country_id') || name.endsWith('_country')
}

const isStateField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'state_id' || name === 'state' || name === 'province_id' || name === 'province' || 
         name.endsWith('_state_id') || name.endsWith('_state') || name.endsWith('_province_id') || name.endsWith('_province') || 
         name.startsWith('state_') || name.startsWith('province_')
}

const isCityField = (field) => {
  const name = (field?.name || '').toLowerCase()
  return name === 'city_id' || name === 'city' || name.endsWith('_city_id') || name.endsWith('_city') || name.startsWith('city_')
}

// Styling functions
const getFieldGridClass = (field) => {
  if (field.field_type === 'textarea' ||
      ['address', 'medicalHistory', 'allergies', 'medications'].includes(field.name)) {
    return 'md:col-span-2'
  }
  return ''
}

const getFieldClasses = (field) => {
  const base = 'mt-1 block w-full rounded-md border shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
  const border = props.hasFieldError(field.name) ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
  return `${base} ${border}`
}

// Event handlers
const handleCountryChange = (fieldName, value) => {
  emit('country-change', fieldName, value)
}

const handleStateChange = (fieldName, value) => {
  emit('state-change', fieldName, value)
}

const handleCityChange = (fieldName, value) => {
  emit('city-change', fieldName, value)
}
</script>