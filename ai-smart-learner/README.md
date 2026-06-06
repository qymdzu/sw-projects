# AI智能学习机

> 学生错题/漏题管理 + 大人读书笔记的 Flutter iOS/iPadOS 自用型学习软件

## 项目状态

- ✅ **v0.1 MVP** 编码完成（推 Gitee，Mac 验证中）
- ⏳ Stage 6 审查 + Stage 7 测试待启动
- ⏳ Stage 8 部署（iPad 真机调试）待启动

## 技术栈

- **框架**：Flutter 3.24+（Dart 3.5+）
- **状态管理**：Riverpod 2.5+
- **路由**：go_router 14+
- **数据库**：sqflite（SQLite）
- **AI**：DeepSeek API（可热插拔抽象层）
- **OCR**：百度云教育版（可热插拔抽象层）
- **iOS 最低**：15.0+

## v0.1 已实现（14 条 P0）

| FR | 标题 | 状态 |
|:-:|:-----|:-----|
| FR-001 | 错题拍照录入 | ✅ |
| FR-002 | 漏题拍照录入 | ✅ |
| FR-005 | 本周错题池 | ✅ |
| FR-006 | 本周漏题池 | ✅ |
| FR-007 | 周末周测触发 | ✅（手动 + 调度）|
| FR-008 | 周测答对→升 | ✅ |
| FR-009 | 周测答错→变式 | ✅ |
| FR-013 | 智能打标 | ✅ |
| FR-017 | 知识彩虹图 | ✅ |
| FR-026 | 启动选身份 | ✅ |
| FR-027 | 工作空间隔离 | ✅ |
| FR-028 | 身份切换 | ✅ |
| FR-029 | OCR 抽象层 | ✅（百度云 + Apple Vision）|
| FR-030 | AI 抽象层 | ✅（DeepSeek）|

## 快速开始（Mac）

参见 [`docs/deploy/Mac环境安装checklist.md`](docs/deploy/Mac环境安装checklist.md)

```bash
# 1. 拉代码
cd ~/gitee-software/sw-projects
git pull origin master
cd ai-smart-learner

# 2. 装依赖
flutter pub get

# 3. iOS pod
cd ios && pod install && cd ..

# 4. 跑
flutter run
```

## 项目结构

```
ai-smart-learner/
├── lib/
│   ├── main.dart                    # 入口
│   ├── app.dart                     # App 根
│   ├── config/                      # 主题、路由
│   ├── core/                        # 业务核心
│   │   ├── auth/                    # 用户 session
│   │   ├── data/                    # 数据层（13 张表）
│   │   ├── ocr/                     # OCR 抽象层
│   │   ├── ai/                      # AI 抽象层
│   │   ├── students/                # 学生业务
│   │   ├── scheduler/               # 调度
│   │   ├── notifications/           # 通知
│   │   ├── config/                  # 配置加载
│   │   ├── errors/                  # 错误类
│   │   └── utils/
│   └── ui/                          # UI 层
│       ├── shared/                  # 共享（user_select / settings）
│       ├── students/                # 学生 7 个页面
│       └── widgets/                 # 通用 widget
├── assets/                          # 资源
├── test/                            # 测试
├── ios/                             # iOS 工程
├── docs/                            # 全部设计文档
└── pubspec.yaml
```

## 设计文档

- [需求规格说明书](docs/requirements/需求规格说明书.md)
- [技术方案说明书](docs/research/技术方案说明书.md)
- [系统架构设计](docs/design/系统架构设计.md)
- [数据模型设计](docs/design/数据模型设计.md)
- [API 设计](docs/design/API设计.md)
- [目录结构](docs/design/目录结构.md)
- [UX 交互流程](docs/ux/交互流程设计.md)

## 元目标（公子 2026-06-06 拍板）

1. 验证模型（翠花）开发软件的能力，实现软件开发集群的工程化开发能力
2. 保证开发出的软件没有问题，真正可使用
3. 验证 pipeline 管道在控制工程化开发软件方面是否存在问题

## License

MIT
