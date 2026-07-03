<script setup lang="ts">
// AppInput — 包装 Element Plus el-input

interface Props {
  modelValue?: string | number
  type?: 'text' | 'password' | 'email' | 'tel' | 'number' | 'textarea'
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  showPassword?: boolean
  prefixIcon?: string
  maxlength?: number
  rows?: number
  autocomplete?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
  (e: 'blur', evt: FocusEvent): void
  (e: 'focus', evt: FocusEvent): void
  (e: 'enter', evt: KeyboardEvent): void
}>()

function onInput(value: string | number) {
  emit('update:modelValue', value)
}

function onBlur(e: FocusEvent) { emit('blur', e) }
function onFocus(e: FocusEvent) { emit('focus', e) }
function onEnter(e: KeyboardEvent) { emit('enter', e) }

// 保留 props 引用避免 vite-tree-shake 移除
void props
</script>

<template>
  <el-input
    :model-value="modelValue"
    :type="type === 'textarea' ? 'textarea' : type"
    :placeholder="placeholder"
    :disabled="disabled"
    :readonly="readonly"
    :clearable="clearable"
    :show-password="showPassword"
    :prefix-icon="prefixIcon"
    :maxlength="maxlength"
    :rows="rows"
    :autocomplete="autocomplete"
    class="app-input"
    @update:model-value="onInput"
    @blur="onBlur"
    @focus="onFocus"
    @keyup.enter="onEnter"
  />
</template>

<style scoped>
.app-input { width: 100%; }
</style>