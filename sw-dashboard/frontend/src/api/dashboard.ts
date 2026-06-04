import api from './index'

export const getDashboard = () => api.get('/dashboard')
export const getHealth = () => api.get('/health')
export const verifyToken = (token: string) => api.post('/auth', { token })
