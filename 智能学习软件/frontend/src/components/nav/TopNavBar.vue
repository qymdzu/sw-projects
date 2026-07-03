<script setup lang="ts">
// TopNavBar — PC 端顶部横向导航

import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const navItems = [
  { path: '/dashboard', label: '仪表盘', icon: 'House' },
  { path: '/subjects', label: '科目', icon: 'Reading' },
  { path: '/exercise', label: '智能选题', icon: 'EditPen' },
  { path: '/mistakes', label: '错题本', icon: 'Warning' },
  { path: '/settings', label: '设置', icon: 'Setting' }
]

const activePath = computed(() => route.path)

function go(path: string) {
  router.push(path)
}

async function onLogout() {
  await auth.logout()
  router.replace('/login')
}
</script>

<template>
  <header class="top-nav">
    <div class="brand" @click="go('/dashboard')">
      <span class="brand-icon">🎓</span>
      <span class="brand-text">智学助手</span>
    </div>

    <nav class="nav-items">
      <button
        v-for="item in navItems"
        :key="item.path"
        class="nav-item"
        :class="{ active: activePath.startsWith(item.path) }"
        @click="go(item.path)"
      >
        {{ item.label }}
      </button>
    </nav>

    <div class="user-area">
      <span v-if="auth.user" class="user-name">{{ auth.user.name }}</span>
      <button class="logout-btn" @click="onLogout">退出</button>
    </div>
  </header>
</template>

<style scoped>
.top-nav {
  display: flex;
  align-items: center;
  height: 60px;
  padding: 0 var(--space-6);
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
}

.brand {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-h3);
  font-weight: var(--fw-semibold);
  color: var(--color-primary);
  cursor: pointer;
  user-select: none;
}

.brand-icon {
  font-size: 24px;
}

.nav-items {
  display: flex;
  margin-left: var(--space-12);
  flex: 1;
  gap: var(--space-2);
}

.nav-item {
  padding: var(--space-2) var(--space-4);
  font-size: var(--font-body);
  color: var(--color-text-regular);
  border-radius: var(--radius-md);
  transition: all 0.15s;
}

.nav-item:hover {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.nav-item.active {
  color: var(--color-primary);
  background: var(--color-primary-light);
  font-weight: var(--fw-medium);
}

.user-area {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.user-name {
  font-size: var(--font-body);
  color: var(--color-text-regular);
}

.logout-btn {
  padding: var(--space-2) var(--space-3);
  font-size: var(--font-caption);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all 0.15s;
}

.logout-btn:hover {
  color: var(--color-error);
  border-color: var(--color-error);
}

@media (max-width: 768px) {
  .top-nav { display: none; }
}
</style>