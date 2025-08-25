<template>
  <div class="bg-white shadow rounded-lg">
    <!-- Section Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-medium text-gray-900">
          {{ title }}
          <span v-if="showFieldCount" class="text-sm font-normal text-gray-500 ml-2">
            ({{ totalFields }} fields)
          </span>
        </h3>
        
        <!-- Virtualization Toggle (Development Only) -->
        <div v-if="isDevelopment && totalFields >= minVirtualFields" class="flex items-center space-x-2">
          <label class="text-sm text-gray-600">Virtual scrolling:</label>
          <button
            @click="virtualEnabled = !virtualEnabled"
            :class="virtualToggleClasses"
            type="button"
          >
            {{ virtualEnabled ? 'ON' : 'OFF' }}
          </button>
        </div>
      </div>
      
      <!-- Performance Info (Development Only) -->
      <div v-if="isDevelopment && virtualEnabled" class="mt-2 text-xs text-gray-500">
        Rendering {{ visibleFieldCount }}/{{ totalFields }} fields
        <span v-if="scrollPosition > 0">(scrolled {{ Math.round(scrollPosition) }}px)</span>
      </div>
    </div>

    <!-- Virtual Scrolled Fields -->
    <div v-if="shouldUseVirtual" class="p-6">
      <VirtualScrollForm
        :items="formFields"
        :item-height="getFieldHeight"
        :container-height="virtualHeight"
        :overscan="overscan"
        :get-item-key="getFieldKey"
        :loading="loading"
        @visible-range="handleVisibleRange"
        @scroll="handleScroll"
        ref="virtualScroll"
      >
        <template #default="{ item: field, index }">
          <div class="virtual-field-wrapper mb-4">
            <FormField
              :key="getFieldKey(field, index)"
              :field="field"
              :model-value="getFieldValue(field.name)"
              :options="getFieldOptions(field)"
              :errors="getFieldError(field.name)"
              :loading="fieldLoading"
              :disabled="fieldDisabled"
              @update:model-value="(value) => handleFieldUpdate(field.name, value)"
              @input="(value) => handleFieldInput(field.name, value)"
              @change="(value) => handleFieldChange(field.name, value)"
            />
          </div>
        </template>
        
        <template #empty>
          <div class="text-center py-8 text-gray-500">
            <p>No fields configured for this section</p>
          </div>
        </template>
      </VirtualScrollForm>
    </div>

    <!-- Regular Non-Virtual Fields -->
    <div v-else class="p-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField
          v-for="field in formFields"
          :key="field.field_id || field.name"
          :field="field"
          :model-value="getFieldValue(field.name)"
          :options="getFieldOptions(field)"
          :errors="getFieldError(field.name)"
          :loading="fieldLoading"
          :disabled="fieldDisabled"
          @update:model-value="(value) => handleFieldUpdate(field.name, value)"
          @input="(value) => handleFieldInput(field.name, value)"
          @change="(value) => handleFieldChange(field.name, value)"
        />
      </div>
      
      <!-- Empty State -->
      <div v-if="formFields.length === 0" class="text-center py-8 text-gray-500">
        <p>No fields configured for this section</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import VirtualScrollForm from './VirtualScrollForm.vue'
import FormField from './FormField.vue'

// Props
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  fields: {
    type: Array,
    required: true,
    default: () => []
  },
  formData: {
    type: Object,
    default: () => ({})
  },
  getFieldOptions: {
    type: Function,
    default: () => []
  },
  getFieldError: {
    type: Function,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  // Virtual scrolling props
  enableVirtual: {
    type: Boolean,
    default: null // null = auto-detect
  },
  minVirtualFields: {
    type: Number,
    default: 20,
    validator: (value) => value > 0
  },
  virtualHeight: {
    type: String,
    default: '500px'
  },
  overscan: {
    type: Number,
    default: 5
  },
  showFieldCount: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits({
  'field-update': (fieldName, value) => typeof fieldName === 'string',
  'field-input': (fieldName, value) => typeof fieldName === 'string',
  'field-change': (fieldName, value) => typeof fieldName === 'string'
})

// Refs
const virtualScroll = ref(null)

// State
const virtualEnabled = ref(false)
const visibleFieldCount = ref(0)
const scrollPosition = ref(0)

// Computed
const isDevelopment = computed(() => import.meta.env.DEV)

const totalFields = computed(() => props.fields.length)

const formFields = computed(() => {
  // Sort fields by sort_order if available
  return [...props.fields].sort((a, b) => {
    const aOrder = a.sort_order ?? 0
    const bOrder = b.sort_order ?? 0
    return aOrder - bOrder
  })
})

const shouldUseVirtual = computed(() => {
  if (props.enableVirtual === false) return false
  if (props.enableVirtual === true) return true
  
  // Auto-detect based on field count and virtual enabled state
  return virtualEnabled.value && totalFields.value >= props.minVirtualFields
})

const fieldLoading = computed(() => props.loading)

const fieldDisabled = computed(() => props.disabled)

const virtualToggleClasses = computed(() => {
  const base = 'px-2 py-1 text-xs rounded font-medium transition-colors'
  if (virtualEnabled.value) {
    return `${base} bg-green-100 text-green-700 hover:bg-green-200`
  } else {
    return `${base} bg-gray-100 text-gray-700 hover:bg-gray-200`
  }
})

// Methods
const getFieldValue = (fieldName) => {
  return props.formData[fieldName] ?? ''
}

const getFieldHeight = (field, index) => {
  // Estimate field height based on field type
  const fieldType = field.field_type
  
  if (fieldType === 'textarea') {
    return 120 // Taller for textarea
  } else if (['select', 'multiselect'].includes(fieldType)) {
    return 90 // Slightly taller for selects
  } else if (fieldType === 'file') {
    return 140 // Tallest for file inputs
  } else if (fieldType === 'checkbox' || fieldType === 'boolean') {
    return 70 // Shorter for checkboxes
  }
  
  return 80 // Default height
}

const getFieldKey = (field, index) => {
  return field.field_id || field.name || `field-${index}`
}

const handleFieldUpdate = (fieldName, value) => {
  emit('field-update', fieldName, value)
}

const handleFieldInput = (fieldName, value) => {
  emit('field-input', fieldName, value)
}

const handleFieldChange = (fieldName, value) => {
  emit('field-change', fieldName, value)
}

const handleVisibleRange = (startIndex, endIndex) => {
  visibleFieldCount.value = endIndex - startIndex + 1
}

const handleScroll = (scrollTop, scrollDirection) => {
  scrollPosition.value = scrollTop
}

// Public methods
const scrollToField = (fieldName) => {
  if (!shouldUseVirtual.value || !virtualScroll.value) return
  
  const fieldIndex = formFields.value.findIndex(f => f.name === fieldName)
  if (fieldIndex >= 0) {
    virtualScroll.value.scrollToIndex(fieldIndex, 'center')
  }
}

const scrollToTop = () => {
  if (virtualScroll.value) {
    virtualScroll.value.scrollToTop()
  }
}

// Auto-enable virtual scrolling for large sections
watch(totalFields, (newCount) => {
  if (props.enableVirtual === null) {
    virtualEnabled.value = newCount >= props.minVirtualFields
  }
}, { immediate: true })

// Initialize virtual scrolling state
onMounted(() => {
  if (props.enableVirtual !== null) {
    virtualEnabled.value = props.enableVirtual
  } else {
    virtualEnabled.value = totalFields.value >= props.minVirtualFields
  }
})

// Expose methods for parent components
defineExpose({
  scrollToField,
  scrollToTop,
  virtualScroll: computed(() => virtualScroll.value)
})
</script>

<style scoped>
.virtual-field-wrapper {
  contain: layout style paint;
}
</style>