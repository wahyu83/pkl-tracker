<template>
  <div>
    <div class="flex items-center justify-between mb-5">
      <div>
        <h2 class="text-lg font-bold text-gray-800">Jurnal PKL</h2>
        <p class="text-xs text-gray-500">{{ journals.length }} entri jurnal</p>
      </div>
      <button
        @click="showForm = !showForm"
        :class="[
          'flex items-center gap-1.5 px-4 py-2 rounded-xl text-sm font-medium transition-colors',
          showForm ? 'bg-gray-200 text-gray-700' : 'bg-accent text-white hover:bg-accent-dark'
        ]"
      >
        <PlusIcon v-if="!showForm" :size="18" />
        <XIcon v-else :size="18" />
        {{ showForm ? 'Batal' : 'Tulis' }}
      </button>
    </div>

    <!-- Form -->
    <div v-if="showForm" class="mb-5">
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div class="bg-white rounded-2xl p-5 border border-gray-100">
          <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Kegiatan <span class="text-xs text-gray-400">(maks. 10 hari kebelakang)</span></label>
          <input
            v-model="form.date"
            type="date"
            :max="todayStr"
            :min="maxBackStr"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
          />
        </div>

        <div class="bg-white rounded-2xl p-5 border border-gray-100">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Uraian Kegiatan <span class="text-danger">*</span>
          </label>
          <textarea
            v-model="form.activity"
            rows="5"
            placeholder="Ceritakan kegiatan PKL kamu hari ini..."
            required
            class="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none"
            maxlength="2000"
          />
          <p class="text-xs text-gray-400 text-right mt-1">{{ form.activity.length }}/2000</p>
        </div>

        <div class="bg-white rounded-2xl p-5 border border-gray-100">
          <label class="block text-sm font-medium text-gray-700 mb-1">Refleksi / Komentar Pribadi</label>
          <textarea
            v-model="form.reflection"
            rows="3"
            placeholder="Apa yang kamu pelajari hari ini? Ada kendala?"
            class="w-full px-4 py-3 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none"
            maxlength="1000"
          />
          <p class="text-xs text-gray-400 text-right mt-1">{{ form.reflection.length }}/1000</p>
        </div>

        <button
          type="submit"
          :disabled="loading || !form.activity"
          class="w-full py-3 bg-accent text-white rounded-2xl text-sm font-bold hover:bg-accent-dark transition-colors disabled:opacity-50 flex items-center justify-center gap-2 shadow-lg shadow-accent/20"
        >
          <LoaderIcon v-if="loading" :size="20" class="animate-spin" />
          <span>{{ loading ? 'Menyimpan...' : 'Simpan Jurnal' }}</span>
        </button>
      </form>

      <!-- Success Banner -->
      <div v-if="success" class="mt-4 bg-accent/10 rounded-xl p-3 flex items-center gap-2 text-sm text-accent font-medium">
        <CheckCircleIcon :size="18" />
        Jurnal berhasil disimpan!
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loadingList" class="text-center py-8 text-gray-400 text-sm">Memuat...</div>

    <!-- Empty -->
    <div v-else-if="journals.length === 0" class="text-center py-8">
      <p class="text-sm text-gray-400 mb-3">Belum ada jurnal. Mulai tulis sekarang!</p>
    </div>

    <!-- List -->
    <div v-else class="space-y-3">
      <div
        v-for="j in journals"
        :key="j.ID || j.id"
        class="bg-white rounded-xl border border-gray-100 overflow-hidden"
      >
        <div class="p-4">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs text-gray-400">{{ formatDate(j.Date || j.date) }}</span>
            <div class="flex items-center gap-2">
              <span :class="['text-[10px] font-medium px-2 py-0.5 rounded-full', hasComments(j) ? 'bg-accent/10 text-accent' : 'bg-gray-100 text-gray-500']">
                {{ hasComments(j) ? 'Dikomentari' : 'Menunggu' }}
              </span>
            </div>
          </div>

          <p class="text-sm text-gray-700 mb-2 line-clamp-3">{{ j.Activity || j.activity }}</p>

          <p v-if="j.Reflection || j.reflection" class="text-xs text-gray-500 italic mb-3">
            "{{ j.Reflection || j.reflection }}"
          </p>

          <!-- Comments -->
          <div v-if="hasComments(j)" class="space-y-1.5">
            <div v-if="j.TeacherComment || j.teacherComment" class="bg-info/5 rounded-lg px-3 py-2">
              <p class="text-[10px] text-gray-400">Guru Pembimbing:</p>
              <p class="text-xs text-gray-700">{{ j.TeacherComment || j.teacherComment }}</p>
            </div>
            <div v-if="j.DudiComment || j.dudiComment" class="bg-warning/5 rounded-lg px-3 py-2">
              <p class="text-[10px] text-gray-400">Instruktur:</p>
              <p class="text-xs text-gray-700">{{ j.DudiComment || j.dudiComment }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { PlusIcon, XIcon, LoaderIcon, CheckCircleIcon } from 'lucide-vue-next'

const showForm = ref(false)
const loading = ref(false)
const success = ref(false)
const loadingList = ref(true)
const journals = ref([])

const todayStr = computed(() => new Date().toISOString().split('T')[0])
const maxBackStr = computed(() => {
  const d = new Date()
  d.setDate(d.getDate() - 10)
  return d.toISOString().split('T')[0]
})

const form = reactive({
  date: new Date().toISOString().split('T')[0],
  activity: '',
  reflection: ''
})

function formatDate(d) {
  return new Date(d).toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

function hasComments(j) {
  return !!(j.TeacherComment || j.DudiComment || j.teacherComment || j.dudiComment)
}

async function fetchJournals() {
  loadingList.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/jurnal', {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Gagal memuat jurnal')
    const json = await res.json()
    journals.value = json.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingList.value = false
  }
}

async function handleSubmit() {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await fetch('/api/jurnal', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        date: form.date,
        activity: form.activity,
        reflection: form.reflection
      })
    })

    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'Gagal menyimpan jurnal')
    }

    success.value = true
    form.activity = ''
    form.reflection = ''
    form.date = new Date().toISOString().split('T')[0]
    showForm.value = false
    fetchJournals()
    setTimeout(() => { success.value = false }, 3000)
  } catch (e) {
    alert('Gagal: ' + e.message)
  } finally {
    loading.value = false
  }
}

onMounted(fetchJournals)
</script>
