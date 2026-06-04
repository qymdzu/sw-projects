import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useTokenStore } from '@/stores/token'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { title: '仪表盘', icon: 'Monitor' }
  },
  {
    path: '/storage',
    name: 'Storage',
    component: () => import('@/views/Storage.vue'),
    meta: { title: '存储浏览', icon: 'FolderOpened' }
  },
  {
    path: '/configs',
    name: 'Configs',
    component: () => import('@/views/Configs.vue'),
    meta: { title: '配置管理', icon: 'Setting' }
  },
  {
    path: '/logs',
    name: 'Logs',
    component: () => import('@/views/Logs.vue'),
    meta: { title: '日志查看', icon: 'Document' }
  },
  {
    path: '/skills',
    name: 'Skills',
    component: () => import('@/views/Skills.vue'),
    meta: { title: '技能浏览', icon: 'MagicStick' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '404', noAuth: true }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const tokenStore = useTokenStore()

  if (to.meta.noAuth) {
    if (tokenStore.isAuthenticated && to.name === 'Login') {
      next({ name: 'Dashboard' })
    } else {
      next()
    }
    return
  }

  if (!tokenStore.isAuthenticated) {
    next({ name: 'Login' })
    return
  }

  document.title = `${to.meta.title || ''} - sw-dashboard`
  next()
})

export default router
