// lib/ui/shared/settings_page.dart
// 设置页
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/user_session.dart';
import '../../config/theme.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(currentUserProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('设置')),
      body: ListView(
        children: [
          // 用户
          ListTile(
            leading: Text(user?.role == 'student' ? '🎒' : '👨', style: const TextStyle(fontSize: 32)),
            title: Text(user?.name ?? '未选择'),
            subtitle: Text(user?.role == 'student' ? '学生' : '大人'),
            trailing: const Icon(Icons.swap_horiz),
            onTap: () async {
              await ref.read(currentUserProvider.notifier).clear();
              if (context.mounted) context.go('/');
            },
          ),
          const Divider(),

          // OCR 服务
          ListTile(
            leading: const Icon(Icons.document_scanner, color: AppTheme.primaryLight),
            title: const Text('OCR 服务'),
            subtitle: const Text('默认：百度云教育版'),
            onTap: () {
              // v0.1 简化为提示
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('OCR 服务设置 v0.2 推')),
              );
            },
          ),

          // AI 服务
          ListTile(
            leading: const Icon(Icons.psychology, color: AppTheme.primaryLight),
            title: const Text('AI 服务'),
            subtitle: const Text('默认：DeepSeek'),
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('AI 服务设置 v0.2 推')),
              );
            },
          ),

          // 图片保留
          ListTile(
            leading: const Icon(Icons.image, color: AppTheme.primaryLight),
            title: const Text('图片保留'),
            subtitle: const Text('默认：识别后立即删除 ⚠️'),
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('图片保留设置 v0.2 推')),
              );
            },
          ),

          const Divider(),

          // 关于
          ListTile(
            leading: const Icon(Icons.info_outline),
            title: const Text('关于'),
            subtitle: const Text('v0.1.0 (Stage 5 编码中)'),
          ),
        ],
      ),
    );
  }
}
