<template>
  <div class="bg-white shadow rounded-lg p-6 mb-6">
    <div class="max-w-md">
      <label for="search" class="block text-sm font-medium text-gray-700">
        {{ $t('patients.searchLabel') }}
      </label>
      <input
        v-model="searchQuery"
        type="text"
        id="search"
        :placeholder="$t('patients.searchPlaceholder')"
        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        @input="handleSearch"
      />
      <p class="mt-1 text-xs text-gray-500">{{ $t('patients.searchHelp') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => typeof value === 'string',
  'search': (query) => typeof query === 'string'
})

// Local state
const searchQuery = ref(props.modelValue)

// Methods
const handleSearch = () => {
  emit('update:modelValue', searchQuery.value)
  emit('search', searchQuery.value)
}

// Watch for external changes
watch(() => props.modelValue, (newVal) => {
  searchQuery.value = newVal
})
</script>