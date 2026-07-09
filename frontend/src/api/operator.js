// System Command Center + maintenance-action API client.
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

// ── Backups ─────────────────────────────────────────────────────────

export const listBackups = () =>
  apiClient.get('/api/v1/admin/system/backups').then(r => r.data?.data ?? r.data)

export const getBackupHealth = () =>
  apiClient.get('/api/v1/admin/system/backups/health').then(r => r.data?.data ?? r.data)

export const queueBackupJob = (action, backupId, confirmation) =>
  apiClient.post('/api/v1/admin/system/backups/jobs', { action, backup_id: backupId, confirmation }).then(r => r.data?.data ?? r.data)

export const deleteBackup = (id) =>
  apiClient.delete(`/api/v1/admin/system/backups/${id}`).then(r => r.data)

export const downloadBackup = (id) =>
  apiClient.get(`/api/v1/admin/system/backups/${id}/download`, { responseType: 'blob' })

export const exportBackupJobs = () =>
  apiClient.get('/api/v1/admin/system/backups/jobs/export', { responseType: 'blob' })

export const clearBackupJobs = (confirmation) =>
  apiClient.delete('/api/v1/admin/system/backups/jobs', { data: { confirmation } }).then(r => r.data)

export const deleteBackupJob = (jobId) =>
  apiClient.delete(`/api/v1/admin/system/backups/jobs/${jobId}`).then(r => r.data)

export const downloadSchema = () =>
  apiClient.get('/api/v1/admin/system/backups/schema', { responseType: 'blob' })

export const importSchema = (file, confirmation) => {
  const fd = new FormData()
  fd.append('schema', file)
  fd.append('confirmation', confirmation)
  return apiClient.post('/api/v1/admin/system/backups/schema/import', fd, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }).then(r => r.data)
}

export const deleteSchema = (confirmation) =>
  apiClient.delete('/api/v1/admin/system/backups/schema', { data: { confirmation } }).then(r => r.data)
