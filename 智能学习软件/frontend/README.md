# 智学助手前端

> Phase B 新建 — Vue 3 + Vite + TypeScript + Pinia + Vue Router 4 + Axios

## 开发

```bash
npm install
npm run dev      # 启动 Vite (默认 http://localhost:5173)
```

后端 API 走 Vite 代理转发 `/api/v1` → `http://localhost:8080`。

## 构建

```bash
npm run build    # vue-tsc + vite build → dist/
npm run preview  # 预览构建产物
```

## 测试

```bash
npm run test         # 一次性跑（CI 用）
npm run test:watch   # watch 模式
```

## 技术栈

- **Vue 3** + Composition API + `<script setup>`
- **TypeScript** strict 模式
- **Pinia** 状态管理
- **Vue Router 4** 路由 + 鉴权守卫
- **Axios** HTTP 客户端 + 拦截器（自动 refresh）
- **Element Plus** PC 端 UI（按需 unplugin）
- **Vitest** + happy-dom 单元测试

## 目录结构

```
src/
├── api/           # Axios 客户端 + 各模块 API 封装
├── assets/styles/ # 全局样式（reset / tokens / global）
├── components/    # 通用组件 + 业务组件
│   ├── common/    # 基础原子组件
│   ├── nav/       # 导航
│   └── exercise/  # 题目卡片
├── layouts/       # 路由布局（Blank / Main）
├── router/        # 路由表 + 守卫
├── stores/        # Pinia store（auth / setting）
├── types/         # 类型声明
├── utils/         # 工具函数（storage / format）
└── views/         # 页面级组件（7 个）
```