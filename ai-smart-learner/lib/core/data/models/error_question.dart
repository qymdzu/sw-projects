// lib/core/data/models/error_question.dart
// 错题/漏题模型
class ErrorQuestion {
  static const String typeError = 'error';
  static const String typeMissing = 'missing';

  static const String statusPendingWeekly = 'pendingWeekly';
  static const String statusPassWeekly = 'passWeekly';
  static const String statusPendingMonthly = 'pendingMonthly';
  static const String statusArchived = 'archived';

  final String id;
  final String userId;
  final String type;  // 'error' / 'missing'
  final String? imagePath;
  final String ocrText;
  final String? ocrLatex;
  final String? knowledgePoint;
  final String? errorType;
  final String? chapter;
  final String status;
  final int retryCount;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? archivedAt;

  const ErrorQuestion({
    required this.id,
    required this.userId,
    required this.type,
    this.imagePath,
    required this.ocrText,
    this.ocrLatex,
    this.knowledgePoint,
    this.errorType,
    this.chapter,
    required this.status,
    this.retryCount = 0,
    required this.createdAt,
    required this.updatedAt,
    this.archivedAt,
  });

  bool get isError => type == typeError;
  bool get isMissing => type == typeMissing;
  bool get isPendingWeekly => status == statusPendingWeekly;
  bool get isPassWeekly => status == statusPassWeekly;
  bool get isPendingMonthly => status == statusPendingMonthly;
  bool get isArchived => status == statusArchived;

  Map<String, Object?> toMap() => {
    'id': id,
    'user_id': userId,
    'type': type,
    'image_path': imagePath,
    'ocr_text': ocrText,
    'ocr_latex': ocrLatex,
    'knowledge_point': knowledgePoint,
    'error_type': errorType,
    'chapter': chapter,
    'status': status,
    'retry_count': retryCount,
    'created_at': createdAt.millisecondsSinceEpoch,
    'updated_at': updatedAt.millisecondsSinceEpoch,
    'archived_at': archivedAt?.millisecondsSinceEpoch,
  };

  factory ErrorQuestion.fromMap(Map<String, Object?> m) => ErrorQuestion(
    id: m['id']! as String,
    userId: m['user_id']! as String,
    type: m['type']! as String,
    imagePath: m['image_path'] as String?,
    ocrText: m['ocr_text']! as String,
    ocrLatex: m['ocr_latex'] as String?,
    knowledgePoint: m['knowledge_point'] as String?,
    errorType: m['error_type'] as String?,
    chapter: m['chapter'] as String?,
    status: m['status']! as String,
    retryCount: m['retry_count'] as int? ?? 0,
    createdAt: DateTime.fromMillisecondsSinceEpoch(m['created_at']! as int),
    updatedAt: DateTime.fromMillisecondsSinceEpoch(m['updated_at']! as int),
    archivedAt: m['archived_at'] != null
        ? DateTime.fromMillisecondsSinceEpoch(m['archived_at']! as int)
        : null,
  );

  ErrorQuestion copyWith({
    String? imagePath,
    String? ocrText,
    String? ocrLatex,
    String? knowledgePoint,
    String? errorType,
    String? chapter,
    String? status,
    int? retryCount,
    DateTime? updatedAt,
    DateTime? archivedAt,
  }) => ErrorQuestion(
    id: id,
    userId: userId,
    type: type,
    imagePath: imagePath ?? this.imagePath,
    ocrText: ocrText ?? this.ocrText,
    ocrLatex: ocrLatex ?? this.ocrLatex,
    knowledgePoint: knowledgePoint ?? this.knowledgePoint,
    errorType: errorType ?? this.errorType,
    chapter: chapter ?? this.chapter,
    status: status ?? this.status,
    retryCount: retryCount ?? this.retryCount,
    createdAt: createdAt,
    updatedAt: updatedAt ?? DateTime.now(),
    archivedAt: archivedAt ?? this.archivedAt,
  );
}
