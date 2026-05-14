<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Rekap & Laporan</h2>
      <p class="text-sm text-gray-500 mt-0.5">Export dan lihat laporan PKL</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      <div
        v-for="report in reports"
        :key="report.title"
        class="bg-white rounded-2xl p-5 border border-gray-100 card-hover"
      >
        <div :class="['w-10 h-10 rounded-xl flex items-center justify-center mb-3', report.bg]">
          <component :is="report.icon" :size="20" :class="report.color" />
        </div>
        <h3 class="font-semibold text-gray-800 mb-1">{{ report.title }}</h3>
        <p class="text-xs text-gray-500 mb-4">{{ report.desc }}</p>
        <div class="flex items-center gap-2">
          <button
            @click="viewReport(report.key)"
            class="flex-1 py-2 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors"
          >
            Lihat
          </button>
          <button
            @click="exportCsv(report.key)"
            :disabled="!report.key"
            class="flex-1 py-2 text-xs font-medium text-gray-600 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors flex items-center justify-center gap-1 disabled:opacity-30"
          >
            <DownloadIcon :size="14" />
            Export
          </button>
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

    <div class="bg-white rounded-2xl p-5 border border-gray-100">
      <h3 class="font-semibold text-gray-800 mb-4">Log Sistem</h3>
      <div class="space-y-3">
        <div v-for="log in logs" :key="log.id" class="flex items-start gap-3 py-2 border-b border-gray-50 last:border-0">
          <div :class="['w-2 h-2 rounded-full mt-1.5 flex-shrink-0', log.level === 'error' ? 'bg-danger' : log.level === 'warning' ? 'bg-warning' : 'bg-accent']" />
          <div class="min-w-0 flex-1">
            <p class="text-sm text-gray-700">{{ log.message }}</p>
            <p class="text-xs text-gray-400 mt-0.5">{{ log.time }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { FileText, ClipboardCheck, BookOpen, Award, DownloadIcon, MapPinIcon, CalendarIcon, XIcon } from 'lucide-vue-next'
import { downloadCsv, get } from '../../api'

const viewing = ref(null)
const loadingView = ref(false)
const viewData = ref([])

const reports = [
  { key: 'absensi', title: 'Rekap Absensi', desc: 'Laporan kehadiran siswa per periode PKL', icon: ClipboardCheck, bg: 'bg-primary/10', color: 'text-primary' },
  { key: 'jurnal', title: 'Rekap Jurnal', desc: 'Laporan jurnal harian siswa PKL', icon: BookOpen, bg: 'bg-accent/10', color: 'text-accent' },
  { key: 'nilai', title: 'Rekap Nilai', desc: 'Laporan penilaian dari DUDI per siswa', icon: Award, bg: 'bg-warning/10', color: 'text-warning' },
  { key: null, title: 'Distribusi Lokasi', desc: 'Heatmap lokasi absensi siswa', icon: MapPinIcon, bg: 'bg-info/10', color: 'text-info' },
  { key: null, title: 'Export Lengkap', desc: 'Export semua data (PDF/Excel)', icon: FileText, bg: 'bg-gray-100', color: 'text-gray-600' },
  { key: null, title: 'Ringkasan Periode', desc: 'Laporan ringkasan per periode PKL', icon: CalendarIcon, bg: 'bg-purple-100', color: 'text-purple-600' },
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
  'Kegiatan': 'Kegiatan', 'Refleksi': 'Refleksi', 'Komentar Guru': 'Komentar Guru', 'Komentar DUDI': 'Komentar DUDI'
}

const nilaiKeys = {
  'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS', 'DUDI': 'DUDI', 'Kehadiran (%)': 'Kehadiran (%)',
  'Disiplin (1-5)': 'Disiplin (1-5)', 'Tanggung Jawab (1-5)': 'Tanggung Jawab (1-5)',
  'Kerjasama (1-5)': 'Kerjasama (1-5)', 'Inisiatif (1-5)': 'Inisiatif (1-5)',
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
        if (val === undefined) val = row.student?.full_name || row.student?.email || ''
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
  const name = key === 'absensi' ? 'absensi' : key === 'jurnal' ? 'jurnal' : 'nilai'
  try {
    await downloadCsv(`/export/${name}`, `${name}_export.csv`)
  } catch (e) {
    alert('Export gagal: ' + e.message)
  }
}

function exportCurrentView() {
  if (!viewing.value) return
  exportCsv(viewing.value)
}

const logs = [
  { id: 1, message: 'Backup database harian berhasil disimpan', time: '23:00 - Hari ini', level: 'info' },
  { id: 2, message: '12 siswa belum absensi hari ini (15:00)', time: '15:05 - Hari ini', level: 'warning' },
  { id: 3, message: 'Upload ke Google Drive gagal untuk user 20230005 - quota penuh', time: '10:30 - Hari ini', level: 'error' },
  { id: 4, message: 'Export laporan absensi berhasil (format CSV)', time: '08:45 - Hari ini', level: 'info' },
  { id: 5, message: 'Periode PKL "Semester Genap 2025/2026" akan berakhir 7 hari lagi', time: '07:00 - Hari ini', level: 'warning' },
]
</script>
