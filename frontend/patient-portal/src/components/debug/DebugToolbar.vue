<template>
  <div class="debug-toolbar">
    <div class="flex items-center space-x-4 text-xs">
      <span class="font-semibold">Debug</span>
      <!-- Logging toggle -->
      <label class="flex items-center space-x-1 cursor-pointer select-none">
        <input type="checkbox" class="hidden" :checked="debug.logsEnabled" @change="toggleLogs" />
        <span
          class="w-9 h-5 flex items-center bg-gray-600 rounded-full p-0.5 transition-colors"
          :class="{ 'bg-green-500': debug.logsEnabled }"
        >
          <span
            class="bg-white w-4 h-4 rounded-full shadow transform transition-transform"
            :class="{ 'translate-x-4': debug.logsEnabled }"
          ></span>
        </span>
        <span>{{ debug.logsEnabled ? 'Logs On' : 'Logs Off' }}</span>
      </label>
    </div>
  </div>
</template>

<script>
import { watch } from 'vue'
import { useDebugStore } from '../../stores/debug'
import logger from '../../utils/logger'

export default {
  name: 'DebugToolbar',
  setup() {
    const debug = useDebugStore()

    const toggleLogs = () => {
      const newVal = debug.toggleLogs()
      logger.info('Debug logging toggled', newVal)
    }

    // Log when backend debug enabled flag changes
    watch(() => debug.isEnabled, (val) => {
      if (debug.logsEnabled) {
        logger.debug('Backend debug enabled flag changed', val)
      }
    })

    return { debug, toggleLogs }
  }
}
</script>

<style scoped>
.debug-toolbar {
  position: fixed;
  bottom: 0.5rem;
  right: 0.5rem;
  background: rgba(31, 41, 55, 0.9); /* gray-800 */
  color: #f9fafb; /* gray-50 */
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  z-index: 2000;
  font-size: 0.75rem;
  line-height: 1rem;
}
</style>
