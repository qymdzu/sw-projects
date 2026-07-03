<script setup lang="ts">
// LoginView — 登录页（Phase B）
//
// 设计：单一卡片 + 智学蓝主色，输入校验失败 Toast 提示。
// 支持手机号 / 邮箱登录。

import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const appTitle = import.meta.env.VITE_APP_TITLE || '智学助手'
const loading = ref(false)
const mode = ref<'login' | 'register'>('login')

const loginForm = reactive({
  account: '',
  password: ''
})

const registerForm = reactive({
  name: '',
  account: '',
  password: ''
})

function validateAccount(a: string): boolean {
  return /^1[3-9]\d{9}$/.test(a) || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(a)
}

async function onLogin() {
  if (!loginForm.account || !loginForm.password) {
    ElMessage.warning('请填写账号和密码')
    return
  }
  if (!validateAccount(loginForm.account)) {
    ElMessage.warning('账号必须是手机号或邮箱')
    return
  }
  loading.value = true
  try {
    await auth.login(loginForm.account, loginForm.password)
    const redirect = (route.query.redirect as string) || '/dashboard'
    ElMessage.success('登录成功')
    router.replace(redirect)
  } catch (e) {
    const msg = (e as { response?: { data?: { message?: string } } })?.response?.data?.message
      || (e instanceof Error ? e.message : '登录失败')
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}

async function onRegister() {
  if (!registerForm.name || !registerForm.account || !registerForm.password) {
    ElMessage.warning('请完整填写注册信息')
    return
  }
  if (registerForm.password.length < 6) {
    ElMessage.warning('密码至少 6 位')
    return
  }
  loading.value = true
  try {
    const isEmail = registerForm.account.includes('@')
    await auth.register({
      name: registerForm.name,
      phone: isEmail ? undefined : registerForm.account,
      email: isEmail ? registerForm.account : undefined,
      password: registerForm.password
    })
    ElMessage.success('注册成功')
    router.replace('/dashboard')
  } catch (e) {
    const msg = (e as { response?: { data?: { message?: string } } })?.response?.data?.message
      || (e instanceof Error ? e.message : '注册失败')
    ElMessage.error(msg)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <AppCard padding="lg" shadow="lg" class="login-card">
      <header class="brand">
        <span class="logo">🎓</span>
        <h1 class="title">{{ appTitle }}</h1>
        <p class="subtitle">智能学习，从这里开始</p>
      </header>

      <div class="tab-switch">
        <button :class="{ active: mode === 'login' }" @click="mode = 'login'">登 录</button>
        <button :class="{ active: mode === 'register' }" @click="mode = 'register'">注 册</button>
      </div>

      <form v-if="mode === 'login'" class="form" @submit.prevent="onLogin">
        <AppInput v-model="loginForm.account" placeholder="手机号 / 邮箱" autocomplete="username" />
        <AppInput v-model="loginForm.password" type="password" placeholder="密码" show-password autocomplete="current-password" />
        <AppButton type="primary" block native-type="submit" :loading="loading">登 录</AppButton>
      </form>

      <form v-else class="form" @submit.prevent="onRegister">
        <AppInput v-model="registerForm.name" placeholder="昵称" maxlength="20" />
        <AppInput v-model="registerForm.account" placeholder="手机号 / 邮箱" />
        <AppInput v-model="registerForm.password" type="password" placeholder="密码（至少 6 位）" show-password />
        <AppButton type="primary" block native-type="submit" :loading="loading">注 册</AppButton>
      </form>

      <p class="hint">
        {{ mode === 'login' ? '还没有账号？' : '已有账号？' }}
        <a href="#" @click.prevent="mode = mode === 'login' ? 'register' : 'login'">
          {{ mode === 'login' ? '立即注册' : '直接登录' }}
        </a>
      </p>
    </AppCard>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: var(--space-4);
  background: linear-gradient(135deg, var(--color-primary-light) 0%, #ffffff 100%);
}

.login-card {
  width: 100%;
  max-width: 420px;
}

.brand {
  text-align: center;
  margin-bottom: var(--space-6);
}

.logo {
  font-size: 48px;
  display: block;
  margin-bottom: var(--space-2);
}

.title {
  font-size: var(--font-h1);
  font-weight: var(--fw-bold);
  color: var(--color-primary);
  margin-bottom: var(--space-1);
}

.subtitle {
  font-size: var(--font-body);
  color: var(--color-text-secondary);
}

.tab-switch {
  display: flex;
  background: var(--color-bg-page);
  border-radius: var(--radius-md);
  padding: var(--space-1);
  margin-bottom: var(--space-4);
}

.tab-switch button {
  flex: 1;
  padding: var(--space-2) 0;
  font-size: var(--font-body);
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  transition: all 0.15s;
}

.tab-switch button.active {
  background: var(--color-bg-card);
  color: var(--color-primary);
  font-weight: var(--fw-medium);
  box-shadow: var(--shadow-sm);
}

.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.hint {
  margin-top: var(--space-4);
  text-align: center;
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
}

.hint a {
  color: var(--color-primary);
  margin-left: var(--space-1);
}
</style>