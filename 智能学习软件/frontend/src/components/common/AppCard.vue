<script setup lang="ts">
// AppCard — 通用卡片容器，PC/Mobile 通用
interface Props {
  hoverable?: boolean
  padding?: 'none' | 'sm' | 'md' | 'lg'
  shadow?: 'none' | 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  padding: 'md',
  shadow: 'md'
})

const emit = defineEmits<{ (e: 'click', evt: MouseEvent): void }>()

const padCls = {
  none: 'p-0',
  sm: 'p-3',
  md: 'p-4',
  lg: 'p-6'
}[props.padding]

const shadowCls = {
  none: 'shadow-none',
  sm: 'shadow-sm',
  md: 'shadow-md',
  lg: 'shadow-lg'
}[props.shadow]

function onClick(e: MouseEvent) { emit('click', e) }
</script>

<template>
  <div
    class="app-card"
    :class="[padCls, shadowCls, { hoverable }]"
    @click="onClick"
  >
    <slot />
  </div>
</template>

<style scoped>
.app-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  transition: all 0.15s ease;
  cursor: default;
}
.shadow-none { box-shadow: none; }
.shadow-sm { box-shadow: var(--shadow-sm); }
.shadow-md { box-shadow: var(--shadow-md); }
.shadow-lg { box-shadow: var(--shadow-lg); }

.p-0 { padding: 0; }
.p-3 { padding: var(--space-3); }
.p-4 { padding: var(--space-4); }
.p-6 { padding: var(--space-6); }

.hoverable {
  cursor: pointer;
}
.hoverable:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}
</style>