<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-dark p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-accent/30">
          <GraduationCap :size="32" class="text-white" />
        </div>
        <h1 class="text-2xl font-bold text-white">Lupa Password</h1>
        <p class="text-gray-300 text-sm mt-1">Reset password akun Anda</p>
      </div>

      <div class="bg-white rounded-2xl shadow-xl p-6 md:p-8">
        <template v-if="!sent">
          <h2 class="text-xl font-bold text-gray-800 text-center mb-2">Reset Password</h2>
          <p class="text-sm text-gray-500 text-center mb-6">
            Masukkan email terdaftar untuk menerima link reset password.
          </p>

          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                v-model="email"
                type="email"
                placeholder="Masukkan email terdaftar"
                required
                class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              />
            </div>

            <button
              type="submit"
              :disabled="loading"
              class="w-full py-2.5 bg-primary text-white rounded-xl font-medium text-sm hover:bg-primary-light transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <LoaderIcon v-if="loading" :size="18" class="animate-spin" />
              <span>{{ loading ? 'Mengirim...' : 'Kirim Link Reset' }}</span>
            </button>
          </form>
        </template>

        <template v-else>
          <div class="text-center py-4">
            <div class="w-14 h-14 bg-accent/10 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircleIcon :size="32" class="text-accent" />
            </div>
            <h3 class="text-lg font-bold text-gray-800 mb-2">Email Terkirim!</h3>
            <p class="text-sm text-gray-500 mb-6">
              Link reset password telah dikirim ke <strong>{{ email }}</strong>. Silakan periksa inbox Anda.
            </p>
            <button
              @click="sent = false; email = ''"
              class="text-primary font-medium text-sm hover:underline"
            >
              Kirim ulang
            </button>
          </div>
        </template>

        <p class="text-center text-sm text-gray-500 mt-5">
          <router-link to="/login" class="text-primary font-medium hover:underline">
            Kembali ke login
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { GraduationCap, LoaderIcon, CheckCircleIcon } from 'lucide-vue-next'

const email = ref('')
const loading = ref(false)
const sent = ref(false)

function handleSubmit() {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    sent.value = true
  }, 1000)
}
</script>
