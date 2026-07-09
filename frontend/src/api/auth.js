import apiClient from './client.js'

export function registerAccount(payload) {
  return apiClient.post('/api/v1/auth/register', payload)
}

export function loginAccount(payload) {
  return apiClient.post('/api/v1/auth/login', payload)
}

export function loginWithGoogleCredential(credential) {
  return apiClient.post('/api/v1/auth/google', { credential })
}

export function getCurrentUser() {
  return apiClient.get('/api/v1/auth/me')
}

export function updateMyProfile(data) {
  return apiClient.put('/api/v1/auth/me', data)
}
