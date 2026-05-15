<template>
  <div>
    <div class="flex items-center gap-3 mb-5">
      <router-link to="/siswa/jurnal" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
        <ArrowLeftIcon :size="20" />
      </router-link>
      <div>
        <h2 class="text-lg font-bold text-gray-800">Tulis Jurnal</h2>
        <p class="text-xs text-gray-500">{{ todayDate }}</p>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="bg-white rounded-2xl p-5 border border-gray-100">
        <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Kegiatan</label>
        <input
          v-model="form.date"
          type="date"
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

    <!-- Success -->
    <div
      v-if="success"
      class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl p-6 max-w-sm w-full text-center">
        <div class="w-16 h-16 bg-accent/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircleIcon :size="36" class="text-accent" />
        </div>
        <h3 class="text-lg font-bold text-gray-800 mb-2">Jurnal Tersimpan!</h3>
        <p class="text-sm text-gray-500 mb-4">Jurnal harian kamu telah disimpan</p>
        <router-link
          to="/siswa/jurnal"
          class="block w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors"
        >
          Lihat Jurnal
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ArrowLeftIcon, LoaderIcon, CheckCircleIcon } from 'lucide-vue-next'

const loading = ref(false)
const success = ref(false)

const todayDate = new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })

const form = reactive({
  date: new Date().toISOString().split('T')[0],
  activity: '',
  reflection: ''
})

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
  } catch (e) {
    alert('Gagal: ' + e.message)
  } finally {
    loading.value = false
  }
}
</script>
