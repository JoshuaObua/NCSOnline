import apiClient from './client'

export async function listClubAcademyLicenses(params = {}) {
  const res = await apiClient.get('/api/v1/admin/club-academy-licenses', { params })
  return res.data
}

export async function getClubAcademyLicenseKPIs() {
  const res = await apiClient.get('/api/v1/admin/club-academy-licenses/kpis')
  return res.data?.data || res.data
}

export async function getClubAcademyLicenseByID(id) {
  const res = await apiClient.get(`/api/v1/admin/club-academy-licenses/${id}`)
  return res.data?.data || res.data
}

export async function createClubAcademyLicense(payload) {
  const res = await apiClient.post('/api/v1/admin/club-academy-licenses', payload)
  return res.data?.data || res.data
}

export async function extendClubAcademyLicense(id, payload) {
  const res = await apiClient.post(`/api/v1/admin/club-academy-licenses/${id}/extend`, payload)
  return res.data?.data || res.data
}

export async function revokeClubAcademyLicense(id, payload) {
  const res = await apiClient.post(`/api/v1/admin/club-academy-licenses/${id}/revoke`, payload)
  return res.data?.data || res.data
}

export async function reinstateClubAcademyLicense(id, payload) {
  const res = await apiClient.post(`/api/v1/admin/club-academy-licenses/${id}/reinstate`, payload)
  return res.data?.data || res.data
}

export async function getClubActiveLicense(clubId) {
  const res = await apiClient.get(`/api/v1/clubs/${clubId}/license`)
  return res.data?.data || res.data
}
