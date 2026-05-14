<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 p-4">
    <div class="text-center">
      <div class="w-20 h-20 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-6">
        <FrownIcon :size="40" class="text-gray-400" />
      </div>
      <h1 class="text-3xl font-bold text-gray-800 mb-2">404</h1>
      <p class="text-gray-500 mb-6">Halaman tidak ditemukan</p>
      <router-link
        to="/login"
        class="inline-flex items-center gap-2 px-6 py-2.5 bg-primary text-white rounded-xl font-medium text-sm hover:bg-primary-light transition-colors"
      >
        <ArrowLeftIcon :size="18" />
        Kembali ke Login
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { FrownIcon, ArrowLeftIcon } from 'lucide-vue-next'
</script>
