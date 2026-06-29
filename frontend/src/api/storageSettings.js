import apiClient from './client.js'

export function getStorageSettings() {
  return apiClient.get('/api/v1/admin/storage-settings').then(r => r.data.data)
}

export function updateStorageSettings(payload) {
  return apiClient.put('/api/v1/admin/storage-settings', payload).then(r => r.data.data)
}
