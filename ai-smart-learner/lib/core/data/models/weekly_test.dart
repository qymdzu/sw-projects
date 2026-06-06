// lib/core/data/models/weekly_test.dart
// 周测记录模型
class WeeklyTest {
  final String id;
  final String userId;
  final String type;  // 'error' / 'missing'
  final DateTime weekStart;
  final List<String> questionIds;
  final DateTime? completedAt;
  final int? correctCount;

  const WeeklyTest({
    required this.id,
    required this.userId,
    required this.type,
    required this.weekStart,
    required this.questionIds,
    this.completedAt,
    this.correctCount,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'user_id': userId,
    'type': type,
    'week_start': weekStart.millisecondsSinceEpoch,
    'question_ids': questionIds.join(','),
    'completed_at': completedAt?.millisecondsSinceEpoch,
    'correct_count': correctCount,
  };

  factory WeeklyTest.fromMap(Map<String, Object?> m) => WeeklyTest(
    id: m['id']! as String,
    userId: m['user_id']! as String,
    type: m['type']! as String,
    weekStart: DateTime.fromMillisecondsSinceEpoch(m['week_start']! as int),
    questionIds: (m['question_ids']! as String).split(',').where((s) => s.isNotEmpty).toList(),
    completedAt: m['completed_at'] != null
        ? DateTime.fromMillisecondsSinceEpoch(m['completed_at']! as int)
        : null,
    correctCount: m['correct_count'] as int?,
  );

  WeeklyTest copyWith({DateTime? completedAt, int? correctCount}) => WeeklyTest(
    id: id,
    userId: userId,
    type: type,
    weekStart: weekStart,
    questionIds: questionIds,
    completedAt: completedAt ?? this.completedAt,
    correctCount: correctCount ?? this.correctCount,
  );
}
