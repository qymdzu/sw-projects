// axios 客户端单元测试 — Phase B
//
// 验证：
//   - 业务 code != 0 → reject ApiBusinessError
//   - 5xx → reject（自动 toast）
//   - 业务 code == 0 → resolve

import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock pinia store
vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    refresh: vi.fn().mockResolvedValue(null),
    logout: vi.fn()
  })
}))

// Mock element-plus toast
vi.mock('element-plus', () => ({
  ElMessage: { error: vi.fn(), warning: vi.fn(), success: vi.fn() }
}))

import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import { http, ApiBusinessError } from '@/api/client'

let mock: MockAdapter

describe('http client', () => {
  beforeEach(() => {
    mock = new MockAdapter(http)
  })

  it('业务 code == 0 → resolve', async () => {
    mock.onGet('/test/ok').reply(200, { code: 0, message: 'success', data: { x: 1 } })

    const { data } = await http.get('/test/ok')
    expect(data.data).toEqual({ x: 1 })
  })

  it('业务 code != 0 → reject ApiBusinessError', async () => {
    mock.onGet('/test/bad').reply(200, { code: 10001, message: '参数错误', data: null })

    await expect(http.get('/test/bad')).rejects.toThrowError(ApiBusinessError)
    await expect(http.get('/test/bad')).rejects.toMatchObject({
      code: 10001,
      message: '参数错误'
    })
  })

  it('5xx → reject', async () => {
    mock.onGet('/test/500').reply(500, { code: 50001, message: 'server error' })

    await expect(http.get('/test/500')).rejects.toThrow()
  })

  it('timeout → reject', async () => {
    mock.onGet('/test/timeout').timeout()

    await expect(http.get('/test/timeout')).rejects.toThrow()
  })
})