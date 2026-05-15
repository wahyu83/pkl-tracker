<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <!-- Mobile overlay -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-20 bg-black/50 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed lg:static inset-y-0 left-0 z-30 w-64 bg-primary-dark text-white flex flex-col transition-transform duration-300',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
      ]"
    >
      <div class="flex items-center gap-3 px-6 py-5 border-b border-white/10">
        <div class="w-9 h-9 bg-accent rounded-lg flex items-center justify-center flex-shrink-0">
          <span class="text-white font-bold text-sm">PKL</span>
        </div>
        <div class="min-w-0">
          <h1 class="font-bold text-sm truncate">PKL Tracker</h1>
          <p class="text-xs text-gray-400 truncate capitalize">{{ roleLabel }}</p>
        </div>
      </div>

      <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1">
        <template v-for="item in menuItems" :key="item.name">
          <p
            v-if="item.group"
            class="px-3 pt-4 pb-1 text-xs font-semibold text-gray-400 uppercase tracking-wider"
          >
            {{ item.group }}
          </p>
          <router-link
            v-else
            :to="item.to"
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors"
            :class="$route.path === item.to ? 'bg-primary text-white' : 'text-gray-300 hover:bg-white/10'"
          >
            <component :is="item.icon" :size="20" />
            <span>{{ item.label }}</span>
            <span
              v-if="item.badge"
              class="ml-auto bg-red-500 text-white text-xs px-1.5 py-0.5 rounded-full"
            >
              {{ item.badge }}
            </span>
          </router-link>
        </template>
      </nav>

      <div class="p-4 border-t border-white/10">
        <router-link
          :to="profilePath"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-white/10 transition-colors"
        >
          <UserIcon :size="20" />
          <span>{{ authStore.userName }}</span>
        </router-link>
        <button
          @click="handleLogout"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-400 hover:text-red-400 hover:bg-white/5 transition-colors w-full mt-1"
        >
          <LogOutIcon :size="20" />
          <span>Keluar</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top bar -->
      <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-4 flex-shrink-0">
        <button
          class="lg:hidden p-1.5 rounded-lg hover:bg-gray-100 text-gray-600"
          @click="sidebarOpen = !sidebarOpen"
        >
          <MenuIcon :size="22" />
        </button>

        <h2 class="text-sm font-semibold text-gray-800 truncate ml-2 lg:ml-0">
          {{ pageTitle }}
        </h2>

        <div class="flex items-center gap-3">
          <button class="relative p-1.5 rounded-lg hover:bg-gray-100 text-gray-500">
            <BellIcon :size="20" />
            <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>
          <div class="hidden sm:flex items-center gap-2 text-sm text-gray-600">
            <div class="w-7 h-7 rounded-full bg-primary text-white flex items-center justify-center text-xs font-bold">
              {{ initial }}
            </div>
            <span class="truncate max-w-[120px]">{{ authStore.userName }}</span>
          </div>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-4 md:p-6">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import {
  LayoutDashboard, Users, Building2, FileText, ClipboardCheck,
  BookOpen, Award, BarChart3, UserIcon, LogOutIcon, MenuIcon, BellIcon,
  GraduationCap, MapPin, PenTool, Calendar
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const sidebarOpen = ref(false)

const roleLabel = computed(() => {
  const map = { admin: 'Admin', teacher: 'Guru' }
  return map[authStore.userRole] || ''
})

const profilePath = computed(() => {
  const map = { admin: '/admin/profile', teacher: '/guru/profile' }
  return map[authStore.userRole] || '/admin/profile'
})

const initial = computed(() => {
  return authStore.userName?.charAt(0)?.toUpperCase() || 'U'
})

const role = computed(() => authStore.userRole)

const menuItems = computed(() => {
  if (role.value === 'admin') return adminMenu
  if (role.value === 'teacher') return guruMenu
  return []
})

const pageTitle = computed(() => {
  return route.meta.title || route.name || 'Dashboard'
})

const adminMenu = [
  { group: 'Utama' },
  { to: '/admin', label: 'Dashboard', icon: LayoutDashboard },
  { group: 'Manajemen' },
  { to: '/admin/users', label: 'Pengguna', icon: Users },
  { to: '/admin/dudi', label: 'Data DUDI', icon: Building2 },
  { to: '/admin/periode', label: 'Periode', icon: Calendar },
  { group: 'Laporan' },
  { to: '/admin/reports', label: 'Rekap & Laporan', icon: FileText },
]

const guruMenu = [
  { group: 'Utama' },
  { to: '/guru', label: 'Dashboard', icon: LayoutDashboard },
  { group: 'Monitoring' },
  { to: '/guru/absensi', label: 'Absensi Siswa', icon: ClipboardCheck },
  { to: '/guru/jurnal', label: 'Jurnal Siswa', icon: BookOpen },
  { to: '/guru/nilai', label: 'Nilai PKL', icon: Award },
  { group: 'Laporan' },
  { to: '/guru/reports', label: 'Rekap & Laporan', icon: FileText },
]

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
