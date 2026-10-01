<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { tokenService } from '@/services/token.service'
import { Activity, Key, Bot, Zap } from 'lucide-vue-next'

const authStore = useAuthStore()
const stats = ref({
  totalTokens: 0,
  activeTokens: 0,
  totalModels: 0,
  requestsToday: 0,
})

const loading = ref(true)

onMounted(async () => {
  try {
    const response = await tokenService.list()
    if (response.data.success) {
      const tokens = response.data.data
      stats.value.totalTokens = tokens.length
      stats.value.activeTokens = tokens.filter((t: any) => t.is_active).length
    }
  } catch (error) {
    console.error('Failed to load stats', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Dashboard</h1>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="i in 4" :key="i" class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow animate-pulse">
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24 mb-2"></div>
        <div class="h-8 bg-gray-200 dark:bg-gray-700 rounded w-16"></div>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">Total Token</p>
            <p class="text-3xl font-bold mt-1">{{ stats.totalTokens }}</p>
          </div>
          <Key class="w-10 h-10 text-blue-500" />
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">Token Aktif</p>
            <p class="text-3xl font-bold mt-1">{{ stats.activeTokens }}</p>
          </div>
          <Zap class="w-10 h-10 text-green-500" />
        </div>
      </div>

      <div v-if="authStore.isAdmin" class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">Total Model</p>
            <p class="text-3xl font-bold mt-1">{{ stats.totalModels }}</p>
          </div>
          <Bot class="w-10 h-10 text-purple-500" />
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">Request Hari Ini</p>
            <p class="text-3xl font-bold mt-1">{{ stats.requestsToday }}</p>
          </div>
          <Activity class="w-10 h-10 text-orange-500" />
        </div>
      </div>
    </div>

    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
      <h2 class="text-xl font-bold mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <router-link
          to="/playground"
          class="p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        >
          <h3 class="font-semibold mb-1">AI Playground</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">Test AI generation</p>
        </router-link>

        <router-link
          to="/tokens"
          class="p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        >
          <h3 class="font-semibold mb-1">Manage Tokens</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">Create & manage API tokens</p>
        </router-link>

        <router-link
          v-if="authStore.isAdmin"
          to="/admin/models"
          class="p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition"
        >
          <h3 class="font-semibold mb-1">Manage Models</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">Add & configure AI models</p>
        </router-link>
      </div>
    </div>

    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
      <h2 class="text-xl font-bold mb-4">Welcome, {{ authStore.user?.full_name }}</h2>
      <p class="text-gray-600 dark:text-gray-400">
        Selamat datang di AI Generate Dashboard. Gunakan menu di samping untuk mengakses fitur yang tersedia.
      </p>
    </div>
  </div>
</template>
