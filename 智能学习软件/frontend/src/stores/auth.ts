// Pinia 鉴权 store — Phase B

import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { http } from '@/api/client'
import { authApi } from '@/api/auth'
import { storage } from '@/utils/storage'
import type { User } from '@/types/domain'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string>(storage.get('access_token') || '')
  const refreshToken = ref<string>(storage.get('refresh_token') || '')
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(account: string, password: string): Promise<void> {
    const { data } = await authApi.login({ account, password })
    const body = data.data
    accessToken.value = body.access_token
    refreshToken.value = body.refresh_token
    user.value = body.user
    storage.set('access_token', body.access_token)
    storage.set('refresh_token', body.refresh_token)
  }

  async function register(payload: {
    name: string
    phone?: string
    email?: string
    password: string
  }): Promise<void> {
    const { data } = await authApi.register(payload)
    const body = data.data
    accessToken.value = body.access_token
    refreshToken.value = body.refresh_token
    user.value = body.user
    storage.set('access_token', body.access_token)
    storage.set('refresh_token', body.refresh_token)
  }

  async function fetchMe(): Promise<void> {
    const { data } = await http.get<{ data: User }>('/users/me')
    user.value = data.data
  }

  /**
   * 用 refresh_token 换新 token 对。
   * 成功返回新 access_token；失败返回 null。
   */
  async function refresh(): Promise<string | null> {
    if (!refreshToken.value) return null
    try {
      const { data } = await authApi.refresh({ refresh_token: refreshToken.value })
      const body = data.data
      accessToken.value = body.access_token
      refreshToken.value = body.refresh_token
      user.value = body.user
      storage.set('access_token', body.access_token)
      storage.set('refresh_token', body.refresh_token)
      return body.access_token
    } catch {
      return null
    }
  }

  async function logout(): Promise<void> {
    accessToken.value = ''
    refreshToken.value = ''
    user.value = null
    storage.remove('access_token')
    storage.remove('refresh_token')
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    login,
    register,
    fetchMe,
    refresh,
    logout
  }
})