// lib/ui/students/weekly_pool_page.dart
// 周题池列表
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/error_question.dart';
import '../../core/data/repositories/error_question_repository.dart';
import '../../core/providers.dart';
import '../../config/theme.dart';
import '../shared/error_boundary.dart';

final _weeklyPoolProvider = FutureProvider.family<List<ErrorQuestion>, String>((ref, type) async {
  final userId = ref.watch(currentUserIdProvider);
  if (userId == null) return [];
  final repo = ErrorQuestionRepository(ref.watch(databaseProvider));
  return repo.findByWeeklyPool(userId: userId, type: type);
});

class WeeklyPoolPage extends ConsumerWidget {
  const WeeklyPoolPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final pool = ref.watch(_weeklyPoolProvider('error'));

    return Scaffold(
      appBar: AppBar(title: const Text('本周错题池')),
      body: pool.when(
        data: (questions) {
          if (questions.isEmpty) {
            return const EmptyState(
              emoji: '🎉',
              title: '本周还没有错题',
              hint: '拍 1 张题照自动加入',
            );
          }
          return ListView.separated(
            padding: const EdgeInsets.all(16),
            itemCount: questions.length,
            separatorBuilder: (_, __) => const SizedBox(height: 8),
            itemBuilder: (context, idx) {
              final q = questions[idx];
              return _QuestionCard(question: q);
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, st) => ErrorState(
          title: '加载失败',
          detail: e.toString(),
          onRetry: () => ref.invalidate(_weeklyPoolProvider),
        ),
      ),
    );
  }
}

class _QuestionCard extends StatelessWidget {
  final ErrorQuestion question;
  const _QuestionCard({required this.question});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              question.ocrText,
              maxLines: 3,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontSize: 15),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                if (question.knowledgePoint != null)
                  _Tag(text: question.knowledgePoint!),
                const SizedBox(width: 4),
                if (question.errorType != null) _Tag(text: question.errorType!),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _Tag extends StatelessWidget {
  final String text;
  const _Tag({required this.text});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppTheme.primaryBackground,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(text, style: const TextStyle(fontSize: 11, color: AppTheme.primaryLight)),
    );
  }
}
