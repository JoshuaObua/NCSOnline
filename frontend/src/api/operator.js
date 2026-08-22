import apiClient from './client.js'

export const getSystemStatus = () =>
  apiClient.get('/api/v1/admin/system/status').then(r => r.data?.data ?? r.data)

export const getSystemResources = () =>
  apiClient.get('/api/v1/admin/system/resources').then(r => r.data?.data ?? r.data)

export const getServiceLogs = (service, lines = 150) =>
  apiClient.get('/api/v1/admin/system/service-logs', { params: { service, lines } }).then(r => r.data?.data ?? r.data)

export const runServiceAction = (service, action, confirmation) =>
  apiClient.post('/api/v1/admin/system/services/action', { service, action, confirmation }).then(r => r.data?.data ?? r.data)

export const runMaintenanceAction = (action, confirmation) =>
  apiClient.post('/api/v1/admin/system/actions', { action, confirmation }).then(r => r.data?.data ?? r.data)

export const flushCache = () =>
  apiClient.post('/api/v1/admin/system/cache/flush').then(r => r.data?.data ?? r.data)

export const setMaintenance = (payload) =>
  apiClient.put('/api/v1/admin/system/maintenance', payload).then(r => r.data?.data ?? r.data)

// ── Backup management ────────────────────────────────────────────────
export const listBackups = () =>
  apiClient.get('/api/v1/admin/backups').then(r => r.data?.data ?? r.data)

export const queueBackupJob = (payload) =>
  apiClient.post('/api/v1/admin/backups/queue', payload).then(r => r.data?.data ?? r.data)

export const deleteBackup = (id) =>
  apiClient.delete(`/api/v1/admin/backups/${id}`).then(r => r.data?.data ?? r.data)

export const downloadBackup = (id) =>
  apiClient.get(`/api/v1/admin/backups/${id}/download`, { responseType: 'blob' })

export const exportBackupJobs = () =>
  apiClient.get('/api/v1/admin/backups/jobs/export', { responseType: 'blob' })

export const clearBackupJobs = () =>
  apiClient.delete('/api/v1/admin/backups/jobs').then(r => r.data?.data ?? r.data)

export const getBackupHealth = () =>
  apiClient.get('/api/v1/admin/backups/health').then(r => r.data?.data ?? r.data)

// ── Schema management ────────────────────────────────────────────────
export const importSchema = (formData) =>
  apiClient.post('/api/v1/admin/schema/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } }).then(r => r.data?.data ?? r.data)

export const deleteSchema = (confirmation) =>
  apiClient.delete('/api/v1/admin/schema', { data: { confirmation } }).then(r => r.data?.data ?? r.data)

export const downloadSchema = () =>
  apiClient.get('/api/v1/admin/schema/export', { responseType: 'blob' })
