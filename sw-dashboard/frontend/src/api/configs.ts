import api from './index'

export const getConfigs = () => api.get('/configs')

export const getConfig = (name: string, masked?: boolean) =>
  api.get(`/configs/${name}`, { params: { masked } })

export const saveConfig = (name: string, content: string) =>
  api.put(`/configs/${name}`, { content })

export const getCronJobs = () => api.get('/configs/cron')

export const toggleCronJob = (command: string, enable: boolean) =>
  api.post('/configs/cron/toggle', { command, enable })
