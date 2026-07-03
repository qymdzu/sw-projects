<script setup lang="ts">
// MainLayout — 含 TopNavBar + 内容区 + 移动端底部 TabBar
//
// 适配：
//   - PC (>= 768px)：顶部横向导航
//   - Mobile (< 768px)：底部 TabBar

import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const title = computed(() => (route.meta.title as string) || '智学助手')
</script>

<template>
  <div class="main-layout">
    <TopNavBar />

    <main class="main-content">
      <header class="page-header">
        <h1 class="page-title">{{ title }}</h1>
      </header>
      <div class="page-body">
        <router-view />
      </div>
    </main>

    <MobileTabBar />
  </div>
</template>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: var(--space-6);
  padding-bottom: 96px; /* 给移动端 TabBar 留位置 */
}

.page-header {
  margin-bottom: var(--space-6);
}

.page-title {
  font-size: var(--font-h1);
  font-weight: var(--fw-semibold);
  color: var(--color-text-primary);
}

.page-body {
  min-height: 60vh;
}

@media (max-width: 768px) {
  .main-content {
    padding: var(--space-4);
  }
}
</style>