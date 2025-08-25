import { createNamespaceClient } from './utils.js'

// Namespaced client for debug endpoints
const debugClient = createNamespaceClient('/debug')

export const debugApi = {
  // Returns whether debug features are enabled (currently always true server-side)
  async isEnabled() {
    const { data } = await debugClient.get('/enabled')
    return !!data?.enabled
  },
}
