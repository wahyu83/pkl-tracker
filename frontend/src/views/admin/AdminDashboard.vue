<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">{{ dashboardTitle }}</h2>
      <p class="text-sm text-gray-500 mt-0.5">Selamat datang, {{ authStore.userName }}</p>
      <p v-if="authStore.userRole === 'admin_jurusan' && authStore.userJurusan" class="text-xs text-accent/70 mt-1">Jurusan: {{ authStore.userJurusan }}</p>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat dashboard...</div>

    <template v-else>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <StatsCard label="Total Siswa" :value="stats.total_siswa" :icon="Users" iconBg="bg-primary/10" iconColor="text-primary" />
        <StatsCard label="Total Guru" :value="stats.total_guru" :icon="UserCircle" iconBg="bg-accent/10" iconColor="text-accent" />
        <StatsCard label="Total DUDI" :value="stats.total_dudi" :icon="Building2" iconBg="bg-warning/10" iconColor="text-warning" />
        <StatsCard label="Periode Aktif" :value="stats.active_period" :icon="Calendar" iconBg="bg-info/10" iconColor="text-info" />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded-2xl p-5 border border-gray-100">
          <h3 class="font-semibold text-gray-800 mb-4">Aktivitas Terbaru</h3>
          <div v-if="recentActivities.length === 0" class="text-center py-6 text-gray-400 text-sm">Belum ada aktivitas</div>
          <div v-else class="space-y-3">
            <div v-for="item in recentActivities" :key="item.text + item.time" class="flex items-start gap-3 py-2 border-b border-gray-50 last:border-0">
              <div :class="['w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0', activityStyle(item.type).bg]">
                <component :is="activityStyle(item.type).icon" :size="16" :class="activityStyle(item.type).color" />
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
          <div v-if="distributions.every(d => d.count === 0)" class="text-center py-6 text-gray-400 text-sm">Belum ada pengguna</div>
          <div v-else class="space-y-4">
            <div v-for="item in distributions" :key="item.label">
              <div class="flex justify-between text-sm mb-1">
                <span class="text-gray-600">{{ item.label }}</span>
                <span class="font-medium text-gray-800">{{ item.count }}</span>
              </div>
              <div class="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                <div
                  :class="distributionBar(item.role)"
                  class="h-full rounded-full transition-all duration-500"
                  :style="{ width: item.percent + '%' }"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { get } from '@/api'
import StatsCard from '../../components/StatsCard.vue'
import { Users, UserCircle, Building2, Calendar, ClipboardCheck, BookOpen, Award, UserPlus } from 'lucide-vue-next'

const authStore = useAuthStore()
const loading = ref(true)

const stats = ref({ total_siswa: 0, total_guru: 0, total_dudi: 0, total_admin: 0, active_period: 0 })
const recentActivities = ref([])
const distributions = ref([])

const dashboardTitle = computed(() => {
  if (authStore.userRole === 'admin_jurusan') return 'Dashboard Jurusan'
  return 'Dashboard Admin'
})

function activityStyle(type) {
  switch (type) {
    case 'absensi': return { icon: ClipboardCheck, bg: 'bg-accent/10', color: 'text-accent' }
    case 'jurnal': return { icon: BookOpen, bg: 'bg-info/10', color: 'text-info' }
    case 'penilaian': return { icon: Award, bg: 'bg-warning/10', color: 'text-warning' }
    default: return { icon: UserPlus, bg: 'bg-primary/10', color: 'text-primary' }
  }
}

function distributionBar(role) {
  switch (role) {
    case 'student': return 'bg-primary'
    case 'teacher': return 'bg-accent'
    case 'dudi': return 'bg-warning'
    default: return 'bg-gray-400'
  }
}

async function fetchDashboard() {
  loading.value = true
  try {
    const res = await get('/admin/dashboard')
    stats.value = res.stats
    recentActivities.value = res.recent_activities || []
    distributions.value = res.distributions || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchDashboard)
</script>
