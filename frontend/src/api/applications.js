import apiClient from './client.js'

export function getAdminDashboard() {
  return apiClient.get('/api/v1/admin/dashboard')
}

export function adminListApplications(params = {}) {
  return apiClient.get('/api/v1/admin/applications', { params })
}

export function adminGetApplication(id) {
  return apiClient.get(`/api/v1/admin/applications/${id}`)
}

export function adminReviewApplication(id, action, notes = '') {
  return apiClient.post(`/api/v1/admin/applications/${id}/${action}`, { notes })
}

export function listMyLegacyApplications(params = {}) {
  return apiClient.get('/api/v1/applications', { params })
}

export function listMyTransactions(params = {}) {
  return apiClient.get('/api/v1/transactions', { params })
}
