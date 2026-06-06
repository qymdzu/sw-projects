// lib/core/ai/ai_provider.dart
// AI 抽象接口
abstract class AiProvider {
  String get name;
  List<String> get supportedModels;

  /// 同步调用
  Future<AiResponse> chat({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
    int maxTokens = 1000,
    Duration? timeout,
  });

  /// 流式调用
  Stream<String> chatStream({
    required List<Message> messages,
    String? model,
    double temperature = 0.7,
  });

  /// 释放
  Future<void> close();
}

class Message {
  final String role;  // 'system' / 'user' / 'assistant'
  final String content;
  const Message({required this.role, required this.content});

  Map<String, String> toJson() => {'role': role, 'content': content};
}

class AiResponse {
  final String content;
  final int promptTokens;
  final int completionTokens;
  final int totalTokens;
  final Duration latency;
  final String model;
  final String finishReason;  // 'stop' / 'length' / 'content_filter'
  final DateTime createdAt;

  const AiResponse({
    required this.content,
    required this.promptTokens,
    required this.completionTokens,
    required this.totalTokens,
    required this.latency,
    required this.model,
    required this.finishReason,
    required this.createdAt,
  });
}
