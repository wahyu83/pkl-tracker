<template>
  <div>
    <div class="mb-6 text-center">
      <div class="w-20 h-20 rounded-full bg-primary text-white flex items-center justify-center text-2xl font-bold mx-auto mb-3">
        {{ initial }}
      </div>
      <h2 class="text-lg font-bold text-gray-800">{{ authStore.userName }}</h2>
      <p class="text-sm text-gray-500">{{ authStore.userEmail }}</p>
      <span :class="['inline-flex px-3 py-1 rounded-full text-xs font-medium mt-2', roleBadge]">
        {{ roleLabel }}
      </span>
    </div>

    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <h3 class="font-semibold text-gray-800 mb-4">Informasi Akun</h3>
      <div class="space-y-3">
        <div class="flex justify-between text-sm">
          <span class="text-gray-500">NIS/NIP/NIK</span>
          <span class="font-medium text-gray-800 font-mono">{{ authStore.user?.nis_nip_nik || '-' }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-500">Email</span>
          <span class="font-medium text-gray-800">{{ authStore.userEmail }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-gray-500">Role</span>
          <span class="font-medium text-gray-800">{{ roleLabel }}</span>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-2xl p-5 border border-gray-100 mb-4">
      <h3 class="font-semibold text-gray-800 mb-4">Ubah Password</h3>
      <form @submit.prevent="handleChangePassword" class="space-y-3">
        <div>
          <label class="block text-xs font-medium text-gray-600 mb-1">Password Lama</label>
          <input
            v-model="oldPassword"
            type="password"
            placeholder="Masukkan password lama"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-gray-600 mb-1">Password Baru</label>
          <input
            v-model="newPassword"
            type="password"
            placeholder="Minimal 6 karakter"
            required
            minlength="6"
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-gray-600 mb-1">Konfirmasi Password Baru</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="Ulangi password baru"
            required
            class="w-full px-4 py-2.5 rounded-xl border border-gray-200 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
            :class="confirmPassword && newPassword !== confirmPassword ? 'border-danger' : ''"
          />
        </div>
        <p v-if="confirmPassword && newPassword !== confirmPassword" class="text-xs text-danger">
          Password tidak cocok
        </p>
        <p v-if="successMsg" class="text-xs text-accent font-medium">{{ successMsg }}</p>
        <p v-if="errorMsg" class="text-xs text-danger">{{ errorMsg }}</p>
        <button
          type="submit"
          :disabled="changing || (confirmPassword && newPassword !== confirmPassword)"
          class="w-full py-2.5 bg-primary text-white rounded-xl text-sm font-medium hover:bg-primary-light transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
        >
          <LoaderIcon v-if="changing" :size="16" class="animate-spin" />
          <span>{{ changing ? 'Menyimpan...' : 'Simpan Password' }}</span>
        </button>
      </form>
    </div>

    <button
      @click="handleLogout"
      class="w-full py-3 bg-danger text-white rounded-2xl text-sm font-bold hover:bg-red-600 transition-colors flex items-center justify-center gap-2"
    >
      <LogOutIcon :size="20" />
      Keluar
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { post } from '../../api'
import { LoaderIcon, LogOutIcon } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const changing = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const initial = computed(() => authStore.userName?.charAt(0)?.toUpperCase() || 'U')

const roleLabel = computed(() => {
  const map = { student: 'Siswa', teacher: 'Guru', dudi: 'DUDI', admin: 'Admin' }
  return map[authStore.userRole] || ''
})

const roleBadge = computed(() => {
  const map = {
    student: 'bg-primary/10 text-primary',
    teacher: 'bg-accent/10 text-accent',
    dudi: 'bg-warning/10 text-warning',
    admin: 'bg-gray-100 text-gray-600'
  }
  return map[authStore.userRole] || 'bg-gray-100 text-gray-500'
})

async function handleChangePassword() {
  if (newPassword.value !== confirmPassword.value) return
  changing.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await post('/change-password', {
      old_password: oldPassword.value,
      new_password: newPassword.value
    })
    successMsg.value = 'Password berhasil diubah!'
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e) {
    errorMsg.value = e.message
  } finally {
    changing.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
