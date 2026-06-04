<template>
  <div class="not-found">
    <div class="not-found-card">
      <div class="error-code">404</div>
      <h1>页面不存在</h1>
      <p>您访问的页面不存在或已被移除。</p>
      <el-button type="primary" @click="goDashboard">返回仪表盘</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goDashboard() {
  router.push('/dashboard')
}
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f7fa;
}
.not-found-card {
  text-align: center;
  background: #fff;
  border-radius: 12px;
  padding: 60px 80px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}
.error-code {
  font-size: 80px;
  font-weight: 700;
  color: #409eff;
  line-height: 1;
  margin-bottom: 16px;
}
h1 {
  font-size: 24px;
  color: #303133;
  margin: 0 0 12px;
}
p {
  color: #909399;
  margin: 0 0 32px;
}
</style>
