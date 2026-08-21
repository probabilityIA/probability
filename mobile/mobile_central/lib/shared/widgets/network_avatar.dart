import 'package:flutter/material.dart';
import '../utils/brand_assets.dart';
import '../utils/formatters.dart';

class NetworkAvatar extends StatelessWidget {
  const NetworkAvatar({
    super.key,
    this.imageUrl,
    this.fallbackText,
    this.fallbackIcon = Icons.person,
    this.radius = 24,
    this.backgroundColor,
    this.foregroundColor,
  });

  final String? imageUrl;
  final String? fallbackText;
  final IconData fallbackIcon;
  final double radius;
  final Color? backgroundColor;
  final Color? foregroundColor;

  String? get _resolvedUrl => BrandAssets.mediaUrl(imageUrl);

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final background = backgroundColor ?? scheme.primaryContainer;
    final foreground = foregroundColor ?? scheme.primary;
    final url = _resolvedUrl;

    return ClipOval(
      child: Container(
        width: radius * 2,
        height: radius * 2,
        color: background,
        child: url == null
            ? _Fallback(
                text: fallbackText,
                icon: fallbackIcon,
                color: foreground,
                radius: radius,
              )
            : Image.network(
                url,
                width: radius * 2,
                height: radius * 2,
                fit: BoxFit.cover,
                filterQuality: FilterQuality.medium,
                webHtmlElementStrategy: WebHtmlElementStrategy.fallback,
                errorBuilder: (context, error, stackTrace) => _Fallback(
                  text: fallbackText,
                  icon: fallbackIcon,
                  color: foreground,
                  radius: radius,
                ),
                loadingBuilder: (context, child, progress) {
                  if (progress == null) return child;
                  return _Fallback(
                    text: fallbackText,
                    icon: fallbackIcon,
                    color: foreground,
                    radius: radius,
                  );
                },
              ),
      ),
    );
  }
}

class _Fallback extends StatelessWidget {
  const _Fallback({
    required this.text,
    required this.icon,
    required this.color,
    required this.radius,
  });

  final String? text;
  final IconData icon;
  final Color color;
  final double radius;

  @override
  Widget build(BuildContext context) {
    final clean = (text ?? '').trim();
    if (clean.isEmpty) {
      return Center(child: Icon(icon, color: color, size: radius));
    }
    return Center(
      child: Text(
        AppFormat.initials(clean),
        style: TextStyle(
          fontFamily: 'Inter',
          color: color,
          fontWeight: FontWeight.w700,
          fontSize: radius * 0.66,
        ),
      ),
    );
  }
}
