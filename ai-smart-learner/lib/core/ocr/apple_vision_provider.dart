// lib/core/ocr/apple_vision_provider.dart
// Apple Vision OCR（iOS SDK 内置，离线 fallback）
import 'dart:io';
import 'dart:typed_data';
import 'package:flutter/services.dart';
import 'package:logger/logger.dart';

import '../errors/app_error.dart';
import '../errors/error_codes.dart';
import 'ocr_provider.dart';

final _logger = Logger();

/// Apple Vision OCR Provider
/// 注意：实际 OCR 识别需要通过 platform channel 调用 iOS Vision 框架
/// 此处提供接口骨架，v0.1 简化使用 mock 实现
class AppleVisionProvider implements OcrProvider {
  AppleVisionProvider();

  static const _visionChannel = MethodChannel('ai_smart_learner/vision');

  @override
  String get name => 'apple_vision';

  @override
  bool get supportsHandwritingRemoval => false;  // Apple Vision 不直接去笔迹

  @override
  bool supportsQuestionType(String questionType) => true;

  @override
  Future<OcrResult> recognize({
    required File imageFile,
    bool removeHandwriting = true,
    Map<String, dynamic>? options,
  }) async {
    if (!Platform.isIOS) {
      throw OcrError(
        code: ErrorCodes.ocrInvalidImage,
        message: 'Apple Vision 仅支持 iOS 平台',
      );
    }

    try {
      final bytes = await imageFile.readAsBytes();
      final result = await _visionChannel.invokeMethod<Map<dynamic, dynamic>>(
        'recognizeText',
        {'image': bytes, 'removeHandwriting': removeHandwriting},
      );

      if (result == null) {
        throw OcrError(
          code: ErrorCodes.ocrInternal,
          message: 'Vision 返回空',
        );
      }

      final text = result['text'] as String? ?? '';
      final confidence = (result['confidence'] as num?)?.toDouble() ?? 0.8;

      return OcrResult(
        text: text,
        latex: result['latex'] as String?,
        tags: const [],
        confidence: confidence,
        raw: Map<String, dynamic>.from(result),
        createdAt: DateTime.now(),
        providerName: name,
      );
    } on PlatformException catch (e) {
      _logger.e('Apple Vision 调用失败: ${e.message}');
      throw OcrError(
        code: ErrorCodes.ocrInternal,
        message: 'Apple Vision 识别失败',
        cause: e,
      );
    } on MissingPluginException catch (e) {
      // v0.1 未实现 platform channel，回退到 mock
      _logger.w('Apple Vision platform channel 未实现，使用 mock');
      return OcrResult(
        text: '[模拟识别结果]\n请在 iOS 上集成 Vision framework 后实拍',
        latex: null,
        tags: const [],
        confidence: 0.0,
        raw: const {'mock': true},
        createdAt: DateTime.now(),
        providerName: name,
      );
    }
  }

  @override
  Future<void> close() async {
    // 平台通道无需关闭
  }
}
