<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-xl font-bold text-gray-800">Data DUDI</h2>
        <p class="text-sm text-gray-500 mt-0.5">Kelola data perusahaan/industri mitra PKL</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="showImport = true"
          class="flex items-center gap-2 px-4 py-2.5 border border-gray-200 text-gray-700 rounded-xl text-sm font-medium hover:bg-gray-50 transition-colors"
        >
          <UploadIcon :size="18" />
          <span class="hidden sm:inline">Import CSV</span>
        </button>
        <button
          @click="showModal = true"
          class="flex items-center gap-2 px-4 py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
        >
          <PlusIcon :size="18" />
          <span class="hidden sm:inline">Tambah DUDI</span>
        </button>
      </div>
    </div>

    <div class="bg-white rounded-2xl p-4 border border-gray-100 mb-4">
      <div class="flex flex-wrap gap-3">
        <div class="flex-1 min-w-[200px]">
          <input
            v-model="search"
            type="text"
            placeholder="Cari nama perusahaan..."
            class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
          />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="d in filteredDudi"
        :key="d.id"
        class="bg-white rounded-2xl p-5 border border-gray-100 card-hover"
      >
        <div class="flex items-start gap-3 mb-3">
          <div class="w-10 h-10 rounded-xl bg-warning/10 text-warning flex items-center justify-center flex-shrink-0">
            <Building2 :size="20" />
          </div>
          <div class="min-w-0">
            <h3 class="font-semibold text-gray-800 truncate">{{ d.name }}</h3>
            <p class="text-xs text-gray-500 mt-0.5">PIC: {{ d.pic }}</p>
          </div>
        </div>

        <div class="space-y-2 mb-4">
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <MapPinIcon :size="14" />
            <span class="truncate">{{ d.address }}</span>
          </div>
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <PhoneIcon :size="14" />
            <span>{{ d.phone }}</span>
          </div>
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <UsersIcon :size="14" />
            <span>{{ d.studentCount }} siswa magang</span>
          </div>
          <div class="flex items-center gap-2 text-xs text-gray-500">
            <TargetIcon :size="14" />
            <span>Radius: {{ d.radius }}m</span>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-3 border-t border-gray-50">
          <button class="flex-1 py-2 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors">
            Lihat Detail
          </button>
          <button class="px-3 py-2 text-xs font-medium text-gray-500 hover:text-red-500 rounded-lg hover:bg-red-50 transition-colors">
            <TrashIcon :size="15" />
          </button>
        </div>
      </div>
    </div>

    <CsvImport
      :show="showImport"
      title="Import Data DUDI"
      label="DUDI"
      endpoint="/import/dudi"
      csv-columns="company_name, address, latitude, longitude, radius_allowed, pic_name, phone"
      :sample-csv="dudiSample"
      @close="showImport = false"
      @done="showImport = false"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { PlusIcon, Building2, MapPinIcon, PhoneIcon, UsersIcon, TargetIcon, TrashIcon, UploadIcon } from 'lucide-vue-next'
import CsvImport from '../../components/CsvImport.vue'

const dudiSample = `company_name,address,latitude,longitude,radius_allowed,pic_name,phone
PT. Teknologi Maju,Jl. Sudirman 123 Jakarta,-6.2088,106.8456,500,Hendra Gunawan,021-5551234
PT. Sejahtera Abadi,Jl. Gatot Subroto 45 Jakarta,-6.2297,106.8243,300,Ratna Dewi,021-5555678`

const search = ref('')
const showModal = ref(false)
const showImport = ref(false)

const dudiList = [
  { id: 1, name: 'PT. Teknologi Maju', pic: 'Hendra Gunawan', address: 'Jl. Sudirman No. 123, Jakarta Pusat', phone: '021-5551234', studentCount: 45, radius: 500, lat: -6.2088, lng: 106.8456 },
  { id: 2, name: 'PT. Sejahtera Abadi', pic: 'Ratna Dewi', address: 'Jl. Gatot Subroto No. 45, Jakarta Selatan', phone: '021-5555678', studentCount: 32, radius: 300, lat: -6.2297, lng: 106.8243 },
  { id: 3, name: 'CV. Kreatif Digital', pic: 'Andi Pratama', address: 'Jl. Merdeka No. 78, Bandung', phone: '022-5559012', studentCount: 18, radius: 400, lat: -6.9175, lng: 107.6191 },
  { id: 4, name: 'UD. Mandiri Jaya', pic: 'Susi Susanti', address: 'Jl. Ahmad Yani No. 10, Surabaya', phone: '031-5553456', studentCount: 28, radius: 250, lat: -7.2575, lng: 112.7521 },
  { id: 5, name: 'PT. Inovasi Nusantara', pic: 'Bambang Hermanto', address: 'Jl. Thamrin No. 56, Jakarta Pusat', phone: '021-5557890', studentCount: 52, radius: 600, lat: -6.1823, lng: 106.8229 },
]

const filteredDudi = computed(() => {
  return dudiList.filter(d =>
    d.name.toLowerCase().includes(search.value.toLowerCase())
  )
})
</script>
