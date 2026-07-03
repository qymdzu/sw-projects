<script setup lang="ts">
// AppEmpty — 空状态展示
interface Props {
  description?: string
  icon?: string
  image?: string
  size?: 'small' | 'default' | 'large'
}

withDefaults(defineProps<Props>(), {
  description: '暂无数据',
  icon: '📭',
  size: 'default'
})
</script>

<template>
  <div class="app-empty" :class="size">
    <div v-if="image" class="image"><img :src="image" /></div>
    <div v-else class="icon">{{ icon }}</div>
    <p class="description">{{ description }}</p>
    <div v-if="$slots.default" class="extra">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.app-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12) var(--space-4);
  color: var(--color-text-secondary);
  text-align: center;
}

.icon {
  font-size: 56px;
  margin-bottom: var(--space-3);
  opacity: 0.6;
}

.description {
  font-size: var(--font-body);
  margin: 0;
}

.app-empty.small { padding: var(--space-6) var(--space-3); }
.app-empty.small .icon { font-size: 36px; }
.app-empty.large { padding: var(--space-16) var(--space-6); }
.app-empty.large .icon { font-size: 80px; }
</style>