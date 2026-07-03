<script setup lang="ts">
// SubjectsView — 科目浏览（Phase B）
//
// 数据：GET /subjects
// 交互：点击科目卡片 → 跳转到 /exercise?subject_id=...

import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { subjectApi } from '@/api/subject'
import type { Subject } from '@/types/domain'

const subjects = ref<Subject[]>([])
const loading = ref(true)
const router = useRouter()

onMounted(async () => {
  try {
    const { data } = await subjectApi.list()
    subjects.value = data.data.items || []
  } catch (e) {
    ElMessage.error('加载科目失败')
    console.error(e)
  } finally {
    loading.value = false
  }
})

function openSubject(id: number) {
  router.push({ path: '/exercise', query: { subject_id: id } })
}
</script>

<template>
  <div class="subjects">
    <section v-if="loading" class="skeleton-grid">
      <AppCard v-for="i in 6" :key="i" class="skeleton-card">
        <div class="skel-bar" />
      </AppCard>
    </section>

    <AppEmpty v-else-if="!subjects.length" description="暂无科目数据" icon="📚" />

    <ul v-else class="subject-list">
      <li v-for="s in subjects" :key="s.id" class="subject-item">
        <AppCard hoverable padding="md" shadow="md" @click="openSubject(s.id)">
          <div class="subject-head">
            <span class="subject-icon">{{ s.icon || '📘' }}</span>
            <span class="subject-name">{{ s.name }}</span>
          </div>
          <p class="subject-desc">{{ s.description || '点击进入练习' }}</p>
        </AppCard>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.subjects { display: flex; flex-direction: column; gap: var(--space-4); }

.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
}
@media (max-width: 768px) {
  .skeleton-grid { grid-template-columns: 1fr; }
}
.skel-bar { height: 80px; background: var(--color-border-light); border-radius: var(--radius-md); }

.subject-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
}
@media (max-width: 768px) {
  .subject-list { grid-template-columns: 1fr; }
}

.subject-head {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-bottom: var(--space-2);
}

.subject-icon { font-size: 24px; }

.subject-name {
  font-size: var(--font-h3);
  font-weight: var(--fw-semibold);
  color: var(--color-text-primary);
}

.subject-desc {
  font-size: var(--font-body);
  color: var(--color-text-secondary);
  margin: 0;
}
</style>