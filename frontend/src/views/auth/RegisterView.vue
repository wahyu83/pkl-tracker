<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-dark p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-accent/30">
          <GraduationCap :size="32" class="text-white" />
        </div>
        <h1 class="text-2xl font-bold text-white">Daftar Akun</h1>
        <p class="text-gray-300 text-sm mt-1">Buat akun PKL Tracker</p>
      </div>

      <div class="bg-white rounded-2xl shadow-xl p-6 md:p-8">
        <h2 class="text-xl font-bold text-gray-800 text-center mb-6">Pendaftaran</h2>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="r in roles"
              :key="r.value"
              type="button"
              @click="role = r.value"
              :class="[
                'flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-xs font-medium border-2 transition-all',
                role === r.value
                  ? 'border-primary bg-primary/5 text-primary'
                  : 'border-gray-200 text-gray-500 hover:border-gray-300'
              ]"
            >
              <component :is="r.icon" :size="16" />
              <span>{{ r.label }}</span>
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
            <input
              v-model="fullName"
              type="text"
              placeholder="Masukkan nama lengkap"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ credentialLabel }}</label>
            <input
              v-model="credential"
              type="text"
              :placeholder="'Masukkan ' + credentialLabel"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input
              v-model="email"
              type="email"
              placeholder="Masukkan email"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input
              v-model="password"
              type="password"
              placeholder="Minimal 8 karakter"
              required
              minlength="8"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Konfirmasi Password</label>
            <input
              v-model="confirmPassword"
              type="password"
              placeholder="Ulangi password"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              :class="confirmPassword && password !== confirmPassword ? 'border-danger focus:ring-danger/20' : ''"
            />
            <p v-if="confirmPassword && password !== confirmPassword" class="text-danger text-xs mt-1">
              Password tidak cocok
            </p>
          </div>

          <button
            type="submit"
            :disabled="loading || (!!confirmPassword && password !== confirmPassword)"
            class="w-full py-2.5 bg-primary text-white rounded-xl font-medium text-sm hover:bg-primary-light transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <LoaderIcon v-if="loading" :size="18" class="animate-spin" />
            <span>{{ loading ? 'Mendaftar...' : 'Daftar' }}</span>
          </button>
        </form>

        <p class="text-center text-sm text-gray-500 mt-5">
          Sudah punya akun?
          <router-link to="/login" class="text-primary font-medium hover:underline">
            Masuk sekarang
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { GraduationCap, LoaderIcon, Users, UserCircle, Building2, Shield } from 'lucide-vue-next'

const router = useRouter()

const role = ref('student')
const fullName = ref('')
const credential = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

const roles = [
  { value: 'student', label: 'Siswa', icon: GraduationCap },
  { value: 'teacher', label: 'Guru', icon: UserCircle },
  { value: 'dudi', label: 'DUDI', icon: Building2 },
  { value: 'admin', label: 'Admin', icon: Shield },
]

const credentialLabel = computed(() => {
  const map = { student: 'NIS', teacher: 'NIP', dudi: 'NIK', admin: 'NIP/Username' }
  return map[role.value]
})

function handleRegister() {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    router.push('/login')
  }, 800)
}
</script>
