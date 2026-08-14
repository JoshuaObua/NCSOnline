import axios from 'axios'

// Resolve the API base URL at runtime so the same build works on any host.
// Priority: VITE_API_BASE_URL env var -> same host as the browser on port 9081.
function resolveApiBase() {
  const envBase = import.meta.env?.VITE_API_BASE_URL
  if (envBase) return envBase.replace(/\/$/, '')
  if (typeof window !== 'undefined' && window.location?.hostname) {
    return `${window.location.protocol}//${window.location.hostname}:9081`
  }
  return 'http://localhost:9081'
}

export const API_BASE_URL = resolveApiBase()

/**
 * Resolve a possibly-relative media path (e.g. "/uploads/images/x.jpg") to an
 * absolute URL the browser can fetch. nginx serves uploads at the API origin
 * (port 9081), but the frontend is hosted on port 3001 during dev, so a bare
 * "/uploads/..." path would 404 against the Vite host.
 *
 * Accepts:
 *   - "" / null / undefined         → returns ""
 *   - "http(s)://..." absolute URLs → returned unchanged
 *   - "data:..." or "blob:..."      → returned unchanged
 *   - any other path starting with "/" → prefixed with API_BASE_URL
 *   - bare path (no leading slash)  → "/" + path, then prefixed
 */
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

// Request interceptor — attach auth token
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

// Response interceptor — handle 401 and auto-refresh
apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    // Skip refresh logic for auth endpoints themselves
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
  // Use window.location for hard redirect outside of Vue context
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

export default apiClient

// Named convenience helpers used by departmental dashboards
export const apiGet = (url, params = {}) => apiClient.get(`/api/v1${url}`, { params })
export const apiPost = (url, data = {}) => apiClient.post(`/api/v1${url}`, data)
export const apiPut = (url, data = {}) => apiClient.put(`/api/v1${url}`, data)
export const apiDelete = (url) => apiClient.delete(`/api/v1${url}`)
