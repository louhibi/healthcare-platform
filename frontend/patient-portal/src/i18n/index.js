import { createI18n } from 'vue-i18n'
import { ref } from 'vue'

// Import locale messages - Morocco focused
import enMA from './locales/en-MA.json'
import frMA from './locales/fr-MA.json'

// Available locales configuration - Morocco focused
export const availableLocales = [
  {
    code: 'en-MA',
    name: 'English (Morocco)',
    nativeName: 'English (Morocco)',
    flag: '🇲🇦'
  },
  {
    code: 'fr-MA',
    name: 'French (Morocco)', 
    nativeName: 'Français (Maroc)',
    flag: '🇲🇦'
  }
]

// Create i18n instance
export const i18n = createI18n({
  legacy: false, // Use Composition API mode
  locale: 'en-MA', // Default locale - English Morocco
  fallbackLocale: 'en-MA',
  globalInjection: true,
  allowComposition: true,
  messages: {
    'en-MA': enMA,
    'fr-MA': frMA
  },
  // Enable number and datetime formatting
  numberFormats: {
    'en-MA': {
      currency: {
        style: 'currency',
        currency: 'MAD'
      }
    },
    'fr-MA': {
      currency: {
        style: 'currency',
        currency: 'MAD'
      }
    }
  },
  datetimeFormats: {
    'en-MA': {
      short: {
        year: 'numeric',
        month: 'short', 
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'long',
        hour: 'numeric',
        minute: 'numeric'
      }
    },
    'fr-MA': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long', 
        day: 'numeric',
        weekday: 'long',
        hour: 'numeric',
        minute: 'numeric'
      }
    }
  }
})

// Reactive current locale state
export const currentLocale = ref('en-MA')

// Helper functions
export function setLocale(locale) {
  if (availableLocales.some(l => l.code === locale)) {
    i18n.global.locale.value = locale
    currentLocale.value = locale
    
    // Update HTML lang attribute
    document.documentElement.lang = locale
    
    // Update HTML dir attribute for RTL languages
    document.documentElement.dir = locale.startsWith('ar') ? 'rtl' : 'ltr'
    
    // Save to localStorage
    localStorage.setItem('healthcare-locale', locale)
    
    return true
  }
  return false
}

export function getLocale() {
  return i18n.global.locale.value
}

export function initializeLocale() {
  // Try to get locale from localStorage first
  const savedLocale = localStorage.getItem('healthcare-locale')
  if (savedLocale && setLocale(savedLocale)) {
    return savedLocale
  }
  
  // Fallback to browser locale
  const browserLocale = navigator.language || navigator.languages?.[0]
  if (browserLocale) {
    // Try exact match first
    if (setLocale(browserLocale)) {
      return browserLocale
    }
    
    // Try language code match (e.g., 'en' for 'en-US')
    const languageCode = browserLocale.split('-')[0]
    const matchingLocale = availableLocales.find(l => l.code.startsWith(languageCode))
    if (matchingLocale && setLocale(matchingLocale.code)) {
      return matchingLocale.code
    }
  }
  
  // Final fallback
  setLocale('en-MA')
  return 'en-MA'
}

export function getLocalizedFieldName(fieldName, formType) {
  const key = `forms.${formType}.fields.${fieldName}.label`
  const translated = i18n.global.t(key)
  
  // If translation not found, return the original field name formatted
  if (translated === key) {
    return fieldName.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())
  }
  
  return translated
}

// Export i18n instance as default
export default i18n