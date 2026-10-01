<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { configService } from '@/services/config.service'
import { Plus, Edit, Trash2, Eye, EyeOff } from 'lucide-vue-next'

const configs = ref<any[]>([])
const loading = ref(true)
const showDialog = ref(false)
const editMode = ref(false)
const showValues = ref<Record<number, boolean>>({})

const form = ref({
  id: 0,
  key: '',
  value: '',
  description: '',
})

const loadConfigs = async () => {
  loading.value = true
  try {
    const response = await configService.list()
    if (response.data.success) {
      configs.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load configs', error)
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  form.value = { id: 0, key: '', value: '', description: '' }
  editMode.value = false
  showDialog.value = true
}

const openEditDialog = (config: any) => {
  form.value = { ...config }
  editMode.value = true
  showDialog.value = true
}

const handleSubmit = async () => {
  try {
    if (editMode.value) {
      await configService.update(form.value.id, {
        value: form.value.value,
        description: form.value.description,
      })
    } else {
      await configService.create({
        key: form.value.key,
        value: form.value.value,
        description: form.value.description,
      })
    }
    showDialog.value = false
    await loadConfigs()
  } catch (error) {
    console.error('Failed to save config', error)
  }
}

const deleteConfig = async (id: number) => {
  if (!confirm('Yakin ingin menghapus config ini?')) return

  try {
    await configService.remove(id)
    await loadConfigs()
  } catch (error) {
    console.error('Failed to delete config', error)
  }
}

const toggleShowValue = (id: number) => {
  showValues.value[id] = !showValues.value[id]
}

const maskValue = (value: string) => {
  return '•'.repeat(Math.min(value.length, 20))
}

const isSensitive = (key: string) => {
  return key.includes('KEY') || key.includes('SECRET') || key.includes('PASSWORD')
}

onMounted(loadConfigs)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold">Configs</h1>
      <button
        @click="openCreateDialog"
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition flex items-center"
      >
        <Plus class="w-4 h-4 mr-2" />
        Tambah Config
      </button>
    </div>

    <!-- Dialog -->
    <div v-if="showDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">{{ editMode ? 'Edit Config' : 'Tambah Config' }}</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Key</label>
            <input
              v-model="form.key"
              type="text"
              required
              :disabled="editMode"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 disabled:opacity-50"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Value</label>
            <textarea
              v-model="form.value"
              rows="3"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            ></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Description</label>
            <input
              v-model="form.description"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>
          <div class="flex gap-2">
            <button type="submit" class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md">
              {{ editMode ? 'Update' : 'Buat' }}
            </button>
            <button type="button" @click="showDialog = false" class="flex-1 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 font-medium py-2 px-4 rounded-md">
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
      <div v-if="loading" class="p-6">
        <div class="animate-pulse space-y-3">
          <div v-for="i in 3" :key="i" class="h-12 bg-gray-200 dark:bg-gray-700 rounded"></div>
        </div>
      </div>

      <div v-else-if="configs.length === 0" class="p-6 text-center text-gray-500">
        Belum ada config. Tambah config pertama!
      </div>

      <table v-else class="w-full">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Key</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Value</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="config in configs" :key="config.id">
            <td class="px-6 py-4 font-medium">{{ config.key }}</td>
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <code class="text-sm">
                  {{ isSensitive(config.key) && !showValues[config.id] ? maskValue(config.value) : config.value }}
                </code>
                <button
                  v-if="isSensitive(config.key)"
                  @click="toggleShowValue(config.id)"
                  class="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
                >
                  <Eye v-if="!showValues[config.id]" class="w-4 h-4" />
                  <EyeOff v-else class="w-4 h-4" />
                </button>
              </div>
            </td>
            <td class="px-6 py-4 text-sm">{{ config.description }}</td>
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <button @click="openEditDialog(config)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                  <Edit class="w-4 h-4" />
                </button>
                <button @click="deleteConfig(config.id)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded text-red-600">
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
