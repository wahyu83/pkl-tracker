<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Nilai PKL Siswa</h2>
      <p class="text-sm text-gray-500 mt-0.5">Penilaian dari Instruktur DUDI untuk siswa bimbingan Anda</p>
    </div>

    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 bg-gray-50/50">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <FilterIcon :size="16" class="text-primary" />
          </div>
          <div class="flex-1 flex flex-wrap gap-2">
            <select v-model="filterDudi" class="px-3 py-1.5 rounded-lg border border-gray-200 bg-white text-sm outline-none focus:border-primary">
              <option value="">Semua DUDI</option>
              <option v-for="d in dudiList" :key="d" :value="d">{{ d }}</option>
            </select>
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
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Siswa</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">DUDI</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Kehadiran</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Disiplin</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Tanggung Jawab</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Kerjasama</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Inisiatif</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Nilai Akhir</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Grade</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="n in filteredNilai" :key="n.id" class="border-b border-gray-50 hover:bg-gray-50/50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                    {{ n.student.charAt(0) }}
                  </div>
                  <span class="text-sm font-medium text-gray-800">{{ n.student }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ n.dudi }}</td>
              <td class="px-4 py-3 text-center text-sm font-mono" :class="n.scores.attendance >= 80 ? 'text-accent' : n.scores.attendance >= 70 ? 'text-warning' : 'text-danger'">
                {{ n.scores.attendance }}
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.scores.discipline" class="text-sm">{{ n.scores.discipline }}/5</span>
                <span v-else class="text-xs text-gray-300">-</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.scores.responsibility" class="text-sm">{{ n.scores.responsibility }}/5</span>
                <span v-else class="text-xs text-gray-300">-</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.scores.teamwork" class="text-sm">{{ n.scores.teamwork }}/5</span>
                <span v-else class="text-xs text-gray-300">-</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.scores.initiative" class="text-sm">{{ n.scores.initiative }}/5</span>
                <span v-else class="text-xs text-gray-300">-</span>
              </td>
              <td class="px-4 py-3 text-center text-sm font-bold" :class="n.finalScore >= 80 ? 'text-accent' : n.finalScore >= 70 ? 'text-warning' : 'text-gray-600'">
                {{ n.finalScore || '-' }}
              </td>
              <td class="px-4 py-3 text-center">
                <span v-if="n.grade" :class="['inline-flex w-8 h-8 rounded-full items-center justify-center text-sm font-bold', gradeStyle(n.grade)]">
                  {{ n.grade }}
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
import { ref, computed } from 'vue'
import { FilterIcon, DownloadIcon } from 'lucide-vue-next'
import { downloadCsv } from '../../api'

const filterDudi = ref('')
const filterGrade = ref('')

function handleExport() {
  downloadCsv('/export/nilai', 'nilai_pkl.csv').catch(e => alert('Export gagal: ' + e.message))
}

const dudiList = ['PT. Teknologi Maju', 'PT. Sejahtera Abadi', 'CV. Kreatif Digital', 'UD. Mandiri Jaya', 'PT. Inovasi Nusantara']

const nilaiData = [
  { id: 1, student: 'Ahmad Rizky', dudi: 'PT. Teknologi Maju', scores: { attendance: 95, discipline: 5, responsibility: 5, teamwork: 4, initiative: 5 }, finalScore: 94, grade: 'A' },
  { id: 2, student: 'Siti Nurhaliza', dudi: 'PT. Sejahtera Abadi', scores: { attendance: 88, discipline: 4, responsibility: 4, teamwork: 5, initiative: 4 }, finalScore: 84, grade: 'B+' },
  { id: 3, student: 'Dian Permata', dudi: 'CV. Kreatif Digital', scores: { attendance: 72, discipline: 3, responsibility: 3, teamwork: 4, initiative: 3 }, finalScore: 72, grade: 'B' },
  { id: 4, student: 'Rudi Hartono', dudi: 'UD. Mandiri Jaya', scores: { attendance: 100, discipline: 5, responsibility: 5, teamwork: 5, initiative: 5 }, finalScore: 98, grade: 'A' },
  { id: 5, student: 'Maya Sari', dudi: 'PT. Inovasi Nusantara', scores: { attendance: 60, discipline: null, responsibility: null, teamwork: null, initiative: null }, finalScore: null, grade: null },
]

const filteredNilai = computed(() => {
  return nilaiData.filter(n => {
    const matchDudi = !filterDudi.value || n.dudi === filterDudi.value
    const matchGrade = !filterGrade.value || (filterGrade.value === 'belum' ? !n.grade : n.grade?.startsWith(filterGrade.value))
    return matchDudi && matchGrade
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
</script>
