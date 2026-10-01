<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { modelService } from '@/services/model.service'
import { Plus, Edit, Trash2, Power } from 'lucide-vue-next'

const models = ref<any[]>([])
const loading = ref(true)
const showDialog = ref(false)
const editMode = ref(false)

const form = ref({
  id: 0,
  name: '',
  slug: '',
  description: '',
  is_active: true,
})

const loadModels = async () => {
  loading.value = true
  try {
    const response = await modelService.list()
    if (response.data.success) {
      models.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load models', error)
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  form.value = { id: 0, name: '', slug: '', description: '', is_active: true }
  editMode.value = false
  showDialog.value = true
}

const openEditDialog = (model: any) => {
  form.value = { ...model }
  editMode.value = true
  showDialog.value = true
}

const handleSubmit = async () => {
  try {
    if (editMode.value) {
      await modelService.update(form.value.id, {
        name: form.value.name,
        slug: form.value.slug,
        description: form.value.description,
        is_active: form.value.is_active,
      })
    } else {
      await modelService.create({
        name: form.value.name,
        slug: form.value.slug,
        description: form.value.description,
        is_active: form.value.is_active,
      })
    }
    showDialog.value = false
    await loadModels()
  } catch (error) {
    console.error('Failed to save model', error)
  }
}

const deleteModel = async (id: number) => {
  if (!confirm('Yakin ingin menghapus model ini?')) return

  try {
    await modelService.remove(id)
    await loadModels()
  } catch (error) {
    console.error('Failed to delete model', error)
  }
}

const toggleActive = async (model: any) => {
  try {
    await modelService.update(model.id, { is_active: !model.is_active })
    await loadModels()
  } catch (error) {
    console.error('Failed to toggle active', error)
  }
}

onMounted(loadModels)
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-3xl font-bold">AI Models</h1>
      <button
        @click="openCreateDialog"
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition flex items-center"
      >
        <Plus class="w-4 h-4 mr-2" />
        Tambah Model
      </button>
    </div>

    <!-- Dialog -->
    <div v-if="showDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">{{ editMode ? 'Edit Model' : 'Tambah Model' }}</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Name</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Slug</label>
            <input
              v-model="form.slug"
              type="text"
              required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Description</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            ></textarea>
          </div>
          <div class="flex items-center">
            <input
              v-model="form.is_active"
              type="checkbox"
              id="is_active"
              class="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
            />
            <label for="is_active" class="ml-2 text-sm font-medium">Active</label>
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

      <div v-else-if="models.length === 0" class="p-6 text-center text-gray-500">
        Belum ada model. Tambah model pertama!
      </div>

      <table v-else class="w-full">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Slug</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="model in models" :key="model.id">
            <td class="px-6 py-4 font-medium">{{ model.name }}</td>
            <td class="px-6 py-4 text-sm">{{ model.slug }}</td>
            <td class="px-6 py-4 text-sm">{{ model.description }}</td>
            <td class="px-6 py-4">
              <span
                :class="model.is_active ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'"
                class="px-2 py-1 rounded-full text-xs font-medium"
              >
                {{ model.is_active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <button @click="openEditDialog(model)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                  <Edit class="w-4 h-4" />
                </button>
                <button @click="toggleActive(model)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded">
                  <Power class="w-4 h-4" :class="model.is_active ? 'text-green-600' : 'text-gray-400'" />
                </button>
                <button @click="deleteModel(model.id)" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded text-red-600">
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
