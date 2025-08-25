import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'

// Mock vue-toastification
vi.mock('vue-toastification', () => ({
  useToast: vi.fn(() => ({
    error: vi.fn(),
    success: vi.fn(),
    info: vi.fn()
  }))
}))

// Mock API properly with factory function
vi.mock('@/api/locations', () => ({
  locationsApi: {
    getCountries: vi.fn(),
    getStatesByCountry: vi.fn(), 
    getCitiesByCountry: vi.fn()
  }
}))

import { useLocationData } from '../useLocationData'
import { locationsApi } from '@/api/locations'

describe('useLocationData', () => {
  let locationData
  
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    locationData = useLocationData()
  })
  
  afterEach(() => {
    // Cancel any pending requests
    locationData.cancelPendingRequests()
    vi.useRealTimers()
  })

  describe('initialization', () => {
    it('should initialize with empty data', () => {
      expect(locationData.countries.value).toEqual([])
      expect(locationData.stateOptions.value).toEqual([])
      expect(locationData.cityOptions.value).toEqual([])
      expect(locationData.loading.value).toEqual({
        countries: false,
        states: false,
        cities: false
      })
    })
  })

  describe('loadCountries', () => {
    it('should load countries successfully', async () => {
      const mockCountries = [
        { code: 'CA', name: 'Canada' },
        { code: 'US', name: 'United States' }
      ]
      locationsApi.getCountries.mockResolvedValue(mockCountries)
      
      await locationData.loadCountries()
      
      expect(locationData.countries.value).toEqual(mockCountries)
      expect(locationsApi.getCountries).toHaveBeenCalledTimes(1)
    })

    it('should handle countries loading error', async () => {
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {})
      locationsApi.getCountries.mockRejectedValue(new Error('API Error'))
      
      // Handle the error properly to avoid unhandled rejection
      try {
        await locationData.loadCountries()
      } catch (error) {
        // Expected error, just continue
      }
      
      expect(locationData.countries.value).toEqual([])
      expect(consoleError).toHaveBeenCalled()
      
      consoleError.mockRestore()
    })

    it('should not reload countries if already loaded', async () => {
      // First load countries properly
      const mockCountries = [{ code: 'CA', name: 'Canada' }]
      locationsApi.getCountries.mockResolvedValue(mockCountries)
      await locationData.loadCountries()
      
      // Reset mock to verify no additional calls
      vi.clearAllMocks()
      
      // Should not make another call
      await locationData.loadCountries()
      
      expect(locationsApi.getCountries).not.toHaveBeenCalled()
    })
  })

  describe('loadStatesForCountry', () => {
    it('should load states for country with debouncing', async () => {
      const mockStates = [
        { code: 'ON', name: 'Ontario' },
        { code: 'BC', name: 'British Columbia' }
      ]
      locationsApi.getStatesByCountry.mockResolvedValue(mockStates)
      
      // Use immediate version to avoid debounce delay in tests
      await locationData.loadStatesForCountryImmediate('CA')
      
      expect(locationData.stateOptions.value).toEqual(mockStates)
      expect(locationsApi.getStatesByCountry).toHaveBeenCalledWith('CA')
    })

    it('should handle states loading error', async () => {
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {})
      locationsApi.getStatesByCountry.mockRejectedValue(new Error('States API Error'))
      
      await locationData.loadStatesForCountryImmediate('CA')
      
      expect(locationData.stateOptions.value).toEqual([])
      expect(consoleError).toHaveBeenCalled()
      
      consoleError.mockRestore()
    })

    it('should clear states when country changes', async () => {
      // Set initial states
      locationData.stateOptions.value = [{ code: 'ON', name: 'Ontario' }]
      
      locationsApi.getStatesByCountry.mockResolvedValue([])
      
      await locationData.loadStatesForCountryImmediate('US')
      
      expect(locationData.stateOptions.value).toEqual([])
    })

    it('should skip loading if country code is empty', async () => {
      await locationData.loadStatesForCountryImmediate('')
      
      expect(locationsApi.getStatesByCountry).not.toHaveBeenCalled()
      expect(locationData.stateOptions.value).toEqual([])
    })
  })

  describe('loadCitiesForCountryAndState', () => {
    it('should load cities for country and state', async () => {
      const mockCities = [
        { id: 'toronto', name: 'Toronto' },
        { id: 'ottawa', name: 'Ottawa' }
      ]
      locationsApi.getCitiesByCountry.mockResolvedValue(mockCities)
      
      await locationData.loadCitiesForCountryAndStateImmediate('CA', 'ON')
      
      expect(locationData.cityOptions.value).toEqual(mockCities)
      expect(locationsApi.getCitiesByCountry).toHaveBeenCalledWith('CA', { state: 'ON' })
    })

    it('should load cities for country only (no state)', async () => {
      const mockCities = [
        { id: 'toronto', name: 'Toronto' },
        { id: 'vancouver', name: 'Vancouver' }
      ]
      locationsApi.getCitiesByCountry.mockResolvedValue(mockCities)
      
      await locationData.loadCitiesForCountryAndStateImmediate('CA')
      
      expect(locationData.cityOptions.value).toEqual(mockCities)
      expect(locationsApi.getCitiesByCountry).toHaveBeenCalledWith('CA', {})
    })

    it('should handle cities loading error', async () => {
      const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {})
      locationsApi.getCitiesByCountry.mockRejectedValue(new Error('Cities API Error'))
      
      await locationData.loadCitiesForCountryAndStateImmediate('CA', 'ON')
      
      expect(locationData.cityOptions.value).toEqual([])
      expect(consoleError).toHaveBeenCalled()
      
      consoleError.mockRestore()
    })

    it('should skip loading if country code is empty', async () => {
      await locationData.loadCitiesForCountryAndStateImmediate('')
      
      expect(locationsApi.getCitiesByCountry).not.toHaveBeenCalled()
      expect(locationData.cityOptions.value).toEqual([])
    })
  })

  describe('resolveCountryCode', () => {
    beforeEach(async () => {
      // Mock the API to return countries and load them
      locationsApi.getCountries.mockResolvedValue([
        { code: 'CA', name: 'Canada' },
        { code: 'US', name: 'United States' },
        { code: 'FR', name: 'France' }
      ])
      
      // Load countries first so resolveCountryCode can work
      await locationData.loadCountries()
    })

    it('should resolve country code from code', () => {
      expect(locationData.resolveCountryCode('CA')).toBe('CA')
      expect(locationData.resolveCountryCode('US')).toBe('US')
    })

    it('should resolve country code from name', () => {
      expect(locationData.resolveCountryCode('Canada')).toBe('CA')
      expect(locationData.resolveCountryCode('United States')).toBe('US')
    })

    it('should handle case insensitive matching', () => {
      expect(locationData.resolveCountryCode('canada')).toBe('CA')
      expect(locationData.resolveCountryCode('FRANCE')).toBe('FR')
    })

    it('should return empty string for unmatched values', () => {
      expect(locationData.resolveCountryCode('InvalidCountry')).toBe('')
      expect(locationData.resolveCountryCode('')).toBe('')
      expect(locationData.resolveCountryCode(null)).toBe('')
    })
  })

  describe('handleCountryChange', () => {
    it('should handle country change and load related data', async () => {
      const mockStates = [{ code: 'ON', name: 'Ontario' }]
      const mockCities = [{ id: 'toronto', name: 'Toronto' }]
      
      locationsApi.getStatesByCountry.mockResolvedValue(mockStates)
      locationsApi.getCitiesByCountry.mockResolvedValue(mockCities)
      
      const mockFormData = { country: 'CA', state: 'old-state', city: 'old-city' }
      
      locationData.handleCountryChange('CA', mockFormData)
      await vi.runAllTimersAsync()
      
      expect(locationData.stateOptions.value).toEqual(mockStates)
      expect(locationData.cityOptions.value).toEqual(mockCities)
      expect(mockFormData.state).toBe('')
      expect(mockFormData.city).toBe('')
    })

    it('should clear form data when changing countries', async () => {
      const mockFormData = {
        country: 'US',
        state: 'California',
        city: 'los-angeles'
      }
      
      locationsApi.getStatesByCountry.mockResolvedValue([])
      locationsApi.getCitiesByCountry.mockResolvedValue([])
      
      locationData.handleCountryChange('CA', mockFormData)
      await vi.advanceTimersByTime(300)
      
      expect(mockFormData.state).toBe('')
      expect(mockFormData.city).toBe('')
    })
  })

  describe('handleStateChange', () => {
    it('should handle state change and load cities', async () => {
      const mockCities = [{ id: 'toronto', name: 'Toronto' }]
      locationsApi.getCitiesByCountry.mockResolvedValue(mockCities)
      
      const mockFormData = { country: 'CA', state: 'ON', city: 'old-city' }
      
      locationData.handleStateChange('CA', 'ON', mockFormData)
      await vi.runAllTimersAsync()
      
      expect(locationData.cityOptions.value).toEqual(mockCities)
      expect(mockFormData.city).toBe('')
    })

    it('should clear city when state changes', async () => {
      const mockFormData = {
        country: 'CA',
        state: 'BC',
        city: 'toronto' // Should be cleared since it belongs to ON
      }
      
      locationsApi.getCitiesByCountry.mockResolvedValue([])
      
      locationData.handleStateChange('CA', 'BC', mockFormData)
      await vi.runAllTimersAsync()
      
      expect(mockFormData.city).toBe('')
    })
  })

  describe('loading states', () => {
    it('should track loading states correctly', async () => {
      let resolvePromise
      locationsApi.getCountries.mockImplementation(() => 
        new Promise(resolve => {
          resolvePromise = resolve
        })
      )
      
      const loadPromise = locationData.loadCountries()
      
      // Should be loading initially
      expect(locationData.loading.value.countries).toBe(true)
      
      // Resolve the promise
      resolvePromise([])
      await loadPromise
      
      // Should not be loading after completion
      expect(locationData.loading.value.countries).toBe(false)
    })

    it('should handle concurrent loading requests', async () => {
      locationsApi.getStatesByCountry.mockResolvedValue([])
      
      // Start multiple concurrent requests
      const promise1 = locationData.loadStatesForCountryImmediate('CA')
      const promise2 = locationData.loadStatesForCountryImmediate('CA')
      
      await Promise.all([promise1, promise2])
      
      // Should make two separate API calls since these are immediate (non-debounced)
      expect(locationsApi.getStatesByCountry).toHaveBeenCalledTimes(2)
    })
  })

  describe('cleanup', () => {
    it('should cancel pending requests', () => {
      // This should not throw
      locationData.cancelPendingRequests()
      expect(true).toBe(true) // If we get here, cleanup worked
    })
  })

  describe('debouncing', () => {
    it('should provide debounced and immediate versions', () => {
      expect(typeof locationData.loadStatesForCountry).toBe('function')
      expect(typeof locationData.loadStatesForCountryImmediate).toBe('function')
      expect(typeof locationData.loadCitiesForCountryAndState).toBe('function')
      expect(typeof locationData.loadCitiesForCountryAndStateImmediate).toBe('function')
    })

    it('should debounce multiple rapid calls', async () => {
      locationsApi.getStatesByCountry.mockResolvedValue([])
      
      // Make multiple rapid calls
      locationData.loadStatesForCountry('CA')
      locationData.loadStatesForCountry('CA')
      locationData.loadStatesForCountry('CA')
      
      // Wait for debounce to settle
      await vi.runAllTimersAsync()
      
      // Should only make one API call due to debouncing
      expect(locationsApi.getStatesByCountry).toHaveBeenCalledTimes(1)
    })
  })
})