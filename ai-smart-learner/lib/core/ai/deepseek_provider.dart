// lib/core/ai/deepseek_provider.dart
// DeepSeek API 实现
import 'dart:async';
import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:logger/logger.dart';

import '../errors/app_error.dart';
import '../errors/error_codes.dart';
import 'ai_provider.dart';

final _logger = Logger();

class DeepSeekProvider implements AiProvider {
  final Dio _dio;
  final String _apiKey;

  DeepSeekProvider({required String apiKey, Dio? dio})
      : _apiKey = apiKey,
        _dio = dio ?? Dio() {
    _dio.options.baseUrl = 'https://api.deepseek.com';
    _dio.options.connectTimeout = const Duration(seconds: 10);
    _dio.options.receiveTimeout = const Duration(seconds: 60);
  }

  @override
  String get name => 'deepseek';

  @override
  List<String> get supportedModels => ['deepseek-chat', 'deepseek-reasoner'];

  @override
  Future<AiResponse> chat({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
    int maxTokens = 1000,
    Duration? timeout,
  }) async {
    if (_apiKey.isEmpty) {
      throw AiError(
        code: ErrorCodes.aiAuth,
        message: 'DeepSeek API Key 未配置',
      );
    }

    final start = DateTime.now();
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/chat/completions',
        data: {
          'model': model ?? 'deepseek-chat',
          'messages': messages.map((m) => m.toJson()).toList(),
          'temperature': temperature,
          'max_tokens': maxTokens,
          'stream': false,
        },
        options: Options(
          headers: {
            'Authorization': 'Bearer $_apiKey',
            'Content-Type': 'application/json',
          },
          receiveTimeout: timeout,
        ),
      );

      final data = response.data;
      if (data == null) {
        throw AiError(
          code: ErrorCodes.aiInternal,
          message: 'DeepSeek 返回空',
        );
      }

      final choice = (data['choices'] as List?)?.firstOrNull;
      if (choice == null) {
        throw AiError(
          code: ErrorCodes.aiContent,
          message: 'DeepSeek 返回无 choices',
          detail: data.toString(),
        );
      }

      final message = (choice as Map)['message'] as Map;
      final content = message['content'] as String? ?? '';
      final finishReason = choice['finish_reason'] as String? ?? 'stop';
      final usage = data['usage'] as Map? ?? {};
      final modelName = data['model'] as String? ?? (model ?? 'deepseek-chat');

      return AiResponse(
        content: content,
        promptTokens: usage['prompt_tokens'] as int? ?? 0,
        completionTokens: usage['completion_tokens'] as int? ?? 0,
        totalTokens: usage['total_tokens'] as int? ?? 0,
        latency: DateTime.now().difference(start),
        model: modelName,
        finishReason: finishReason,
        createdAt: DateTime.now(),
      );
    } on DioException catch (e) {
      _logger.e('DeepSeek 调用失败: ${e.message}');
      if (e.type == DioExceptionType.connectionTimeout || e.type == DioExceptionType.receiveTimeout) {
        throw NetworkError(
          code: ErrorCodes.netTimeout,
          message: 'AI 服务超时',
          cause: e,
        );
      }
      final code = e.response?.statusCode ?? 0;
      if (code == 401) {
        throw AiError(code: ErrorCodes.aiAuth, message: 'API Key 无效', cause: e);
      }
      if (code == 429) {
        throw AiError(code: ErrorCodes.aiRate, message: 'AI 调用过快', cause: e);
      }
      throw AiError(code: ErrorCodes.aiInternal, message: 'AI 服务异常', cause: e);
    } catch (e) {
      _logger.e('DeepSeek 异常: $e');
      throw AiError(code: ErrorCodes.aiInternal, message: 'AI 调用失败', cause: e);
    }
  }

  @override
  Stream<String> chatStream({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
  }) async* {
    // v0.1 简化：不实现流式（v0.2 推）
    final response = await chat(messages: messages, model: model, temperature: temperature);
    yield response.content;
  }

  @override
  Future<void> close() async {
    _dio.close();
  }
}

extension _ListExt<T> on Iterable<T> {
  T? get firstOrNull => isEmpty ? null : first;
}
