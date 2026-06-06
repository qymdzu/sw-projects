// lib/core/data/repositories/user_repository.dart
// 用户仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/user.dart';

class UserRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  UserRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 创建用户
  Future<User> create({required String name, required String role, String? grade}) async {
    final user = User(
      id: _uuid.v4(),
      name: name,
      role: role,
      grade: grade,
      createdAt: DateTime.now(),
    );
    await _db.insert('users', user.toMap());
    return user;
  }

  /// 按 ID 查
  Future<User?> findById(String id) async {
    final rows = await _db.query('users', where: 'id = ?', whereArgs: [id], limit: 1);
    if (rows.isEmpty) return null;
    return User.fromMap(rows.first);
  }

  /// 查所有用户
  Future<List<User>> findAll() async {
    final rows = await _db.query('users', orderBy: 'created_at ASC');
    return rows.map(User.fromMap.new).toList();
  }

  /// 按角色查
  Future<List<User>> findByRole(String role) async {
    final rows = await _db.query('users', where: 'role = ?', whereArgs: [role], orderBy: 'created_at ASC');
    return rows.map(User.fromMap.new).toList();
  }

  /// 更新用户
  Future<void> update(User user) async {
    await _db.update('users', user.toMap(), where: 'id = ?', whereArgs: [user.id]);
  }

  /// 更新最后活跃时间
  Future<void> touch(String id) async {
    await _db.update('users', {'last_active_at': DateTime.now().millisecondsSinceEpoch}, where: 'id = ?', whereArgs: [id]);
  }

  /// 删除用户
  Future<void> delete(String id) async {
    await _db.delete('users', where: 'id = ?', whereArgs: [id]);
  }
}
