// Crypto polyfill module to ensure getRandomValues exists on imported 'crypto'
import * as nodeCrypto from 'crypto'

if (!nodeCrypto.getRandomValues) {
  if (nodeCrypto.webcrypto && nodeCrypto.webcrypto.getRandomValues) {
    nodeCrypto.getRandomValues = (typedArray) => nodeCrypto.webcrypto.getRandomValues(typedArray)
  } else {
    nodeCrypto.getRandomValues = (typedArray) => {
      for (let i = 0; i < typedArray.length; i++) {
        typedArray[i] = Math.floor(Math.random() * 256)
      }
      return typedArray
    }
  }
}

export * from 'crypto'
export default nodeCrypto
