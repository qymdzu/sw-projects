// lib/ui/students/weekly_test_page.dart
// 周测答题页
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/error_question.dart';
import '../../core/data/repositories/error_question_repository.dart';
import '../../core/data/repositories/weekly_test_repository.dart';
import '../../core/providers.dart';
import '../../config/theme.dart';
import '../shared/error_boundary.dart';

class WeeklyTestPage extends ConsumerStatefulWidget {
  final String type;
  const WeeklyTestPage({required this.type, super.key});

  @override
  ConsumerState<WeeklyTestPage> createState() => _WeeklyTestPageState();
}

class _WeeklyTestPageState extends ConsumerState<WeeklyTestPage> {
  late Future<List<ErrorQuestion>> _questionsFuture;
  final Map<int, bool> _answers = {};  // questionIndex -> isCorrect
  int _currentIndex = 0;
  bool _submitting = false;

  @override
  void initState() {
    super.initState();
    _questionsFuture = _loadQuestions();
  }

  Future<List<ErrorQuestion>> _loadQuestions() async {
    final userId = ref.read(currentUserIdProvider);
    if (userId == null) return [];
    final repo = ErrorQuestionRepository(ref.read(databaseProvider));
    return repo.pickRandom(userId: userId, type: widget.type, count: 5);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('本周错题歼灭战')),
      body: FutureBuilder<List<ErrorQuestion>>(
        future: _questionsFuture,
        builder: (context, snap) {
          if (snap.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ErrorState(title: '加载失败', detail: snap.error.toString());
          }
          final questions = snap.data ?? [];
          if (questions.isEmpty) {
            return const EmptyState(emoji: '🎉', title: '本周没有题', hint: '先去拍几道错题');
          }
          return _buildQuiz(context, questions);
        },
      ),
    );
  }

  Widget _buildQuiz(BuildContext context, List<ErrorQuestion> questions) {
    final q = questions[_currentIndex];
    final isLast = _currentIndex == questions.length - 1;

    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            // 进度
            LinearProgressIndicator(
              value: (_currentIndex + 1) / questions.length,
              backgroundColor: AppTheme.divider,
              color: AppTheme.primaryLight,
            ),
            const SizedBox(height: 16),
            Text(
              '第 ${_currentIndex + 1} / ${questions.length} 题',
              style: const TextStyle(fontSize: 14, color: AppTheme.textSecondary),
            ),
            const SizedBox(height: 24),

            // 题目
            Expanded(
              child: SingleChildScrollView(
                child: Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: AppTheme.cardWhite,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    q.ocrText,
                    style: const TextStyle(fontSize: 18, height: 1.6),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 选项（v0.1 简化：让孩子在脑中答 + 答对/答错按钮）
            Text('你做对了吗？', style: const TextStyle(fontSize: 14, color: AppTheme.textSecondary)),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: _submitting ? null : () => _answer(false, questions),
                    icon: const Icon(Icons.close, color: AppTheme.textSecondary),
                    label: const Text('答错了'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: _submitting ? null : () => _answer(true, questions),
                    icon: const Icon(Icons.check),
                    label: const Text('答对了'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (!isLast)
              TextButton(
                onPressed: _submitting ? null : () => setState(() => _currentIndex++),
                child: const Text('跳过 →'),
              ),
          ],
        ),
      ),
    );
  }

  Future<void> _answer(bool isCorrect, List<ErrorQuestion> questions) async {
    _answers[_currentIndex] = isCorrect;
    if (_currentIndex < questions.length - 1) {
      setState(() => _currentIndex++);
    } else {
      // 提交
      setState(() => _submitting = true);
      try {
        final userId = ref.read(currentUserIdProvider)!;
        final testRepo = WeeklyTestRepository(ref.read(databaseProvider));
        final test = await testRepo.create(
          userId: userId,
          type: widget.type,
          weekStart: DateTime.now(),
          questionIds: questions.map((q) => q.id).toList(),
        );

        // 简化：直接跳转结果页（实际周测服务在 Stage 5 完整实现后做状态机）
        final correctCount = _answers.values.where((v) => v).length;
        if (mounted) {
          context.go('/student/test/result?correct=$correctCount&total=${questions.length}');
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('提交失败: $e')));
        }
      } finally {
        if (mounted) setState(() => _submitting = false);
      }
    }
  }
}
