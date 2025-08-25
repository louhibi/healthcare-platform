import api from './index.js'

/**
 * Create a lightweight namespaced API client that automatically prefixes
 * all request paths with the provided basePath. Returns a set of verb
 * helpers plus the raw base path (base) should you need it.
 *
 * Usage:
 *   const entities = createNamespaceClient('/api/entities')
 *   const { data } = await entities.get('/123')
 */
export function createNamespaceClient(basePath, client = api) {
  const join = (p = '') => `${basePath}${p}`
  return {
    base: basePath,
    get: (p, config) => client.get(join(p), config),
    delete: (p, config) => client.delete(join(p), config),
    head: (p, config) => client.head(join(p), config),
    options: (p, config) => client.options(join(p), config),
    post: (p, data, config) => client.post(join(p), data, config),
    put: (p, data, config) => client.put(join(p), data, config),
    patch: (p, data, config) => client.patch(join(p), data, config),
  }
}
