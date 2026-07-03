<script setup lang="ts">
// App.vue 是根组件。
//
// 设计：使用 <router-view /> 让布局（MainLayout / BlankLayout）由路由 meta.layout 决定。
// 这里不持有任何业务状态，所有状态都在各 Pinia store。
</script>

<template>
  <router-view />
</template>

<style>
#app {
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC',
    'Microsoft YaHei', sans-serif;
  color: var(--color-text-primary);
  background: var(--color-bg-page);
}
</style>