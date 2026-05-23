<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Daftar Siswa PKL</h2>
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

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else-if="filteredStudents.length === 0" class="text-center py-8 text-gray-400 text-sm">
      Belum ada siswa PKL yang terdaftar di DUDI Anda
    </div>

    <div v-else class="space-y-3">
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

        <div class="flex items-center gap-2 mb-2">
          <router-link
            :to="'/dudi/penilaian/' + s.id"
            :class="[
              'flex-1 block text-center py-2.5 rounded-xl text-sm font-medium transition-colors',
              s.nilai
                ? 'text-warning bg-warning/5 hover:bg-warning/10'
                : 'text-primary bg-primary/5 hover:bg-primary/10'
            ]"
          >
            {{ s.nilai ? 'Edit Nilai' : 'Beri Nilai' }}
          </router-link>
          <button
            @click="openIzin(s)"
            class="flex-1 text-center py-2.5 rounded-xl text-sm font-medium text-info bg-info/5 hover:bg-info/10 transition-colors"
          >
            Izin/Sakit
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Izin/Sakit -->
    <div v-if="izinModal.open" class="fixed inset-0 z-50 flex items-end justify-center bg-black/30" @click.self="izinModal.open = false">
      <div class="bg-white rounded-t-2xl w-full max-w-lg p-5 pb-8">
        <h3 class="text-sm font-bold text-gray-800 mb-3">Catat Izin / Sakit</h3>
        <p class="text-xs text-gray-500 mb-4">{{ izinModal.studentName }} — NIS: {{ izinModal.studentNis }}</p>

        <div class="space-y-3">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Status</label>
              <select v-model="izinModal.status" class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none bg-white">
                <option value="">-- Pilih --</option>
                <option value="izin">Izin</option>
                <option value="sakit">Sakit</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-600 mb-1">Tanggal</label>
              <input v-model="izinModal.tanggal" type="date" class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary outline-none" />
            </div>
          </div>
          <div class="flex gap-2">
            <button @click="izinModal.open = false" class="flex-1 py-2.5 rounded-xl text-sm font-medium text-gray-600 bg-gray-100 hover:bg-gray-200 transition-colors">
              Batal
            </button>
            <button
              @click="submitIzin"
              :disabled="submittingIzin"
              class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-white bg-primary hover:bg-primary/90 disabled:opacity-50 transition-colors"
            >
              {{ submittingIzin ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
          <p v-if="izinMsg" :class="['text-xs text-center', izinErr ? 'text-red-500' : 'text-accent']">{{ izinMsg }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { get, post } from '../../api'
import { SearchIcon } from 'lucide-vue-next'

const authStore = useAuthStore()
const search = ref('')
const loading = ref(true)
const students = ref([])

const submittingIzin = ref(false)
const izinMsg = ref('')
const izinErr = ref(false)
const izinModal = reactive({
  open: false,
  studentId: '',
  studentName: '',
  studentNis: '',
  status: '',
  tanggal: new Date().toISOString().slice(0, 10)
})

function openIzin(s) {
  izinModal.open = true
  izinModal.studentId = s.id
  izinModal.studentName = s.name
  izinModal.studentNis = s.nis
  izinModal.status = ''
  izinModal.tanggal = new Date().toISOString().slice(0, 10)
  izinMsg.value = ''
  izinErr.value = false
}

async function submitIzin() {
  izinMsg.value = ''
  izinErr.value = false

  if (!izinModal.status) {
    izinMsg.value = 'Pilih status izin/sakit'
    izinErr.value = true
    return
  }

  submittingIzin.value = true
  try {
    await post('/absensi/izin', {
      student_id: izinModal.studentId,
      status: izinModal.status,
      tanggal: izinModal.tanggal || undefined
    })
    izinMsg.value = 'Absensi berhasil dicatat!'
    izinErr.value = false
    setTimeout(() => {
      izinModal.open = false
    }, 800)
  } catch (e) {
    izinMsg.value = e.message || 'Gagal mencatat absensi'
    izinErr.value = true
  } finally {
    submittingIzin.value = false
  }
}

const filteredStudents = computed(() => {
  return students.value.filter(s =>
    s.name.toLowerCase().includes(search.value.toLowerCase()) ||
    s.nis.includes(search.value)
  )
})

async function fetchData() {
  loading.value = true
  try {
    const dashRes = await get('/dudi/dashboard')
    const studentList = dashRes.students || []

    const enriched = await Promise.all(studentList.map(async (s) => {
      let attendance = 0
      let journals = 0
      let nilai = null

      try {
        const [nilaiRes, jrnRes] = await Promise.all([
          get('/nilai/' + s.id).catch(() => ({ data: null })),
          get('/jurnal?student_id=' + s.id).catch(() => ({ data: [] }))
        ])

        if (nilaiRes.data) {
          nilai = nilaiRes.data.final_grade || nilaiRes.data.FinalGrade
          attendance = Math.round(nilaiRes.data.attendance_score_auto || nilaiRes.data.AttendanceScoreAuto || 0)
        }

        journals = (jrnRes.data || []).length
      } catch (e) {
        // ignore individual errors
      }

      return {
        id: s.id,
        name: s.name,
        nis: s.nis,
        jurusan: s.jurusan || '-',
        attendance,
        journals,
        nilai
      }
    }))

    students.value = enriched
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>
