import apiClient from './client.js'

export function listNsmisDomain(resource, params = {}) {
  return apiClient.get(`/api/v1/nsmis/${resource}`, { params })
}

export function createNsmisDomain(resource, data) {
  return apiClient.post(`/api/v1/nsmis/${resource}`, data)
}

export function updateNsmisDomain(resource, id, data) {
  return apiClient.put(`/api/v1/nsmis/${resource}/${id}`, data)
}

export function deleteNsmisDomain(resource, id) {
  return apiClient.delete(`/api/v1/nsmis/${resource}/${id}`)
}

export function getAthleteDashboard() {
  return apiClient.get('/api/v1/nsmis/dashboards/athletes')
}

export function getPerformanceDashboard() {
  return apiClient.get('/api/v1/nsmis/dashboards/performance')
}

export function getFinanceDashboard() {
  return apiClient.get('/api/v1/nsmis/dashboards/finance')
}

export function getTalentDashboard() {
  return apiClient.get('/api/v1/nsmis/dashboards/talent')
}
