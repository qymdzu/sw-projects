// lib/core/scheduler/scheduler_service.dart
// 调度服务：周末/月末自动推送
import 'package:logger/logger.dart';

import 'notification_service.dart';

final _logger = Logger();

class SchedulerService {
  final NotificationService _notifications;

  SchedulerService(this._notifications);

  /// 安排周测推送（周六 8:00）
  Future<void> scheduleWeeklyTest({required String userId, required String userName}) async {
    final next = _nextSaturday8am();
    await _notifications.schedule(
      id: 'weekly_test_$userId',
      title: '📝 本周错题歼灭战',
      body: '嗨 $userName，本周的错题已就绪，5 道题等你挑战！',
      scheduledTime: next,
    );
    _logger.i('✅ 周测推送已安排: $userId @ $next');
  }

  /// 取消周测推送
  Future<void> cancelWeeklyTest({required String userId}) async {
    await _notifications.cancel('weekly_test_$userId');
  }

  /// 安排月测推送（月末 20:00，v0.2 推）
  Future<void> scheduleMonthlyTest({required String userId, required String userName}) async {
    final next = _nextMonthEnd20pm();
    await _notifications.schedule(
      id: 'monthly_test_$userId',
      title: '📊 本月错题巩固',
      body: '嗨 $userName，本月的观察题已就绪，3-5 道题等你巩固！',
      scheduledTime: next,
    );
  }

  /// 下周六 8:00
  DateTime _nextSaturday8am() {
    final now = DateTime.now();
    var daysUntilSat = (DateTime.saturday - now.weekday) % 7;
    if (daysUntilSat == 0 && now.hour >= 8) {
      daysUntilSat = 7;
    }
    final next = now.add(Duration(days: daysUntilSat));
    return DateTime(next.year, next.month, next.day, 8, 0);
  }

  /// 下月末 20:00
  DateTime _nextMonthEnd20pm() {
    final now = DateTime.now();
    // 下个月最后一天
    final nextMonth = now.month == 12 ? DateTime(now.year + 1, 1, 1) : DateTime(now.year, now.month + 1, 1);
    final lastDay = nextMonth.subtract(const Duration(days: 1));
    return DateTime(lastDay.year, lastDay.month, lastDay.day, 20, 0);
  }
}
