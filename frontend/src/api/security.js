// Self-service security + activity log endpoints.
import apiClient from './client'

export const getMyActivities = (params = {}) => {
  const qs = new URLSearchParams(params).toString()
  return apiClient.get(`/api/v1/me/activities${qs ? '?' + qs : ''}`).then(r => r.data)
}

export const getMySecurity = () =>
  apiClient.get('/api/v1/me/security').then(r => r.data?.data ?? r.data)

export const addIpWhitelist = (ipOrCidr, label = '') =>
  apiClient.post('/api/v1/me/security/ip-whitelist', { ip_or_cidr: ipOrCidr, label }).then(r => r.data?.data ?? r.data)

export const removeIpWhitelist = (id) =>
  apiClient.delete(`/api/v1/me/security/ip-whitelist/${id}`).then(r => r.data)

export const enroll2FA = () =>
  apiClient.post('/api/v1/me/security/2fa/enroll').then(r => r.data?.data ?? r.data)

export const verify2FA = (code) =>
  apiClient.post('/api/v1/me/security/2fa/verify', { code }).then(r => r.data)

export const disable2FA = () =>
  apiClient.post('/api/v1/me/security/2fa/disable').then(r => r.data)

// Existing auth surface re-exported for convenience
export const changePassword = (currentPassword, newPassword) =>
  apiClient.post('/api/v1/auth/change-password', { current_password: currentPassword, new_password: newPassword }).then(r => r.data)

export const setPin = (pin) =>
  apiClient.post('/api/v1/auth/pin/set', { pin }).then(r => r.data)

export const changePin = (currentPin, newPin) =>
  apiClient.put('/api/v1/auth/pin/change', { current_pin: currentPin, new_pin: newPin }).then(r => r.data)

// Institutional departments (10-tier) with staff counts
export const listInstitutionalDepartments = () =>
  apiClient.get('/api/v1/cms/departments').then(r => r.data?.data ?? r.data)
