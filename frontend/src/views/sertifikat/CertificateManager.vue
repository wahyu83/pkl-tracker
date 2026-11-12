<template>
  <div>
    <div class="flex items-center justify-between mb-6 flex-wrap gap-3">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Sertifikat PKL</h2>
        <p class="text-sm text-gray-500 mt-0.5">Cetak sertifikat praktik kerja lapangan siswa</p>
      </div>
      <button
        @click="printAll"
        :disabled="!filtered.length"
        class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark disabled:opacity-50"
      >
        <PrinterIcon :size="16" /> Cetak Semua ({{ filtered.length }})
      </button>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <template v-else>
      <!-- Info sekolah & periode -->
      <div class="bg-white rounded-2xl border border-gray-100 p-4 mb-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-gray-50 flex items-center justify-center flex-shrink-0 overflow-hidden">
          <img v-if="school.logo_url" :src="school.logo_url" alt="Logo" class="w-full h-full object-contain" />
          <BuildingIcon v-else :size="22" class="text-gray-300" />
        </div>
        <div class="min-w-0">
          <p class="text-sm font-semibold text-gray-800 truncate">{{ school.name || 'Nama sekolah belum diatur' }}</p>
          <p class="text-xs text-gray-500 truncate">{{ school.address || school.city || 'Lengkapi di Pengaturan Sekolah' }}</p>
        </div>
      </div>

      <!-- Filter -->
      <div class="bg-white rounded-2xl border border-gray-100 p-4 mb-4 flex flex-wrap gap-3 items-center">
        <select v-model="selectedPeriod" class="px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none bg-white min-w-[200px]">
          <option value="">Periode Aktif</option>
          <option v-for="p in periods" :key="p.id" :value="p.id">
            {{ p.tahun_pelajaran }} {{ p.semester === 'ganjil' ? 'Ganjil' : 'Genap' }} {{ p.is_active ? '(Aktif)' : '' }}
          </option>
        </select>
        <div class="flex-1 min-w-[180px]">
          <input v-model="search" type="text" placeholder="Cari nama atau NIS..." class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
        </div>
      </div>

      <!-- Daftar siswa -->
      <div v-if="filtered.length === 0" class="bg-white rounded-2xl border border-gray-100 text-center py-10">
        <AwardIcon :size="28" class="text-gray-300 mx-auto mb-2" />
        <p class="text-sm text-gray-400">Tidak ada siswa yang bisa dicetak sertifikatnya.</p>
      </div>

      <div v-else class="bg-white rounded-2xl border border-gray-100 divide-y divide-gray-50 overflow-hidden">
        <div
          v-for="c in filtered"
          :key="c.student.id"
          class="p-4 flex items-center gap-3 hover:bg-gray-50/50"
        >
          <div class="w-9 h-9 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary flex-shrink-0">
            {{ c.student.full_name.charAt(0) }}
          </div>

          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-800 truncate">{{ c.student.full_name }}</p>
            <p class="text-xs text-gray-400 truncate">
              NIS {{ c.student.nis_nip_nik }} • {{ c.dudi?.company_name || 'Belum ada DUDI' }}
            </p>
          </div>

          <div class="text-right flex-shrink-0 hidden sm:block">
            <p class="text-sm font-bold" :class="c.score.assessed ? 'text-accent' : 'text-gray-300'">
              {{ c.score.assessed ? Number(c.score.final_score).toFixed(2) : '-' }}
            </p>
            <p class="text-[10px] text-gray-400">
              {{ c.score.assessed ? 'Nilai ' + (c.score.final_grade || '-') : 'Belum dinilai' }}
            </p>
          </div>

          <button
            @click="printOne(c)"
            class="flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-medium text-primary bg-primary/5 hover:bg-primary/10 flex-shrink-0"
          >
            <PrinterIcon :size="14" /> Cetak
          </button>
        </div>
      </div>

      <p class="text-xs text-gray-400 mt-3">
        Halaman sertifikat terbuka di tab baru. Gunakan tombol Cetak atau Ctrl+P, lalu pilih
        <span class="font-medium">Simpan sebagai PDF</span>. Ukuran kertas: A4 landscape.
      </p>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { PrinterIcon, AwardIcon, BuildingIcon } from 'lucide-vue-next'
import { get } from '@/api'

const loading = ref(true)
const search = ref('')
const periods = ref([])
const selectedPeriod = ref('')
const school = ref({})
const certificates = ref([])

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return certificates.value
  return certificates.value.filter(c =>
    c.student.full_name.toLowerCase().includes(q) ||
    c.student.nis_nip_nik.toLowerCase().includes(q)
  )
})

async function fetchPeriods() {
  try {
    const res = await get('/periodes')
    periods.value = res.data || []
  } catch (e) { /* abaikan */ }
}

async function fetchCertificates() {
  loading.value = true
  try {
    const qs = selectedPeriod.value ? `?periode_id=${selectedPeriod.value}` : ''
    const res = await get('/certificate' + qs)
    certificates.value = res.data || []
    school.value = res.school || {}
  } catch (e) {
    certificates.value = []
  } finally {
    loading.value = false
  }
}

function printUrl(params) {
  const qs = new URLSearchParams(params).toString()
  return `/sertifikat/cetak${qs ? '?' + qs : ''}`
}

function printOne(c) {
  const params = { student_id: c.student.id }
  if (selectedPeriod.value) params.periode_id = selectedPeriod.value
  window.open(printUrl(params), '_blank')
}

function printAll() {
  const params = {}
  if (selectedPeriod.value) params.periode_id = selectedPeriod.value
  window.open(printUrl(params), '_blank')
}

watch(selectedPeriod, fetchCertificates)

onMounted(async () => {
  await fetchPeriods()
  await fetchCertificates()
})
</script>
