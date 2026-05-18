<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Dashboard Instruktur</h2>
      <p class="text-xs text-gray-500 mt-0.5">{{ authStore.userName }}</p>
    </div>

    <div class="grid grid-cols-3 gap-3 mb-5">
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <UsersIcon :size="20" class="text-primary mx-auto mb-1" />
        <p class="text-2xl font-bold text-gray-800">{{ stats.totalStudents }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Siswa Magang</p>
      </div>
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <ClipboardCheck :size="20" class="text-accent mx-auto mb-1" />
        <p class="text-2xl font-bold text-accent">{{ stats.ratedStudents }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Sudah Dinilai</p>
      </div>
      <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
        <BookOpen :size="20" class="text-warning mx-auto mb-1" />
        <p class="text-2xl font-bold text-warning">{{ stats.totalJournals }}</p>
        <p class="text-[10px] text-gray-500 mt-1">Jurnal Masuk</p>
      </div>
    </div>

    <!-- Students list preview -->
    <div class="mb-3 flex items-center justify-between">
      <h3 class="font-semibold text-gray-800">Siswa Magang</h3>
      <router-link to="/dudi/siswa" class="text-xs text-primary font-medium hover:underline">
        Lihat Semua
      </router-link>
    </div>

    <div class="space-y-3 mb-5">
      <div
        v-for="s in students"
        :key="s.id"
        class="bg-white rounded-xl p-4 border border-gray-100"
      >
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <div class="w-9 h-9 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
              {{ s.name.charAt(0) }}
            </div>
            <div>
              <p class="text-sm font-medium text-gray-800">{{ s.name }}</p>
              <p class="text-[10px] text-gray-400">NIS: {{ s.nis }}</p>
            </div>
          </div>
          <span :class="['text-[10px] font-medium px-2 py-0.5 rounded-full', s.nilai ? 'bg-accent/10 text-accent' : 'bg-gray-100 text-gray-500']">
            {{ s.nilai ? 'Dinilai' : 'Belum' }}
          </span>
        </div>

        <div class="flex items-center gap-2">
          <router-link
            v-if="s.nilai"
            :to="'/dudi/penilaian/' + s.id"
            class="flex-1 text-center py-2 text-xs font-medium text-warning bg-warning/5 rounded-lg hover:bg-warning/10 transition-colors"
          >
            Lihat/Edit Nilai
          </router-link>
          <router-link
            v-else
            :to="'/dudi/penilaian/' + s.id"
            class="flex-1 text-center py-2 text-xs font-medium text-primary bg-primary/5 rounded-lg hover:bg-primary/10 transition-colors"
          >
            Beri Nilai
          </router-link>
          <router-link
            to="/dudi/jurnal"
            class="flex-1 text-center py-2 text-xs font-medium text-gray-500 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            Lihat Jurnal
          </router-link>
        </div>
      </div>
    </div>

    <!-- Unrated alert -->
    <div v-if="unratedCount > 0" class="bg-warning/10 rounded-xl p-4 flex items-start gap-3">
      <AlertCircleIcon :size="20" class="text-warning flex-shrink-0 mt-0.5" />
      <div>
        <p class="text-sm font-medium text-gray-800">Perhatian!</p>
        <p class="text-xs text-gray-600 mt-0.5">
          {{ unratedCount }} siswa belum dinilai. Silakan lengkapi penilaian untuk semua siswa magang.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { UsersIcon, ClipboardCheck, BookOpen, AlertCircleIcon } from 'lucide-vue-next'

const authStore = useAuthStore()

const students = [
  { id: 1, name: 'Ahmad Rizky', nis: '20230001', nilai: true },
  { id: 2, name: 'Siti Nurhaliza', nis: '20230002', nilai: true },
  { id: 3, name: 'Dian Permata', nis: '20230004', nilai: false },
  { id: 4, name: 'Rudi Hartono', nis: '20230005', nilai: true },
  { id: 5, name: 'Maya Sari', nis: '20230007', nilai: false },
]

const stats = {
  totalStudents: 45,
  ratedStudents: 30,
  totalJournals: 128
}

const unratedCount = computed(() => students.filter(s => !s.nilai).length)
</script>
