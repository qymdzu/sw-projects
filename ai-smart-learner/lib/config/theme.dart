// lib/config/theme.dart
// 主题：浅色 + 暗色（v0.1 仅浅色），按学段切换主色
import 'package:flutter/material.dart';

class AppTheme {
  // 默认主色（小学/初中）
  static const Color primaryLight = Color(0xFF1976D2);
  static const Color primaryDark = Color(0xFF0D47A1);
  static const Color primaryBackground = Color(0xFFE3F2FD);

  // 高中主色
  static const Color highSchoolPrimary = Color(0xFF0D47A1);
  static const Color highSchoolBackground = Color(0xFFE8EAF6);

  // 状态色
  static const Color success = Color(0xFF34A853);
  static const Color warning = Color(0xFFFBBC04);
  static const Color error = Color(0xFFE53935);
  static const Color info = Color(0xFFF57C00);

  // 掌握度色
  static const Color masteryRed = Color(0xFFE53935);
  static const Color masteryYellow = Color(0xFFFBBC04);
  static const Color masteryGreen = Color(0xFF34A853);

  // 中性色
  static const Color textPrimary = Color(0xFF212121);
  static const Color textSecondary = Color(0xFF757575);
  static const Color textDisabled = Color(0xFF9E9E9E);
  static const Color divider = Color(0xFFE0E0E0);
  static const Color backgroundGrey = Color(0xFFF5F5F5);
  static const Color cardWhite = Colors.white;

  // 浅色主题
  static ThemeData light() {
    final colorScheme = ColorScheme.fromSeed(
      seedColor: primaryLight,
      brightness: Brightness.light,
      primary: primaryLight,
      secondary: primaryDark,
      error: error,
      surface: cardWhite,
    );

    return ThemeData(
      useMaterial3: true,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: backgroundGrey,
      appBarTheme: AppBarTheme(
        backgroundColor: primaryLight,
        foregroundColor: Colors.white,
        elevation: 0,
        centerTitle: true,
        titleTextStyle: const TextStyle(
          color: Colors.white,
          fontSize: 18,
          fontWeight: FontWeight.w600,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primaryLight,
          foregroundColor: Colors.white,
          minimumSize: const Size.fromHeight(56),  // 主按钮 56pt
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          textStyle: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: primaryLight,
          side: const BorderSide(color: primaryLight, width: 2),
          minimumSize: const Size.fromHeight(56),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          textStyle: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w500,
          ),
        ),
      ),
      cardTheme: CardThemeData(
        color: cardWhite,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: const BorderSide(color: divider, width: 1),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        contentPadding: const EdgeInsets.all(12),
      ),
      textTheme: const TextTheme(
        displayLarge: TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: textPrimary),
        headlineLarge: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: textPrimary),
        headlineMedium: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: textPrimary),
        titleLarge: TextStyle(fontSize: 18, fontWeight: FontWeight.w600, color: textPrimary),
        bodyLarge: TextStyle(fontSize: 16, color: textPrimary),
        bodyMedium: TextStyle(fontSize: 14, color: textPrimary),
        bodySmall: TextStyle(fontSize: 12, color: textSecondary),
        labelSmall: TextStyle(fontSize: 10, color: textDisabled),
      ),
    );
  }

  // 暗色主题（v0.3 推，v0.1 简单占位）
  static ThemeData dark() {
    return light();  // v0.1 复用浅色
  }

  /// 根据学段返回主色（学段适配 v0.1 简化）
  static Color primaryForGrade(String? grade) {
    if (grade == null) return primaryLight;
    if (grade.startsWith('senior_')) return highSchoolPrimary;
    return primaryLight;  // 小学/初中都用主蓝
  }
}
