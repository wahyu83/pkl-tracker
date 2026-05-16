<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Manajemen DUDI</h2>
        <p class="text-sm text-gray-500 mt-0.5">Kelola data perusahaan/world usaha tempat PKL</p>
      </div>
      <div class="flex items-center gap-2">
        <button @click="showImport = true" class="flex items-center gap-1.5 px-4 py-2 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">
          <UploadIcon :size="16" /> Import CSV
        </button>
        <button @click="openCreateModal" class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark">
          <PlusIcon :size="16" /> Tambah
        </button>
      </div>
    </div>

    <!-- Search -->
    <div class="flex items-center gap-3 mb-4">
      <div class="flex-1 relative">
        <SearchIcon :size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input v-model="search" @input="fetchDudi" placeholder="Cari nama perusahaan..." class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
      </div>
      <select v-model="jurusanFilter" @change="fetchDudi" class="px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
        <option value="">Semua Jurusan</option>
        <option value="RPL">RPL</option>
        <option value="TKJ">TKJ</option>
        <option value="MM">MM</option>
        <option value="AKL">AKL</option>
        <option value="OTKP">OTKP</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <!-- Card Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="d in dudiList" :key="d.id" class="bg-white rounded-2xl border border-gray-100 p-5 hover:shadow-sm transition-shadow">
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="font-semibold text-gray-800">{{ d.company_name }}</h3>
            <p class="text-xs text-gray-500 mt-0.5">{{ d.pic_name || 'Belum ada PIC' }}</p>
            <p v-if="d.jurusan" class="text-xs text-gray-400 mt-0.5">Jurusan: {{ d.jurusan }}</p>
          </div>
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium bg-blue-50 text-blue-600">
            <UsersIcon :size="12" /> {{ d.student_count || 0 }}
          </span>
        </div>
        <div class="space-y-1 mb-4">
          <p class="text-xs text-gray-500 flex items-start gap-1"><MapPinIcon :size="12" class="mt-0.5 flex-shrink-0" /> {{ d.address || 'Alamat belum diisi' }}</p>
          <p v-if="d.phone" class="text-xs text-gray-500 flex items-center gap-1"><PhoneIcon :size="12" /> {{ d.phone }}</p>
          <p class="text-[11px] text-gray-400">Radius: {{ d.radius_allowed }}m | Lat: {{ d.latitude?.toFixed(4) }} Lng: {{ d.longitude?.toFixed(4) }}</p>
        </div>
        <div class="flex gap-2">
          <button @click="openEditModal(d)" class="flex-1 py-2 rounded-xl text-xs font-medium border border-gray-200 text-gray-600 hover:bg-gray-50 flex items-center justify-center gap-1">
            <PencilIcon :size="13" /> Edit
          </button>
          <button @click="confirmDelete(d)" class="py-2 px-3 rounded-xl text-xs font-medium border border-gray-200 text-gray-600 hover:bg-danger/5 hover:text-danger hover:border-danger/20 flex items-center gap-1">
            <TrashIcon :size="13" />
          </button>
        </div>
      </div>
      <div v-if="dudiList.length === 0" class="col-span-full text-center py-8 text-gray-400 text-sm">Tidak ada data DUDI</div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="closeModal">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-bold text-gray-800 mb-4">{{ editingDudi ? 'Edit DUDI' : 'Tambah DUDI' }}</h3>
        <form @submit.prevent="saveDudi" class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Nama Perusahaan</label>
            <input v-model="form.company_name" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Alamat</label>
            <textarea v-model="form.address" rows="2" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Latitude</label>
              <input v-model.number="form.latitude" type="number" step="any" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Longitude</label>
              <input v-model.number="form.longitude" type="number" step="any" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Radius (meter)</label>
              <input v-model.number="form.radius_allowed" type="number" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Telepon</label>
              <input v-model="form.phone" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
            </div>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Nama PIC</label>
            <input v-model="form.pic_name" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Jurusan</label>
            <select v-model="form.jurusan" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
              <option value="">-- Pilih Jurusan --</option>
              <option value="RPL">RPL</option>
              <option value="TKJ">TKJ</option>
              <option value="MM">MM</option>
              <option value="AKL">AKL</option>
              <option value="OTKP">OTKP</option>
            </select>
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
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus DUDI?</h3>
        <p class="text-sm text-gray-500 mb-4">{{ deleteTarget?.company_name }} akan dihapus permanen.</p>
        <div class="flex gap-2">
          <button @click="showDelete = null" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="deleteDudi" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>

    <!-- CSV Import -->
    <CsvImport v-if="showImport" :show="showImport" title="Import DUDI" label="File CSV DUDI" endpoint="/import/dudi"
      :csv-columns="['company_name (Nama Perusahaan)', 'address', 'latitude', 'longitude', 'radius_allowed', 'pic_name', 'phone', 'jurusan']"
      sample-csv="company_name,address,latitude,longitude,radius_allowed,pic_name,phone,jurusan
PT. Teknologi Maju,Jl. Sudirman No. 123,-6.2088,106.8456,500,Hendra Gunawan,021-5551234,RPL"
      @close="showImport = false" @done="showImport = false; fetchDudi()" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { PlusIcon, SearchIcon, UploadIcon, PencilIcon, TrashIcon, AlertCircleIcon, UsersIcon, MapPinIcon, PhoneIcon } from 'lucide-vue-next'
import { get, post, put, del } from '@/api'
import CsvImport from '@/components/CsvImport.vue'

const dudiList = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const jurusanFilter = ref('')
const showModal = ref(false)
const editingDudi = ref(null)
const showDelete = ref(null)
const deleteTarget = ref(null)
const showImport = ref(false)

const form = reactive({
  company_name: '', address: '', latitude: 0, longitude: 0, radius_allowed: 500, pic_name: '', phone: '', jurusan: ''
})

async function fetchDudi() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (search.value) params.set('search', search.value)
    if (jurusanFilter.value) params.set('jurusan', jurusanFilter.value)
    const res = await get('/admin/dudi?' + params.toString())
    dudiList.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingDudi.value = null
  Object.assign(form, { company_name: '', address: '', latitude: 0, longitude: 0, radius_allowed: 500, pic_name: '', phone: '', jurusan: '' })
  showModal.value = true
}

function openEditModal(d) {
  editingDudi.value = d
  Object.assign(form, {
    company_name: d.company_name, address: d.address || '', latitude: d.latitude || 0,
    longitude: d.longitude || 0, radius_allowed: d.radius_allowed || 500,
    pic_name: d.pic_name || '', phone: d.phone || '', jurusan: d.jurusan || ''
  })
  showModal.value = true
}

function closeModal() { showModal.value = false; editingDudi.value = null }

async function saveDudi() {
  saving.value = true
  try {
    if (editingDudi.value) {
      await put('/admin/dudi/' + editingDudi.value.id, { ...form })
    } else {
      await post('/admin/dudi', { ...form })
    }
    closeModal()
    fetchDudi()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

function confirmDelete(d) { deleteTarget.value = d; showDelete.value = true }

async function deleteDudi() {
  saving.value = true
  try {
    await del('/admin/dudi/' + deleteTarget.value.id)
    showDelete.value = null
    fetchDudi()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchDudi)
</script>
