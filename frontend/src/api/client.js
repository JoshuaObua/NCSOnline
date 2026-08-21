import axios from 'axios'

// Resolve the API base URL at runtime so the same build works on any host.
// Priority: explicit base URL, then the browser host plus VITE_API_PORT.
function resolveApiBase() {
  const envBase = import.meta.env?.VITE_API_BASE_URL
  if (envBase) return envBase.replace(/\/$/, '')
  if (typeof window !== 'undefined' && window.location?.hostname) {
    const isLocal = ['localhost', '127.0.0.1', '0.0.0.0'].includes(window.location.hostname)
    const port = import.meta.env?.VITE_API_PORT
    if (isLocal && port) {
      return `${window.location.protocol}//${window.location.hostname}:${port}`
    }
    return window.location.origin
  }
  return ''
}

export const API_BASE_URL = resolveApiBase()

/**
 * Resolve a possibly-relative media path (e.g. "/uploads/images/x.jpg") to an
 * absolute URL the browser can fetch. nginx serves uploads at the API origin,
 * but the frontend can be hosted on a separate dev origin, so a bare
 * "/uploads/..." path may 404 against the Vite host.
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

const AUTH_ENDPOINTS = ['/auth/login', '/auth/google', '/auth/refresh', '/auth/logout']
const isAuthEndpoint = url => AUTH_ENDPOINTS.some(path => url?.includes(path))
const hasAccessToken = () => Boolean(localStorage.getItem('ncsms_access_token'))
const isPublicCmsRequest = config =>
  String(config?.method || 'get').toLowerCase() === 'get' &&
  config?.url?.includes('/api/v1/cms/')
const skipsAuth = config => Boolean(config?.skipAuth) || isPublicCmsRequest(config)

function decodeJwtExpMs(token) {
  try {
    const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    const payload = JSON.parse(atob(base64))
    return typeof payload.exp === 'number' ? payload.exp * 1000 : null
  } catch {
    return null
  }
}

// Single-flight refresh: concurrent callers all await the same in-flight
// request instead of each firing their own POST /auth/refresh.
let refreshPromise = null
function refreshAccessToken() {
  if (!refreshPromise) {
    refreshPromise = axios.post(`${API_BASE_URL}/api/v1/auth/refresh`, {}, { withCredentials: true })
      .then(res => {
        const { access_token } = res.data.data
        localStorage.setItem('ncsms_access_token', access_token)
        apiClient.defaults.headers.common.Authorization = `Bearer ${access_token}`
        return access_token
      })
      .finally(() => { refreshPromise = null })
  }
  return refreshPromise
}

// Request interceptor — attach the auth token, refreshing it first if it's
// within 30s of expiring. Access tokens are short-lived (15m), so without
// this every CMS action after that window would fire with a stale token,
// hit a 401, and only then refresh-and-retry — visible in the browser
// console as a failed request even though the retry silently succeeds.
apiClient.interceptors.request.use(
  async config => {
    let token = localStorage.getItem('ncsms_access_token')
    if (token && token !== 'local-cms-preview-token' && !isAuthEndpoint(config.url) && !skipsAuth(config)) {
      const expMs = decodeJwtExpMs(token)
      if (expMs && expMs - Date.now() < 30_000) {
        try { token = await refreshAccessToken() } catch { /* let the request go; the response interceptor will handle the 401 */ }
      }
    }
    if (token && !skipsAuth(config)) config.headers.Authorization = `Bearer ${token}`
    return config
  },
  error => Promise.reject(error)
)

// Response interceptor — fallback for the rare case a request still hits a
// 401 (e.g. clock skew, or a request already in flight when the token
// expired): refresh once and retry, sharing the same single-flight refresh.
apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config
    if (error.response?.status === 401 && hasAccessToken() && !skipsAuth(originalRequest) && !originalRequest._retry && !isAuthEndpoint(originalRequest.url)) {
      originalRequest._retry = true
      try {
        const token = await refreshAccessToken()
        originalRequest.headers.Authorization = `Bearer ${token}`
        return apiClient(originalRequest)
      } catch (refreshError) {
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
  if (window.location.pathname !== '/cms/login') {
    window.location.href = '/cms/login'
  }
}

export default apiClient
