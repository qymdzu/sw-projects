// lib/core/data/repositories/error_question_repository.dart
// 错题/漏题仓库（带 user_id 强制过滤）
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/error_question.dart';

class ErrorQuestionRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  ErrorQuestionRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 保存错题/漏题
  Future<ErrorQuestion> save({
    required String userId,
    required String type,
    String? imagePath,
    required String ocrText,
    String? ocrLatex,
    String? knowledgePoint,
    String? errorType,
    String? chapter,
    String status = ErrorQuestion.statusPendingWeekly,
  }) async {
    final now = DateTime.now();
    final q = ErrorQuestion(
      id: _uuid.v4(),
      userId: userId,
      type: type,
      imagePath: imagePath,
      ocrText: ocrText,
      ocrLatex: ocrLatex,
      knowledgePoint: knowledgePoint,
      errorType: errorType,
      chapter: chapter,
      status: status,
      createdAt: now,
      updatedAt: now,
    );
    await _db.insert('error_questions', q.toMap());
    return q;
  }

  /// 按 ID 查
  Future<ErrorQuestion?> findById(String id) async {
    final rows = await _db.query('error_questions', where: 'id = ?', whereArgs: [id], limit: 1);
    if (rows.isEmpty) return null;
    return ErrorQuestion.fromMap(rows.first);
  }

  /// 查某用户的周题池（v0.1 用：本周所有 status=pendingWeekly 的题）
  Future<List<ErrorQuestion>> findByWeeklyPool({
    required String userId,
    required String type,
  }) async {
    final rows = await _db.query(
      'error_questions',
      where: 'user_id = ? AND type = ? AND status = ?',
      whereArgs: [userId, type, ErrorQuestion.statusPendingWeekly],
      orderBy: 'created_at ASC',
    );
    return rows.map(ErrorQuestion.fromMap.new).toList();
  }

  /// 查某用户的月观察池
  Future<List<ErrorQuestion>> findByMonthlyPool({required String userId}) async {
    final rows = await _db.query(
      'error_questions',
      where: 'user_id = ? AND status = ?',
      whereArgs: [userId, ErrorQuestion.statusPassWeekly],
      orderBy: 'updated_at ASC',
    );
    return rows.map(ErrorQuestion.fromMap.new).toList();
  }

  /// 查某用户的归档
  Future<List<ErrorQuestion>> findArchived({required String userId, int limit = 50}) async {
    final rows = await _db.query(
      'error_questions',
      where: 'user_id = ? AND status = ?',
      whereArgs: [userId, ErrorQuestion.statusArchived],
      orderBy: 'archived_at DESC',
      limit: limit,
    );
    return rows.map(ErrorQuestion.fromMap.new).toList();
  }

  /// 随机抽 N 道题（周测用）
  Future<List<ErrorQuestion>> pickRandom({
    required String userId,
    required String type,
    int count = 5,
  }) async {
    final pool = await findByWeeklyPool(userId: userId, type: type);
    pool.shuffle();
    return pool.take(count).toList();
  }

  /// 查某知识点的所有题
  Future<List<ErrorQuestion>> findByKnowledgePoint({required String userId, required String knowledgePoint}) async {
    final rows = await _db.query(
      'error_questions',
      where: 'user_id = ? AND knowledge_point = ?',
      whereArgs: [userId, knowledgePoint],
      orderBy: 'created_at DESC',
    );
    return rows.map(ErrorQuestion.fromMap.new).toList();
  }

  /// 更新题目
  Future<void> update(ErrorQuestion q) async {
    await _db.update('error_questions', q.toMap(), where: 'id = ?', whereArgs: [q.id]);
  }

  /// 更新状态
  Future<void> updateStatus({required String id, required String status, DateTime? archivedAt}) async {
    final values = <String, Object?>{
      'status': status,
      'updated_at': DateTime.now().millisecondsSinceEpoch,
    };
    if (archivedAt != null) {
      values['archived_at'] = archivedAt.millisecondsSinceEpoch;
    }
    await _db.update('error_questions', values, where: 'id = ?', whereArgs: [id]);
  }

  /// 增加重试次数
  Future<void> incrementRetryCount(String id) async {
    await _db.rawUpdate(
      'UPDATE error_questions SET retry_count = retry_count + 1, updated_at = ? WHERE id = ?',
      [DateTime.now().millisecondsSinceEpoch, id],
    );
  }

  /// 删除
  Future<void> delete(String id) async {
    await _db.delete('error_questions', where: 'id = ?', whereArgs: [id]);
  }

  /// 删除图片路径（识别后清空原图）
  Future<void> clearImagePath(String id) async {
    await _db.update('error_questions', {'image_path': null, 'updated_at': DateTime.now().millisecondsSinceEpoch}, where: 'id = ?', whereArgs: [id]);
  }
}
