// src/services/auth.service.ts
import api from './api'
import type { LoginRequest, RegisterRequest, AuthResponse, MeResponse } from '@/types/auth'

export const authService = {
  register: (data: RegisterRequest) => api.post<AuthResponse>('/api/register', data),
  login: (data: LoginRequest) => api.post<AuthResponse>('/api/login', data),
  me: () => api.get<MeResponse>('/api/me'),
}
