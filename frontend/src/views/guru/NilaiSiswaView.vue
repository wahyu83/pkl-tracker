<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Nilai PKL Siswa</h2>
      <p class="text-sm text-gray-500 mt-0.5">Penilaian dari Instruktur DUDI untuk siswa bimbingan Anda</p>
    </div>

    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 bg-gray-50/50 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <select v-model="filterGrade" class="px-3 py-1.5 rounded-lg border border-gray-200 bg-white text-sm outline-none focus:border-primary">
            <option value="">Semua Nilai</option>
            <option value="A">A</option>
            <option value="B">B</option>
            <option value="C">C</option>
            <option value="D">D</option>
            <option value="belum">Belum Dinilai</option>
          </select>
        </div>
        <button @click="handleExport" class="px-3 py-1.5 rounded-lg border border-gray-200 text-sm text-gray-500 hover:bg-gray-50 flex items-center gap-1">
          <DownloadIcon :size="14" />
          Export
        </button>
      </div>

      <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

      <div v-else-if="filteredNilai.length === 0" class="text-center py-8 text-gray-400 text-sm">Belum ada data penilaian</div>

      <div v-else class="divide-y divide-gray-50">
        <div v-for="n in filteredNilai" :key="n.id" class="p-4">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-9 h-9 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary flex-shrink-0">
              {{ (n.student?.full_name || '?').charAt(0) }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-800 truncate">{{ n.student?.full_name || '-' }}</p>
              <p class="text-xs text-gray-400 truncate">{{ n.student?.dudi?.company_name || 'Belum ada DUDI' }}</p>
            </div>
            <div class="text-right flex-shrink-0">
              <p class="text-base font-bold" :class="n.final_score >= 80 ? 'text-accent' : n.final_score >= 70 ? 'text-warning' : 'text-gray-600'">
                {{ n.final_score ? Math.round(n.final_score) : '-' }}
              </p>
              <span v-if="n.final_grade" :class="['inline-flex px-2 py-0.5 rounded-full text-[10px] font-bold', gradeStyle(n.final_grade)]">
                {{ n.final_grade }}
              </span>
              <span v-else class="text-[10px] text-gray-300">Belum dinilai</span>
            </div>
          </div>

          <div class="space-y-1.5">
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Kehadiran</span>
              <span class="font-mono font-medium" :class="n.attendance_score_auto >= 80 ? 'text-accent' : n.attendance_score_auto >= 70 ? 'text-warning' : 'text-danger'">
                {{ Math.round(n.attendance_score_auto) }}%
              </span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Memahami alur bisnis</span>
              <span class="font-medium text-gray-700">{{ n.discipline || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan soft skills</span>
              <span class="font-medium text-gray-700">{{ n.responsibility || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan kompetensi teknis</span>
              <span class="font-medium text-gray-700">{{ n.teamwork || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan POS dan K3LH</span>
              <span class="font-medium text-gray-700">{{ n.initiative || '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Menu Sertifikat -->
    <router-link
      to="/guru/sertifikat"
      class="mt-4 w-full flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-semibold text-primary bg-primary/5 hover:bg-primary/10 transition-colors"
    >
      <AwardIcon :size="18" /> Cetak Sertifikat
    </router-link>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { DownloadIcon, Award as AwardIcon } from 'lucide-vue-next'
import { get, downloadCsv } from '../../api'

const filterGrade = ref('')
const loading = ref(true)
const nilaiData = ref([])

function handleExport() {
  downloadCsv('/export/nilai', 'nilai_pkl.csv').catch(e => alert('Export gagal: ' + e.message))
}

const filteredNilai = computed(() => {
  return nilaiData.value.filter(n => {
    const matchGrade = !filterGrade.value || (filterGrade.value === 'belum' ? !n.final_grade : (n.final_grade || '').startsWith(filterGrade.value))
    return matchGrade
  })
})

function gradeStyle(grade) {
  if (!grade) return 'bg-gray-100 text-gray-300'
  const g = grade.charAt(0)
  const map = {
    A: 'bg-accent/10 text-accent',
    B: 'bg-info/10 text-info',
    C: 'bg-warning/10 text-warning',
    D: 'bg-danger/10 text-danger'
  }
  return map[g] || 'bg-gray-100 text-gray-500'
}

onMounted(async () => {
  try {
    const res = await get('/nilai')
    nilaiData.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
