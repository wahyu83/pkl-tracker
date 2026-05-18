<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Nilai PKL</h2>
      <p class="text-xs text-gray-500 mt-0.5">Penilaian dari Instruktur DUDI tempat magang</p>
    </div>

    <template v-if="nilai">
      <!-- Final grade -->
      <div class="bg-white rounded-2xl p-6 border border-gray-100 mb-5 text-center">
        <p class="text-sm text-gray-500 mb-2">Nilai Akhir PKL</p>
        <div :class="['w-20 h-20 rounded-full mx-auto flex items-center justify-center mb-2', gradeCircle(nilai.grade)]">
          <span :class="['text-3xl font-black', gradeTextColor(nilai.grade)]">{{ nilai.grade }}</span>
        </div>
        <div class="flex items-center justify-center gap-2">
          <span class="text-2xl font-bold text-gray-800">{{ nilai.finalScore }}</span>
          <span class="text-xs text-gray-500">/ 100</span>
        </div>
        <p class="text-xs text-gray-400 mt-1">dinilai oleh {{ nilai.dudiName }}</p>
      </div>

      <!-- Score details -->
      <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-5">
        <h3 class="font-semibold text-gray-800 mb-4">Rincian Penilaian</h3>

        <div class="space-y-4">
          <div v-for="item in scoreItems" :key="item.label">
            <div class="flex justify-between items-center mb-1">
              <div class="flex items-center gap-2">
                <component :is="item.icon" :size="16" class="text-gray-400" />
                <span class="text-sm text-gray-600">{{ item.label }}</span>
              </div>
              <span class="text-sm font-bold text-gray-800">{{ item.value }}{{ item.suffix }}</span>
            </div>
            <div class="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                :class="item.barColor"
                class="h-full rounded-full transition-all duration-500"
                :style="{ width: item.percent + '%' }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="nilai.notes" class="bg-white rounded-2xl p-5 border border-gray-100">
        <h3 class="font-semibold text-gray-800 mb-2">Catatan dari Instruktur</h3>
        <p class="text-sm text-gray-600 italic">"{{ nilai.notes }}"</p>
      </div>
    </template>

    <template v-else>
      <div class="bg-white rounded-2xl p-6 border border-gray-100 text-center">
        <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <Award :size="32" class="text-gray-300" />
        </div>
        <h3 class="font-semibold text-gray-800 mb-2">Belum Ada Nilai</h3>
        <p class="text-sm text-gray-500">Instruktur DUDI belum mengisi penilaian PKL Anda.</p>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ClipboardCheck, ShieldCheck, HandshakeIcon, LightbulbIcon } from 'lucide-vue-next'

const nilai = ref({
  grade: 'A',
  finalScore: 94,
  dudiName: 'PT. Teknologi Maju',
  notes: 'Siswa menunjukkan performa yang sangat baik. Disiplin, proaktif, dan cepat belajar. Direkomendasikan untuk mendapat nilai A.',
  attendance: 95,
  discipline: 5,
  responsibility: 5,
  teamwork: 4,
  initiative: 5
})

const scoreItems = [
  { label: 'Kehadiran', value: nilai.value.attendance, suffix: '%', percent: nilai.value.attendance, barColor: 'bg-primary', icon: ClipboardCheck },
  { label: 'Kedisiplinan', value: nilai.value.discipline, suffix: '/5', percent: (nilai.value.discipline / 5) * 100, barColor: 'bg-accent', icon: ShieldCheck },
  { label: 'Tanggung Jawab', value: nilai.value.responsibility, suffix: '/5', percent: (nilai.value.responsibility / 5) * 100, barColor: 'bg-warning', icon: ShieldCheck },
  { label: 'Kerjasama', value: nilai.value.teamwork, suffix: '/5', percent: (nilai.value.teamwork / 5) * 100, barColor: 'bg-info', icon: HandshakeIcon },
  { label: 'Inisiatif', value: nilai.value.initiative, suffix: '/5', percent: (nilai.value.initiative / 5) * 100, barColor: 'bg-purple-500', icon: LightbulbIcon },
]

function gradeCircle(grade) {
  const g = grade.charAt(0)
  const map = { A: 'bg-accent/10', B: 'bg-info/10', C: 'bg-warning/10', D: 'bg-danger/10' }
  return map[g] || 'bg-gray-100'
}

function gradeTextColor(grade) {
  const g = grade.charAt(0)
  const map = { A: 'text-accent', B: 'text-info', C: 'text-warning', D: 'text-danger' }
  return map[g] || 'text-gray-400'
}
</script>
