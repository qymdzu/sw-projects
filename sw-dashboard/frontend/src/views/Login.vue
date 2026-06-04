<template>
  <el-container style="height: 100vh; justify-content: center; align-items: center; background: #141414">
    <div class="login-card">
      <div class="login-header">
        <span style="font-size: 40px">🌼</span>
        <h2>翠花集群管理面板</h2>
        <p style="color: var(--text-secondary); font-size: 13px; margin-top: 4px">sw-dashboard v0.1.0</p>
      </div>
      <el-form @submit.prevent="handleLogin" style="margin-top: 24px">
        <el-form-item>
          <el-input
            v-model="tokenInput"
            type="password"
            placeholder="请输入访问 Token"
            size="large"
            show-password
            @keyup.enter="handleLogin"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            style="width: 100%"
            :loading="loading"
            :disabled="!tokenInput.trim()"
            @click="handleLogin"
          >
            {{ loading ? '验证中...' : '登录' }}
          </el-button>
        </el-form-item>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
      </el-form>
      <div class="login-footer">
        <span>通过 Tailscale 内网安全连接</span>
      </div>
    </div>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTokenStore } from '@/stores/token'
import { verifyToken } from '@/api/dashboard'
import { ElMessage } from 'element-plus'

const router = useRouter()
const tokenStore = useTokenStore()

const tokenInput = ref('')
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  if (!tokenInput.value.trim()) return

  loading.value = true
  errorMsg.value = ''

  try {
    await verifyToken(tokenInput.value.trim())
    tokenStore.setToken(tokenInput.value.trim())
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (err: any) {
    if (err.response?.status === 401) {
      errorMsg.value = 'Token 无效，请确认后重试'
    } else {
      errorMsg.value = '网络连接失败，请检查网络'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-card {
  width: 400px;
  padding: 40px;
  background: #1e1e1e;
  border: 1px solid #333;
  border-radius: 8px;
}
.login-header {
  text-align: center;
}
.login-header h2 {
  margin-top: 12px;
  font-size: 20px;
  font-weight: 600;
}
.error-msg {
  color: var(--color-danger);
  font-size: 13px;
  text-align: center;
  margin-top: 8px;
}
.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
