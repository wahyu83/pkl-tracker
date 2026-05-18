<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Manajemen Pengguna</h2>
        <p class="text-sm text-gray-500 mt-0.5">Kelola akun siswa, guru, instruktur DUDI, dan admin</p>
      </div>
      <div class="flex items-center gap-2">
        <div class="relative" ref="importDropdownRef">
          <button @click="showImportDropdown = !showImportDropdown" class="flex items-center gap-1.5 px-4 py-2 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">
            <UploadIcon :size="16" /> Import CSV
          </button>
          <div v-if="showImportDropdown" class="absolute right-0 mt-2 w-64 bg-white rounded-xl shadow-lg border border-gray-100 z-10 py-1">
            <button v-for="(cfg, key) in importConfigs" :key="key" @click="openImport(key)" class="w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2">
              <FileTextIcon :size="14" class="text-gray-400" /> {{ cfg.title }}
            </button>
          </div>
        </div>
        <button @click="openCreateModal" class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark">
          <PlusIcon :size="16" /> Tambah
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-3 mb-4">
      <div class="flex-1 relative">
        <SearchIcon :size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input v-model="search" @input="fetchUsers" placeholder="Cari nama, email, NIS..." class="w-full pl-9 pr-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
      </div>
      <select v-model="roleFilter" @change="fetchUsers" class="px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
        <option value="">Semua Role</option>
        <option value="student">Siswa</option>
        <option value="teacher">Guru</option>
        <option value="dudi">Instruktur DUDI</option>
        <option value="admin">Admin</option>
        <option value="admin_jurusan">Admin Jurusan</option>
      </select>
      <select v-model="jurusanFilter" @change="fetchUsers" class="px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
        <option value="">Semua Jurusan</option>
        <option v-for="j in jurusanOptions" :key="j.kode" :value="j.kode">{{ j.kode }} - {{ j.nama }}</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <!-- Table -->
    <div v-else class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <!-- Bulk toolbar -->
      <div v-if="selected.length > 0" class="flex items-center gap-3 px-4 py-2 bg-danger/5 border-b border-danger/10">
        <span class="text-sm text-danger font-medium">{{ selected.length }} dipilih</span>
        <button @click="showBulkDelete = true" class="px-3 py-1.5 bg-danger text-white rounded-lg text-xs font-medium hover:bg-red-600">Hapus Massal</button>
        <button @click="selected = []" class="px-3 py-1.5 border border-gray-200 rounded-lg text-xs text-gray-600 hover:bg-gray-50">Batal</button>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-100 bg-gray-50/50 text-left">
              <th class="px-4 py-3 w-10">
                <input type="checkbox" :checked="allSelected" @change="toggleAll" class="rounded border-gray-300" />
              </th>
              <th class="px-4 py-3 font-medium text-gray-500">Nama</th>
              <th class="px-4 py-3 font-medium text-gray-500">Role</th>
              <th class="px-4 py-3 font-medium text-gray-500">Jurusan</th>
              <th class="px-4 py-3 font-medium text-gray-500">Email</th>
              <th class="px-4 py-3 font-medium text-gray-500">NIS/NIP/NIK</th>
              <th class="px-4 py-3 font-medium text-gray-500 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="u in users" :key="u.id" :class="['hover:bg-gray-50/50', selected.includes(u.id) ? 'bg-primary/5' : '']">
              <td class="px-4 py-3">
                <input type="checkbox" :checked="selected.includes(u.id)" @change="toggleSelect(u.id)" class="rounded border-gray-300" />
              </td>
              <td class="px-4 py-3 font-medium text-gray-800">{{ u.full_name }}</td>
              <td class="px-4 py-3">
                <span :class="['inline-flex px-2 py-0.5 rounded-full text-[11px] font-medium', roleBadge(u.role)]">{{ roleLabel(u.role) }}</span>
              </td>
              <td class="px-4 py-3 text-gray-600 text-xs">{{ u.jurusan || '-' }}</td>
              <td class="px-4 py-3 text-gray-600">{{ u.email }}</td>
              <td class="px-4 py-3 text-gray-600 font-mono text-xs">{{ u.nis_nip_nik }}</td>
              <td class="px-4 py-3 text-center">
                <div class="flex items-center justify-center gap-1">
                  <button @click="openEditModal(u)" class="p-1.5 rounded-lg text-gray-400 hover:text-primary hover:bg-primary/10"><PencilIcon :size="14" /></button>
                  <button @click="confirmDelete(u)" class="p-1.5 rounded-lg text-gray-400 hover:text-danger hover:bg-danger/10"><TrashIcon :size="14" /></button>
                </div>
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="7" class="px-4 py-8 text-center text-gray-400">Tidak ada data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="closeModal">
      <div class="bg-white rounded-2xl p-6 max-w-md w-full max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-bold text-gray-800 mb-4">{{ editingUser ? 'Edit Pengguna' : 'Tambah Pengguna' }}</h3>
        <form @submit.prevent="saveUser" class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Nama Lengkap</label>
            <input v-model="form.full_name" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Email</label>
            <input v-model="form.email" type="email" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">NIS/NIP/NIK</label>
            <input v-model="form.nis_nip_nik" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Role</label>
            <select v-model="form.role" required class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
              <option value="student">Siswa</option>
              <option value="teacher">Guru</option>
        <option value="dudi">Instruktur DUDI</option>
              <option value="admin">Admin</option>
              <option value="admin_jurusan">Admin Jurusan</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Jurusan</label>
            <select v-model="form.jurusan" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
              <option value="">-- Pilih Jurusan --</option>
              <option v-for="j in jurusanOptions" :key="j.kode" :value="j.kode">{{ j.kode }} - {{ j.nama }}</option>
            </select>
          </div>
          <div v-if="form.role === 'student' || form.role === 'dudi'">
            <label class="block text-xs font-medium text-gray-600 mb-1">DUDI</label>
            <select v-model="form.dudi_id" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none">
              <option value="">-- Pilih DUDI --</option>
              <option v-for="d in dudiList" :key="d.id" :value="d.id">{{ d.company_name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Password {{ editingUser ? '(kosongkan jika tidak diubah)' : '' }}</label>
            <input v-model="form.password" type="password" :required="!editingUser" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none" />
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
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus Pengguna?</h3>
        <p class="text-sm text-gray-500 mb-4">{{ deleteTarget?.full_name }} akan dihapus permanen.</p>
        <div class="flex gap-2">
          <button @click="showDelete = null" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="deleteUser" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>

    <!-- Bulk Delete Confirm -->
    <div v-if="showBulkDelete" class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" @click.self="showBulkDelete = false">
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-12 h-12 bg-danger/10 rounded-full flex items-center justify-center mx-auto mb-3">
          <AlertCircleIcon :size="24" class="text-danger" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-1">Hapus {{ selected.length }} Pengguna?</h3>
        <p class="text-sm text-gray-500 mb-4">Data akan dihapus permanen dan tidak bisa dikembalikan.</p>
        <div class="flex gap-2">
          <button @click="showBulkDelete = false" class="flex-1 py-2.5 border border-gray-200 rounded-xl text-sm text-gray-600 hover:bg-gray-50">Batal</button>
          <button @click="bulkDelete" :disabled="saving" class="flex-1 py-2.5 bg-danger text-white rounded-xl text-sm font-medium hover:bg-red-600 disabled:opacity-50">{{ saving ? 'Menghapus...' : 'Hapus Semua' }}</button>
        </div>
      </div>
    </div>

    <!-- CsvImport modal -->
    <CsvImport v-if="importKey" :show="!!importKey" :title="importConfigs[importKey]?.title" :label="importConfigs[importKey]?.label" :endpoint="importConfigs[importKey]?.endpoint" :csv-columns="importConfigs[importKey]?.columns" :sample-csv="importConfigs[importKey]?.sample" @close="importKey = null" @done="importKey = null; fetchUsers()" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { PlusIcon, SearchIcon, UploadIcon, FileTextIcon, PencilIcon, TrashIcon, AlertCircleIcon } from 'lucide-vue-next'
import { get, post, put, del } from '@/api'
import CsvImport from '@/components/CsvImport.vue'

const users = ref([])
const dudiList = ref([])
const jurusanOptions = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const roleFilter = ref('')
const jurusanFilter = ref('')
const showModal = ref(false)
const editingUser = ref(null)
const showDelete = ref(null)
const deleteTarget = ref(null)
const showImportDropdown = ref(false)
const importKey = ref(null)
const selected = ref([])
const showBulkDelete = ref(false)

const form = reactive({
  full_name: '', email: '', nis_nip_nik: '', role: 'student', password: '', jurusan: '', dudi_id: ''
})

const importConfigs = {
  siswa: {
    title: 'Import Siswa', label: 'File CSV Siswa', endpoint: '/import/siswa',
    columns: ['full_name (Nama Lengkap)', 'email', 'nis (NIS)', 'password', 'jurusan', 'dudi_nik'],
    sample: 'full_name,email,nis,password,jurusan,dudi_nik\nAhmad Rizky,ahmad@sekolah.sch.id,20230001,rahasia123,RPL,D-001'
  },
  guru: {
    title: 'Import Guru', label: 'File CSV Guru', endpoint: '/import/guru',
    columns: ['full_name (Nama Lengkap)', 'email', 'nip (NIP)', 'password', 'jurusan'],
    sample: 'full_name,email,nip,password,jurusan\nBudi Santoso,budi@sekolah.sch.id,198501012025011001,rahasia123,RPL'
  },
  'instruktur-dudi': {
    title: 'Import Instruktur DUDI', label: 'File CSV Instruktur DUDI', endpoint: '/import/instruktur-dudi',
    columns: ['full_name (Nama Lengkap)', 'email', 'nik', 'password', 'dudi_id'],
    sample: 'full_name,email,nik,password,dudi_id\nHendra Gunawan,hendra@dudi.id,D-001,rahasia123\n'
  }
}

function roleLabel(r) { return { student: 'Siswa', teacher: 'Guru', dudi: 'Instruktur DUDI', admin: 'Admin', admin_jurusan: 'Admin Jurusan' }[r] || r }
function roleBadge(r) { return { student: 'bg-blue-50 text-blue-600', teacher: 'bg-purple-50 text-purple-600', dudi: 'bg-orange-50 text-orange-600', admin: 'bg-accent/10 text-accent', admin_jurusan: 'bg-green-50 text-green-600' }[r] || '' }

const allSelected = computed(() => users.value.length > 0 && selected.value.length === users.value.length)

function toggleSelect(id) {
  const idx = selected.value.indexOf(id)
  if (idx >= 0) selected.value.splice(idx, 1)
  else selected.value.push(id)
}

function toggleAll() {
  if (allSelected.value) selected.value = []
  else selected.value = users.value.map(u => u.id)
}

async function fetchUsers() {
  loading.value = true
  selected.value = []
  try {
    const params = new URLSearchParams()
    if (roleFilter.value) params.set('role', roleFilter.value)
    if (jurusanFilter.value) params.set('jurusan', jurusanFilter.value)
    if (search.value) params.set('search', search.value)
    const res = await get('/admin/users?' + params.toString())
    users.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchJurusan() {
  try { const res = await get('/admin/jurusan'); jurusanOptions.value = res.data } catch (e) { /* ignore */ }
}

async function fetchDudi() {
  try { const res = await get('/admin/dudi'); dudiList.value = res.data } catch (e) { /* ignore */ }
}

function openCreateModal() {
  editingUser.value = null
  Object.assign(form, { full_name: '', email: '', nis_nip_nik: '', role: 'student', password: '', jurusan: '', dudi_id: '' })
  showModal.value = true
  fetchDudi()
}

function openEditModal(u) {
  editingUser.value = u
  Object.assign(form, {
    full_name: u.full_name, email: u.email, nis_nip_nik: u.nis_nip_nik, role: u.role,
    password: '', jurusan: u.jurusan || '', dudi_id: u.dudi_id || ''
  })
  showModal.value = true
  fetchDudi()
}

function closeModal() { showModal.value = false; editingUser.value = null }

async function saveUser() {
  saving.value = true
  try {
    const payload = {
      full_name: form.full_name, email: form.email, nis_nip_nik: form.nis_nip_nik,
      role: form.role, jurusan: form.jurusan, dudi_id: form.dudi_id || ''
    }
    if (!editingUser.value || form.password) payload.password = form.password

    if (editingUser.value) {
      await put('/admin/users/' + editingUser.value.id, payload)
    } else {
      await post('/admin/users', payload)
    }
    closeModal()
    fetchUsers()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

function confirmDelete(u) { deleteTarget.value = u; showDelete.value = true }

async function deleteUser() {
  saving.value = true
  try {
    await del('/admin/users/' + deleteTarget.value.id)
    showDelete.value = null
    fetchUsers()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

async function bulkDelete() {
  saving.value = true
  try {
    await post('/admin/users/bulk-delete', { ids: selected.value })
    showBulkDelete.value = false
    selected.value = []
    fetchUsers()
  } catch (e) {
    alert(e.message)
  } finally {
    saving.value = false
  }
}

function openImport(key) { importKey.value = key; showImportDropdown.value = false }

onMounted(() => { fetchUsers(); fetchJurusan() })
</script>
