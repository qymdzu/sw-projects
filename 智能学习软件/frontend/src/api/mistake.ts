// 错题本 API — Phase B

import { http } from './client'
import type { Mistake } from '@/types/domain'

export const mistakeApi = {
  list(params?: { knowledge_point_id?: number; mastered?: boolean; page?: number; page_size?: number }) {
    return http.get<{ data: { items: Mistake[]; total: number } }>('/mistakes', { params })
  },
  groups() {
    return http.get<{ data: unknown[] }>('/mistakes/groups')
  },
  review(body: { mistake_ids: number[] }) {
    return http.post<{ code: number }>('/mistakes/review', body)
  },
  markMastered(id: number) {
    return http.put<{ code: number }>(`/mistakes/${id}/master`)
  },
  remove(id: number) {
    return http.delete<{ code: number }>(`/mistakes/${id}`)
  }
}