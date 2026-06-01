import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client.js'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const accessToken = ref(null)
  const refreshToken = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)

  const isAdmin = computed(() => {
    if (!user.value || !user.value.roles) return false
    return user.value.roles.some(r => {
      const name = typeof r === 'string' ? r : r.name
      return name === 'admin' || name === 'super_admin'
    })
  })

  const isSuperAdmin = computed(() => {
    if (!user.value || !user.value.roles) return false
    return user.value.roles.some(r => {
      const name = typeof r === 'string' ? r : r.name
      return name === 'super_admin'
    })
  })

  const isContentManager = computed(() => {
    if (!user.value || !user.value.roles) return false
    return user.value.roles.some(r => {
      const name = typeof r === 'string' ? r : r.name
      return name === 'content_manager' || name === 'admin' || name === 'super_admin'
    })
  })

  const isApplicant = computed(() => {
    if (!user.value) return false
    const roles = user.value.roles || []
    if (!roles.length) return true // no roles → treat as plain user
    return roles.some(r => {
      const name = typeof r === 'string' ? r : r.name
      return name === 'applicant' || name === 'user'
    })
  })

  async function refreshUser() {
    try {
      const res = await apiClient.get('/api/v1/auth/me')
      const fresh = res.data?.data || res.data
      if (fresh) {
        user.value = { ...user.value, ...fresh }
        persistToStorage()
      }
    } catch {
      // ignore — silent refresh
    }
  }

  function loadFromStorage() {
    const token = localStorage.getItem('ncsms_access_token')
    const refresh = localStorage.getItem('ncsms_refresh_token')
    const storedUser = localStorage.getItem('ncsms_user')

    if (token) accessToken.value = token
    if (refresh) refreshToken.value = refresh
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        user.value = null
      }
    }
  }

  function persistToStorage() {
    if (accessToken.value) {
      localStorage.setItem('ncsms_access_token', accessToken.value)
    } else {
      localStorage.removeItem('ncsms_access_token')
    }
    if (refreshToken.value) {
      localStorage.setItem('ncsms_refresh_token', refreshToken.value)
    } else {
      localStorage.removeItem('ncsms_refresh_token')
    }
    if (user.value) {
      localStorage.setItem('ncsms_user', JSON.stringify(user.value))
    } else {
      localStorage.removeItem('ncsms_user')
    }
  }

  async function register(firstName, lastName, email, password) {
    loading.value = true
    try {
      const response = await apiClient.post('/api/v1/auth/register', {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
      })
      const data = response.data.data
      // Auto-login if backend returns tokens
      if (data?.access_token) {
        accessToken.value = data.access_token
        refreshToken.value = data.refresh_token
        user.value = data.user
        persistToStorage()
      }
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error?.message || 'Registration failed. Please try again.'
      return { success: false, message }
    } finally {
      loading.value = false
    }
  }

  async function login(email, password) {
    loading.value = true
    try {
      const response = await apiClient.post('/api/v1/auth/login', { email, password })
      const data = response.data.data
      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      user.value = data.user
      persistToStorage()
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error?.message || 'Login failed. Please check your credentials.'
      return { success: false, message }
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    try {
      if (refreshToken.value) {
        await apiClient.post('/api/v1/auth/logout', { refresh_token: refreshToken.value })
      }
    } catch {
      // ignore logout errors
    } finally {
      accessToken.value = null
      refreshToken.value = null
      user.value = null
      persistToStorage()
      loading.value = false
    }
  }

  async function refreshAccessToken() {
    const storedRefresh = refreshToken.value || localStorage.getItem('ncsms_refresh_token')
    if (!storedRefresh) throw new Error('No refresh token')

    const response = await apiClient.post('/api/v1/auth/refresh', { refresh_token: storedRefresh })
    const data = response.data.data
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    if (data.user) user.value = data.user
    persistToStorage()
    return data.access_token
  }

  function clearAuth() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem('ncsms_access_token')
    localStorage.removeItem('ncsms_refresh_token')
    localStorage.removeItem('ncsms_user')
  }

  return {
    user,
    accessToken,
    refreshToken,
    loading,
    isAuthenticated,
    isAdmin,
    isSuperAdmin,
    isContentManager,
    isApplicant,
    register,
    login,
    logout,
    refreshAccessToken,
    refreshUser,
    loadFromStorage,
    clearAuth
  }
})
