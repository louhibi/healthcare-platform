import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

import App from './App.vue'
import router from './router'
import i18n, { initializeLocale } from './i18n'
import './style.css'

const app = createApp(App)

// Pinia store
const pinia = createPinia()
app.use(pinia)

// Router
app.use(router)

// Internationalization
app.use(i18n)

// Toast notifications
app.use(Toast, {
  position: 'top-right',
  timeout: 5000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: 'button',
  icon: true,
  rtl: false
})

// Initialize auth store and i18n before mounting
import { useAuthStore } from './stores/auth'
import { useBootstrapStore } from './stores/bootstrap'
import { useI18n } from './composables/useI18n'

async function initializeApp() {
  const bootstrapStore = useBootstrapStore()
  await bootstrapStore.fetchBootstrap()
  const authStore = useAuthStore()
  
  // Initialize auth state from localStorage if available
  await authStore.initialize()
  
  // Initialize i18n system
  initializeLocale()
  
  // Mount the app
  app.mount('#app')
  
  // Initialize full i18n after mounting (when composables are available)
  setTimeout(async () => {
    try {
      const { initialize } = useI18n()
      await initialize()
    } catch (error) {
      console.warn('Failed to initialize i18n system:', error)
      // Continue without backend i18n - use client-side translations only
    }
  }, 100)
}

initializeApp()