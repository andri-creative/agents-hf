<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/auth.service'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    const response = await authService.login(form.value)
    
    if (response.data.success) {
      authStore.setAuth({
        token: response.data.data.jwt_token,
        apiToken: response.data.data.api_token,
        user: response.data.data.user,
      })
      router.push('/dashboard')
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Login gagal'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-8">
        <h1 class="text-2xl font-bold text-center mb-6">Login</h1>
        
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div v-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Username</label>
            <input
              v-model="form.username"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Password</label>
            <input
              v-model="form.password"
              type="password"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-medium py-2 px-4 rounded-md transition"
          >
            {{ loading ? 'Loading...' : 'Login' }}
          </button>
        </form>

        <p class="text-center mt-4 text-sm">
          Belum punya akun?
          <router-link to="/register" class="text-blue-600 hover:underline">Register</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
