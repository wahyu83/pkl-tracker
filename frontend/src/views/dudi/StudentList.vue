<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Daftar Siswa Magang</h2>
      <p class="text-xs text-gray-500 mt-0.5">{{ authStore.userName }}</p>
    </div>

    <div class="bg-white rounded-xl p-3 border border-gray-100 mb-4">
      <div class="relative">
        <SearchIcon :size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          v-model="search"
          type="text"
          placeholder="Cari nama siswa..."
          class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
        />
      </div>
    </div>

    <div class="space-y-3">
      <div
        v-for="s in filteredStudents"
        :key="s.id"
        class="bg-white rounded-xl p-4 border border-gray-100"
      >
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-sm font-bold text-primary">
            {{ s.name.charAt(0) }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-800">{{ s.name }}</p>
            <p class="text-xs text-gray-400">NIS: {{ s.nis }} | {{ s.jurusan }}</p>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-2 mb-3">
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-sm font-bold" :class="s.attendance >= 80 ? 'text-accent' : s.attendance >= 60 ? 'text-warning' : 'text-danger'">{{ s.attendance }}%</p>
            <p class="text-[9px] text-gray-400">Kehadiran</p>
          </div>
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-sm font-bold text-info">{{ s.journals }}</p>
            <p class="text-[9px] text-gray-400">Jurnal</p>
          </div>
          <div class="text-center p-2 bg-gray-50 rounded-lg">
            <p class="text-sm font-bold" :class="s.nilai ? 'text-warning' : 'text-gray-300'">{{ s.nilai || '-' }}</p>
            <p class="text-[9px] text-gray-400">Nilai</p>
          </div>
        </div>

        <router-link
          :to="'/dudi/penilaian/' + s.id"
          :class="[
            'block w-full text-center py-2.5 rounded-xl text-sm font-medium transition-colors',
            s.nilai
              ? 'text-warning bg-warning/5 hover:bg-warning/10'
              : 'text-primary bg-primary/5 hover:bg-primary/10'
          ]"
        >
          {{ s.nilai ? 'Edit Nilai' : 'Beri Nilai' }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { SearchIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const search = ref('')

const students = [
  { id: 1, name: 'Ahmad Rizky', nis: '20230001', jurusan: 'RPL', attendance: 95, journals: 42, nilai: 'A' },
  { id: 2, name: 'Siti Nurhaliza', nis: '20230002', jurusan: 'TKJ', attendance: 88, journals: 40, nilai: 'B+' },
  { id: 3, name: 'Dian Permata', nis: '20230004', jurusan: 'RPL', attendance: 72, journals: 35, nilai: null },
  { id: 4, name: 'Rudi Hartono', nis: '20230005', jurusan: 'MM', attendance: 100, journals: 45, nilai: 'A' },
  { id: 5, name: 'Maya Sari', nis: '20230007', jurusan: 'TKJ', attendance: 60, journals: 28, nilai: null },
  { id: 6, name: 'Bambang Kusumo', nis: '20230008', jurusan: 'RPL', attendance: 85, journals: 38, nilai: 'B+' },
]

const filteredStudents = computed(() => {
  return students.filter(s =>
    s.name.toLowerCase().includes(search.value.toLowerCase()) ||
    s.nis.includes(search.value)
  )
})
</script>
