<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Rekap & Laporan</h2>
      <p class="text-sm text-gray-500 mt-0.5">Export laporan dan ringkasan monitoring</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
      <div v-for="r in reports" :key="r.title" class="bg-white rounded-2xl p-5 border border-gray-100 card-hover">
        <div class="flex items-start gap-3">
          <div :class="['w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0', r.bg]">
            <component :is="r.icon" :size="20" :class="r.color" />
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-gray-800 mb-1">{{ r.title }}</h3>
            <p class="text-xs text-gray-500 mb-3">{{ r.desc }}</p>
            <div class="flex items-center gap-2">
              <button
                @click="viewReport(r.key)"
                class="px-3 py-1.5 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors"
              >
                Lihat
              </button>
              <button
                @click="exportCsv(r.key)"
                :disabled="!r.key"
                class="px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors flex items-center gap-1 disabled:opacity-30"
              >
                <DownloadIcon :size="14" />
                Export
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Data Modal -->
    <div v-if="viewing" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="viewing = null">
      <div class="fixed inset-0 bg-black/50" />
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-3xl max-h-[80vh] flex flex-col">
        <div class="flex items-center justify-between p-5 border-b border-gray-100">
          <h3 class="text-lg font-bold text-gray-800">{{ viewingTitle }}</h3>
          <button @click="viewing = null" class="p-1.5 rounded-lg hover:bg-gray-100 text-gray-400">
            <XIcon :size="20" />
          </button>
        </div>
        <div class="overflow-auto p-5 flex-1">
          <div v-if="loadingView" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>
          <div v-else-if="viewData.length === 0" class="text-center py-8 text-gray-400 text-sm">Belum ada data.</div>
          <div v-else class="overflow-x-auto">
            <table class="w-full text-xs">
              <thead>
                <tr class="border-b border-gray-100">
                  <th v-for="col in viewColumns" :key="col" class="text-left px-3 py-2 font-semibold text-gray-500 whitespace-nowrap">{{ col }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, i) in viewData" :key="i" class="border-b border-gray-50 hover:bg-gray-50/50">
                  <td v-for="col in viewColumns" :key="col" class="px-3 py-2 text-gray-700 whitespace-nowrap max-w-[200px] truncate">{{ row[col] || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="flex items-center gap-3 p-4 border-t border-gray-100">
          <button @click="viewing = null" class="flex-1 py-2 rounded-xl text-sm font-medium border border-gray-200 text-gray-600 hover:bg-gray-50">Tutup</button>
          <button @click="exportCurrentView" class="flex-1 py-2 rounded-xl text-sm font-medium bg-primary text-white hover:bg-primary-light flex items-center justify-center gap-1">
            <DownloadIcon :size="16" />
            Export CSV
          </button>
        </div>
      </div>
    </div>

    <!-- Summary by company -->
    <div class="bg-white rounded-2xl p-5 border border-gray-100">
      <h3 class="font-semibold text-gray-800 mb-4">Ringkasan per DUDI</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <div v-for="item in summarByDudi" :key="item.name" class="p-4 rounded-xl bg-gray-50">
          <h4 class="font-medium text-gray-800 text-sm mb-2">{{ item.name }}</h4>
          <div class="space-y-1.5 text-xs text-gray-500">
            <div class="flex justify-between">
              <span>Siswa</span>
              <span class="font-medium text-gray-700">{{ item.students }}</span>
            </div>
            <div class="flex justify-between">
              <span>Rata-rata Kehadiran</span>
              <span class="font-medium" :class="item.avgAttendance >= 80 ? 'text-accent' : 'text-warning'">{{ item.avgAttendance }}%</span>
            </div>
            <div class="flex justify-between">
              <span>Sudah Dinilai</span>
              <span class="font-medium text-gray-700">{{ item.rated }}/{{ item.students }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { FileText, ClipboardCheck, BookOpen, Award, DownloadIcon, XIcon } from 'lucide-vue-next'
import { downloadCsv, get } from '../../api'

const viewing = ref(null)
const loadingView = ref(false)
const viewData = ref([])

const reports = [
  { key: 'absensi', title: 'Rekap Absensi Per Siswa', desc: 'Laporan detail kehadiran harian lengkap dengan lokasi', icon: ClipboardCheck, bg: 'bg-primary/10', color: 'text-primary' },
  { key: 'jurnal', title: 'Rekap Jurnal Per Siswa', desc: 'Laporan isi jurnal dan komentar', icon: BookOpen, bg: 'bg-accent/10', color: 'text-accent' },
  { key: 'nilai', title: 'Rekap Nilai Per Siswa', desc: 'Laporan penilaian lengkap dari Instruktur', icon: Award, bg: 'bg-warning/10', color: 'text-warning' },
  { key: null, title: 'Gabungan Bulanan', desc: 'Gabungan absensi + jurnal + nilai', icon: FileText, bg: 'bg-gray-100', color: 'text-gray-600' },
]

const viewingTitle = computed(() => {
  const titles = { absensi: 'Rekap Absensi', jurnal: 'Rekap Jurnal', nilai: 'Rekap Nilai' }
  return titles[viewing.value] || 'Data'
})

const viewColumns = computed(() => {
  if (viewData.value.length > 0) return Object.keys(viewData.value[0])
  return []
})

const absensiKeys = {
  'Tanggal': 'Tanggal', 'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS', 'DUDI': 'DUDI',
  'Status': 'Status', 'Latitude': 'Latitude', 'Longitude': 'Longitude', 'Terverifikasi': 'Terverifikasi'
}

const jurnalKeys = {
  'Tanggal': 'Tanggal', 'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS',
  'Kegiatan': 'Kegiatan', 'Refleksi': 'Refleksi', 'Komentar Guru': 'Komentar Guru',   'Komentar Instruktur': 'Komentar Instruktur'
}

const nilaiKeys = {
  'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS', 'DUDI': 'DUDI', 'Kehadiran (%)': 'Kehadiran (%)',
  'Alur Bisnis (0-100)': 'Alur Bisnis (0-100)', 'Soft Skills (0-100)': 'Soft Skills (0-100)',
  'Kompetensi Teknis (0-100)': 'Kompetensi Teknis (0-100)', 'POS & K3LH (0-100)': 'POS & K3LH (0-100)',
  'Nilai Akhir': 'Nilai Akhir', 'Grade': 'Grade', 'Catatan': 'Catatan'
}

async function viewReport(key) {
  if (!key) return
  viewing.value = key
  loadingView.value = true
  viewData.value = []
  try {
    const data = await get(`/report/${key}`)
    const list = data.data || []
    const keyMap = key === 'absensi' ? absensiKeys : key === 'jurnal' ? jurnalKeys : nilaiKeys
    viewData.value = list.map(row => {
      const mapped = {}
      for (const [label, field] of Object.entries(keyMap)) {
        let val = row[field]
        if (val === undefined) val = row[field.toLowerCase().replace(/ /g, '_')]
        if (typeof val === 'boolean') val = val ? 'Ya' : 'Tidak'
        if (val === '' || val === null || val === undefined) val = '-'
        if (typeof val === 'string' && val.length > 60) val = val.substring(0, 57) + '...'
        mapped[label] = val
      }
      return mapped
    })
  } catch (e) {
    viewData.value = []
  } finally {
    loadingView.value = false
  }
}

async function exportCsv(key) {
  if (!key) return
  try {
    await downloadCsv(`/export/${key}`, `${key}_export.csv`)
  } catch (e) {
    alert('Export gagal: ' + e.message)
  }
}

function exportCurrentView() {
  if (!viewing.value) return
  exportCsv(viewing.value)
}

const summarByDudi = ref([])

onMounted(() => { fetchSummary() })

async function fetchSummary() {
  try {
    const res = await get('/guru/students')
    const data = res.data || []
    const dudiMap = {}
    for (const s of data) {
      const dudiName = s.dudi?.company_name || 'Tanpa DUDI'
      if (!dudiMap[dudiName]) {
        dudiMap[dudiName] = { students: 0, totalAttendance: 0, rated: 0 }
      }
      dudiMap[dudiName].students++
      dudiMap[dudiName].totalAttendance += s.attendance_percent || 0
      if (s.nilai) dudiMap[dudiName].rated++
    }
    summarByDudi.value = Object.entries(dudiMap).map(([name, info]) => ({
      name,
      students: info.students,
      avgAttendance: info.students > 0 ? Math.round(info.totalAttendance / info.students) : 0,
      rated: info.rated
    }))
  } catch (e) {
    console.error(e)
  }
}
</script>
