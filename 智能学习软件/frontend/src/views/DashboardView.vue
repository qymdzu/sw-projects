<script setup lang="ts">
// DashboardView — 仪表盘（Phase B）
//
// 展示用户学习概览：总练习、正确率、连续天数、未掌握错题、今日时长。
// 数据来源：GET /reports/summary

import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { reportApi } from '@/api/report'
import type { ReportSummary } from '@/types/domain'
import { formatPercent } from '@/utils/format'

const summary = ref<ReportSummary | null>(null)
const loading = ref(true)

const cards = [
  { key: 'total_exercises', label: '总练习', color: '#4F7CFF' },
  { key: 'overall_correct_rate', label: '正确率', color: '#52C41A', format: formatPercent },
  { key: 'streak_days', label: '连续天数', color: '#FF9F45' },
  { key: 'unmastered_mistakes', label: '未掌握错题', color: '#FF4D4F' }
] as const

function valueOf(key: string): string {
  if (!summary.value) return '-'
  const v = (summary.value as unknown as Record<string, unknown>)[key]
  if (typeof v !== 'number') return String(v ?? '-')
  const card = cards.find(c => c.key === key)
  return card && 'format' in card ? (card as any).format(v) : String(v)
}

onMounted(async () => {
  try {
    const { data } = await reportApi.summary()
    summary.value = data.data
  } catch (e) {
    ElMessage.error('加载概览失败')
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard">
    <p class="page-subtitle">今日已学习 {{ summary?.today_duration_min ?? 0 }} 分钟</p>

    <section v-if="loading" class="skeleton-grid">
      <AppCard v-for="i in 4" :key="i" shadow="none" class="skeleton-card">
        <div class="skel-bar" />
      </AppCard>
    </section>

    <section v-else class="stat-grid">
      <AppCard
        v-for="card in cards"
        :key="card.key"
        padding="md"
        shadow="md"
        hoverable
      >
        <div class="stat-num" :style="{ color: card.color }">{{ valueOf(card.key) }}</div>
        <div class="stat-label">{{ card.label }}</div>
      </AppCard>
    </section>

    <section class="quick-entries">
      <h2 class="section-title">快捷入口</h2>
      <div class="entries">
        <router-link to="/exercise" class="entry">
          <span class="entry-icon">✏️</span>
          <span class="entry-label">开始练习</span>
        </router-link>
        <router-link to="/mistakes" class="entry">
          <span class="entry-icon">📝</span>
          <span class="entry-label">查看错题</span>
        </router-link>
        <router-link to="/subjects" class="entry">
          <span class="entry-icon">📚</span>
          <span class="entry-label">浏览科目</span>
        </router-link>
        <router-link to="/settings" class="entry">
          <span class="entry-icon">⚙️</span>
          <span class="entry-label">模型设置</span>
        </router-link>
      </div>
    </section>
  </div>
</template>

<style scoped>
.dashboard { display: flex; flex-direction: column; gap: var(--space-6); }

.skeleton-grid, .stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}

@media (max-width: 768px) {
  .skeleton-grid, .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.stat-num {
  font-size: var(--font-display);
  font-weight: var(--fw-bold);
  margin-bottom: var(--space-2);
}

.stat-label {
  font-size: var(--font-body);
  color: var(--color-text-secondary);
}

.skel-bar {
  height: 60px;
  background: linear-gradient(90deg, var(--color-border-light), var(--color-border), var(--color-border-light));
  border-radius: var(--radius-md);
}

.entries {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}

@media (max-width: 768px) {
  .entries { grid-template-columns: repeat(2, 1fr); }
}

.entry {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-6) var(--space-4);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: all 0.15s;
  cursor: pointer;
}

.entry:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.entry-icon {
  font-size: 32px;
  margin-bottom: var(--space-2);
}

.entry-label {
  font-size: var(--font-body);
  color: var(--color-text-regular);
}
</style>