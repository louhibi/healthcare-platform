import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { debugApi } from '../api/debug.js'

// Store for handling debug/feature flag style toggles originating from backend
export const useDebugStore = defineStore('debug', () => {
  // State
  const enabled = ref(false)
  const lastChecked = ref(null)
  const isLoading = ref(false)
  const loadError = ref(null)
  // Local frontend-only debug logging toggle (persisted in localStorage)
  const logsEnabled = ref(localStorage.getItem('debug_logs_enabled') === 'true')

  // Getters
  const isEnabled = computed(() => enabled.value === true)
  const wasCheckedRecently = (maxAgeSeconds = 300) => {
    if (!lastChecked.value) return false
    const age = (Date.now() - new Date(lastChecked.value).getTime()) / 1000
    return age < maxAgeSeconds
  }

  // Actions
  const refresh = async (force = false) => {
    if (!force && wasCheckedRecently()) return enabled.value
    isLoading.value = true
    loadError.value = null
    try {
      const status = await debugApi.isEnabled()
      enabled.value = status
      lastChecked.value = new Date().toISOString()
      return enabled.value
    } catch (err) {
      console.error('Failed to refresh debug status', err)
      loadError.value = err
      return false
    } finally {
      isLoading.value = false
    }
  }

  const initialize = async () => {
    if (!lastChecked.value) {
      await refresh(true)
    }
    // Sync logsEnabled from localStorage on start
    logsEnabled.value = localStorage.getItem('debug_logs_enabled') === 'true'
  }

  const toggleLogs = (value) => {
    // Accept explicit boolean or invert
    if (typeof value === 'boolean') {
      logsEnabled.value = value
    } else {
      logsEnabled.value = !logsEnabled.value
    }
    localStorage.setItem('debug_logs_enabled', logsEnabled.value ? 'true' : 'false')
    return logsEnabled.value
  }

  return {
    // State
    enabled,
    lastChecked,
    isLoading,
    loadError,

    // Getters
    isEnabled,
  logsEnabled,

    // Actions
    refresh,
    initialize
  ,toggleLogs
  }
})
