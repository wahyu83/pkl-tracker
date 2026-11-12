const BASE_URL = '/api'

// Endpoint yang tidak perlu (dan tidak boleh) memakai refresh token.
const NO_REFRESH_ENDPOINTS = ['/login', '/register', '/refresh']

let refreshPromise = null

function authHeader() {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function storeSession(data) {
  if (data.access_token) localStorage.setItem('token', data.access_token)
  if (data.refresh_token) localStorage.setItem('refreshToken', data.refresh_token)
  if (data.user) localStorage.setItem('user', JSON.stringify(data.user))
}

function clearSession() {
  localStorage.removeItem('token')
  localStorage.removeItem('refreshToken')
  localStorage.removeItem('user')
}

function redirectToLogin() {
  clearSession()
  if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

async function performRefresh() {
  const refreshToken = localStorage.getItem('refreshToken')
  if (!refreshToken) return false
  try {
    const res = await fetch(`${BASE_URL}/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    })
    if (!res.ok) return false
    const data = await res.json()
    if (!data.access_token) return false
    storeSession(data)
    return true
  } catch {
    return false
  }
}

// Pastikan hanya satu refresh berjalan bersamaan (hindari request paralel beruntun).
function refreshSession() {
  if (!refreshPromise) {
    refreshPromise = performRefresh().finally(() => { refreshPromise = null })
  }
  return refreshPromise
}

function canRefresh(endpoint) {
  return !NO_REFRESH_ENDPOINTS.some((e) => endpoint.startsWith(e))
}

async function send(endpoint, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...authHeader(),
    ...options.headers
  }
  return fetch(`${BASE_URL}${endpoint}`, { ...options, headers })
}

async function parse(res) {
  let data = {}
  try {
    data = await res.json()
  } catch {
    data = {}
  }
  if (!res.ok) {
    const err = new Error(data.error || 'Request failed')
    err.status = res.status
    throw err
  }
  return data
}

async function request(endpoint, options = {}) {
  const res = await send(endpoint, options)

  // Access token expired -> coba refresh sekali, lalu ulangi request.
  if (res.status === 401 && canRefresh(endpoint)) {
    const ok = await refreshSession()
    if (ok) {
      return parse(await send(endpoint, options))
    }
    redirectToLogin()
    const err = new Error('Sesi berakhir, silakan login kembali')
    err.status = 401
    throw err
  }

  return parse(res)
}

export function get(endpoint) {
  return request(endpoint)
}

export function post(endpoint, body) {
  return request(endpoint, {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function put(endpoint, body) {
  return request(endpoint, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function del(endpoint) {
  return request(endpoint, { method: 'DELETE' })
}

// Unggah file (multipart). Content-Type sengaja tidak diset agar browser
// menuliskan boundary multipart sendiri.
export async function upload(endpoint, formData) {
  const sendForm = () => fetch(`${BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: { ...authHeader() },
    body: formData
  })

  let res = await sendForm()

  if (res.status === 401 && canRefresh(endpoint)) {
    const ok = await refreshSession()
    if (ok) {
      res = await sendForm()
    } else {
      redirectToLogin()
      const err = new Error('Sesi berakhir, silakan login kembali')
      err.status = 401
      throw err
    }
  }

  return parse(res)
}

async function fetchFile(endpoint) {
  let res = await fetch(`${BASE_URL}${endpoint}`, { headers: authHeader() })

  if (res.status === 401 && canRefresh(endpoint)) {
    const ok = await refreshSession()
    if (ok) {
      res = await fetch(`${BASE_URL}${endpoint}`, { headers: authHeader() })
    }
  }

  if (res.status === 401) {
    redirectToLogin()
    throw new Error('Sesi berakhir, silakan login kembali')
  }
  if (!res.ok) throw new Error('Download failed')
  return res.blob()
}

export async function downloadCsv(endpoint, filename) {
  const blob = await fetchFile(endpoint)
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

export async function downloadFile(endpoint, filename) {
  return downloadCsv(endpoint, filename)
}

export { refreshSession }
export default { get, post, put, del, upload, downloadCsv, downloadFile, refreshSession }
