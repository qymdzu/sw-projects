// lib/core/ai/prompts/variant_prompt.dart
// 变式题生成 prompt
import '../ai_provider.dart';

class VariantPrompt {
  static const systemPrompt = '''
你是一个出题老师。根据用户提供的原题，生成 3 道变式题。

**要求**：
1. 核心知识点保持不变
2. 数字、情景、提问方式可变化
3. 难度与原题相当（不能更简单或更难）
4. 每道题独立完整，可直接做

**输出格式**（严格 JSON 数组）：
[
  {"index": 1, "text": "变式题 1 的完整内容"},
  {"index": 2, "text": "变式题 2 的完整内容"},
  {"index": 3, "text": "变式题 3 的完整内容"}
]
''';

  static List<Message> build({required String originalQuestion, String? knowledgePoint, String? difficulty}) {
    final userContent = StringBuffer('原题：$originalQuestion');
    if (knowledgePoint != null) {
      userContent.write('\n\n核心知识点：$knowledgePoint');
    }
    if (difficulty != null) {
      userContent.write('\n难度：$difficulty');
    }
    userContent.write('\n\n请生成 3 道变式题，输出 JSON 数组。');

    return [
      const Message(role: 'system', content: systemPrompt),
      Message(role: 'user', content: userContent.toString()),
    ];
  }

  /// 解析 AI 返回
  static List<String> parseVariants(String aiResponse) {
    final cleaned = aiResponse.trim();
    // 找 JSON 数组
    final arrayMatch = RegExp(r'\[[^\]]*\]', multiLine: false).firstMatch(cleaned);
    final arr = arrayMatch?.group(0) ?? cleaned;
    if (!arr.startsWith('[') || !arr.endsWith(']')) return [];

    final inner = arr.substring(1, arr.length - 1);
    final texts = <String>[];
    for (final item in _splitObjects(inner)) {
      // 提取 "text" 字段
      final textMatch = RegExp(r'"text"\s*:\s*"([^"]+)"').firstMatch(item);
      if (textMatch != null) {
        var t = textMatch.group(1)!;
        t = t.replaceAll('\\"', '"').replaceAll('\\n', '\n');
        texts.add(t);
      }
    }
    return texts;
  }

  static List<String> _splitObjects(String s) {
    final objects = <String>[];
    int depth = 0;
    final buf = StringBuffer();
    for (final c in s.split('')) {
      if (c == '{') depth++;
      if (c == '}') depth--;
      buf.write(c);
      if (depth == 0 && buf.isNotEmpty) {
        objects.add(buf.toString().trim());
        buf.clear();
      }
    }
    return objects.where((o) => o.isNotEmpty).toList();
  }
}
