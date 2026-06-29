// lib/core/students/tagging.dart
// 智能打标：调 AI 给题目打标签
import 'package:logger/logger.dart';

import '../ai/ai_provider.dart';
import '../ai/ai_factory.dart';
import '../ai/prompts/tagging_prompt.dart';
import '../data/models/error_question.dart';
import '../data/repositories/ai_call_log_repository.dart';
import '../errors/app_error.dart';
import '../errors/error_codes.dart';

final _logger = Logger();

class TaggingResult {
  final String? knowledgePoint;
  final String? errorType;
  final String? subject;
  final String? chapter;
  final String? difficulty;

  const TaggingResult({
    this.knowledgePoint,
    this.errorType,
    this.subject,
    this.chapter,
    this.difficulty,
  });
}

class TaggingService {
  final AiProvider _ai;
  final AiCallLogRepository _logRepo;

  TaggingService(this._ai, this._logRepo);

  /// 给题目打标
  Future<TaggingResult> tag({
    required String userId,
    required String questionText,
    String? subjectHint,
    String? stage,
  }) async {
    final messages = TaggingPrompt.build(
      questionText: questionText,
      subjectHint: subjectHint,
      stage: stage,
    );

    final start = DateTime.now();
    try {
      final response = await _ai.chat(
        messages: messages,
        temperature: 0.3,
        maxTokens: 200,
      );

      await _logRepo.log(
        userId: userId,
        provider: _ai.name,
        purpose: 'tagging',
        promptTokens: response.promptTokens,
        completionTokens: response.completionTokens,
        totalTokens: response.totalTokens,
        costCents: _estimateCostCents(response),
        latencyMs: response.latency.inMilliseconds,
        success: true,
      );

      final parsed = TaggingPrompt.parseJson(response.content);
      if (parsed == null) {
        _logger.w('AI 打标返回解析失败: ${response.content}');
        return const TaggingResult();  // 失败不阻塞
      }

      return TaggingResult(
        knowledgePoint: parsed['knowledgePoint'] as String?,
        errorType: parsed['errorType'] as String?,
        subject: parsed['subject'] as String?,
        chapter: parsed['chapter'] as String?,
        difficulty: parsed['difficulty'] as String?,
      );
    } catch (e, st) {
      _logger.e('打标失败: $e', stackTrace: st);
      await _logRepo.log(
        userId: userId,
        provider: _ai.name,
        purpose: 'tagging',
        success: false,
        errorMsg: e.toString(),
      );
      // 打标失败不抛错（业务允许无标签入库）
      return const TaggingResult();
    }
  }

  /// DeepSeek 估算成本（按 token 数，¥1/百万 tokens → 0.1 分/千 token）
  int _estimateCostCents(AiResponse r) {
    // 1 元 = 100 分，1 百万 token = 1 元 = 100 分，0.0001 分/token
    return ((r.totalTokens * 0.0001)).round();
  }
}
