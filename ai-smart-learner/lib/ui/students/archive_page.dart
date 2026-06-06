// lib/ui/students/archive_page.dart
// 已归档页
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/error_question.dart';
import '../../core/data/repositories/error_question_repository.dart';
import '../../core/providers.dart';
import '../shared/error_boundary.dart';

final _archivedProvider = FutureProvider<List<ErrorQuestion>>((ref) async {
  final userId = ref.watch(currentUserIdProvider);
  if (userId == null) return [];
  final repo = ErrorQuestionRepository(ref.watch(databaseProvider));
  return repo.findArchived(userId: userId);
});

class ArchivePage extends ConsumerWidget {
  const ArchivePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final archived = ref.watch(_archivedProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('📦 已归档')),
      body: archived.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, st) => ErrorState(title: '加载失败', detail: e.toString()),
        data: (questions) {
          if (questions.isEmpty) {
            return const EmptyState(emoji: '🌟', title: '还没有归档题目', hint: '答对月测后会归档');
          }
          return ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: questions.length,
            itemBuilder: (ctx, i) {
              final q = questions[i];
              return Card(
                child: ListTile(
                  title: Text(q.ocrText, maxLines: 2, overflow: TextOverflow.ellipsis),
                  subtitle: Text('${q.knowledgePoint ?? ''} · ${q.archivedAt?.toString().split(' ').first ?? ''}'),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
