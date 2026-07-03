<script setup lang="ts">
// MistakesView — 错题本（Phase B）
//
// 数据：GET /mistakes?knowledge_point_id=&mastered=
// 操作：标记掌握 / 重做 / 复习

import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mistakeApi } from '@/api/mistake'
import type { Mistake } from '@/types/domain'

const mistakes = ref<Mistake[]>([])
const loading = ref(true)

const filters = reactive<{ knowledge_point_id?: number; mastered: boolean }>({
  mastered: false
})

const masteredOptions = [
  { label: '未掌握', value: false },
  { label: '已掌握', value: true }
]

const kpOptions = ref<{ label: string; value: number }[]>([])

async function loadMistakes() {
  loading.value = true
  try {
    const { data } = await mistakeApi.list({
      knowledge_point_id: filters.knowledge_point_id,
      mastered: filters.mastered
    })
    mistakes.value = data.data.items || []

    // 提取唯一知识点 → options
    const seen = new Set<number>()
    kpOptions.value = []
    for (const m of mistakes.value) {
      if (!seen.has(m.knowledge_point.id)) {
        seen.add(m.knowledge_point.id)
        kpOptions.value.push({ label: m.knowledge_point.name, value: m.knowledge_point.id })
      }
    }
  } catch (e) {
    ElMessage.error('加载错题失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function onMarkMastered(m: Mistake) {
  try {
    await mistakeApi.markMastered(m.id)
    ElMessage.success('已标记掌握')
    await loadMistakes()
  } catch {
    ElMessage.error('操作失败')
  }
}

async function onReplay(m: Mistake) {
  await ElMessageBox.confirm(`重新练习这道错题？`, '提示', { type: 'info' })
    .catch(() => null)
  // 简化：跳转练习页并预填题目 id
  ElMessage.info('请到"智能选题"页通过筛选定位题目')
}

onMounted(loadMistakes)
</script>

<template>
  <div class="mistakes">
    <header class="bar">
      <AppSelect
        v-model="filters.knowledge_point_id"
        :options="kpOptions"
        placeholder="知识点"
        clearable
      />
      <AppSelect
        v-model="filters.mastered"
        :options="masteredOptions"
        placeholder="掌握状态"
      />
      <AppButton @click="loadMistakes">筛 选</AppButton>
    </header>

    <section v-if="loading" class="skeleton-list">
      <AppCard v-for="i in 3" :key="i" class="skel-m" />
    </section>

    <AppEmpty v-else-if="!mistakes.length" description="暂无错题，继续保持！🎉" icon="🎉" />

    <ul v-else class="mistake-list">
      <li v-for="m in mistakes" :key="m.id">
        <AppCard padding="md" shadow="md">
          <div class="head">
            <span class="kp">{{ m.knowledge_point.name }}</span>
            <span class="count">错 {{ m.mistake_count }} 次</span>
          </div>
          <div class="q-content" v-html="m.question.content"></div>
          <div class="meta">
            <span class="wrong">你的答案：{{ m.wrong_answer }}</span>
            <span class="correct">正确答案：{{ m.question.answer }}</span>
          </div>
          <div class="actions">
            <AppButton size="small" @click="onReplay(m)">重做</AppButton>
            <AppButton size="small" type="success" @click="onMarkMastered(m)">标记掌握</AppButton>
          </div>
        </AppCard>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.mistakes { display: flex; flex-direction: column; gap: var(--space-4); }
.bar { display: flex; gap: var(--space-2); align-items: center; flex-wrap: wrap; }
.bar .app-select { width: 180px; }

.skeleton-list { display: flex; flex-direction: column; gap: var(--space-3); }
.skel-m { height: 140px; }

.mistake-list { display: flex; flex-direction: column; gap: var(--space-3); }

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-2);
}
.kp {
  font-size: var(--font-caption);
  color: var(--color-primary);
  background: var(--color-primary-light);
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
}
.count { font-size: var(--font-caption); color: var(--color-error); }

.q-content {
  font-size: var(--font-body-lg);
  line-height: var(--line-normal);
  color: var(--color-text-primary);
  margin-bottom: var(--space-2);
}

.meta {
  display: flex;
  gap: var(--space-4);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  margin-bottom: var(--space-3);
  flex-wrap: wrap;
}
.wrong { color: var(--color-error); }
.correct { color: var(--color-success); }

.actions { display: flex; gap: var(--space-2); justify-content: flex-end; }
</style>