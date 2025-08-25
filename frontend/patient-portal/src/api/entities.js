import { createNamespaceClient } from './utils.js'

// Namespaced axios helper for entities endpoints. Changing the base path only happens here.
const entities = createNamespaceClient('/api/entities')

// Entities related API calls (healthcare entities, settings, etc.)
export const entitiesApi = {
  // Fetch a single healthcare entity by id
  async getEntity(entityId) {
    if (!entityId) throw new Error('entityId is required')
    const { data } = await entities.get(`/${entityId}`)
    return data
  },
  // Room requirement for a healthcare entity
  async getEntityRoomRequirement(entityId) {
    const { data } = await entities.get(`/${entityId}/room-requirement`)
    return data
  },
}
