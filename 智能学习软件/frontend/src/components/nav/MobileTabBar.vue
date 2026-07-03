<script setup lang="ts">
// MobileTabBar — 移动端底部 TabBar

import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const tabs = [
  { path: '/dashboard', label: '首页', icon: '🏠' },
  { path: '/subjects', label: '科目', icon: '📚' },
  { path: '/exercise', label: '练习', icon: '✏️' },
  { path: '/mistakes', label: '错题', icon: '📝' },
  { path: '/settings', label: '设置', icon: '⚙️' }
]

function go(path: string) {
  router.push(path)
}
</script>

<template>
  <nav class="mobile-tab-bar">
    <button
      v-for="t in tabs"
      :key="t.path"
      class="tab"
      :class="{ active: route.path.startsWith(t.path) }"
      @click="go(t.path)"
    >
      <span class="icon">{{ t.icon }}</span>
      <span class="label">{{ t.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.mobile-tab-bar {
  display: none;
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  height: 56px;
  background: var(--color-bg-card);
  border-top: 1px solid var(--color-border);
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.06);
  z-index: var(--z-dropdown);
}

.tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  color: var(--color-text-secondary);
  font-size: var(--font-tiny);
  transition: color 0.15s;
}

.tab.active {
  color: var(--color-primary);
}

.icon {
  font-size: 20px;
}

.label {
  font-size: var(--font-tiny);
}

@media (max-width: 768px) {
  .mobile-tab-bar { display: flex; }
}
</style>