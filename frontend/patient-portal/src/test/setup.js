/**
 * Vitest Test Setup
 * Global test configuration and utilities
 */

// Minimal crypto.getRandomValues fallback (before other imports)
if (!globalThis.crypto || !globalThis.crypto.getRandomValues) {
  try {
    const { webcrypto } = await import('crypto')
    globalThis.crypto = webcrypto
  } catch (e) {
    // Very small fallback (NOT cryptographically secure) just to satisfy libs in tests
    globalThis.crypto = {
      getRandomValues(arr) {
        for (let i = 0; i < arr.length; i++) {
          arr[i] = Math.floor(Math.random() * 256)
        }
        return arr
      }
    }
  }
}

import { config } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { vi } from 'vitest'

// Setup global test plugins
config.global.plugins = [createPinia()]

// Mock vue-toastification
const mockToast = {
  success: vi.fn(),
  error: vi.fn(),
  warning: vi.fn(),
  info: vi.fn()
}

vi.mock('vue-toastification', () => ({
  useToast: () => mockToast,
  POSITION: {
    TOP_RIGHT: 'top-right'
  }
}))

// Mock window.performance if not available
if (typeof window !== 'undefined') {
  if (!window.performance) {
    window.performance = {
      now: () => Date.now(),
      memory: {
        usedJSHeapSize: 1000000,
        totalJSHeapSize: 2000000,
        jsHeapSizeLimit: 4000000,
      },
    }
  }
}

// Mock ResizeObserver if not available
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock IntersectionObserver if not available
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock matchMedia if not available
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock

// Mock sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.sessionStorage = sessionStorageMock

// Mock fetch if needed
global.fetch = vi.fn()

// Mock console methods for cleaner test output
const consoleMock = {
  log: vi.fn(),
  error: vi.fn(),
  warn: vi.fn(),
  info: vi.fn(),
  debug: vi.fn(),
  group: vi.fn(),
  groupEnd: vi.fn(),
  groupCollapsed: vi.fn(),
}

// Only mock console in test environment
if (process.env.NODE_ENV === 'test') {
  global.console = consoleMock
}

// Crypto polyfill for environments lacking webcrypto (happy-dom)
// (Secondary crypto patch retained above; legacy section removed)