<template>
  <div>
    <div class="flex items-center justify-between mb-5">
      <div>
        <h2 class="text-lg font-bold text-gray-800">Jurnal PKL</h2>
        <p class="text-xs text-gray-500">{{ journals.length }} entri jurnal</p>
      </div>
      <router-link
        to="/siswa/jurnal/tulis"
        class="flex items-center gap-1.5 px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent-dark transition-colors"
      >
        <PlusIcon :size="18" />
        Tulis
      </router-link>
    </div>

    <div class="space-y-3">
      <div
        v-for="j in journals"
        :key="j.id"
        class="bg-white rounded-xl border border-gray-100 overflow-hidden"
      >
        <div class="p-4">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs text-gray-400">{{ j.date }}</span>
            <div class="flex items-center gap-2">
              <span :class="['text-[10px] font-medium px-2 py-0.5 rounded-full', j.verified ? 'bg-accent/10 text-accent' : 'bg-gray-100 text-gray-500']">
                {{ j.verified ? 'Terverifikasi' : 'Menunggu' }}
              </span>
            </div>
          </div>

          <p class="text-sm text-gray-700 mb-2 line-clamp-3">{{ j.activity }}</p>

          <p v-if="j.reflection" class="text-xs text-gray-500 italic mb-3">
            "{{ j.reflection }}"
          </p>

          <!-- Comments -->
          <div v-if="j.teacherComment || j.dudiComment" class="space-y-1.5">
            <div v-if="j.teacherComment" class="bg-info/5 rounded-lg px-3 py-2">
              <p class="text-[10px] text-gray-400">Guru Pembimbing:</p>
              <p class="text-xs text-gray-700">{{ j.teacherComment }}</p>
            </div>
            <div v-if="j.dudiComment" class="bg-warning/5 rounded-lg px-3 py-2">
              <p class="text-[10px] text-gray-400">DUDI:</p>
              <p class="text-xs text-gray-700">{{ j.dudiComment }}</p>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { PlusIcon } from 'lucide-vue-next'

const journals = [
  {
    id: 1, date: '14 Mei 2026', activity: 'Mempelajari framework Laravel dan membuat CRUD sederhana untuk modul inventaris. Dibimbing oleh Pak Hendra.', reflection: 'Hari ini saya belajar banyak tentang MVC pattern.',
    teacherComment: 'Bagus, terus tingkatkan!', dudiComment: 'Siswa menunjukkan antusiasme yang baik.', verified: true
  },
  {
    id: 2, date: '13 Mei 2026', activity: 'Debugging aplikasi internal perusahaan. Memperbaiki bug pada modul pelaporan. Berhasil fix 3 bugs.', reflection: null,
    teacherComment: null, dudiComment: null, verified: true
  },
  {
    id: 3, date: '12 Mei 2026', activity: 'Meeting dengan tim developer, membahas sprint planning untuk minggu depan. Saya ditugaskan untuk modul reporting.',
    reflection: null, teacherComment: 'Terus komunikasi dengan tim ya.', dudiComment: null, verified: true
  },
  {
    id: 4, date: '11 Mei 2026', activity: 'Izin tidak masuk dikarenakan ada keperluan keluarga.', reflection: null,
    teacherComment: null, dudiComment: null, verified: false
  },
]
</script>
