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
