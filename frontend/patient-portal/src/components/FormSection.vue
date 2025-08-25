<template>
  <div class="bg-white shadow rounded-lg p-6">
    <h3 class="text-lg font-medium text-gray-900 mb-4 border-b border-gray-200 pb-2">
      {{ title }}
    </h3>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <slot />
    </div>
  </div>
</template>

<script setup>
// Props with validation
defineProps({
  title: {
    type: String,
    required: true,
    validator: (value) => {
      // Title should be a non-empty string
      return typeof value === 'string' && value.trim().length > 0
    }
  }
})
</script>