// Smart-updates API client — wraps the admin GitHub-release sentinel
// + Docker-Compose in-place deploy endpoints.
import apiClient from './client'

export const getUpdateStatus = () =>
  apiClient.get('/api/v1/admin/system/updates').then(r => r.data?.data ?? r.data)

export const checkForUpdates = () =>
  apiClient.post('/api/v1/admin/system/updates/check').then(r => r.data?.data ?? r.data)

export const getDeployStatus = () =>
  apiClient.get('/api/v1/admin/system/updates/deploy').then(r => r.data?.data ?? r.data)

export const startDeploy = () =>
  apiClient.post('/api/v1/admin/system/updates/deploy').then(r => r.data?.data ?? r.data)

export const startRollback = () =>
  apiClient.post('/api/v1/admin/system/updates/rollback').then(r => r.data?.data ?? r.data)
