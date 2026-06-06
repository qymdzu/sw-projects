// lib/core/ocr/baidu_edu_provider.dart
// 百度云 OCR 教育版实现
import 'dart:convert';
import 'dart:io';
import 'package:dio/dio.dart';
import 'package:logger/logger.dart';

import '../errors/app_error.dart';
import '../errors/error_codes.dart';
import 'ocr_provider.dart';

final _logger = Logger();

class BaiduEduProvider implements OcrProvider {
  final Dio _dio;
  final String _apiKey;
  final String _secretKey;
  String? _accessToken;
  DateTime? _tokenExpiresAt;

  BaiduEduProvider({
    required String apiKey,
    required String secretKey,
    Dio? dio,
  })  : _apiKey = apiKey,
        _secretKey = secretKey,
        _dio = dio ?? Dio() {
    _dio.options.connectTimeout = const Duration(seconds: 10);
    _dio.options.receiveTimeout = const Duration(seconds: 30);
  }

  @override
  String get name => 'baidu_edu';

  @override
  bool get supportsHandwritingRemoval => true;

  @override
  bool supportsQuestionType(String questionType) => true;  // 教育版支持所有

  @override
  Future<OcrResult> recognize({
    required File imageFile,
    bool removeHandwriting = true,
    Map<String, dynamic>? options,
  }) async {
    if (_apiKey.isEmpty || _secretKey.isEmpty) {
      throw OcrError(
        code: ErrorCodes.ocrAuth,
        message: '百度云 OCR Key 未配置',
      );
    }

    try {
      // 1. 获取 access_token（缓存 29 天）
      final token = await _getAccessToken();

      // 2. 调 OCR API（教育版去笔迹）
      final imageBytes = await imageFile.readAsBytes();
      final base64Image = base64Encode(imageBytes);

      final response = await _dio.post<Map<String, dynamic>>(
        'https://aip.baidubce.com/rest/2.0/ocr/v1/handwriting',
        queryParameters: {
          'access_token': token,
        },
        data: {
          'image': base64Image,
          'recognize_granularity': 'small',
          'probability': 'true',
        },
        options: Options(
          headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        ),
      );

      final data = response.data;
      if (data == null || data['words_result'] == null) {
        throw OcrError(
          code: ErrorCodes.ocrInternal,
          message: 'OCR 返回结果为空',
          detail: data?.toString(),
        );
      }

      // 3. 解析结果
      final words = (data['words_result'] as List)
          .map((w) => (w as Map)['words'] as String)
          .where((s) => s.isNotEmpty)
          .toList();
      final text = words.join('\n');

      return OcrResult(
        text: text,
        latex: null,  // 百度云 OCR 不直接返回 LaTeX
        tags: const [],
        confidence: 0.9,
        raw: data,
        createdAt: DateTime.now(),
        providerName: name,
      );
    } on DioException catch (e) {
      _logger.e('百度云 OCR 调用失败: ${e.message}');
      if (e.type == DioExceptionType.connectionTimeout || e.type == DioExceptionType.receiveTimeout) {
        throw NetworkError(
          code: ErrorCodes.netTimeout,
          message: 'OCR 服务超时',
          cause: e,
        );
      }
      throw OcrError(
        code: ErrorCodes.ocrInternal,
        message: 'OCR 服务异常',
        cause: e,
      );
    } catch (e) {
      _logger.e('百度云 OCR 异常: $e');
      throw OcrError(
        code: ErrorCodes.ocrInternal,
        message: 'OCR 服务异常',
        cause: e,
      );
    }
  }

  Future<String> _getAccessToken() async {
    if (_accessToken != null && _tokenExpiresAt != null && DateTime.now().isBefore(_tokenExpiresAt!)) {
      return _accessToken!;
    }

    final response = await _dio.post<Map<String, dynamic>>(
      'https://aip.baidubce.com/oauth/2.0/token',
      queryParameters: {
        'grant_type': 'client_credentials',
        'client_id': _apiKey,
        'client_secret': _secretKey,
      },
    );

    final data = response.data;
    if (data == null || data['access_token'] == null) {
      throw OcrError(
        code: ErrorCodes.ocrAuth,
        message: '获取 access_token 失败',
        detail: data?.toString(),
      );
    }

    _accessToken = data['access_token'] as String;
    // 提前 1 天刷新
    _tokenExpiresAt = DateTime.now().add(const Duration(days: 29));
    return _accessToken!;
  }

  @override
  Future<void> close() async {
    _dio.close();
  }
}
