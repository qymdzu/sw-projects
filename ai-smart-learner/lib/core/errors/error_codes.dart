// lib/core/errors/error_codes.dart
// 错误码常量
class ErrorCodes {
  // OCR
  static const ocrAuth = 'OCR-AUTH-001';
  static const ocrRate = 'OCR-RATE-002';
  static const ocrInternal = 'OCR-INT-003';
  static const ocrInvalidImage = 'OCR-INV-004';

  // AI
  static const aiAuth = 'AI-AUTH-001';
  static const aiRate = 'AI-RATE-002';
  static const aiContent = 'AI-CONTENT-003';
  static const aiInternal = 'AI-INT-004';
  static const aiTimeout = 'AI-TIMEOUT-005';

  // DB
  static const dbConstraint = 'DB-CONST-001';
  static const dbInternal = 'DB-INT-002';
  static const dbNotFound = 'DB-NOT-003';

  // State
  static const stInvalid = 'ST-INV-001';
  static const stMissing = 'ST-MISS-002';

  // Network
  static const netTimeout = 'NET-TIMEOUT-001';
  static const netOffline = 'NET-OFFLINE-002';
  static const netServerError = 'NET-SRV-003';

  // File
  static const fileNotFound = 'FILE-NOT-FOUND-001';
  static const filePermission = 'FILE-PERM-002';
  static const fileTooLarge = 'FILE-LARGE-003';

  // Auth
  static const authNoUser = 'AUTH-NO-USER-001';
  static const authUserNotFound = 'AUTH-NOT-FOUND-002';

  // Config
  static const cfgInvalid = 'CFG-INV-001';
  static const cfgMissing = 'CFG-MISS-002';
}

/// 用户友好的错误消息
class UserFacingMessages {
  static String forCode(String code) {
    switch (code) {
      case ErrorCodes.ocrAuth:
      case ErrorCodes.aiAuth:
        return '服务异常，请联系开发者';
      case ErrorCodes.ocrRate:
        return '识别太快，请稍后再试';
      case ErrorCodes.aiRate:
        return 'AI 调用过快，请稍后再试';
      case ErrorCodes.ocrInternal:
        return '识别失败，请重试';
      case ErrorCodes.aiInternal:
        return 'AI 调用失败，已自动重试';
      case ErrorCodes.aiContent:
        return 'AI 返回内容异常，自动重试中';
      case ErrorCodes.aiTimeout:
        return 'AI 响应超时';
      case ErrorCodes.dbInternal:
        return '数据库异常，请重启 App';
      case ErrorCodes.netTimeout:
        return '网络超时，请检查网络';
      case ErrorCodes.netOffline:
        return '当前无网络，请联网后重试';
      case ErrorCodes.netServerError:
        return '服务异常，请稍后再试';
      case ErrorCodes.fileNotFound:
        return '文件不存在';
      case ErrorCodes.filePermission:
        return '权限不足，请在设置中开启';
      case ErrorCodes.fileTooLarge:
        return '文件过大，请压缩后再试';
      case ErrorCodes.authNoUser:
        return '请先选择用户';
      case ErrorCodes.stInvalid:
        return '题目状态异常';
      default:
        return '发生未知错误';
    }
  }
}
