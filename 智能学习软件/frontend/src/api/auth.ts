// 认证 API — Phase B

import { http } from './client'
import type { AuthResponse } from '@/types/domain'

export interface RegisterPayload {
  name: string
  phone?: string
  email?: string
  password: string
  role?: string
}

export interface LoginPayload {
  account: string
  password: string
}

export interface RefreshPayload {
  refresh_token: string
}

export const authApi = {
  register(body: RegisterPayload) {
    return http.post<{ data: AuthResponse }>('/auth/register', body)
  },
  login(body: LoginPayload) {
    return http.post<{ data: AuthResponse }>('/auth/login', body)
  },
  refresh(body: RefreshPayload) {
    return http.post<{ data: AuthResponse }>('/auth/refresh', body)
  }
}