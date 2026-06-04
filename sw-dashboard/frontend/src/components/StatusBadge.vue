<template>
  <span class="status-badge">
    <span class="status-dot" :class="statusClass"></span>
    <span class="status-text">{{ statusText }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  status?: string
}>(), {
  status: 'unknown'
})

const statusConfig: Record<string, { color: string; text: string }> = {
  active: { color: '#67C23A', text: '运行中' },
  completed: { color: '#67C23A', text: '已完成' },
  inactive: { color: '#909399', text: '已停止' },
  pending: { color: '#909399', text: '待定' },
  failed: { color: '#F56C6C', text: '异常' },
  running: { color: '#409EFF', text: '执行中' },
  success: { color: '#67C23A', text: '成功' },
  unknown: { color: '#909399', text: '未知' }
}

const statusClass = computed(() => {
  const s = props.status?.toLowerCase() || 'unknown'
  if (['active', 'completed', 'success'].includes(s)) return 'active'
  if (s === 'failed') return 'failed'
  if (s === 'running') return 'running'
  if (['inactive', 'pending'].includes(s)) return 'inactive'
  return 'inactive'
})

const statusText = computed(() => {
  const s = props.status?.toLowerCase() || 'unknown'
  return statusConfig[s]?.text || '未知'
})
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.status-text {
  font-size: 13px;
}
</style>
