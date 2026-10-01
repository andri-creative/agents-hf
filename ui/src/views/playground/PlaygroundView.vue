<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { aiService } from '@/services/ai.service'
import { modelService } from '@/services/model.service'
import { Send, Copy, Loader2 } from 'lucide-vue-next'

const models = ref<any[]>([])
const selectedModel = ref('')
const prompt = ref('')
const result = ref('')
const loading = ref(false)
const loadingModels = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const response = await modelService.list()
    if (response.data.success) {
      models.value = response.data.data.filter((m: any) => m.is_active)
      if (models.value.length > 0) {
        selectedModel.value = models.value[0].slug
      }
    }
  } catch (err) {
    console.error('Failed to load models', err)
  } finally {
    loadingModels.value = false
  }
})

const handleGenerate = async () => {
  if (!prompt.value.trim()) return

  error.value = ''
  result.value = ''
  loading.value = true

  try {
    const response = await aiService.generate({
      model: selectedModel.value,
      prompt: prompt.value,
    })

    if (response.data.success) {
      result.value = response.data.data.text
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Generate gagal'
  } finally {
    loading.value = false
  }
}

const copyToClipboard = () => {
  navigator.clipboard.writeText(result.value)
}
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">AI Playground</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input -->
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <h2 class="text-xl font-bold mb-4">Input</h2>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Model</label>
            <select
              v-model="selectedModel"
              :disabled="loadingModels"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            >
              <option v-if="loadingModels" value="">Loading models...</option>
              <option v-for="model in models" :key="model.id" :value="model.slug">
                {{ model.name }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Prompt</label>
            <textarea
              v-model="prompt"
              rows="10"
              placeholder="Masukkan prompt Anda di sini..."
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700"
            ></textarea>
          </div>

          <button
            @click="handleGenerate"
            :disabled="loading || !prompt.trim() || !selectedModel"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-medium py-2 px-4 rounded-md transition flex items-center justify-center"
          >
            <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
            <Send v-else class="w-4 h-4 mr-2" />
            {{ loading ? 'Generating...' : 'Generate' }}
          </button>
        </div>
      </div>

      <!-- Output -->
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 shadow">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-xl font-bold">Output</h2>
          <button
            v-if="result"
            @click="copyToClipboard"
            class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md transition"
          >
            <Copy class="w-4 h-4" />
          </button>
        </div>

        <div v-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded mb-4">
          {{ error }}
        </div>

        <div
          v-if="!result && !loading && !error"
          class="h-64 flex items-center justify-center text-gray-400 dark:text-gray-600"
        >
          Hasil akan muncul di sini
        </div>

        <div
          v-else-if="loading"
          class="h-64 flex items-center justify-center"
        >
          <Loader2 class="w-8 h-8 animate-spin text-blue-500" />
        </div>

        <div v-else class="prose dark:prose-invert max-w-none">
          <pre class="whitespace-pre-wrap bg-gray-50 dark:bg-gray-900 p-4 rounded-md">{{ result }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>
