// src/services/config.service.ts
import api from './api'
import type { CreateConfigRequest, UpdateConfigRequest, ConfigListResponse, ConfigDetailResponse } from '@/types/config'

export const configService = {
  list: () => api.get<ConfigListResponse>('/api/admin/configs'),
  create: (data: CreateConfigRequest) => api.post('/api/admin/configs', data),
  detail: (id: number) => api.get<ConfigDetailResponse>(`/api/admin/configs/${id}`),
  update: (id: number, data: UpdateConfigRequest) => api.put(`/api/admin/configs/${id}`, data),
  remove: (id: number) => api.delete(`/api/admin/configs/${id}`),
}
