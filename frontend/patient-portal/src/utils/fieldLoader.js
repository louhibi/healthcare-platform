/**
 * Field component loader utility with performance tracking
 * Provides metrics and caching for dynamically loaded field components
 */

// Component load cache
const componentCache = new Map()
const loadPromises = new Map()
const loadMetrics = {
  totalLoads: 0,
  cacheHits: 0,
  loadTimes: new Map(),
  failedLoads: new Set()
}

/**
 * Load a field component with caching and performance tracking
 * @param {string} fieldType - The field type to load
 * @param {Function} loader - The dynamic import function
 * @returns {Promise<Object>} The loaded component
 */
export async function loadFieldComponent(fieldType, loader) {
  const startTime = performance.now()
  
  // Return cached component if available
  if (componentCache.has(fieldType)) {
    loadMetrics.cacheHits++
    return componentCache.get(fieldType)
  }
  
  // Return existing load promise if component is already being loaded
  if (loadPromises.has(fieldType)) {
    return await loadPromises.get(fieldType)
  }
  
  // Create new load promise
  const loadPromise = (async () => {
    try {
      loadMetrics.totalLoads++
      const module = await loader()
      const component = module && (module.default || module)
      
      // Cache the loaded component (even if null/undefined)
      componentCache.set(fieldType, component)
      
      // Record load time
      const loadTime = performance.now() - startTime
      loadMetrics.loadTimes.set(fieldType, loadTime)
      
      
      return component
    } catch (error) {
      loadMetrics.failedLoads.add(fieldType)
      console.error(`[FieldLoader] Failed to load component for field type: ${fieldType}`, error)
      throw error
    } finally {
      // Remove from pending loads
      loadPromises.delete(fieldType)
    }
  })()
  
  // Store the promise to prevent duplicate loads
  loadPromises.set(fieldType, loadPromise)
  
  return await loadPromise
}

/**
 * Preload field components for better performance
 * @param {Array<string>} fieldTypes - Array of field types to preload
 * @param {Object} loaderMap - Map of field type to loader function
 */
export async function preloadFieldComponents(fieldTypes, loaderMap) {
  const preloadPromises = fieldTypes
    .filter(type => loaderMap[type] && !componentCache.has(type))
    .map(type => loadFieldComponent(type, loaderMap[type]).catch(() => null))
  
  await Promise.all(preloadPromises)
  
  if (import.meta.env.DEV) {
  }
}

/**
 * Get performance metrics for loaded components
 * @returns {Object} Performance metrics object
 */
export function getLoadMetrics() {
  return {
    ...loadMetrics,
    cacheHitRate: loadMetrics.totalLoads > 0 ? (loadMetrics.cacheHits / loadMetrics.totalLoads * 100).toFixed(2) : '0',
    averageLoadTime: loadMetrics.loadTimes.size > 0 
      ? (Array.from(loadMetrics.loadTimes.values()).reduce((sum, time) => sum + time, 0) / loadMetrics.loadTimes.size).toFixed(2)
      : '0'
  }
}

/**
 * Clear component cache (useful for development)
 */
export function clearComponentCache() {
  componentCache.clear()
  loadPromises.clear()
  loadMetrics.failedLoads.clear()
  loadMetrics.totalLoads = 0
  loadMetrics.cacheHits = 0
  loadMetrics.loadTimes.clear()
}

/**
 * Check if a component is cached
 * @param {string} fieldType - The field type to check
 * @returns {boolean} True if component is cached
 */
export function isComponentCached(fieldType) {
  return componentCache.has(fieldType)
}

/**
 * Get cached component names
 * @returns {Array<string>} Array of cached component names
 */
export function getCachedComponents() {
  return Array.from(componentCache.keys())
}

export default {
  loadFieldComponent,
  preloadFieldComponents,
  getLoadMetrics,
  clearComponentCache,
  isComponentCached,
  getCachedComponents
}