// lib/core/students/weekly_test.dart
// 周测：抽题 + 答题 + 提交
import 'package:logger/logger.dart';

import '../data/models/error_question.dart';
import '../data/models/knowledge_point.dart';
import '../data/models/variant_question.dart';
import '../data/models/weekly_test.dart';
import '../data/repositories/error_question_repository.dart';
import '../data/repositories/knowledge_point_repository.dart';
import '../data/repositories/variant_question_repository.dart';
import '../data/repositories/weekly_test_repository.dart';
import 'state_machine.dart';
import 'variant_gen.dart';

final _logger = Logger();

class WeeklyTestResult {
  final int correctCount;
  final int totalCount;
  final List<VariantQuestion> newVariants;
  final List<ErrorQuestion> promoted;

  const WeeklyTestResult({
    required this.correctCount,
    required this.totalCount,
    this.newVariants = const [],
    this.promoted = const [],
  });
}

class WeeklyTestService {
  final ErrorQuestionRepository _eqRepo;
  final KnowledgePointRepository _kpRepo;
  final VariantQuestionRepository _variantRepo;
  final WeeklyTestRepository _testRepo;
  final VariantGenService _variantGen;

  WeeklyTestService(this._eqRepo, this._kpRepo, this._variantRepo, this._testRepo, this._variantGen);

  /// 抽 5 道题（默认）
  Future<List<ErrorQuestion>> pickQuestions({required String userId, required String type, int count = 5}) async {
    return _eqRepo.pickRandom(userId: userId, type: type, count: count);
  }

  /// 创建周测记录
  Future<WeeklyTest> createTest({
    required String userId,
    required String type,
    required List<String> questionIds,
  }) async {
    return _testRepo.create(
      userId: userId,
      type: type,
      weekStart: _weekStart(),
      questionIds: questionIds,
    );
  }

  /// 提交答案
  Future<WeeklyTestResult> submit({
    required String userId,
    required WeeklyTest test,
    required Map<String, bool> answers,  // questionId -> isCorrect
  }) async {
    int correct = 0;
    final newVariants = <VariantQuestion>[];
    final promoted = <ErrorQuestion>[];

    for (final questionId in test.questionIds) {
      final isCorrect = answers[questionId] ?? false;
      final question = await _eqRepo.findById(questionId);
      if (question == null) continue;

      if (isCorrect) {
        correct++;
        // 状态机：pendingWeekly → passWeekly
        final newStatus = StateMachine.next(question.status, isCorrect: true);
        await _eqRepo.updateStatus(id: question.id, status: newStatus);
        if (newStatus == ErrorQuestion.statusPassWeekly) {
          promoted.add(question.copyWith(status: newStatus));
        }
        // 知识点 → 升到 yellow
        if (question.knowledgePoint != null) {
          final kp = await _kpRepo.findByName(
            userId: userId,
            name: question.knowledgePoint!,
            chapter: question.chapter,
          );
          if (kp != null) {
            await _kpRepo.updateMastery(id: kp.id, mastery: KnowledgePoint.masteryYellow);
          }
        }
      } else {
        // 答错：状态保持 pendingWeekly
        await _eqRepo.incrementRetryCount(question.id);
        // 生成 3 道变式题
        final variants = await _variantGen.generate(userId: userId, parent: question);
        newVariants.addAll(variants);
      }
    }

    // 更新周测记录
    await _testRepo.complete(testId: test.id, correctCount: correct);

    return WeeklyTestResult(
      correctCount: correct,
      totalCount: test.questionIds.length,
      newVariants: newVariants,
      promoted: promoted,
    );
  }

  /// 本周开始时间（周一 00:00）
  DateTime _weekStart() {
    final now = DateTime.now();
    final monday = now.subtract(Duration(days: now.weekday - 1));
    return DateTime(monday.year, monday.month, monday.day);
  }
}
