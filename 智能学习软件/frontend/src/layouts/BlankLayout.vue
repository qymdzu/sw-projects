<script setup lang="ts">
// BlankLayout — 用于登录、404 等无侧边栏页面
</script>

<template>
  <div class="blank-layout">
    <slot />
  </div>
</template>

<style scoped>
.blank-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-primary-light) 0%, #ffffff 100%);
}
</style>