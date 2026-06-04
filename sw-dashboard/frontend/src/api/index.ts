import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('dashboard_token')
  if (token) {
    config.headers['X-Token'] = token
  }
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    const status = error.response?.status
    const msg = error.response?.data?.detail || error.message

    if (status === 401) {
      localStorage.removeItem('dashboard_token')
      ElMessage.error('Token 已失效，请重新登录')
      router.push('/login')
    } else if (status === 403) {
      ElMessage.error('权限不足: ' + msg)
    } else if (status === 404) {
      ElMessage.warning('资源不存在: ' + msg)
    } else if (status === 413) {
      ElMessage.warning('文件过大: ' + msg)
    } else if (status === 400) {
      ElMessage.error('请求错误: ' + msg)
    } else {
      ElMessage.error(msg)
    }
    return Promise.reject(error)
  }
)

export default api
