// lib/core/errors/logger_init.dart
// Logger 初始化
import 'package:logger/logger.dart';
import 'package:path_provider/path_provider.dart';
import 'dart:io';

class LoggerInit {
  static final List<LogEvent> _memoryLogs = [];
  static const int _maxMemoryLogs = 1000;

  static void init() {
    // v0.1 简化：仅控制台输出
    // v0.2 加文件日志：FileOutput(file)
    Logger.level = Level.info;
  }

  /// 记录到内存（v0.1 用于调试，v0.2 写文件）
  static void capture(LogEvent event) {
    _memoryLogs.add(event);
    if (_memoryLogs.length > _maxMemoryLogs) {
      _memoryLogs.removeAt(0);
    }
  }

  /// 获取最近的日志（v0.1 调试用）
  static List<LogEvent> recent() {
    return List.unmodifiable(_memoryLogs);
  }
}
