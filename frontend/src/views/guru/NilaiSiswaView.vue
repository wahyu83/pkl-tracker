<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Nilai PKL Siswa</h2>
      <p class="text-sm text-gray-500 mt-0.5">Isi dan pantau penilaian siswa bimbingan Anda</p>
    </div>

    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 bg-gray-50/50 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <select v-model="filterGrade" class="px-3 py-1.5 rounded-lg border border-gray-200 bg-white text-sm outline-none focus:border-primary">
            <option value="">Semua Siswa</option>
            <option value="belum">Belum Dinilai</option>
            <option value="A">A</option>
            <option value="B">B</option>
            <option value="C">C</option>
            <option value="D">D</option>
          </select>
        </div>
        <button @click="handleExport" class="px-3 py-1.5 rounded-lg border border-gray-200 text-sm text-gray-500 hover:bg-gray-50 flex items-center gap-1">
          <DownloadIcon :size="14" />
          Export
        </button>
      </div>

      <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

      <div v-else-if="filteredRows.length === 0" class="text-center py-8 text-gray-400 text-sm">
        {{ rows.length === 0 ? 'Belum ada siswa bimbingan' : 'Tidak ada siswa sesuai filter' }}
      </div>

      <div v-else class="divide-y divide-gray-50">
        <div v-for="row in filteredRows" :key="row.student.id" class="p-4">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-9 h-9 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary flex-shrink-0">
              {{ (row.student.full_name || '?').charAt(0) }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-800 truncate">{{ row.student.full_name }}</p>
              <p class="text-xs text-gray-400 truncate">{{ row.student.dudi?.company_name || 'Belum ada DUDI' }}</p>
            </div>
            <div class="text-right flex-shrink-0">
              <template v-if="row.grade">
                <p class="text-base font-bold" :class="scoreColor(row.grade.final_score)">
                  {{ Math.round(row.grade.final_score) }}
                </p>
                <span :class="['inline-flex px-2 py-0.5 rounded-full text-[10px] font-bold', gradeStyle(row.grade.final_grade)]">
                  {{ row.grade.final_grade || '-' }}
                </span>
              </template>
              <span v-else class="text-[10px] text-gray-400">Belum dinilai</span>
            </div>
          </div>

          <div v-if="row.grade" class="space-y-1.5 mb-3">
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Kehadiran</span>
              <span class="font-mono font-medium" :class="attendanceColor(row.grade.attendance_score_auto)">
                {{ Math.round(row.grade.attendance_score_auto) }}%
              </span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Memahami alur bisnis</span>
              <span class="font-medium text-gray-700">{{ row.grade.discipline || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan soft skills</span>
              <span class="font-medium text-gray-700">{{ row.grade.responsibility || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan kompetensi teknis</span>
              <span class="font-medium text-gray-700">{{ row.grade.teamwork || '-' }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Menerapkan POS dan K3LH</span>
              <span class="font-medium text-gray-700">{{ row.grade.initiative || '-' }}</span>
            </div>
            <div v-if="row.grade.attendance_start_date && row.grade.attendance_end_date" class="flex items-center justify-between text-xs">
              <span class="text-gray-500">Rentang PKL</span>
              <span class="font-medium text-gray-700 text-right">
                {{ formatDate(row.grade.attendance_start_date) }} – {{ formatDate(row.grade.attendance_end_date) }}
              </span>
            </div>
          </div>

          <router-link
            :to="'/guru/penilaian/' + row.student.id"
            :class="[
              'block w-full text-center py-2.5 rounded-xl text-sm font-medium transition-colors',
              row.grade ? 'text-gray-600 bg-gray-50 hover:bg-gray-100' : 'text-primary bg-primary/5 hover:bg-primary/10'
            ]"
          >
            {{ row.grade ? 'Edit Nilai' : 'Beri Nilai' }}
          </router-link>
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
const rows = ref([])

const filteredRows = computed(() => {
  if (!filterGrade.value) return rows.value
  return rows.value.filter(r => {
    if (filterGrade.value === 'belum') return !r.grade
    if (!r.grade?.final_grade) return false
    return r.grade.final_grade.startsWith(filterGrade.value)
  })
})

function scoreColor(score) {
  if (score >= 80) return 'text-accent'
  if (score >= 70) return 'text-warning'
  return 'text-gray-600'
}

function attendanceColor(score) {
  if (score >= 80) return 'text-accent'
  if (score >= 70) return 'text-warning'
  return 'text-danger'
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function gradeStyle(grade) {
  if (!grade) return 'bg-gray-100 text-gray-400'
  const map = {
    A: 'bg-accent/10 text-accent',
    B: 'bg-info/10 text-info',
    C: 'bg-warning/10 text-warning',
    D: 'bg-danger/10 text-danger'
  }
  return map[grade.charAt(0)] || 'bg-gray-100 text-gray-500'
}

function handleExport() {
  downloadCsv('/export/nilai', 'nilai_pkl.csv').catch(e => alert('Export gagal: ' + e.message))
}

onMounted(async () => {
  try {
    const [studentsRes, nilaiRes] = await Promise.all([
      get('/guru/students'),
      get('/nilai').catch(() => ({ data: [] }))
    ])

    const students = studentsRes.data || []
    const gradeMap = new Map((nilaiRes.data || []).map(g => [g.student_id, g]))

    rows.value = students.map(s => ({
      student: s,
      grade: gradeMap.get(s.id) || null
    }))
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
