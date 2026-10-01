<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/auth.service'

const router = useRouter()

const form = ref({
  full_name: '',
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const success = ref(false)

const handleRegister = async () => {
  error.value = ''
  loading.value = true

  try {
    const response = await authService.register(form.value)
    
    if (response.data.success) {
      success.value = true
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Register gagal'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg p-8">
        <h1 class="text-2xl font-bold text-center mb-6">Register</h1>
        
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div v-if="success" class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-green-700 dark:text-green-400 px-4 py-3 rounded">
            Registrasi berhasil! Redirect ke login...
          </div>

          <div v-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Nama Lengkap</label>
            <input
              v-model="form.full_name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
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
              minlength="6"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>

          <button
            type="submit"
            :disabled="loading || success"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-medium py-2 px-4 rounded-md transition"
          >
            {{ loading ? 'Loading...' : 'Register' }}
          </button>
        </form>

        <p class="text-center mt-4 text-sm">
          Sudah punya akun?
          <router-link to="/login" class="text-blue-600 hover:underline">Login</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
