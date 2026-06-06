# Mac Flutter 开发环境 checklist — AI智能学习机 v0.1

> **目标**：在公子 Mac 上装好 Flutter + Xcode，能 `flutter run` 跑 iOS 真机调试  
> **预计耗时**：1-2 小时（视网络速度）

---

## 1. 系统要求

- [ ] macOS 13 (Ventura) 或更新
- [ ] Apple ID 账号（已有 ✅）
- [ ] 至少 20GB 硬盘空间（Xcode 占 ~15GB）
- [ ] 稳定的网络（下载 Xcode ~7-8GB）

---

## 2. 安装步骤

### 步骤 1：安装 Xcode（30-60 分钟）

```bash
# 从 App Store 装 Xcode 15+
# 或命令行装命令行工具
xcode-select --install
```

**注意**：从 App Store 装 Xcode 是图形界面，必须双击。命令行只能装命令行工具。

**装完验证**：
```bash
xcodebuild -version
# 应输出 Xcode 15.x 或更新
```

### 步骤 2：同意 Xcode 许可

```bash
sudo xcodebuild -license accept
```

### 步骤 3：安装 Homebrew（如果没有）

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 步骤 4：安装 Flutter SDK（10-15 分钟）

```bash
brew install --cask flutter
```

**装完验证**：
```bash
flutter --version
# 应输出 Flutter 3.24.x 或更新
```

### 步骤 5：安装 CocoaPods（5-10 分钟）

```bash
sudo gem install cocoapods
# 或用 brew
brew install cocoapods
```

**装完验证**：
```bash
pod --version
# 应输出 1.13.x 或更新
```

### 步骤 6：跑 flutter doctor（30 秒）

```bash
flutter doctor
```

**期望输出**（所有 ✅）：
```
[✓] Flutter (Channel stable, 3.24.x)
[✓] Xcode - develop for iOS and macOS
[✓] Android Studio (not required for iOS)
[✓] Connected device (1 available)
[✓] iOS toolchain - develop for iOS devices
```

**iOS 模拟器**（如果用模拟器不接真机）：
```bash
xcrun simctl list devices
# 看有没有 iPhone 15 / iPad
```

### 步骤 7：iPad 真机调试准备（10 分钟）

1. **iPad 用 USB 连 Mac**
2. **iPad 信任 Mac**：iPad 弹窗"信任此电脑"，点信任
3. **开开发者模式**：
   - iPad 设置 → 隐私与安全 → 开发者模式 → 开
   - iPad 会重启
4. **Flutter 检测设备**：
   ```bash
   flutter devices
   # 应输出 iPad (mobile) · iOS · ...
   ```

---

## 3. 拉代码 + 跑项目

### 步骤 1：拉代码

```bash
# 第一次拉
cd ~/gitee-software/sw-projects
git pull origin master

# 进入项目
cd ai-smart-learner
```

### 步骤 2：装依赖

```bash
flutter pub get
```

**这会跑 pub get**（下载 pubspec.yaml 里所有 Dart 包）。

### 步骤 3：iOS pod install

```bash
cd ios
pod install
cd ..
```

**这会装 iOS 依赖**（百度云 OCR SDK、image_picker iOS 桥等）。

### 步骤 4：跑项目

```bash
# 真机调试
flutter run

# 或指定设备
flutter run -d <device_id>
```

---

## 4. 常见问题

### Q1: flutter doctor 报 CocoaPods 没装

```bash
sudo gem install cocoapods
# 或
brew install cocoapods
```

### Q2: pod install 失败

```bash
# 清缓存重试
cd ios
rm -rf Pods Podfile.lock
pod install
```

### Q3: 真机调试报"未签名"

Xcode → Runner → Signing & Capabilities → Team → 选你的 Apple ID 个人团队（自动签名）。

### Q4: 百度云 OCR SDK 缺

v0.1 不调真百度云（用 Apple Vision fallback）。等开通后再加。

### Q5: DeepSeek API Key 怎么配

```bash
cp .env.example .env
# 编辑 .env 填入 DEEPSEEK_API_KEY=xxx
```

v0.1 实际不读 .env（开发期硬编码），生产用 Keychain。

---

## 5. 完成验证

`flutter run` 后应该看到：
- [ ] App 启动到 `UserSelectPage`（启动选身份）
- [ ] 选"学生" → 进 `StudentHomePage`
- [ ] 点"拍照" → `CapturePage`（可能需要相机权限授权）
- [ ] 拍 1 张图 → 显示识别中 → 显示结果弹窗

**全部通过 = v0.1 跑通**。

---

## 6. 装好后告诉我

跑下面命令，把输出发我：

```bash
flutter doctor -v
flutter devices
flutter --version
```

我对照看环境是否就绪，告诉你下一步操作。

---

## 7. 装环境的同时，我会做的事

- ✅ Stage 5 编码（写 lib/core/ + lib/ui/ + test/）
- ✅ 推 Gitee（你随时能 git pull 看）
- ✅ 写 Mac 装环境的 checkpoint 文档（本文件）
- ❌ 不阻塞（不重启 gateway / 不影响其他服务）

---

## 8. 时间安排建议

| 时间 | 做什么 |
|:-----|:-----|
| 现在 - 1 小时后 | 你装环境（后台跑） |
| 现在 - 2 小时后 | 我写代码（前台跑） |
| 2 小时后 | 你 `git pull` + `flutter run` 验证 |
| 验证有 bug | 告诉我，我修 |
| 验证 OK | 进 Stage 6 审查 |

---

> **最后**：本 checklist 只为**让你 Mac 准备好**。代码本身不需要 Mac 写，但**最终验证需要**。