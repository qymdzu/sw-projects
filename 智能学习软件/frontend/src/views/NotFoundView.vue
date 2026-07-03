<script setup lang="ts">
// NotFoundView — 404 页面
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <div class="not-found">
    <AppEmpty icon="🤔" description="页面不存在 (404)">
      <AppButton type="primary" @click="router.replace('/dashboard')">回到首页</AppButton>
    </AppEmpty>
  </div>
</template>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>