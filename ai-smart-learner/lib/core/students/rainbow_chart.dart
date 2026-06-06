// lib/core/students/rainbow_chart.dart
// 知识彩虹图数据
import 'package:logger/logger.dart';

import '../data/models/knowledge_point.dart';
import '../data/repositories/knowledge_point_repository.dart';

final _logger = Logger();

class RainbowChartData {
  final int redCount;
  final int yellowCount;
  final int greenCount;
  final List<KnowledgePoint> redPoints;
  final List<KnowledgePoint> yellowPoints;
  final List<KnowledgePoint> greenPoints;

  const RainbowChartData({
    required this.redCount,
    required this.yellowCount,
    required this.greenCount,
    this.redPoints = const [],
    this.yellowPoints = const [],
    this.greenPoints = const [],
  });

  int get totalCount => redCount + yellowCount + greenCount;
}

class RainbowChartService {
  final KnowledgePointRepository _kpRepo;

  RainbowChartService(this._kpRepo);

  /// 查用户的彩虹图数据
  Future<RainbowChartData> loadForUser({required String userId}) async {
    final all = await _kpRepo.findAll(userId: userId);
    final red = all.where((k) => k.mastery == KnowledgePoint.masteryRed).toList();
    final yellow = all.where((k) => k.mastery == KnowledgePoint.masteryYellow).toList();
    final green = all.where((k) => k.mastery == KnowledgePoint.masteryGreen).toList();

    return RainbowChartData(
      redCount: red.length,
      yellowCount: yellow.length,
      greenCount: green.length,
      redPoints: red,
      yellowPoints: yellow,
      greenPoints: green,
    );
  }
}
