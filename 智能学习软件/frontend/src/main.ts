// 应用入口（Phase B）
//
// 启动顺序：
//   1. 创建 Vue 应用
//   2. 注册 Pinia（状态管理）
//   3. 注册 Vue Router（路由 + 鉴权守卫）
//   4. 注册 Element Plus（按需 unplugin 自动注入，无需手动 use）
//   5. 挂载到 #app

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/styles/reset.css'
import './assets/styles/tokens.css'
import './assets/styles/global.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')