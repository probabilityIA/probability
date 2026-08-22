import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

class CoreButton extends StatefulWidget {
  const CoreButton({super.key, required this.onTap, this.size = 62});

  final VoidCallback onTap;
  final double size;

  @override
  State<CoreButton> createState() => _CoreButtonState();
}

class _CoreButtonState extends State<CoreButton>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 2600),
  )..repeat();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final brand = Theme.of(context).colorScheme.primary;
    final size = widget.size;

    return SizedBox(
      width: size + 26,
      height: size + 26,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, child) {
          return CustomPaint(
            painter: _ElectricPainter(
              progress: _controller.value,
              color: brand,
            ),
            child: child,
          );
        },
        child: Center(
          child: Material(
            color: Colors.transparent,
            shape: const CircleBorder(),
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              onTap: widget.onTap,
              customBorder: const CircleBorder(),
              child: Container(
                width: size,
                height: size,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      Color.lerp(brand, Colors.white, 0.18)!,
                      brand,
                      Color.lerp(brand, Colors.black, 0.22)!,
                    ],
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: brand.withValues(alpha: 0.45),
                      blurRadius: 18,
                      spreadRadius: 1,
                      offset: const Offset(0, 5),
                    ),
                  ],
                  border: Border.all(color: AppColors.surface, width: 3),
                ),
                alignment: Alignment.center,
                child: Image.asset(
                  'assets/images/logo_mark_white.png',
                  width: size * 0.44,
                  height: size * 0.44,
                  errorBuilder: (context, error, stack) => Icon(
                    Icons.bolt_rounded,
                    size: size * 0.5,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _ElectricPainter extends CustomPainter {
  _ElectricPainter({required this.progress, required this.color});

  final double progress;
  final Color color;

  static const int _arcs = 5;
  static const int _segments = 22;

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final baseRadius = size.width / 2 - 6;

    final halo = Paint()
      ..shader = RadialGradient(
        colors: [
          color.withValues(alpha: 0.22 + 0.10 * math.sin(progress * math.pi * 2)),
          color.withValues(alpha: 0),
        ],
      ).createShader(Rect.fromCircle(center: center, radius: baseRadius));
    canvas.drawCircle(center, baseRadius, halo);

    for (var arc = 0; arc < _arcs; arc++) {
      final seed = arc * 97.13;
      final phase = (progress * (1 + arc * 0.23) + arc / _arcs) % 1.0;
      final alpha = math.sin(phase * math.pi);
      if (alpha <= 0.05) continue;

      final paint = Paint()
        ..color = color.withValues(alpha: alpha * 0.9)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.6
        ..strokeCap = StrokeCap.round
        ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 1.2);

      final start = (seed % 6.28) + progress * 2.4;
      final sweep = 1.1 + (arc % 3) * 0.5;
      final path = Path();

      for (var i = 0; i <= _segments; i++) {
        final t = i / _segments;
        final angle = start + sweep * t;
        final jitter = math.sin((t * 14) + seed + progress * 12) * 3.2 +
            math.cos((t * 23) + seed * 1.7) * 1.8;
        final radius = baseRadius + jitter;
        final point = Offset(
          center.dx + math.cos(angle) * radius,
          center.dy + math.sin(angle) * radius,
        );
        if (i == 0) {
          path.moveTo(point.dx, point.dy);
        } else {
          path.lineTo(point.dx, point.dy);
        }
      }
      canvas.drawPath(path, paint);
    }

    final sparkPaint = Paint()..color = color.withValues(alpha: 0.85);
    for (var i = 0; i < 3; i++) {
      final phase = (progress * (1.6 + i * 0.4) + i / 3) % 1.0;
      final angle = phase * math.pi * 2 + i * 2.1;
      final radius = baseRadius + math.sin(phase * math.pi * 4) * 3;
      canvas.drawCircle(
        Offset(
          center.dx + math.cos(angle) * radius,
          center.dy + math.sin(angle) * radius,
        ),
        1.6 * math.sin(phase * math.pi).clamp(0.0, 1.0),
        sparkPaint,
      );
    }
  }

  @override
  bool shouldRepaint(covariant _ElectricPainter oldDelegate) =>
      oldDelegate.progress != progress || oldDelegate.color != color;
}
