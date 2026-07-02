import apiClient from './client.js'

export function loginWithGoogleCredential(credential) {
  return apiClient.post('/api/v1/auth/google', { credential })
}

export function getCurrentUser() {
  return apiClient.get('/api/v1/auth/me')
}
