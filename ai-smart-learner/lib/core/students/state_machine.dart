// lib/core/students/state_machine.dart
// 错题状态机
import '../data/models/error_question.dart';
import '../errors/app_error.dart';
import '../errors/error_codes.dart';

class StateMachine {
  /// 计算下一状态
  static String next(String current, {required bool isCorrect}) {
    switch (current) {
      case ErrorQuestion.statusPendingWeekly:
        // 答对 → 升；答错 → 保持（生成变式题）
        return isCorrect ? ErrorQuestion.statusPassWeekly : ErrorQuestion.statusPendingWeekly;
      case ErrorQuestion.statusPassWeekly:
        // 升到月观察
        return ErrorQuestion.statusPendingMonthly;
      case ErrorQuestion.statusPendingMonthly:
        // 答对 → 归档；答错 → 打回（"假性掌握"）
        return isCorrect ? ErrorQuestion.statusArchived : ErrorQuestion.statusPendingWeekly;
      case ErrorQuestion.statusArchived:
        // 归档是终态
        return ErrorQuestion.statusArchived;
      default:
        throw BusinessError(
          code: ErrorCodes.stInvalid,
          message: '未知状态: $current',
        );
    }
  }

  /// 校验状态转换合法
  static bool isValidTransition(String from, String to) {
    const validTransitions = {
      'pendingWeekly|passWeekly',
      'pendingWeekly|pendingWeekly',  // 重答
      'passWeekly|pendingMonthly',
      'pendingMonthly|archived',
      'pendingMonthly|pendingWeekly',  // 回炉
      'archived|archived',  // 终态
    };
    return validTransitions.contains('$from|$to');
  }
}
