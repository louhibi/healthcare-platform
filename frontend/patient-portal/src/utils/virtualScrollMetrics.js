/**
 * Virtual Scrolling Performance Metrics
 * Tracks and analyzes virtual scrolling performance for optimization
 */

// Performance tracking state
const metrics = {
  renders: 0,
  totalScrollTime: 0,
  scrollEvents: 0,
  averageVisibleItems: 0,
  maxVisibleItems: 0,
  minVisibleItems: Infinity,
  scrollDistances: [],
  renderTimes: [],
  lastScrollTime: 0,
  scrollSessions: 0,
  memoryUsage: [],
  startTime: Date.now()
}

// Performance observers
let renderObserver = null
let memoryObserver = null

/**
 * Initialize performance tracking
 */
export function initVirtualScrollMetrics() {
  metrics.startTime = Date.now()
  
  // Track memory usage if available
  if (performance.memory) {
    memoryObserver = setInterval(() => {
      metrics.memoryUsage.push({
        timestamp: Date.now(),
        used: performance.memory.usedJSHeapSize,
        total: performance.memory.totalJSHeapSize,
        limit: performance.memory.jsHeapSizeLimit
      })
      
      // Keep only last 100 measurements
      if (metrics.memoryUsage.length > 100) {
        metrics.memoryUsage = metrics.memoryUsage.slice(-100)
      }
    }, 5000) // Every 5 seconds
  }
  
}

/**
 * Track a render event
 * @param {number} visibleCount - Number of visible items
 * @param {number} renderTime - Time taken to render (ms)
 */
export function trackRender(visibleCount, renderTime = 0) {
  metrics.renders++
  
  if (renderTime > 0) {
    metrics.renderTimes.push(renderTime)
    
    // Keep only last 100 render times
    if (metrics.renderTimes.length > 100) {
      metrics.renderTimes = metrics.renderTimes.slice(-100)
    }
  }
  
  // Update visible item stats
  const totalVisible = metrics.averageVisibleItems * (metrics.renders - 1) + visibleCount
  metrics.averageVisibleItems = totalVisible / metrics.renders
  metrics.maxVisibleItems = Math.max(metrics.maxVisibleItems, visibleCount)
  metrics.minVisibleItems = Math.min(metrics.minVisibleItems, visibleCount)
}

/**
 * Track a scroll event
 * @param {number} scrollTop - Current scroll position
 * @param {number} scrollDistance - Distance scrolled since last event
 */
export function trackScroll(scrollTop, scrollDistance = 0) {
  metrics.scrollEvents++
  
  if (scrollDistance > 0) {
    metrics.scrollDistances.push(scrollDistance)
    
    // Keep only last 100 scroll distances
    if (metrics.scrollDistances.length > 100) {
      metrics.scrollDistances = metrics.scrollDistances.slice(-100)
    }
  }
  
  // Track scroll session timing
  const now = Date.now()
  if (metrics.lastScrollTime === 0) {
    // First scroll event starts the first session
    metrics.scrollSessions = 1
  } else {
    const timeDiff = now - metrics.lastScrollTime
    if (timeDiff > 1000) { // New session if gap > 1 second
      metrics.scrollSessions++
    }
    metrics.totalScrollTime += timeDiff
  }
  metrics.lastScrollTime = now
}

/**
 * Track scroll session start
 */
export function trackScrollSessionStart() {
  metrics.scrollSessions++
  metrics.lastScrollTime = Date.now()
}

/**
 * Track scroll session end
 */
export function trackScrollSessionEnd() {
  if (metrics.lastScrollTime > 0) {
    const sessionTime = Date.now() - metrics.lastScrollTime
    metrics.totalScrollTime += sessionTime
    metrics.lastScrollTime = 0
  }
}

/**
 * Get current performance metrics
 * @returns {Object} Performance metrics object
 */
export function getVirtualScrollMetrics() {
  const runtimeMinutes = (Date.now() - metrics.startTime) / 60000
  
  return {
    // Basic metrics
    renders: metrics.renders,
    scrollEvents: metrics.scrollEvents,
    scrollSessions: metrics.scrollSessions,
    runtimeMinutes: runtimeMinutes.toFixed(2),
    
    // Scroll performance
    averageScrollDistance: metrics.scrollDistances.length > 0 
      ? (metrics.scrollDistances.reduce((sum, d) => sum + d, 0) / metrics.scrollDistances.length).toFixed(2)
      : '0',
    totalScrollTime: metrics.totalScrollTime,
    averageScrollSessionTime: metrics.scrollSessions > 0 
      ? (metrics.totalScrollTime / metrics.scrollSessions).toFixed(2)
      : '0',
    
    // Render performance
    averageRenderTime: metrics.renderTimes.length > 0
      ? (metrics.renderTimes.reduce((sum, t) => sum + t, 0) / metrics.renderTimes.length).toFixed(2)
      : '0',
    maxRenderTime: metrics.renderTimes.length > 0
      ? Math.max(...metrics.renderTimes).toFixed(2)
      : '0',
    minRenderTime: metrics.renderTimes.length > 0
      ? Math.min(...metrics.renderTimes).toFixed(2)
      : '0',
    
    // Visible items stats
    averageVisibleItems: metrics.averageVisibleItems.toFixed(2),
    maxVisibleItems: metrics.maxVisibleItems === -Infinity ? 0 : metrics.maxVisibleItems,
    minVisibleItems: metrics.minVisibleItems === Infinity ? 0 : metrics.minVisibleItems,
    
    // Performance rates
    rendersPerMinute: runtimeMinutes > 0 ? (metrics.renders / runtimeMinutes).toFixed(2) : '0',
    scrollEventsPerMinute: runtimeMinutes > 0 ? (metrics.scrollEvents / runtimeMinutes).toFixed(2) : '0',
    
    // Memory usage (if available)
    currentMemoryUsage: performance.memory ? {
      used: formatBytes(performance.memory.usedJSHeapSize),
      total: formatBytes(performance.memory.totalJSHeapSize),
      limit: formatBytes(performance.memory.jsHeapSizeLimit)
    } : null,
    
    memoryTrend: getMemoryTrend()
  }
}

/**
 * Get memory usage trend
 * @returns {Object|null} Memory trend information
 */
function getMemoryTrend() {
  if (metrics.memoryUsage.length < 2) return null
  
  const recent = metrics.memoryUsage.slice(-10) // Last 10 measurements
  const first = recent[0]
  const last = recent[recent.length - 1]
  
  const usedDiff = last.used - first.used
  const timeDiff = last.timestamp - first.timestamp
  
  return {
    usedChange: formatBytes(usedDiff),
    timeSpan: `${(timeDiff / 1000).toFixed(1)}s`,
    trend: usedDiff > 0 ? 'increasing' : usedDiff < 0 ? 'decreasing' : 'stable'
  }
}

/**
 * Format bytes to human readable format
 * @param {number} bytes - Number of bytes
 * @returns {string} Formatted bytes string
 */
function formatBytes(bytes) {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Reset all metrics
 */
export function resetVirtualScrollMetrics() {
  Object.assign(metrics, {
    renders: 0,
    totalScrollTime: 0,
    scrollEvents: 0,
    averageVisibleItems: 0,
    maxVisibleItems: 0,
    minVisibleItems: Infinity,
    scrollDistances: [],
    renderTimes: [],
    lastScrollTime: 0,
    scrollSessions: 0,
    memoryUsage: [],
    startTime: Date.now()
  })
  
  if (import.meta.env.DEV) {
  }
}

/**
 * Log performance report to console (development only)
 */
export function logVirtualScrollReport() {
  if (!import.meta.env.DEV) return
  
  const report = getVirtualScrollMetrics()
  
  console.group('🚀 Virtual Scroll Performance Report')
  
  if (report.currentMemoryUsage) {
    if (report.memoryTrend) {
    }
  }
  
  console.groupEnd()
}

/**
 * Cleanup performance tracking
 */
export function cleanupVirtualScrollMetrics() {
  if (renderObserver) {
    clearInterval(renderObserver)
    renderObserver = null
  }
  
  if (memoryObserver) {
    clearInterval(memoryObserver)
    memoryObserver = null
  }
  
  if (import.meta.env.DEV) {
    logVirtualScrollReport()
  }
}

// Auto-initialize in development
if (import.meta.env.DEV) {
  initVirtualScrollMetrics()
  
  // Auto-report every 5 minutes in development
  setInterval(logVirtualScrollReport, 5 * 60 * 1000)
}

export default {
  initVirtualScrollMetrics,
  trackRender,
  trackScroll,
  trackScrollSessionStart,
  trackScrollSessionEnd,
  getVirtualScrollMetrics,
  resetVirtualScrollMetrics,
  logVirtualScrollReport,
  cleanupVirtualScrollMetrics
}