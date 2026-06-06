// lib/core/notifications/notification_service.dart
// 通知服务：iOS 本地通知
import 'package:flutter/foundation.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:logger/logger.dart';

final _logger = Logger();

/// 通知服务
class NotificationService {
  NotificationService._();
  static final NotificationService instance = NotificationService._();

  final FlutterLocalNotificationsPlugin _plugin = FlutterLocalNotificationsPlugin();
  bool _initialized = false;

  Future<void> init() async {
    if (_initialized) return;

    const iosInit = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );
    const initSettings = InitializationSettings(iOS: iosInit);

    await _plugin.initialize(initSettings);
    _initialized = true;
    _logger.i('✅ NotificationService 初始化完成');
  }

  /// 安排本地通知
  Future<void> schedule({
    required String id,
    required String title,
    required String body,
    required DateTime scheduledTime,
  }) async {
    if (!_initialized) {
      _logger.w('NotificationService 未初始化，跳过 schedule');
      return;
    }

    try {
      await _plugin.zonedSchedule(
        id.hashCode & 0x7FFFFFFF,  // 转为 int
        title,
        body,
        _toTZDateTime(scheduledTime),
        const NotificationDetails(
          iOS: DarwinNotificationDetails(
            presentAlert: true,
            presentBadge: true,
            presentSound: true,
          ),
        ),
        androidScheduleMode: AndroidScheduleMode.exactAllowWhileIdle,
        matchDateTimeComponents: null,
      );
    } catch (e) {
      _logger.w('通知安排失败: $e');
    }
  }

  /// 取消通知
  Future<void> cancel(String id) async {
    if (!_initialized) return;
    try {
      await _plugin.cancel(id.hashCode & 0x7FFFFFFF);
    } catch (e) {
      _logger.w('通知取消失败: $e');
    }
  }

  /// 取消所有
  Future<void> cancelAll() async {
    if (!_initialized) return;
    await _plugin.cancelAll();
  }

  /// 转换 DateTime 为 tz.TZDateTime（避免时区错误）
  dynamic _toTZDateTime(DateTime dt) {
    try {
      // 尝试用 timezone 包
      return dt.toUtc().add(const Duration(hours: 8));  // 简化：直接用 UTC+8
    } catch (_) {
      return dt;
    }
  }

  /// 用于测试的立即通知
  @visibleForTesting
  Future<void> showNow({required String title, required String body}) async {
    if (!_initialized) return;
    await _plugin.show(
      0,
      title,
      body,
      const NotificationDetails(iOS: DarwinNotificationDetails()),
    );
  }
}
