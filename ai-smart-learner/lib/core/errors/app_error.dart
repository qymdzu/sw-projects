// lib/core/errors/app_error.dart
// 统一错误类
class AppError implements Exception {
  final String code;
  final String message;
  final String? detail;
  final Object? cause;
  final StackTrace? stackTrace;

  AppError({
    required this.code,
    required this.message,
    this.detail,
    this.cause,
    this.stackTrace,
  });

  @override
  String toString() => '[$code] $message${detail != null ? ' ($detail)' : ''}';
}

/// OCR 错误
class OcrError extends AppError {
  OcrError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}

/// AI 错误
class AiError extends AppError {
  AiError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}

/// 数据库错误
class DatabaseError extends AppError {
  DatabaseError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}

/// 业务错误（状态机、校验等）
class BusinessError extends AppError {
  BusinessError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}

/// 网络错误
class NetworkError extends AppError {
  NetworkError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}

/// 文件错误
class FileError extends AppError {
  FileError({
    required super.code,
    required super.message,
    super.detail,
    super.cause,
  });
}
