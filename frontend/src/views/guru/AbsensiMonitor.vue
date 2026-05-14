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
        <option value="alpa">Alpa</option>
      </select>
      <input v-model="filterDate" type="date" class="px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none" />
      <button class="px-4 py-2 rounded-xl bg-primary text-white text-sm font-medium hover:bg-primary-light transition-colors">
        Filter
      </button>
    </div>

    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Siswa</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">DUDI</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Tanggal</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Jam</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Lokasi</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Status</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Validasi</th>
              <th class="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Bukti</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in filteredAbsensi" :key="a.id" class="border-b border-gray-50 hover:bg-gray-50/50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                    {{ a.student.charAt(0) }}
                  </div>
                  <span class="text-sm font-medium text-gray-800">{{ a.student }}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ a.dudi }}</td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ a.date }}</td>
              <td class="px-4 py-3 text-sm text-gray-600 font-mono">{{ a.time }}</td>
              <td class="px-4 py-3 text-sm text-gray-500 max-w-[120px] truncate">{{ a.location }}</td>
              <td class="px-4 py-3">
                <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium', statusStyle(a.status)]">
                  {{ statusLabel(a.status) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span :class="['inline-flex items-center gap-1 text-xs font-medium', a.verified ? 'text-accent' : 'text-danger']">
                  <CheckCircleIcon v-if="a.verified" :size="14" />
                  <AlertCircleIcon v-else :size="14" />
                  {{ a.verified ? 'Valid' : 'Di luar radius' }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <button
                  v-if="a.photo"
                  class="text-xs text-primary font-medium hover:underline"
                  @click="showPhoto = a.photo"
                >
                  Lihat Foto
                </button>
                <span v-else class="text-xs text-gray-400">-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Photo Modal -->
    <div
      v-if="showPhoto"
      class="fixed inset-0 z-50 bg-black/60 flex items-center justify-center p-4"
      @click.self="showPhoto = null"
    >
      <div class="bg-white rounded-2xl p-4 max-w-sm w-full">
        <img :src="showPhoto" class="w-full rounded-xl" alt="Bukti Absensi" />
        <button
          @click="showPhoto = null"
          class="mt-3 w-full py-2 bg-gray-100 rounded-xl text-sm text-gray-600 hover:bg-gray-200 transition-colors"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { CheckCircleIcon, AlertCircleIcon } from 'lucide-vue-next'

const search = ref('')
const filterStatus = ref('')
const filterDate = ref('')
const showPhoto = ref(null)

const absensiData = [
  { id: 1, student: 'Ahmad Rizky', dudi: 'PT. Teknologi Maju', date: '14 Mei 2026', time: '07:45', location: '-6.2088, 106.8456', status: 'hadir', verified: true, photo: 'https://placehold.co/400x300/e2e8f0/64748b?text=Selfie+Absensi' },
  { id: 2, student: 'Siti Nurhaliza', dudi: 'PT. Sejahtera Abadi', date: '14 Mei 2026', time: '08:15', location: '-6.2297, 106.8243', status: 'terlambat', verified: true, photo: 'https://placehold.co/400x300/e2e8f0/64748b?text=Selfie+Absensi' },
  { id: 3, student: 'Dian Permata', dudi: 'CV. Kreatif Digital', date: '14 Mei 2026', time: '-', location: '-6.9175, 107.6191', status: 'izin', verified: false, photo: null },
  { id: 4, student: 'Rudi Hartono', dudi: 'UD. Mandiri Jaya', date: '14 Mei 2026', time: '07:30', location: '-7.2575, 112.7521', status: 'hadir', verified: true, photo: 'https://placehold.co/400x300/e2e8f0/64748b?text=Selfie+Absensi' },
  { id: 5, student: 'Maya Sari', dudi: 'PT. Inovasi Nusantara', date: '14 Mei 2026', time: '-', location: '-6.1823, 106.8229', status: 'sakit', verified: false, photo: null },
  { id: 6, student: 'Bambang Kusumo', dudi: 'PT. Teknologi Maju', date: '14 Mei 2026', time: '11:20', location: '-6.5000, 106.7000', status: 'hadir', verified: false, photo: 'https://placehold.co/400x300/e2e8f0/64748b?text=Selfie+Absensi' },
]

const filteredAbsensi = computed(() => {
  return absensiData.filter(a => {
    const matchSearch = a.student.toLowerCase().includes(search.value.toLowerCase())
    const matchStatus = !filterStatus.value || a.status === filterStatus.value
    return matchSearch && matchStatus
  })
})

function statusLabel(status) {
  return { hadir: 'Hadir', terlambat: 'Terlambat', izin: 'Izin', sakit: 'Sakit', alpa: 'Alpa' }[status]
}

function statusStyle(status) {
  const map = {
    hadir: 'bg-accent/10 text-accent',
    terlambat: 'bg-warning/10 text-warning',
    izin: 'bg-info/10 text-info',
    sakit: 'bg-purple-100 text-purple-600',
    alpa: 'bg-danger/10 text-danger'
  }
  return map[status] || 'bg-gray-100 text-gray-500'
}
</script>
