// lib/ui/students/weekly_test_result_page.dart
// 周测结果页
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../config/theme.dart';

class WeeklyTestResultPage extends StatelessWidget {
  final int correctCount;
  final int totalCount;

  const WeeklyTestResultPage({
    required this.correctCount,
    required this.totalCount,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    final wrong = totalCount - correctCount;
    final variantsCreated = wrong * 3;

    return Scaffold(
      appBar: AppBar(title: const Text('周测完成')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              const SizedBox(height: 40),
              const Text('🎉', style: TextStyle(fontSize: 80)),
              const SizedBox(height: 16),
              Text(
                '答对 $correctCount / $totalCount',
                style: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 32),
              if (wrong > 0)
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: AppTheme.primaryBackground,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Column(
                    children: [
                      const Text('💪 已自动生成变式题', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                      const SizedBox(height: 8),
                      Text(
                        '$variantsCreated 道新题已加入下周错题池',
                        style: const TextStyle(color: AppTheme.textSecondary),
                      ),
                    ],
                  ),
                ),
              const Spacer(),
              ElevatedButton(
                onPressed: () => context.go('/student'),
                child: const Text('返回首页'),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }
}
