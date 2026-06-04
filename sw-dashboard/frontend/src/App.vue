<template>
  <el-container style="height: 100vh">
    <!-- 登录页不显示侧边栏 -->
    <template v-if="route.name !== 'Login'">
      <el-aside :width="appStore.sidebarCollapsed ? '64px' : '220px'">
        <div class="sidebar-logo" @click="router.push('/dashboard')" style="cursor: pointer">
          <span style="font-size: 20px">🌼</span>
          <span v-show="!appStore.sidebarCollapsed">sw-dashboard</span>
        </div>
        <el-menu
          :default-active="route.path"
          :collapse="appStore.sidebarCollapsed"
          background-color="#191a1b"
          text-color="#a0a0a0"
          active-text-color="#409eff"
          router
        >
          <el-menu-item index="/dashboard">
            <el-icon><Monitor /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/storage">
            <el-icon><FolderOpened /></el-icon>
            <template #title>存储浏览</template>
          </el-menu-item>
          <el-menu-item index="/configs">
            <el-icon><Setting /></el-icon>
            <template #title>配置管理</template>
          </el-menu-item>
          <el-menu-item index="/logs">
            <el-icon><Document /></el-icon>
            <template #title>日志查看</template>
          </el-menu-item>
          <el-menu-item index="/skills">
            <el-icon><MagicStick /></el-icon>
            <template #title>技能浏览</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
    </template>
    <el-container direction="vertical">
      <template v-if="route.name !== 'Login'">
        <el-header>
          <div style="display: flex; align-items: center; gap: 12px">
            <el-button text style="color: #a0a0a0; font-size: 18px" @click="appStore.toggleSidebar">
              <el-icon><Fold v-if="!appStore.sidebarCollapsed" /><Expand v-else /></el-icon>
            </el-button>
            <span style="font-size: 15px; font-weight: 500">{{ route.meta?.title }}</span>
          </div>
          <div style="display: flex; align-items: center; gap: 10px; margin-left: auto">
            <StatusBadge :status="tokenStore.isAuthenticated ? 'active' : 'inactive'" />
            <span style="font-size: 12px; color: var(--text-secondary)">
              {{ tokenStore.isAuthenticated ? '已认证' : '未认证' }}
            </span>
            <el-button text size="small" style="color: var(--text-secondary)" @click="handleLogout">
              退出
            </el-button>
          </div>
        </el-header>
      </template>
      <el-main>
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useTokenStore } from '@/stores/token'
import { useAppStore } from '@/stores/app'
import StatusBadge from '@/components/StatusBadge.vue'
import {
  Monitor, FolderOpened, Setting, Document, MagicStick,
  Fold, Expand
} from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const tokenStore = useTokenStore()
const appStore = useAppStore()

function handleLogout() {
  ElMessageBox.confirm('确认退出登录？', '提示', {
    confirmButtonText: '退出',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    tokenStore.clearToken()
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.el-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
