<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Rekap & Laporan</h2>
      <p class="text-sm text-gray-500 mt-0.5">Export dan lihat laporan PKL per periode</p>
    </div>

    <!-- Periode Selector -->
    <div class="flex items-center gap-3 mb-4">
      <label class="text-sm font-medium text-gray-600">Periode:</label>
      <select v-model="selectedPeriod" class="px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none min-w-[240px]">
        <option value="">Semua</option>
        <option v-for="p in periods" :key="p.id" :value="p.id">
          {{ p.tahun_pelajaran }} {{ p.semester === 'ganjil' ? 'Ganjil' : 'Genap' }} {{ p.is_active ? '(Aktif)' : '' }}
        </option>
      </select>
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
            :disabled="!report.key"
            class="flex-1 py-2 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors disabled:opacity-30"
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

    <!-- Ringkasan Periode -->
    <div v-if="periodSummary" class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <h3 class="font-semibold text-gray-800 mb-4">Ringkasan Periode {{ activeTitle }}</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="bg-primary/5 rounded-xl p-4 text-center">
          <p class="text-2xl font-bold text-primary">{{ periodSummary.total_hadir }}</p>
          <p class="text-xs text-gray-500 mt-1">Hadir</p>
        </div>
        <div class="bg-warning/5 rounded-xl p-4 text-center">
          <p class="text-2xl font-bold text-warning">{{ periodSummary.total_terlambat }}</p>
          <p class="text-xs text-gray-500 mt-1">Terlambat</p>
        </div>
        <div class="bg-info/5 rounded-xl p-4 text-center">
          <p class="text-2xl font-bold text-info">{{ periodSummary.total_izin }}</p>
          <p class="text-xs text-gray-500 mt-1">Izin</p>
        </div>
        <div class="bg-gray-100 rounded-xl p-4 text-center">
          <p class="text-2xl font-bold text-gray-600">{{ periodSummary.total_sakit }}</p>
          <p class="text-xs text-gray-500 mt-1">Sakit</p>
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
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { FileText, ClipboardCheck, BookOpen, Award, DownloadIcon, XIcon } from 'lucide-vue-next'
import { downloadCsv, get } from '../../api'

const periods = ref([])
const selectedPeriod = ref('')
const periodSummary = ref(null)
const viewing = ref(null)
const loadingView = ref(false)
const viewData = ref([])

const reports = [
  { key: 'absensi', title: 'Rekap Absensi', desc: 'Laporan kehadiran siswa per periode PKL', icon: ClipboardCheck, bg: 'bg-primary/10', color: 'text-primary' },
  { key: 'jurnal', title: 'Rekap Jurnal', desc: 'Laporan jurnal harian siswa PKL', icon: BookOpen, bg: 'bg-accent/10', color: 'text-accent' },
  { key: 'nilai', title: 'Rekap Nilai', desc: 'Laporan penilaian dari Instruktur DUDI per siswa', icon: Award, bg: 'bg-warning/10', color: 'text-warning' },
  { key: 'absensi', title: 'Export Absensi CSV', desc: 'Download data absensi per periode', icon: FileText, bg: 'bg-gray-100', color: 'text-gray-600' },
  { key: 'nilai', title: 'Export Nilai CSV', desc: 'Download data penilaian lengkap', icon: FileText, bg: 'bg-gray-100', color: 'text-gray-600' },
  { key: 'jurnal', title: 'Export Jurnal CSV', desc: 'Download data jurnal harian', icon: FileText, bg: 'bg-gray-100', color: 'text-gray-600' },
]

const activeTitle = computed(() => {
  if (!selectedPeriod.value) return ''
  const p = periods.value.find(p => p.id === selectedPeriod.value)
  return p ? `${p.tahun_pelajaran} ${p.semester}` : ''
})

const viewingTitle = computed(() => {
  const t = { absensi: 'Rekap Absensi', jurnal: 'Rekap Jurnal', nilai: 'Rekap Nilai' }
  return t[viewing.value] || 'Data'
})

const viewColumns = computed(() => {
  if (viewData.value.length > 0) return Object.keys(viewData.value[0])
  return []
})

const absensiKeys = {
  'Tanggal': 'Tanggal', 'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS', 'Status': 'Status',
  'Latitude': 'Latitude', 'Longitude': 'Longitude', 'Terverifikasi': 'Terverifikasi'
}

const jurnalKeys = {
  'Tanggal': 'Tanggal', 'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS',
  'Kegiatan': 'Kegiatan', 'Refleksi': 'Refleksi'
}

const nilaiKeys = {
  'Nama Siswa': 'Nama Siswa', 'NIS': 'NIS', 'Kehadiran (%)': 'Kehadiran (%)',
  'Disiplin (1-5)': 'Disiplin (1-5)', 'Tanggung Jawab (1-5)': 'Tanggung Jawab (1-5)',
  'Kerjasama (1-5)': 'Kerjasama (1-5)', 'Inisiatif (1-5)': 'Inisiatif (1-5)',
  'Nilai Akhir': 'Nilai Akhir', 'Grade': 'Grade'
}

async function fetchPeriods() {
  try {
    const res = await get('/admin/periode')
    periods.value = res.data || []
    const active = periods.value.find(p => p.is_active)
    if (active) selectedPeriod.value = active.id
  } catch (e) { /* ignore */ }
}

async function fetchSummary() {
  if (!selectedPeriod.value) { periodSummary.value = null; return }
  try {
    const data = await get(`/report/absensi?periode_id=${selectedPeriod.value}`)
    periodSummary.value = data.summary || null
  } catch (e) { periodSummary.value = null }
}

watch(selectedPeriod, fetchSummary)

async function viewReport(key) {
  if (!key) return
  viewing.value = key
  loadingView.value = true
  viewData.value = []
  try {
    let url = `/report/${key}`
    if (key === 'absensi' && selectedPeriod.value) url += `?periode_id=${selectedPeriod.value}`
    const data = await get(url)
    const list = data.data || []
    const keyMap = key === 'absensi' ? absensiKeys : key === 'jurnal' ? jurnalKeys : nilaiKeys
    viewData.value = list.map(row => {
      const mapped = {}
      for (const [label, field] of Object.entries(keyMap)) {
        let val = row[field]
        if (val === undefined) val = row[field.toLowerCase().replace(/ /g, '_')]
        if (val === undefined && row.student) val = row.student.full_name || row.student.email || ''
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
  let url = `/export/${name}`
  if (key === 'absensi' && selectedPeriod.value) url += `?periode_id=${selectedPeriod.value}`
  try {
    await downloadCsv(url, `${name}_export.csv`)
  } catch (e) {
    alert('Export gagal: ' + e.message)
  }
}

function exportCurrentView() {
  if (!viewing.value) return
  exportCsv(viewing.value)
}

onMounted(fetchPeriods)
</script>
