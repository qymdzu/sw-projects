// lib/core/ai/prompts/tagging_prompt.dart
// 智能打标 prompt
import '../ai_provider.dart';

class TaggingPrompt {
  /// 系统 prompt
  static const systemPrompt = '''
你是一个小学/初中/高中全科老师，擅长分析错题。
请根据用户提交的题目内容，给出以下信息：
1. knowledgePoint：知识点（不超过 15 字）
2. errorType：错误类型（"计算错"/"概念错"/"审题错"/"跳题"/"粗心"）
3. subject：学科（"数学"/"语文"/"英语"/"物理"/"化学"/"生物"/"政治"/"历史"/"地理"）
4. chapter：建议章节（不超过 10 字）
5. difficulty：难度（"easy"/"medium"/"hard"）

**严格用 JSON 格式输出**，不要任何额外说明。示例：
{
  "knowledgePoint": "异分母分数加法",
  "errorType": "计算错",
  "subject": "数学",
  "chapter": "分数运算",
  "difficulty": "medium"
}
''';

  /// 构建打标 messages
  static List<Message> build({required String questionText, String? subjectHint, String? stage}) {
    final userContent = StringBuffer('题目：$questionText');
    if (subjectHint != null) {
      userContent.write('\n\n提示学科：$subjectHint');
    }
    if (stage != null) {
      userContent.write('\n提示学段：$stage');
    }
    userContent.write('\n\n请输出 JSON。');

    return [
      const Message(role: 'system', content: systemPrompt),
      Message(role: 'user', content: userContent.toString()),
    ];
  }

  /// 解析 AI 返回
  static Map<String, dynamic>? parseJson(String aiResponse) {
    final cleaned = aiResponse.trim();
    // 尝试提取 JSON（AI 可能包在 ```json 中）
    final jsonMatch = RegExp(r'\{[^{}]*\}', multiLine: false).firstMatch(cleaned);
    final jsonStr = jsonMatch?.group(0) ?? cleaned;
    try {
      return _parseMap(jsonStr);
    } catch (_) {
      return null;
    }
  }

  static Map<String, dynamic>? _parseMap(String s) {
    // 简化版 JSON 解析（生产用 dart:convert）
    s = s.trim();
    if (!s.startsWith('{') || !s.endsWith('}')) return null;
    final inner = s.substring(1, s.length - 1);
    final result = <String, dynamic>{};
    for (final pair in inner.split(',')) {
      final kv = pair.split(':');
      if (kv.length != 2) continue;
      final key = kv[0].trim().replaceAll('"', '').replaceAll("'", '');
      final value = kv[1].trim().replaceAll('"', '').replaceAll("'", '');
      result[key] = value;
    }
    return result.isEmpty ? null : result;
  }
}
