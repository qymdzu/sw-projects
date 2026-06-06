// lib/ui/students/capture_page.dart
// 拍照页
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import 'package:logger/logger.dart';

import '../../core/auth/user_session.dart';
import '../../core/data/repositories/error_question_repository.dart';
import '../../core/providers.dart';
import '../../core/students/capture.dart';
import '../../config/theme.dart';
import '../shared/error_boundary.dart';

final _captureServiceProvider = Provider<CaptureService>((ref) {
  final aiFactory = ref.watch(_aiFactoryProvider);
  final aiProvider = aiFactory.create('deepseek');
  final ocrFactory = ref.watch(_ocrFactoryProvider);
  // 简单用第一个 provider；fallback 由调用方处理
  // 实际上需要 router
  throw UnimplementedError('Need proper DI - v0.1 simplified');
});

final _aiFactoryProvider = Provider((ref) => _AiFactoryStub());
final _ocrFactoryProvider = Provider((ref) => _OcrFactoryStub());

class _AiFactoryStub {
  dynamic create(String name) => null;
}

class _OcrFactoryStub {
  dynamic create(String name) => null;
}

class CapturePage extends ConsumerStatefulWidget {
  final String type;  // 'error' / 'missing'
  const CapturePage({required this.type, super.key});

  @override
  ConsumerState<CapturePage> createState() => _CapturePageState();
}

class _CapturePageState extends ConsumerState<CapturePage> {
  final ImagePicker _picker = ImagePicker();
  File? _image;
  bool _processing = false;
  String? _result;
  String? _error;
  final _logger = Logger();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.type == 'error' ? '拍错题' : '拍漏题')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              // 类型 Tab
              Row(
                children: [
                  Expanded(
                    child: _TypeTab(
                      label: '错题',
                      selected: widget.type == 'error',
                      onTap: () => context.replace('/student/capture?type=error'),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: _TypeTab(
                      label: '漏题',
                      selected: widget.type == 'missing',
                      onTap: () => context.replace('/student/capture?type=missing'),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),

              // 拍照预览
              Expanded(
                child: _image != null
                    ? ClipRRect(
                        borderRadius: BorderRadius.circular(12),
                        child: Image.file(_image!, fit: BoxFit.contain),
                      )
                    : Container(
                        decoration: BoxDecoration(
                          color: AppTheme.backgroundGrey,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text('📷', style: TextStyle(fontSize: 64)),
                              SizedBox(height: 12),
                              Text('点击下方按钮拍照', style: TextStyle(color: AppTheme.textSecondary)),
                            ],
                          ),
                        ),
                      ),
              ),
              const SizedBox(height: 16),

              // 错误/结果
              if (_error != null)
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: AppTheme.error.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(_error!, style: const TextStyle(color: AppTheme.error)),
                ),
              if (_result != null)
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: AppTheme.success.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(_result!, style: const TextStyle(color: AppTheme.success)),
                ),
              const SizedBox(height: 16),

              // 操作按钮
              if (_image == null) ...[
                ElevatedButton.icon(
                  onPressed: _processing ? null : () => _pickImage(ImageSource.camera),
                  icon: const Text('📷', style: TextStyle(fontSize: 24)),
                  label: const Text('拍照'),
                ),
                const SizedBox(height: 8),
                OutlinedButton.icon(
                  onPressed: _processing ? null : () => _pickImage(ImageSource.gallery),
                  icon: const Icon(Icons.photo_library),
                  label: const Text('从相册选'),
                ),
              ] else ...[
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: _processing ? null : () => setState(() => _image = null),
                        child: const Text('重新拍'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: _processing ? null : _submit,
                        child: _processing
                            ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                            : const Text('提交识别'),
                      ),
                    ),
                  ],
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _pickImage(ImageSource source) async {
    try {
      final picked = await _picker.pickImage(source: source, imageQuality: 85);
      if (picked != null) {
        setState(() => _image = File(picked.path));
      }
    } catch (e) {
      setState(() => _error = '拍照失败: $e');
    }
  }

  Future<void> _submit() async {
    if (_image == null) return;
    setState(() {
      _processing = true;
      _error = null;
      _result = null;
    });
    try {
      // v0.1 简化：占位流程（实际 capture service 需要在 ProviderScope 正确注入）
      await Future.delayed(const Duration(seconds: 2));
      setState(() {
        _result = '✅ 识别成功！已加入本周池（v0.1 占位，请按 Stage 6 测试）';
        _processing = false;
      });
      Future.delayed(const Duration(seconds: 2), () {
        if (mounted) context.go('/student');
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _processing = false;
      });
    }
  }
}

class _TypeTab extends StatelessWidget {
  final String label;
  final bool selected;
  final VoidCallback onTap;

  const _TypeTab({required this.label, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: selected ? AppTheme.primaryLight : AppTheme.cardWhite,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: selected ? AppTheme.primaryLight : AppTheme.divider,
          ),
        ),
        child: Center(
          child: Text(
            label,
            style: TextStyle(
              color: selected ? Colors.white : AppTheme.textPrimary,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
    );
  }
}
