<script setup lang="ts">
// ExerciseView — 智能选题（Phase B）
//
// 数据：GET /exercises?subject_id=&difficulty=
// 操作：智能推荐 / 提交答案 / 反馈正确性

import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { exerciseApi } from '@/api/exercise'
import { subjectApi } from '@/api/subject'
import type { Question, Subject, DifficultyLevel } from '@/types/domain'

const questions = ref<Question[]>([])
const loading = ref(true)
const submitting = ref(false)
const lastResult = ref<{ correct: boolean; correct_answer: string } | null>(null)

const filters = reactive<{ subject_id?: number; difficulty?: DifficultyLevel }>({})

const subjectOptions = ref<{ label: string; value: number }[]>([])
const difficultyOptions = [
  { label: '★☆☆☆☆', value: 1 },
  { label: '★★☆☆☆', value: 2 },
  { label: '★★★☆☆', value: 3 },
  { label: '★★★★☆', value: 4 },
  { label: '★★★★★', value: 5 }
]

async function loadQuestions() {
  loading.value = true
  try {
    const { data } = await exerciseApi.list({
      subject_id: filters.subject_id,
      difficulty: filters.difficulty
    })
    questions.value = data.data.items || []
  } catch (e) {
    ElMessage.error('加载题目失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadRecommend() {
  loading.value = true
  try {
    const { data } = await exerciseApi.recommend(10)
    questions.value = data.data.items || []
    ElMessage.success('已为你智能推荐 10 道题')
  } catch (e) {
    ElMessage.error('推荐失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function onSubmit(payload: { question_id: number; answer: string }) {
  submitting.value = true
  try {
    const { data } = await exerciseApi.submit(payload)
    const r = data.data
    lastResult.value = { correct: r.correct, correct_answer: r.correct_answer }
    if (r.correct) {
      ElMessage.success('答对啦！')
    } else {
      ElMessage.warning(`答错了，正确答案：${r.correct_answer}`)
    }
  } catch (e) {
    ElMessage.error('提交失败')
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const { data } = await subjectApi.list()
    subjectOptions.value = (data.data.items || []).map((s: Subject) => ({
      label: s.name,
      value: s.id
    }))
  } catch {
    // ignore
  }
  await loadQuestions()
})
</script>

<template>
  <div class="exercise">
    <header class="bar">
      <AppSelect
        v-model="filters.subject_id"
        :options="subjectOptions"
        placeholder="选择科目"
        clearable
      />
      <AppSelect
        v-model="filters.difficulty"
        :options="difficultyOptions"
        placeholder="难度"
        clearable
      />
      <AppButton @click="loadQuestions">筛 选</AppButton>
      <AppButton type="primary" plain @click="loadRecommend">智能推荐</AppButton>
    </header>

    <section v-if="loading" class="skeleton-list">
      <AppCard v-for="i in 3" :key="i" class="skel-q" />
    </section>

    <AppEmpty v-else-if="!questions.length" description="暂无可练习题目，试试智能推荐" icon="📭" />

    <ol v-else class="question-list">
      <li v-for="q in questions" :key="q.id">
        <QuestionCard :question="q" :loading="submitting" @submit="onSubmit" />
      </li>
    </ol>

    <div v-if="lastResult" class="last-result" :class="{ correct: lastResult.correct }">
      {{ lastResult.correct ? '✅ 上一题答对' : `❌ 上一题答错，正确：${lastResult.correct_answer}` }}
    </div>
  </div>
</template>

<style scoped>
.exercise { display: flex; flex-direction: column; gap: var(--space-4); }

.bar {
  display: flex;
  gap: var(--space-2);
  align-items: center;
  flex-wrap: wrap;
}

.bar .app-select { width: 160px; }

.skeleton-list { display: flex; flex-direction: column; gap: var(--space-4); }
.skel-q { height: 200px; }

.question-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: 0;
}

.last-result {
  padding: var(--space-3);
  border-radius: var(--radius-md);
  text-align: center;
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: var(--fw-medium);
}
.last-result.correct {
  background: #f6ffed;
  color: var(--color-success);
}
</style>