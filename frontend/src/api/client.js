import axios from 'axios'

const apiClient = axios.create({
  baseURL: 'http://localhost:9080',
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

      const refreshToken = localStorage.getItem('ncsms_refresh_token')
      if (!refreshToken) {
        isRefreshing = false
        clearAuthAndRedirect()
        return Promise.reject(error)
      }

      try {
        const response = await axios.post('http://localhost:9080/api/v1/auth/refresh', {
          refresh_token: refreshToken
        })
        const { access_token, refresh_token: new_refresh } = response.data.data

        localStorage.setItem('ncsms_access_token', access_token)
        if (new_refresh) localStorage.setItem('ncsms_refresh_token', new_refresh)

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
  localStorage.removeItem('ncsms_refresh_token')
  localStorage.removeItem('ncsms_user')
  // Use window.location for hard redirect outside of Vue context
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

export default apiClient
