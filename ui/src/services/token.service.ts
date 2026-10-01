// src/services/token.service.ts
import api from './api'
import type { CreateTokenRequest, UpdateTokenStatusRequest, TokenListResponse, TokenDetailResponse, CreateTokenResponse } from '@/types/token'

export const tokenService = {
  list: () => api.get<TokenListResponse>('/api/tokens'),
  create: (data: CreateTokenRequest) => api.post<CreateTokenResponse>('/api/tokens', data),
  detail: (id: number) => api.get<TokenDetailResponse>(`/api/tokens/${id}`),
  toggleStatus: (id: number, isActive: boolean) =>
    api.put(`/api/tokens/${id}/status`, { is_active: isActive } as UpdateTokenStatusRequest),
  remove: (id: number) => api.delete(`/api/tokens/${id}`),
}
