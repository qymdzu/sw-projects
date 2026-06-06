// lib/core/config/config_loader.dart
// 配置加载：YAML 路由配置 + 环境变量
import 'dart:io';
import 'package:flutter/services.dart' show rootBundle;
import 'package:yaml/yaml.dart';
import 'package:logger/logger.dart';
import '../errors/app_error.dart';
import '../errors/error_codes.dart';

final _logger = Logger();

/// OCR 路由配置
class OcrRoutingConfig {
  final String defaultProvider;
  final List<OcrRoutingRule> rules;

  OcrRoutingConfig({required this.defaultProvider, required this.rules});

  static OcrRoutingConfig empty() => OcrRoutingConfig(defaultProvider: 'baidu_edu', rules: []);
}

/// OCR 路由规则
class OcrRoutingRule {
  final Map<String, dynamic> when;
  final String use;

  OcrRoutingRule({required this.when, required this.use});
}

/// AI 路由配置
class AiRoutingConfig {
  final String defaultProvider;
  final String defaultModel;
  final String? reasonerModel;
  final List<AiRoutingRule> rules;

  AiRoutingConfig({
    required this.defaultProvider,
    required this.defaultModel,
    this.reasonerModel,
    required this.rules,
  });

  static AiRoutingConfig empty() => AiRoutingConfig(
    defaultProvider: 'deepseek',
    defaultModel: 'deepseek-chat',
    reasonerModel: 'deepseek-reasoner',
    rules: [],
  );
}

/// AI 路由规则
class AiRoutingRule {
  final Map<String, dynamic> when;
  final String use;

  AiRoutingRule({required this.when, required this.use});
}

/// 配置加载器
class ConfigLoader {
  OcrRoutingConfig _ocrConfig = OcrRoutingConfig.empty();
  AiRoutingConfig _aiConfig = AiRoutingConfig.empty();
  bool _loaded = false;

  OcrRoutingConfig get ocrConfig => _ocrConfig;
  AiRoutingConfig get aiConfig => _aiConfig;
  bool get isLoaded => _loaded;

  /// 加载配置（YAML 资源 + 环境变量）
  Future<void> load() async {
    try {
      // 加载 OCR 路由
      _ocrConfig = await _loadOcrRouting();

      // 加载 AI 路由
      _aiConfig = await _loadAiRouting();

      // 环境变量覆盖
      _applyEnvOverrides();

      _loaded = true;
      _logger.i('✅ Config loaded: OCR=${_ocrConfig.defaultProvider}, AI=${_aiConfig.defaultProvider}');
    } catch (e, st) {
      _logger.w('⚠️ Config load partially failed, using defaults: $e', stackTrace: st);
      _loaded = true;  // 默认值也算加载
    }
  }

  Future<OcrRoutingConfig> _loadOcrRouting() async {
    try {
      final yamlStr = await rootBundle.loadString('assets/configs/ocr_routing.yaml');
      return _parseOcrRouting(yamlStr);
    } catch (e) {
      // 资源不存在时用默认
      return OcrRoutingConfig.empty();
    }
  }

  OcrRoutingConfig _parseOcrRouting(String yamlStr) {
    final yaml = loadYaml(yamlStr);
    final defaultProvider = yaml['default'] as String? ?? 'baidu_edu';
    final rulesList = (yaml['rules'] as List?) ?? [];
    final rules = rulesList.map((r) {
      final m = r as Map;
      return OcrRoutingRule(
        when: Map<String, dynamic>.from(m['when'] as Map),
        use: m['use'] as String,
      );
    }).toList();
    return OcrRoutingConfig(defaultProvider: defaultProvider, rules: rules);
  }

  Future<AiRoutingConfig> _loadAiRouting() async {
    try {
      final yamlStr = await rootBundle.loadString('assets/configs/ai_routing.yaml');
      return _parseAiRouting(yamlStr);
    } catch (e) {
      return AiRoutingConfig.empty();
    }
  }

  AiRoutingConfig _parseAiRouting(String yamlStr) {
    final yaml = loadYaml(yamlStr);
    final defaultProvider = yaml['default'] as String? ?? 'deepseek';
    final models = yaml['models'] as Map?;
    final defaultModel = models?['chat'] as String? ?? 'deepseek-chat';
    final reasonerModel = models?['reasoner'] as String?;
    final rulesList = (yaml['rules'] as List?) ?? [];
    final rules = rulesList.map((r) {
      final m = r as Map;
      return AiRoutingRule(
        when: Map<String, dynamic>.from(m['when'] as Map),
        use: m['use'] as String,
      );
    }).toList();
    return AiRoutingConfig(
      defaultProvider: defaultProvider,
      defaultModel: defaultModel,
      reasonerModel: reasonerModel,
      rules: rules,
    );
  }

  void _applyEnvOverrides() {
    // 从环境变量覆盖（v0.1 简化：仅支持 .env 文件）
    final envFile = File('.env');
    if (!envFile.existsSync()) return;

    final lines = envFile.readAsLinesSync();
    final env = <String, String>{};
    for (final line in lines) {
      final trimmed = line.trim();
      if (trimmed.isEmpty || trimmed.startsWith('#')) continue;
      final eq = trimmed.indexOf('=');
      if (eq > 0) {
        env[trimmed.substring(0, eq).trim()] = trimmed.substring(eq + 1).trim();
      }
    }

    if (env['OCR_DEFAULT'] != null) {
      _ocrConfig = OcrRoutingConfig(
        defaultProvider: env['OCR_DEFAULT']!,
        rules: _ocrConfig.rules,
      );
    }
    if (env['AI_DEFAULT'] != null) {
      _aiConfig = AiRoutingConfig(
        defaultProvider: env['AI_DEFAULT']!,
        defaultModel: _aiConfig.defaultModel,
        reasonerModel: _aiConfig.reasonerModel,
        rules: _aiConfig.rules,
      );
    }
  }
}
