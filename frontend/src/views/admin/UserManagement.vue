<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Manajemen Pengguna</h2>
        <p class="text-sm text-gray-500 mt-0.5">Kelola data siswa, guru, DUDI, dan admin</p>
      </div>
      <div class="flex items-center gap-2">
        <div class="relative">
          <button
            @click="showImportMenu = !showImportMenu"
            class="flex items-center gap-2 px-4 py-2.5 border border-gray-200 text-gray-700 rounded-xl text-sm font-medium hover:bg-gray-50 transition-colors"
          >
            <UploadIcon :size="18" />
            <span class="hidden sm:inline">Import</span>
            <ChevronDownIcon :size="14" />
          </button>
          <div
            v-if="showImportMenu"
            class="absolute right-0 top-full mt-1 w-56 bg-white rounded-xl shadow-lg border border-gray-100 py-1 z-40"
          >
            <button
              @click="openImport('siswa'); showImportMenu = false"
              class="w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
            >
              <GraduationCap :size="16" class="text-primary" />
              Import Siswa
            </button>
            <button
              @click="openImport('guru'); showImportMenu = false"
              class="w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
            >
              <UserCircle :size="16" class="text-accent" />
              Import Guru
            </button>
            <button
              @click="openImport('instruktur'); showImportMenu = false"
              class="w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
            >
              <Building2 :size="16" class="text-warning" />
              Import Instruktur DUDI
            </button>
          </div>
        </div>
        <button
          @click="showModal = true"
          class="flex items-center gap-2 px-4 py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
        >
          <PlusIcon :size="18" />
          <span class="hidden sm:inline">Tambah Pengguna</span>
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-2xl p-4 border border-gray-100 mb-4 flex flex-wrap gap-3">
      <div class="flex-1 min-w-[200px]">
        <input
          v-model="search"
          type="text"
          placeholder="Cari nama, email, NIS/NIP/NIK..."
          class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
        />
      </div>
      <select
        v-model="filterRole"
        class="px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none bg-white"
      >
        <option value="">Semua Role</option>
        <option value="student">Siswa</option>
        <option value="teacher">Guru</option>
        <option value="dudi">DUDI</option>
        <option value="admin">Admin</option>
      </select>
      <button class="px-4 py-2 rounded-xl border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 flex items-center gap-2">
        <FilterIcon :size="16" />
        Filter
      </button>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Nama</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Role</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Email</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">NIS/NIP/NIK</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Status</th>
              <th class="text-right px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in filteredUsers" :key="u.id" class="border-b border-gray-50 hover:bg-gray-50/50 transition-colors">
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-bold">
                    {{ u.name.charAt(0) }}
                  </div>
                  <span class="text-sm font-medium text-gray-800">{{ u.name }}</span>
                </div>
              </td>
              <td class="px-4 py-3">
                <span :class="['inline-flex px-2.5 py-0.5 rounded-full text-xs font-medium', roleBadge(u.role)]">
                  {{ roleLabel(u.role) }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ u.email }}</td>
              <td class="px-4 py-3 text-sm text-gray-500 font-mono">{{ u.nis }}</td>
              <td class="px-4 py-3">
                <span :class="['inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium', u.active ? 'bg-accent/10 text-accent' : 'bg-gray-100 text-gray-500']">
                  <span :class="['w-1.5 h-1.5 rounded-full', u.active ? 'bg-accent' : 'bg-gray-400']" />
                  {{ u.active ? 'Aktif' : 'Nonaktif' }}
                </span>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button class="p-1.5 rounded-lg hover:bg-gray-100 text-gray-500" title="Edit">
                    <PencilIcon :size="16" />
                  </button>
                  <button class="p-1.5 rounded-lg hover:bg-red-50 text-gray-500 hover:text-red-500" title="Hapus">
                    <TrashIcon :size="16" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-4 py-3 border-t border-gray-100">
        <p class="text-sm text-gray-500">Menampilkan {{ filteredUsers.length }} dari {{ users.length }} pengguna</p>
        <div class="flex items-center gap-1">
          <button class="px-3 py-1.5 rounded-lg text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-40" disabled>
            <ChevronLeftIcon :size="16" />
          </button>
          <button class="px-3 py-1.5 rounded-lg text-sm bg-primary text-white">1</button>
          <button class="px-3 py-1.5 rounded-lg text-sm text-gray-600 hover:bg-gray-100">2</button>
          <button class="px-3 py-1.5 rounded-lg text-sm text-gray-600 hover:bg-gray-100">3</button>
          <button class="px-3 py-1.5 rounded-lg text-sm text-gray-500 hover:bg-gray-100">
            <ChevronRightIcon :size="16" />
          </button>
        </div>
      </div>
    </div>

    <CsvImport
      :show="showImport"
      :title="importTitle"
      :label="importLabel"
      :endpoint="importEndpoint"
      :csv-columns="importColumns"
      :sample-csv="importSample"
      @close="showImport = false"
      @done="showImport = false"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { PlusIcon, FilterIcon, PencilIcon, TrashIcon, ChevronLeftIcon, ChevronRightIcon, UploadIcon, ChevronDownIcon, GraduationCap, UserCircle, Building2 } from 'lucide-vue-next'
import CsvImport from '../../components/CsvImport.vue'

const search = ref('')
const filterRole = ref('')
const showModal = ref(false)
const showImport = ref(false)
const showImportMenu = ref(false)
const importType = ref('')

const users = [
  { id: 1, name: 'Ahmad Rizky', email: 'ahmad@pkl.local', role: 'student', nis: '20230001', active: true },
  { id: 2, name: 'Siti Nurhaliza', email: 'siti@pkl.local', role: 'student', nis: '20230002', active: true },
  { id: 3, name: 'Budi Santoso, S.Kom', email: 'budi@pkl.local', role: 'teacher', nis: '19850101', active: true },
  { id: 4, name: 'PT. Teknologi Maju', email: 'info@teknologimaju.id', role: 'dudi', nis: 'D-001', active: true },
  { id: 5, name: 'Admin Utama', email: 'admin@pkl.local', role: 'admin', nis: 'ADM-001', active: true },
  { id: 6, name: 'Rina Marlina', email: 'rina@pkl.local', role: 'student', nis: '20230003', active: false },
  { id: 7, name: 'Dr. Ahmad Fauzi', email: 'fauzi@pkl.local', role: 'teacher', nis: '19870203', active: true },
  { id: 8, name: 'PT. Sejahtera Abadi', email: 'info@sejahtera.id', role: 'dudi', nis: 'D-002', active: true },
]

const filteredUsers = computed(() => {
  return users.filter(u => {
    const matchSearch = u.name.toLowerCase().includes(search.value.toLowerCase()) ||
      u.email.toLowerCase().includes(search.value.toLowerCase()) ||
      u.nis.includes(search.value)
    const matchRole = !filterRole.value || u.role === filterRole.value
    return matchSearch && matchRole
  })
})

function roleLabel(role) {
  const map = { student: 'Siswa', teacher: 'Guru', dudi: 'DUDI', admin: 'Admin' }
  return map[role]
}

function roleBadge(role) {
  const map = {
    student: 'bg-primary/10 text-primary',
    teacher: 'bg-accent/10 text-accent',
    dudi: 'bg-warning/10 text-warning',
    admin: 'bg-gray-100 text-gray-600'
  }
  return map[role]
}

const importConfigs = {
  siswa: {
    title: 'Import Siswa',
    label: 'Siswa',
    endpoint: '/import/siswa',
    columns: 'full_name, email, nis, password, dudi_nik',
    sample: 'full_name,email,nis,password,dudi_nik\nAhmad Rizky,ahmad@pkl.local,20230001,siswa123,D-001\nSiti Nurhaliza,siti@pkl.local,20230002,siswa123,D-002'
  },
  guru: {
    title: 'Import Guru Pembimbing',
    label: 'Guru',
    endpoint: '/import/guru',
    columns: 'full_name, email, nip, password',
    sample: 'full_name,email,nip,password\nBudi Santoso,budi@pkl.local,19850101,guru123\nDewi Lestari,dewi@pkl.local,19880102,guru123'
  },
  instruktur: {
    title: 'Import Instruktur DUDI',
    label: 'Instruktur DUDI',
    endpoint: '/import/instruktur-dudi',
    columns: 'full_name, email, nik, password, dudi_nik',
    sample: 'full_name,email,nik,password,dudi_nik\nPak Hendra,hendra@pkl.local,D-001,dudi123,D-001\nIbu Ratna,ratna@pkl.local,D-002,dudi123,D-002'
  }
}

const importTitle = computed(() => importConfigs[importType.value]?.title || 'Import')
const importLabel = computed(() => importConfigs[importType.value]?.label || '')
const importEndpoint = computed(() => importConfigs[importType.value]?.endpoint || '')
const importColumns = computed(() => importConfigs[importType.value]?.columns || '')
const importSample = computed(() => importConfigs[importType.value]?.sample || '')

function openImport(type) {
  importType.value = type
  showImport.value = true
}
</script>
