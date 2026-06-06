// lib/ui/shared/user_select_page.dart
// 启动选身份页
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/models/user.dart';
import '../../core/data/repositories/user_repository.dart';
import '../../core/providers.dart';
import '../../config/theme.dart';

class UserSelectPage extends ConsumerWidget {
  const UserSelectPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userRepo = ref.watch(userRepositoryProvider);

    return Scaffold(
      backgroundColor: AppTheme.backgroundGrey,
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('📚', style: TextStyle(fontSize: 80)),
              const SizedBox(height: 16),
              const Text(
                'AI智能学习机',
                style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: AppTheme.textPrimary),
              ),
              const SizedBox(height: 8),
              const Text(
                '你好，请选择身份',
                style: TextStyle(fontSize: 16, color: AppTheme.textSecondary),
              ),
              const SizedBox(height: 48),

              // 学生按钮
              _RoleButton(
                emoji: '🎒',
                label: '我是学生',
                onTap: () => _selectOrCreate(context, ref, userRepo, 'student'),
              ),
              const SizedBox(height: 16),

              // 大人按钮
              _RoleButton(
                emoji: '👨',
                label: '我是大人',
                onTap: () => _selectOrCreate(context, ref, userRepo, 'adult'),
              ),

              const Spacer(),
              TextButton(
                onPressed: () => context.push('/settings'),
                child: const Text('⚙ 设置'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _selectOrCreate(
    BuildContext context,
    WidgetRef ref,
    UserRepository userRepo,
    String role,
  ) async {
    // 查该角色下的用户
    final users = await userRepo.findByRole(role);

    if (!context.mounted) return;

    if (users.isEmpty) {
      // 没用户：弹窗输入名字
      _showCreateDialog(context, ref, userRepo, role);
    } else if (users.length == 1) {
      // 单用户：直接进
      await _setCurrentUser(ref, users.first, role);
    } else {
      // 多用户：选一个
      _showUserPicker(context, ref, userRepo, users);
    }
  }

  Future<void> _setCurrentUser(WidgetRef ref, User user, String role) async {
    await ref.read(currentUserProvider.notifier).setUser(user);
    if (role == 'student') {
      // ignore: use_build_context_synchronously
      // (context used in caller)
    }
  }

  void _showCreateDialog(BuildContext context, WidgetRef ref, UserRepository userRepo, String role) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(role == 'student' ? '你的名字是？' : '你的称呼？'),
        content: TextField(
          controller: controller,
          autofocus: true,
          decoration: const InputDecoration(hintText: '如：儿子 / 公子'),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('取消')),
          ElevatedButton(
            onPressed: () async {
              final name = controller.text.trim();
              if (name.isEmpty) return;
              final user = await userRepo.create(
                name: name,
                role: role,
                grade: role == 'student' ? 'primary_3' : null,
              );
              if (!ctx.mounted) return;
              Navigator.pop(ctx);
              await _setCurrentUser(ref, user, role);
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _showUserPicker(BuildContext context, WidgetRef ref, UserRepository userRepo, List<User> users) {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Padding(
              padding: EdgeInsets.all(16),
              child: Text('选择用户', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            ),
            ...users.map((u) => ListTile(
              leading: Text(u.role == 'student' ? '🎒' : '👨', style: const TextStyle(fontSize: 32)),
              title: Text(u.name),
              subtitle: Text(u.role == 'student' ? '学生' : '大人'),
              onTap: () async {
                Navigator.pop(ctx);
                await _setCurrentUser(ref, u, u.role);
              },
            )),
          ],
        ),
      ),
    );
  }
}

class _RoleButton extends StatelessWidget {
  final String emoji;
  final String label;
  final VoidCallback onTap;

  const _RoleButton({required this.emoji, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Container(
        height: 100,
        decoration: BoxDecoration(
          color: AppTheme.cardWhite,
          borderRadius: BorderRadius.circular(16),
          boxShadow: const [
            BoxShadow(color: Color(0x14000000), blurRadius: 8, offset: Offset(0, 2)),
          ],
        ),
        child: Row(
          children: [
            const SizedBox(width: 24),
            Text(emoji, style: const TextStyle(fontSize: 40)),
            const SizedBox(width: 16),
            Expanded(
              child: Text(
                label,
                style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w600, color: AppTheme.textPrimary),
              ),
            ),
            const Icon(Icons.arrow_forward_ios, color: AppTheme.textSecondary),
            const SizedBox(width: 16),
          ],
        ),
      ),
    );
  }
}
