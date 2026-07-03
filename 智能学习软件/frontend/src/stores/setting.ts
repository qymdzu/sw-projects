// Pinia 模型配置 store — Phase B（核心新增）

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { settingApi, type SaveSettingPayload } from '@/api/setting'
import type { ModelSettingFull, Provider } from '@/types/domain'

export const useSettingStore = defineStore('setting', () => {
  const active = ref<ModelSettingFull | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadActive(): Promise<void> {
    loading.value = true
    error.value = null
    try {
      active.value = await settingApi.getActive()
    } catch (e) {
      const msg = e instanceof Error ? e.message : '加载失败'
      error.value = msg
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * 保存并切换为 default；本地立即回填明文 key 以便用户继续使用。
   */
  async function save(input: SaveSettingPayload): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const dto = await settingApi.save(input)
      // 后端创建时 api_key 是 *** 掩码；用前端刚填的明文回填即可
      active.value = {
        provider: dto.provider,
        api_endpoint: dto.api_endpoint,
        model: dto.model,
        extra_config: dto.extra_config,
        is_default: dto.is_default,
        updated_at: dto.updated_at,
        api_key: input.api_key
      }
    } catch (e) {
      const msg = e instanceof Error ? e.message : '保存失败'
      error.value = msg
      throw e
    } finally {
      loading.value = false
    }
  }

  async function activate(provider: Provider): Promise<void> {
    await settingApi.setActive(provider)
    await loadActive()
  }

  async function remove(provider: Provider): Promise<void> {
    loading.value = true
    try {
      await settingApi.remove(provider)
      if (active.value?.provider === provider) {
        active.value = null
      }
    } finally {
      loading.value = false
    }
  }

  return { active, loading, error, loadActive, save, activate, remove }
})