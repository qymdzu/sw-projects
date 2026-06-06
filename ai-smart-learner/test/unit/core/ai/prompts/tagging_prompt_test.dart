// test/unit/core/ai/prompts/tagging_prompt_test.dart
// 智能打标 prompt 解析单测
import 'package:flutter_test/flutter_test.dart';
import 'package:ai_smart_learner/core/ai/prompts/tagging_prompt.dart';

void main() {
  group('TaggingPrompt.parseJson', () {
    test('解析标准 JSON', () {
      const input = '''
{
  "knowledgePoint": "异分母分数加法",
  "errorType": "计算错",
  "subject": "数学",
  "chapter": "分数运算",
  "difficulty": "medium"
}
''';
      final result = TaggingPrompt.parseJson(input);
      expect(result, isNotNull);
      expect(result!['knowledgePoint'], '异分母分数加法');
      expect(result['errorType'], '计算错');
      expect(result['subject'], '数学');
    });

    test('从 markdown 中提取 JSON', () {
      const input = '''
好的，分析结果如下：
```json
{"knowledgePoint": "长方形周长", "errorType": "概念错"}
```
''';
      final result = TaggingPrompt.parseJson(input);
      expect(result, isNotNull);
      expect(result!['knowledgePoint'], '长方形周长');
    });

    test('无效 JSON 返回 null', () {
      const input = '这不是 JSON';
      final result = TaggingPrompt.parseJson(input);
      expect(result, isNull);
    });
  });
}
