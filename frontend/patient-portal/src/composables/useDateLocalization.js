import { computed, ref, watchEffect, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useEntityStore } from '@/stores/entity'

/**
 * Composable for locale-aware date handling in healthcare contexts
 * Supports multiple date formats, locales, and healthcare entity preferences
 */
export function useDateLocalization() {
  const { locale, d: formatDate } = useI18n()
  const entityStore = useEntityStore()
  
  // Force reactivity to locale changes
  // Use the actual locale ref from useI18n instead of creating a new one
  const currentLocale = computed(() => locale.value)
  
  // Debug locale changes
  watch(locale, (newLocale) => {
  }, { immediate: true })

  // Date format configurations per country/locale
  const dateFormats = {
    // Morocco - supports Arabic, French, and English
    'MA': {
      'en-MA': {
        display: 'DD/MM/YYYY',        // British-style dates common in Morocco
        input: 'YYYY-MM-DD',          // HTML5 date input format
        placeholder: 'DD/MM/YYYY',
        example: '15/06/1985'
      },
      'fr-MA': {
        display: 'DD/MM/YYYY',        // French format
        input: 'YYYY-MM-DD', 
        placeholder: 'JJ/MM/AAAA',
        example: '15/06/1985'
      },
      'ar-MA': {
        display: 'DD/MM/YYYY',        // Common Arabic format
        input: 'YYYY-MM-DD',
        placeholder: 'ي ي/ش ش/س س س س',
        example: '١٥/٠٦/١٩٨٥'
      }
    },
    // Canada - French and English
    'CA': {
      'en-CA': {
        display: 'MM/DD/YYYY',        // North American format
        input: 'YYYY-MM-DD',
        placeholder: 'MM/DD/YYYY', 
        example: '06/15/1985'
      },
      'fr-CA': {
        display: 'DD/MM/YYYY',        // Quebec French format
        input: 'YYYY-MM-DD',
        placeholder: 'JJ/MM/AAAA',
        example: '15/06/1985'
      }
    },
    // United States - English only
    'US': {
      'en-US': {
        display: 'MM/DD/YYYY',        // US format
        input: 'YYYY-MM-DD',
        placeholder: 'MM/DD/YYYY',
        example: '06/15/1985'
      }
    },
    // France - French only
    'FR': {
      'fr-FR': {
        display: 'DD/MM/YYYY',        // European format
        input: 'YYYY-MM-DD', 
        placeholder: 'JJ/MM/AAAA',
        example: '15/06/1985'
      }
    }
  }

  // Get current format configuration
  const currentFormat = computed(() => {
    const countryCode = entityStore.entity?.country || 'MA'
    const localeCode = currentLocale.value || 'en-MA'
    
    
    // Try exact locale match first
    let formatConfig = dateFormats[countryCode]?.[localeCode]
    
    // Fallback to language-only match
    if (!formatConfig && localeCode.includes('-')) {
      const languageOnly = localeCode.split('-')[0]
      const fallbackLocale = Object.keys(dateFormats[countryCode] || {})
        .find(key => key.startsWith(languageOnly))
      
      if (fallbackLocale) {
        formatConfig = dateFormats[countryCode][fallbackLocale]
      }
    }
    
    // Final fallback to first available format for country
    if (!formatConfig && dateFormats[countryCode]) {
      const firstAvailable = Object.keys(dateFormats[countryCode])[0]
      formatConfig = dateFormats[countryCode][firstAvailable]
    }
    
    // Ultimate fallback
    return formatConfig || dateFormats.MA['en-MA']
  })

  // Convert ISO date string (YYYY-MM-DD) to display format
  const formatDateForDisplay = (isoDateString, options = {}) => {
    if (!isoDateString) return ''
    
    try {
      // Parse the ISO date string
      const date = new Date(isoDateString + 'T00:00:00') // Avoid timezone issues
      
      if (isNaN(date.getTime())) {
        console.warn('Invalid date provided to formatDateForDisplay:', isoDateString)
        return isoDateString
      }

      // Use Vue i18n datetime formatting if available
      if (options.useI18n) {
        return formatDate(date, 'short')
      }

      // Use manual formatting based on current format
      const format = currentFormat.value
      const day = String(date.getDate()).padStart(2, '0')
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const year = date.getFullYear()

      return format.display
        .replace('DD', day)
        .replace('MM', month) 
        .replace('YYYY', year)
    } catch (error) {
      console.error('Error formatting date:', error)
      return isoDateString
    }
  }

  // Convert display format back to ISO string (YYYY-MM-DD) for API
  const parseDisplayDate = (displayDate) => {
    if (!displayDate || typeof displayDate !== 'string') return ''
    
    const format = currentFormat.value
    let day, month, year

    try {
      if (format.display === 'DD/MM/YYYY') {
        const parts = displayDate.split('/')
        if (parts.length === 3) {
          day = parts[0]
          month = parts[1] 
          year = parts[2]
        }
      } else if (format.display === 'MM/DD/YYYY') {
        const parts = displayDate.split('/')
        if (parts.length === 3) {
          month = parts[0]
          day = parts[1]
          year = parts[2]
        }
      }

      if (day && month && year) {
        // Validate the parsed components
        const dayNum = parseInt(day, 10)
        const monthNum = parseInt(month, 10)
        const yearNum = parseInt(year, 10)

        if (dayNum >= 1 && dayNum <= 31 && 
            monthNum >= 1 && monthNum <= 12 && 
            yearNum >= 1900 && yearNum <= 2100) {
          
          // Return ISO format
          return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`
        }
      }
    } catch (error) {
      console.error('Error parsing display date:', error)
    }

    return ''
  }

  // Validate date string in current locale format
  const validateDateFormat = (dateString) => {
    const errors = []
    
    if (!dateString) return errors
    
    const format = currentFormat.value
    
    // Check basic format pattern
    const formatRegex = format.display === 'DD/MM/YYYY' 
      ? /^\d{2}\/\d{2}\/\d{4}$/
      : /^\d{2}\/\d{2}\/\d{4}$/
    
    if (!formatRegex.test(dateString)) {
      errors.push(`Please enter date in ${format.placeholder} format`)
      return errors
    }

    // Validate the actual date
    const isoDate = parseDisplayDate(dateString)
    if (!isoDate) {
      errors.push('Please enter a valid date')
      return errors
    }

    // Additional validations
    const date = new Date(isoDate)
    if (isNaN(date.getTime())) {
      errors.push('Please enter a valid date')
      return errors
    }

    // Check age constraints for healthcare context
    const today = new Date()
    const age = Math.floor((today - date) / (365.25 * 24 * 60 * 60 * 1000))
    
    if (age < 0) {
      errors.push('Birth date cannot be in the future')
    } else if (age > 150) {
      errors.push('Please enter a reasonable birth date')
    }

    return errors
  }

  // Get localized error messages
  const getLocalizedErrorMessage = (errorType) => {
    const format = currentFormat.value
    const messages = {
      invalid_format: `Please enter date in ${format.placeholder} format (e.g., ${format.example})`,
      invalid_date: 'Please enter a valid date',
      future_date: 'Birth date cannot be in the future',
      too_old: 'Please enter a reasonable birth date (within last 150 years)'
    }
    
    return messages[errorType] || 'Invalid date'
  }

  return {
    // State
    currentFormat,
    
    // Methods
    formatDateForDisplay,
    parseDisplayDate,
    validateDateFormat,
    getLocalizedErrorMessage,
    
    // Computed helpers
    placeholder: computed(() => currentFormat.value.placeholder),
    example: computed(() => currentFormat.value.example),
    displayFormat: computed(() => currentFormat.value.display),
    inputFormat: computed(() => currentFormat.value.input)
  }
}