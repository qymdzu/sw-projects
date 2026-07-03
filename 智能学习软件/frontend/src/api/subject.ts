// 科目 + 知识点 API — Phase B

import { http } from './client'
import type { Subject, KnowledgePoint } from '@/types/domain'

export const subjectApi = {
  list(params?: { q?: string }) {
    return http.get<{ data: { items: Subject[]; total: number } }>('/subjects', { params })
  },
  tree(params?: { subject_id?: number }) {
    return http.get<{ data: KnowledgePoint[] }>('/knowledge-points', { params })
  }
}