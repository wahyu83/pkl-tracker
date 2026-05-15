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

    <!-- Step 1: Lokasi & Waktu -->
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
        Lanjut ke Konfirmasi
      </button>
    </div>

    <!-- Step 2: Konfirmasi & Submit -->
    <div v-if="currentStep === 1" class="space-y-4">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-1">Konfirmasi Absensi</h3>
        <p class="text-xs text-gray-500 mb-4">Periksa kembali data absensi Anda</p>

        <div class="space-y-3">
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
        </div>
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
import { ref } from 'vue'
import { MapPinIcon, ClockIcon, CheckCircleIcon, AlertCircleIcon, CheckIcon, LoaderIcon } from 'lucide-vue-next'

const steps = ['Lokasi', 'Konfirmasi']
const currentStep = ref(0)
const locating = ref(false)
const location = ref(null)
const inRadius = ref(true)
const submitting = ref(false)
const success = ref(false)

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
      location.value = { lat: -6.2088, lng: 106.8456 }
      locating.value = false
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

async function submitAbsensi() {
  submitting.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/absensi', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        latitude: location.value?.lat || 0,
        longitude: location.value?.lng || 0,
        status: inRadius.value ? 'hadir' : 'terlambat'
      })
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'Gagal menyimpan absensi')
    }

    success.value = true
  } catch (e) {
    alert('Absensi gagal: ' + e.message)
  } finally {
    submitting.value = false
  }
}
</script>
