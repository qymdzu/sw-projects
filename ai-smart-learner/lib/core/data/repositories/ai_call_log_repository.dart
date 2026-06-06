// lib/core/data/repositories/ai_call_log_repository.dart
// AI 调用日志仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/ai_call_log.dart';

class AiCallLogRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  AiCallLogRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

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
  }) async {
    final log = AiCallLog(
      id: _uuid.v4(),
      userId: userId,
      provider: provider,
      purpose: purpose,
      promptTokens: promptTokens,
      completionTokens: completionTokens,
      totalTokens: totalTokens,
      costCents: costCents,
      latencyMs: latencyMs,
      success: success,
      errorMsg: errorMsg,
      createdAt: DateTime.now(),
    );
    await _db.insert('ai_call_logs', log.toMap());
  }

  /// 查用户的最近日志
  Future<List<AiCallLog>> findRecent({required String userId, int limit = 50}) async {
    final rows = await _db.query(
      'ai_call_logs',
      where: 'user_id = ?',
      whereArgs: [userId],
      orderBy: 'created_at DESC',
      limit: limit,
    );
    return rows.map(AiCallLog.fromMap.new).toList();
  }

  /// 聚合统计
  Future<AiCallStats> getStats({required String userId, required DateTime since}) async {
    final rows = await _db.rawQuery('''
      SELECT
        COUNT(*) as total_calls,
        SUM(CASE WHEN success = 1 THEN 1 ELSE 0 END) as success_count,
        SUM(COALESCE(cost_cents, 0)) as total_cost_cents,
        SUM(COALESCE(total_tokens, 0)) as total_tokens,
        AVG(COALESCE(latency_ms, 0)) as avg_latency_ms
      FROM ai_call_logs
      WHERE user_id = ? AND created_at >= ?
    ''', [userId, since.millisecondsSinceEpoch]);
    if (rows.isEmpty) return AiCallStats.empty();
    final r = rows.first;
    return AiCallStats(
      totalCalls: r['total_calls'] as int? ?? 0,
      successCount: r['success_count'] as int? ?? 0,
      totalCostCents: r['total_cost_cents'] as int? ?? 0,
      totalTokens: r['total_tokens'] as int? ?? 0,
      avgLatencyMs: (r['avg_latency_ms'] as num?)?.toInt() ?? 0,
    );
  }
}

class AiCallStats {
  final int totalCalls;
  final int successCount;
  final int totalCostCents;
  final int totalTokens;
  final int avgLatencyMs;

  const AiCallStats({
    required this.totalCalls,
    required this.successCount,
    required this.totalCostCents,
    required this.totalTokens,
    required this.avgLatencyMs,
  });

  factory AiCallStats.empty() => const AiCallStats(
    totalCalls: 0,
    successCount: 0,
    totalCostCents: 0,
    totalTokens: 0,
    avgLatencyMs: 0,
  );
}
