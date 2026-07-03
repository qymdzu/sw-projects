<script setup lang="ts">
// AppButton — 包装 Element Plus el-button，统一 design tokens

import { computed } from 'vue'

interface Props {
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  size?: 'large' | 'default' | 'small'
  loading?: boolean
  disabled?: boolean
  block?: boolean
  plain?: boolean
  round?: boolean
  nativeType?: 'button' | 'submit' | 'reset'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  size: 'default',
  nativeType: 'button'
})

const emit = defineEmits<{ (e: 'click', evt: MouseEvent): void }>()

const elType = computed(() => {
  if (props.type === 'default' || props.type === 'info') return props.type === 'info' ? 'info' : undefined
  return props.type
})

function onClick(e: MouseEvent) {
  emit('click', e)
}
</script>

<template>
  <el-button
    :type="elType"
    :size="size"
    :loading="loading"
    :disabled="disabled"
    :plain="plain"
    :round="round"
    :native-type="nativeType"
    class="app-button"
    :class="{ 'is-block': block }"
    @click="onClick"
  >
    <slot />
  </el-button>
</template>

<style scoped>
.app-button {
  font-weight: var(--fw-medium);
}
.app-button.is-block { display: block; width: 100%; }
</style>