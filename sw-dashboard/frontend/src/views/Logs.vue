<template>
  <div>
    <div class="page-header">
      <h2>日志查看</h2>
      <div class="page-header-actions">
        <el-select v-model="selectedFile" size="small" @change="switchLogFile" style="width: 180px">
          <el-option
            v-for="f in logFiles"
            :key="f.name"
            :label="f.name"
            :value="f.name"
          />
        </el-select>
        <el-button size="small" @click="fetchTailLog" :icon="Refresh">刷新</el-button>
        <el-switch v-model="autoScroll" active-text="自动滚动" size="small" />
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索日志关键字..."
        size="small"
        style="width: 300px"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button size="small" type="primary" @click="handleSearch">搜索</el-button>
      <el-button v-if="searchMode" size="small" @click="clearSearch">清空搜索</el-button>
      <span v-if="searchResult" style="font-size: 13px; color: var(--text-secondary); margin-left: 8px">
        共找到 {{ searchResult.total_matches }} 条匹配
      </span>
    </div>

    <!-- 暂停提示条 -->
    <div v-if="isPaused && !searchMode" class="pause-bar">
      <span>新日志已暂停 · </span>
      <el-button text size="small" type="primary" @click="scrollToBottom">回到底部</el-button>
    </div>

    <!-- 日志内容区 -->
    <div
      ref="logContainerRef"
      class="log-container"
      @scroll="handleScroll"
      v-loading="loading"
    >
      <div v-if="!logContent && !loading" style="text-align: center; color: var(--text-secondary); padding: 40px">
        <el-empty description="暂无日志内容" />
      </div>
      <div v-else>
        <div
          v-for="(line, index) in logLines"
          :key="index"
          :class="{ 'log-line-highlight': searchMode && line.includes(searchKeyword) }"
        >{{ line }}</div>
      </div>
    </div>

    <!-- 历史归档 -->
    <el-collapse style="margin-top: 12px">
      <el-collapse-item title="📦 历史归档" name="archives">
        <el-table :data="archives" size="small" empty-text="暂无历史归档" @row-click="loadArchive">
          <el-table-column prop="name" label="文件名" min-width="200" />
          <el-table-column label="大小" width="100">
            <template #default="{ row }">
              {{ row.size ? (row.size / 1024).toFixed(1) + ' KB' : '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="mtime" label="修改时间" width="180" />
        </el-table>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { getLogFiles, tailLog, searchLogs, getArchives } from '@/api/logs'
import { getFile } from '@/api/storage'
import { Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const logFiles = ref<any[]>([])
const selectedFile = ref('hermes.log')
const logContent = ref('')
const loading = ref(false)
const autoScroll = ref(true)
const isPaused = ref(false)
const searchMode = ref(false)
const searchKeyword = ref('')
const searchResult = ref<any>(null)
const archives = ref<any[]>([])
const logContainerRef = ref<HTMLElement | null>(null)

let pollingTimer: ReturnType<typeof setInterval> | null = null

const logLines = computed(() => {
  if (!logContent.value) return []
  return logContent.value.split('\n')
})

async function fetchTailLog() {
  try {
    const res = await tailLog(selectedFile.value, 200)
    const data = res.data?.data || res.data
    logContent.value = data?.content || ''
    await nextTick()
    if (autoScroll.value && !isPaused.value) {
      scrollToBottom()
    }
  } catch {
    // 静默重试
  }
}

function startPolling() {
  stopPolling()
  fetchTailLog()
  pollingTimer = setInterval(fetchTailLog, 3000)
}

function stopPolling() {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

function switchLogFile(file: string) {
  selectedFile.value = file
  logContent.value = ''
  searchMode.value = false
  searchResult.value = null
  startPolling()
}

async function handleSearch() {
  if (!searchKeyword.value.trim()) return
  loading.value = true
  searchMode.value = true
  stopPolling()
  try {
    const res = await searchLogs(searchKeyword.value, selectedFile.value)
    const data = res.data?.data || res.data
    searchResult.value = data
    logContent.value = (data?.matches || []).join('\n')
  } catch {
    ElMessage.error('搜索失败')
    searchMode.value = false
  } finally {
    loading.value = false
    await nextTick()
    if (logContainerRef.value) {
      logContainerRef.value.scrollTop = 0
    }
  }
}

function clearSearch() {
  searchMode.value = false
  searchResult.value = null
  searchKeyword.value = ''
  startPolling()
}

function handleScroll() {
  const el = logContainerRef.value
  if (!el) return
  const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
  if (distanceFromBottom < 50) {
    isPaused.value = false
  } else if (!searchMode.value) {
    isPaused.value = true
  }
}

function scrollToBottom() {
  const el = logContainerRef.value
  if (el) {
    el.scrollTop = el.scrollHeight
    isPaused.value = false
  }
}

async function loadArchive(archive: any) {
  stopPolling()
  searchMode.value = true
  try {
    const res = await getFile(archive.path)
    const data = res.data?.data || res.data
    logContent.value = data?.content || ''
  } catch {
    ElMessage.error('加载归档文件失败')
  }
}

async function loadLogFiles() {
  try {
    const res = await getLogFiles()
    const data = res.data?.data || res.data
    logFiles.value = data?.files || []
    if (data?.default_file) {
      selectedFile.value = data.default_file
    }
  } catch {
    logFiles.value = []
  }
}

async function loadArchives() {
  try {
    const res = await getArchives()
    archives.value = res.data?.data || res.data || []
  } catch {
    archives.value = []
  }
}

onMounted(() => {
  loadLogFiles()
  loadArchives()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.pause-bar {
  background: #1a2a3a;
  border: 1px solid #2a4a6a;
  border-radius: 4px;
  padding: 6px 12px;
  font-size: 13px;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}
</style>
