import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

// Store for minimal bootstrap config (currently only environment)
export const useBootstrapStore = defineStore('bootstrap', () => {
  const environment = ref(null)
  const loaded = ref(false)
  const error = ref(null)
  const loading = ref(false)

  const isProd = computed(() => environment.value === 'production')
  const isDev = computed(() => !isProd.value)

  async function fetchBootstrap() {
    if (loaded.value || loading.value) return
    loading.value = true
    error.value = null
    try {
      const res = await api.get('/api/config/bootstrap')
      environment.value = res.data.environment || 'development'
      loaded.value = true
    } catch (e) {
      error.value = e.message || 'Failed to load bootstrap'
    } finally {
      loading.value = false
    }
  }

  return { environment, loaded, loading, error, isProd, isDev, fetchBootstrap }
})
