import 'package:flutter/material.dart';
import '../../theme/app_colors.dart';
import '../../utils/image_memory.dart';
import '../../utils/brand_assets.dart';
import '../../utils/formatters.dart';

class BrandLogo extends StatelessWidget {
  const BrandLogo({
    super.key,
    required this.name,
    this.imageUrl,
    this.size = 40,
    this.radius = 10,
    this.padding = 6,
    this.background = AppColors.surface,
    this.bordered = true,
  });

  const BrandLogo.carrier({
    super.key,
    required this.name,
    this.size = 40,
    this.radius = 10,
    this.padding = 6,
    this.background = AppColors.surface,
    this.bordered = true,
  }) : imageUrl = null;

  final String name;
  final String? imageUrl;
  final double size;
  final double radius;
  final double padding;
  final Color background;
  final bool bordered;

  String? get _resolvedUrl {
    final direct = BrandAssets.mediaUrl(imageUrl);
    if (direct != null) return direct;
    return BrandAssets.integrationLogo(name) ?? BrandAssets.carrierLogo(name);
  }

  @override
  Widget build(BuildContext context) {
    final url = _resolvedUrl;
    return Container(
      width: size,
      height: size,
      padding: EdgeInsets.all(padding),
      decoration: BoxDecoration(
        color: background,
        borderRadius: BorderRadius.circular(radius),
        border: bordered ? Border.all(color: AppColors.border) : null,
      ),
      child: url == null
          ? _Fallback(name: name, size: size)
          : Image.network(
              url,
              fit: BoxFit.contain,
              cacheWidth: ImageMemory.decodePixels(context, size),
              cacheHeight: ImageMemory.decodePixels(context, size),
              filterQuality: FilterQuality.medium,
              webHtmlElementStrategy: WebHtmlElementStrategy.fallback,
              errorBuilder: (context, error, stack) => _Fallback(name: name, size: size),
              loadingBuilder: (context, child, progress) {
                if (progress == null) return child;
                return _Fallback(name: name, size: size, muted: true);
              },
            ),
    );
  }
}

class _Fallback extends StatelessWidget {
  const _Fallback({required this.name, required this.size, this.muted = false});

  final String name;
  final double size;
  final bool muted;

  static const List<Color> _palette = [
    AppColors.primary,
    Color(0xFF0EA5E9),
    Color(0xFFF97316),
    Color(0xFF10B981),
    Color(0xFFEC4899),
    Color(0xFF6366F1),
  ];

  @override
  Widget build(BuildContext context) {
    final color = muted
        ? AppColors.borderStrong
        : _palette[name.isEmpty ? 0 : name.codeUnitAt(0) % _palette.length];
    return Container(
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(size * 0.16),
      ),
      child: Text(
        AppFormat.initials(name),
        style: TextStyle(
          fontFamily: 'Inter',
          fontSize: size * 0.34,
          fontWeight: FontWeight.w700,
          color: color,
          height: 1,
        ),
      ),
    );
  }
}
