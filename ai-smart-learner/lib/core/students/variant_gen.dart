// lib/core/students/variant_gen.dart
// 变式题生成
import 'package:logger/logger.dart';

import '../ai/ai_provider.dart';
import '../ai/prompts/variant_prompt.dart';
import '../data/models/error_question.dart';
import '../data/models/variant_question.dart';
import '../data/repositories/ai_call_log_repository.dart';
import '../data/repositories/variant_question_repository.dart';

final _logger = Logger();

class VariantGenService {
  final AiProvider _ai;
  final VariantQuestionRepository _variantRepo;
  final AiCallLogRepository _logRepo;

  VariantGenService(this._ai, this._variantRepo, this._logRepo);

  /// 生成 3 道变式题
  Future<List<VariantQuestion>> generate({
    required String userId,
    required ErrorQuestion parent,
  }) async {
    final messages = VariantPrompt.build(
      originalQuestion: parent.ocrText,
      knowledgePoint: parent.knowledgePoint,
      difficulty: parent.errorType,
    );

    final start = DateTime.now();
    try {
      final response = await _ai.chat(
        messages: messages,
        temperature: 0.7,
        maxTokens: 800,
      );

      await _logRepo.log(
        userId: userId,
        provider: _ai.name,
        purpose: 'variant',
        promptTokens: response.promptTokens,
        completionTokens: response.completionTokens,
        totalTokens: response.totalTokens,
        latencyMs: response.latency.inMilliseconds,
        success: true,
      );

      final texts = VariantPrompt.parseVariants(response.content);
      if (texts.isEmpty) {
        _logger.w('变式题解析失败: ${response.content}');
        return [];
      }

      // 取前 3 道保存
      final saved = await _variantRepo.saveBatch(
        parentQuestionId: parent.id,
        texts: texts.take(3).toList(),
        latex: null,
      );
      _logger.i('生成了 ${saved.length} 道变式题（父题: ${parent.id}）');
      return saved;
    } catch (e, st) {
      _logger.e('变式题生成失败: $e', stackTrace: st);
      await _logRepo.log(
        userId: userId,
        provider: _ai.name,
        purpose: 'variant',
        success: false,
        errorMsg: e.toString(),
      );
      return [];
    }
  }
}
