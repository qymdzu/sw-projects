# UI 组件规范 — 翠花集群管理面板 (sw-dashboard)

> 版本：v1.0.0
> 日期：2026-06-03
> 对应阶段：Stage 3 — UI/UX 设计

---

## 1. 页面清单

| 序号 | 页面 | 路由 | 组件层级 | 说明 |
|:-----|:-----|:-----|:---------|:-----|
| 1 | **LoginPage** | `/login` | 独立页面（无 Layout） | Token 认证登录 |
| 2 | **DashboardPage** | `/dashboard` | Layout → 内容区 | 仪表盘首页 |
| 3 | **StoragePage** | `/storage` | Layout → 内容区 | 存储结构浏览器 |
| 4 | **ConfigsPage** | `/configs` | Layout → 内容区 | 配置管理 |
| 5 | **LogsPage** | `/logs` | Layout → 内容区 | 日志查看 |
| 6 | **SkillsPage** | `/skills` | Layout → 内容区 | 技能浏览器 |

### 1.1 通用组件

| 组件名 | 用途 | 复用范围 |
|:-------|:-----|:---------|
| `AppLayout` | 整体布局（导航 + 顶栏 + 内容） | 页面 2~6 |
| `NavSidebar` | 左侧导航菜单 | 页面 2~6 |
| `StatusBadge` | 状态指示圆点（绿/红/黄/灰） | 页面 2, 4, 5 |
| `FileTree` | 目录树（基于 el-tree） | 页面 3, 6 |
| `CodeEditor` | CodeMirror 编辑器封装 | 页面 3, 4 |
| `ConfirmDialog` | 二次确认弹窗 | 页面 3, 4 |
| `Breadcrumb` | 路径面包屑导航 | 页面 3, 6 |
| `SkeletonLoader` | 骨架屏加载占位 | 页面 2, 3, 4, 5, 6 |
| `ErrorState` | 错误状态占位 + 重试按钮 | 页面 2, 3, 4, 5, 6 |
| `EmptyState` | 空状态占位 + 引导文字 | 页面 3, 5, 6 |
| `PageHeader` | 页面标题 + 操作按钮 | 页面 2~6 |
| `AutoRefreshToggle` | 自动刷新开关 | 页面 2, 5 |
| `SearchBar` | 搜索输入框 | 页面 5 |

---

## 2. LoginPage 组件规范

### 2.1 组件结构

```
LoginPage
└── el-container (居中布局)
    ├── LoginCard
    │   ├── Logo / 标题 "翠花集群管理面板"
    │   ├── el-input (Token 输入框, type=password)
    │   ├── el-button "登录" (type=primary, :loading)
    │   └── 错误文字提示
    └── Footer "通过 Tailscale 内网安全连接"
```

### 2.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **空** | 首次加载 | Token 输入框为空，登录按钮禁用 |
| **输入中** | 用户输入 Token | 登录按钮随输入变为可用 |
| **加载中** | 点击登录，等待响应 | 按钮显示 loading 动画，输入框禁用 |
| **正常** | 登录成功 | 跳转 `/dashboard` |
| **异常** | Token 无效 (401) | 输入框变红，显示 "Token 无效" 错误提示 |
| **异常** | 网络错误 | 显示 "网络连接失败，请检查网络" |
| **边界** | Token 过长 (>1024) | 输入时自动截断 |

### 2.3 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 用户在输入框回车 | 触发登录（如非空） |
| 点击登录按钮 | 显示 loading 状态，禁用重复点击 |
| 登录失败 | 按钮恢复，输入框保留输入内容，显示错误 |
| 按 Tab 键 | 聚焦下一个元素（输入框 → 按钮） |

---

## 3. AppLayout 组件规范

### 3.1 组件结构

```
AppLayout
├── el-container (方向: horizontal)
│   ├── NavSidebar (左侧导航)
│   │   ├── Logo / 应用名 "sw-dashboard"
│   │   ├── el-menu (导航菜单)
│   │   │   ├── MenuItem "仪表盘"   → icon: Monitor    → /dashboard
│   │   │   ├── MenuItem "存储"     → icon: FolderOpened → /storage
│   │   │   ├── MenuItem "配置"     → icon: Setting     → /configs
│   │   │   ├── MenuItem "日志"     → icon: Document    → /logs
│   │   │   └── MenuItem "技能"     → icon: MagicStick  → /skills
│   │   └── Divider + "退出登录" 按钮
│   │
│   └── el-container (方向: vertical)
│       ├── TopBar (顶部状态条)
│       │   ├── 页面标题 (根据路由动态) 
│       │   ├── 服务器时间 (实时时钟)
│       │   ├── Token 状态指示 (绿色圆点)
│       │   └── 折叠/展开 导航按钮 (响应式)
│       │
│       └── ContentArea (路由视图: <router-view />)
│               └── 各页面内容
```

### 3.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **正常** | Token 有效，路由匹配 | 导航高亮当前页，内容区显示对应页面 |
| **空** | Token 有值但路由不匹配 | 内容区显示 404 占位 |
| **异常** | Token 过期 (401) | 清除 Token，跳转 `/login` |
| **响应式** | 窗口 < 900px | 导航收起为汉堡菜单 |
| **响应式** | 窗口 900~1199px | 导航折叠为图标模式 (64px) |

### 3.3 导航菜单状态

| 状态 | UI 表现 |
|:-----|:---------|
| 默认 | 灰色图标 + 文字 |
| hover | 背景色微变 (var(--el-menu-hover-bg-color)) |
| 选中 | 高亮背景 + 高亮图标颜色 |
| 禁用 | 无（所有菜单项均可点击） |

### 3.4 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 点击菜单项 | 路由切换，内容区淡入动画 |
| 当前页面菜单项 | 高亮 + 不可点击 |
| 点击退出登录 | `ElMessageBox.confirm("确认退出登录？")` → 清除 Token → 跳转 `/login` |
| 鼠标悬停菜单 | 背景色微变，150ms 过渡 |
| 折叠按钮点击 | 导航宽度动画切换 (240px ↔ 64px) |

---

## 4. DashboardPage 组件规范

### 4.1 组件结构

```
DashboardPage
├── PageHeader
│   ├── 标题 "仪表盘"
│   └── 操作栏
│       ├── 刷新按钮 (el-button icon=Refresh)
│       └── 自动刷新开关 (el-switch, 默认开启, 30s)
│
├── StatsGrid (卡片网格, el-row :gutter="20")
│   ├── GatewayCard (el-card)
│   │   ├── 卡片标题 "Gateway 状态"
│   │   ├── StatusBadge (绿色/红色圆点 + "运行中"/"已停止")
│   │   ├── InfoRows: Uptime / PID / Memory / CPU
│   │   └── 卡片尺寸: 1/3 行宽 (col-span-4)
│   │
│   ├── SessionCard (el-card)
│   │   ├── 卡片标题 "会话统计"
│   │   ├── 统计数字展示 (大号字体): 总会话数 / 活跃会话
│   │   ├── 辅助信息: 总消息数 / 平均时长
│   │   └── 卡片尺寸: 1/3 行宽
│   │
│   └── CronJobsCard (el-card)
│       ├── 卡片标题 "定时任务"
│       ├── 任务列表 (el-table 精简模式)
│       │   ├── 列: 名称 / 调度 / 上次运行 / 状态 / 启停开关
│       │   └── 最大显示 5 条
│       └── 卡片尺寸: 1/3 行宽
│
└── PipelineCard (el-card, 全行宽)
    ├── 卡片标题 "Pipeline 进度"
    ├── 整体进度条 (el-progress :percentage="overall_progress * 100")
    ├── StageList (阶段列表, 水平排列)
    │   └── StageItem × N
    │       ├── 状态图标: ✓ completed / ⟳ running / ○ pending / ✗ failed
    │       ├── 阶段名称
    │       └── (可选) 阶段进度条
    └── 当前阶段高亮 + "执行中" 标签
```

### 4.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **加载中** | 页面首次加载 / 手动刷新 | 3 张骨架卡片 + 进度条骨架 |
| **正常** | 数据返回 | 渲染所有卡片 |
| **部分异常** | 某个子系统无响应 | 对应卡片显示"数据不可用" + 重试按钮，其他卡片正常 |
| **全部异常** | 所有数据不可用 | 整体 ErrorState 组件显示 + 重试按钮 |
| **空** | 无定时任务 | CronJobsCard 显示 "暂无定时任务" (EmptyState) |
| **空** | 无 Pipeline 数据 | PipelineCard 显示 "暂无 Pipeline 任务" |
| **超时** | API > 5s | 降级显示部分已获取数据，超时部分显示 "加载超时" |
| **离线** | 网络断开 | ErrorState "无法连接到服务器" + 重试按钮 |

### 4.3 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 手动点击刷新 | 按钮 loading 旋转，3 张卡片单独重新加载 |
| 自动刷新 | 静默更新，无视觉闪烁（数据稳定时） |
| Pipeline 阶段悬停 | tooltip 显示阶段详情 |
| Cron 启停开关切换 | `ElMessageBox.confirm` → `POST /api/configs/cron/toggle` |
| 页面隐藏 | 暂停自动刷新 (Page Visibility API) |

### 4.4 Pipeline 阶段颜色映射

| 状态 | 颜色 | 图标 |
|:-----|:-----|:-----|
| `completed` | `#67C23A` (绿色) | ✓ |
| `running` | `#409EFF` (蓝色) | ⟳ (旋转动画) |
| `pending` | `#909399` (灰色) | ○ |
| `failed` | `#F56C6C` (红色) | ✗ |

---

## 5. StoragePage 组件规范

### 5.1 组件结构

```
StoragePage
├── PageHeader
│   ├── 标题 "存储结构浏览器"
│   └── 仓库切换 Tabs (el-tabs)
│       ├── Tab "全部" (repo=all)
│       ├── Tab "Workspace"
│       ├── Tab "Skills Library"
│       └── Tab "Projects"
│
├── el-container (方向: horizontal, 高度: calc(100vh - 顶部 - header))
│   ├── FileTreePanel (左侧目录树, 可拖拽宽度)
│   │   ├── 标题 "目录树" + 刷新图标
│   │   └── el-tree (懒加载, :load="loadNode")
│   │       └── TreeNode
│   │           ├── 文件夹图标 / 文件图标 (根据 type)
│   │           └── 文件名
│   │
│   └── FileContentPanel (右侧文件内容)
│       ├── Breadcrumb (路径导航)
│       ├── FileInfoBar
│       │   ├── 文件名
│       │   ├── 文件大小 (字节格式化: B/KB/MB)
│       │   └── 修改时间
│       ├── ActionBar
│       │   ├── 编辑按钮 / 取消编辑按钮
│       │   └── 保存按钮 (仅编辑模式)
│       ├── CodeEditor (vue-codemirror)
│       │   ├── readonly 模式 (默认)
│       │   └── writable 模式 (点击编辑后)
│       └── 文件未选中占位 (EmptyState)
```

### 5.2 目录树节点状态

| 状态 | UI 表现 |
|:-----|:---------|
| **未展开** | 文件夹图标 + 名称，箭头向右 |
| **展开中（加载）** | 箭头旋转动画 |
| **已展开** | 文件夹图标打开 + 子节点列表，箭头向下 |
| **选中** | 背景色高亮 |
| **叶子（文件）** | 文件图标 (根据扩展名区分) |

### 5.3 文件内容区状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **空** | 未选择文件 | 居中 EmptyState "选择文件以预览" + 文件图标 |
| **加载中** | 正在获取文件内容 | CodeMirror 位置显示 SkeletonLoader |
| **正常（只读）** | 内容加载完成 | CodeMirror 渲染，readonly=true，显示"编辑"按钮 |
| **正常（编辑）** | 用户点击编辑 | CodeMirror readonly=false，显示"保存"/"取消" |
| **异常（404）** | 文件不存在 | 内容区 ErrorState "文件不存在或已被删除" |
| **异常（403）** | 权限不足 | ErrorState "无权限访问此文件" |
| **异常（413）** | 文件过大 | ErrorState "文件过大，无法预览（最大 10MB）" |
| **异常（500）** | 服务器错误 | ErrorState "读取文件失败" + 重试按钮 |
| **保存成功** | PUT 成功 | ElMessage.success + 切换回只读模式 |
| **保存失败(400)** | 格式错误 | ElMessage.error + 保持编辑模式 |
| **保存失败(500)** | 写入失败 | ElMessage.error + 保持编辑模式 |

### 5.4 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 点击目录展开箭头 | 懒加载子节点，箭头旋转 |
| 点击文件节点 | 高亮 + 加载文件内容 |
| 悬停文件节点 | tooltip 显示完整路径 |
| 点击编辑按钮 | CodeMirror 切换编辑模式 |
| 内容修改 | 保存按钮变为可用 (el-button :disabled=false) |
| 点击保存 | 二次确认弹窗 → PUT API |
| 点击取消 | 二次确认 "放弃修改？" → 重新加载原始内容 |
| 切换仓库 Tab | 重新加载树（当前选中状态清除） |
| 切换文件 | 如有未保存修改，先提示确认 |

### 5.5 文件图标映射

| 文件扩展名 | 图标 (Element Plus) | 颜色 |
|:-----------|:--------------------|:-----|
| `.md` | `Document` | 蓝色 |
| `.yaml` / `.yml` | `Setting` | 橙色 |
| `.py` | `Code` | 绿色 |
| `.json` | `DataBoard` | 紫色 |
| `.toml` | `Setting` | 橙色 |
| `.env` | `Key` | 红色 |
| `.sh` | `Terminal` | 灰色 |
| `.log` | `Tickets` | 棕色 |
| 通用文本 | `Document` | 默认 |
| 二进制文件 | `Lock` | 灰色（不可预览） |

---

## 6. ConfigsPage 组件规范

### 6.1 组件结构

```
ConfigsPage
├── PageHeader
│   └── 标题 "配置管理"
│
├── el-tabs (标签页, :model-value="activeTab")
│   ├── TabPane "config.yaml"  → icon=Setting, 语言=yaml
│   ├── TabPane ".env"         → icon=Key, 语言=env
│   ├── TabPane "CLAUDE.md"   → icon=Document, 语言=markdown
│   ├── TabPane "SOUL.md"     → icon=Document, 语言=markdown
│   └── TabPane "Cron 任务"   → icon=Timer
│
├── ConfigEditor (Tab 内容: config.yaml / .env / CLAUDE.md / SOUL.md)
│   ├── InfoBar
│   │   ├── 配置文件名 + 路径
│   │   ├── 最后修改时间
│   │   └── (仅 .env) "显示原始值" 按钮
│   ├── CodeEditor (vue-codemirror, always editable)
│   │   └── 语法高亮根据 language 字段
│   ├── SaveBar
│   │   ├── 保存按钮 (Ctrl+S 快捷键)
│   │   └── 提示: 未保存修改数量指示
│   └── (仅 config.yaml) "重启服务" 提示框
│
└── CronEditor (Tab 内容: Cron 任务)
    └── el-table
        ├── 列1: 任务名称 (从 command 或 comment 提取)
        ├── 列2: 调度表达式 (人类可读)
        ├── 列3: 上次执行时间
        ├── 列4: 执行状态 (success/failed)
        ├── 列5: 启用/禁用开关 (el-switch)
        └── (每行) 展开可查看完整命令和输出
```

### 6.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **加载中** | 切换 Tab | CodeMirror 位置显示 SkeletonLoader |
| **正常（未修改）** | 内容加载完成 | CodeMirror 可编辑，保存按钮禁用 |
| **正常（已修改）** | 内容发生变化 | 保存按钮可用，Tab 标签显示小红点 |
| **正常（.env 掩码）** | 默认加载 | Content 显示掩码版本，显示"显示原始值"按钮 |
| **正常（.env 明文）** | 点击显示原始值 | 弹窗确认后显示明文 |
| **保存成功** | PUT 成功 | `ElMessage.success` |
| **保存失败(400)** | YAML 格式错误 | `ElMessage.error` + 错误信息 |
| **保存失败(500)** | 写入失败 | `ElMessage.error("保存失败")` |
| **空（Cron）** | 无 Cron 任务 | Table 显示 EmptyState "暂无定时任务" |

### 6.3 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 切换 Tab | 检测未保存修改 → 确认弹窗 → 加载新配置 |
| 编辑内容 | 保存按钮立即可用，Tab 标签出现小红点 |
| Ctrl+S / 保存按钮 | `ElMessage.info("正在保存...")` → PUT → 成功/失败反馈 |
| .env 显示原始值 | 弹窗警告 "显示原始值会暴露敏感信息" → 确认后 `?masked=false` |
| Tab 小红点 | 内容已修改未保存 |
| Cron 启停切换 | 无确认弹窗，直接 toggle，显示加载状态 |

---

## 7. LogsPage 组件规范

### 7.1 组件结构

```
LogsPage
├── PageHeader
│   ├── 标题 "日志查看"
│   └── 操作栏
│       ├── 日志文件选择器 (el-select)
│       │   └── 通过 fetchLogFiles 填充
│       ├── 刷新按钮 (手动刷新)
│       └── 自动滚动开关 (默认开启)
│
├── SearchBar (搜索区域)
│   ├── el-input (关键字搜索, v-model="keyword")
│   ├── el-button "搜索" (type=primary)
│   └── 搜索结果指示 "共找到 N 条匹配 | 清空搜索"
│
├── LogContent (日志内容区, 等高性能: flex:1 overflow:auto)
│   ├── (搜索模式) 匹配行高亮显示 + 上下文行
│   ├── (tail 模式) 实时增量显示
│   ├── 自动滚动到底部 (跟随新日志)
│   └── 暂停指示条 "新日志已暂停 · 回到底部"
│
└── el-collapse (折叠面板: 历史归档)
    └── ArchiveList
        └── el-table (精简)
            ├── 列1: 归档文件名
            ├── 列2: 大小
            ├── 列3: 修改时间
            └── (点击行) 只读加载内容
```

### 7.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **空（Tail）** | 首次加载 | 日志区 EmptyState "暂无日志内容" |
| **加载中（Tail）** | 首次请求 | SkeletonLoader |
| **正常（Tail）** | 数据返回 | 显示日志行，自动滚动到底部 |
| **暂停滚动** | 用户手动上滚 | 顶部显示暂停指示条 "新日志已暂停 · 回到底部" |
| **搜索中** | 点击搜索 | 搜索按钮 loading，日志区显示等待 |
| **搜索正常** | 搜索结果返回 | 显示匹配行 + 高亮关键字 + 总数 |
| **搜索空** | 无匹配 | 日志区 EmptyState "未找到匹配 '{keyword}' 的内容" |
| **搜索异常** | 搜索失败 | `ElMessage.error` |
| **轮询失败** | 网络抖动 | 静默重试，显示上次成功内容 |
| **空（归档）** | 无历史归档 | 折叠面板 EmptyState "暂无历史归档" |
| **查看归档** | 点击归档文件 | 以只读模式显示归档内容 |

### 7.3 日志行样式

| 日志级别 | 前缀颜色 | 示例 |
|:---------|:---------|:------|
| INFO | 灰色 (#909399) | `[2026-06-03 14:29:00] INFO: Session started` |
| WARNING | 橙色 (#E6A23C) | `[2026-06-03 14:29:00] WARNING: High memory usage` |
| ERROR | 红色 (#F56C6C) | `[2026-06-03 14:29:00] ERROR: Connection timeout` |
| DEBUG | 蓝色 (#409EFF) | `[2026-06-03 14:29:00] DEBUG: Loaded config` |
| CRITICAL | 深红 (#C03639) | `[2026-06-03 14:29:00] CRITICAL: Service down` |

### 7.4 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 切换日志文件 | 清空当前内容，重新开始 tail 轮询 |
| 输入关键字 + 回车 | 触发搜索 |
| 搜索完成点击 "清空" | 回到 tail 实时模式 |
| 用户手动滚动到顶部以外的位置 | 暂停自动滚动，显示暂停指示条 |
| 点击 "回到底部" | 恢复自动滚动 |
| 点击归档行 | 以只读模式加载归档内容（停止 tail） |

---

## 8. SkillsPage 组件规范

### 8.1 组件结构

```
SkillsPage
├── PageHeader
│   └── 标题 "技能浏览器"
│
├── el-container (方向: horizontal, 高度: calc(100vh - 顶部 - header))
│   ├── SkillTreePanel (左侧技能树)
│   │   ├── 标题 "已安装技能" + 刷新图标
│   │   └── el-tree (懒加载)
│   │       └── TreeNode
│   │           ├── 技能图标 (MagicStick 或 文件夹图标)
│   │           ├── 文件名
│   │           └── (目录节点) 如有 SKILL.md 显示特殊标记
│   │
│   └── SkillPreviewPanel (右侧预览区)
│       ├── Breadcrumb
│       ├── FileInfoBar
│       │   ├── 文件名
│       │   ├── 大小
│       │   └── 修改时间
│       ├── MarkdownRenderer (Markdown 渲染, 非 CodeMirror)
│       │   └── 渲染 SKILL.md 内容为格式化文档
│       └── 未选中占位 (EmptyState)
```

### 8.2 状态矩阵

| 状态 | 触发条件 | UI 表现 |
|:-----|:---------|:--------|
| **空** | 未选择文件 | EmptyState "选择一个技能查看详情" |
| **加载中** | 点击文件 | 右侧内容区 SkeletonLoader |
| **正常** | 内容加载完成 | Markdown 渲染展示 |
| **空（无技能）** | 技能目录为空 | 树区域 EmptyState "暂无已安装的技能" |
| **异常** | 文件读取失败 | 内容区 ErrorState + 重试按钮 |

### 8.3 交互行为

| 操作 | 反馈 |
|:-----|:------|
| 展开目录 | 懒加载子节点 |
| 点击 SKILL.md | 渲染为 Markdown 文档（非源码） |
| 点击非 SKILL.md 文件 | 显示源码内容 |
| 悬停 has_skill_md 目录 | tooltip "包含 SKILL.md" |

---

## 9. 通用组件规范

### 9.1 SkeletonLoader

| 属性 | 说明 |
|:-----|:------|
| **用途** | 页面/组件首次加载占位 |
| **实现** | Element Plus `el-skeleton` |
| **行数** | 根据内容区域高度动态 (3~8 行) |
| **动画** | 波纹闪烁 (animated) |
| **变体** | 文本骨架 / 卡片骨架 / 表格骨架 |

### 9.2 ErrorState

| 属性 | 说明 |
|:-----|:------|
| **用途** | API 请求失败时展示 |
| **结构** | 错误图标 + 错误文字 + 重试按钮 |
| **文字** | 默认 "加载失败" + 错误详情 |
| **操作** | "重试" 按钮 → 重新发起请求 |
| **变体** | 无权限 (403) / 不存在 (404) / 服务错误 (500) / 网络断开 |

### 9.3 EmptyState

| 属性 | 说明 |
|:-----|:------|
| **用途** | 数据为空时引导用户 |
| **结构** | 大图标 (Element Plus icon) + 描述文字 |
| **示例** | "选择文件以预览" / "暂无日志内容" / "暂无定时任务" |
| **样式** | 居中，半透明图标 (opacity: 0.4) |

### 9.4 StatusBadge

| 属性 | 说明 |
|:-----|:------|
| **用途** | 显示系统/组件的运行状态 |
| **形状** | 8px 圆点 + 文字 |
| **颜色映射** | 绿色=正常/活跃, 红色=停止/异常, 黄色=警告, 灰色=离线/禁用 |
| **动画** | 活跃状态圆点有呼吸动画 (pulse) |

### 9.5 ConfirmDialog

| 属性 | 说明 |
|:-----|:------|
| **用途** | 保存/删除/退出等危险操作确认 |
| **实现** | `ElMessageBox.confirm` |
| **标题** | "确认操作" |
| **内容** | 描述具体操作和后果 |
| **按钮** | "取消" (次要) + "确认" (主要/危险) |
| **取消文案** | "放弃未保存的修改？" → 确认/取消 |
| **保存文案** | "确认保存对 {filename} 的修改？" → 确认/取消 |

### 9.6 Breadcrumb

| 属性 | 说明 |
|:-----|:------|
| **用途** | 显示当前文件/页面的路径导航 |
| **实现** | Element Plus `el-breadcrumb` |
| **分隔符** | `>` |
| **最大深度** | 5 层，超出用 `...` 截断 |
| **点击** | 可点击路径段跳转到对应目录 |

---

## 10. 全局交互规范

### 10.1 加载时序

```
用户操作 → 即时反馈(100ms内) → API 请求 → 结果反馈
    │                           │            │
    └─ 按钮loading/骨架屏       └─ <300ms    └─ 成功: 更新内容
                                        可跳过loading   失败: 错误提示
```

### 10.2 Toast 消息规范

| 类型 | 触发场景 | 展示位置 | 自动关闭 |
|:-----|:---------|:---------|:---------|
| `success` | 保存成功 / 操作成功 | 右上角 | 3s |
| `warning` | 即将超时 / 文件过大 | 右上角 | 4s |
| `error` | API 错误 / 保存失败 | 右上角 | 5s |
| `info` | 正在保存 / 操作提示 | 右上角 | 3s |

### 10.3 键盘导航

| 场景 | 行为 |
|:-----|:------|
| 弹窗打开 | 自动聚焦「确认」按钮 |
| 弹窗关闭 | 恢复焦点到触发元素 |
| Tab 键顺序 | 从左到右，从上到下 |
| Enter 键 | 触发默认操作（提交/确认） |
| Escape 键 | 关闭弹窗/取消编辑 |

### 10.4 无障碍

| 要求 | 措施 |
|:-----|:------|
| 语义化标签 | 使用正确的 HTML5 语义元素 |
| ARIA 标签 | 图标按钮添加 `aria-label` |
| 键盘可操作 | 所有功能均可通过键盘完成 |
| 对比度 | 文本/背景对比度 ≥ 4.5:1 |
| 字号 | 最小字号 12px，支持浏览器缩放 |
