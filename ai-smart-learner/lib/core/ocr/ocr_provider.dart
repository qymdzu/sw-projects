// lib/core/ocr/ocr_provider.dart
// OCR 抽象接口
import 'dart:io';

abstract class OcrProvider {
  String get name;
  bool get supportsHandwritingRemoval;
  bool supportsQuestionType(String questionType);

  /// 识别图片
  Future<OcrResult> recognize({
    required File imageFile,
    bool removeHandwriting = true,
    Map<String, dynamic>? options,
  });

  /// 释放资源
  Future<void> close();
}

class OcrResult {
  final String text;                      // 纯文本
  final String? latex;                    // LaTeX 公式
  final List<String> tags;                // 关键词
  final double confidence;                // 0-1
  final Map<String, dynamic> raw;         // 原始响应
  final DateTime createdAt;
  final String providerName;              // 哪个 provider 出的

  const OcrResult({
    required this.text,
    this.latex,
    this.tags = const [],
    this.confidence = 1.0,
    this.raw = const {},
    required this.createdAt,
    required this.providerName,
  });

  Map<String, Object?> toMap() => {
    'text': text,
    'latex': latex,
    'tags': tags,
    'confidence': confidence,
    'raw': raw,
    'provider': providerName,
  };

  factory OcrResult.fromMap(Map<String, Object?> m) => OcrResult(
    text: m['text']! as String,
    latex: m['latex'] as String?,
    tags: (m['tags'] as List?)?.cast<String>() ?? [],
    confidence: (m['confidence'] as num?)?.toDouble() ?? 1.0,
    raw: Map<String, dynamic>.from(m['raw'] as Map? ?? {}),
    createdAt: DateTime.now(),
    providerName: m['provider'] as String? ?? 'unknown',
  );
}
