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

export function adminVerifyApplicationPayment(id) {
  return apiClient.post(`/api/v1/admin/applications/${id}/verify-payment`)
}

export function adminRejectApplicationPayment(id, notes = '') {
  return apiClient.post(`/api/v1/admin/applications/${id}/reject-payment`, { notes })
}

export function listMyLegacyApplications(params = {}) {
  return apiClient.get('/api/v1/applications', { params })
}

export function getMyLegacyApplication(id) {
  return apiClient.get(`/api/v1/applications/${id}`)
}

export function listMyTransactions(params = {}) {
  return apiClient.get('/api/v1/transactions', { params })
}
