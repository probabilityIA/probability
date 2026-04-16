import 'package:flutter/material.dart';

/// A [CircleAvatar] that loads a network image with graceful error handling.
///
/// When [imageUrl] is null, empty, or the image fails to load, the avatar
/// displays a fallback: either the first letter of [fallbackText] or a
/// generic [Icons.person] icon.
class NetworkAvatar extends StatelessWidget {
  final String? imageUrl;
  final String? fallbackText;
  final IconData fallbackIcon;
  final double radius;
  final Color? backgroundColor;
  final Color? foregroundColor;

  const NetworkAvatar({
    super.key,
    this.imageUrl,
    this.fallbackText,
    this.fallbackIcon = Icons.person,
    this.radius = 24,
    this.backgroundColor,
    this.foregroundColor,
  });

  bool get _hasValidUrl =>
      imageUrl != null && imageUrl!.trim().isNotEmpty;

  @override
  Widget build(BuildContext context) {
    final bgColor =
        backgroundColor ?? Theme.of(context).colorScheme.primary;
    final fgColor =
        foregroundColor ?? Theme.of(context).colorScheme.onPrimary;

    return ClipOval(
      child: Container(
        width: radius * 2,
        height: radius * 2,
        color: bgColor,
        child: _hasValidUrl
            ? Image.network(
                imageUrl!,
                width: radius * 2,
                height: radius * 2,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) {
                  return _buildFallback(fgColor);
                },
                loadingBuilder: (context, child, loadingProgress) {
                  if (loadingProgress == null) return child;
                  return Center(
                    child: SizedBox(
                      width: radius * 0.8,
                      height: radius * 0.8,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: fgColor,
                      ),
                    ),
                  );
                },
              )
            : _buildFallback(fgColor),
      ),
    );
  }

  Widget _buildFallback(Color color) {
    final text = fallbackText ?? '';
    if (text.isNotEmpty) {
      return Center(
        child: Text(
          text[0].toUpperCase(),
          style: TextStyle(
            color: color,
            fontWeight: FontWeight.bold,
            fontSize: radius * 0.8,
          ),
        ),
      );
    }
    return Center(
      child: Icon(fallbackIcon, color: color, size: radius),
    );
  }
}
