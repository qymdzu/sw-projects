// lib/core/auth/user_session.dart
// 用户 session（当前用户 + Riverpod 状态）
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../data/models/user.dart';
import '../data/repositories/user_repository.dart';
import '../providers.dart';

const _kCurrentUserIdKey = 'current_user_id';

/// 当前用户 provider（监听 user_id 变化重新查 User）
final currentUserProvider = StateNotifierProvider<CurrentUserNotifier, User?>((ref) {
  final repo = ref.watch(userRepositoryProvider);
  return CurrentUserNotifier(repo);
});

class CurrentUserNotifier extends StateNotifier<User?> {
  final UserRepository _repo;
  static const _storage = FlutterSecureStorage();

  CurrentUserNotifier(this._repo) : super(null) {
    _loadFromStorage();
  }

  Future<void> _loadFromStorage() async {
    final userId = await _storage.read(key: _kCurrentUserIdKey);
    if (userId == null) return;
    final user = await _repo.findById(userId);
    if (user != null) {
      state = user;
    }
  }

  Future<void> setUser(User user) async {
    await _storage.write(key: _kCurrentUserIdKey, value: user.id);
    state = user;
    await _repo.touch(user.id);
  }

  Future<void> clear() async {
    await _storage.delete(key: _kCurrentUserIdKey);
    state = null;
  }
}

/// User Repository provider
final userRepositoryProvider = Provider<UserRepository>((ref) {
  final db = ref.watch(databaseProvider);
  return UserRepository(db);
});

/// 当前用户 ID（异步读取）
final currentUserIdProvider = Provider<String?>((ref) {
  return ref.watch(currentUserProvider)?.id;
});
