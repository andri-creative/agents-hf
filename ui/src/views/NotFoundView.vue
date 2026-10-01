<script setup lang="ts">
import { Home } from 'lucide-vue-next'
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-gray-300 dark:text-gray-700 mb-4">404</h1>
      <h2 class="text-2xl font-semibold mb-4">Halaman Tidak Ditemukan</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-6">
        Maaf, halaman yang Anda cari tidak dapat ditemukan.
      </p>
      <router-link
        to="/dashboard"
        class="inline-flex items-center bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-6 rounded-md transition"
      >
        <Home class="w-4 h-4 mr-2" />
        Kembali ke Dashboard
      </router-link>
    </div>
  </div>
</template>
