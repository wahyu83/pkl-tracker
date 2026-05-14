<template>
  <div>
    <!-- Step Navigation -->
    <div class="flex items-center gap-2 mb-5">
      <template v-for="(step, i) in steps" :key="i">
        <div class="flex items-center gap-2">
          <div
            :class="[
              'w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold',
              currentStep > i ? 'bg-accent text-white' : currentStep === i ? 'bg-primary text-white' : 'bg-gray-200 text-gray-400'
            ]"
          >
            <CheckIcon v-if="currentStep > i" :size="14" />
            <span v-else>{{ i + 1 }}</span>
          </div>
          <span :class="['text-xs font-medium hidden sm:block', currentStep >= i ? 'text-gray-800' : 'text-gray-400']">
            {{ step }}
          </span>
        </div>
        <div v-if="i < steps.length - 1" class="flex-1 h-px bg-gray-200" :class="currentStep > i ? 'bg-accent' : ''" />
      </template>
    </div>

    <!-- Step 1: Lokasi & Timestamp -->
    <div v-if="currentStep === 0" class="space-y-4">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-1">Lokasi & Waktu</h3>
        <p class="text-xs text-gray-500 mb-4">Sistem akan mengambil lokasi GPS dan timestamp otomatis</p>

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
              <p class="text-xs text-gray-500">Waktu Absen</p>
              <p class="text-sm font-medium text-gray-800 font-mono">{{ currentTime }}</p>
            </div>
          </div>
        </div>

        <!-- Map placeholder -->
        <div class="mt-4 h-48 bg-gray-200 rounded-xl flex items-center justify-center">
          <div class="text-center text-gray-400">
            <MapPinIcon :size="28" class="mx-auto mb-1" />
            <p class="text-xs">Peta lokasi absensi</p>
          </div>
        </div>

        <!-- Radius validation -->
        <div
          v-if="location"
          :class="[
            'mt-4 p-3 rounded-xl text-xs font-medium flex items-center gap-2',
            inRadius ? 'bg-accent/10 text-accent' : 'bg-danger/10 text-danger'
          ]"
        >
          <CheckCircleIcon v-if="inRadius" :size="16" />
          <AlertCircleIcon v-else :size="16" />
          <span>{{ inRadius ? 'Lokasi valid (dalam radius DUDI)' : 'Peringatan: Anda di luar radius DUDI!' }}</span>
        </div>
      </div>

      <button
        @click="currentStep = 1"
        :disabled="!location"
        class="w-full py-3 bg-primary text-white rounded-2xl text-sm font-bold hover:bg-primary-light transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Lanjut ke Foto
      </button>
    </div>

    <!-- Step 2: Foto Selfie -->
    <div v-if="currentStep === 1" class="space-y-4">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-1">Foto Selfie</h3>
        <p class="text-xs text-gray-500 mb-4">Ambil foto selfie dengan latar tempat PKL</p>

        <div class="relative bg-black rounded-xl overflow-hidden aspect-[3/4] mb-4">
          <div
            v-if="!photoTaken && !streaming && !cameraError"
            class="absolute inset-0 flex flex-col items-center justify-center bg-gray-900"
          >
            <CameraIcon :size="48" class="text-gray-600 mb-3" />
            <p class="text-gray-400 text-sm">Kamera belum aktif</p>
          </div>
          <div
            v-if="cameraError && !photoTaken"
            class="absolute inset-0 flex flex-col items-center justify-center bg-gray-900 p-4"
          >
            <AlertCircleIcon :size="40" class="text-warning mb-3" />
            <p class="text-warning text-sm font-medium text-center">Kamera tidak tersedia</p>
            <p class="text-gray-400 text-xs text-center mt-2">{{ cameraError }}</p>
          </div>
          <video
            v-if="streaming"
            ref="videoRef"
            autoplay
            playsinline
            class="w-full h-full object-cover"
          />
          <img
            v-if="photoTaken"
            :src="photoData"
            class="w-full h-full object-cover"
            alt="Selfie"
          />
        </div>

        <div class="flex gap-3">
          <button
            v-if="!photoTaken"
            @click="startCamera"
            :disabled="streaming"
            class="flex-1 py-3 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <CameraIcon :size="18" />
            {{ streaming ? 'Kamera Aktif' : 'Buka Kamera' }}
          </button>
          <button
            v-if="streaming && !photoTaken"
            @click="capturePhoto"
            class="flex-1 py-3 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark transition-colors flex items-center justify-center gap-2"
          >
            <CameraIcon :size="18" />
            Ambil Foto
          </button>
          <button
            v-if="photoTaken"
            @click="retakePhoto"
            class="flex-1 py-3 bg-gray-100 text-gray-700 rounded-xl text-sm font-medium hover:bg-gray-200 transition-colors flex items-center justify-center gap-2"
          >
            <RefreshCcwIcon :size="18" />
            Ulangi
          </button>
        </div>
      </div>

      <button
        @click="currentStep = 2"
        :disabled="!photoTaken"
        class="w-full py-3 bg-primary text-white rounded-2xl text-sm font-bold hover:bg-primary-light transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Lanjut ke Konfirmasi
      </button>
    </div>

    <!-- Step 3: Konfirmasi & Submit -->
    <div v-if="currentStep === 2" class="space-y-4">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-1">Konfirmasi Absensi</h3>
        <p class="text-xs text-gray-500 mb-4">Periksa kembali data absensi Anda</p>

        <div class="space-y-3 mb-4">
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Waktu</span>
            <span class="font-medium text-gray-800">{{ currentTime }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Lokasi</span>
            <span class="font-medium text-gray-800">{{ location?.lat.toFixed(6) }}, {{ location?.lng.toFixed(6) }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Validasi Radius</span>
            <span :class="['font-medium', inRadius ? 'text-accent' : 'text-danger']">
              {{ inRadius ? 'Valid' : 'Di luar radius' }}
            </span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Foto Bukti</span>
            <span class="font-medium text-accent">✓ Terlampir</span>
          </div>
        </div>

        <img
          v-if="photoData"
          :src="photoData"
          class="w-full rounded-xl object-cover max-h-48"
          alt="Preview"
        />
      </div>

      <button
        @click="submitAbsensi"
        :disabled="submitting"
        class="w-full py-3 bg-accent text-white rounded-2xl text-sm font-bold hover:bg-accent-dark transition-colors disabled:opacity-60 flex items-center justify-center gap-2 shadow-lg shadow-accent/20"
      >
        <LoaderIcon v-if="submitting" :size="20" class="animate-spin" />
        <span>{{ submitting ? 'Menyimpan...' : 'Simpan Absensi' }}</span>
      </button>
    </div>

    <!-- Success Modal -->
    <div
      v-if="success"
      class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-16 h-16 bg-accent/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircleIcon :size="36" class="text-accent" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-2">Absensi Berhasil!</h3>
        <p class="text-sm text-gray-500 mb-4">
          Kehadiran Anda telah tercatat pada {{ currentTime }}
        </p>
        <router-link
          to="/siswa"
          class="block w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors mb-2"
        >
          Kembali ke Dashboard
        </router-link>
        <router-link
          to="/siswa/absensi/history"
          class="block text-sm text-primary font-medium hover:underline"
        >
          Lihat Riwayat Absensi
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import {
  MapPinIcon, ClockIcon, CheckCircleIcon, AlertCircleIcon,
  CameraIcon, CheckIcon, RefreshCcwIcon, LoaderIcon
} from 'lucide-vue-next'

const steps = ['Lokasi', 'Foto', 'Konfirmasi']
const currentStep = ref(0)
const locating = ref(false)
const location = ref(null)
const inRadius = ref(true)
const streaming = ref(false)
const photoTaken = ref(false)
const photoData = ref(null)
const submitting = ref(false)
const success = ref(false)
const videoRef = ref(null)
const cameraError = ref('')

const currentTime = ref(new Date().toLocaleString('id-ID', {
  weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
  hour: '2-digit', minute: '2-digit', second: '2-digit'
}))

function getLocation() {
  locating.value = true
  if (!navigator.geolocation) {
    locating.value = false
    return
  }
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      location.value = { lat: pos.coords.latitude, lng: pos.coords.longitude }
      locating.value = false
    },
    () => {
      // Fallback mock location for demo
      location.value = { lat: -6.2088, lng: 106.8456 }
      locating.value = false
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

async function startCamera() {
  cameraError.value = ''
  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user', width: { ideal: 720 }, height: { ideal: 1280 } }
    })
    if (videoRef.value) {
      videoRef.value.srcObject = stream
    }
    streaming.value = true
  } catch (e) {
    if (e.name === 'NotAllowedError') {
      cameraError.value = 'Izin kamera ditolak. Buka pengaturan browser untuk mengizinkan akses kamera.'
    } else if (e.name === 'NotFoundError') {
      cameraError.value = 'Tidak ada kamera ditemukan di perangkat ini.'
    } else if (e.name === 'NotReadableError') {
      cameraError.value = 'Kamera sedang digunakan aplikasi lain. Tutup aplikasi lain terlebih dahulu.'
    } else {
      cameraError.value = 'Kamera gagal diakses. Pastikan menggunakan HTTPS (bukan HTTP). Error: ' + e.message
    }
    streaming.value = false
  }
}

function capturePhoto() {
  if (!videoRef.value) return
  const canvas = document.createElement('canvas')
  canvas.width = videoRef.value.videoWidth || 360
  canvas.height = videoRef.value.videoHeight || 640
  const ctx = canvas.getContext('2d')
  ctx.drawImage(videoRef.value, 0, 0)

  photoData.value = canvas.toDataURL('image/jpeg', 0.8)
  photoTaken.value = true

  // Stop camera
  if (videoRef.value?.srcObject) {
    videoRef.value.srcObject.getTracks().forEach(t => t.stop())
  }
  streaming.value = false
}

function retakePhoto() {
  photoTaken.value = false
  photoData.value = null
  startCamera()
}

function submitAbsensi() {
  submitting.value = true
  setTimeout(() => {
    submitting.value = false
    success.value = true
  }, 1500)
}
</script>
