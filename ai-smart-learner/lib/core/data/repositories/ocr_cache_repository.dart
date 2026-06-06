// lib/core/data/repositories/ocr_cache_repository.dart
// OCR 缓存仓库
import 'package:sqflite/sqflite.dart';
import 'package:uuid/uuid.dart';

import '../database.dart';
import '../models/ocr_cache.dart';

class OcrCacheRepository {
  final AppDatabase _appDb;
  final Uuid _uuid;

  OcrCacheRepository(this._appDb, [Uuid? uuid]) : _uuid = uuid ?? const Uuid();

  Database get _db => _appDb.db;

  /// 按图片 hash 查
  Future<OcrCacheEntry?> findByImageHash(String imageHash) async {
    final rows = await _db.query('ocr_cache', where: 'image_hash = ?', whereArgs: [imageHash], limit: 1);
    if (rows.isEmpty) return null;
    return OcrCacheEntry.fromMap(rows.first);
  }

  /// 保存
  Future<void> save({required String imageHash, required String ocrResultJson}) async {
    final existing = await findByImageHash(imageHash);
    if (existing != null) return;  // 已存在
    final entry = OcrCacheEntry(
      id: _uuid.v4(),
      imageHash: imageHash,
      ocrResultJson: ocrResultJson,
      createdAt: DateTime.now(),
    );
    await _db.insert('ocr_cache', entry.toMap());
  }

  /// 清过期（默认保留 30 天）
  Future<int> cleanExpired({int days = 30}) async {
    final cutoff = DateTime.now().subtract(Duration(days: days));
    return _db.delete('ocr_cache', where: 'created_at < ?', whereArgs: [cutoff.millisecondsSinceEpoch]);
  }
}
