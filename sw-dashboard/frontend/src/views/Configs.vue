<template>
  <div>
    <div class="page-header">
      <h2>配置管理</h2>
    </div>

    <el-tabs v-model="activeTab" @tab-click="handleTabSwitch">
      <el-tab-pane label="config.yaml" name="config.yaml">
        <template #label>
          <span>
            <el-icon><Setting /></el-icon> config.yaml
            <el-tag v-if="isDirty && activeTab === 'config.yaml'" size="small" type="danger" style="margin-left: 4px">未保存</el-tag>
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane label=".env" name=".env">
        <template #label>
          <span>
            <el-icon><Key /></el-icon> .env
            <el-tag v-if="isDirty && activeTab === '.env'" size="small" type="danger" style="margin-left: 4px">未保存</el-tag>
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane label="CLAUDE.md" name="CLAUDE.md">
        <template #label>
          <span>
            <el-icon><Document /></el-icon> CLAUDE.md
            <el-tag v-if="isDirty && activeTab === 'CLAUDE.md'" size="small" type="danger" style="margin-left: 4px">未保存</el-tag>
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane label="SOUL.md" name="SOUL.md">
        <template #label>
          <span>
            <el-icon><Document /></el-icon> SOUL.md
            <el-tag v-if="isDirty && activeTab === 'SOUL.md'" size="small" type="danger" style="margin-left: 4px">未保存</el-tag>
          </span>
        </template>
      </el-tab-pane>
      <el-tab-pane label="Cron 任务" name="cron">
        <template #label>
          <span>
            <el-icon><Timer /></el-icon> Cron 任务
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!-- 配置编辑器 -->
    <div v-if="activeTab !== 'cron'" v-loading="configLoading">
      <div class="file-info-bar">
        <div style="display: flex; align-items: center; gap: 10px">
          <span style="font-weight: 500">{{ currentConfig?.name }}</span>
          <span v-if="currentConfig?.last_modified" style="font-size: 12px; color: var(--text-secondary)">
            最后修改: {{ currentConfig.last_modified }}
          </span>
        </div>
        <div style="display: flex; gap: 8px; align-items: center">
          <el-button v-if="activeTab === '.env'" size="small" text @click="showRawEnv" :disabled="!masked">
            显示原始值
          </el-button>
          <el-button size="small" type="primary" :disabled="!isDirty" @click="handleSave">
            保存
          </el-button>
          <span style="font-size: 12px; color: var(--text-secondary)">Ctrl+S</span>
        </div>
      </div>
      <CodeEditor
        v-if="currentConfig"
        :content="editedContent"
        :language="currentConfig.language || 'text'"
        :readonly="false"
        height="calc(100vh - 320px)"
        @update:content="onContentChange"
      />
    </div>

    <!-- Cron 任务管理 -->
    <div v-if="activeTab === 'cron'">
      <h3 style="margin: 16px 0 12px">Cron 定时任务</h3>
      <el-table :data="cronJobs" v-loading="cronLoading" empty-text="暂无定时任务">
        <el-table-column prop="name" label="任务名称" min-width="120" />
        <el-table-column prop="schedule" label="调度表达式" min-width="120" />
        <el-table-column prop="command" label="命令" min-width="200" show-overflow-tooltip />
        <el-table-column label="上次运行" width="170">
          <template #default="{ row }">
            {{ row.last_run || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <StatusBadge :status="row.enabled ? 'active' : 'inactive'" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-switch
              :model-value="row.enabled"
              @change="(val: boolean) => toggleCronJob(row, val)"
              size="small"
            />
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getConfig, saveConfig, getCronJobs, toggleCronJob as toggleCronApi } from '@/api/configs'
import CodeEditor from '@/components/CodeEditor.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { showConfirm } from '@/components/ConfirmDialog'
import { ElMessage } from 'element-plus'
import { Setting, Key, Document, Timer } from '@element-plus/icons-vue'

const activeTab = ref('config.yaml')
const currentConfig = ref<any>(null)
const editedContent = ref('')
const isDirty = ref(false)
const configLoading = ref(false)
const cronLoading = ref(false)
const cronJobs = ref<any[]>([])
const masked = ref(true)

async function loadConfig(name: string, mask = true) {
  configLoading.value = true
  try {
    const res = await getConfig(name, name === '.env' ? mask : undefined)
    currentConfig.value = res.data?.data || res.data
    editedContent.value = currentConfig.value?.content || ''
    isDirty.value = false
    masked.value = mask
  } catch {
    currentConfig.value = null
  } finally {
    configLoading.value = false
  }
}

async function loadCronJobs() {
  cronLoading.value = true
  try {
    const res = await getCronJobs()
    cronJobs.value = res.data?.data || res.data || []
  } catch {
    cronJobs.value = []
  } finally {
    cronLoading.value = false
  }
}

async function handleTabSwitch(tab: any) {
  if (isDirty.value) {
    const confirmed = await showConfirm({
      title: '放弃修改',
      message: '当前配置有未保存的修改，是否放弃？'
    })
    if (!confirmed) {
      activeTab.value = tab.props.name
      return
    }
  }

  if (tab.props.name === 'cron') {
    loadCronJobs()
  } else {
    loadConfig(tab.props.name)
  }
}

function onContentChange(newContent: string) {
  editedContent.value = newContent
  isDirty.value = editedContent.value !== currentConfig.value?.content
}

async function handleSave() {
  try {
    await saveConfig(activeTab.value, editedContent.value)
    ElMessage.success('保存成功')
    isDirty.value = false
    currentConfig.value = { ...currentConfig.value, content: editedContent.value, last_modified: new Date().toISOString() }
  } catch (err: any) {
    const detail = err.response?.data?.detail || '保存失败'
    ElMessage.error('保存失败: ' + detail)
  }
}

async function showRawEnv() {
  const confirmed = await showConfirm({
    title: '敏感信息警告',
    message: '显示原始值会暴露敏感信息。确认显示？',
    type: 'warning'
  })
  if (confirmed) {
    await loadConfig('.env', false)
  }
}

async function toggleCronJob(row: any, enable: boolean) {
  try {
    await toggleCronApi(row.command, enable)
    row.enabled = enable
    ElMessage.success(`Cron 任务已${enable ? '启用' : '禁用'}`)
  } catch {
    ElMessage.error('操作失败')
  }
}

function handleKeydown(event: KeyboardEvent) {
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault()
    if (isDirty.value && activeTab.value !== 'cron') {
      handleSave()
    }
  }
}

onMounted(() => {
  loadConfig('config.yaml')
  loadCronJobs()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.file-info-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  margin: 8px 0;
}
</style>
