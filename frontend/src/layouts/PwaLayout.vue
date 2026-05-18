<template>
  <div class="flex flex-col min-h-screen bg-gray-50">
    <!-- Top bar -->
    <header class="sticky top-0 z-30 bg-white border-b border-gray-200 safe-top">
      <div class="flex items-center justify-between px-4 h-14">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-xs">PKL</span>
          </div>
          <h1 class="font-bold text-sm text-gray-800 truncate">
            {{ pageTitle }}
          </h1>
        </div>
        <div class="flex items-center gap-2">
          <button class="relative p-1.5 rounded-lg hover:bg-gray-100 text-gray-500">
            <BellIcon :size="20" />
            <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>
          <div
            class="w-7 h-7 rounded-full bg-primary text-white flex items-center justify-center text-xs font-bold cursor-pointer hover:opacity-80 transition-opacity"
            @click="goToProfile"
          >
            {{ initial }}
          </div>
        </div>
      </div>
    </header>

    <!-- Main scrollable content -->
    <main class="flex-1 overflow-y-auto p-4 pb-24">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Bottom Navigation Bar -->
    <nav class="fixed bottom-0 inset-x-0 z-30 bg-white border-t border-gray-200 pb-safe">
      <div class="flex items-center justify-around h-16">
        <router-link
          v-for="item in menuItems"
          :key="item.to"
          :to="item.to"
          class="flex flex-col items-center justify-center gap-0.5 flex-1 h-full transition-colors"
          :class="$route.path === item.to ? 'bottom-nav-active' : 'bottom-nav-inactive'"
        >
          <component :is="item.icon" :size="22" :stroke-width=" $route.path === item.to ? 2.5 : 2 " />
          <span class="text-[10px] font-medium">{{ item.label }}</span>
        </router-link>
      </div>
    </nav>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import {
  LayoutDashboard, ClipboardCheck, BookOpen,
  Users, PenTool, BellIcon, UserIcon, Award
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const initial = computed(() => {
  return authStore.userName?.charAt(0)?.toUpperCase() || 'U'
})

const role = computed(() => authStore.userRole)

const pageTitle = computed(() => {
  const titles = {
    'SiswaDashboard': 'Dashboard',
    'SiswaAbsensi': 'Absensi',
    'SiswaAbsensiHistory': 'Riwayat Absensi',
    'SiswaJurnal': 'Jurnal',
    'SiswaJurnalTulis': 'Tulis Jurnal',
    'SiswaProfile': 'Profil',
    'DudiDashboard': 'Dashboard',
    'DudiSiswa': 'Daftar Siswa',
    'DudiPenilaian': 'Penilaian',
    'DudiJurnal': 'Jurnal Siswa',
    'DudiProfile': 'Profil',
    'GuruDashboard': 'Dashboard',
    'GuruAbsensi': 'Absensi',
    'GuruJurnal': 'Jurnal Siswa',
    'GuruNilai': 'Nilai PKL',
    'GuruReports': 'Rekap',
    'GuruProfile': 'Profil'
  }
  return titles[route.name] || 'PKL Tracker'
})

const siswaMenu = [
  { to: '/siswa', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/siswa/absensi', label: 'Absensi', icon: ClipboardCheck },
  { to: '/siswa/jurnal', label: 'Jurnal', icon: BookOpen },
  { to: '/siswa/profile', label: 'Profil', icon: UserIcon },
]

const dudiMenu = [
  { to: '/dudi', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/dudi/siswa', label: 'Siswa', icon: Users },
  { to: '/dudi/jurnal', label: 'Jurnal', icon: BookOpen },
  { to: '/dudi/penilaian', label: 'Nilai', icon: PenTool },
  { to: '/dudi/profile', label: 'Profil', icon: UserIcon },
]

const guruMenu = [
  { to: '/guru', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/guru/absensi', label: 'Absensi', icon: ClipboardCheck },
  { to: '/guru/jurnal', label: 'Jurnal', icon: BookOpen },
  { to: '/guru/nilai', label: 'Nilai', icon: Award },
  { to: '/guru/profile', label: 'Profil', icon: UserIcon },
]

const menuItems = computed(() => {
  if (role.value === 'student') return siswaMenu
  if (role.value === 'dudi') return dudiMenu
  if (role.value === 'teacher') return guruMenu
  return siswaMenu
})

function goToProfile() {
  const path = role.value === 'dudi' ? '/dudi/profile' : role.value === 'teacher' ? '/guru/profile' : '/siswa/profile'
  router.push(path)
}
</script>

<style scoped>
.safe-top {
  padding-top: env(safe-area-inset-top, 0px);
}
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
