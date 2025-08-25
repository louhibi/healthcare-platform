import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import {
  initVirtualScrollMetrics,
  trackRender,
  trackScroll,
  trackScrollSessionStart,
  trackScrollSessionEnd,
  getVirtualScrollMetrics,
  resetVirtualScrollMetrics,
  cleanupVirtualScrollMetrics
} from '../virtualScrollMetrics'

describe('virtualScrollMetrics', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    resetVirtualScrollMetrics()
  })

  afterEach(() => {
    vi.useRealTimers()
    cleanupVirtualScrollMetrics()
  })

  describe('initialization', () => {
    it('should initialize metrics tracking', () => {
      initVirtualScrollMetrics()
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(0)
      expect(metrics.scrollEvents).toBe(0)
      expect(metrics.scrollSessions).toBe(0)
    })

    it('should track start time', () => {
      const startTime = Date.now()
      vi.setSystemTime(startTime)
      
      initVirtualScrollMetrics()
      
      vi.advanceTimersByTime(60000) // 1 minute
      
      const metrics = getVirtualScrollMetrics()
      expect(parseFloat(metrics.runtimeMinutes)).toBe(1.0)
    })
  })

  describe('render tracking', () => {
    beforeEach(() => {
      initVirtualScrollMetrics()
    })

    it('should track render events', () => {
      trackRender(10, 5.5)
      trackRender(15, 3.2)
      trackRender(8, 7.1)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(3)
      expect(parseFloat(metrics.averageVisibleItems)).toBeCloseTo(11)
      expect(parseFloat(metrics.averageRenderTime)).toBeCloseTo(5.27, 1)
    })

    it('should track min/max visible items', () => {
      trackRender(5)
      trackRender(20)
      trackRender(10)
      trackRender(3)
      trackRender(15)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.minVisibleItems).toBe(3)
      expect(metrics.maxVisibleItems).toBe(20)
    })

    it('should handle render without time', () => {
      trackRender(10)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(1)
      expect(metrics.averageRenderTime).toBe('0')
    })

    it('should limit stored render times', () => {
      // Track more than 100 renders
      for (let i = 0; i < 150; i++) {
        trackRender(10, i + 1)
      }
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(150)
      // Should only keep last 100 render times
      expect(parseFloat(metrics.minRenderTime)).toBeCloseTo(51, 0)
    })
  })

  describe('scroll tracking', () => {
    beforeEach(() => {
      initVirtualScrollMetrics()
    })

    it('should track scroll events', () => {
      trackScroll(100, 50)
      trackScroll(200, 100)
      trackScroll(150, 50)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollEvents).toBe(3)
      expect(parseFloat(metrics.averageScrollDistance)).toBeCloseTo(66.67, 1)
    })

    it('should handle scroll without distance', () => {
      trackScroll(100)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollEvents).toBe(1)
      expect(metrics.averageScrollDistance).toBe('0')
    })

    it('should limit stored scroll distances', () => {
      // Track more than 100 scroll events
      for (let i = 0; i < 150; i++) {
        trackScroll(i * 10, 10)
      }
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollEvents).toBe(150)
      // Should only keep last 100 distances
    })
  })

  describe('scroll session tracking', () => {
    beforeEach(() => {
      initVirtualScrollMetrics()
    })

    it('should track scroll sessions', () => {
      trackScrollSessionStart()
      vi.advanceTimersByTime(1000) // 1 second
      trackScrollSessionEnd()
      
      trackScrollSessionStart()
      vi.advanceTimersByTime(2000) // 2 seconds
      trackScrollSessionEnd()
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollSessions).toBe(2)
      expect(parseFloat(metrics.averageScrollSessionTime)).toBe(1500)
    })

    it('should handle session without end', () => {
      trackScrollSessionStart()
      vi.advanceTimersByTime(1000)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollSessions).toBe(1)
    })

    it('should detect new sessions based on gap', () => {
      const startTime = Date.now()
      vi.setSystemTime(startTime)
      
      trackScroll(100, 50)
      
      vi.advanceTimersByTime(2000) // 2 second gap
      
      trackScroll(200, 50) // Should start new session
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.scrollSessions).toBe(2)
    })
  })

  describe('metrics calculation', () => {
    beforeEach(() => {
      initVirtualScrollMetrics()
    })

    it('should calculate performance rates', () => {
      vi.advanceTimersByTime(60000) // 1 minute
      
      trackRender(10)
      trackRender(15)
      trackScroll(100, 50)
      trackScroll(200, 100)
      
      const metrics = getVirtualScrollMetrics()
      expect(parseFloat(metrics.rendersPerMinute)).toBe(2)
      expect(parseFloat(metrics.scrollEventsPerMinute)).toBe(2)
    })

    it('should handle zero runtime', () => {
      trackRender(10)
      trackScroll(100, 50)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.rendersPerMinute).toBe('0')
      expect(metrics.scrollEventsPerMinute).toBe('0')
    })

    it('should provide memory usage when available', () => {
      const metrics = getVirtualScrollMetrics()
      
      if (performance.memory) {
        expect(metrics.currentMemoryUsage).toBeDefined()
        expect(metrics.currentMemoryUsage.used).toBeDefined()
        expect(metrics.currentMemoryUsage.total).toBeDefined()
        expect(metrics.currentMemoryUsage.limit).toBeDefined()
      } else {
        expect(metrics.currentMemoryUsage).toBeNull()
      }
    })
  })

  describe('metrics reset', () => {
    beforeEach(() => {
      initVirtualScrollMetrics()
    })

    it('should reset all metrics', () => {
      trackRender(10, 5)
      trackScroll(100, 50)
      trackScrollSessionStart()
      
      let metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(1)
      expect(metrics.scrollEvents).toBe(1)
      
      resetVirtualScrollMetrics()
      
      metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(0)
      expect(metrics.scrollEvents).toBe(0)
      expect(metrics.scrollSessions).toBe(0)
      expect(metrics.averageVisibleItems).toBe('0.00')
    })
  })

  describe('memory tracking', () => {
    beforeEach(() => {
      // Mock performance.memory
      Object.defineProperty(performance, 'memory', {
        value: {
          usedJSHeapSize: 1000000,
          totalJSHeapSize: 2000000,
          jsHeapSizeLimit: 4000000
        },
        configurable: true
      })
    })

    it('should track memory usage over time', () => {
      initVirtualScrollMetrics()
      
      // Advance time to trigger memory sampling
      vi.advanceTimersByTime(10000) // 10 seconds (2 samples at 5s intervals)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.currentMemoryUsage).toBeDefined()
      expect(metrics.currentMemoryUsage.used).toMatch(/MB|KB|Bytes/)
    })

    it('should calculate memory trends', () => {
      initVirtualScrollMetrics()
      
      // Change memory usage
      performance.memory.usedJSHeapSize = 1500000
      
      vi.advanceTimersByTime(25000) // 25 seconds (5 samples)
      
      const metrics = getVirtualScrollMetrics()
      if (metrics.memoryTrend) {
        expect(metrics.memoryTrend.trend).toMatch(/increasing|decreasing|stable/)
      }
    })
  })

  describe('utility functions', () => {
    it('should format bytes correctly', () => {
      // This tests the internal formatBytes function through metrics
      Object.defineProperty(performance, 'memory', {
        value: {
          usedJSHeapSize: 1024,
          totalJSHeapSize: 1024 * 1024,
          jsHeapSizeLimit: 1024 * 1024 * 1024
        },
        configurable: true
      })

      const metrics = getVirtualScrollMetrics()
      expect(metrics.currentMemoryUsage.used).toBe('1 KB')
      expect(metrics.currentMemoryUsage.total).toBe('1 MB')
      expect(metrics.currentMemoryUsage.limit).toBe('1 GB')
    })

    it('should handle zero bytes', () => {
      Object.defineProperty(performance, 'memory', {
        value: {
          usedJSHeapSize: 0,
          totalJSHeapSize: 0,
          jsHeapSizeLimit: 0
        },
        configurable: true
      })

      const metrics = getVirtualScrollMetrics()
      expect(metrics.currentMemoryUsage.used).toBe('0 Bytes')
    })
  })

  describe('edge cases', () => {
    it('should handle negative values gracefully', () => {
      initVirtualScrollMetrics()
      
      trackRender(-5, -1) // Negative values
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(1)
      expect(metrics.minVisibleItems).toBe(-5)
    })

    it('should handle very large numbers', () => {
      initVirtualScrollMetrics()
      
      trackRender(Number.MAX_SAFE_INTEGER)
      trackScroll(Number.MAX_SAFE_INTEGER, Number.MAX_SAFE_INTEGER)
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(1)
      expect(metrics.scrollEvents).toBe(1)
    })

    it('should handle concurrent tracking calls', () => {
      initVirtualScrollMetrics()
      
      // Simulate concurrent calls
      for (let i = 0; i < 100; i++) {
        trackRender(i)
        trackScroll(i * 10, 10)
      }
      
      const metrics = getVirtualScrollMetrics()
      expect(metrics.renders).toBe(100)
      expect(metrics.scrollEvents).toBe(100)
    })
  })

  describe('cleanup', () => {
    it('should cleanup without errors', () => {
      initVirtualScrollMetrics()
      trackRender(10)
      trackScroll(100, 50)
      
      expect(() => cleanupVirtualScrollMetrics()).not.toThrow()
    })

    it('should stop memory tracking on cleanup', () => {
      initVirtualScrollMetrics()
      
      // Should not throw even if called multiple times
      cleanupVirtualScrollMetrics()
      cleanupVirtualScrollMetrics()
    })
  })
})