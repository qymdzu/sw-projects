import api from './index'

export const getTree = (repo?: string, path?: string) =>
  api.get('/storage/tree', { params: { repo, path } })

export const getFile = (path: string) =>
  api.get('/storage/file', { params: { path } })

export const saveFile = (path: string, content: string, create_backup = true) =>
  api.put('/storage/file', { path, content, create_backup })
