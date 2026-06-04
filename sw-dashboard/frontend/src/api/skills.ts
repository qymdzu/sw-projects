import api from './index'

export const getTree = (path?: string) =>
  api.get('/skills/tree', { params: { path } })

export const getFile = (path: string) =>
  api.get('/skills/file', { params: { path } })
