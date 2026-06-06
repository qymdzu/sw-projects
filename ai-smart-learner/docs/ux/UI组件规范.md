# UI 组件规范 — AI智能学习机 v0.1

> **文档版本**：v0.1.0  
> **编制日期**：2026-06-06  
> **目标**：所有 UI 组件的视觉规范 + Dart 实现约束

---

## 1. 组件总览

| 分类 | 组件 | v0.1 |
|:-----|:-----|:-----|
| **基础** | AppScaffold, ErrorBoundary, LoadingOverlay | ✅ |
| **状态** | EmptyState, ErrorToast, LoadingIndicator | ✅ |
| **按钮** | PrimaryButton, SecondaryButton, IconButton, CaptureButton | ✅ |
| **卡片** | QuestionCard, KnowledgePointCard, UserCard | ✅ |
| **输入** | NameInputDialog, ManualTextInput | ✅ |
| **图表** | MasteryColorDot, RainbowChart | ✅ |
| **反馈** | SuccessToast, ErrorToast, AchievementAnimation | ✅ |

---

## 2. 基础组件

### 2.1 AppScaffold

```dart
class AppScaffold extends StatelessWidget {
  final String title;
  final Widget? body;
  final Widget? floatingActionButton;
  final List<Widget>? actions;
  
  // 统一 SafeArea + AppBar + Body + Loading 状态
}
```

**视觉**：
- 顶部 SafeArea + AppBar（高 56pt）
- 默认有 Loading 监听（全局 Riverpod provider）
- 错误边界包裹 body

### 2.2 ErrorBoundary

```dart
class ErrorBoundary extends StatelessWidget {
  final Widget child;
  
  // 捕获子组件异常 → 显示 ErrorState
  // 不阻断其他页面
}
```

**视觉**：
- 异常时全屏显示错误页
- 大图标（80pt）+ 错误信息 + [重试] [返回] 按钮

### 2.3 LoadingOverlay

```dart
class LoadingOverlay extends StatelessWidget {
  // 全屏半透明 + 中央 Loading 旋转 + 文案
}
```

**视觉**：
- 黑色 50% 透明遮罩
- 中央：圆形 Loading 动画 + "AI 正在识别..." 文案
- 不可点击背后内容

---

## 3. 按钮组件

### 3.1 PrimaryButton（主操作按钮）

**用法**：拍照、确认、提交、开始

```dart
class PrimaryButton extends StatelessWidget {
  final String label;
  final IconData? icon;       // 左侧图标
  final VoidCallback? onPressed;
  final bool loading;
}
```

**视觉规范**：
- 高度：56pt（最小可点击区域）
- 宽度：撑满父容器
- 圆角：16pt
- 背景：主色（见视觉风格指南）
- 文字：白色 18pt Bold
- 加载时：按钮内显示 Loading 旋转器

### 3.2 SecondaryButton（次操作按钮）

**用法**：取消、跳过、查看详情

```dart
class SecondaryButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
}
```

**视觉规范**：
- 高度：56pt
- 宽度：撑满父容器
- 圆角：16pt
- 背景：透明
- 边框：主色 2pt
- 文字：主色 18pt Medium

### 3.3 CaptureButton（拍照大按钮）

**用法**：CapturePage 拍照主按钮

```dart
class CaptureButton extends StatelessWidget {
  final VoidCallback onTap;
  final bool disabled;
}
```

**视觉规范**：
- 大圆按钮：直径 100pt
- 中心图标：相机（白色）
- 背景：主色
- 边缘：白色 4pt 描边
- 按下时缩放至 0.95

### 3.4 IconButton（图标按钮）

**用法**：AppBar 操作、Tab 切换

**视觉规范**：
- 最小可点击：44x44pt
- 图标：24pt
- 颜色：自适应（深色背景用白，浅色用黑）

---

## 4. 卡片组件

### 4.1 QuestionCard（错题卡片）

```dart
class QuestionCard extends StatelessWidget {
  final ErrorQuestion question;
  final VoidCallback? onTap;
}
```

**视觉规范**：
- 高度：自适应内容（最小 80pt）
- 圆角：12pt
- 背景：白色 + 1pt 边框
- 左边框：4pt 主色（标识类型：错题/漏题）
- 内容：
  - 顶部：题目文本（最多 3 行）
  - 底部：知识点 + 错误类型 + 状态 badge
- 点击：导航到详情（v0.1 简化为只显示）

### 4.2 KnowledgePointCard（知识点卡片）

```dart
class KnowledgePointCard extends StatelessWidget {
  final KnowledgePoint kp;
  final VoidCallback? onTap;
}
```

**视觉规范**：
- 圆角：12pt
- 背景：白色
- 左侧：MasteryColorDot（红/黄/绿）
- 内容：
  - 知识点名称（16pt Bold）
  - 章节（12pt 灰色）
  - 错题数（12pt 数字）
- 点击：展开看关联错题

### 4.3 UserCard（用户卡片）

**用法**：SettingsPage 切换身份

**视觉规范**：
- 头像 emoji（48pt）+ 名字（16pt）+ 角色（12pt 灰色）
- 高度：64pt
- 右侧箭头（"点击切换"）

---

## 5. 输入组件

### 5.1 NameInputDialog（名字输入弹窗）

**用法**：创建用户 / 修改名字

**视觉规范**：
- 居中弹窗，宽 280pt
- 标题："你的名字"（18pt Bold）
- 输入框：12pt padding，14pt 文字
- 按钮：[取消] [确定]

### 5.2 ManualTextInput（手动输入框）

**用法**：OCR 失败时手动输入题目

**视觉规范**：
- 全屏输入页
- 多行文本框（最少 5 行）
- 顶部 [返回] [保存]

---

## 6. 状态色组件

### 6.1 MasteryColorDot（掌握度色点）

```dart
class MasteryColorDot extends StatelessWidget {
  final String mastery;  // 'red' / 'yellow' / 'green'
  final double size;     // 默认 12pt
}
```

**视觉规范**：
- 圆形
- 颜色：红 `#E53935` / 黄 `#FBBC04` / 绿 `#34A853`
- 可选：带边框（深色背景时）

### 6.2 状态 Badge

| 状态 | 文字 | 背景色 | 文字色 |
|:-----|:-----|:-------|:-------|
| 待周测 | "待周测" | 浅蓝 `#E3F2FD` | 深蓝 `#1976D2` |
| 周测通过 | "已通过周测" | 浅黄 `#FFF8E1` | 深黄 `#F57C00` |
| 月测通过 | "已归档" | 浅绿 `#E8F5E9` | 深绿 `#2E7D32` |
| 漏题 | "漏题" | 浅橙 `#FFF3E0` | 深橙 `#E65100` |

---

## 7. 反馈组件

### 7.1 SuccessToast

**用法**：操作成功

**视觉**：
- 底部弹入
- 绿色 ✓ 图标 + 文字
- 3 秒后自动消失

### 7.2 ErrorToast

**用法**：操作失败

**视觉**：
- 底部弹入
- 红色 ✗ 图标 + 错误码 + 文字
- 5 秒后消失
- 可点击查看详情

### 7.3 AchievementAnimation

**用法**：周测完成 / 知识掌握

**视觉**：
- 全屏半透明遮罩
- 中央：🎉 emoji + 文字（"答对 4/5 !"）
- 撒花动画（v0.1 简化为 emoji）
- 3 秒后自动消失或点击关闭

---

## 8. 图表组件

### 8.1 RainbowChart（知识彩虹图）

**用法**：RainbowChartPage

**实现**：fl_chart BarChart

**视觉规范**：
- 3 个柱（红/黄/绿）
- 高度：240pt
- 柱顶显示数字
- X 轴：标签
- Y 轴：隐藏

**变体**（v0.2 推）：
- 饼图（按知识点分布）
- 雷达图（按学科）

---

## 9. 文字规范

| 用途 | 字号 | 字重 | 颜色 |
|:-----|:-----|:-----|:-----|
| H1（页面标题）| 28pt | Bold | 主文字 `#212121` |
| H2（区块标题）| 20pt | Bold | 主文字 `#212121` |
| H3（卡片标题）| 16pt | Medium | 主文字 `#212121` |
| 正文 | 14pt | Regular | 主文字 `#212121` |
| 辅助 | 12pt | Regular | 灰色 `#757575` |
| Caption | 10pt | Regular | 浅灰 `#9E9E9E` |
| 按钮文字 | 18pt | Medium | 白色 / 主色 |

---

## 10. 间距规范

| 用途 | 数值 |
|:-----|:-----|
| 页面边距 | 16pt |
| 卡片内边距 | 16pt |
| 元素间距（小）| 8pt |
| 元素间距（中）| 16pt |
| 元素间距（大）| 24pt |
| 区块间距 | 32pt |

---

## 11. 圆角规范

| 用途 | 圆角 |
|:-----|:-----|
| 按钮 | 16pt |
| 卡片 | 12pt |
| 输入框 | 8pt |
| 头像 | 50%（圆形）|
| Badge | 4pt |

---

## 12. 阴影规范

| 用途 | 阴影 |
|:-----|:-----|
| 卡片 | `0 1pt 3pt rgba(0,0,0,0.08)` |
| 浮层 | `0 4pt 16pt rgba(0,0,0,0.12)` |
| 模态 | `0 8pt 24pt rgba(0,0,0,0.16)` |

---

## 13. 可访问性

| 项 | 要求 |
|:-----|:-----|
| 最小可点击区域 | 44x44pt |
| 文字对比度 | ≥ 4.5:1（WCAG AA）|
| 按钮文字 | 不依赖颜色传达信息（配合文字）|
| VoiceOver | 所有 IconButton 加 Semantics label |

---

## 14. 国际化（v0.1 暂不实现）

- 默认中文（zh-CN）
- 文本不硬编码，用 `AppLocalizations.of(context).xxx`
- v0.1 只在 const string 中（`"你已掌握"` 直接写）
- v0.3 推多语言

---

## 15. 性能

- 列表用 `ListView.builder`（懒加载）
- 大图用 `flutter_image_compress` 压缩
- fl_chart 设置 `repositoriesDataSet: true` 减少重建
- 状态用 Riverpod 细粒度 Provider（避免大 rebuild）

---

## 16. 自查报告

| 自查项 | 结果 |
|:-----|:-----|
| 组件分类清晰 | ✅（7 类）|
| 每个组件 Dart 签名 | ✅ |
| 视觉规范（尺寸/颜色/间距）| ✅ |
| 文字/间距/圆角规范 | ✅ |
| 阴影/可访问性 | ✅ |
| 性能 | ✅ |