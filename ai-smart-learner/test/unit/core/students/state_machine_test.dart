// test/unit/core/students/state_machine_test.dart
// 错题状态机单测
import 'package:flutter_test/flutter_test.dart';
import 'package:ai_smart_learner/core/students/state_machine.dart';
import 'package:ai_smart_learner/core/data/models/error_question.dart';

void main() {
  group('StateMachine', () {
    test('pendingWeekly 答对 → passWeekly', () {
      final result = StateMachine.next(ErrorQuestion.statusPendingWeekly, isCorrect: true);
      expect(result, ErrorQuestion.statusPassWeekly);
    });

    test('pendingWeekly 答错 → 保持 pendingWeekly', () {
      final result = StateMachine.next(ErrorQuestion.statusPendingWeekly, isCorrect: false);
      expect(result, ErrorQuestion.statusPendingWeekly);
    });

    test('passWeekly → pendingMonthly', () {
      final result = StateMachine.next(ErrorQuestion.statusPassWeekly, isCorrect: true);
      expect(result, ErrorQuestion.statusPendingMonthly);
    });

    test('pendingMonthly 答对 → archived', () {
      final result = StateMachine.next(ErrorQuestion.statusPendingMonthly, isCorrect: true);
      expect(result, ErrorQuestion.statusArchived);
    });

    test('pendingMonthly 答错 → 回 pendingWeekly（假性掌握）', () {
      final result = StateMachine.next(ErrorQuestion.statusPendingMonthly, isCorrect: false);
      expect(result, ErrorQuestion.statusPendingWeekly);
    });

    test('archived 任何操作 → 保持 archived', () {
      expect(StateMachine.next(ErrorQuestion.statusArchived, isCorrect: true), ErrorQuestion.statusArchived);
      expect(StateMachine.next(ErrorQuestion.statusArchived, isCorrect: false), ErrorQuestion.statusArchived);
    });
  });
}
