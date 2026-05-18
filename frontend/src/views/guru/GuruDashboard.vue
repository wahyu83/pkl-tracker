<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Dashboard Guru Pembimbing</h2>
      <p class="text-sm text-gray-500 mt-0.5">Selamat datang, {{ authStore.userName }}</p>
    </div>

    <div class="grid grid-cols-2 gap-3 mb-6">
      <StatsCard label="Siswa Bimbingan" :value="stats.total_students" :icon="Users" iconBg="bg-primary/10" iconColor="text-primary" />
      <StatsCard label="Hadir Hari Ini" :value="stats.hadir_hari_ini" :icon="ClipboardCheck" iconBg="bg-accent/10" iconColor="text-accent" />
      <StatsCard label="Jurnal Terisi" :value="stats.total_jurnal" :icon="BookOpen" iconBg="bg-info/10" iconColor="text-info" />
      <StatsCard label="Nilai Tersedia" :value="stats.total_nilai" :icon="Award" iconBg="bg-warning/10" iconColor="text-warning" />
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <template v-else>
      <div class="mb-4">
        <h3 class="font-semibold text-gray-800 mb-3">Siswa Bimbingan Anda</h3>
      </div>
      <div v-if="students.length === 0" class="text-center py-8 text-gray-400 text-sm">
        Belum ada siswa yang ditugaskan kepada Anda. Hubungi admin jurusan.
      </div>
      <div class="grid grid-cols-1 gap-4">
        <div
          v-for="s in students"
          :key="s.id"
          class="bg-white rounded-2xl p-4 border border-gray-100"
        >
          <div class="flex items-start gap-3 mb-3">
            <div class="w-10 h-10 rounded-full bg-primary/10 text-primary flex items-center justify-center font-bold text-sm flex-shrink-0">
              {{ s.full_name?.charAt(0) }}
            </div>
            <div class="min-w-0 flex-1">
              <h4 class="font-semibold text-gray-800 truncate">{{ s.full_name }}</h4>
              <p class="text-xs text-gray-500">NIS: {{ s.nis_nip_nik }} | {{ s.dudi?.company_name || '-' }}</p>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-2">
            <div class="text-center p-2 bg-gray-50 rounded-lg">
              <p class="text-lg font-bold" :class="s.attendance_percent >= 80 ? 'text-accent' : s.attendance_percent >= 60 ? 'text-warning' : 'text-danger'">
                {{ Math.round(s.attendance_percent) }}%
              </p>
              <p class="text-[10px] text-gray-500">Kehadiran</p>
            </div>
            <div class="text-center p-2 bg-gray-50 rounded-lg">
              <p class="text-lg font-bold text-info">{{ s.journal_count }}</p>
              <p class="text-[10px] text-gray-500">Jurnal</p>
            </div>
            <div class="text-center p-2 bg-gray-50 rounded-lg">
              <p class="text-lg font-bold" :class="s.nilai ? 'text-warning' : 'text-gray-300'">
                {{ s.nilai || '-' }}
              </p>
              <p class="text-[10px] text-gray-500">Nilai</p>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import StatsCard from '../../components/StatsCard.vue'
import { Users, ClipboardCheck, BookOpen, Award } from 'lucide-vue-next'
import { get } from '../../api'

const authStore = useAuthStore()
const students = ref([])
const loading = ref(true)
const stats = ref({ total_students: 0, hadir_hari_ini: 0, total_jurnal: 0, total_nilai: 0 })

onMounted(async () => {
  try {
    const [dashRes, studentsRes] = await Promise.all([
      get('/guru/dashboard'),
      get('/guru/students')
    ])
    stats.value = dashRes.stats || stats.value
    students.value = studentsRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
