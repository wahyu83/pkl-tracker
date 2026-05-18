<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Monitoring Absensi</h2>
      <p class="text-sm text-gray-500 mt-0.5">Pantau kehadiran siswa bimbingan</p>
    </div>

    <div class="bg-white rounded-2xl p-4 border border-gray-100 mb-4 flex flex-wrap gap-3 items-center">
      <div class="flex-1 min-w-[180px]">
        <input
          v-model="search"
          type="text"
          placeholder="Cari nama siswa..."
          class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
        />
      </div>
      <select v-model="filterStatus" class="px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none bg-white">
        <option value="">Semua Status</option>
        <option value="hadir">Hadir</option>
        <option value="terlambat">Terlambat</option>
        <option value="izin">Izin</option>
        <option value="sakit">Sakit</option>
      </select>
      <input v-model="filterDate" type="date" class="px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none" />
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Siswa</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Tanggal</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Jam</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Tipe</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Status</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Validasi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredAbsensi.length === 0">
              <td colspan="6" class="px-4 py-8 text-center text-gray-400 text-sm">Belum ada data absensi</td>
            </tr>
            <tr v-for="a in filteredAbsensi" :key="a.id" class="border-b border-gray-50 hover:bg-gray-50/50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                    {{ (a.student?.full_name || '?').charAt(0) }}
                  </div>
                  <span class="text-sm font-medium text-gray-800">{{ a.student?.full_name || '-' }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ formatDate(a.timestamp) }}</td>
              <td class="px-4 py-3 text-sm text-gray-600 font-mono">{{ formatTime(a.timestamp) }}</td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ a.type === 'masuk' ? 'Masuk' : 'Pulang' }}</td>
              <td class="px-4 py-3">
                <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium', statusStyle(a.status)]">
                  {{ statusLabel(a.status) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <button
                  v-if="!a.is_verified && a.type === 'masuk'"
                  @click="verifyAbsensi(a)"
                  :disabled="verifying === a.id"
                  class="text-xs text-primary font-medium hover:underline disabled:opacity-50"
                >
                  {{ verifying === a.id ? '...' : 'Verifikasi' }}
                </button>
                <span v-else-if="a.is_verified" class="inline-flex items-center gap-1 text-xs font-medium text-accent">
                  <CheckCircleIcon :size="14" />
                  Valid
                </span>
                <span v-else class="text-xs text-gray-400">-</span>
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
import { CheckCircleIcon } from 'lucide-vue-next'
import { get, put } from '../../api'

const search = ref('')
const filterStatus = ref('')
const filterDate = ref('')
const loading = ref(true)
const verifying = ref(null)
const absensiData = ref([])

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatTime(d) {
  if (!d) return '-'
  return new Date(d).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

const filteredAbsensi = computed(() => {
  return absensiData.value.filter(a => {
    const name = (a.student?.full_name || '').toLowerCase()
    const matchSearch = name.includes(search.value.toLowerCase())
    const matchStatus = !filterStatus.value || a.status === filterStatus.value
    const matchDate = !filterDate.value || (a.timestamp && a.timestamp.startsWith(filterDate.value))
    return matchSearch && matchStatus && matchDate
  })
})

async function fetchData() {
  loading.value = true
  try {
    const res = await get('/absensi/history')
    absensiData.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function verifyAbsensi(a) {
  verifying.value = a.id
  try {
    await put(`/absensi/${a.id}/verify`)
    a.is_verified = true
  } catch (e) {
    alert('Gagal verifikasi: ' + e.message)
  } finally {
    verifying.value = null
  }
}

function statusLabel(status) {
  return { hadir: 'Hadir', terlambat: 'Terlambat', izin: 'Izin', sakit: 'Sakit' }[status] || status
}

function statusStyle(status) {
  const map = {
    hadir: 'bg-accent/10 text-accent',
    terlambat: 'bg-warning/10 text-warning',
    izin: 'bg-info/10 text-info',
    sakit: 'bg-purple-100 text-purple-600'
  }
  return map[status] || 'bg-gray-100 text-gray-500'
}

onMounted(fetchData)
</script>
