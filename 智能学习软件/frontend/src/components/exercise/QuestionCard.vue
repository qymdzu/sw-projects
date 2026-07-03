<script setup lang="ts">
// QuestionCard — 单道题目卡片，支持单选/多选/填空/简答
//
// 视觉规范（docs/ux/UI组件规范.md §3）：
//   - 难度用 1-5 颗星 ★☆☆☆☆
//   - 答题反馈：✅ / ❌ Toast

import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { Question, DifficultyLevel } from '@/types/domain'

interface Props {
  question: Question
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), { loading: false })

const emit = defineEmits<{
  (e: 'submit', payload: { question_id: number; answer: string }): void
}>()

const userAnswer = ref('')
const selectedOption = ref<string>('')

// 难度 ★ 渲染
const difficultyStars = computed(() => {
  const n = props.question.difficulty as DifficultyLevel
  return '★'.repeat(n) + '☆'.repeat(5 - n)
})

function onSelectOption(opt: string) {
  selectedOption.value = opt
  userAnswer.value = opt
}

function onSubmit() {
  const ans = userAnswer.value.trim()
  if (!ans) {
    ElMessage.warning('请输入或选择答案')
    return
  }
  emit('submit', { question_id: props.question.id, answer: ans })
}
</script>

<template>
  <AppCard padding="md" shadow="md" class="question-card">
    <header class="header">
      <span class="difficulty" :title="`难度 ${question.difficulty}/5`">{{ difficultyStars }}</span>
      <span class="type">{{ question.type }}</span>
    </header>

    <div class="content" v-html="question.content"></div>

    <div v-if="question.type === 'single' && question.options?.length" class="options">
      <button
        v-for="opt in question.options"
        :key="opt"
        class="option"
        :class="{ selected: selectedOption === opt }"
        @click="onSelectOption(opt)"
      >
        {{ opt }}
      </button>
    </div>

    <div v-else-if="question.type === 'multiple' && question.options?.length" class="options">
      <label v-for="opt in question.options" :key="opt" class="option-check">
        <input type="checkbox" :value="opt" v-model="userAnswer" />
        {{ opt }}
      </label>
    </div>

    <div v-else class="answer-area">
      <AppInput
        v-model="userAnswer"
        type="textarea"
        :rows="3"
        :placeholder="question.type === 'essay' ? '请输入答案…' : '请填空…'"
      />
    </div>

    <footer class="actions">
      <AppButton type="primary" :loading="loading" @click="onSubmit">提交答案</AppButton>
    </footer>
  </AppCard>
</template>

<style scoped>
.question-card { margin-bottom: var(--space-4); }

.header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.difficulty { color: var(--color-warning); font-size: 14px; }
.type {
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  padding: 2px var(--space-2);
  background: var(--color-primary-light);
  color: var(--color-primary);
  border-radius: var(--radius-sm);
}

.content {
  font-size: var(--font-body-lg);
  line-height: var(--line-relaxed);
  margin-bottom: var(--space-4);
  color: var(--color-text-primary);
}

.options {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.option {
  text-align: left;
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-card);
  transition: all 0.15s;
  font-size: var(--font-body);
  color: var(--color-text-regular);
}

.option:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.option.selected {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: var(--fw-medium);
}

.option-check {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
}

.answer-area { margin-bottom: var(--space-4); }

.actions { display: flex; justify-content: flex-end; }
</style>