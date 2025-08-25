<template>
  <div class="flex space-x-4">
    <button
      @click="$emit('view', patient)"
      class="text-blue-600 hover:text-blue-900 transition-colors duration-150"
      :title="$t('common.view')"
    >
      {{ $t('common.view') }}
    </button>
    <button
      @click="$emit('edit', patient)"
      class="text-green-600 hover:text-green-900 transition-colors duration-150"
      :title="$t('common.edit')"
    >
      {{ $t('common.edit') }}
    </button>
    <button
      @click="$emit('delete', patient)"
      class="text-red-600 hover:text-red-900 transition-colors duration-150"
      :title="$t('common.delete')"
    >
      {{ $t('common.delete') }}
    </button>
  </div>
</template>

<script setup>
// Props with validation
const props = defineProps({
  patient: {
    type: Object,
    required: true,
    validator: (value) => {
      return value && typeof value === 'object' && 'id' in value
    }
  }
})

// Emits with validation
defineEmits({
  'view': (patient) => patient && typeof patient === 'object',
  'edit': (patient) => patient && typeof patient === 'object',
  'delete': (patient) => patient && typeof patient === 'object'
})
</script>