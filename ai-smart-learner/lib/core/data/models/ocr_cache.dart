// lib/core/data/models/ocr_cache.dart
// OCR 结果缓存
class OcrCacheEntry {
  final String id;
  final String imageHash;
  final String ocrResultJson;
  final DateTime createdAt;

  const OcrCacheEntry({
    required this.id,
    required this.imageHash,
    required this.ocrResultJson,
    required this.createdAt,
  });

  Map<String, Object?> toMap() => {
    'id': id,
    'image_hash': imageHash,
    'ocr_result_json': ocrResultJson,
    'created_at': createdAt.millisecondsSinceEpoch,
  };

  factory OcrCacheEntry.fromMap(Map<String, Object?> m) => OcrCacheEntry(
    id: m['id']! as String,
    imageHash: m['image_hash']! as String,
    ocrResultJson: m['ocr_result_json']! as String,
    createdAt: DateTime.fromMillisecondsSinceEpoch(m['created_at']! as int),
  );
}
