// src/services/model.service.ts
import api from './api'
import type { CreateModelRequest, UpdateModelRequest, ModelListResponse, ModelDetailResponse } from '@/types/model'

export const modelService = {
  list: () => api.get<ModelListResponse>('/api/admin/models'),
  create: (data: CreateModelRequest) => api.post('/api/admin/models', data),
  detail: (id: number) => api.get<ModelDetailResponse>(`/api/admin/models/${id}`),
  update: (id: number, data: UpdateModelRequest) => api.put(`/api/admin/models/${id}`, data),
  remove: (id: number) => api.delete(`/api/admin/models/${id}`),
}
