import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { post, get } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const refreshToken = ref(localStorage.getItem('refreshToken') || null)
  const loading = ref(false)
  const error = ref(null)

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || null)
  const userName = computed(() => user.value?.full_name || 'Pengguna')
  const userEmail = computed(() => user.value?.email || '')
  const userId = computed(() => user.value?.id || null)

  const isPwaUser = computed(() => {
    const role = userRole.value
    return role === 'student' || role === 'dudi' || role === 'teacher'
  })

  const isAdmin = computed(() => {
    return userRole.value === 'admin' || userRole.value === 'admin_jurusan'
  })

  const userJurusan = computed(() => user.value?.jurusan || '')

  async function login(credentials) {
    loading.value = true
    error.value = null
    try {
      const data = await post('/login', credentials)
      setAuth(data)
      return data
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function register(userData) {
    loading.value = true
    error.value = null
    try {
      const data = await post('/register', userData)
      return data
    } catch (e) {
      error.value = e.message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchMe() {
    try {
      const data = await get('/me')
      user.value = data
      localStorage.setItem('user', JSON.stringify(data))
    } catch {
      logout()
    }
  }

  function setAuth(data) {
    token.value = data.access_token
    refreshToken.value = data.refresh_token
    user.value = data.user
    localStorage.setItem('token', data.access_token)
    localStorage.setItem('refreshToken', data.refresh_token)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  function restoreSession() {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    if (savedToken && savedUser) {
      token.value = savedToken
      refreshToken.value = localStorage.getItem('refreshToken')
      try {
        user.value = JSON.parse(savedUser)
      } catch {
        user.value = null
        logout()
      }
    }
  }

  function logout() {
    token.value = null
    refreshToken.value = null
    user.value = null
    error.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  return {
    user, token, refreshToken, loading, error,
    isAuthenticated, userRole, userName, userEmail, userId, isPwaUser, isAdmin, userJurusan,
    login, register, fetchMe, setAuth, restoreSession, logout
  }
})
