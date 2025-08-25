// Polyfills for test environment
import { webcrypto } from 'crypto'
if (!globalThis.crypto) {
  globalThis.crypto = webcrypto
}
