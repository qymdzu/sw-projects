// lib/core/data/models/ai_call_log.dart
// AI 调用日志
class AiCallLog {
  final String id;
  final String userId;
  final String provider;
  final String purpose;
  final int? promptTokens;
  final int? completionTokens;
  final int? totalTokens;
  final int? costCents;
  final int? latencyMs;
  final bool success;
  final String? errorMsg;
  final DateTime createdAt;

  const AiCallLog({
    required this.id,
    required this.userId,
    required this.provider,
    required this.purpose,
    this.promptTokens,
    this.completionTokens,
    this.totalTokens,
    this.costCents,
    this.latencyMs,
    required this.success,
    this.errorMsg,
    required this.createdAt,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'user_id': userId,
    'provider': provider,
    'purpose': purpose,
    'prompt_tokens': promptTokens,
    'completion_tokens': completionTokens,
    'total_tokens': totalTokens,
    'cost_cents': costCents,
    'latency_ms': latencyMs,
    'success': success ? 1 : 0,
    'error_msg': errorMsg,
    'created_at': createdAt.millisecondsSinceEpoch,
  };

  factory AiCallLog.fromMap(Map<String, Object?> m) => AiCallLog(
    id: m['id']! as String,
    userId: m['user_id']! as String,
    provider: m['provider']! as String,
    purpose: m['purpose']! as String,
    promptTokens: m['prompt_tokens'] as int?,
    completionTokens: m['completion_tokens'] as int?,
    totalTokens: m['total_tokens'] as int?,
    costCents: m['cost_cents'] as int?,
    latencyMs: m['latency_ms'] as int?,
    success: (m['success']! as int) == 1,
    errorMsg: m['error_msg'] as String?,
    createdAt: DateTime.fromMillisecondsSinceEpoch(m['created_at']! as int),
  );
}
