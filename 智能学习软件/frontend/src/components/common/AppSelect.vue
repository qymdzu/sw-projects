<script setup lang="ts">
// AppSelect — 包装 Element Plus el-select

interface Option { label: string; value: string | number }

interface Props {
  modelValue?: string | number
  options: Option[]
  placeholder?: string
  disabled?: boolean
  clearable?: boolean
  multiple?: boolean
}

const props = withDefaults(defineProps<Props>(), {})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number | (string | number)[]): void
  (e: 'change', value: string | number | (string | number)[]): void
}>()

function onChange(value: string | number | (string | number)[]) {
  emit('update:modelValue', value)
  emit('change', value)
}

void props
</script>

<template>
  <el-select
    :model-value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :clearable="clearable"
    :multiple="multiple"
    class="app-select"
    @update:model-value="onChange"
  >
    <el-option
      v-for="opt in options"
      :key="opt.value"
      :label="opt.label"
      :value="opt.value"
    />
  </el-select>
</template>

<style scoped>
.app-select { width: 100%; }
</style>