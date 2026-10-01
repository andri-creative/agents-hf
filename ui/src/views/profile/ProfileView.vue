<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/auth.service'

const authStore = useAuthStore()
const user = ref<any>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await authService.me()
    if (response.data.success) {
      user.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load profile', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Profile</h1>

    <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow animate-pulse">
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-32 mb-4"></div>
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-48 mb-2"></div>
      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-40"></div>
    </div>

    <div v-else class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">ID</label>
          <p class="text-lg">{{ user?.id }}</p>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Nama Lengkap</label>
          <p class="text-lg">{{ user?.full_name }}</p>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Username</label>
          <p class="text-lg">{{ user?.username }}</p>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Role</label>
          <span
            :class="user?.role === 'admin' ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400' : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'"
            class="px-3 py-1 rounded-full text-sm font-medium"
          >
            {{ user?.role }}
          </span>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Created At</label>
          <p class="text-lg">{{ user?.created_at ? new Date(user.created_at).toLocaleString() : '-' }}</p>
        </div>
      </div>
    </div>

    <div class="mt-6 bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
      <h2 class="text-xl font-bold mb-4">Quick Links</h2>
      <div class="space-y-2">
        <router-link
          to="/tokens"
          class="block p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        >
          Kelola API Tokens
        </router-link>
        <router-link
          to="/playground"
          class="block p-3 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        >
          Test AI Playground
        </router-link>
      </div>
    </div>
  </div>
</template>
