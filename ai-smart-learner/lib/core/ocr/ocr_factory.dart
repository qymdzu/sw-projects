// lib/core/ocr/ocr_factory.dart
// OCR 工厂：按配置创建 provider
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../errors/app_error.dart';
import '../errors/error_codes.dart';
import 'apple_vision_provider.dart';
import 'baidu_edu_provider.dart';
import 'ocr_provider.dart';

class OcrFactory {
  final Map<String, OcrProvider> _providers = {};

  /// 按名称创建 provider
  OcrProvider create(String name) {
    final existing = _providers[name];
    if (existing != null) return existing;

    OcrProvider provider;
    switch (name) {
      case 'baidu_edu':
        // v0.1 简化：从环境变量读，v0.2 用 Keychain
        const apiKey = String.fromEnvironment('BAIDU_OCR_API_KEY', defaultValue: '');
        const secretKey = String.fromEnvironment('BAIDU_OCR_SECRET_KEY', defaultValue: '');
        provider = BaiduEduProvider(apiKey: apiKey, secretKey: secretKey);
        break;
      case 'apple_vision':
        provider = AppleVisionProvider();
        break;
      default:
        throw ArgumentError('Unknown OCR provider: $name');
    }
    _providers[name] = provider;
    return provider;
  }

  /// 释放所有
  Future<void> closeAll() async {
    for (final p in _providers.values) {
      await p.close();
    }
    _providers.clear();
  }
}

/// OcrFactory provider（单例）
final ocrFactoryProvider = Provider<OcrFactory>((ref) {
  return OcrFactory();
});

/// 当前 OCR provider（按配置动态选）
final currentOcrProviderProvider = FutureProvider<OcrProvider>((ref) async {
  final config = ref.watch(ocrConfigProvider);
  final factory = ref.watch(ocrFactoryProvider);
  return factory.create(config.defaultProvider);
});

/// OCR 配置 provider（用 default 即可，复杂规则 v0.2 推）
final ocrConfigProvider = Provider((ref) {
  // v0.1 默认 baidu_edu；fallback 由调用方 try-catch 处理
  return _OcrConfigHolder(defaultProvider: 'baidu_edu');
});

class _OcrConfigHolder {
  final String defaultProvider;
  _OcrConfigHolder({required this.defaultProvider});
}
