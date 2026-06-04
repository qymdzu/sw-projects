<template>
  <div>
    <div class="page-header">
      <h2>仪表盘</h2>
      <div class="page-header-actions">
        <el-switch v-model="autoRefresh" active-text="自动刷新" size="small" style="margin-right: 8px" />
        <el-button size="small" @click="fetchDashboard" :loading="loading" :icon="Refresh">
          刷新
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && !dashboardData" v-loading="true" style="height: 400px"></div>

    <!-- 数据 -->
    <template v-if="dashboardData">
      <el-row :gutter="16">
        <!-- Gateway 状态 -->
        <el-col :span="8">
          <el-card>
            <template #header>
              <div style="display: flex; align-items: center; justify-content: space-between">
                <span>🚀 Gateway 状态</span>
                <StatusBadge :status="dashboardData.gateway?.status || 'unknown'" />
              </div>
            </template>
            <div v-if="dashboardData.gateway">
              <el-descriptions :column="1" size="small" border>
                <el-descriptions-item label="运行时长">
                  {{ dashboardData.gateway.uptime || 'N/A' }}
                </el-descriptions-item>
                <el-descriptions-item label="PID">
                  {{ dashboardData.gateway.pid || 'N/A' }}
                </el-descriptions-item>
                <el-descriptions-item label="内存">
                  {{ dashboardData.gateway.memory_mb ? dashboardData.gateway.memory_mb + ' MB' : 'N/A' }}
                </el-descriptions-item>
                <el-descriptions-item label="CPU">
                  {{ dashboardData.gateway.cpu_percent != null ? dashboardData.gateway.cpu_percent + '%' : 'N/A' }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            <el-empty v-else description="数据不可用" />
          </el-card>
        </el-col>

        <!-- 会话统计 -->
        <el-col :span="8">
          <el-card>
            <template #header>
              <span>💬 会话统计</span>
            </template>
            <div v-if="dashboardData.sessions" style="display: flex; gap: 16px; flex-wrap: wrap">
              <div style="flex: 1; text-align: center">
                <div class="stat-number">{{ dashboardData.sessions.total_sessions ?? 0 }}</div>
                <div class="stat-label">总会话</div>
              </div>
              <div style="flex: 1; text-align: center">
                <div class="stat-number">{{ dashboardData.sessions.active_sessions ?? 0 }}</div>
                <div class="stat-label">活跃会话</div>
              </div>
              <div style="flex: 1; text-align: center">
                <div class="stat-number">{{ dashboardData.sessions.total_messages ?? 0 }}</div>
                <div class="stat-label">总消息</div>
              </div>
              <div style="flex: 1; text-align: center">
                <div class="stat-number">{{ Math.round(dashboardData.sessions.avg_duration_sec ?? 0) }}</div>
                <div class="stat-label">平均时长(s)</div>
              </div>
            </div>
            <el-empty v-else description="数据不可用" />
          </el-card>
        </el-col>

        <!-- Cron 任务 -->
        <el-col :span="8">
          <el-card>
            <template #header>
              <span>⏰ 定时任务</span>
            </template>
            <div v-if="dashboardData.cron_jobs && dashboardData.cron_jobs.length > 0">
              <el-table :data="dashboardData.cron_jobs" size="small" max-height="220">
                <el-table-column prop="name" label="名称" min-width="80" />
                <el-table-column prop="schedule" label="调度" min-width="100" />
                <el-table-column label="状态" width="70">
                  <template #default="{ row }">
                    <StatusBadge :status="row.enabled ? 'active' : 'inactive'" />
                  </template>
                </el-table-column>
              </el-table>
            </div>
            <el-empty v-else description="暂无定时任务" />
          </el-card>
        </el-col>
      </el-row>

      <!-- Pipeline 进度 -->
      <el-card style="margin-top: 16px">
        <template #header>
          <span>📋 Pipeline 进度</span>
        </template>
        <div v-if="dashboardData.pipeline && dashboardData.pipeline.stages">
          <el-progress
            :percentage="Math.round((dashboardData.pipeline.overall_progress || 0) * 100)"
            :status="dashboardData.pipeline.overall_progress >= 1 ? 'success' : undefined"
            style="margin-bottom: 20px"
          />
          <div style="display: flex; gap: 8px; flex-wrap: wrap">
            <div
              v-for="stage in dashboardData.pipeline.stages"
              :key="stage.stage_id"
              class="pipeline-stage"
            >
              <span class="stage-icon">
                <template v-if="stage.status === 'completed'">✅</template>
                <template v-else-if="stage.status === 'running'">🔄</template>
                <template v-else-if="stage.status === 'failed'">❌</template>
                <template v-else>⭕</template>
              </span>
              <span style="font-size: 12px; text-align: center">{{ stage.name }}</span>
              <el-progress
                :percentage="Math.round((stage.progress || 0) * 100)"
                :width="40"
                type="circle"
                :stroke-width="4"
                size="small"
              />
            </div>
          </div>
        </div>
        <el-empty v-else description="暂无 Pipeline 任务" />
      </el-card>
    </template>

    <el-empty v-else-if="!loading" description="加载失败，请重试" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { getDashboard } from '@/api/dashboard'
import StatusBadge from '@/components/StatusBadge.vue'
import { Refresh } from '@element-plus/icons-vue'

interface DashboardData {
  gateway?: any
  sessions?: any
  pipeline?: any
  cron_jobs?: any[]
  server_time?: string
  version?: string
}

const dashboardData = ref<DashboardData | null>(null)
const loading = ref(true)
const autoRefresh = ref(true)
let refreshTimer: ReturnType<typeof setInterval> | null = null

async function fetchDashboard() {
  try {
    const res = await getDashboard()
    dashboardData.value = res.data?.data || res.data
    loading.value = false
  } catch {
    loading.value = false
  }
}

function startAutoRefresh() {
  stopAutoRefresh()
  if (autoRefresh.value) {
    refreshTimer = setInterval(fetchDashboard, 30000)
  }
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

watch(autoRefresh, (val) => {
  if (val) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
})

onMounted(() => {
  fetchDashboard()
  startAutoRefresh()
  document.addEventListener('visibilitychange', handleVisibility)
})

onUnmounted(() => {
  stopAutoRefresh()
  document.removeEventListener('visibilitychange', handleVisibility)
})

function handleVisibility() {
  if (document.hidden) {
    stopAutoRefresh()
  } else if (autoRefresh.value) {
    fetchDashboard()
    startAutoRefresh()
  }
}
</script>
