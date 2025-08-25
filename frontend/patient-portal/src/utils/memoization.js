import { computed, unref } from 'vue'

/**
 * Memoization utilities for Vue.js composables and components
 */

// Simple memoization for pure functions
const memoCache = new Map()

export function memoize(fn, keyGenerator = (...args) => JSON.stringify(args)) {
  return function (...args) {
    const key = keyGenerator(...args)
    
    if (memoCache.has(key)) {
      return memoCache.get(key)
    }
    
    const result = fn.apply(this, args)
    memoCache.set(key, result)
    
    return result
  }
}

// Clear memoization cache
export function clearMemoCache() {
  memoCache.clear()
}

// Memoized computed for expensive operations
export function memoizedComputed(getter, keyFn = null) {
  const cache = new Map()
  
  return computed(() => {
    // Generate cache key from dependencies
    const key = keyFn ? keyFn() : getter.toString()
    
    // Check if we have a cached result for the same dependencies
    if (cache.has(key)) {
      return cache.get(key)
    }
    
    // Compute new result and cache it
    const result = getter()
    cache.set(key, result)
    
    // Clean up old cache entries (keep only last 5)
    if (cache.size > 5) {
      const firstKey = cache.keys().next().value
      cache.delete(firstKey)
    }
    
    return result
  })
}

// Memoization for array operations
export function memoizedArraySort(array, compareFn, dependencies = []) {
  const cache = new WeakMap()
  const lastDeps = new Map()
  
  return computed(() => {
    const arr = unref(array)
    const deps = dependencies.map(dep => unref(dep))
    const depsKey = JSON.stringify(deps)
    
    // Check if dependencies changed
    if (lastDeps.get(arr) === depsKey && cache.has(arr)) {
      return cache.get(arr)
    }
    
    // Sort the array
    const sorted = [...arr].sort(compareFn)
    
    // Cache the result
    cache.set(arr, sorted)
    lastDeps.set(arr, depsKey)
    
    return sorted
  })
}

// Memoization for filtering operations
export function memoizedArrayFilter(array, filterFn, dependencies = []) {
  const cache = new WeakMap()
  const lastDeps = new Map()
  
  return computed(() => {
    const arr = unref(array)
    const deps = dependencies.map(dep => unref(dep))
    const depsKey = JSON.stringify(deps)
    
    // Check if dependencies changed
    if (lastDeps.get(arr) === depsKey && cache.has(arr)) {
      return cache.get(arr)
    }
    
    // Filter the array
    const filtered = arr.filter(filterFn)
    
    // Cache the result
    cache.set(arr, filtered)
    lastDeps.set(arr, depsKey)
    
    return filtered
  })
}

// Memoization for grouping operations
export function memoizedGroupBy(array, keyFn, dependencies = []) {
  const cache = new WeakMap()
  const lastDeps = new Map()
  
  return computed(() => {
    const arr = unref(array)
    const deps = dependencies.map(dep => unref(dep))
    const depsKey = JSON.stringify(deps)
    
    // Check if dependencies changed
    if (lastDeps.get(arr) === depsKey && cache.has(arr)) {
      return cache.get(arr)
    }
    
    // Group the array
    const grouped = arr.reduce((groups, item) => {
      const key = keyFn(item)
      if (!groups[key]) {
        groups[key] = []
      }
      groups[key].push(item)
      return groups
    }, {})
    
    // Cache the result
    cache.set(arr, grouped)
    lastDeps.set(arr, depsKey)
    
    return grouped
  })
}

export default {
  memoize,
  clearMemoCache,
  memoizedComputed,
  memoizedArraySort,
  memoizedArrayFilter,
  memoizedGroupBy
}