// src/types/ai.ts
export interface GenerateRequest {
  model?: string
  prompt: string
}

export interface GenerateResponse {
  success: boolean
  data: {
    text: string
    model: string
  }
}

export interface ChatMessage {
  role: 'system' | 'user' | 'assistant'
  content: string
}

export interface ChatCompletionsRequest {
  model: string
  messages: ChatMessage[]
  temperature?: number
  max_tokens?: number
  stream?: boolean
}
