// lib/config/router.dart
// 路由：go_router
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/auth/user_session.dart';
import '../ui/shared/user_select_page.dart';
import '../ui/shared/settings_page.dart';
import '../ui/students/student_home_page.dart';
import '../ui/students/capture_page.dart';
import '../ui/students/weekly_pool_page.dart';
import '../ui/students/weekly_test_page.dart';
import '../ui/students/weekly_test_result_page.dart';
import '../ui/students/rainbow_chart_page.dart';
import '../ui/students/archive_page.dart';

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final user = ref.read(currentUserProvider);
      final loc = state.matchedLocation;
      // 启动选身份：未选用户
      if (user == null && loc != '/') return '/';
      return null;
    },
    routes: [
      GoRoute(
        path: '/',
        name: 'user_select',
        builder: (context, state) => const UserSelectPage(),
      ),
      GoRoute(
        path: '/settings',
        name: 'settings',
        builder: (context, state) => const SettingsPage(),
      ),
      // 学生
      GoRoute(
        path: '/student',
        name: 'student_home',
        builder: (context, state) => const StudentHomePage(),
        routes: [
          GoRoute(
            path: 'capture',
            name: 'capture',
            builder: (context, state) {
              final type = state.uri.queryParameters['type'] ?? 'error';
              return CapturePage(type: type);
            },
          ),
          GoRoute(
            path: 'pool',
            name: 'weekly_pool',
            builder: (context, state) => const WeeklyPoolPage(),
          ),
          GoRoute(
            path: 'test',
            name: 'weekly_test',
            builder: (context, state) {
              final type = state.uri.queryParameters['type'] ?? 'error';
              return WeeklyTestPage(type: type);
            },
          ),
          GoRoute(
            path: 'test/result',
            name: 'weekly_test_result',
            builder: (context, state) {
              final correctCount = int.tryParse(state.uri.queryParameters['correct'] ?? '0') ?? 0;
              final totalCount = int.tryParse(state.uri.queryParameters['total'] ?? '0') ?? 0;
              return WeeklyTestResultPage(
                correctCount: correctCount,
                totalCount: totalCount,
              );
            },
          ),
          GoRoute(
            path: 'chart',
            name: 'rainbow_chart',
            builder: (context, state) => const RainbowChartPage(),
          ),
          GoRoute(
            path: 'archive',
            name: 'archive',
            builder: (context, state) => const ArchivePage(),
          ),
        ],
      ),
    ],
  );
});
