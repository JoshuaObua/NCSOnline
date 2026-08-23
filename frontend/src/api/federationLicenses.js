import apiClient from './client'

export const listFederationLicenses = (params = {}) => {
  const qs = new URLSearchParams(params).toString()
  return apiClient.get(`/api/v1/admin/federation-licenses${qs ? '?' + qs : ''}`).then(r => r.data)
}

export const getFederationLicense = (id) =>
  apiClient.get(`/api/v1/admin/federation-licenses/${id}`).then(r => r.data?.data ?? r.data)

export const createFederationLicense = (data) =>
  apiClient.post('/api/v1/admin/federation-licenses', data).then(r => r.data?.data ?? r.data)

export const extendFederationLicense = (id, data) =>
  apiClient.post(`/api/v1/admin/federation-licenses/${id}/extend`, data).then(r => r.data?.data ?? r.data)

export const revokeFederationLicense = (id, data) =>
  apiClient.post(`/api/v1/admin/federation-licenses/${id}/revoke`, data).then(r => r.data?.data ?? r.data)

export const reinstateFederationLicense = (id, data) =>
  apiClient.post(`/api/v1/admin/federation-licenses/${id}/reinstate`, data).then(r => r.data?.data ?? r.data)

export const getFederationLicenseKPIs = () =>
  apiClient.get('/api/v1/admin/federation-licenses/kpis').then(r => r.data?.data ?? r.data)

export const getFederationActiveLicense = (federationId) =>
  apiClient.get(`/api/v1/federations/${federationId}/license`).then(r => r.data?.data ?? r.data)
