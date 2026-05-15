<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Periode Tahun Pelajaran</h2>
        <p class="text-sm text-gray-500 mt-0.5">Atur tahun pelajaran dan aktivasi semester ganjil/genap</p>
      </div>
      <button @click="openCreateModal" class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark">
        <PlusIcon :size="16" /> Tambah Periode
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <!-- Period cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="p in periods" :key="p.id"
        :class="['bg-white rounded-2xl border-2 p-5 transition-all', p.is_active ? 'border-accent shadow-md shadow-accent/10' : 'border-gray-100 hover:shadow-sm']">
        <div class="flex items-start justify-between mb-3">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="text-xs font-mono text-gray-500">{{ p.tahun_pelajaran }}</span>
              <span :class="['px-2 py-0.5 rounded-full text-[11px] font-medium', p.semester === 'ganjil' ? 'bg-blue-50 text-blue-600' : 'bg-orange-50 text-orange-600']">
                {{ p.semester === 'ganjil' ? 'Ganjil' : 'Genap' }}
              </span>
            </div>
            <p class="text-xs text-gray-400">{{ formatDate(p.start_date) }} — {{ formatDate(p.end_date) }}</p>
          </div>
          <span v-if="p.is_active" class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-accent/10 text-accent text-[11px] font-bold">
            <CheckCircleIcon :size="12" /> Aktif
          </span>
        </div>

        <div class="flex gap-2">
          <button v-if="!p.is_active" @click="activatePeriod(p)" :disabled="saving"
            class="flex-1 py-2 rounded-xl text-xs font-medium bg-accent text-white hover:bg-accent-dark disabled:opacity-50 flex items-center justify-center gap-1">
            <PlayIcon :size="12" /> Aktifkan
          </button>
          <span v-else class="flex-1 py-2 rounded-xl text-xs font-medium bg-accent/5 text-accent text-center">
            Sedang Aktif
          </span>
          <button @click="openEditModal(p)" class="py-2 px-3 rounded-xl text-xs font-medium border border-gray-200 text-gray-600 hover:bg-gray-50">
            <PencilIcon :size="13" />
          </button>
          <button @click="confirmDelete(p)" class="py-2 px-3 rounded-xl text-xs font-medium border border-gray-200 text-gray-600 hover:bg-danger/5 hover:text-danger hover:border-danger/20">
            <TrashIcon :size="13" />
          </button>
        </div>
      </div>
      <div v-if="periods.length === 0" class="col-span-full text-center py-8 text-gray-400 text-sm">Belum ada periode</div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="closeModal">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full">
        <h3 class="text-lg font-bold text-gray-800 mb-4">{{ editing ? 'Edit Periode' : 'Tambah Periode' }}</h3>
        <form @submit.prevent="savePeriode" class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Tahun Pelajaran</label>
            <input v-model="form.tahun_pelajaran" placeholder="2025/2026" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Semester</label>
            <select v-model="form.semester" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
              <option value="ganjil">Ganjil</option>
              <option value="genap">Genap</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Tanggal Mulai</label>
              <input v-model="form.start_date" type="date" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Tanggal Selesai</label>
              <input v-model="form.end_date" type="date" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
          </div>
          <div class="flex gap-2 pt-2">
            <button type="button" @click="closeModal" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
            <button type="submit" :disabled="saving" class="flex-1 py-2.5 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark disabled:opacity-50">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirm -->
    <div v-if="showDelete" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="showDelete = null">
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-12 h-12 bg-danger/10 rounded-full flex items-center justify-center mx-auto mb-3">
          <AlertCircleIcon :size="24" class="text-danger" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus Periode?</h3>
        <p class="text-sm text-gray-500 mb-4">{{ deleteTarget?.tahun_pelajaran }} - {{ deleteTarget?.semester }}</p>
        <div class="flex gap-2">
          <button @click="showDelete = null" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="deletePeriode" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { PlusIcon, PencilIcon, TrashIcon, PlayIcon, CheckCircleIcon, AlertCircleIcon } from 'lucide-vue-next'
import { get, post, put, del } from '@/api'

const periods = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const editing = ref(null)
const showDelete = ref(null)
const deleteTarget = ref(null)

const form = reactive({
  tahun_pelajaran: '', semester: 'ganjil', start_date: '', end_date: ''
})

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchPeriods() {
  loading.value = true
  try {
    const res = await get('/admin/periode')
    periods.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editing.value = null
  Object.assign(form, { tahun_pelajaran: '', semester: 'ganjil', start_date: '', end_date: '' })
  showModal.value = true
}

function openEditModal(p) {
  editing.value = p
  Object.assign(form, {
    tahun_pelajaran: p.tahun_pelajaran,
    semester: p.semester,
    start_date: p.start_date?.split('T')[0] || '',
    end_date: p.end_date?.split('T')[0] || ''
  })
  showModal.value = true
}

function closeModal() { showModal.value = false; editing.value = null }

async function savePeriode() {
  saving.value = true
  try {
    if (editing.value) {
      await put('/admin/periode/' + editing.value.id, { ...form })
    } else {
      await post('/admin/periode', { ...form })
    }
    closeModal()
    fetchPeriods()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

async function activatePeriod(p) {
  saving.value = true
  try {
    await put('/admin/periode/' + p.id + '/activate')
    fetchPeriods()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

function confirmDelete(p) { deleteTarget.value = p; showDelete.value = true }

async function deletePeriode() {
  saving.value = true
  try {
    await del('/admin/periode/' + deleteTarget.value.id)
    showDelete.value = null
    fetchPeriods()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchPeriods)
</script>
