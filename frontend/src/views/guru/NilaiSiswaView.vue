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

      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Siswa</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">DUDI</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Kehadiran</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Disiplin</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Tanggung Jawab</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Kerja Sama</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Inisiatif</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Nilai Akhir</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Grade</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredNilai.length === 0">
              <td colspan="9" class="px-4 py-8 text-center text-gray-400 text-sm">Belum ada data penilaian</td>
            </tr>
            <tr v-for="n in filteredNilai" :key="n.id" class="border-b border-gray-50 hover:bg-gray-50/50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                    {{ (n.student?.full_name || '?').charAt(0) }}
                  </div>
                  <span class="text-sm font-medium text-gray-800">{{ n.student?.full_name || '-' }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ n.student?.dudi?.company_name || '-' }}</td>
              <td class="px-4 py-3 text-center text-sm font-mono" :class="n.attendance_score_auto >= 80 ? 'text-accent' : n.attendance_score_auto >= 70 ? 'text-warning' : 'text-danger'">
                {{ Math.round(n.attendance_score_auto) }}
              </td>
              <td class="px-4 py-3 text-center text-sm">{{ n.discipline || '-' }}</td>
              <td class="px-4 py-3 text-center text-sm">{{ n.responsibility || '-' }}</td>
              <td class="px-4 py-3 text-center text-sm">{{ n.teamwork || '-' }}</td>
              <td class="px-4 py-3 text-center text-sm">{{ n.initiative || '-' }}</td>
              <td class="px-4 py-3 text-center text-sm font-bold" :class="n.final_score >= 80 ? 'text-accent' : n.final_score >= 70 ? 'text-warning' : 'text-gray-600'">
                {{ n.final_score ? Math.round(n.final_score) : '-' }}
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.final_grade" :class="['inline-flex w-8 h-8 rounded-full items-center justify-center text-sm font-bold', gradeStyle(n.final_grade)]">
                  {{ n.final_grade }}
                </span>
                <span v-else class="text-xs text-gray-300">-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { DownloadIcon } from 'lucide-vue-next'
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
