// test/unit/core/ai/prompts/variant_prompt_test.dart
// 变式题 prompt 解析单测
import 'package:flutter_test/flutter_test.dart';
import 'package:ai_smart_learner/core/ai/prompts/variant_prompt.dart';

void main() {
  group('VariantPrompt.parseVariants', () {
    test('解析 3 道变式题', () {
      const input = '''
[
  {"index": 1, "text": "变式题 1"},
  {"index": 2, "text": "变式题 2"},
  {"index": 3, "text": "变式题 3"}
]
''';
      final result = VariantPrompt.parseVariants(input);
      expect(result.length, 3);
      expect(result[0], '变式题 1');
    });

    test('从 markdown 提取数组', () {
      const input = '''
变式题如下：
```json
[{"index": 1, "text": "题 1"}, {"index": 2, "text": "题 2"}]
```
''';
      final result = VariantPrompt.parseVariants(input);
      expect(result.length, 2);
    });

    test('空数组返回空列表', () {
      const input = '[]';
      final result = VariantPrompt.parseVariants(input);
      expect(result.isEmpty, true);
    });
  });
}
