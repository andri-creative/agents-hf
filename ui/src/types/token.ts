// src/types/token.ts
export interface APIToken {
  id: number
  user_id: number
  token: string
  name: string
  is_active: boolean
  last_used_at: string | null
  expires_at: string | null
  created_at: string
  updated_at: string
}

export interface CreateTokenRequest {
  name: string
  expires_at?: string
}

export interface UpdateTokenStatusRequest {
  is_active: boolean
}

export interface TokenListResponse {
  success: boolean
  data: APIToken[]
}

export interface TokenDetailResponse {
  success: boolean
  data: APIToken
}

export interface CreateTokenResponse {
  success: boolean
  message: string
  data: APIToken
}
