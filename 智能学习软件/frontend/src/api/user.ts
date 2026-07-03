// 用户 API — Phase B

import { http } from './client'
import type { User } from '@/types/domain'

export const userApi = {
  me() {
    return http.get<{ data: User }>('/users/me')
  },
  updateMe(body: Partial<Pick<User, 'name' | 'avatar_url'>>) {
    return http.put<{ data: User }>('/users/me', body)
  },
  changePassword(body: { old_password: string; new_password: string }) {
    return http.put<{ code: number }>('/users/me/password', body)
  },
  updateAvatar(body: { avatar_url: string }) {
    return http.post<{ data: User }>('/users/me/avatar', body)
  }
}