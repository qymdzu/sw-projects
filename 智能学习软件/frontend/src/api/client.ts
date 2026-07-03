// Axios 实例 — Phase B
//
// 功能：
//   - 自动注入 Authorization Bearer Token
//   - 401 时自动尝试 refresh（去重并发）
//   - 5xx 时 Toast 提示
//   - 业务 code != 0 时统一 reject

import axios, { AxiosError, type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { storage } from '@/utils/storage'
import type { ApiResponse } from '@/types/api'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

export const http: AxiosInstance = axios.create({
  baseURL,
  timeout: 15_000,
  withCredentials: false
})

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = storage.get('access_token')
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})

// 并发 refresh 去重：多个 401 共用同一个 promise
let refreshing: Promise<string | null> | null = null

http.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse
    if (body && typeof body === 'object' && 'code' in body && body.code !== 0) {
      // 业务错误：抛错给调用方处理
      return Promise.reject(new ApiBusinessError(body.code, body.message, body.detail))
    }
    return response
  },
  async (err: AxiosError<ApiResponse>) => {
    const status = err.response?.status
    const original = err.config as InternalAxiosRequestConfig & { _retry?: boolean }
    const auth = useAuthStore()

    if (status === 401 && !original._retry && storage.get('refresh_token')) {
      original._retry = true
      refreshing ??= auth.refresh().finally(() => {
        refreshing = null
      })
      const newToken = await refreshing
      if (newToken) {
        original.headers?.set('Authorization', `Bearer ${newToken}`)
        return http.request(original)
      }
    }

    // refresh 失败 → 强制登出
    if (status === 401) {
      await auth.logout()
      ElMessage.error('登录已过期，请重新登录')
    } else if (status && status >= 500) {
      ElMessage.error('服务器开小差，请稍后再试')
    }

    return Promise.reject(err)
  }
)

/**
 * ApiBusinessError 是业务错误（HTTP 200 但 code != 0）。
 * 携带 code/message/detail 便于上层精细处理。
 */
export class ApiBusinessError extends Error {
  readonly code: number
  readonly detail: unknown
  constructor(code: number, message: string, detail?: unknown) {
    super(message)
    this.name = 'ApiBusinessError'
    this.code = code
    this.detail = detail
  }
}