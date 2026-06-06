// lib/core/data/repositories/variant_question_repository.dart
// 变式题仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/variant_question.dart';
import '../models/error_question.dart';

class VariantQuestionRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  VariantQuestionRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 批量保存变式题（3 道一组）
  Future<List<VariantQuestion>> saveBatch({
    required String parentQuestionId,
    required List<String> texts,
    String? latex,
  }) async {
    if (texts.length > 3) {
      throw ArgumentError('At most 3 variants per parent');
    }
    final now = DateTime.now();
    final result = <VariantQuestion>[];
    for (var i = 0; i < texts.length; i++) {
      final v = VariantQuestion(
        id: _uuid.v4(),
        parentQuestionId: parentQuestionId,
        text: texts[i],
        latex: latex,
        variantIndex: i + 1,
        status: ErrorQuestion.statusPendingWeekly,
        createdAt: now,
      );
      await _db.insert('variant_questions', v.toMap());
      result.add(v);
    }
    return result;
  }

  /// 查父题的所有变式题
  Future<List<VariantQuestion>> findByParent(String parentQuestionId) async {
    final rows = await _db.query(
      'variant_questions',
      where: 'parent_question_id = ?',
      whereArgs: [parentQuestionId],
      orderBy: 'variant_index ASC',
    );
    return rows.map(VariantQuestion.fromMap.new).toList();
  }

  /// 查某用户的变式题（v0.1 简化：所有变式题都进入下周池）
  Future<List<VariantQuestion>> findByUser({required String userId, String? status}) async {
    // 通过 JOIN error_questions 拿到 user_id
    final where = status != null
        ? 'eq.user_id = ? AND vq.status = ?'
        : 'eq.user_id = ?';
    final args = status != null ? [userId, status] : [userId];
    final rows = await _db.rawQuery('''
      SELECT vq.* FROM variant_questions vq
      INNER JOIN error_questions eq ON vq.parent_question_id = eq.id
      WHERE $where
      ORDER BY vq.created_at ASC
    ''', args);
    return rows.map(VariantQuestion.fromMap.new).toList();
  }

  /// 更新状态
  Future<void> updateStatus({required String id, required String status}) async {
    await _db.update('variant_questions', {'status': status}, where: 'id = ?', whereArgs: [id]);
  }

  /// 删除
  Future<void> delete(String id) async {
    await _db.delete('variant_questions', where: 'id = ?', whereArgs: [id]);
  }
}
