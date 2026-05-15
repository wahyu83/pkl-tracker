<template>
  <div>
    <div class="mb-6">
      <h2 class="text-xl font-bold text-gray-800">Monitoring Jurnal</h2>
      <p class="text-sm text-gray-500 mt-0.5">Baca dan beri komentar pada jurnal siswa</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Student list -->
      <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
        <div class="p-4 border-b border-gray-100">
          <h3 class="font-semibold text-gray-800">Pilih Siswa</h3>
        </div>
        <div class="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
          <button
            v-for="s in students"
            :key="s.id"
            @click="selectedStudent = s.id"
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
        </div>
      </div>

      <!-- Journal list -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-gray-100">
        <div class="p-4 border-b border-gray-100">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-gray-800">
              {{ selectedStudentName }} - Jurnal Harian
            </h3>
            <span class="text-xs text-gray-400">{{ journals.length }} entri</span>
          </div>
        </div>

        <div class="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
          <div v-if="journals.length === 0" class="p-8 text-center text-gray-400 text-sm">
            Pilih siswa untuk melihat jurnal
          </div>

          <div
            v-for="j in journals"
            :key="j.id"
            class="p-4"
          >
            <div class="flex items-start justify-between mb-2">
              <div>
                <span class="text-sm font-semibold text-gray-800">{{ j.date }}</span>
                <span :class="['ml-2 inline-flex px-2 py-0.5 rounded-full text-[10px] font-medium', j.verified ? 'bg-accent/10 text-accent' : 'bg-gray-100 text-gray-500']">
                  {{ j.verified ? 'Terverifikasi' : 'Menunggu' }}
                </span>
              </div>
            </div>

            <p class="text-sm text-gray-700 mb-3">{{ j.activity }}</p>

            <p v-if="j.reflection" class="text-sm text-gray-500 italic mb-3">
              "{{ j.reflection }}"
            </p>

            <!-- Teacher comment -->
            <div v-if="j.teacherComment" class="mb-3 ml-3 pl-3 border-l-2 border-info/30">
              <p class="text-xs text-gray-400 mb-0.5">Komentar Anda:</p>
              <p class="text-sm text-gray-700">{{ j.teacherComment }}</p>
            </div>

            <!-- DUDI comment -->
            <div v-if="j.dudiComment" class="mb-3 ml-3 pl-3 border-l-2 border-warning/30">
              <p class="text-xs text-gray-400 mb-0.5">Komentar DUDI:</p>
              <p class="text-sm text-gray-700">{{ j.dudiComment }}</p>
            </div>

            <!-- Comment input -->
            <div v-if="!j.teacherComment" class="flex items-center gap-2 mt-2">
              <input
                v-model="j._newComment"
                type="text"
                placeholder="Tulis komentar..."
                class="flex-1 px-3 py-2 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                @keyup.enter="addComment(j)"
              />
              <button
                @click="addComment(j)"
                :disabled="!j._newComment?.trim()"
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
import { ref, computed } from 'vue'

const selectedStudent = ref(1)

const students = [
  { id: 1, name: 'Ahmad Rizky', totalJurnal: 42 },
  { id: 2, name: 'Siti Nurhaliza', totalJurnal: 40 },
  { id: 3, name: 'Dian Permata', totalJurnal: 35 },
]

const selectedStudentName = computed(() => {
  return students.find(s => s.id === selectedStudent.value)?.name || ''
})

const journals = ref([
  {
    id: 1, date: '14 Mei 2026', activity: 'Mempelajari framework Laravel dan membuat CRUD sederhana untuk modul inventaris. Dibimbing oleh Pak Hendra.',
    reflection: 'Hari ini saya belajar banyak tentang MVC pattern dan bagaimana Laravel mempermudah development.',
    teacherComment: null, dudiComment: 'Siswa menunjukkan antusiasme yang baik.', verified: true,
    _newComment: ''
  },
  {
    id: 2, date: '13 Mei 2026', activity: 'Debugging aplikasi internal perusahaan. Memperbaiki bug pada modul pelaporan.',
    reflection: 'Debugging ternyata butuh ketelitian tinggi. Saya belajar menggunakan breakpoint di VS Code.',
    teacherComment: null, dudiComment: null, verified: true,
    _newComment: ''
  },
  {
    id: 3, date: '12 Mei 2026', activity: 'Meeting dengan tim developer, membahas sprint planning minggu depan.',
    reflection: null, teacherComment: 'Bagus, terus tingkatkan komunikasi dengan tim.',
    dudiComment: null, verified: true, _newComment: ''
  },
])

function addComment(j) {
  if (j._newComment?.trim()) {
    j.teacherComment = j._newComment.trim()
    j._newComment = ''
  }
}
</script>
