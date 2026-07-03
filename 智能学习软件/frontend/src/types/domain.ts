// 领域类型 — Phase B

export type Role = 'student' | 'teacher' | 'admin' | 'parent'

export interface User {
  id: string
  name: string
  role: Role
  phone?: string | null
  email?: string | null
  avatar_url?: string | null
  created_at: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface Subject {
  id: number
  name: string
  description?: string
  icon?: string
}

export interface KnowledgePoint {
  id: number
  subject_id: number
  name: string
  parent_id?: number | null
}

export type DifficultyLevel = 1 | 2 | 3 | 4 | 5

export interface Question {
  id: number
  subject_id: number
  knowledge_point_id?: number
  type: 'single' | 'multiple' | 'fill' | 'essay'
  difficulty: DifficultyLevel
  content: string
  options?: string[]
  answer: string
  analysis?: string
}

export interface StudyPlan {
  id: number
  user_id: string
  title: string
  description?: string
  start_date: string
  end_date: string
  status: 'active' | 'completed' | 'paused'
  progress?: number
}

export interface Mistake {
  id: number
  user_id: string
  question_id: number
  question: Question
  knowledge_point: { id: number; name: string }
  wrong_answer: string
  mistake_count: number
  mastered: boolean
  created_at: string
  last_wrong_at: string
}

export interface ReportSummary {
  total_exercises: number
  overall_correct_rate: number
  streak_days: number
  unmastered_mistakes: number
  today_duration_min: number
}

export type Provider = 'openai' | 'anthropic' | 'qwen' | 'deepseek' | 'ollama' | 'custom'

export const ProviderLabels: Record<Provider, string> = {
  openai: 'OpenAI',
  anthropic: 'Anthropic (Claude)',
  qwen: '通义千问 (Qwen)',
  deepseek: 'DeepSeek',
  ollama: 'Ollama (本地)',
  custom: '自定义'
}

export const ProviderEndpoints: Record<Provider, string> = {
  openai: 'https://api.openai.com/v1',
  anthropic: 'https://api.anthropic.com',
  qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
  deepseek: 'https://api.deepseek.com/v1',
  ollama: 'http://localhost:11434',
  custom: ''
}

export const ProviderModels: Record<Provider, string[]> = {
  openai: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-3.5-turbo'],
  anthropic: ['claude-3-5-sonnet-20241022', 'claude-3-opus-20240229', 'claude-3-haiku-20240307'],
  qwen: ['qwen-plus', 'qwen-turbo', 'qwen-max', 'qwen-long'],
  deepseek: ['deepseek-chat', 'deepseek-coder', 'deepseek-reasoner'],
  ollama: ['llama3.1', 'qwen2', 'mistral', 'phi3'],
  custom: []
}

export interface ModelSetting {
  provider: Provider
  api_endpoint: string
  model: string
  extra_config?: Record<string, unknown>
  is_default: boolean
  updated_at: string
}

export interface ModelSettingFull extends ModelSetting {
  api_key: string // 仅 GET 时返回明文
}