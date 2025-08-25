import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import {
  loadFieldComponent,
  preloadFieldComponents,
  getLoadMetrics,
  clearComponentCache,
  isComponentCached,
  getCachedComponents
} from '../fieldLoader'

describe('fieldLoader', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    clearComponentCache()
  })

  afterEach(() => {
    vi.useRealTimers()
    clearComponentCache()
  })

  describe('loadFieldComponent', () => {
    it('should load and cache component successfully', async () => {
      const mockComponent = { name: 'TestComponent' }
      const mockLoader = vi.fn().mockResolvedValue({ default: mockComponent })
      
      const result = await loadFieldComponent('test', mockLoader)
      
      expect(result).toEqual(mockComponent)
      expect(mockLoader).toHaveBeenCalledTimes(1)
      expect(isComponentCached('test')).toBe(true)
    })

    it('should return cached component on subsequent calls', async () => {
      const mockComponent = { name: 'TestComponent' }
      const mockLoader = vi.fn().mockResolvedValue({ default: mockComponent })
      
      const result1 = await loadFieldComponent('test', mockLoader)
      const result2 = await loadFieldComponent('test', mockLoader)
      
      expect(result1).toEqual(result2)
      expect(mockLoader).toHaveBeenCalledTimes(1) // Should only load once
    })

    it('should handle component without default export', async () => {
      const mockComponent = { name: 'TestComponent' }
      const mockLoader = vi.fn().mockResolvedValue(mockComponent) // No .default
      
      const result = await loadFieldComponent('test', mockLoader)
      
      expect(result).toEqual(mockComponent)
    })

    it('should handle loading errors', async () => {
      const mockLoader = vi.fn().mockRejectedValue(new Error('Loading failed'))
      
      await expect(loadFieldComponent('test', mockLoader)).rejects.toThrow('Loading failed')
      expect(isComponentCached('test')).toBe(false)
    })

    it('should prevent duplicate loading of same component', async () => {
      const mockLoader = vi.fn().mockImplementation(() => 
        new Promise(resolve => setTimeout(() => resolve({ default: { name: 'Test' } }), 100))
      )
      
      // Start multiple concurrent loads
      const promise1 = loadFieldComponent('test', mockLoader)
      const promise2 = loadFieldComponent('test', mockLoader)
      const promise3 = loadFieldComponent('test', mockLoader)
      
      vi.advanceTimersByTime(100)
      
      const [result1, result2, result3] = await Promise.all([promise1, promise2, promise3])
      
      expect(result1).toEqual(result2)
      expect(result2).toEqual(result3)
      expect(mockLoader).toHaveBeenCalledTimes(1) // Should only call loader once
    })

    it('should track load times in development', async () => {
      // Mock import.meta.env.DEV
      vi.stubGlobal('import.meta', { env: { DEV: true } })
      
      const mockComponent = { name: 'TestComponent' }
      const mockLoader = vi.fn().mockResolvedValue({ default: mockComponent })
      
      await loadFieldComponent('test', mockLoader)
      
      const metrics = getLoadMetrics()
      expect(metrics.totalLoads).toBe(1)
      expect(metrics.loadTimes).toBeDefined()
    })
  })

  describe('preloadFieldComponents', () => {
    it('should preload multiple components', async () => {
      const mockComponent1 = { name: 'Component1' }
      const mockComponent2 = { name: 'Component2' }
      const mockLoader1 = vi.fn().mockResolvedValue({ default: mockComponent1 })
      const mockLoader2 = vi.fn().mockResolvedValue({ default: mockComponent2 })
      
      const loaderMap = {
        type1: mockLoader1,
        type2: mockLoader2
      }
      
      await preloadFieldComponents(['type1', 'type2'], loaderMap)
      
      expect(isComponentCached('type1')).toBe(true)
      expect(isComponentCached('type2')).toBe(true)
      expect(mockLoader1).toHaveBeenCalledTimes(1)
      expect(mockLoader2).toHaveBeenCalledTimes(1)
    })

    it('should skip already cached components', async () => {
      const mockComponent1 = { name: 'Component1' }
      const mockLoader1 = vi.fn().mockResolvedValue({ default: mockComponent1 })
      
      // Pre-cache one component
      await loadFieldComponent('type1', mockLoader1)
      mockLoader1.mockClear()
      
      const mockLoader2 = vi.fn().mockResolvedValue({ default: { name: 'Component2' } })
      const loaderMap = {
        type1: mockLoader1,
        type2: mockLoader2
      }
      
      await preloadFieldComponents(['type1', 'type2'], loaderMap)
      
      expect(mockLoader1).not.toHaveBeenCalled() // Should skip cached
      expect(mockLoader2).toHaveBeenCalledTimes(1) // Should load new
    })

    it('should filter out types without loaders', async () => {
      const mockLoader = vi.fn().mockResolvedValue({ default: { name: 'Component' } })
      const loaderMap = {
        type1: mockLoader
      }
      
      await preloadFieldComponents(['type1', 'type2', 'type3'], loaderMap)
      
      expect(mockLoader).toHaveBeenCalledTimes(1)
      expect(isComponentCached('type1')).toBe(true)
      expect(isComponentCached('type2')).toBe(false)
      expect(isComponentCached('type3')).toBe(false)
    })

    it('should handle preload errors gracefully', async () => {
      const mockLoader1 = vi.fn().mockResolvedValue({ default: { name: 'Component1' } })
      const mockLoader2 = vi.fn().mockRejectedValue(new Error('Load failed'))
      
      const loaderMap = {
        type1: mockLoader1,
        type2: mockLoader2
      }
      
      // Should not throw even if some components fail
      await expect(preloadFieldComponents(['type1', 'type2'], loaderMap)).resolves.toBeUndefined()
      
      expect(isComponentCached('type1')).toBe(true)
      expect(isComponentCached('type2')).toBe(false)
    })
  })

  describe('cache management', () => {
    it('should track cached components', async () => {
      const mockLoader = vi.fn().mockResolvedValue({ default: { name: 'Test' } })
      
      expect(getCachedComponents()).toEqual([])
      
      await loadFieldComponent('test1', mockLoader)
      await loadFieldComponent('test2', mockLoader)
      
      const cached = getCachedComponents()
      expect(cached).toContain('test1')
      expect(cached).toContain('test2')
      expect(cached).toHaveLength(2)
    })

    it('should clear component cache', async () => {
      const mockLoader = vi.fn().mockResolvedValue({ default: { name: 'Test' } })
      
      await loadFieldComponent('test', mockLoader)
      expect(isComponentCached('test')).toBe(true)
      
      clearComponentCache()
      
      expect(isComponentCached('test')).toBe(false)
      expect(getCachedComponents()).toEqual([])
      
      const metrics = getLoadMetrics()
      expect(metrics.totalLoads).toBe(0)
      expect(metrics.cacheHits).toBe(0)
    })
  })

  describe('metrics tracking', () => {
    beforeEach(() => {
      clearComponentCache()
    })

    it('should track load metrics', async () => {
      const mockLoader = vi.fn().mockResolvedValue({ default: { name: 'Test' } })
      
      await loadFieldComponent('test1', mockLoader)
      await loadFieldComponent('test2', mockLoader)
      await loadFieldComponent('test1', mockLoader) // Should hit cache
      
      const metrics = getLoadMetrics()
      expect(metrics.totalLoads).toBe(2)
      expect(metrics.cacheHits).toBe(1)
      expect(parseFloat(metrics.cacheHitRate)).toBeCloseTo(50.0, 1) // 1 cache hit out of 2 total loads = 50%
    })

    it('should calculate average load time', async () => {
      const mockLoader = vi.fn()
        .mockResolvedValueOnce({ default: { name: 'Test1' } })
        .mockResolvedValueOnce({ default: { name: 'Test2' } })
      
      // Mock performance.now to control timing
      const originalNow = performance.now
      let callCount = 0
      Object.defineProperty(performance, 'now', {
        value: vi.fn(() => {
          callCount++
          // First call is start time, second is end time for each load
          // Load 1: start=10, end=20 -> diff=10
          // Load 2: start=30, end=40 -> diff=10
          // Average = (10+10)/2 = 10
          return callCount * 10 
        }),
        configurable: true
      })
      
      try {
        await loadFieldComponent('test1', mockLoader)
        await loadFieldComponent('test2', mockLoader)
        
        const metrics = getLoadMetrics()
        expect(parseFloat(metrics.averageLoadTime)).toBe(10) // (10+10)/2 = 10
      } finally {
        Object.defineProperty(performance, 'now', {
          value: originalNow,
          configurable: true
        })
      }
    })

    it('should track failed loads', async () => {
      const mockLoader1 = vi.fn().mockResolvedValue({ default: { name: 'Success' } })
      const mockLoader2 = vi.fn().mockRejectedValue(new Error('Failed'))
      
      await loadFieldComponent('success', mockLoader1)
      
      try {
        await loadFieldComponent('failed', mockLoader2)
      } catch (e) {
        // Expected to fail
      }
      
      const metrics = getLoadMetrics()
      expect(metrics.totalLoads).toBe(2)
      expect(metrics.failedLoads).toContain('failed')
      expect(metrics.failedLoads).not.toContain('success')
    })

    it('should handle zero loads gracefully', () => {
      const metrics = getLoadMetrics()
      
      expect(metrics.totalLoads).toBe(0)
      expect(metrics.cacheHits).toBe(0)
      expect(metrics.cacheHitRate).toBe('0')
      expect(metrics.averageLoadTime).toBe('0')
    })
  })

  describe('error handling', () => {
    it('should handle loader that throws synchronously', async () => {
      const mockLoader = vi.fn().mockImplementation(() => {
        throw new Error('Sync error')
      })
      
      await expect(loadFieldComponent('test', mockLoader)).rejects.toThrow('Sync error')
    })

    it('should handle loader that returns null', async () => {
      const mockLoader = vi.fn().mockResolvedValue(null)
      
      const result = await loadFieldComponent('test', mockLoader)
      
      expect(result).toBeNull()
      expect(isComponentCached('test')).toBe(true) // Should still cache
    })

    it('should handle loader that returns undefined', async () => {
      const mockLoader = vi.fn().mockResolvedValue(undefined)
      
      const result = await loadFieldComponent('test', mockLoader)
      
      expect(result).toBeUndefined()
      expect(isComponentCached('test')).toBe(true)
    })
  })

  describe('performance', () => {
    it('should not block on concurrent loads of different components', async () => {
      const mockLoader1 = vi.fn().mockImplementation(() => 
        new Promise(resolve => setTimeout(() => resolve({ default: { name: 'Component1' } }), 100))
      )
      const mockLoader2 = vi.fn().mockImplementation(() => 
        new Promise(resolve => setTimeout(() => resolve({ default: { name: 'Component2' } }), 100))
      )
      
      const startTime = Date.now()
      const promise1 = loadFieldComponent('test1', mockLoader1)
      const promise2 = loadFieldComponent('test2', mockLoader2)
      
      vi.advanceTimersByTime(100)
      
      await Promise.all([promise1, promise2])
      
      // Both should complete in parallel, not sequential
      expect(Date.now() - startTime).toBeLessThan(150)
    })
  })
})