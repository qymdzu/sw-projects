// Vue Router 配置 — Phase B
//
// 路由结构：
//   /login                       → BlankLayout（公开）
//   /404                         → BlankLayout（公开）
//   /                            → MainLayout（鉴权后）
//   ├─ /dashboard                → 仪表盘
//   ├─ /subjects                 → 科目浏览
//   ├─ /exercise                 → 智能选题
//   ├─ /mistakes                 → 错题本
//   └─ /settings                 → 设置（含模型配置）
//   /:pathMatch(.*)*             → 重定向到 /404

import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('@/views/LoginView.vue'),
    meta: { layout: 'blank', public: true, title: '登录' }
  },
  {
    path: '/404',
    component: () => import('@/views/NotFoundView.vue'),
    meta: { layout: 'blank', public: true, title: '404' }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/DashboardView.vue'),
        meta: { title: '仪表盘', icon: 'House' }
      },
      {
        path: 'subjects',
        name: 'subjects',
        component: () => import('@/views/SubjectsView.vue'),
        meta: { title: '科目', icon: 'Reading' }
      },
      {
        path: 'exercise',
        name: 'exercise',
        component: () => import('@/views/ExerciseView.vue'),
        meta: { title: '智能选题', icon: 'EditPen' }
      },
      {
        path: 'mistakes',
        name: 'mistakes',
        component: () => import('@/views/MistakesView.vue'),
        meta: { title: '错题本', icon: 'Warning' }
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/SettingsView.vue'),
        meta: { title: '设置', icon: 'Setting' }
      }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/404' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 鉴权守卫：未登录访问受保护路由 → /login?redirect=...
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.isAuthenticated) {
    return { path: '/dashboard' }
  }
  return true
})

export default router