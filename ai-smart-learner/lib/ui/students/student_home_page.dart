// lib/ui/students/student_home_page.dart
// 学生首页
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/error_question.dart';
import '../../core/data/repositories/error_question_repository.dart';
import '../../core/providers.dart';
import '../../config/theme.dart';

final weeklyPoolCountProvider = FutureProvider.family<int, String>((ref, type) async {
  final userId = ref.watch(currentUserIdProvider);
  if (userId == null) return 0;
  final repo = ref.watch(_errorQuestionRepoProvider);
  final pool = await repo.findByWeeklyPool(userId: userId, type: type);
  return pool.length;
});

final _errorQuestionRepoProvider = Provider<ErrorQuestionRepository>((ref) {
  return ErrorQuestionRepository(ref.watch(databaseProvider));
});

class StudentHomePage extends ConsumerWidget {
  const StudentHomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(currentUserProvider);
    final errorCount = ref.watch(weeklyPoolCountProvider('error'));
    final missingCount = ref.watch(weeklyPoolCountProvider('missing'));

    return Scaffold(
      appBar: AppBar(
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('🎒', style: TextStyle(fontSize: 22)),
            const SizedBox(width: 8),
            Text(user?.name ?? '学生'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.rainbow),
            tooltip: '我的掌握',
            onPressed: () => context.push('/student/chart'),
          ),
          IconButton(
            icon: const Icon(Icons.settings),
            tooltip: '设置',
            onPressed: () => context.push('/settings'),
          ),
        ],
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '你好 ${user?.name ?? '同学'}',
                style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: AppTheme.textPrimary),
              ),
              const SizedBox(height: 8),
              Text(
                '本周新增 ${errorCount.value ?? 0} 道错题 · 待周测 ${errorCount.value ?? 0} 道',
                style: const TextStyle(fontSize: 14, color: AppTheme.textSecondary),
              ),
              const SizedBox(height: 32),

              // 拍照主按钮
              ElevatedButton.icon(
                onPressed: () => context.push('/student/capture?type=error'),
                icon: const Text('📷', style: TextStyle(fontSize: 24)),
                label: const Text('拍错题（或拍漏题）'),
              ),
              const SizedBox(height: 24),

              // 快捷入口
              _MenuCard(
                emoji: '📚',
                title: '本周错题池',
                count: errorCount.value ?? 0,
                onTap: () => context.push('/student/pool'),
              ),
              const SizedBox(height: 12),
              _MenuCard(
                emoji: '📝',
                title: '本周周测（错题）',
                count: errorCount.value ?? 0,
                onTap: () => context.push('/student/test?type=error'),
              ),
              const SizedBox(height: 12),
              _MenuCard(
                emoji: '📝',
                title: '本周周测（漏题）',
                count: missingCount.value ?? 0,
                onTap: () => context.push('/student/test?type=missing'),
              ),
              const SizedBox(height: 12),
              _MenuCard(
                emoji: '🌈',
                title: '知识彩虹图',
                onTap: () => context.push('/student/chart'),
              ),
              const SizedBox(height: 12),
              _MenuCard(
                emoji: '📦',
                title: '已归档',
                onTap: () => context.push('/student/archive'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _MenuCard extends StatelessWidget {
  final String emoji;
  final String title;
  final int? count;
  final VoidCallback onTap;

  const _MenuCard({required this.emoji, required this.title, this.count, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Text(emoji, style: const TextStyle(fontSize: 28)),
              const SizedBox(width: 16),
              Expanded(
                child: Text(
                  title,
                  style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
                ),
              ),
              if (count != null) ...[
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppTheme.primaryBackground,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    '$count',
                    style: const TextStyle(color: AppTheme.primaryLight, fontWeight: FontWeight.bold),
                  ),
                ),
                const SizedBox(width: 8),
              ],
              const Icon(Icons.arrow_forward_ios, size: 16, color: AppTheme.textSecondary),
            ],
          ),
        ),
      ),
    );
  }
}
