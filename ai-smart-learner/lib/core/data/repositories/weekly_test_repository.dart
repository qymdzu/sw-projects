// lib/core/data/repositories/weekly_test_repository.dart
// 周测记录仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/weekly_test.dart';

class WeeklyTestRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  WeeklyTestRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 创建周测
  Future<WeeklyTest> create({
    required String userId,
    required String type,
    required DateTime weekStart,
    required List<String> questionIds,
  }) async {
    final test = WeeklyTest(
      id: _uuid.v4(),
      userId: userId,
      type: type,
      weekStart: weekStart,
      questionIds: questionIds,
    );
    await _db.insert('weekly_tests', test.toMap());
    return test;
  }

  /// 查用户的最近周测
  Future<WeeklyTest?> findLatest({required String userId, required String type}) async {
    final rows = await _db.query(
      'weekly_tests',
      where: 'user_id = ? AND type = ?',
      whereArgs: [userId, type],
      orderBy: 'week_start DESC',
      limit: 1,
    );
    if (rows.isEmpty) return null;
    return WeeklyTest.fromMap(rows.first);
  }

  /// 完成周测
  Future<void> complete({required String testId, required int correctCount}) async {
    await _db.update(
      'weekly_tests',
      {
        'completed_at': DateTime.now().millisecondsSinceEpoch,
        'correct_count': correctCount,
      },
      where: 'id = ?',
      whereArgs: [testId],
    );
  }

  /// 查历史
  Future<List<WeeklyTest>> findHistory({required String userId, int limit = 20}) async {
    final rows = await _db.query(
      'weekly_tests',
      where: 'user_id = ?',
      whereArgs: [userId],
      orderBy: 'week_start DESC',
      limit: limit,
    );
    return rows.map(WeeklyTest.fromMap.new).toList();
  }
}
