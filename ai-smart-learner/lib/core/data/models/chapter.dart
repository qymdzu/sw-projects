// lib/core/data/models/chapter.dart
// 教材章节模型
class Chapter {
  final String id;
  final String name;
  final String subject;
  final String stage;  // 'primary' / 'junior' / 'senior'
  final String? parentId;
  final int? displayOrder;

  const Chapter({
    required this.id,
    required this.name,
    required this.subject,
    required this.stage,
    this.parentId,
    this.displayOrder,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'name': name,
    'subject': subject,
    'stage': stage,
    'parent_id': parentId,
    'display_order': displayOrder,
  };

  factory Chapter.fromMap(Map<String, Object?> m) => Chapter(
    id: m['id']! as String,
    name: m['name']! as String,
    subject: m['subject']! as String,
    stage: m['stage']! as String,
    parentId: m['parent_id'] as String?,
    displayOrder: m['display_order'] as int?,
  );
}
