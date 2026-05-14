<template>
  <div>
    <div class="flex items-center gap-3 mb-5">
      <router-link to="/siswa/absensi" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ArrowLeftIcon :size="20" />
      </router-link>
      <div>
        <h2 class="text-lg font-bold text-gray-800">Riwayat Absensi</h2>
        <p class="text-xs text-gray-500">Bulan {{ currentMonth }}</p>
      </div>
    </div>

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
    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <div class="grid grid-cols-4 gap-2 text-center">
        <div>
          <p class="text-xl font-bold text-accent">{{ summary.hadir }}</p>
          <p class="text-[10px] text-gray-500">Hadir</p>
        </div>
        <div>
          <p class="text-xl font-bold text-warning">{{ summary.terlambat }}</p>
          <p class="text-[10px] text-gray-500">Terlambat</p>
        </div>
        <div>
          <p class="text-xl font-bold text-info">{{ summary.izin }}</p>
          <p class="text-[10px] text-gray-500">Izin</p>
        </div>
        <div>
          <p class="text-xl font-bold text-danger">{{ summary.sakit }}</p>
          <p class="text-[10px] text-gray-500">Sakit</p>
        </div>
      </div>
    </div>

    <!-- List -->
    <div class="space-y-2">
      <div
        v-for="a in history"
        :key="a.id"
        class="bg-white rounded-xl p-4 border border-gray-100 flex items-center gap-3"
      >
        <div :class="['w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0', statusIconBg(a.status)]">
          <component :is="statusIcon(a.status)" :size="18" :class="statusIconColor(a.status)" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-800">{{ a.date }}</p>
          <p class="text-xs text-gray-500">{{ a.time }} - {{ a.location?.slice(0, 20) }}...</p>
        </div>
        <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-[10px] font-medium', statusBadgeStyle(a.status)]">
          {{ statusLabel(a.status) }}
        </span>
        <ChevronRightIcon :size="16" class="text-gray-300" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  ArrowLeftIcon, ChevronLeftIcon, ChevronRightIcon,
  CheckCircleIcon, ClockIcon, FileTextIcon, HeartIcon, AlertCircleIcon
} from 'lucide-vue-next'

const monthOffset = ref(0)

const currentMonth = computed(() => {
  const d = new Date()
  d.setMonth(d.getMonth() - monthOffset.value)
  return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

function prevMonth() { monthOffset.value++ }
function nextMonth() { if (monthOffset.value > 0) monthOffset.value-- }

const summary = { hadir: 20, terlambat: 2, izin: 1, sakit: 1 }

const history = [
  { id: 1, date: '14 Mei 2026', time: '07:45 WIB', location: '-6.2088, 106.8456', status: 'hadir' },
  { id: 2, date: '13 Mei 2026', time: '08:15 WIB', location: '-6.2088, 106.8456', status: 'terlambat' },
  { id: 3, date: '12 Mei 2026', time: '07:30 WIB', location: '-6.2088, 106.8456', status: 'hadir' },
  { id: 4, date: '11 Mei 2026', time: '-', location: '-', status: 'izin' },
  { id: 5, date: '10 Mei 2026', time: '-', location: '-', status: 'sakit' },
  { id: 6, date: '9 Mei 2026', time: '07:50 WIB', location: '-6.2088, 106.8456', status: 'hadir' },
]

function statusLabel(s) {
  return { hadir: 'Hadir', terlambat: 'Terlambat', izin: 'Izin', sakit: 'Sakit' }[s]
}

function statusBadgeStyle(s) {
  return {
    hadir: 'bg-accent/10 text-accent',
    terlambat: 'bg-warning/10 text-warning',
    izin: 'bg-info/10 text-info',
    sakit: 'bg-purple-100 text-purple-600'
  }[s]
}

function statusIcon(s) {
  return {
    hadir: CheckCircleIcon,
    terlambat: ClockIcon,
    izin: FileTextIcon,
    sakit: HeartIcon
  }[s]
}

function statusIconBg(s) {
  return {
    hadir: 'bg-accent/10',
    terlambat: 'bg-warning/10',
    izin: 'bg-info/10',
    sakit: 'bg-purple-100'
  }[s]
}

function statusIconColor(s) {
  return {
    hadir: 'text-accent',
    terlambat: 'text-warning',
    izin: 'text-info',
    sakit: 'text-purple-600'
  }[s]
}
</script>
