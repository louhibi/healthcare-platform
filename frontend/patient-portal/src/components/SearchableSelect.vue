<template>
  <div class="relative">
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <!-- Search Input with Dropdown Toggle -->
    <div class="relative">
      <input
        :id="inputId"
        ref="searchInput"
        v-model="searchQuery"
        type="text"
        :placeholder="placeholder"
        :required="required"
        :disabled="disabled"
        @focus="handleFocus"
        @blur="handleBlur"
        @input="handleInput"
        @keydown="handleKeydown"
        :class="[
          'block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
          'pr-10',
          disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-white',
          error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : 'border-gray-300'
        ]"
      />
      
      <!-- Dropdown Toggle Icon -->
      <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
        <ChevronUpDownIcon class="h-5 w-5 text-gray-400" />
      </div>
    </div>

    <!-- Error Message -->
    <p v-if="error" class="mt-1 text-sm text-red-600">{{ error }}</p>

    <!-- Loading State -->
    <div v-if="loading" class="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-md shadow-lg z-50 p-3">
      <div class="flex items-center justify-center">
        <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600"></div>
        <span class="ml-2 text-sm text-gray-600">{{ loadingText }}</span>
      </div>
    </div>

    <!-- Dropdown Options -->
    <div 
      v-show="showDropdown && !loading"
      class="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-md shadow-lg z-50 max-h-60 overflow-y-auto"
    >
      <!-- No Results -->
      <div v-if="filteredItems.length === 0" class="px-4 py-3 text-sm text-gray-500 text-center">
        <component :is="noResultsIcon" class="mx-auto h-6 w-6 text-gray-300 mb-1" />
        <div v-if="searchQuery.length > 0">
          {{ noResultsText || `No ${itemType} found for "${searchQuery}"` }}
        </div>
        <div v-else>
          {{ emptyText || `No ${itemType} available` }}
        </div>
      </div>

      <!-- Items List -->
      <div
        v-for="(item, index) in filteredItems"
        :key="getItemKey(item)"
        @mousedown.prevent="selectItem(item)"
        :class="[
          'flex items-center px-4 py-3 cursor-pointer',
          index === highlightedIndex ? 'bg-indigo-50 text-indigo-900' : 'text-gray-900 hover:bg-gray-50'
        ]"
      >
        <!-- Custom Item Template -->
        <slot name="option" :item="item" :selected="isSelected(item)">
          <!-- Default Item Display -->
          <div class="flex items-center w-full">
            <!-- Icon -->
            <div v-if="showIcon" class="flex-shrink-0">
              <div :class="iconClass">
                <component :is="itemIcon" :class="iconSize" />
              </div>
            </div>

            <!-- Item Content -->
            <div class="flex-1 min-w-0" :class="showIcon ? 'ml-3' : ''">
              <div class="text-sm font-medium truncate">
                {{ getDisplayText(item) }}
              </div>
              <div v-if="getSecondaryText(item)" class="text-xs text-gray-500 truncate">
                {{ getSecondaryText(item) }}
              </div>
            </div>

            <!-- Selected Indicator -->
            <div v-if="isSelected(item)" class="flex-shrink-0">
              <CheckIcon class="h-4 w-4 text-indigo-600" />
            </div>
          </div>
        </slot>
      </div>

      <!-- Load More Button -->
      <div v-if="hasMore && filteredItems.length > 0" class="border-t border-gray-200">
        <button
          @mousedown.prevent="loadMore"
          class="w-full px-4 py-2 text-sm text-indigo-600 hover:bg-indigo-50 font-medium"
          :disabled="loadingMore"
        >
          <span v-if="loadingMore">Loading more...</span>
          <span v-else>Load more {{ itemType }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { 
  ChevronUpDownIcon, 
  CheckIcon,
  UserIcon,
  DocumentIcon,
  BeakerIcon,
  BuildingOfficeIcon
} from '@heroicons/vue/24/outline'

// Props
const props = defineProps({
  // v-model value
  modelValue: {
    type: [Number, String, Object],
    default: null
  },
  
  // Basic configuration
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Search and select...'
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: null
  },
  
  // Data configuration
  items: {
    type: Array,
    default: () => []
  },
  itemType: {
    type: String,
    default: 'items'
  },
  
  // Display configuration
  displayKey: {
    type: [String, Function],
    default: 'name'
  },
  secondaryKey: {
    type: [String, Function],
    default: null
  },
  valueKey: {
    type: String,
    default: 'id'
  },
  searchKeys: {
    type: Array,
    default: () => ['name']
  },
  
  // Icon configuration
  showIcon: {
    type: Boolean,
    default: true
  },
  itemIcon: {
    type: [String, Object],
    default: 'UserIcon'
  },
  iconClass: {
    type: String,
    default: 'h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center'
  },
  iconSize: {
    type: String,
    default: 'h-4 w-4 text-indigo-600'
  },
  noResultsIcon: {
    type: [String, Object],
    default: 'DocumentIcon'
  },
  
  // Loading and async configuration
  loading: {
    type: Boolean,
    default: false
  },
  loadingText: {
    type: String,
    default: 'Searching...'
  },
  loadingMore: {
    type: Boolean,
    default: false
  },
  hasMore: {
    type: Boolean,
    default: false
  },
  
  // Behavior configuration
  searchable: {
    type: Boolean,
    default: true
  },
  clearable: {
    type: Boolean,
    default: true
  },
  closeOnSelect: {
    type: Boolean,
    default: true
  },
  minSearchLength: {
    type: Number,
    default: 0
  },
  
  // Text customization
  noResultsText: {
    type: String,
    default: null
  },
  emptyText: {
    type: String,
    default: null
  }
})

// Emits
const emit = defineEmits([
  'update:modelValue', 
  'select', 
  'search', 
  'focus', 
  'blur',
  'load-more'
])

// Reactive state
const searchInput = ref(null)
const searchQuery = ref('')
const selectedItem = ref(null)
const showDropdown = ref(false)
const highlightedIndex = ref(-1)

// Computed
const inputId = computed(() => `searchable-select-${Math.random().toString(36).substr(2, 9)}`)

const iconComponents = {
  UserIcon,
  DocumentIcon,
  BeakerIcon,
  BuildingOfficeIcon
}

const filteredItems = computed(() => {
  if (!props.searchable || !searchQuery.value.trim()) {
    return props.items
  }

  if (searchQuery.value.trim().length < props.minSearchLength) {
    return props.items
  }

  const query = searchQuery.value.toLowerCase().trim()
  
  return props.items.filter(item => {
    return props.searchKeys.some(key => {
      const value = getNestedValue(item, key)
      return value && value.toString().toLowerCase().includes(query)
    })
  })
})

// Methods
const getNestedValue = (obj, path) => {
  return path.split('.').reduce((current, key) => current?.[key], obj)
}

const getItemKey = (item) => {
  return getNestedValue(item, props.valueKey) || item
}

const getDisplayText = (item) => {
  if (typeof props.displayKey === 'function') {
    return props.displayKey(item)
  }
  return getNestedValue(item, props.displayKey) || item
}

const getSecondaryText = (item) => {
  if (!props.secondaryKey) return null
  
  if (typeof props.secondaryKey === 'function') {
    return props.secondaryKey(item)
  }
  return getNestedValue(item, props.secondaryKey)
}

const isSelected = (item) => {
  if (!selectedItem.value) return false
  return getItemKey(item) === getItemKey(selectedItem.value)
}

const handleFocus = () => {
  if (!props.disabled) {
    showDropdown.value = true
    highlightedIndex.value = -1
    emit('focus')
  }
}

const handleBlur = () => {
  setTimeout(() => {
    showDropdown.value = false
    emit('blur')
  }, 150)
}

const handleInput = () => {
  highlightedIndex.value = -1
  
  // Clear selection if search doesn't match
  if (selectedItem.value && searchQuery.value.trim()) {
    const selectedText = getDisplayText(selectedItem.value)
    if (!selectedText.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      updateSelection(null)
    }
  }
  
  emit('search', searchQuery.value)
}

const handleKeydown = (event) => {
  if (!showDropdown.value) return

  switch (event.key) {
    case 'ArrowDown':
      event.preventDefault()
      highlightedIndex.value = Math.min(highlightedIndex.value + 1, filteredItems.value.length - 1)
      break
    
    case 'ArrowUp':
      event.preventDefault()
      highlightedIndex.value = Math.max(highlightedIndex.value - 1, -1)
      break
    
    case 'Enter':
      event.preventDefault()
      if (highlightedIndex.value >= 0 && filteredItems.value[highlightedIndex.value]) {
        selectItem(filteredItems.value[highlightedIndex.value])
      }
      break
    
    case 'Escape':
      event.preventDefault()
      showDropdown.value = false
      searchInput.value.blur()
      break
  }
}

const selectItem = (item) => {
  selectedItem.value = item
  searchQuery.value = getDisplayText(item)
  
  if (props.closeOnSelect) {
    showDropdown.value = false
    highlightedIndex.value = -1
    
    nextTick(() => {
      searchInput.value.blur()
    })
  }
  
  updateSelection(getItemKey(item))
  emit('select', item)
}

const updateSelection = (value) => {
  emit('update:modelValue', value)
}

const loadMore = () => {
  emit('load-more')
}

const findItemByValue = (value) => {
  return props.items.find(item => getItemKey(item) === value)
}

// Watch for external value changes
watch(() => props.modelValue, (newValue) => {
  if (newValue && newValue !== getItemKey(selectedItem.value)) {
    const item = findItemByValue(newValue)
    if (item) {
      selectedItem.value = item
      searchQuery.value = getDisplayText(item)
    }
  } else if (!newValue) {
    selectedItem.value = null
    searchQuery.value = ''
  }
}, { immediate: true })

// Expose methods for parent component
defineExpose({
  focus: () => searchInput.value?.focus(),
  clear: () => {
    selectedItem.value = null
    searchQuery.value = ''
    updateSelection(null)
  },
  select: selectItem
})
</script>

<style scoped>
/* Custom scrollbar for dropdown */
.overflow-y-auto::-webkit-scrollbar {
  width: 4px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>