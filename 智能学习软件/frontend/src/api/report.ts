// 看板/报告 API — Phase B

import { http } from './client'
import type { ReportSummary } from '@/types/domain'

export const reportApi = {
  summary() {
    return http.get<{ data: ReportSummary }>('/reports/summary')
  },
  detail(params?: { days?: number }) {
    return http.get<{ data: unknown }>('/reports/detail', { params })
  },
  mastery() {
    return http.get<{ data: unknown }>('/reports/mastery')
  },
  trend(params?: { days?: number }) {
    return http.get<{ data: unknown }>('/reports/trend', { params })
  }
}