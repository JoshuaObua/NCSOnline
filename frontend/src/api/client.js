import axios from 'axios'

function resolveApiBase() {
  const envBase = import.meta.env?.VITE_API_BASE_URL
  if (envBase) return envBase.replace(/\/$/, '')
  return ''
}

export const API_BASE_URL = resolveApiBase()

export function mediaUrl(path) {
  if (!path) return ''
  if (/^(https?:|data:|blob:)/i.test(path)) return path
  const normalized = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE_URL}${normalized}`
}

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json'
  }
})

let isRefreshing = false
let failedQueue = []

function processQueue(error, token = null) {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  failedQueue = []
}

apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('ncsms_access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !originalRequest.url?.includes('/auth/login') &&
      !originalRequest.url?.includes('/auth/refresh') &&
      !originalRequest.url?.includes('/auth/logout')
    ) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then(token => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            return apiClient(originalRequest)
          })
          .catch(err => Promise.reject(err))
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const response = await axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {}, { withCredentials: true })
        const { access_token } = response.data.data

        localStorage.setItem('ncsms_access_token', access_token)

        apiClient.defaults.headers.common.Authorization = `Bearer ${access_token}`
        originalRequest.headers.Authorization = `Bearer ${access_token}`

        processQueue(null, access_token)
        isRefreshing = false

        return apiClient(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        isRefreshing = false
        clearAuthAndRedirect()
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

function clearAuthAndRedirect() {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

export default apiClient

export const apiGet = (url, params = {}) => apiClient.get(`/api/v1${url}`, { params })
export const apiPost = (url, data = {}) => apiClient.post(`/api/v1${url}`, data)
export const apiPut = (url, data = {}) => apiClient.put(`/api/v1${url}`, data)
export const apiDelete = (url) => apiClient.delete(`/api/v1${url}`)
