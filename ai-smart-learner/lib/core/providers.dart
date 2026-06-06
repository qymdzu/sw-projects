// lib/core/providers.dart
// 全局 Riverpod Providers
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'config/config_loader.dart';
import 'data/database.dart';

/// Config loader provider
final configLoaderProvider = Provider<ConfigLoader>((ref) {
  throw UnimplementedError('Must be overridden in main()');
});

/// Database provider
final databaseProvider = Provider<AppDatabase>((ref) {
  throw UnimplementedError('Must be overridden in main()');
});
