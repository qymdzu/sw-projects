// lib/ui/shared/error_boundary.dart
// 错误边界
import 'package:flutter/material.dart';

import '../../config/theme.dart';

class ErrorBoundary extends StatelessWidget {
  final Widget child;
  const ErrorBoundary({required this.child, super.key});

  @override
  Widget build(BuildContext context) {
    return child;
  }
}

/// 全屏错误页
class ErrorState extends StatelessWidget {
  final String title;
  final String? detail;
  final VoidCallback? onRetry;
  final VoidCallback? onBack;

  const ErrorState({
    required this.title,
    this.detail,
    this.onRetry,
    this.onBack,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('⚠️', style: TextStyle(fontSize: 80)),
            const SizedBox(height: 16),
            Text(
              title,
              textAlign: TextAlign.center,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600, color: AppTheme.textPrimary),
            ),
            if (detail != null) ...[
              const SizedBox(height: 8),
              Text(
                detail!,
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 14, color: AppTheme.textSecondary),
              ),
            ],
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                if (onBack != null)
                  OutlinedButton(onPressed: onBack, child: const Text('返回')),
                if (onBack != null && onRetry != null) const SizedBox(width: 12),
                if (onRetry != null)
                  ElevatedButton(onPressed: onRetry, child: const Text('重试')),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

/// 全屏空状态
class EmptyState extends StatelessWidget {
  final String emoji;
  final String title;
  final String? hint;

  const EmptyState({
    required this.emoji,
    required this.title,
    this.hint,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(emoji, style: const TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text(
            title,
            style: const TextStyle(fontSize: 16, color: AppTheme.textSecondary),
          ),
          if (hint != null) ...[
            const SizedBox(height: 8),
            Text(
              hint!,
              textAlign: TextAlign.center,
              style: const TextStyle(fontSize: 13, color: AppTheme.textDisabled),
            ),
          ],
        ],
      ),
    );
  }
}
