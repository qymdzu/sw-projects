// 与后端一致的统一响应格式

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  detail?: unknown
}

export interface PaginationData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// 业务错误码（与 pkg/response/response.go 一致）
export const BusinessCode = {
  Success: 0,
  BadRequest: 10001,
  ResourceExists: 10002,
  ResourceState: 10003,
  Unauthorized: 20001,
  TokenExpired: 20002,
  TokenInvalid: 20003,
  Forbidden: 20004,
  NotFound: 30001,
  Conflict: 30002,
  RateLimited: 40001,
  ServerError: 50001,
  UpstreamError: 50002
} as const

export type BusinessCodeValue = (typeof BusinessCode)[keyof typeof BusinessCode]