import api from './index'

export const getLogFiles = () => api.get('/logs/files')

export const tailLog = (file?: string, lines?: number) =>
  api.get('/logs/tail', { params: { file, lines } })

export const searchLogs = (keyword: string, file?: string) =>
  api.get('/logs/search', { params: { keyword, file } })

export const getArchives = () => api.get('/logs/archive')
