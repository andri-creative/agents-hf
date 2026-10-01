// src/types/config.ts
export interface Config {
  id: number
  key: string
  value: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateConfigRequest {
  key: string
  value: string
  description?: string
}

export interface UpdateConfigRequest {
  value?: string
  description?: string
}

export interface ConfigListResponse {
  success: boolean
  data: Config[]
}

export interface ConfigDetailResponse {
  success: boolean
  data: Config
}
