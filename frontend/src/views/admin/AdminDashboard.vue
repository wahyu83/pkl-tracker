<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">{{ dashboardTitle }}</h2>
      <p class="text-sm text-gray-500 mt-0.5">Selamat datang, {{ authStore.userName }}</p>
      <p v-if="authStore.userRole === 'admin_jurusan' && authStore.userJurusan" class="text-xs text-accent/70 mt-1">Jurusan: {{ authStore.userJurusan }}</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <StatsCard label="Total Siswa" :value="1250" :icon="Users" iconBg="bg-primary/10" iconColor="text-primary" :change="12" />
      <StatsCard label="Total Guru" :value="48" :icon="UserCircle" iconBg="bg-accent/10" iconColor="text-accent" :change="5" />
      <StatsCard label="Total DUDI" :value="32" :icon="Building2" iconBg="bg-warning/10" iconColor="text-warning" :change="8" />
      <StatsCard label="Periode Aktif" :value="3" :icon="Calendar" iconBg="bg-info/10" iconColor="text-info" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-4">Aktivitas Terbaru</h3>
        <div class="space-y-3">
          <div v-for="item in activities" :key="item.id" class="flex items-start gap-3 py-2 border-b border-gray-50 last:border-0">
            <div :class="['w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0', item.bg]">
              <component :is="item.icon" :size="16" :class="item.color" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm text-gray-700 truncate">{{ item.text }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ item.time }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-4">Distribusi Pengguna</h3>
        <div class="space-y-4">
          <div v-for="item in distributions" :key="item.label">
            <div class="flex justify-between text-sm mb-1">
              <span class="text-gray-600">{{ item.label }}</span>
              <span class="font-medium text-gray-800">{{ item.count }}</span>
            </div>
            <div class="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                :class="item.barColor"
                class="h-full rounded-full transition-all duration-500"
                :style="{ width: item.percent + '%' }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../../stores/auth'
import { computed } from 'vue'
import StatsCard from '../../components/StatsCard.vue'
import { Users, UserCircle, Building2, Calendar, ClipboardCheck, BookOpen, Award, UserPlus } from 'lucide-vue-next'

const authStore = useAuthStore()

const dashboardTitle = computed(() => {
  if (authStore.userRole === 'admin_jurusan') return 'Dashboard Jurusan'
  return 'Dashboard Admin'
})

const activities = [
  { id: 1, text: 'Siswa "Ahmad Rizky" melakukan absensi hari ini', time: '5 menit lalu', icon: ClipboardCheck, bg: 'bg-accent/10', color: 'text-accent' },
  { id: 2, text: 'Guru "Budi Santoso" menambah komentar jurnal', time: '15 menit lalu', icon: BookOpen, bg: 'bg-info/10', color: 'text-info' },
  { id: 3, text: 'DUDI "PT. Maju" memberikan nilai PKL', time: '1 jam lalu', icon: Award, bg: 'bg-warning/10', color: 'text-warning' },
  { id: 4, text: 'Registrasi siswa baru: "Siti Nurhaliza"', time: '2 jam lalu', icon: UserPlus, bg: 'bg-primary/10', color: 'text-primary' },
  { id: 5, text: 'Admin memperbarui data DUDI "PT. Sejahtera"', time: '3 jam lalu', icon: Building2, bg: 'bg-gray-100', color: 'text-gray-500' },
]

const distributions = [
  { label: 'Siswa', count: 1250, percent: 78, barColor: 'bg-primary' },
  { label: 'Guru', count: 48, percent: 15, barColor: 'bg-accent' },
  { label: 'DUDI', count: 32, percent: 10, barColor: 'bg-warning' },
  { label: 'Admin', count: 3, percent: 2, barColor: 'bg-gray-400' },
]
</script>
