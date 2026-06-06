// lib/main.dart
// 入口：初始化服务 + 启动 App
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import 'app.dart';
import 'core/config/config_loader.dart';
import 'core/data/database.dart';
import 'core/errors/logger_init.dart';
import 'core/notifications/notification_service.dart';
import 'core/providers.dart';

final _logger = Logger();

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // 锁定竖屏
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);

  // 初始化 logger
  LoggerInit.init();

  // 加载配置
  final configLoader = ConfigLoader();
  try {
    await configLoader.load();
    _logger.i('✅ Config loaded');
  } catch (e, st) {
    _logger.e('❌ Config load failed', error: e, stackTrace: st);
    rethrow;
  }

  // 初始化数据库
  late final AppDatabase database;
  try {
    database = await AppDatabase.open();
    _logger.i('✅ Database opened (version ${AppDatabase.dbVersion})');
  } catch (e, st) {
    _logger.e('❌ Database open failed', error: e, stackTrace: st);
    rethrow;
  }

  // 初始化通知服务
  try {
    await NotificationService.instance.init();
    _logger.i('✅ Notification service initialized');
  } catch (e, st) {
    _logger.w('⚠️ Notification service init failed (non-fatal)', error: e, stackTrace: st);
  }

  runApp(
    ProviderScope(
      overrides: [
        configLoaderProvider.overrideWithValue(configLoader),
        databaseProvider.overrideWithValue(database),
      ],
      child: const AiSmartLearnerApp(),
    ),
  );
}
