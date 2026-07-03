// 模型配置 API — Phase B（核心新增模块）

import { http } from './client'
import type { Provider, ModelSetting, ModelSettingFull } from '@/types/domain'

export interface SaveSettingPayload {
  provider: Provider
  api_endpoint: string
  api_key: string
  model: string
  extra_config?: Record<string, unknown>
}

export interface ActiveSettingResponse {
  setting: ModelSettingFull | null
  message?: string
}

export const settingApi = {
  /** POST /settings/model — 保存/更新（自动设为 default） */
  async save(body: SaveSettingPayload): Promise<ModelSetting> {
    const { data } = await http.post<{ data: ModelSetting }>('/settings/model', body)
    return data.data
  },

  /** GET /settings/model — 读当前 default，未配置时返回 null */
  async getActive(): Promise<ModelSettingFull | null> {
    const { data } = await http.get<{ data: ActiveSettingResponse }>('/settings/model')
    return data.data.setting
  },

  /** PUT /settings/model — 切换 default 到已存在的 provider */
  async setActive(provider: Provider): Promise<void> {
    await http.put<{ code: number }>('/settings/model', { provider })
  },

  /** DELETE /settings/model?provider=xxx */
  async remove(provider: Provider): Promise<void> {
    await http.delete<{ code: number }>(`/settings/model?provider=${encodeURIComponent(provider)}`)
  }
}