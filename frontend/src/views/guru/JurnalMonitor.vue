<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Monitoring Jurnal</h2>
      <p class="text-sm text-gray-500 mt-0.5">Baca dan beri komentar pada jurnal siswa</p>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-400 text-sm">Memuat data...</div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
        <div class="p-4 border-b border-gray-100">
          <h3 class="font-semibold text-gray-800">Pilih Siswa</h3>
        </div>
        <div class="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
          <button
            v-for="s in students"
            :key="s.id"
            @click="selectStudent(s.id)"
            :class="[
              'w-full text-left px-4 py-3 flex items-center gap-3 transition-colors',
              selectedStudent === s.id ? 'bg-primary/5 border-l-2 border-primary' : 'hover:bg-gray-50 border-l-2 border-transparent'
            ]"
          >
            <div class="w-9 h-9 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-bold flex-shrink-0">
              {{ s.name.charAt(0) }}
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium text-gray-800 truncate">{{ s.name }}</p>
              <p class="text-xs text-gray-500">{{ s.totalJurnal }} jurnal</p>
            </div>
          </button>
          <div v-if="students.length === 0" class="p-4 text-center text-gray-400 text-sm">
            Belum ada siswa
          </div>
        </div>
      </div>

      <div class="lg:col-span-2 bg-white rounded-2xl border border-gray-100">
        <div class="p-4 border-b border-gray-100">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-gray-800">
              {{ selectedStudentName || 'Pilih Siswa' }} - Jurnal Harian
            </h3>
            <span class="text-xs text-gray-400">{{ journals.length }} entri</span>
          </div>
        </div>

        <div class="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
          <div v-if="!selectedStudent" class="p-8 text-center text-gray-400 text-sm">
            Pilih siswa untuk melihat jurnal
          </div>
          <div v-else-if="journals.length === 0" class="p-8 text-center text-gray-400 text-sm">
            Belum ada jurnal dari siswa ini
          </div>

          <div
            v-for="j in journals"
            :key="j.ID || j.id"
            class="p-4"
          >
            <div class="flex items-start justify-between mb-2">
              <div>
                <span class="text-sm font-semibold text-gray-800">{{ formatDate(j.Date) }}</span>
              </div>
            </div>

            <p class="text-sm text-gray-700 mb-3">{{ j.Activity }}</p>

            <p v-if="j.Reflection" class="text-sm text-gray-500 italic mb-3">
              "{{ j.Reflection }}"
            </p>

            <div v-if="j.TeacherComment" class="mb-3 ml-3 pl-3 border-l-2 border-info/30">
              <p class="text-xs text-gray-400 mb-0.5">Komentar Anda:</p>
              <p class="text-sm text-gray-700">{{ j.TeacherComment }}</p>
            </div>

            <div v-if="j.DudiComment" class="mb-3 ml-3 pl-3 border-l-2 border-warning/30">
              <p class="text-xs text-gray-400 mb-0.5">Komentar Instruktur:</p>
              <p class="text-sm text-gray-700">{{ j.DudiComment }}</p>
            </div>

            <div v-if="!j.TeacherComment" class="flex items-center gap-2 mt-2">
              <input
                v-model="commentInputs[j.ID || j.id]"
                type="text"
                placeholder="Tulis komentar..."
                class="flex-1 px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                @keyup.enter="addComment(j)"
              />
              <button
                @click="addComment(j)"
                :disabled="!commentInputs[j.ID || j.id]?.trim() || commenting"
                class="px-4 py-2 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors disabled:opacity-50"
              >
                Kirim
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { get, post } from '@/api'

const students = ref([])
const journals = ref([])
const selectedStudent = ref(null)
const loading = ref(true)
const commenting = ref(false)
const commentInputs = reactive({})

const selectedStudentName = computed(() => {
  return students.value.find(s => s.id === selectedStudent.value)?.name || ''
})

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

async function fetchStudents() {
  try {
    const res = await get('/guru/students')
    const users = res.data || []
    const result = []
    for (const u of users) {
      result.push({
        id: u.id,
        name: u.full_name,
        totalJurnal: u.journal_count || 0
      })
    }
    students.value = result
    if (result.length > 0 && !selectedStudent.value) {
      selectStudent(result[0].id)
    }
  } catch (e) {
    console.error(e)
  }
}

async function selectStudent(id) {
  selectedStudent.value = id
  loading.value = true
  try {
    const res = await get('/jurnal?student_id=' + id)
    journals.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function addComment(j) {
  const jid = j.ID || j.id
  const comment = commentInputs[jid]?.trim()
  if (!comment) return
  commenting.value = true
  try {
    await post('/jurnal/comment', { jurnal_id: jid, comment })
    j.TeacherComment = comment
    commentInputs[jid] = ''
  } catch (e) {
    alert(e.message)
  } finally {
    commenting.value = false
  }
}

onMounted(() => { fetchStudents() })
</script>
