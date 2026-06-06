# API 设计 — AI智能学习机 v0.1

> **文档版本**：v0.1.0  
> **编制日期**：2026-06-06  
> **前置依赖**：`docs/design/系统架构设计.md` + `docs/design/数据模型设计.md`  
> **目标**：业务层与 UI 层、抽象层接口的完整契约（30 条需求全考虑）

---

## 1. API 设计原则

| 原则 | 实现 |
|:-----|:-----|
| **业务/UI 解耦** | UI 只调 Provider，Provider 暴露业务方法 |
| **抽象层稳定** | OCR/AI Provider 接口稳定，实现可换 |
| **数据隔离** | 所有 Repo 方法签名带 userId |
| **错误统一** | 抛 BusinessError（含 code + message）|
| **可测试** | 接口 + 依赖注入，便于 mock |

---

## 2. Repository API（数据层 → 业务层）

### 2.1 UserRepository

```dart
abstract class UserRepository {
  Future<User> create({required String name, required String role});
  Future<User?> findById(String id);
  Future<List<User>> findAll();
  Future<User?> getCurrentUser();
  Future<void> setCurrentUser(String userId);
  Future<void> delete(String userId);
}
```

### 2.2 ErrorQuestionRepository（学生题库）

```dart
abstract class ErrorQuestionRepository {
  Future<ErrorQuestion> save({
    required String userId,
    required String type,           // 'error' / 'missing'
    String? imagePath,
    required String ocrText,
    String? ocrLatex,
    String? knowledgePoint,
    String? errorType,
    String? chapter,
  });
  Future<ErrorQuestion?> findById(String id);
  Future<List<ErrorQuestion>> findByWeeklyPool({
    required String userId,
    required String type,           // 'error' / 'missing'
    required DateTime weekStart,
  });
  Future<List<ErrorQuestion>> findByMonthlyPool({required String userId});
  Future<List<ErrorQuestion>> findArchived({required String userId});
  Future<void> updateStatus({required String id, required String status});
  Future<void> incrementRetryCount(String id);
}

class ErrorQuestion {
  String id, userId, type, ocrText;
  String? imagePath, ocrLatex, knowledgePoint, errorType, chapter;
  String status;  // 'pendingWeekly' / 'passWeekly' / 'pendingMonthly' / 'archived'
  int retryCount;
  DateTime createdAt, updatedAt;
  DateTime? archivedAt;
}
```

### 2.3 VariantQuestionRepository

```dart
abstract class VariantQuestionRepository {
  Future<List<VariantQuestion>> saveBatch({
    required String parentQuestionId,
    required List<String> texts,
  });
  Future<List<VariantQuestion>> findByParent(String parentQuestionId);
  Future<void> updateStatus({required String id, required String status});
}

class VariantQuestion {
  String id, parentQuestionId, text;
  String? latex;
  int variantIndex;  // 1/2/3
  String status;
  DateTime createdAt;
}
```

### 2.4 KnowledgePointRepository

```dart
abstract class KnowledgePointRepository {
  Future<KnowledgePoint> upsert({
    required String userId,
    required String name,
    String? chapter,
    String? mastery,
    int? errorCount,
  });
  Future<KnowledgePoint?> findByName({required String userId, required String name, String? chapter});
  Future<List<KnowledgePoint>> findAll({required String userId});
  Future<List<KnowledgePoint>> findByMastery({required String userId, required String mastery});
  Future<List<KnowledgePoint>> findByChapter({required String userId, required String chapterId});
  Future<void> updateMastery({required String id, required String mastery});
}

class KnowledgePoint {
  String id, userId, name;
  String? chapter;
  String mastery;  // 'red' / 'yellow' / 'green'
  int errorCount;
  DateTime updatedAt;
}
```

### 2.5 ChapterRepository（v0.1 预置数据）

```dart
abstract class ChapterRepository {
  Future<List<Chapter>> findAll({String? subject});
  Future<Chapter?> findById(String id);
  Future<List<Chapter>> findByParent(String? parentId);
}

class Chapter {
  String id, name, subject;
  String? parentId;
  int? displayOrder;
}
```

### 2.6 OCRCacheRepository

```dart
abstract class OCRCacheRepository {
  Future<OCRResult?> findByImageHash(String imageHash);
  Future<void> save({required String imageHash, required OCRResult result});
  Future<void> cleanExpired({required DateTime olderThan});
}
```

### 2.7 AICallLogRepository（监控）

```dart
abstract class AICallLogRepository {
  Future<void> log({
    required String userId,
    required String provider,
    required String purpose,
    int? promptTokens,
    int? completionTokens,
    int? totalTokens,
    int? costCents,
    int? latencyMs,
    required bool success,
    String? errorMsg,
  });
  Future<AICallLogStats> getStats({required String userId, required DateTime since});
  Future<List<AICallLog>> findRecent({required String userId, int limit = 50});
}

class AICallLogStats {
  int totalCalls;
  int successCount;
  int totalCostCents;
  int totalTokens;
  Duration avgLatency;
}
```

---

## 3. 抽象层 API（业务层 → 抽象层）

### 3.1 OCRProvider 抽象

```dart
abstract class OCRProvider {
  String get name;
  bool get supportsHandwritingRemoval;
  bool supportsQuestionType(String questionType);

  Future<OCRResult> recognize({
    required File imageFile,
    bool removeHandwriting = true,
    Map<String, dynamic>? options,
  });

  Future<void> close();  // 释放资源
}

class OCRResult {
  String text;                      // 纯文本
  String? latex;                    // 数学公式 LaTeX（如果有）
  List<String> tags;                // 关键词标签
  double confidence;                // 0-1
  Map<String, dynamic> raw;         // 原始响应
  DateTime createdAt;
}
```

**实现类**（v0.1）：

| 类 | 来源 | 特点 |
|:-----|:-----|:-----|
| `BaiduEduProvider` | 百度云 OCR 教育版 | 默认，去笔迹能力强 |
| `AppleVisionProvider` | iOS Vision SDK | 离线 fallback，识别率一般 |

**实现类**（v0.2+）：`MathpixProvider`

### 3.2 AIProvider 抽象

```dart
abstract class AIProvider {
  String get name;
  List<String> get supportedModels;  // e.g. ['deepseek-chat', 'deepseek-reasoner']

  Future<AIResponse> chat({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
    int maxTokens = 1000,
    Duration? timeout,
  });

  Stream<String> chatStream({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
  });

  Future<void> close();
}

class Message {
  final String role;  // 'system' / 'user' / 'assistant'
  final String content;
  Message({required this.role, required this.content});
}

class AIResponse {
  String content;
  int promptTokens;
  int completionTokens;
  int totalTokens;
  Duration latency;
  String model;
  String finishReason;  // 'stop' / 'length' / 'content_filter'
  DateTime createdAt;
}
```

**实现类**（v0.1）：

| 类 | 来源 | 模型 |
|:-----|:-----|:-----|
| `DeepSeekProvider` | DeepSeek API | `deepseek-chat` / `deepseek-reasoner` |

**实现类**（v0.2+）：`ClaudeProvider` / `GPT4Provider` / `TongyiProvider`

### 3.3 NotificationProvider 抽象

```dart
abstract class NotificationProvider {
  Future<void> init();
  Future<void> schedule({
    required String id,
    required String title,
    required String body,
    required DateTime scheduledTime,
  });
  Future<void> cancel(String id);
  Future<List<PendingNotification>> getPending();
}

class PendingNotification {
  String id, title, body;
  DateTime scheduledTime;
}
```

**实现类**：`IOSLocalNotificationProvider`（iOS 本地通知，v0.1）

---

## 4. 业务引擎 API（业务层内部）

### 4.1 学生业务引擎

```dart
abstract class StudentEngine {
  // 录入
  Future<ErrorQuestion> captureQuestion({
    required String userId,
    required String type,             // 'error' / 'missing'
    required File imageFile,
  });
  
  // 智能打标
  Future<ErrorQuestion> autoTagging(String questionId);
  
  // 周测
  Future<WeeklyTest> generateWeeklyTest({required String userId, required String type});
  Future<WeeklyTestResult> submitWeeklyTest({
    required String testId,
    required Map<String, bool> answers,  // questionId -> isCorrect
  });
  
  // 月测（v0.2）
  Future<MonthlyTest> generateMonthlyTest({required String userId});
  Future<MonthlyTestResult> submitMonthlyTest({...});
  
  // 知识图谱（v0.2）
  Future<KnowledgeGraph> buildKnowledgeGraph({required String userId});
  
  // AI 私教（v0.2）
  Future<TutorCard> generateTutorCard({required String knowledgePointId});
}

class WeeklyTest {
  String id, userId, type;
  DateTime weekStart;
  List<String> questionIds;
  DateTime? completedAt;
  int? correctCount;
}

class WeeklyTestResult {
  int correctCount;
  int totalCount;
  List<VariantQuestion> newVariants;   // 答错时生成的变式题
  List<ErrorQuestion> promoted;       // 答对时升入月观察
}
```

### 4.2 大人业务引擎（v0.2 推）

```dart
abstract class AdultEngine {
  Future<BookNote> captureBookPage({required String userId, required File imageFile});
  Future<NoteSummary> generateSummary({required String bookNoteId, required File imageFile});
  Future<BookNote> generateNotes({required String bookNoteId});
  Future<File> exportNotes({
    required String bookNoteId,
    required String format,  // 'pdf' / 'markdown' / 'feishu_doc'
  });
}
```

### 4.3 共享业务引擎

```dart
abstract class SchedulerEngine {
  Future<void> scheduleWeeklyTest({required String userId});  // 周六 8:00
  Future<void> scheduleMonthlyTest({required String userId});  // 月末 20:00
  Future<void> scheduleSemesterReport({required String userId});  // 学期末
}

abstract class WorkspaceEngine {
  Future<void> setCurrentUser(String userId);
  Future<User> getCurrentUser();
  Stream<User> watchCurrentUser();
}

abstract class ConfigLoader {
  Future<void> load();
  OCRRoutingConfig get ocrConfig;
  AIRoutingConfig get aiConfig;
}
```

---

## 5. Provider API（UI 层 → 业务层，Riverpod）

```dart
// lib/ui/providers.dart

// 学生
final userIdProvider = StateProvider<String?>((ref) => null);

final captureProvider = Provider<StudentEngine>((ref) => 
  StudentEngineImpl(
    ocr: ref.watch(ocrFactoryProvider),
    ai: ref.watch(aiFactoryProvider),
    eqRepo: ref.watch(errorQuestionRepoProvider),
    kpRepo: ref.watch(knowledgePointRepoProvider),
  ),
);

final weeklyTestProvider = FutureProvider.family<WeeklyTest, String>((ref, type) async {
  final userId = ref.watch(userIdProvider);
  if (userId == null) throw 'No user';
  return ref.watch(captureProvider).generateWeeklyTest(userId: userId, type: type);
});

final rainbowChartProvider = FutureProvider<RainbowChartData>((ref) async {
  final userId = ref.watch(userIdProvider);
  if (userId == null) throw 'No user';
  final kps = await ref.watch(knowledgePointRepoProvider).findAll(userId: userId);
  return RainbowChartData(
    redCount: kps.where((k) => k.mastery == 'red').length,
    yellowCount: kps.where((k) => k.mastery == 'yellow').length,
    greenCount: kps.where((k) => k.mastery == 'green').length,
  );
});

// 大人（v0.2）
final adultEngineProvider = Provider<AdultEngine>(...);
```

---

## 6. 错误码规范

### 6.1 错误码格式

`{LAYER}-{CATEGORY}-{NUMBER}`

| Layer | 前缀 | 例 |
|:-----|:-----|:-----|
| OCR | `OCR-` | `OCR-AUTH-001` |
| AI | `AI-` | `AI-RATE-002` |
| DB | `DB-` | `DB-CONST-003` |
| State | `ST-` | `ST-INV-004` |
| Network | `NET-` | `NET-TIMEOUT-005` |
| Unknown | `UNK-` | `UNK-INT-006` |

### 6.2 错误类

```dart
class AppError implements Exception {
  final String code;     // 'OCR-AUTH-001'
  final String message;  // 用户友好消息
  final String? detail;  // 技术细节（给日志）
  final Object? cause;   // 原始异常
  
  AppError({required this.code, required this.message, this.detail, this.cause});
  
  @override
  String toString() => '[$code] $message';
}

class OCRError extends AppError {
  OCRError({required super.code, required super.message, super.detail, super.cause});
}

class AIError extends AppError {
  AIError({required super.code, required super.message, super.detail, super.cause});
}

class BusinessError extends AppError {
  BusinessError({required super.code, required super.message, super.detail, super.cause});
}
```

---

## 7. 错误码清单（v0.1）

| 错误码 | 含义 | 触发场景 | 用户提示 |
|:-----|:-----|:---------|:---------|
| `OCR-AUTH-001` | OCR 鉴权失败 | 百度云 API Key 失效 | "OCR 服务异常，请联系开发者" |
| `OCR-RATE-002` | OCR 限流 | 1 分钟内调超 10 次 | "识别太快，请稍后再试" |
| `OCR-INT-003` | OCR 服务异常 | 百度云 5xx | "识别失败，请重试" |
| `AI-AUTH-001` | AI 鉴权失败 | DeepSeek Key 失效 | "AI 服务异常，请联系开发者" |
| `AI-RATE-002` | AI 限流 | DeepSeek RPM 超限 | "AI 调用过快，请稍后再试" |
| `AI-CONTENT-003` | AI 输出内容过滤 | 变式题被 deepseek 拒 | "AI 返回内容异常，自动重试中" |
| `AI-INT-004` | AI 服务异常 | DeepSeek 5xx | "AI 调用失败，已重试 2 次" |
| `DB-CONST-001` | DB 约束违反 | FK / UNIQUE 失败 | "数据错误，请联系开发者" |
| `DB-INT-002` | DB 内部错误 | sqflite 抛错 | "数据库异常，请重启 App" |
| `ST-INV-001` | 状态机非法转换 | archived 状态想变 passWeekly | "题目状态异常" |
| `NET-TIMEOUT-001` | 网络超时 | dio 超时 | "网络超时，请检查网络" |
| `NET-OFFLINE-002` | 无网络 | dio ConnectionError | "当前无网络，请联网后重试" |
| `FILE-NOT-FOUND-001` | 文件不存在 | image_picker 返回空 | "拍照失败，请重试" |
| `FILE-PERM-002` | 文件权限拒绝 | iOS 沙盒权限 | "无法访问照片，请检查权限" |

---

## 8. API 演进与版本管理

### 8.1 兼容原则

- **接口签名**：加可选参数 ✅，改参数类型 ❌，改方法名 ❌
- **返回值**：加字段 ✅，删字段 ❌，改类型 ❌
- **错误码**：加新码 ✅，删旧码（v0.x → v1.0 时才能删）

### 8.2 版本号

- **API 版本跟随 App 版本**（v0.1.0 = API v1）
- **数据模型版本**用 sqflite `_dbVersion`（v0.1 = 1）

---

## 9. 自查报告

| 自查项 | 结果 |
|:-----|:-----|
| Repository API 完整 | ✅（7 个 Repo）|
| 抽象层 API 完整 | ✅（3 个抽象）|
| 业务引擎 API 完整 | ✅（3 个 Engine）|
| Riverpod Provider 完整 | ✅ |
| 错误码规范 | ✅ |
| 错误码清单 v0.1 | ✅（14 个）|
| API 演进原则 | ✅ |
| 30 条需求 API 全考虑 | ✅ |