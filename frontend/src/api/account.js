import apiClient from './client.js'

export function getProfile() {
  return apiClient.get('/api/v1/auth/me')
}

export function updatePassword(data) {
  return apiClient.post('/api/v1/auth/update-password', {
    current_password: data.current_password,
    new_password: data.new_password,
  })
}

export function listSessions() {
  return apiClient.get('/api/v1/auth/sessions')
}

export function revokeSession(id) {
  return apiClient.delete(`/api/v1/auth/sessions/${id}`)
}

export function revokeOtherSessions() {
  return apiClient.delete('/api/v1/auth/sessions/clear-all')
}

export function listMyAuditLogs(params = {}) {
  return apiClient.get('/api/v1/account/audit-logs', { params })
}
