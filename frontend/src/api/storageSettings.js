import apiClient from './client.js'

export function getStorageSettings() {
  return apiClient.get('/api/v1/admin/storage-settings').then(r => r.data.data)
}

export function updateStorageSettings(payload) {
  return apiClient.put('/api/v1/admin/storage-settings', payload).then(r => r.data.data)
}

export function connectGoogleDrive(redirectUri) {
  return apiClient
    .post('/api/v1/admin/storage-settings/google-drive/connect', { redirect_uri: redirectUri })
    .then(r => r.data.data)
}

export function exchangeGoogleDriveCode({ code, state, redirectUri }) {
  return apiClient
    .post('/api/v1/admin/storage-settings/google-drive/exchange', { code, state, redirect_uri: redirectUri })
    .then(r => r.data.data)
}

export function disconnectGoogleDrive() {
  return apiClient.post('/api/v1/admin/storage-settings/google-drive/disconnect').then(r => r.data.data)
}
