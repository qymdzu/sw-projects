// lib/core/data/repositories/knowledge_point_repository.dart
// 知识点仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/knowledge_point.dart';

class KnowledgePointRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  KnowledgePointRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 按 (user, name, chapter) upsert 知识点
  Future<KnowledgePoint> upsert({
    required String userId,
    required String name,
    String? chapter,
    String mastery = KnowledgePoint.masteryRed,
  }) async {
    final existing = await findByName(userId: userId, name: name, chapter: chapter);
    if (existing != null) {
      return existing;
    }
    final kp = KnowledgePoint(
      id: _uuid.v4(),
      userId: userId,
      name: name,
      chapter: chapter,
      mastery: mastery,
      updatedAt: DateTime.now(),
    );
    await _db.insert('knowledge_points', kp.toMap());
    return kp;
  }

  /// 按 (user, name, chapter) 查
  Future<KnowledgePoint?> findByName({required String userId, required String name, String? chapter}) async {
    final rows = await _db.query(
      'knowledge_points',
      where: 'user_id = ? AND name = ? AND chapter ${chapter == null ? "IS NULL" : "= ?"}',
      whereArgs: chapter == null ? [userId, name] : [userId, name, chapter],
      limit: 1,
    );
    if (rows.isEmpty) return null;
    return KnowledgePoint.fromMap(rows.first);
  }

  /// 查某用户所有知识点
  Future<List<KnowledgePoint>> findAll({required String userId}) async {
    final rows = await _db.query('knowledge_points', where: 'user_id = ?', whereArgs: [userId], orderBy: 'error_count DESC');
    return rows.map(KnowledgePoint.fromMap.new).toList();
  }

  /// 按掌握度查
  Future<List<KnowledgePoint>> findByMastery({required String userId, required String mastery}) async {
    final rows = await _db.query('knowledge_points', where: 'user_id = ? AND mastery = ?', whereArgs: [userId, mastery]);
    return rows.map(KnowledgePoint.fromMap.new).toList();
  }

  /// 更新掌握度
  Future<void> updateMastery({required String id, required String mastery}) async {
    await _db.update(
      'knowledge_points',
      {'mastery': mastery, 'updated_at': DateTime.now().millisecondsSinceEpoch},
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  /// 增加错题数
  Future<void> incrementErrorCount({required String id, int by = 1}) async {
    await _db.rawUpdate(
      'UPDATE knowledge_points SET error_count = error_count + ?, last_error_at = ?, updated_at = ? WHERE id = ?',
      [by, DateTime.now().millisecondsSinceEpoch, DateTime.now().millisecondsSinceEpoch, id],
    );
  }

  /// 删除
  Future<void> delete(String id) async {
    await _db.delete('knowledge_points', where: 'id = ?', whereArgs: [id]);
  }
}
