// lib/core/data/models/knowledge_point.dart
// 知识点模型
class KnowledgePoint {
  static const String masteryRed = 'red';     // 未掌握
  static const String masteryYellow = 'yellow'; // 不稳定
  static const String masteryGreen = 'green';   // 已掌握

  final String id;
  final String userId;
  final String name;
  final String? chapter;
  final String mastery;
  final int errorCount;
  final DateTime? lastErrorAt;
  final DateTime updatedAt;

  const KnowledgePoint({
    required this.id,
    required this.userId,
    required this.name,
    this.chapter,
    required this.mastery,
    this.errorCount = 0,
    this.lastErrorAt,
    required this.updatedAt,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'user_id': userId,
    'name': name,
    'chapter': chapter,
    'mastery': mastery,
    'error_count': errorCount,
    'last_error_at': lastErrorAt?.millisecondsSinceEpoch,
    'updated_at': updatedAt.millisecondsSinceEpoch,
  };

  factory KnowledgePoint.fromMap(Map<String, Object?> m) => KnowledgePoint(
    id: m['id']! as String,
    userId: m['user_id']! as String,
    name: m['name']! as String,
    chapter: m['chapter'] as String?,
    mastery: m['mastery']! as String,
    errorCount: m['error_count'] as int? ?? 0,
    lastErrorAt: m['last_error_at'] != null
        ? DateTime.fromMillisecondsSinceEpoch(m['last_error_at']! as int)
        : null,
    updatedAt: DateTime.fromMillisecondsSinceEpoch(m['updated_at']! as int),
  );

  KnowledgePoint copyWith({
    String? mastery,
    int? errorCount,
    DateTime? lastErrorAt,
    DateTime? updatedAt,
  }) => KnowledgePoint(
    id: id,
    userId: userId,
    name: name,
    chapter: chapter,
    mastery: mastery ?? this.mastery,
    errorCount: errorCount ?? this.errorCount,
    lastErrorAt: lastErrorAt ?? this.lastErrorAt,
    updatedAt: updatedAt ?? DateTime.now(),
  );
}
