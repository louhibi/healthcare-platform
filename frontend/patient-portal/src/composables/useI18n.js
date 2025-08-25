import { ref, computed, watch } from 'vue'
import { useI18n as useVueI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth.js'
import i18nAPI from '../api/i18n.js'
import { availableLocales, setLocale, getLocale, currentLocale } from '../i18n/index.js'

// Global reactive state
const supportedLocales = ref([])
const isInitializing = ref(false)
const backendTranslations = ref({})

export function useI18n() {
  const { t, locale, availableLocales: vueI18nLocales } = useVueI18n()
  const authStore = useAuthStore()

  // Computed properties
  const currentLocaleName = computed(() => {
    const localeInfo = availableLocales.find(l => l.code === currentLocale.value)
    return localeInfo?.name || currentLocale.value
  })

  const currentLocaleNativeName = computed(() => {
    const localeInfo = availableLocales.find(l => l.code === currentLocale.value)
    return localeInfo?.nativeName || currentLocale.value
  })

  const isRTL = computed(() => {
    return currentLocale.value.startsWith('ar')
  })

  // Initialize i18n system
  async function initialize() {
    if (isInitializing.value) return
    
    isInitializing.value = true
    try {
      // Load supported locales from backend
      await loadSupportedLocales()
      
      // Initialize user's preferred locale
      await initializeUserLocale()
      
      // Load backend translations for current locale
      await loadBackendTranslations(currentLocale.value)
      
    } catch (error) {
      console.error('Failed to initialize i18n:', error)
    } finally {
      isInitializing.value = false
    }
  }

  // Load supported locales from backend
  async function loadSupportedLocales() {
    try {
      const response = await i18nAPI.getSupportedLocales()
      supportedLocales.value = response.data || []
    } catch (error) {
      console.error('Failed to load supported locales:', error)
      // Use client-side locales as fallback
      supportedLocales.value = availableLocales
    }
  }

  // Initialize user's preferred locale
  async function initializeUserLocale() {
    let targetLocale = 'en-US'
    
    try {
      // If user is authenticated, get their preferred locale from store
      if (authStore.isAuthenticated && authStore.user) {
        targetLocale = i18nAPI.getUserPreferredLocale(authStore.user)
      } else {
        // Use browser locale or localStorage
        targetLocale = localStorage.getItem('healthcare-locale') || 
                     getBrowserLocale() || 
                     'en-US'
      }
      
      // Validate and set locale
      if (isLocaleSupported(targetLocale)) {
        await changeLocale(targetLocale, false) // Don't save to backend for initialization
      } else {
        await changeLocale('en-US', false)
      }
    } catch (error) {
      console.error('Failed to initialize user locale:', error)
      await changeLocale('en-US', false)
    }
  }

  // Change current locale
  async function changeLocale(newLocale, saveToBackend = true) {
    try {
      if (!isLocaleSupported(newLocale)) {
        throw new Error(`Locale ${newLocale} is not supported`)
      }

      // Update Vue I18n locale
      const success = setLocale(newLocale)
      if (!success) {
        throw new Error(`Failed to set locale to ${newLocale}`)
      }

      // Load backend translations for new locale
      await loadBackendTranslations(newLocale)

      // Save to backend if user is authenticated and requested
      if (authStore.isAuthenticated && saveToBackend) {
        try {
          await i18nAPI.updateUserLocale(newLocale)
        } catch (error) {
          console.warn('Failed to save locale preference to backend:', error)
          // Continue anyway - local change is still valid
        }
      }

      return true
    } catch (error) {
      console.error('Failed to change locale:', error)
      return false
    }
  }

  // Load translations from backend for a locale
  async function loadBackendTranslations(locale) {
    try {
      const response = await i18nAPI.getTranslations(locale)
      if (response.data && typeof response.data === 'object') {
        backendTranslations.value[locale] = response.data
        
        // Merge backend translations with client-side translations
        mergeBackendTranslations(locale, response.data)
      }
    } catch (error) {
      console.warn(`Failed to load backend translations for ${locale}:`, error)
      // Continue with client-side translations only
    }
  }

  // Merge backend translations with Vue I18n
  function mergeBackendTranslations(locale, translations) {
    try {
      // Validate translations data
      if (!translations || typeof translations !== 'object') {
        console.warn('Invalid translations data:', translations)
        return
      }

      // Convert flat key-value pairs to nested objects
      const nested = {}
      
      Object.entries(translations).forEach(([key, value]) => {
        try {
          // Skip invalid entries
          if (!key || typeof key !== 'string' || value === undefined) {
            return
          }

          const parts = key.split('.')
          let current = nested
          
          for (let i = 0; i < parts.length - 1; i++) {
            if (!current[parts[i]]) {
              current[parts[i]] = {}
            }
            current = current[parts[i]]
          }
          
          current[parts[parts.length - 1]] = value
        } catch (entryError) {
          console.warn('Failed to process translation entry:', key, value, entryError)
        }
      })

      // Merge with existing messages
      const existingMessages = vueI18nLocales[locale] || {}
      const mergedMessages = deepMerge(existingMessages, nested)
      
      // Update Vue I18n messages
      vueI18nLocales[locale] = mergedMessages
    } catch (error) {
      console.warn('Failed to merge backend translations:', error)
    }
  }

  // Deep merge utility
  function deepMerge(target, source) {
    const result = { ...target }
    
    Object.keys(source).forEach(key => {
      if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
        result[key] = deepMerge(target[key] || {}, source[key])
      } else {
        result[key] = source[key]
      }
    })
    
    return result
  }

  // Check if locale is supported
  function isLocaleSupported(locale) {
    return availableLocales.some(l => l.code === locale)
  }

  // Get browser locale
  function getBrowserLocale() {
    const browserLocale = navigator.language || navigator.languages?.[0]
    if (!browserLocale) return null
    
    // Try exact match first
    if (isLocaleSupported(browserLocale)) {
      return browserLocale
    }
    
    // Try language code match (e.g., 'en' for 'en-US')
    const languageCode = browserLocale.split('-')[0]
    const matchingLocale = availableLocales.find(l => l.code.startsWith(languageCode))
    return matchingLocale?.code || null
  }

  // Get localized field name
  function getLocalizedFieldName(fieldName, formType) {
    const key = `forms.${formType}.fields.${fieldName}.label`
    const translated = t(key)
    
    // If translation not found, return formatted field name
    if (translated === key) {
      return fieldName.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())
    }
    
    return translated
  }

  // Get localized form metadata from backend
  async function getLocalizedFormMetadata(formType, locale = null) {
    try {
      const targetLocale = locale || currentLocale.value
      const response = await i18nAPI.getLocalizedFormMetadata(formType, targetLocale)
      return response.data
    } catch (error) {
      console.error('Failed to get localized form metadata:', error)
      throw error
    }
  }

  // Watch for authentication changes
  watch(() => authStore.isAuthenticated, async (isAuthenticated) => {
    if (isAuthenticated) {
      // Re-initialize with user's preferred locale
      await initializeUserLocale()
    }
  })

  // Watch for locale changes to apply CSS direction
  watch(currentLocale, (newLocale) => {
    document.documentElement.lang = newLocale
    document.documentElement.dir = newLocale.startsWith('ar') ? 'rtl' : 'ltr'
  }, { immediate: true })

  return {
    // State
    currentLocale: computed(() => currentLocale.value),
    currentLocaleName,
    currentLocaleNativeName,
    availableLocales,
    supportedLocales: computed(() => supportedLocales.value),
    isInitializing: computed(() => isInitializing.value),
    isRTL,
    
    // Actions
    initialize,
    changeLocale,
    loadSupportedLocales,
    loadBackendTranslations,
    getLocalizedFieldName,
    getLocalizedFormMetadata,
    isLocaleSupported,
    
    // Vue I18n functions
    t,
    locale
  }
}