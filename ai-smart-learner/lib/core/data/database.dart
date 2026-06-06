// lib/core/data/database.dart
// sqflite 数据库初始化 + 表创建 + 迁移
import 'dart:async';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:logger/logger.dart';

import '../errors/app_error.dart';
import '../errors/error_codes.dart';

final _logger = Logger();

class AppDatabase {
  static const int dbVersion = 1;
  static const String _dbName = 'ai_smart_learner.db';

  final Database _db;
  AppDatabase._(this._db);

  Database get db => _db;

  /// 打开数据库（应用启动时调一次）
  static Future<AppDatabase> open() async {
    final dir = await getApplicationDocumentsDirectory();
    final path = p.join(dir.path, _dbName);
    _logger.i('📦 Database path: $path');

    final db = await openDatabase(
      path,
      version: dbVersion,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
    return AppDatabase._(db);
  }

  static Future<void> _onCreate(Database db, int version) async {
    _logger.i('📦 Creating database tables (version $version)');
    await db.transaction((txn) async {
      for (final ddl in _allTableDdls) {
        await txn.execute(ddl);
      }
    });
    _logger.i('✅ All tables created (${_allTableDdls.length} tables)');
  }

  static Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    _logger.i('📦 Migrating database: $oldVersion -> $newVersion');
    // v0.1 = v1，无迁移；v0.2+ 写迁移脚本
  }

  /// 关闭数据库
  Future<void> close() async {
    await _db.close();
  }

  /// 清空所有数据（v0.1 调试用）
  Future<void> clearAll() async {
    await _db.transaction((txn) async {
      for (final table in _allTables) {
        await txn.delete(table);
      }
    });
  }

  /// 所有表名（清空/迁移用）
  static const List<String> _allTables = [
    'users',
    'error_questions',
    'variant_questions',
    'knowledge_points',
    'chapters',
    'weekly_tests',
    'ocr_cache',
    'ai_call_logs',
  ];

  /// 所有表 DDL（按顺序创建）
  static const List<String> _allTableDdls = [
    '''
    CREATE TABLE users (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      role TEXT NOT NULL CHECK (role IN ('student', 'adult')),
      grade TEXT,
      created_at INTEGER NOT NULL,
      last_active_at INTEGER
    )
    ''',
    'CREATE INDEX idx_users_role ON users(role)',
    '''
    CREATE TABLE error_questions (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      type TEXT NOT NULL CHECK (type IN ('error', 'missing')),
      image_path TEXT,
      ocr_text TEXT NOT NULL,
      ocr_latex TEXT,
      knowledge_point TEXT,
      error_type TEXT,
      chapter TEXT,
      status TEXT NOT NULL CHECK (status IN ('pendingWeekly', 'passWeekly', 'pendingMonthly', 'archived')),
      retry_count INTEGER DEFAULT 0,
      created_at INTEGER NOT NULL,
      updated_at INTEGER NOT NULL,
      archived_at INTEGER,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )
    ''',
    'CREATE INDEX idx_eq_user_status ON error_questions(user_id, status)',
    'CREATE INDEX idx_eq_user_type ON error_questions(user_id, type)',
    'CREATE INDEX idx_eq_kp ON error_questions(knowledge_point)',
    '''
    CREATE TABLE variant_questions (
      id TEXT PRIMARY KEY,
      parent_question_id TEXT NOT NULL,
      text TEXT NOT NULL,
      latex TEXT,
      variant_index INTEGER NOT NULL CHECK (variant_index IN (1, 2, 3)),
      status TEXT NOT NULL,
      created_at INTEGER NOT NULL,
      FOREIGN KEY (parent_question_id) REFERENCES error_questions(id) ON DELETE CASCADE
    )
    ''',
    'CREATE INDEX idx_vq_parent ON variant_questions(parent_question_id)',
    'CREATE INDEX idx_vq_status ON variant_questions(status)',
    '''
    CREATE TABLE knowledge_points (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      name TEXT NOT NULL,
      chapter TEXT,
      mastery TEXT NOT NULL CHECK (mastery IN ('red', 'yellow', 'green')),
      error_count INTEGER DEFAULT 0,
      last_error_at INTEGER,
      updated_at INTEGER NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      UNIQUE (user_id, name, chapter)
    )
    ''',
    'CREATE INDEX idx_kp_user_mastery ON knowledge_points(user_id, mastery)',
    '''
    CREATE TABLE chapters (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      subject TEXT NOT NULL,
      stage TEXT NOT NULL CHECK (stage IN ('primary', 'junior', 'senior')),
      parent_id TEXT,
      display_order INTEGER,
      FOREIGN KEY (parent_id) REFERENCES chapters(id)
    )
    ''',
    'CREATE INDEX idx_chapters_subject_stage ON chapters(subject, stage)',
    '''
    CREATE TABLE weekly_tests (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      type TEXT NOT NULL CHECK (type IN ('error', 'missing')),
      week_start INTEGER NOT NULL,
      question_ids TEXT NOT NULL,
      completed_at INTEGER,
      correct_count INTEGER,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )
    ''',
    'CREATE INDEX idx_wt_user_week ON weekly_tests(user_id, week_start)',
    '''
    CREATE TABLE ocr_cache (
      id TEXT PRIMARY KEY,
      image_hash TEXT NOT NULL UNIQUE,
      ocr_result_json TEXT NOT NULL,
      created_at INTEGER NOT NULL
    )
    ''',
    'CREATE INDEX idx_ocr_hash ON ocr_cache(image_hash)',
    '''
    CREATE TABLE ai_call_logs (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      provider TEXT NOT NULL,
      purpose TEXT NOT NULL,
      prompt_tokens INTEGER,
      completion_tokens INTEGER,
      total_tokens INTEGER,
      cost_cents INTEGER,
      latency_ms INTEGER,
      success INTEGER NOT NULL,
      error_msg TEXT,
      created_at INTEGER NOT NULL,
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )
    ''',
    'CREATE INDEX idx_ai_user_created ON ai_call_logs(user_id, created_at)',
  ];
}
