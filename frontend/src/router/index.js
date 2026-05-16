import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/LoginView.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/RegisterView.vue'),
    meta: { guest: true }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('../views/auth/ForgotPasswordView.vue'),
    meta: { guest: true }
  },

  // Admin Jurusan Routes (Web Layout)
  {
    path: '/jurusan',
    component: () => import('../layouts/WebLayout.vue'),
    meta: { requiresAuth: true, role: 'admin_jurusan' },
    children: [
      {
        path: '',
        name: 'JurusanDashboard',
        component: () => import('../views/admin/AdminDashboard.vue'),
        meta: { title: 'Dashboard Jurusan' }
      },
      {
        path: 'users',
        name: 'JurusanUsers',
        component: () => import('../views/admin/UserManagement.vue'),
        meta: { title: 'Manajemen Pengguna' }
      },
      {
        path: 'dudi',
        name: 'JurusanDudi',
        component: () => import('../views/admin/DudiManagement.vue'),
        meta: { title: 'Data DUDI' }
      },
      {
        path: 'periode',
        name: 'JurusanPeriode',
        component: () => import('../views/admin/PeriodeManagement.vue'),
        meta: { title: 'Periode' }
      },
      {
        path: 'reports',
        name: 'JurusanReports',
        component: () => import('../views/admin/ReportsView.vue'),
        meta: { title: 'Rekap & Laporan' }
      },
      {
        path: 'profile',
        name: 'JurusanProfile',
        component: () => import('../views/siswa/ProfileView.vue'),
        meta: { title: 'Profil' }
      }
    ]
  },

  // Admin Routes (Web Layout)
  {
    path: '/admin',
    component: () => import('../layouts/WebLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('../views/admin/AdminDashboard.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('../views/admin/UserManagement.vue')
      },
      {
        path: 'dudi',
        name: 'AdminDudi',
        component: () => import('../views/admin/DudiManagement.vue')
      },
      {
        path: 'periode',
        name: 'AdminPeriode',
        component: () => import('../views/admin/PeriodeManagement.vue')
      },
      {
        path: 'reports',
        name: 'AdminReports',
        component: () => import('../views/admin/ReportsView.vue')
      },
      {
        path: 'profile',
        name: 'AdminProfile',
        component: () => import('../views/siswa/ProfileView.vue')
      }
    ]
  },

  // Guru Pembimbing Routes (Web Layout)
  {
    path: '/guru',
    component: () => import('../layouts/WebLayout.vue'),
    meta: { requiresAuth: true, role: 'teacher' },
    children: [
      {
        path: '',
        name: 'GuruDashboard',
        component: () => import('../views/guru/GuruDashboard.vue')
      },
      {
        path: 'absensi',
        name: 'GuruAbsensi',
        component: () => import('../views/guru/AbsensiMonitor.vue')
      },
      {
        path: 'jurnal',
        name: 'GuruJurnal',
        component: () => import('../views/guru/JurnalMonitor.vue')
      },
      {
        path: 'nilai',
        name: 'GuruNilai',
        component: () => import('../views/guru/NilaiSiswaView.vue')
      },
      {
        path: 'reports',
        name: 'GuruReports',
        component: () => import('../views/guru/ReportsView.vue')
      },
      {
        path: 'profile',
        name: 'GuruProfile',
        component: () => import('../views/siswa/ProfileView.vue')
      }
    ]
  },

  // Siswa Routes (PWA Layout)
  {
    path: '/siswa',
    component: () => import('../layouts/PwaLayout.vue'),
    meta: { requiresAuth: true, role: 'student' },
    children: [
      {
        path: '',
        name: 'SiswaDashboard',
        component: () => import('../views/siswa/SiswaDashboard.vue')
      },
      {
        path: 'absensi',
        name: 'SiswaAbsensi',
        component: () => import('../views/siswa/AbsensiView.vue')
      },
      {
        path: 'absensi/history',
        name: 'SiswaAbsensiHistory',
        component: () => import('../views/siswa/AbsensiHistory.vue')
      },
      {
        path: 'jurnal',
        name: 'SiswaJurnal',
        component: () => import('../views/siswa/JurnalView.vue')
      },
      {
        path: 'jurnal/tulis',
        name: 'SiswaJurnalTulis',
        component: () => import('../views/siswa/JurnalForm.vue')
      },
      {
        path: 'nilai',
        name: 'SiswaNilai',
        component: () => import('../views/siswa/NilaiView.vue')
      },
      {
        path: 'profile',
        name: 'SiswaProfile',
        component: () => import('../views/siswa/ProfileView.vue')
      }
    ]
  },

  // DUDI Routes (PWA Layout)
  {
    path: '/dudi',
    component: () => import('../layouts/PwaLayout.vue'),
    meta: { requiresAuth: true, role: 'dudi' },
    children: [
      {
        path: '',
        name: 'DudiDashboard',
        component: () => import('../views/dudi/DudiDashboard.vue')
      },
      {
        path: 'siswa',
        name: 'DudiSiswa',
        component: () => import('../views/dudi/StudentList.vue')
      },
      {
        path: 'penilaian/:studentId?',
        name: 'DudiPenilaian',
        component: () => import('../views/dudi/PenilaianForm.vue')
      },
      {
        path: 'jurnal',
        name: 'DudiJurnal',
        component: () => import('../views/dudi/JurnalReview.vue')
      },
      {
        path: 'profile',
        name: 'DudiProfile',
        component: () => import('../views/siswa/ProfileView.vue')
      }
    ]
  },

  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  authStore.restoreSession()

  if (to.meta.guest && authStore.isAuthenticated) {
    return redirectByRole(authStore.userRole, next)
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next({ name: 'Login' })
  }

  if (to.meta.role && authStore.userRole !== to.meta.role) {
    return redirectByRole(authStore.userRole, next)
  }

  next()
})

function redirectByRole(role, next) {
  const map = {
    admin: { name: 'AdminDashboard' },
    admin_jurusan: { name: 'JurusanDashboard' },
    teacher: { name: 'GuruDashboard' },
    student: { name: 'SiswaDashboard' },
    dudi: { name: 'DudiDashboard' }
  }
  const route = map[role]
  if (route) next(route)
  else next({ name: 'Login' })
}

export default router
