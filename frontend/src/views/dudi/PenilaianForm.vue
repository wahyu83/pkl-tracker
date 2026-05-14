<template>
  <div>
    <div class="flex items-center gap-3 mb-5">
      <router-link to="/dudi/siswa" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ArrowLeftIcon :size="20" />
      </router-link>
      <div>
        <h2 class="text-lg font-bold text-gray-800">Penilaian PKL</h2>
        <p class="text-xs text-gray-500">{{ studentName }}</p>
      </div>
    </div>

    <!-- Auto attendance score -->
    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <h3 class="font-semibold text-gray-800 mb-3">Kehadiran (Otomatis)</h3>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500">Total Kehadiran</p>
          <p class="text-2xl font-bold text-accent">{{ attendanceScore }}%</p>
        </div>
        <div class="w-16 h-16 rounded-full bg-accent/10 flex items-center justify-center">
          <ClipboardCheck :size="28" class="text-accent" />
        </div>
      </div>
    </div>

    <!-- Manual scoring -->
    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <h3 class="font-semibold text-gray-800 mb-4">Penilaian Manual (Skala 1-5)</h3>

      <div class="space-y-4">
        <div v-for="item in criteria" :key="item.key">
          <div class="flex items-center justify-between mb-1.5">
            <div class="flex items-center gap-2">
              <component :is="item.icon" :size="16" class="text-gray-400" />
              <span class="text-sm text-gray-700">{{ item.label }}</span>
            </div>
            <span class="text-sm font-bold" :class="scores[item.key] >= 4 ? 'text-accent' : scores[item.key] >= 3 ? 'text-warning' : 'text-danger'">
              {{ scores[item.key] }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-[10px] text-gray-400 w-6">1</span>
            <input
              type="range"
              min="1"
              max="5"
              v-model.number="scores[item.key]"
              class="flex-1 h-2 rounded-full appearance-none bg-gray-200 accent-primary [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-5 [&::-webkit-slider-thumb]:h-5 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-primary [&::-webkit-slider-thumb]:cursor-pointer"
            />
            <span class="text-[10px] text-gray-400 w-6 text-right">5</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Notes -->
    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">Catatan</label>
      <textarea
        v-model="notes"
        rows="4"
        placeholder="Tulis catatan tambahan tentang performa siswa..."
        class="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none"
        maxlength="500"
      />
      <p class="text-xs text-gray-400 text-right">{{ notes.length }}/500</p>
    </div>

    <!-- Preview Final Score -->
    <div class="bg-primary/5 rounded-2xl p-5 mb-4">
      <div class="flex items-center justify-between mb-2">
        <span class="text-sm text-gray-600">Nilai Akhir (Kalkulasi)</span>
        <span class="text-lg font-bold text-primary">{{ finalScore }}</span>
      </div>
      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-600">Grade</span>
        <span class="text-xl font-black" :class="gradeColor">{{ finalGrade }}</span>
      </div>
    </div>

    <button
      @click="handleSubmit"
      :disabled="submitting"
      class="w-full py-3 bg-primary text-white rounded-2xl text-sm font-bold hover:bg-primary-light transition-colors disabled:opacity-60 flex items-center justify-center gap-2"
    >
      <LoaderIcon v-if="submitting" :size="20" class="animate-spin" />
      <span>{{ submitting ? 'Menyimpan...' : 'Submit Penilaian' }}</span>
    </button>

    <!-- Success -->
    <div
      v-if="success"
      class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-16 h-16 bg-accent/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircleIcon :size="36" class="text-accent" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-2">Nilai Tersimpan!</h3>
        <p class="text-sm text-gray-500 mb-4">Penilaian untuk {{ studentName }} telah tersimpan dan terkunci.</p>
        <router-link
          to="/dudi/siswa"
          class="block w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
        >
          Kembali ke Daftar Siswa
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeftIcon, ClipboardCheck, ShieldCheck, HandshakeIcon, LightbulbIcon, LoaderIcon, CheckCircleIcon } from 'lucide-vue-next'

const route = useRoute()
const submitting = ref(false)
const success = ref(false)
const notes = ref('')
const attendanceScore = ref(95)

const studentId = computed(() => route.params.studentId || '1')

const studentNames = {
  '1': 'Ahmad Rizky',
  '2': 'Siti Nurhaliza',
  '3': 'Dian Permata',
  '4': 'Rudi Hartono',
  '5': 'Maya Sari',
}

const studentName = computed(() => studentNames[studentId.value] || 'Siswa')

const scores = reactive({
  discipline: 4,
  responsibility: 5,
  teamwork: 4,
  initiative: 5
})

const criteria = [
  { key: 'discipline', label: 'Kedisiplinan', icon: ShieldCheck },
  { key: 'responsibility', label: 'Tanggung Jawab', icon: ShieldCheck },
  { key: 'teamwork', label: 'Kerjasama', icon: HandshakeIcon },
  { key: 'initiative', label: 'Inisiatif', icon: LightbulbIcon },
]

const manualAvg = computed(() => {
  const vals = Object.values(scores)
  return vals.reduce((a, b) => a + b, 0) / vals.length
})

const finalScore = computed(() => {
  const attendanceWeight = 0.3
  const manualWeight = 0.7
  return Math.round((attendanceScore.value * attendanceWeight) + (manualAvg.value / 5 * 100 * manualWeight))
})

const finalGrade = computed(() => {
  const s = finalScore.value
  if (s >= 90) return 'A'
  if (s >= 80) return 'B+'
  if (s >= 70) return 'B'
  if (s >= 60) return 'C'
  return 'D'
})

const gradeColor = computed(() => {
  const g = finalGrade.value.charAt(0)
  const map = { A: 'text-accent', B: 'text-info', C: 'text-warning', D: 'text-danger' }
  return map[g] || 'text-gray-400'
})

function handleSubmit() {
  submitting.value = true
  setTimeout(() => {
    submitting.value = false
    success.value = true
  }, 1200)
}
</script>
