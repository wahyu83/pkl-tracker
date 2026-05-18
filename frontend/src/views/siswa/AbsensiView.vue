<template>
  <div>
    <!-- Tab Switch -->
    <div class="flex bg-gray-100 rounded-xl p-1 mb-5">
      <button
        v-for="t in tabs"
        :key="t.key"
        @click="activeTab = t.key"
        :class="[
          'flex-1 py-2 rounded-lg text-sm font-medium transition-all',
          activeTab === t.key ? 'bg-white text-primary shadow-sm' : 'text-gray-500'
        ]"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- Tab: Absen -->
    <div v-if="activeTab === 'absen'">
      <!-- Status Card -->
      <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
        <h3 class="font-semibold text-gray-800 mb-1">Status Absensi Hari Ini</h3>
        <p class="text-xs text-gray-500 mb-3">Server: {{ serverTime }}</p>

        <!-- Sudah masuk + pulang -->
        <div v-if="status.has_masuk && status.has_pulang" class="bg-accent/10 rounded-xl p-4 text-center">
          <CheckCircleIcon :size="40" class="mx-auto mb-2 text-accent" />
          <p class="text-sm font-medium text-accent">Absensi Hari Ini Selesai</p>
          <p class="text-xs text-accent/70 mt-1">
            Masuk: {{ formatTime(status.masuk_at) }} | Pulang: {{ formatTime(status.pulang_at) }}
          </p>
        </div>

        <!-- Sudah masuk, pulang belum tersedia (menunggu 7 jam) -->
        <div v-else-if="status.has_masuk && !status.pulang_available" class="bg-warning/10 rounded-xl p-4 text-center">
          <ClockIcon :size="40" class="mx-auto mb-2 text-warning" />
          <p class="text-sm font-medium text-warning">Menunggu Waktu Absen Pulang</p>
          <p class="text-xs text-warning/70 mt-1">Masuk: {{ formatTime(status.masuk_at) }}</p>
          <p class="text-xs text-warning/70">Pulang tersedia: {{ formatTime(status.pulang_available_at) }}</p>
          <p class="text-lg font-bold text-warning mt-2">{{ countdownText }}</p>
        </div>

        <!-- Sudah masuk, pulang sudah tersedia -->
        <div v-else-if="status.has_masuk && status.pulang_available && !status.has_pulang" class="bg-info/10 rounded-xl p-4 text-center">
          <AlertCircleIcon :size="40" class="mx-auto mb-2 text-info" />
          <p class="text-sm font-medium text-info">Anda sudah absen masuk. Silakan absen pulang.</p>
          <p class="text-xs text-info/70 mt-1">Masuk: {{ formatTime(status.masuk_at) }}</p>
          <button
            @click="startPulang"
            class="mt-3 px-6 py-2.5 bg-info text-white rounded-xl text-sm font-bold hover:bg-info/80 transition-colors"
          >
            Absen Pulang
          </button>
        </div>

        <!-- Belum masuk -->
        <div v-else class="bg-primary/5 rounded-xl p-4 text-center">
          <MapPinIcon :size="40" class="mx-auto mb-2 text-primary" />
          <p class="text-sm font-medium text-primary">Anda belum absen masuk hari ini</p>
          <button
            @click="startMasuk"
            class="mt-3 px-6 py-2.5 bg-primary text-white rounded-xl text-sm font-bold hover:bg-primary-light transition-colors"
          >
            Absen Masuk
          </button>
        </div>
      </div>

      <!-- Form Absen (masuk/pulang) -->
      <div v-if="absenType" class="space-y-4">
        <div class="bg-white rounded-2xl p-5 border border-gray-100">
          <h3 class="font-semibold text-gray-800 mb-1">
            {{ absenType === 'masuk' ? 'Absen Masuk' : 'Absen Pulang' }}
          </h3>
          <p class="text-xs text-gray-500 mb-4">Ambil lokasi GPS Anda</p>

          <div class="bg-gray-50 rounded-xl p-4 space-y-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0">
                <MapPinIcon :size="20" class="text-primary" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-xs text-gray-500">Lokasi Anda</p>
                <p v-if="location" class="text-sm font-medium text-gray-800 font-mono">
                  {{ location.lat.toFixed(6) }}, {{ location.lng.toFixed(6) }}
                </p>
                <p v-else class="text-sm text-gray-400">Mengambil lokasi...</p>
              </div>
              <button
                @click="getLocation"
                :disabled="locating"
                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-primary text-white hover:bg-primary-light disabled:opacity-50 transition-colors"
              >
                {{ locating ? 'Mencari...' : 'Deteksi' }}
              </button>
            </div>

            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-accent/10 flex items-center justify-center flex-shrink-0">
                <ClockIcon :size="20" class="text-accent" />
              </div>
              <div>
                <p class="text-xs text-gray-500">Waktu Server</p>
                <p class="text-sm font-medium text-gray-800 font-mono">{{ serverTime }}</p>
              </div>
            </div>
          </div>

          <!-- Map -->
          <div ref="mapContainer" class="mt-4 h-52 rounded-xl overflow-hidden border border-gray-200"></div>
        </div>

        <button
          @click="submitAbsensi"
          :disabled="!location || submitting"
          class="w-full py-3 bg-accent text-white rounded-2xl text-sm font-bold hover:bg-accent-dark transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2 shadow-lg shadow-accent/20"
        >
          <LoaderIcon v-if="submitting" :size="20" class="animate-spin" />
          <span>{{ submitting ? 'Menyimpan...' : (absenType === 'masuk' ? 'Simpan Absen Masuk' : 'Simpan Absen Pulang') }}</span>
        </button>

        <button
          @click="absenType = null; location = null; destroyMap()"
          class="w-full py-2.5 text-gray-500 text-sm font-medium"
        >
          Batal
        </button>
      </div>

      <!-- Success Banner -->
      <div v-if="successMsg" class="mt-4 bg-accent/10 rounded-xl p-3 flex items-center gap-2 text-sm text-accent font-medium">
        <CheckCircleIcon :size="18" />
        {{ successMsg }}
      </div>

      <!-- Error -->
      <div v-if="errorMsg" class="mt-4 bg-danger/10 rounded-xl p-3 flex items-center gap-2 text-sm text-danger font-medium">
        <XIcon :size="18" />
        {{ errorMsg }}
      </div>
    </div>

    <!-- Tab: Riwayat -->
    <div v-if="activeTab === 'riwayat'">
      <!-- Month selector -->
      <div class="flex items-center gap-2 mb-4">
        <button @click="prevMonth" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
          <ChevronLeftIcon :size="18" />
        </button>
        <span class="flex-1 text-center text-sm font-semibold text-gray-800">{{ currentMonth }}</span>
        <button @click="nextMonth" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
          <ChevronRightIcon :size="18" />
        </button>
      </div>

      <!-- Stats -->
      <div class="bg-white rounded-2xl p-4 border border-gray-100 mb-4">
        <div class="grid grid-cols-3 gap-2 text-center">
          <div>
            <p class="text-lg font-bold text-accent">{{ filteredHistory.length }}</p>
            <p class="text-[10px] text-gray-500">Total Absen</p>
          </div>
          <div>
            <p class="text-lg font-bold text-primary">{{ summary.masuk }}</p>
            <p class="text-[10px] text-gray-500">Masuk</p>
          </div>
          <div>
            <p class="text-lg font-bold text-info">{{ summary.pulang }}</p>
            <p class="text-[10px] text-gray-500">Pulang</p>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loadingHistory" class="text-center py-8 text-gray-400 text-sm">Memuat...</div>

      <!-- Empty -->
      <div v-else-if="filteredHistory.length === 0" class="text-center py-8">
        <ClockIcon :size="40" class="mx-auto mb-2 text-gray-300" />
        <p class="text-sm text-gray-400">Belum ada riwayat absensi</p>
      </div>

      <!-- List -->
      <div v-else class="space-y-2">
        <div
          v-for="a in filteredHistory"
          :key="a.ID || a.id"
          class="bg-white rounded-xl p-4 border border-gray-100 flex items-center gap-3"
        >
          <div :class="['w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0', a.Type === 'pulang' || a.type === 'pulang' ? 'bg-info/10' : 'bg-accent/10']">
            <LogInIcon v-if="a.Type === 'masuk' || a.type === 'masuk'" :size="18" class="text-accent" />
            <LogOutIcon v-else :size="18" class="text-info" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-800">
              {{ (a.Type || a.type) === 'masuk' ? 'Masuk' : 'Pulang' }} - {{ formatDate(a.Timestamp || a.timestamp) }}
            </p>
            <p class="text-xs text-gray-500">
              {{ formatTime(a.Timestamp || a.timestamp) }} - {{ (a.Latitude ?? a.latitude)?.toFixed(4) }}, {{ (a.Longitude ?? a.longitude)?.toFixed(4) }}
            </p>
          </div>
          <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-[10px] font-medium', (a.Type || a.type) === 'pulang' ? 'bg-info/10 text-info' : 'bg-accent/10 text-accent']">
            {{ (a.Type || a.type) === 'masuk' ? 'Masuk' : 'Pulang' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { MapPinIcon, ClockIcon, CheckCircleIcon, AlertCircleIcon, LoaderIcon, XIcon, ChevronLeftIcon, ChevronRightIcon, LogInIcon, LogOutIcon } from 'lucide-vue-next'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const tabs = [
  { key: 'absen', label: 'Absen' },
  { key: 'riwayat', label: 'Riwayat' }
]
const activeTab = ref('absen')

const status = ref({
  has_masuk: false,
  has_pulang: false,
  pulang_available: false,
  masuk_at: null,
  pulang_at: null,
  pulang_available_at: null,
  server_time: ''
})
const serverTime = ref('')
const absenType = ref(null)
const locating = ref(false)
const location = ref(null)
const mapContainer = ref(null)
const map = ref(null)
let mapMarker = null
const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const history = ref([])
const loadingHistory = ref(false)
const monthOffset = ref(0)

const countdownText = ref('')
let countdownTimer = null

const currentMonth = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() - monthOffset.value)
  return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

const filteredHistory = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() - monthOffset.value)
  const targetMonth = d.getMonth()
  const targetYear = d.getFullYear()
  return history.value.filter((a) => {
    const t = new Date(a.Timestamp || a.timestamp)
    return t.getMonth() === targetMonth && t.getFullYear() === targetYear
  })
})

const summary = computed(() => {
  const items = filteredHistory.value
  return {
    masuk: items.filter((a) => (a.Type || a.type) === 'masuk').length,
    pulang: items.filter((a) => (a.Type || a.type) === 'pulang').length,
  }
})

function prevMonth() { monthOffset.value++ }
function nextMonth() { if (monthOffset.value > 0) monthOffset.value-- }

function formatDate(ts) {
  return new Date(ts).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatTime(ts) {
  if (!ts) return '-'
  const d = new Date(ts)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WIB'
}

function initMap() {
  if (!mapContainer.value || map.value) return
  map.value = L.map(mapContainer.value, { zoomControl: false }).setView([-6.2, 106.8], 13)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: 'OSM'
  }).addTo(map.value)
  L.control.zoom({ position: 'bottomright' }).addTo(map.value)
}

function updateMapMarker() {
  if (!map.value || !location.value) return
  const { lat, lng } = location.value
  if (mapMarker) map.value.removeLayer(mapMarker)
  mapMarker = L.marker([lat, lng]).addTo(map.value)
    .bindPopup('Lokasi Anda')
    .openPopup()
  map.value.setView([lat, lng], 16)
}

function destroyMap() {
  if (map.value) {
    map.value.remove()
    map.value = null
    mapMarker = null
  }
}

function getLocation() {
  locating.value = true
  errorMsg.value = ''
  if (!navigator.geolocation) {
    errorMsg.value = 'Geolokasi tidak didukung di browser ini'
    locating.value = false
    return
  }
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      location.value = { lat: pos.coords.latitude, lng: pos.coords.longitude }
      locating.value = false
      if (map.value) updateMapMarker()
    },
    (err) => {
      locating.value = false
      switch (err.code) {
        case err.PERMISSION_DENIED:
          errorMsg.value = 'Izin lokasi ditolak. Aktifkan GPS dan izinkan akses lokasi.'
          break
        case err.POSITION_UNAVAILABLE:
          errorMsg.value = 'Lokasi tidak tersedia. Pastikan GPS aktif dan coba lagi.'
          break
        case err.TIMEOUT:
          errorMsg.value = 'Waktu mencari lokasi habis. Coba lagi di tempat dengan sinyal lebih baik.'
          break
        default:
          errorMsg.value = 'Gagal mendapatkan lokasi. Pastikan GPS aktif.'
      }
    },
    { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
  )
}

function startMasuk() {
  absenType.value = 'masuk'
  errorMsg.value = ''
  nextTick(initMap)
  getLocation()
}

function startPulang() {
  absenType.value = 'pulang'
  errorMsg.value = ''
  nextTick(initMap)
  getLocation()
}

async function submitAbsensi() {
  if (!absenType.value) return
  submitting.value = true
  errorMsg.value = ''
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/absensi', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        type: absenType.value,
        latitude: location.value?.lat || 0,
        longitude: location.value?.lng || 0,
      })
    })

    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || 'Gagal menyimpan absensi')
    }

    successMsg.value = absenType.value === 'masuk' ? 'Absen masuk berhasil tercatat!' : 'Absen pulang berhasil tercatat!'
    absenType.value = null
    location.value = null
    destroyMap()
    fetchStatus()
    setTimeout(() => { successMsg.value = '' }, 3000)
  } catch (e) {
    errorMsg.value = e.message
  } finally {
    submitting.value = false
  }
}

async function fetchStatus() {
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/absensi/status', {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Gagal')
    const data = await res.json()
    status.value = data
    serverTime.value = data.server_time || ''
    startCountdown()
  } catch (e) {
    console.error(e)
  }
}

function startCountdown() {
  clearInterval(countdownTimer)
  if (!status.value.has_masuk || status.value.has_pulang || status.value.pulang_available) {
    countdownText.value = ''
    return
  }
  updateCountdown()
  countdownTimer = setInterval(updateCountdown, 1000)
}

function updateCountdown() {
  if (!status.value.pulang_available_at) {
    countdownText.value = ''
    return
  }
  const now = new Date()
  const target = new Date(status.value.pulang_available_at)
  const diff = target.getTime() - now.getTime()
  if (diff <= 0) {
    countdownText.value = 'Absen pulang sudah tersedia!'
    clearInterval(countdownTimer)
    fetchStatus()
    return
  }
  const h = Math.floor(diff / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  const s = Math.floor((diff % 60000) / 1000)
  countdownText.value = `${h}j ${m}m ${s}d`
}

async function fetchHistory() {
  loadingHistory.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/absensi/history', {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Gagal')
    const json = await res.json()
    history.value = json.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingHistory.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === 'riwayat') {
    fetchHistory()
  }
})

onMounted(fetchStatus)

onUnmounted(() => {
  clearInterval(countdownTimer)
  destroyMap()
})
</script>
