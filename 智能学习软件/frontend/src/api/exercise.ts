// 练习 API — Phase B

import { http } from './client'
import type { Question, DifficultyLevel } from '@/types/domain'

export interface SubmitResult {
  correct: boolean
  correct_answer: string
  analysis?: string
  knowledge_point_id?: number
}

export const exerciseApi = {
  list(params: { subject_id?: number; difficulty?: DifficultyLevel; page?: number; page_size?: number }) {
    return http.get<{ data: { items: Question[]; total: number } }>('/exercises', { params })
  },
  random(params?: { subject_id?: number; difficulty?: DifficultyLevel }) {
    return http.get<{ data: Question }>('/exercises/random', { params })
  },
  submit(body: { question_id: number; answer: string }) {
    return http.post<{ data: SubmitResult }>('/exercises/submit', body)
  },
  recommend(limit = 10) {
    return http.get<{ data: { items: Question[]; total: number } }>('/exercises/recommend', { params: { limit } })
  },
  byKnowledgePoint(kpId: number) {
    return http.get<{ data: { items: Question[] } }>(`/exercises/knowledge-points/${kpId}`)
  },
  history(params?: { page?: number; page_size?: number }) {
    return http.get<{ data: { items: unknown[]; total: number } }>('/exercises/history', { params })
  }
}