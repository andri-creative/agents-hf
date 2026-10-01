<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tokenService } from '@/services/token.service'
import { Plus, Copy, Trash2, Power, Eye, EyeOff } from 'lucide-vue-next'

const tokens = ref<any[]>([])
const loading = ref(true)
const showCreateDialog = ref(false)
const newTokenResult = ref<any>(null)

const form = ref({
  name: '',
  expires_at: '',
})

const loadTokens = async () => {
  loading.value = true
  try {
    const response = await tokenService.list()
    if (response.data.success) {
      tokens.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load tokens', error)
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    const payload: any = { name: form.value.name }
    if (form.value.expires_at) {
      payload.expires_at = new Date(form.value.expires_at).toISOString()
    }

    const response = await tokenService.create(payload)
    if (response.data.success) {
      newTokenResult.value = response.data.data
      form.value.name = ''
      form.value.expires_at = ''
      showCreateDialog.value = false
      await loadTokens()
    }
  } catch (error) {
    console.error('Failed to create token', error)
  }
}

const toggleStatus = async (token: any) => {
  try {
    await tokenService.toggleStatus(token.id, !token.is_active)
    await loadTokens()
  } catch (error) {
    console.error('Failed to toggle status', error)
  }
}

const deleteToken = async (id: number) => {
  if (!confirm('Yakin ingin menghapus token ini?')) return

  try {
    await tokenService.remove(id)
    await loadTokens()
  } catch (error) {
    console.error('Failed to delete token', error)
  }
}

const copyToken = (token: string) => {
  navigator.clipboard.writeText(token)
}

const maskToken = (token: string) => {
  return token.substring(0, 10) + '...' + token.substring(token.length - 8)
}

onMounted(loadTokens)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold">API Tokens</h1>
      <button
        @click="showCreateDialog = true"
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition flex items-center"
      >
        <Plus class="w-4 h-4 mr-2" />
        Buat Token Baru
      </button>
    </div>

    <!-- New Token Result -->
    <div v-if="newTokenResult" class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4 mb-6">
      <h3 class="font-semibold text-green-800 dark:text-green-400 mb-2">Token berhasil dibuat!</h3>
      <p class="text-sm text-green-700 dark:text-green-300 mb-2">Salin token ini sekarang, tidak akan ditampilkan lagi:</p>
      <div class="flex items-center gap-2">
        <code class="flex-1 bg-white dark:bg-gray-800 px-3 py-2 rounded text-sm">{{ newTokenResult.token }}</code>
        <button @click="copyToken(newTokenResult.token)" class="p-2 hover:bg-green-100 dark:hover:bg-green-800 rounded">
          <Copy class="w-4 h-4" />
        </button>
      </div>
      <button @click="newTokenResult = null" class="mt-2 text-sm text-green-700 dark:text-green-300 hover:underline">
        Tutup
      </button>
    </div>

    <!-- Create Dialog -->
    <div v-if="showCreateDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Buat Token Baru</h2>
        <form @submit.prevent="handleCreate" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Nama Token</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Expires At (Opsional)</label>
            <input
              v-model="form.expires_at"
              type="datetime-local"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>
          <div class="flex gap-2">
            <button type="submit" class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md">
              Buat
            </button>
            <button type="button" @click="showCreateDialog = false" class="flex-1 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 font-medium py-2 px-4 rounded-md">
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Tokens Table -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
      <div v-if="loading" class="p-6">
        <div class="animate-pulse space-y-3">
          <div v-for="i in 3" :key="i" class="h-12 bg-gray-200 dark:bg-gray-700 rounded"></div>
        </div>
      </div>

      <div v-else-if="tokens.length === 0" class="p-6 text-center text-gray-500">
        Belum ada token. Buat token pertama Anda!
      </div>

      <table v-else class="w-full">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Token</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Last Used</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="token in tokens" :key="token.id">
            <td class="px-6 py-4">
              <code class="text-sm">{{ maskToken(token.token) }}</code>
            </td>
            <td class="px-6 py-4">{{ token.name }}</td>
            <td class="px-6 py-4">
              <span
                :class="token.is_active ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'"
                class="px-2 py-1 rounded-full text-xs font-medium"
              >
                {{ token.is_active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm">
              {{ token.last_used_at ? new Date(token.last_used_at).toLocaleString() : 'Never' }}
            </td>
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <button @click="copyToken(token.token)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                  <Copy class="w-4 h-4" />
                </button>
                <button @click="toggleStatus(token)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                  <Power class="w-4 h-4" :class="token.is_active ? 'text-green-600' : 'text-gray-400'" />
                </button>
                <button @click="deleteToken(token.id)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded text-red-600">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
