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

  const roleNames = computed(() => (user.value?.roles || []).map(r => typeof r === 'string' ? r : r.name))
  const hasAnyRole = (...names) => names.some(name => roleNames.value.includes(name))
  const canUseNSMIS = computed(() => hasAnyRole(
    'super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department',
    'finance_department', 'federation_president', 'federation_general_secretary',
    'safeguarding_officer', 'auditor'
  ))

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
    const storedUser = localStorage.getItem('ncsms_user')

    if (token) accessToken.value = token
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
	localStorage.removeItem('ncsms_refresh_token')
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
		refreshToken.value = null
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
	  refreshToken.value = null
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
      clearAuth()
      await apiClient.post('/api/v1/auth/logout', {}).catch(() => {})
    } catch {
      // ignore logout errors
    } finally {
      clearAuth()
      loading.value = false
    }
  }

  async function refreshAccessToken() {
	const response = await apiClient.post('/api/v1/auth/refresh', {})
    const data = response.data.data
    accessToken.value = data.access_token
	refreshToken.value = null
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
    roleNames,
    hasAnyRole,
    canUseNSMIS,
    register,
    login,
    logout,
    refreshAccessToken,
    refreshUser,
    loadFromStorage,
    clearAuth
  }
})
