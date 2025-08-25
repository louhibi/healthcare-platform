import { locationsApi } from '@/api/locations'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { entitiesApi } from '../api/entities.js'

export const useEntityStore = defineStore('entity', () => {
  // State
  const entity = ref(null)
  const isLoading = ref(false)
  const lastFetchTime = ref(null)

  // Nationalities (simple direct load)
  const nationalityOptions = ref([])
  const loadNationalitiesImmediate = async () => {
    if (nationalityOptions.value.length) return
    try {
      const data = await locationsApi.getNationalities()
      nationalityOptions.value = Array.isArray(data) ? data.map(n => ({ value: n.id, label: n.name, country_id: n.country_id, is_primary: n.is_primary })) : []
      nationalityOptions.value.sort((a,b)=>a.label.localeCompare(b.label))
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error('[entityStore] Failed to load nationalities', e)
    }
  }

  // Getters
  const isEntityLoaded = computed(() => !!entity.value)
  const entityId = computed(() => entity.value?.id || null)
  const entityName = computed(() => entity.value?.name || '')
  const entityType = computed(() => entity.value?.type || '')
  const entityCountry = computed(() => entity.value?.country || '')
  const entityCountryId = computed(() => entity.value?.country_id || null)
  const entityState = computed(() => entity.value?.state || '')
  const entityStateId = computed(() => entity.value?.state_id || null)
  const entityCity = computed(() => entity.value?.city || '')
  const entityCityId = computed(() => entity.value?.city_id || null)
  const entityTimezone = computed(() => entity.value?.timezone || 'UTC')
  const entityLanguage = computed(() => entity.value?.language || 'en')
  const entityCurrency = computed(() => entity.value?.currency || 'USD')
  const entityAddress = computed(() => {
    if (!entity.value) return ''
    const parts = [
      entity.value.address,
      entity.value.city,
      entity.value.state,
      entity.value.postal_code,
      entity.value.country
    ].filter(Boolean)
    return parts.join(', ')
  })
  const entityContact = computed(() => ({
    phone: entity.value?.phone || '',
    email: entity.value?.email || '',
    website: entity.value?.website || ''
  }))
  const entityRequireRoomAssignment = computed(() => entity.value?.require_room_assignment || false)
  
  // Default nationality based on entity's country
  const entityDefaultNationalityId = computed(() => {
    const countryId = entityCountryId.value
    
    if (!countryId || !nationalityOptions.value.length) {
      return null
    }
    
    const defaultNationality = nationalityOptions.value.find(
      n => n.country_id === countryId && n.is_primary
    )
    
    return defaultNationality ? defaultNationality.value : null
  })

  // Actions
  const fetchEntity = async (entityIdParam) => {
    if (!entityIdParam) {
      console.warn('Cannot fetch entity: entity ID is required')
      return null
    }

    isLoading.value = true
    try {
      const data = await entitiesApi.getEntity(entityIdParam)
      
      entity.value = data
      lastFetchTime.value = new Date().toISOString()
      
      return entity.value
    } catch (error) {
      console.error('💥 Failed to fetch healthcare entity:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const clearEntity = () => {
    entity.value = null
    lastFetchTime.value = null
  }

  const updateEntity = (updatedData) => {
    if (entity.value) {
      entity.value = { ...entity.value, ...updatedData }
    }
  }

  // Helper methods for specific entity features
  const getEntityDisplayInfo = () => ({
    name: entityName.value,
    type: entityType.value,
    address: entityAddress.value,
    timezone: entityTimezone.value,
    currency: entityCurrency.value,
    language: entityLanguage.value
  })

  const getTimezoneInfo = () => ({
    timezone: entityTimezone.value,
    displayName: entityTimezone.value, // Could be enhanced with timezone display name
    abbreviation: entityTimezone.value // Could be enhanced with timezone abbreviation
  })

  const shouldRefreshEntity = (maxAgeMinutes = 60) => {
    if (!lastFetchTime.value) return true
    const ageMs = Date.now() - new Date(lastFetchTime.value).getTime()
    return ageMs > maxAgeMinutes * 60 * 1000
  }

  const ensureNationalitiesLoaded = async () => {
    await loadNationalitiesImmediate()
  }

  return {
    // State
    entity,
    isLoading,
    lastFetchTime,
    
    // Getters
    isEntityLoaded,
    entityId,
    entityName,
    entityType,
    entityCountry,
    entityCountryId,
    entityState,
    entityStateId,
    entityCity,
    entityCityId,
    entityTimezone,
    entityLanguage,
    entityCurrency,
    entityAddress,
    entityContact,
  entityRequireRoomAssignment,
  entityDefaultNationalityId,
  nationalityOptions,
    
    // Actions
    fetchEntity,
    clearEntity,
    updateEntity,
    getEntityDisplayInfo,
    getTimezoneInfo,
    shouldRefreshEntity,
  ensureNationalitiesLoaded,
  loadNationalitiesImmediate
  }
})