<template>
  <div>
    <div class="flex items-center gap-3 mb-5">
      <router-link :to="backPath" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ArrowLeftIcon :size="20" />
      </router-link>
      <div>
        <h2 class="text-lg font-bold text-gray-800">Penilaian PKL</h2>
        <p class="text-xs text-gray-500">{{ studentName }}</p>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else-if="!studentId" class="bg-white rounded-2xl p-6 border border-gray-100 text-center">
      <p class="text-sm text-gray-500">Pilih siswa terlebih dahulu untuk memberi penilaian.</p>
      <router-link :to="backPath" class="inline-block mt-3 px-4 py-2 rounded-xl text-sm font-medium text-primary bg-primary/5 hover:bg-primary/10">
        Buka Daftar Siswa
      </router-link>
    </div>

    <template v-else>
      <!-- Rentang pelaksanaan PKL (khusus guru) -->
      <div v-if="isTeacher" class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
        <h3 class="font-semibold text-gray-800 mb-1">Rentang Pelaksanaan PKL</h3>
        <p class="text-xs text-gray-500 mb-3">
          Menjadi dasar perhitungan persentase kehadiran pada sertifikat. Kosongkan untuk memakai periode aktif.
        </p>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Tanggal Mulai</label>
            <input v-model="startDate" type="date" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none" />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Tanggal Selesai</label>
            <input v-model="endDate" type="date" class="w-full px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none" />
          </div>
        </div>
      </div>

      <!-- Auto attendance score -->
      <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
        <h3 class="font-semibold text-gray-800 mb-3">Kehadiran (Otomatis)</h3>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Total Kehadiran</p>
            <p class="text-2xl font-bold text-accent">{{ Math.round(attendanceScore) }}%</p>
            <p v-if="isTeacher" class="text-[11px] text-gray-400 mt-1">
              {{ rangeLabel }}
            </p>
          </div>
          <div class="w-16 h-16 rounded-full bg-accent/10 flex items-center justify-center">
            <ClipboardCheck :size="28" class="text-accent" />
          </div>
        </div>
      </div>

      <!-- Manual scoring -->
      <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
        <h3 class="font-semibold text-gray-800 mb-4">Penilaian Manual (Skala 0-100)</h3>

        <div class="space-y-4">
          <div v-for="item in criteria" :key="item.key">
            <div class="flex items-center justify-between mb-1.5">
              <div class="flex items-center gap-2">
                <component :is="item.icon" :size="16" class="text-gray-400" />
                <span class="text-sm text-gray-700">{{ item.label }}</span>
              </div>
              <span class="text-sm font-bold" :class="scores[item.key] >= 80 ? 'text-accent' : scores[item.key] >= 60 ? 'text-warning' : 'text-danger'">
                {{ scores[item.key] }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-[10px] text-gray-400 w-6">0</span>
              <input
                type="range"
                min="0"
                max="100"
                v-model.number="scores[item.key]"
                class="flex-1 h-2 rounded-full appearance-none bg-gray-200 accent-primary [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-5 [&::-webkit-slider-thumb]:h-5 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-primary [&::-webkit-slider-thumb]:cursor-pointer"
              />
              <span class="text-[10px] text-gray-400 w-8 text-right">100</span>
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
          <p class="text-sm text-gray-500 mb-4">Penilaian untuk {{ studentName }} telah tersimpan.</p>
          <router-link
            :to="backPath"
            class="block w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
          >
            Kembali
          </router-link>
        </div>
      </div>
    </template>

    <!-- Menu Sertifikat -->
    <router-link
      :to="certPath"
      class="mt-4 w-full flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-semibold text-primary bg-primary/5 hover:bg-primary/10 transition-colors"
    >
      <AwardIcon :size="18" /> Cetak Sertifikat
    </router-link>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { get, post } from '@/api'
import {
  ArrowLeftIcon, ClipboardCheck, ShieldCheck, HandshakeIcon,
  LightbulbIcon, LoaderIcon, CheckCircleIcon, Award as AwardIcon
} from 'lucide-vue-next'

const route = useRoute()
const authStore = useAuthStore()

const isTeacher = computed(() => authStore.userRole === 'teacher')
const backPath = computed(() => (isTeacher.value ? '/guru/nilai' : '/dudi/siswa'))
const certPath = computed(() => (isTeacher.value ? '/guru/sertifikat' : '/dudi/sertifikat'))

const submitting = ref(false)
const success = ref(false)
const loading = ref(true)
const notes = ref('')
const attendanceScore = ref(0)
const studentName = ref('')
const startDate = ref('')
const endDate = ref('')

const studentId = computed(() => route.params.studentId)

const scores = reactive({
  discipline: 60,
  responsibility: 60,
  teamwork: 60,
  initiative: 60
})

const criteria = [
  { key: 'discipline', label: 'Memahami alur bisnis tempat PKL', icon: ShieldCheck },
  { key: 'responsibility', label: 'Menerapkan soft skills', icon: ShieldCheck },
  { key: 'teamwork', label: 'Menerapkan kompetensi teknis', icon: HandshakeIcon },
  { key: 'initiative', label: 'Menerapkan POS dan K3LH', icon: LightbulbIcon },
]

const rangeLabel = computed(() => {
  if (!startDate.value && !endDate.value) return 'Dihitung dari seluruh absensi (rentang belum diisi)'
  const fmt = d => d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '...'
  return `Dihitung dari absensi ${fmt(startDate.value)} – ${fmt(endDate.value)}`
})

const manualAvg = computed(() => {
  const vals = Object.values(scores)
  return vals.reduce((a, b) => a + b, 0) / vals.length
})

const finalScore = computed(() => {
  const attendanceWeight = 0.3
  const manualWeight = 0.7
  return Math.round((attendanceScore.value * attendanceWeight) + (manualAvg.value * manualWeight))
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

function toDateInput(value) {
  if (!value) return ''
  return String(value).slice(0, 10)
}

async function fetchStudentName() {
  try {
    if (isTeacher.value) {
      const res = await get('/guru/students')
      const found = (res.data || []).find(s => s.id === studentId.value)
      studentName.value = found?.full_name || 'Siswa'
    } else {
      const res = await get('/dudi/dashboard')
      const found = (res.students || []).find(s => s.id === studentId.value)
      studentName.value = found?.name || 'Siswa'
    }
  } catch (e) {
    studentName.value = 'Siswa'
  }
}

async function fetchData() {
  if (!studentId.value) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const [nilaiRes] = await Promise.all([
      get('/nilai/' + studentId.value).catch(() => ({ data: null })),
      fetchStudentName()
    ])

    if (nilaiRes.data) {
      const n = nilaiRes.data
      attendanceScore.value = n.AttendanceScoreAuto || n.attendance_score_auto || 0
      const d = n.Discipline || n.discipline || 60
      const r = n.Responsibility || n.responsibility || 60
      const t = n.Teamwork || n.teamwork || 60
      const i = n.Initiative || n.initiative || 60
      // konversi skala lama 1-5 ke 0-100 bila perlu
      scores.discipline = d <= 5 ? d * 20 : d
      scores.responsibility = r <= 5 ? r * 20 : r
      scores.teamwork = t <= 5 ? t * 20 : t
      scores.initiative = i <= 5 ? i * 20 : i
      notes.value = n.Notes || n.notes || ''
      startDate.value = toDateInput(n.attendance_start_date)
      endDate.value = toDateInput(n.attendance_end_date)
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!studentId.value) return
  submitting.value = true
  try {
    const payload = {
      student_id: studentId.value,
      discipline: scores.discipline,
      responsibility: scores.responsibility,
      teamwork: scores.teamwork,
      initiative: scores.initiative,
      notes: notes.value
    }

    if (isTeacher.value) {
      payload.attendance_start_date = startDate.value
      payload.attendance_end_date = endDate.value
    }

    const res = await post('/nilai', payload)
    // Tampilkan skor kehadiran terbaru hasil hitung server.
    const saved = res.data
    if (saved) {
      attendanceScore.value = saved.attendance_score_auto || 0
      startDate.value = toDateInput(saved.attendance_start_date) || startDate.value
      endDate.value = toDateInput(saved.attendance_end_date) || endDate.value
    }
    success.value = true
  } catch (e) {
    alert(e.message)
  } finally {
    submitting.value = false
  }
}

onMounted(fetchData)
</script>
