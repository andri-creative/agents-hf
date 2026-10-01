// src/types/auth.ts
export interface User {
  id: number
  full_name: string
  username: string
  role: string
  created_at: string
}

export interface RegisterRequest {
  full_name: string
  username: string
  password: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface AuthResponse {
  success: boolean
  message: string
  data: {
    jwt_token: string
    api_token: string
    user: User
  }
}

export interface MeResponse {
  success: boolean
  data: User
}
