<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Dashboard Instruktur</h2>
      <p class="text-xs text-gray-500 mt-0.5">{{ authStore.userName }}</p>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat...</div>

    <template v-else>
      <div class="grid grid-cols-3 gap-3 mb-5">
        <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
          <UsersIcon :size="20" class="text-primary mx-auto mb-1" />
          <p class="text-2xl font-bold text-gray-800">{{ stats.total_students }}</p>
          <p class="text-[10px] text-gray-500 mt-1">Siswa PKL</p>
        </div>
        <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
          <ClipboardCheck :size="20" class="text-accent mx-auto mb-1" />
          <p class="text-2xl font-bold text-accent">{{ stats.rated_students }}</p>
          <p class="text-[10px] text-gray-500 mt-1">Sudah Dinilai</p>
        </div>
        <div class="bg-white rounded-xl p-4 border border-gray-100 text-center">
          <BookOpen :size="20" class="text-warning mx-auto mb-1" />
          <p class="text-2xl font-bold text-warning">{{ stats.total_journals }}</p>
          <p class="text-[10px] text-gray-500 mt-1">Jurnal Masuk</p>
        </div>
      </div>

      <!-- Students list preview -->
      <div class="mb-3 flex items-center justify-between">
        <h3 class="font-semibold text-gray-800">Siswa PKL</h3>
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

        <div v-if="students.length === 0" class="text-center py-6 text-gray-400 text-sm">
          Belum ada siswa magang yang terdaftar di DUDI Anda
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
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { get } from '../../api'
import { UsersIcon, ClipboardCheck, BookOpen, AlertCircleIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const loading = ref(true)

const stats = ref({ total_students: 0, rated_students: 0, total_journals: 0 })
const students = ref([])

const unratedCount = computed(() => students.value.filter(s => !s.nilai).length)

async function fetchDashboard() {
  loading.value = true
  try {
    const res = await get('/dudi/dashboard')
    stats.value = res.stats
    students.value = res.students || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchDashboard)
</script>
