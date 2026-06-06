// lib/app.dart
// App 根 widget：MaterialApp + go_router + 主题
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'config/router.dart';
import 'config/theme.dart';
import 'ui/shared/error_boundary.dart';

class AiSmartLearnerApp extends ConsumerWidget {
  const AiSmartLearnerApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'AI智能学习机',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light(),
      darkTheme: AppTheme.dark(),
      routerConfig: router,
      builder: (context, child) {
        return ErrorBoundary(child: child ?? const SizedBox.shrink());
      },
    );
  }
}
