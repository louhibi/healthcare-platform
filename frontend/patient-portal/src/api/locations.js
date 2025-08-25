import api from './index'

export const locationsApi = {
  async getCountries() {
    const res = await api.get('/api/locations/countries')
    return res.data.data
  },
  async getStatesByCountryId(countryId, params = {}) {
    const res = await api.get(`/api/locations/countries/${countryId}/states`, { params })
    return res.data.data
  },
  async getCitiesByCountryId(countryId, params = {}) {
    const res = await api.get(`/api/locations/countries/${countryId}/cities`, { params })
    return res.data.data
  },
  async getCitiesByStateId(stateId, params = {}) {
    const res = await api.get(`/api/locations/states/${stateId}/cities`, { params })
    return res.data.data
  },
  async getNationalities() {
    const res = await api.get('/api/locations/nationalities')
    return res.data.data
  },
  async getNationalitiesByCountryId(countryId, params = {}) {
    const res = await api.get(`/api/locations/countries/${countryId}/nationalities`, { params })
    return res.data.data
  },
  async getInsuranceTypesByCountry(countryCode, params = {}) {
    const res = await api.get(`/api/locations/countries/${countryCode}/insurance-types`, { params })
    return res.data.data
  },
  async getInsuranceProvidersByType(typeId, params = {}) {
    const res = await api.get(`/api/locations/insurance-types/${typeId}/providers`, { params })
    return res.data.data
  }
}
