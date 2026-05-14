<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Dashboard Guru Pembimbing</h2>
      <p class="text-sm text-gray-500 mt-0.5">Selamat datang, {{ authStore.userName }}</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <StatsCard label="Siswa Bimbingan" :value="28" :icon="Users" iconBg="bg-primary/10" iconColor="text-primary" />
      <StatsCard label="Hadir Hari Ini" :value="24" :icon="ClipboardCheck" iconBg="bg-accent/10" iconColor="text-accent" />
      <StatsCard label="Jurnal Terisi" :value="22" :icon="BookOpen" iconBg="bg-info/10" iconColor="text-info" />
      <StatsCard label="Nilai Tersedia" :value="15" :icon="Award" iconBg="bg-warning/10" iconColor="text-warning" />
    </div>

    <!-- Student cards -->
    <div class="mb-4">
      <h3 class="font-semibold text-gray-800 mb-3">Siswa Bimbingan Anda</h3>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div
        v-for="s in students"
        :key="s.id"
        class="bg-white rounded-2xl p-5 border border-gray-100 card-hover"
      >
        <div class="flex items-start gap-3 mb-3">
          <div class="w-10 h-10 rounded-full bg-primary/10 text-primary flex items-center justify-center font-bold text-sm flex-shrink-0">
            {{ s.name.charAt(0) }}
          </div>
          <div class="min-w-0">
            <h4 class="font-semibold text-gray-800 truncate">{{ s.name }}</h4>
            <p class="text-xs text-gray-500">NIS: {{ s.nis }} | {{ s.dudi }}</p>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-2 mb-3">
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-lg font-bold" :class="s.attendancePercent >= 80 ? 'text-accent' : s.attendancePercent >= 60 ? 'text-warning' : 'text-danger'">
              {{ s.attendancePercent }}%
            </p>
            <p class="text-[10px] text-gray-500">Kehadiran</p>
          </div>
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-lg font-bold text-info">{{ s.journalCount }}</p>
            <p class="text-[10px] text-gray-500">Jurnal</p>
          </div>
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-lg font-bold text-gray-800" :class="s.nilai ? 'text-warning' : 'text-gray-300'">
              {{ s.nilai || '-' }}
            </p>
            <p class="text-[10px] text-gray-500">Nilai</p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <router-link
            :to="'/guru/absensi?siswa=' + s.id"
            class="flex-1 text-center py-2 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors"
          >
            Detail
          </router-link>
          <button class="px-3 py-2 text-xs font-medium text-gray-500 hover:text-info rounded-lg hover:bg-info/5 transition-colors">
            Komentar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../../stores/auth'
import StatsCard from '../../components/StatsCard.vue'
import { Users, ClipboardCheck, BookOpen, Award } from 'lucide-vue-next'

const authStore = useAuthStore()

const students = [
  { id: 1, name: 'Ahmad Rizky', nis: '20230001', dudi: 'PT. Teknologi Maju', attendancePercent: 95, journalCount: 42, nilai: 'A' },
  { id: 2, name: 'Siti Nurhaliza', nis: '20230002', dudi: 'PT. Sejahtera Abadi', attendancePercent: 88, journalCount: 40, nilai: 'B+' },
  { id: 3, name: 'Dian Permata', nis: '20230004', dudi: 'CV. Kreatif Digital', attendancePercent: 72, journalCount: 35, nilai: 'B' },
  { id: 4, name: 'Rudi Hartono', nis: '20230005', dudi: 'UD. Mandiri Jaya', attendancePercent: 100, journalCount: 45, nilai: 'A' },
  { id: 5, name: 'Maya Sari', nis: '20230007', dudi: 'PT. Inovasi Nusantara', attendancePercent: 60, journalCount: 28, nilai: null },
  { id: 6, name: 'Bambang Kusumo', nis: '20230008', dudi: 'PT. Teknologi Maju', attendancePercent: 85, journalCount: 38, nilai: 'B+' },
]
</script>
