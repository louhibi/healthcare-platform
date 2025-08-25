/**
 * Centralized date utilities for healthcare platform frontend
 * 
 * This module provides consistent date handling across the entire application
 * to prevent timezone-related bugs and ensure consistent date formatting.
 */

/**
 * Safely parse a date string (YYYY-MM-DD) without timezone conversion issues
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {Date} Date object at local midnight
 */
export function parseDate(dateString) {
  if (!dateString) return null
  
  // Add T00:00:00 to ensure it's treated as local time, not UTC
  return new Date(dateString + 'T00:00:00')
}

/**
 * Format a date string for display in locale-aware format
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @param {string} locale - Locale code (e.g., 'en-MA', 'fr-MA')
 * @param {Object} options - Intl.DateTimeFormat options
 * @returns {string} Formatted date string
 */
export function formatDate(dateString, locale = 'en-MA', options = {}) {
  if (!dateString) return ''
  
  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }
  
  const date = parseDate(dateString)
  return date.toLocaleDateString(locale, { ...defaultOptions, ...options })
}

/**
 * Format a date string for display based on healthcare entity locale
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @param {string} country - Country code (MA, CA, US, FR)
 * @param {string} locale - Locale code
 * @returns {string} Formatted date string
 */
export function formatDateForHealthcare(dateString, country = 'MA', locale = 'en-MA') {
  if (!dateString) return ''
  
  const date = parseDate(dateString)
  
  // Healthcare-specific formatting based on country
  const formatsByCountry = {
    'MA': { day: '2-digit', month: '2-digit', year: 'numeric' }, // DD/MM/YYYY
    'CA': locale.startsWith('fr') 
      ? { day: '2-digit', month: '2-digit', year: 'numeric' } // DD/MM/YYYY for French Canada
      : { month: '2-digit', day: '2-digit', year: 'numeric' }, // MM/DD/YYYY for English Canada
    'US': { month: '2-digit', day: '2-digit', year: 'numeric' }, // MM/DD/YYYY
    'FR': { day: '2-digit', month: '2-digit', year: 'numeric' }  // DD/MM/YYYY
  }
  
  const format = formatsByCountry[country] || formatsByCountry['MA']
  return date.toLocaleDateString(locale, format)
}

/**
 * Get the day of the week for a date string
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {string} Day of week (e.g., "Monday")
 */
export function getDayOfWeek(dateString) {
  if (!dateString) return ''
  
  const date = parseDate(dateString)
  return date.toLocaleDateString('en-US', { weekday: 'long' })
}

/**
 * Get the short day of the week for a date string
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {string} Short day of week (e.g., "Mon")
 */
export function getShortDayOfWeek(dateString) {
  if (!dateString) return ''
  
  const date = parseDate(dateString)
  return date.toLocaleDateString('en-US', { weekday: 'short' })
}

/**
 * Format date for input fields (YYYY-MM-DD)
 * @param {Date|string} date - Date object or date string
 * @returns {string} Date string in YYYY-MM-DD format
 */
export function formatDateForInput(date) {
  if (!date) return ''
  
  if (typeof date === 'string') {
    // If it's already a string, validate format and return as-is
    if (/^\d{4}-\d{2}-\d{2}$/.test(date)) {
      return date
    }
    date = parseDate(date)
  }
  
  if (!(date instanceof Date) || isNaN(date)) return ''
  
  return date.toISOString().split('T')[0]
}

/**
 * Get today's date in YYYY-MM-DD format
 * @returns {string} Today's date string
 */
export function getTodayString() {
  return new Date().toISOString().split('T')[0]
}

/**
 * Check if a date string represents today
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {boolean} True if the date is today
 */
export function isToday(dateString) {
  return dateString === getTodayString()
}

/**
 * Check if a date string represents a future date
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {boolean} True if the date is in the future
 */
export function isFuture(dateString) {
  return dateString > getTodayString()
}

/**
 * Check if a date string represents a past date
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {boolean} True if the date is in the past
 */
export function isPast(dateString) {
  return dateString < getTodayString()
}

/**
 * Add days to a date string
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @param {number} days - Number of days to add (can be negative)
 * @returns {string} New date string in YYYY-MM-DD format
 */
export function addDays(dateString, days) {
  const date = parseDate(dateString)
  date.setDate(date.getDate() + days)
  return formatDateForInput(date)
}

/**
 * Get the difference in days between two date strings
 * @param {string} dateString1 - First date string in YYYY-MM-DD format
 * @param {string} dateString2 - Second date string in YYYY-MM-DD format
 * @returns {number} Difference in days (dateString1 - dateString2)
 */
export function daysDifference(dateString1, dateString2) {
  const date1 = parseDate(dateString1)
  const date2 = parseDate(dateString2)
  const diffTime = date1 - date2
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

/**
 * Get the start and end of a month for a given date string
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {Object} Object with startDate and endDate strings
 */
export function getMonthRange(dateString) {
  const date = parseDate(dateString)
  const year = date.getFullYear()
  const month = date.getMonth()
  
  const startDate = new Date(year, month, 1)
  const endDate = new Date(year, month + 1, 0) // Last day of month
  
  return {
    startDate: formatDateForInput(startDate),
    endDate: formatDateForInput(endDate)
  }
}

/**
 * Format date and time together
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @param {string} timeString - Time string in HH:MM format
 * @param {Object} options - Formatting options
 * @returns {string} Formatted date and time string
 */
export function formatDateTime(dateString, timeString, options = {}) {
  if (!dateString) return ''
  
  const dateOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    ...(timeString && {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    }),
    ...options
  }
  
  if (timeString) {
    const [hours, minutes] = timeString.split(':').map(Number)
    const date = parseDate(dateString)
    date.setHours(hours, minutes, 0, 0)
    return date.toLocaleString('en-US', dateOptions)
  }
  
  return formatDate(dateString, dateOptions)
}

/**
 * Format time from a datetime string (ISO 8601 or similar)
 * @param {string} dateTimeString - Full datetime string (e.g., "2024-01-15T14:30:00Z")
 * @param {Object} options - Formatting options
 * @returns {string} Formatted time string (e.g., "2:30 PM")
 */
export function formatTimeFromDateTime(dateTimeString, options = {}) {
  if (!dateTimeString) return ''
  
  const defaultOptions = {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }
  
  const date = new Date(dateTimeString)
  return date.toLocaleTimeString('en-US', { ...defaultOptions, ...options })
}

/**
 * Format time in 24-hour format from a datetime string
 * @param {string} dateTimeString - Full datetime string
 * @returns {string} Time string in HH:MM format (e.g., "14:30")
 */
export function formatTime24(dateTimeString) {
  if (!dateTimeString) return ''
  
  const date = new Date(dateTimeString)
  return date.toLocaleTimeString('en-US', { 
    hour: '2-digit', 
    minute: '2-digit', 
    hour12: false 
  })
}

/**
 * Validate date string format (YYYY-MM-DD)
 * @param {string} dateString - Date string to validate
 * @returns {boolean} True if valid format
 */
export function isValidDateString(dateString) {
  if (!dateString || typeof dateString !== 'string') return false
  
  const regex = /^\d{4}-\d{2}-\d{2}$/
  if (!regex.test(dateString)) return false
  
  const date = parseDate(dateString)
  return date instanceof Date && !isNaN(date)
}

/**
 * Get age from a birth date string
 * @param {string} birthDateString - Birth date string in YYYY-MM-DD format
 * @param {string} [referenceDate] - Reference date string (defaults to today)
 * @returns {number} Age in years
 */
export function calculateAge(birthDateString, referenceDate = getTodayString()) {
  const birthDate = parseDate(birthDateString)
  const refDate = parseDate(referenceDate)
  
  let age = refDate.getFullYear() - birthDate.getFullYear()
  const monthDiff = refDate.getMonth() - birthDate.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && refDate.getDate() < birthDate.getDate())) {
    age--
  }
  
  return age
}

/**
 * Format relative time (e.g., "2 days ago", "in 3 days")
 * @param {string} dateString - Date string in YYYY-MM-DD format
 * @returns {string} Relative time string
 */
export function formatRelativeTime(dateString) {
  const days = daysDifference(dateString, getTodayString())
  
  if (days === 0) return 'Today'
  if (days === 1) return 'Tomorrow'
  if (days === -1) return 'Yesterday'
  if (days > 1) return `In ${days} days`
  if (days < -1) return `${Math.abs(days)} days ago`
  
  return dateString
}

/**
 * Get an array of date strings for a date range
 * @param {string} startDate - Start date string in YYYY-MM-DD format
 * @param {string} endDate - End date string in YYYY-MM-DD format
 * @returns {string[]} Array of date strings
 */
export function getDateRange(startDate, endDate) {
  const dates = []
  let current = startDate
  
  while (current <= endDate) {
    dates.push(current)
    current = addDays(current, 1)
  }
  
  return dates
}

// Export default object for convenience
export default {
  parseDate,
  formatDate,
  getDayOfWeek,
  getShortDayOfWeek,
  formatDateForInput,
  getTodayString,
  isToday,
  isFuture,
  isPast,
  addDays,
  daysDifference,
  getMonthRange,
  formatDateTime,
  formatTimeFromDateTime,
  formatTime24,
  isValidDateString,
  calculateAge,
  formatRelativeTime,
  getDateRange
}