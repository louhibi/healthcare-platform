/**
 * Timezone utilities for healthcare platform frontend
 * 
 * This module provides timezone conversion functions to handle the healthcare platform's
 * international operations. All datetime storage is in UTC, but display is in the
 * healthcare entity's local timezone.
 * 
 * CRITICAL: Use these functions for ALL datetime operations in the frontend.
 * Uses Luxon library for reliable timezone handling.
 */

import { DateTime } from 'luxon'

/**
 * Convert UTC datetime string to entity's local timezone
 * @param {string} utcDateTimeString - UTC datetime in ISO 8601 format (e.g., "2025-01-15T14:30:00Z")
 * @param {string} entityTimezone - IANA timezone identifier (e.g., "America/Toronto", "Europe/Paris")
 * @returns {Date} Date object in entity's timezone
 */
export function convertUTCToEntityTime(utcDateTimeString, entityTimezone) {
  if (!utcDateTimeString) return null
  
  if (!entityTimezone) {
    throw new Error('Entity timezone is required for timezone conversion. Check your configuration.')
  }
  
  try {
    const utcDateTime = DateTime.fromISO(utcDateTimeString, { zone: 'utc' })
    if (!utcDateTime.isValid) {
      throw new Error(`Invalid UTC datetime string: "${utcDateTimeString}". ${utcDateTime.invalidExplanation}`)
    }
    
    const entityDateTime = utcDateTime.setZone(entityTimezone)
    if (!entityDateTime.isValid) {
      throw new Error(`Invalid timezone: "${entityTimezone}". ${entityDateTime.invalidExplanation}`)
    }
    
    return entityDateTime.toJSDate()
  } catch (error) {
    throw new Error(`Failed to convert UTC to entity time: ${error.message}`)
  }
}

/**
 * Convert UTC datetime string to entity's local date (YYYY-MM-DD)
 * @param {string} utcDateTimeString - UTC datetime in ISO 8601 format
 * @param {string} entityTimezone - IANA timezone identifier
 * @returns {string} Date string in YYYY-MM-DD format in entity timezone
 */
export function convertUTCToEntityDate(utcDateTimeString, entityTimezone) {
  if (!utcDateTimeString) return ''
  
  if (!entityTimezone) {
    throw new Error('Entity timezone is required for timezone conversion. Check your configuration.')
  }
  
  try {
    const utcDateTime = DateTime.fromISO(utcDateTimeString, { zone: 'utc' })
    if (!utcDateTime.isValid) {
      throw new Error(`Invalid UTC datetime string: "${utcDateTimeString}". ${utcDateTime.invalidExplanation}`)
    }
    
    const entityDateTime = utcDateTime.setZone(entityTimezone)
    return entityDateTime.toISODate() // Returns YYYY-MM-DD format
  } catch (error) {
    throw new Error(`Failed to convert UTC to entity date: ${error.message}`)
  }
}

/**
 * Convert entity local datetime to UTC for API storage
 * @param {string} localDateTimeString - Local datetime string in various formats
 * @param {string} entityTimezone - IANA timezone identifier
 * @returns {string} UTC datetime in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)
 */
export function convertEntityTimeToUTC(localDateTimeString, entityTimezone) {
  if (!localDateTimeString) return ''
  
  if (!entityTimezone) {
    throw new Error('Entity timezone is required for timezone conversion. Check your configuration.')
  }
  
  // Normalize input format to YYYY-MM-DDTHH:mm:ss
  let normalizedDateTime
  if (localDateTimeString.includes('T')) {
    normalizedDateTime = localDateTimeString
  } else if (localDateTimeString.includes(' ')) {
    normalizedDateTime = localDateTimeString.replace(' ', 'T')
  } else {
    normalizedDateTime = localDateTimeString + 'T00:00:00'
  }
  
  // Ensure seconds are included
  if (normalizedDateTime.split('T')[1] && normalizedDateTime.split('T')[1].split(':').length === 2) {
    normalizedDateTime += ':00'
  }
  
  try {
    // Use Luxon to parse the local datetime in the entity timezone
    const entityDateTime = DateTime.fromISO(normalizedDateTime, { zone: entityTimezone })
    
    if (!entityDateTime.isValid) {
      throw new Error(`Invalid datetime "${normalizedDateTime}" in timezone "${entityTimezone}". ${entityDateTime.invalidExplanation}`)
    }
    
    // Convert to UTC and return ISO string
    const utcDateTime = entityDateTime.toUTC()
    return utcDateTime.toISO()
    
  } catch (error) {
    throw new Error(`Failed to convert "${localDateTimeString}" from timezone "${entityTimezone}" to UTC: ${error.message}`)
  }
}

/**
 * Get timezone offset in minutes for a specific date and timezone
 * @param {Date} date - Date to check offset for
 * @param {string} timezone - IANA timezone identifier
 * @returns {number} Offset in minutes (positive for ahead of UTC, negative for behind)
 */
function getTimezoneOffset(date, timezone) {
  // Create dates in UTC and in the target timezone
  const utcDate = new Date(date.toLocaleString('en-US', {timeZone: 'UTC'}))
  const tzDate = new Date(date.toLocaleString('en-US', {timeZone: timezone}))
  
  // Calculate difference in minutes
  return (tzDate.getTime() - utcDate.getTime()) / (1000 * 60)
}

/**
 * Format UTC datetime for display in entity's timezone
 * @param {string} utcDateTimeString - UTC datetime in ISO 8601 format
 * @param {string} entityTimezone - IANA timezone identifier
 * @param {Object} options - Formatting options (similar to Intl.DateTimeFormat options)
 * @returns {string} Formatted datetime string in entity timezone
 */
export function formatEntityDateTime(utcDateTimeString, entityTimezone, options = {}) {
  if (!utcDateTimeString) return ''
  
  const utcDate = new Date(utcDateTimeString)
  if (isNaN(utcDate.getTime())) {
    console.warn('Invalid UTC datetime string:', utcDateTimeString)
    return ''
  }
  
  const defaultOptions = {
    timeZone: entityTimezone,
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }
  
  // Clean options to avoid conflicts
  const cleanOptions = { ...defaultOptions }
  
  // Only merge compatible options (avoid timeStyle/dateStyle conflicts)
  Object.keys(options).forEach(key => {
    if (key !== 'timeStyle' && key !== 'dateStyle') {
      cleanOptions[key] = options[key]
    }
  })
  
  try {
    return utcDate.toLocaleString('en-US', cleanOptions)
  } catch (error) {
    console.warn('Date formatting error:', error, 'Options:', cleanOptions)
    // Fallback to basic formatting
    return utcDate.toLocaleString('en-US', { timeZone: entityTimezone })
  }
}

/**
 * Format UTC datetime as time only in entity's timezone
 * @param {string} utcDateTimeString - UTC datetime in ISO 8601 format
 * @param {string} entityTimezone - IANA timezone identifier
 * @param {boolean} hour12 - Whether to use 12-hour format (default: true)
 * @returns {string} Time string in entity timezone (e.g., "2:30 PM")
 */
export function formatEntityTime(utcDateTimeString, entityTimezone, hour12 = true) {
  if (!utcDateTimeString) return ''
  
  const utcDate = new Date(utcDateTimeString)
  if (isNaN(utcDate.getTime())) {
    console.warn('Invalid UTC datetime string:', utcDateTimeString)
    return ''
  }
  
  try {
    return utcDate.toLocaleString('en-US', {
      timeZone: entityTimezone,
      hour: '2-digit',
      minute: '2-digit',
      hour12: hour12
    })
  } catch (error) {
    console.warn('Time formatting error:', error)
    // Fallback to basic time formatting
    return utcDate.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: hour12
    })
  }
}

/**
 * Format UTC datetime as date only in entity's timezone
 * @param {string} utcDateTimeString - UTC datetime in ISO 8601 format
 * @param {string} entityTimezone - IANA timezone identifier
 * @param {Object} options - Date formatting options
 * @returns {string} Date string in entity timezone
 */
export function formatEntityDate(utcDateTimeString, entityTimezone, options = {}) {
  if (!utcDateTimeString) return ''
  
  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }
  
  return formatEntityDateTime(utcDateTimeString, entityTimezone, { 
    ...defaultOptions, 
    ...options,
    hour: undefined,
    minute: undefined,
    second: undefined
  })
}

/**
 * Create a UTC datetime string from date and time inputs in entity timezone
 * @param {string} dateString - Date in YYYY-MM-DD format
 * @param {string} timeString - Time in HH:MM format
 * @param {string} entityTimezone - IANA timezone identifier
 * @returns {string} UTC datetime in ISO 8601 format
 */
export function createUTCDateTime(dateString, timeString, entityTimezone) {
  if (!dateString || !timeString) return ''
  
  if (!entityTimezone) {
    throw new Error('Entity timezone is required for timezone conversion. Check your configuration.')
  }
  
  try {
    // Parse date and time components
    const [year, month, day] = dateString.split('-').map(Number)
    const [hour, minute] = timeString.split(':').map(Number)
    
    // Create DateTime in entity timezone
    const entityDateTime = DateTime.fromObject(
      { year, month, day, hour, minute, second: 0 },
      { zone: entityTimezone }
    )
    
    if (!entityDateTime.isValid) {
      throw new Error(`Invalid date/time "${dateString} ${timeString}" in timezone "${entityTimezone}". ${entityDateTime.invalidExplanation}`)
    }
    
    // Convert to UTC and return ISO string
    return entityDateTime.toUTC().toISO()
    
  } catch (error) {
    throw new Error(`Failed to create UTC datetime: ${error.message}`)
  }
}

/**
 * Get current datetime in UTC format for API requests
 * @returns {string} Current UTC datetime in ISO 8601 format
 */
export function getCurrentUTCDateTime() {
  return new Date().toISOString()
}

/**
 * Get current date in entity timezone (YYYY-MM-DD)
 * @param {string} entityTimezone - IANA timezone identifier
 * @returns {string} Current date in entity timezone
 */
export function getCurrentEntityDate(entityTimezone) {
  const now = new Date()
  return convertUTCToEntityDate(now.toISOString(), entityTimezone)
}

/**
 * Check if a UTC datetime is today in the entity's timezone
 * @param {string} utcDateTimeString - UTC datetime string
 * @param {string} entityTimezone - IANA timezone identifier
 * @returns {boolean} True if the UTC datetime is today in entity timezone
 */
export function isToday(utcDateTimeString, entityTimezone) {
  const entityDate = convertUTCToEntityDate(utcDateTimeString, entityTimezone)
  const todayEntityDate = getCurrentEntityDate(entityTimezone)
  return entityDate === todayEntityDate
}

/**
 * Get timezone display name for UI
 * @param {string} timezone - IANA timezone identifier
 * @returns {string} Human-readable timezone name (e.g., "Eastern Standard Time")
 */
export function getTimezoneDisplayName(timezone) {
  if (!timezone) return ''
  
  try {
    const now = new Date()
    const formatter = new Intl.DateTimeFormat('en-US', {
      timeZone: timezone,
      timeZoneName: 'long'
    })
    
    const parts = formatter.formatToParts(now)
    const timeZonePart = parts.find(part => part.type === 'timeZoneName')
    return timeZonePart ? timeZonePart.value : timezone
  } catch (error) {
    console.warn('Invalid timezone:', timezone)
    return timezone
  }
}

/**
 * Get timezone abbreviation for UI
 * @param {string} timezone - IANA timezone identifier  
 * @returns {string} Timezone abbreviation (e.g., "EST", "PST")
 */
export function getTimezoneAbbreviation(timezone) {
  if (!timezone) return ''
  
  try {
    const now = new Date()
    const formatter = new Intl.DateTimeFormat('en-US', {
      timeZone: timezone,
      timeZoneName: 'short'
    })
    
    const parts = formatter.formatToParts(now)
    const timeZonePart = parts.find(part => part.type === 'timeZoneName')
    return timeZonePart ? timeZonePart.value : timezone
  } catch (error) {
    console.warn('Invalid timezone:', timezone)
    return timezone
  }
}

// Export default object for convenience
export default {
  convertUTCToEntityTime,
  convertUTCToEntityDate,
  convertEntityTimeToUTC,
  formatEntityDateTime,
  formatEntityTime,
  formatEntityDate,
  createUTCDateTime,
  getCurrentUTCDateTime,
  getCurrentEntityDate,
  isToday,
  getTimezoneDisplayName,
  getTimezoneAbbreviation
}