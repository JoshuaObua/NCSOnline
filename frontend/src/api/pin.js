import apiClient from './client.js'

export function setPin(pin) {
  return apiClient.post('/api/v1/auth/pin/set', { pin })
}

export function changePin(current_pin, new_pin) {
  return apiClient.put('/api/v1/auth/pin/change', { current_pin, new_pin })
}

export function verifyPin(pin) {
  return apiClient.post('/api/v1/auth/pin/verify', { pin })
}
