// lib/core/data/models/variant_question.dart
// 变式题模型
class VariantQuestion {
  final String id;
  final String parentQuestionId;
  final String text;
  final String? latex;
  final int variantIndex;
  final String status;
  final DateTime createdAt;

  const VariantQuestion({
    required this.id,
    required this.parentQuestionId,
    required this.text,
    this.latex,
    required this.variantIndex,
    required this.status,
    required this.createdAt,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'parent_question_id': parentQuestionId,
    'text': text,
    'latex': latex,
    'variant_index': variantIndex,
    'status': status,
    'created_at': createdAt.millisecondsSinceEpoch,
  };

  factory VariantQuestion.fromMap(Map<String, Object?> m) => VariantQuestion(
    id: m['id']! as String,
    parentQuestionId: m['parent_question_id']! as String,
    text: m['text']! as String,
    latex: m['latex'] as String?,
    variantIndex: m['variant_index']! as int,
    status: m['status']! as String,
    createdAt: DateTime.fromMillisecondsSinceEpoch(m['created_at']! as int),
  );

  VariantQuestion copyWith({String? status}) => VariantQuestion(
    id: id,
    parentQuestionId: parentQuestionId,
    text: text,
    latex: latex,
    variantIndex: variantIndex,
    status: status ?? this.status,
    createdAt: createdAt,
  );
}
