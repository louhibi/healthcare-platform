<template>
  <div v-if="isDevelopment" class="fixed bottom-4 right-4 z-50">
    <!-- Toggle Button -->
    <button
      v-if="!showDemo"
      @click="showDemo = true"
      class="bg-blue-600 text-white px-4 py-2 rounded-lg shadow-lg hover:bg-blue-700 transition-colors"
      title="Show Virtual Scroll Demo"
    >
      📊 Virtual Scroll
    </button>

    <!-- Demo Panel -->
    <div v-else class="bg-white rounded-lg shadow-xl border p-4 max-w-sm">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-900">Virtual Scroll Demo</h3>
        <button
          @click="showDemo = false"
          class="text-gray-400 hover:text-gray-600"
        >
          ✕
        </button>
      </div>

      <!-- Metrics Display -->
      <div class="space-y-3 text-sm">
        <div class="grid grid-cols-2 gap-2">
          <div class="bg-gray-50 p-2 rounded">
            <div class="text-xs text-gray-500">Renders</div>
            <div class="font-semibold">{{ metrics.renders }}</div>
          </div>
          <div class="bg-gray-50 p-2 rounded">
            <div class="text-xs text-gray-500">Scroll Events</div>
            <div class="font-semibold">{{ metrics.scrollEvents }}</div>
          </div>
        </div>

        <div class="bg-gray-50 p-2 rounded">
          <div class="text-xs text-gray-500">Avg Visible Items</div>
          <div class="font-semibold">{{ metrics.averageVisibleItems }}</div>
        </div>

        <div class="bg-gray-50 p-2 rounded">
          <div class="text-xs text-gray-500">Avg Render Time</div>
          <div class="font-semibold">{{ metrics.averageRenderTime }}ms</div>
        </div>

        <div class="bg-gray-50 p-2 rounded">
          <div class="text-xs text-gray-500">Runtime</div>
          <div class="font-semibold">{{ metrics.runtimeMinutes }}min</div>
        </div>

        <div v-if="metrics.currentMemoryUsage" class="bg-gray-50 p-2 rounded">
          <div class="text-xs text-gray-500">Memory</div>
          <div class="font-semibold">{{ metrics.currentMemoryUsage.used }}</div>
          <div v-if="metrics.memoryTrend" class="text-xs" :class="memoryTrendClass">
            {{ metrics.memoryTrend.trend }} ({{ metrics.memoryTrend.usedChange }})
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex space-x-2 mt-4">
        <button
          @click="refreshMetrics"
          class="flex-1 bg-blue-100 text-blue-700 px-3 py-1 rounded text-xs hover:bg-blue-200"
        >
          Refresh
        </button>
        <button
          @click="resetMetrics"
          class="flex-1 bg-red-100 text-red-700 px-3 py-1 rounded text-xs hover:bg-red-200"
        >
          Reset
        </button>
        <button
          @click="logReport"
          class="flex-1 bg-green-100 text-green-700 px-3 py-1 rounded text-xs hover:bg-green-200"
        >
          Log
        </button>
      </div>

      <!-- Test Actions -->
      <div class="mt-4 pt-4 border-t border-gray-200">
        <div class="text-xs text-gray-500 mb-2">Test Actions:</div>
        <div class="flex space-x-2">
          <button
            @click="generateTestFields"
            class="flex-1 bg-purple-100 text-purple-700 px-2 py-1 rounded text-xs hover:bg-purple-200"
          >
            Add 50 Fields
          </button>
          <button
            @click="clearTestFields"
            class="flex-1 bg-gray-100 text-gray-700 px-2 py-1 rounded text-xs hover:bg-gray-200"
          >
            Clear
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { 
  getVirtualScrollMetrics, 
  resetVirtualScrollMetrics, 
  logVirtualScrollReport 
} from '@/utils/virtualScrollMetrics'

// Props
const props = defineProps({
  onGenerateFields: {
    type: Function,
    default: null
  },
  onClearFields: {
    type: Function,
    default: null
  }
})

// State
const showDemo = ref(false)
const metrics = ref({})
const updateInterval = ref(null)

// Computed
const isDevelopment = computed(() => import.meta.env.DEV)

const memoryTrendClass = computed(() => {
  const trend = metrics.value.memoryTrend?.trend
  if (trend === 'increasing') return 'text-red-600'
  if (trend === 'decreasing') return 'text-green-600'
  return 'text-gray-600'
})

// Methods
const refreshMetrics = () => {
  metrics.value = getVirtualScrollMetrics()
}

const resetMetrics = () => {
  resetVirtualScrollMetrics()
  refreshMetrics()
}

const logReport = () => {
  logVirtualScrollReport()
}

const generateTestFields = () => {
  if (props.onGenerateFields) {
    props.onGenerateFields()
  } else {
    // Default test field generation
    const testFields = []
    for (let i = 1; i <= 50; i++) {
      testFields.push({
        field_id: `test_field_${i}`,
        name: `test_field_${i}`,
        display_name: `Test Field ${i}`,
        field_type: ['text', 'select', 'textarea', 'date', 'number'][Math.floor(Math.random() * 5)],
        is_enabled: true,
        is_required: Math.random() > 0.7,
        sort_order: i,
        placeholder_text: `Enter value for field ${i}`,
        description: `This is test field number ${i} for virtual scrolling demo`
      })
    }
    
  }
}

const clearTestFields = () => {
  if (props.onClearFields) {
    props.onClearFields()
  } else {
    // Clear test fields - implement onClearFields prop
  }
}

// Lifecycle
onMounted(() => {
  refreshMetrics()
  
  // Update metrics every 2 seconds when demo is shown
  updateInterval.value = setInterval(() => {
    if (showDemo.value) {
      refreshMetrics()
    }
  }, 2000)
})

onUnmounted(() => {
  if (updateInterval.value) {
    clearInterval(updateInterval.value)
  }
})
</script>