<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Halo, {{ authStore.userName }}!</h2>
      <p class="text-sm text-gray-500">Selamat ber-PKL hari ini</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-3 gap-3 mb-5">
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <p class="text-2xl font-bold text-accent">{{ stats.hadir }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Hadir</p>
      </div>
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <p class="text-2xl font-bold text-gray-800">{{ stats.jurnal }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Jurnal</p>
      </div>
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <p class="text-2xl font-bold text-warning">{{ stats.nilai || '?' }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Nilai</p>
      </div>
    </div>

    <!-- Status Hari Ini -->
    <div
      :class="[
        'rounded-2xl p-5 mb-5 border-2',
        todayStatus === 'hadir'
          ? 'bg-accent/5 border-accent/20'
          : todayStatus === 'pending'
            ? 'bg-warning/5 border-warning/20'
            : 'bg-gray-50 border-gray-100'
      ]"
    >
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-semibold text-gray-800">Status Kehadiran Hari Ini</p>
          <p class="text-xs text-gray-500 mt-0.5">{{ todayDate }}</p>
        </div>
        <span :class="['inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold', todayStatusStyle]">
          <CheckCircleIcon v-if="todayStatus === 'hadir'" :size="16" />
          <ClockIcon v-else-if="todayStatus === 'pending'" :size="16" />
          <XCircleIcon v-else :size="16" />
          {{ todayStatusLabel }}
        </span>
      </div>

      <router-link
        v-if="todayStatus === 'pending'"
        to="/siswa/absensi"
        class="mt-4 flex items-center justify-center gap-2 w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
      >
        <CameraIcon :size="18" />
        Absen Sekarang
      </router-link>
    </div>

    <!-- Recent journals -->
    <div class="mb-3 flex items-center justify-between">
      <h3 class="font-semibold text-gray-800">Jurnal Terbaru</h3>
      <router-link to="/siswa/jurnal" class="text-xs text-primary font-medium hover:underline">
        Lihat Semua
      </router-link>
    </div>

    <div class="space-y-3 mb-5">
      <div
        v-for="j in recentJournals"
        :key="j.id"
        class="bg-white rounded-xl p-4 border border-gray-100"
      >
        <div class="flex items-start justify-between mb-2">
          <span class="text-xs text-gray-400">{{ j.date }}</span>
          <span :class="['text-[10px] font-medium px-2 py-0.5 rounded-full', j.hasComment ? 'bg-info/10 text-info' : 'bg-gray-100 text-gray-400']">
            {{ j.hasComment ? 'Dikomentari' : 'Belum dikomentari' }}
          </span>
        </div>
        <p class="text-sm text-gray-700 line-clamp-2">{{ j.activity }}</p>
      </div>
    </div>

    <router-link
      to="/siswa/jurnal/tulis"
      class="flex items-center justify-center gap-2 w-full py-3 bg-accent text-white rounded-2xl text-sm font-bold hover:bg-accent-dark transition-colors shadow-lg shadow-accent/20"
    >
      <PlusIcon :size="20" />
      Tulis Jurnal Baru
    </router-link>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { CheckCircleIcon, ClockIcon, XCircleIcon, CameraIcon, PlusIcon } from 'lucide-vue-next'

const authStore = useAuthStore()

const todayStatus = ref('pending') // 'hadir' | 'pending' | 'alpa'

const todayDate = computed(() => {
  return new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

const todayStatusLabel = computed(() => ({
  hadir: 'Sudah Absen',
  pending: 'Belum Absen',
  alpa: 'Tidak Hadir'
}[todayStatus.value]))

const todayStatusStyle = computed(() => ({
  hadir: 'bg-accent/10 text-accent',
  pending: 'bg-warning/10 text-warning',
  alpa: 'bg-danger/10 text-danger'
}[todayStatus.value]))

const stats = { hadir: 42, jurnal: 38, nilai: 'A' }

const recentJournals = [
  { id: 1, date: '14 Mei 2026', activity: 'Mempelajari framework Laravel dan membuat CRUD sederhana untuk modul inventaris.', hasComment: true },
  { id: 2, date: '13 Mei 2026', activity: 'Debugging aplikasi internal perusahaan. Memperbaiki bug pada modul pelaporan.', hasComment: false },
  { id: 3, date: '12 Mei 2026', activity: 'Meeting dengan tim developer, membahas sprint planning minggu depan.', hasComment: true },
]
</script>
