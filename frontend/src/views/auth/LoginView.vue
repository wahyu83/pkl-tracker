<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-dark p-4">
    <div class="w-full max-w-md">
      <!-- Logo & Header -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-accent/30">
          <GraduationCap :size="32" class="text-white" />
        </div>
        <h1 class="text-2xl font-bold text-white">PKL Tracker</h1>
        <p class="text-gray-200 text-base mt-1 font-medium">SMKN 1 Arahan</p>
        <p class="text-gray-300 text-sm mt-0.5">Indramayu - Jawa Barat</p>
      </div>

      <!-- Card -->
      <div class="bg-white rounded-2xl shadow-xl p-6 md:p-8">
        <h2 class="text-xl font-bold text-gray-800 text-center mb-6">Masuk</h2>

        <!-- Role Selector -->
        <div class="mb-5">
          <label class="block text-sm font-medium text-gray-700 mb-2">Masuk sebagai</label>
          <div class="grid grid-cols-4 gap-2">
            <button
              v-for="r in roles"
              :key="r.value"
              @click="role = r.value"
              :class="[
                'flex flex-col items-center gap-1 px-2 py-3 rounded-xl text-xs font-medium border-2 transition-all',
                role === r.value
                  ? 'border-primary bg-primary/5 text-primary'
                  : 'border-gray-200 text-gray-500 hover:border-gray-300'
              ]"
            >
              <component :is="r.icon" :size="18" />
              <span>{{ r.label }}</span>
            </button>
          </div>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              {{ credentialLabel }}
            </label>
            <input
              v-model="credential"
              type="text"
              :placeholder="'Masukkan NIS/NIP/NIK'"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Masukkan password"
                required
                class="w-full px-4 py-2.5 rounded-xl border border-gray-300 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all pr-10"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                <EyeOffIcon v-if="showPassword" :size="18" />
                <EyeIcon v-else :size="18" />
              </button>
            </div>
          </div>

          <div class="flex items-center justify-between text-sm">
            <label class="flex items-center gap-2 text-gray-600 cursor-pointer">
              <input type="checkbox" v-model="rememberMe" class="rounded border-gray-300 text-primary focus:ring-primary" />
              Ingat saya
            </label>
            <router-link to="/forgot-password" class="text-primary hover:underline font-medium">
              Lupa password?
            </router-link>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 bg-primary text-white rounded-xl font-medium text-sm hover:bg-primary-light transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <LoaderIcon v-if="loading" :size="18" class="animate-spin" />
            <span>{{ loading ? 'Memproses...' : 'Masuk' }}</span>
          </button>
        </form>

        <p class="text-center text-xs text-gray-400 mt-5">
          Developed by WA-13
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { GraduationCap, EyeIcon, EyeOffIcon, LoaderIcon, Users, UserCircle, Building2, Shield, ShieldCheck } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const role = ref('student')
const credential = ref('')
const password = ref('')
const showPassword = ref(false)
const rememberMe = ref(false)
const loading = ref(false)

const roles = [
  { value: 'student', label: 'Siswa', icon: GraduationCap },
  { value: 'teacher', label: 'Guru', icon: UserCircle },
  { value: 'dudi', label: 'Instruktur', icon: Building2 },
  { value: 'admin', label: 'Admin', icon: Shield },
  { value: 'admin_jurusan', label: 'Admin Jurusan', icon: ShieldCheck },
]

const credentialLabel = computed(() => {
  const map = { student: 'NIS', teacher: 'NIP', dudi: 'NIK', admin: 'NIP/Username', admin_jurusan: 'NIP/Username' }
  return map[role.value]
})

const roleRedirect = {
  admin: '/admin',
  admin_jurusan: '/jurusan',
  teacher: '/guru',
  student: '/siswa',
  dudi: '/dudi'
}

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login({
      nis_nip_nik: credential.value,
      password: password.value
    })
    router.push(roleRedirect[authStore.userRole] || '/siswa')
  } catch (e) {
    alert('Login gagal: ' + e.message)
  } finally {
    loading.value = false
  }
}
</script>
