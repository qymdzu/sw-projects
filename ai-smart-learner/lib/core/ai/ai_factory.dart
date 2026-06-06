// lib/core/ai/ai_factory.dart
// AI 工厂
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'ai_provider.dart';
import 'deepseek_provider.dart';

class AiFactory {
  final Map<String, AiProvider> _providers = {};

  AiProvider create(String name) {
    final existing = _providers[name];
    if (existing != null) return existing;

    AiProvider provider;
    switch (name) {
      case 'deepseek':
        const apiKey = String.fromEnvironment('DEEPSEEK_API_KEY', defaultValue: '');
        provider = DeepSeekProvider(apiKey: apiKey);
        break;
      default:
        throw ArgumentError('Unknown AI provider: $name');
    }
    _providers[name] = provider;
    return provider;
  }

  Future<void> closeAll() async {
    for (final p in _providers.values) {
      await p.close();
    }
    _providers.clear();
  }
}

final aiFactoryProvider = Provider<AiFactory>((ref) {
  return AiFactory();
});

/// 当前 AI provider（默认 deepseek-chat）
final currentAiProviderProvider = Provider<AiProvider>((ref) {
  final factory = ref.watch(aiFactoryProvider);
  return factory.create('deepseek');
});
