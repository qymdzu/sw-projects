// lib/ui/students/rainbow_chart_page.dart
// 知识彩虹图
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/knowledge_point.dart';
import '../../core/data/repositories/knowledge_point_repository.dart';
import '../../core/providers.dart';
import '../../config/theme.dart';
import '../shared/error_boundary.dart';

final _rainbowDataProvider = FutureProvider<RainbowData>((ref) async {
  final userId = ref.watch(currentUserIdProvider);
  if (userId == null) return RainbowData.empty();
  final repo = KnowledgePointRepository(ref.watch(databaseProvider));
  final all = await repo.findAll(userId: userId);
  return RainbowData(
    red: all.where((k) => k.mastery == KnowledgePoint.masteryRed).length,
    yellow: all.where((k) => k.mastery == KnowledgePoint.masteryYellow).length,
    green: all.where((k) => k.mastery == KnowledgePoint.masteryGreen).length,
    points: all,
  );
});

class RainbowData {
  final int red;
  final int yellow;
  final int green;
  final List<KnowledgePoint> points;

  const RainbowData({
    required this.red,
    required this.yellow,
    required this.green,
    required this.points,
  });

  factory RainbowData.empty() => const RainbowData(red: 0, yellow: 0, green: 0, points: []);

  int get total => red + yellow + green;
}

class RainbowChartPage extends ConsumerWidget {
  const RainbowChartPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final data = ref.watch(_rainbowDataProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('🌈 我的掌握')),
      body: data.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, st) => ErrorState(title: '加载失败', detail: e.toString()),
        data: (d) {
          if (d.total == 0) {
            return const EmptyState(
              emoji: '🌱',
              title: '还没有知识点',
              hint: '拍几道错题后会自动出现',
            );
          }
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                Text(
                  '你已掌握 🌟',
                  style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 4),
                Text(
                  '${d.green} 个知识点',
                  style: const TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: AppTheme.primaryLight),
                ),
                const SizedBox(height: 24),

                SizedBox(
                  height: 200,
                  child: BarChart(
                    BarChartData(
                      alignment: BarChartAlignment.spaceAround,
                      maxY: [d.red, d.yellow, d.green].reduce((a, b) => a > b ? a : b).toDouble() + 2,
                      barGroups: [
                        BarChartGroupData(x: 0, barRods: [BarChartRodData(toY: d.red.toDouble(), color: AppTheme.masteryRed, width: 30)]),
                        BarChartGroupData(x: 1, barRods: [BarChartRodData(toY: d.yellow.toDouble(), color: AppTheme.masteryYellow, width: 30)]),
                        BarChartGroupData(x: 2, barRods: [BarChartRodData(toY: d.green.toDouble(), color: AppTheme.masteryGreen, width: 30)]),
                      ],
                      titlesData: FlTitlesData(
                        leftTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                        rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                        topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                        bottomTitles: AxisTitles(
                          sideTitles: SideTitles(
                            showTitles: true,
                            reservedSize: 30,
                            getTitlesWidget: (value, meta) {
                              const labels = ['🔴 未掌握', '🟡 不稳定', '🟢 已掌握'];
                              return Padding(
                                padding: const EdgeInsets.only(top: 8),
                                child: Text(labels[value.toInt()], style: const TextStyle(fontSize: 11)),
                              );
                            },
                          ),
                        ),
                      ),
                      gridData: const FlGridData(show: false),
                      borderData: FlBorderData(show: false),
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                const Divider(),
                const SizedBox(height: 8),
                ...d.points.map((kp) => ListTile(
                      leading: Container(
                        width: 12,
                        height: 12,
                        decoration: BoxDecoration(
                          color: kp.mastery == KnowledgePoint.masteryRed
                              ? AppTheme.masteryRed
                              : kp.mastery == KnowledgePoint.masteryYellow
                                  ? AppTheme.masteryYellow
                                  : AppTheme.masteryGreen,
                          shape: BoxShape.circle,
                        ),
                      ),
                      title: Text(kp.name),
                      subtitle: Text(kp.chapter ?? ''),
                      trailing: Text('${kp.errorCount} 道错题', style: const TextStyle(fontSize: 12)),
                    )),
              ],
            ),
          );
        },
      ),
    );
  }
}
