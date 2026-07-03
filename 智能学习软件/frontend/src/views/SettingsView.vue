<script setup lang="ts">
// SettingsView — 设置页（含模型配置 — 核心新增）
//
// 业务要点：
//   - API Key 保存后立即回显（service 层 Save 也会传回明文，handler 掩码为 ***；前端用本地的输入回填）
//   - 高级参数必须是合法 JSON
//   - 6 个 Provider：openai / anthropic / qwen / deepseek / ollama / custom
//   - 删除配置时若删的是 default，自动回落最近一条

import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useSettingStore } from '@/stores/setting'
import { useAuthStore } from '@/stores/auth'
import {
  type Provider,
  ProviderLabels,
  ProviderEndpoints,
  ProviderModels
} from '@/types/domain'

const store = useSettingStore()
const auth = useAuthStore()
const router = useRouter()

const providerOptions = (Object.keys(ProviderLabels) as Provider[]).map(p => ({
  label: ProviderLabels[p],
  value: p
}))

const form = reactive<{
  provider: Provider
  api_endpoint: string
  api_key: string
  model: string
}>({
  provider: 'openai',
  api_endpoint: ProviderEndpoints.openai,
  api_key: '',
  model: 'gpt-4o-mini'
})

const extraJSON = ref('{}')
const loading = ref(false)

// 当前选中 Provider 的候选模型列表（用于下拉建议）
const modelOptions = computed<string[]>(() => ProviderModels[form.provider] || [])

function fetchModelSuggestions(q: string, cb: (list: { value: string }[]) => void) {
  cb(modelOptions.value.filter((m: string) => m.includes(q)).map((m: string) => ({ value: m })))
}

// 切换 Provider 时同步默认值
watch(() => form.provider, (p) => {
  form.api_endpoint = ProviderEndpoints[p] || ''
  form.model = (ProviderModels[p] && ProviderModels[p][0]) || ''
})

const active = computed(() => store.active)

// 当远端加载到 active 时回填表单
watch(active, (v) => {
  if (v) {
    form.provider = v.provider
    form.api_endpoint = v.api_endpoint
    form.model = v.model
    form.api_key = v.api_key   // 后端 GetActive 回明文
    extraJSON.value = JSON.stringify(v.extra_config ?? {}, null, 2)
  }
}, { immediate: true })

async function onSave() {
  let extra: Record<string, unknown> = {}
  if (extraJSON.value.trim()) {
    try {
      extra = JSON.parse(extraJSON.value)
    } catch {
      ElMessage.error('高级参数必须是合法 JSON')
      return
    }
  }
  if (!form.api_key || form.api_key.length < 8) {
    ElMessage.warning('API Key 至少 8 个字符')
    return
  }
  loading.value = true
  try {
    await store.save({
      provider: form.provider,
      api_endpoint: form.api_endpoint,
      api_key: form.api_key,
      model: form.model,
      extra_config: extra
    })
    ElMessage.success('已保存，当前 LLM 已激活')
  } catch (e) {
    const msg = (e as { response?: { data?: { message?: string } } })?.response?.data?.message
      || (e instanceof Error ? e.message : '保存失败')
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}

async function onRemove() {
  if (!active.value) return
  try {
    await ElMessageBox.confirm(
      `确认删除 ${ProviderLabels[active.value.provider]} 配置？`,
      '提示',
      { type: 'warning' }
    )
  } catch {
    return
  }
  try {
    await store.remove(active.value.provider)
    ElMessage.success('已删除')
  } catch (e) {
    ElMessage.error('删除失败')
    console.error(e)
  }
}

async function onLogout() {
  await auth.logout()
  router.replace('/login')
}

onMounted(() => {
  store.loadActive()
})
</script>

<template>
  <div class="settings">
    <!-- 模型配置 Section — 核心 -->
    <AppCard padding="lg" shadow="md" class="section">
      <h3 class="section-title">模型配置</h3>
      <p class="section-desc">
        配置你的私有 LLM（用于 AI 推荐与排课）。API Key 通过 AES-256-GCM 加密后存于服务端。
      </p>

      <el-alert
        v-if="!active"
        type="info"
        :closable="false"
        show-icon
        title="尚未配置模型"
        description="将使用平台默认 LLM，配置后可解锁个性化推荐能力。"
        class="alert"
      />

      <el-form label-position="top" class="model-form">
        <el-form-item label="Provider">
          <AppSelect
            v-model="form.provider"
            :options="providerOptions"
            placeholder="选择 LLM 服务商"
          />
        </el-form-item>

        <el-form-item label="API Endpoint">
          <AppInput
            v-model="form.api_endpoint"
            placeholder="https://api.openai.com/v1"
            clearable
          />
        </el-form-item>

        <el-form-item label="API Key">
          <AppInput
            v-model="form.api_key"
            type="password"
            show-password
            placeholder="sk-..."
            clearable
          />
        </el-form-item>

        <el-form-item label="Model">
          <el-autocomplete
            v-model="form.model"
            :fetch-suggestions="fetchModelSuggestions"
            placeholder="gpt-4o-mini / claude-3-5-sonnet-20241022 / ..."
            clearable
            class="model-input"
          />
        </el-form-item>

        <el-form-item label="高级参数（可选，合法 JSON）">
          <AppInput
            v-model="extraJSON"
            type="textarea"
            :rows="4"
            placeholder='{"temperature":0.7,"top_p":0.9,"max_tokens":2048}'
          />
        </el-form-item>

        <div class="form-actions">
          <AppButton type="primary" :loading="loading || store.loading" @click="onSave">
            保 存
          </AppButton>
          <AppButton
            type="danger"
            plain
            :disabled="!active"
            @click="onRemove"
          >
            删除配置
          </AppButton>
        </div>
      </el-form>
    </AppCard>

    <!-- 账号 Section -->
    <AppCard padding="lg" shadow="md" class="section">
      <h3 class="section-title">账号</h3>
      <div class="account-row">
        <div class="account-info">
          <div class="name">{{ auth.user?.name || '未登录' }}</div>
          <div class="role">{{ auth.user?.role || '-' }}</div>
        </div>
        <AppButton @click="onLogout">退出登录</AppButton>
      </div>
    </AppCard>
  </div>
</template>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.section { display: flex; flex-direction: column; }

.alert { margin-bottom: var(--space-4); }

.model-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.model-input { width: 100%; }

.form-actions {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-3);
}

.account-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.account-info .name {
  font-size: var(--font-h3);
  font-weight: var(--fw-medium);
  color: var(--color-text-primary);
}
.account-info .role {
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  margin-top: var(--space-1);
}
</style>