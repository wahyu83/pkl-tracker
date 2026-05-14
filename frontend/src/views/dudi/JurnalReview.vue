<template>
  <div>
    <div class="mb-5">
      <h2 class="text-lg font-bold text-gray-800">Review Jurnal</h2>
      <p class="text-xs text-gray-500 mt-0.5">Baca dan beri komentar pada jurnal siswa</p>
    </div>

    <!-- Student selector -->
    <div class="bg-white rounded-xl p-3 border border-gray-100 mb-4">
      <select
        v-model="selectedStudent"
        class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm bg-white focus:border-primary outline-none"
      >
        <option value="">Semua Siswa</option>
        <option v-for="s in students" :key="s.id" :value="s.id">{{ s.name }}</option>
      </select>
    </div>

    <div class="space-y-3">
      <div
        v-for="j in filteredJournals"
        :key="j.id"
        class="bg-white rounded-xl p-4 border border-gray-100"
      >
        <div class="flex items-start justify-between mb-2">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
              {{ j.student.charAt(0) }}
            </div>
            <div>
              <p class="text-sm font-medium text-gray-800">{{ j.student }}</p>
              <p class="text-[10px] text-gray-400">{{ j.date }}</p>
            </div>
          </div>
        </div>

        <p class="text-sm text-gray-700 mb-3">{{ j.activity }}</p>

        <!-- Teacher comment -->
        <div v-if="j.teacherComment" class="mb-2 ml-2 pl-3 border-l-2 border-info/30">
          <p class="text-[10px] text-gray-400">Komentar Guru:</p>
          <p class="text-xs text-gray-700">{{ j.teacherComment }}</p>
        </div>

        <!-- DUDI comment -->
        <div v-if="j.dudiComment" class="mb-2 ml-2 pl-3 border-l-2 border-warning/30">
          <p class="text-[10px] text-gray-400">Komentar Anda:</p>
          <p class="text-xs text-gray-700">{{ j.dudiComment }}</p>
        </div>

        <!-- Comment form -->
        <div v-if="!j.dudiComment" class="mt-2">
          <div class="flex items-center gap-2">
            <input
              v-model="j._newComment"
              type="text"
              placeholder="Tulis komentar..."
              class="flex-1 px-3 py-2 rounded-lg border border-gray-200 text-xs focus:border-primary outline-none"
              @keyup.enter="addComment(j)"
            />
            <button
              @click="addComment(j)"
              :disabled="!j._newComment?.trim()"
              class="px-3 py-2 bg-primary text-white rounded-lg text-xs font-medium hover:bg-primary-light disabled:opacity-50 transition-colors"
            >
              Kirim
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const selectedStudent = ref('')
const students = [
  { id: 1, name: 'Ahmad Rizky' },
  { id: 2, name: 'Siti Nurhaliza' },
  { id: 3, name: 'Dian Permata' },
  { id: 4, name: 'Rudi Hartono' },
  { id: 5, name: 'Maya Sari' },
]

const journals = ref([
  {
    id: 1, studentId: 1, student: 'Ahmad Rizky', date: '14 Mei 2026',
    activity: 'Mempelajari framework Laravel dan membuat CRUD sederhana untuk modul inventaris. Dibimbing oleh Pak Hendra.',
    teacherComment: 'Bagus, terus tingkatkan!', dudiComment: 'Siswa menunjukkan antusiasme yang baik.', _newComment: ''
  },
  {
    id: 2, studentId: 2, student: 'Siti Nurhaliza', date: '14 Mei 2026',
    activity: 'Membantu tim network melakukan maintenance server dan pengecekan kabel UTP.',
    teacherComment: null, dudiComment: null, _newComment: ''
  },
  {
    id: 3, studentId: 1, student: 'Ahmad Rizky', date: '13 Mei 2026',
    activity: 'Debugging aplikasi internal perusahaan. Berhasil fix 3 bugs pada modul pelaporan.',
    teacherComment: null, dudiComment: null, _newComment: ''
  },
  {
    id: 4, studentId: 3, student: 'Dian Permata', date: '13 Mei 2026',
    activity: 'Membuat desain UI untuk aplikasi mobile client menggunakan Figma.',
    teacherComment: 'Desain sudah bagus, pertimbangkan accessibility.', dudiComment: null, _newComment: ''
  },
])

const filteredJournals = computed(() => {
  return journals.value.filter(j =>
    !selectedStudent.value || j.studentId === Number(selectedStudent.value)
  )
})

function addComment(j) {
  if (j._newComment?.trim()) {
    j.dudiComment = j._newComment.trim()
    j._newComment = ''
  }
}
</script>
