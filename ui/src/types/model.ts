// src/types/model.ts
export interface AIModel {
  id: number
  name: string
  slug: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateModelRequest {
  name: string
  slug: string
  description?: string
  is_active?: boolean
}

export interface UpdateModelRequest {
  name?: string
  slug?: string
  description?: string
  is_active?: boolean
}

export interface ModelListResponse {
  success: boolean
  data: AIModel[]
}

export interface ModelDetailResponse {
  success: boolean
  data: AIModel
}
