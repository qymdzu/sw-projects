// 学习计划 API — Phase B

import { http } from './client'
import type { StudyPlan } from '@/types/domain'

export const planApi = {
  list(params?: { status?: string; page?: number; page_size?: number }) {
    return http.get<{ data: { items: StudyPlan[]; total: number } }>('/plans', { params })
  },
  create(body: Partial<StudyPlan>) {
    return http.post<{ data: StudyPlan }>('/plans', body)
  },
  aiGenerate(body: { subject_ids: number[]; days: number }) {
    return http.post<{ data: StudyPlan }>('/plans/ai-generate', body)
  },
  getByID(id: number) {
    return http.get<{ data: StudyPlan }>(`/plans/${id}`)
  },
  update(id: number, body: Partial<StudyPlan>) {
    return http.put<{ data: StudyPlan }>(`/plans/${id}`, body)
  },
  remove(id: number) {
    return http.delete<{ code: number }>(`/plans/${id}`)
  },
  checkin(id: number) {
    return http.post<{ code: number }>(`/plans/${id}/checkin`)
  }
}