// auth store 单元测试 — Phase B
//
// 验证：login 成功路径 → 写 token 到 storage；refresh 失败返回 null；logout 清空。

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { storage } from '@/utils/storage'

// Mock axios 客户端
vi.mock('@/api/client', () => ({
  http: {
    get: vi.fn(),
    post: vi.fn()
  }
}))

vi.mock('@/api/auth', () => ({
  authApi: {
    login: vi.fn(),
    register: vi.fn(),
    refresh: vi.fn()
  }
}))

import { authApi } from '@/api/auth'
import { http } from '@/api/client'

const mockedLogin = authApi.login as unknown as ReturnType<typeof vi.fn>
const mockedRegister = authApi.register as unknown as ReturnType<typeof vi.fn>
const mockedRefresh = authApi.refresh as unknown as ReturnType<typeof vi.fn>
const mockedHttpGet = http.get as unknown as ReturnType<typeof vi.fn>

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    storage.remove('access_token')
    storage.remove('refresh_token')
    vi.clearAllMocks()
  })

  it('login 成功 → 写入 storage 并更新 store', async () => {
    mockedLogin.mockResolvedValueOnce({
      data: {
        data: {
          user: { id: 'u1', name: 'Alice', role: 'student' },
          access_token: 'AT-xxx',
          refresh_token: 'RT-yyy',
          expires_in: 7200
        }
      }
    })

    const auth = useAuthStore()
    await auth.login('alice@example.com', 'pwd123')

    expect(auth.accessToken).toBe('AT-xxx')
    expect(auth.refreshToken).toBe('RT-yyy')
    expect(auth.user?.name).toBe('Alice')
    expect(auth.isAuthenticated).toBe(true)
    expect(storage.get('access_token')).toBe('AT-xxx')
    expect(storage.get('refresh_token')).toBe('RT-yyy')
  })

  it('register 成功 → 写 token', async () => {
    mockedRegister.mockResolvedValueOnce({
      data: {
        data: {
          user: { id: 'u2', name: 'Bob', role: 'student' },
          access_token: 'AT-bob',
          refresh_token: 'RT-bob',
          expires_in: 7200
        }
      }
    })

    const auth = useAuthStore()
    await auth.register({ name: 'Bob', email: 'bob@example.com', password: 'pwd123' })

    expect(auth.accessToken).toBe('AT-bob')
    expect(auth.isAuthenticated).toBe(true)
  })

  it('refresh 成功 → 返回新 access_token', async () => {
    mockedRefresh.mockResolvedValueOnce({
      data: {
        data: {
          user: { id: 'u1', name: 'Alice', role: 'student' },
          access_token: 'AT-new',
          refresh_token: 'RT-new',
          expires_in: 7200
        }
      }
    })

    const auth = useAuthStore()
    auth.refreshToken = 'RT-old'

    const token = await auth.refresh()
    expect(token).toBe('AT-new')
    expect(auth.accessToken).toBe('AT-new')
  })

  it('refresh 失败 → 返回 null', async () => {
    mockedRefresh.mockRejectedValueOnce(new Error('expired'))

    const auth = useAuthStore()
    auth.refreshToken = 'RT-expired'

    const token = await auth.refresh()
    expect(token).toBeNull()
  })

  it('logout 清空 store + storage', async () => {
    const auth = useAuthStore()
    auth.accessToken = 'old'
    auth.refreshToken = 'old-rt'
    auth.user = { id: 'u1', name: 'Alice', role: 'student', created_at: '' }

    await auth.logout()

    expect(auth.accessToken).toBe('')
    expect(auth.refreshToken).toBe('')
    expect(auth.user).toBeNull()
    expect(storage.get('access_token')).toBeNull()
    expect(storage.get('refresh_token')).toBeNull()
  })

  it('fetchMe 成功 → 更新 user', async () => {
    mockedHttpGet.mockResolvedValueOnce({
      data: { data: { id: 'u1', name: 'Alice-updated', role: 'student', created_at: '' } }
    })

    const auth = useAuthStore()
    await auth.fetchMe()

    expect(auth.user?.name).toBe('Alice-updated')
  })
})