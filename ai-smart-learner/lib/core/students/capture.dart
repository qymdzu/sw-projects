// lib/core/students/capture.dart
// 拍照录入流程（核心业务）
import 'dart:io';
import 'package:logger/logger.dart';

import '../data/models/error_question.dart';
import '../data/models/knowledge_point.dart';
import '../data/repositories/error_question_repository.dart';
import '../data/repositories/knowledge_point_repository.dart';
import '../errors/app_error.dart';
import '../errors/error_codes.dart';
import '../ocr/ocr_provider.dart';
import 'tagging.dart';

final _logger = Logger();

class CaptureService {
  final ErrorQuestionRepository _eqRepo;
  final KnowledgePointRepository _kpRepo;
  final TaggingService _tagging;
  final OcrProvider _ocr;

  CaptureService(this._eqRepo, this._kpRepo, this._tagging, this._ocr);

  /// 拍照录入完整流程：
  /// 1. 调 OCR 识别
  /// 2. 调 AI 打标
  /// 3. 入库
  /// 4. 更新知识点状态
  /// 5. 默认删除原图
  Future<ErrorQuestion> capture({
    required String userId,
    required String type,  // 'error' / 'missing'
    required File imageFile,
  }) async {
    // Step 1: OCR
    OcrResult ocrResult;
    try {
      ocrResult = await _ocr.recognize(imageFile: imageFile, removeHandwriting: true);
    } on AppError {
      rethrow;  // OCR 错误往上抛
    } catch (e) {
      throw OcrError(code: ErrorCodes.ocrInternal, message: 'OCR 识别失败', cause: e);
    }

    if (ocrResult.text.trim().isEmpty) {
      throw OcrError(
        code: ErrorCodes.ocrInvalidImage,
        message: 'OCR 识别结果为空，请手动输入或重新拍照',
      );
    }

    // Step 2: AI 打标
    final tagging = await _tagging.tag(
      userId: userId,
      questionText: ocrResult.text,
    );

    // Step 3: 入库
    final question = await _eqRepo.save(
      userId: userId,
      type: type,
      imagePath: imageFile.path,  // 先存路径
      ocrText: ocrResult.text,
      ocrLatex: ocrResult.latex,
      knowledgePoint: tagging.knowledgePoint,
      errorType: tagging.errorType,
      chapter: tagging.chapter,
    );

    // Step 4: 更新知识点（如果打了标）
    if (tagging.knowledgePoint != null) {
      final kp = await _kpRepo.upsert(
        userId: userId,
        name: tagging.knowledgePoint!,
        chapter: tagging.chapter,
        mastery: KnowledgePoint.masteryRed,
      );
      await _kpRepo.incrementErrorCount(id: kp.id, by: 1);
    }

    // Step 5: 默认删除原图（v0.1 NFR-011 隐私要求）
    try {
      if (imageFile.existsSync()) {
        await imageFile.delete();
        await _eqRepo.clearImagePath(question.id);
        _logger.i('✅ 已删除原图（默认行为）');
      }
    } catch (e) {
      _logger.w('删除原图失败（非阻塞）: $e');
    }

    return question;
  }
}
