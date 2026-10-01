// src/services/ai.service.ts
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import type { GenerateRequest, GenerateResponse } from '@/types/ai'

export const aiService = {
  generate: async (data: GenerateRequest) => {
    const auth = useAuthStore()
    return axios.post<GenerateResponse>(
      `${import.meta.env.VITE_API_BASE_URL}/v1/messages`,
      data,
      { headers: { 'X-API-Key': auth.apiToken } }
    )
  },
}
