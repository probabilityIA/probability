import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key, this.onFinished, this.hold = false});

  final VoidCallback? onFinished;
  final bool hold;

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen>
    with TickerProviderStateMixin {
  late final AnimationController _loop = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 3200),
  )..repeat();

  late final AnimationController _intro = AnimationController(
    vsync: this,
    duration: const Duration(milliseconds: 2000),
  );

  late final Animation<double> _rise = CurvedAnimation(
    parent: _intro,
    curve: const Interval(0, 0.45, curve: Curves.easeOutBack),
  );

  late final Animation<double> _ignite = CurvedAnimation(
    parent: _intro,
    curve: const Interval(0.15, 0.75, curve: Curves.easeOutCubic),
  );

  late final Animation<double> _title = CurvedAnimation(
    parent: _intro,
    curve: const Interval(0.55, 1, curve: Curves.easeOut),
  );

  @override
  void initState() {
    super.initState();
    _intro.forward().whenComplete(() {
      if (!widget.hold) widget.onFinished?.call();
    });
  }

  @override
  void didUpdateWidget(SplashScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.hold && !widget.hold && _intro.isCompleted) {
      widget.onFinished?.call();
    }
  }

  @override
  void dispose() {
    _loop.dispose();
    _intro.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF120326),
      body: Center(
        child: AnimatedBuilder(
          animation: Listenable.merge([_loop, _intro]),
          builder: (context, child) {
            final ignite = _ignite.value;
            final rise = _rise.value;

            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                SizedBox(
                  width: 300,
                  height: 300,
                  child: CustomPaint(
                    painter: _StarPainter(
                      progress: _loop.value,
                      ignite: ignite,
                    ),
                    child: Center(
                      child: Transform.scale(
                        scale: 0.6 + rise * 0.4,
                        child: Opacity(
                          opacity: rise.clamp(0, 1),
                          child: _Core(size: 108, glow: ignite),
                        ),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 26),
                Opacity(
                  opacity: _title.value.clamp(0, 1),
                  child: const Text(
                    'PROBABILITY',
                    style: TextStyle(
                      fontFamily: 'Inter',
                      fontSize: 15,
                      fontWeight: FontWeight.w700,
                      letterSpacing: 7,
                      color: Colors.white,
                    ),
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}

class _Core extends StatelessWidget {
  const _Core({required this.size, required this.glow});

  final double size;
  final double glow;

  @override
  Widget build(BuildContext context) {
    const brand = AppColors.primary;

    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFF8B5CF6), brand, Color(0xFF3B0F91)],
        ),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFFFFB020).withValues(alpha: 0.55 * glow),
            blurRadius: 46 * glow,
            spreadRadius: 6 * glow,
          ),
          BoxShadow(
            color: brand.withValues(alpha: 0.65 * glow),
            blurRadius: 26 * glow,
            spreadRadius: 2 * glow,
          ),
        ],
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
    );
  }
}

class _StarPainter extends CustomPainter {
  _StarPainter({required this.progress, required this.ignite});

  final double progress;
  final double ignite;

  static const int _tongues = 30;
  static const int _bolts = 9;
  static const int _boltSegments = 12;

  static const Color _amber = Color(0xFFFFB020);
  static const Color _orange = Color(0xFFFF6A1A);
  static const Color _violet = Color(0xFF8B5CF6);
  static const Color _cyan = Color(0xFFB794FF);

  @override
  void paint(Canvas canvas, Size size) {
    if (ignite <= 0.01) return;

    final center = Offset(size.width / 2, size.height / 2);
    final core = size.shortestSide * 0.18;
    final span = size.shortestSide * 0.47 - core;

    _paintHalo(canvas, center, core, span);
    _paintCorona(canvas, center, core, span);
    _paintBolts(canvas, center, core, span);
  }

  void _paintHalo(Canvas canvas, Offset center, double core, double span) {
    final breath = 0.5 + 0.5 * math.sin(progress * 2 * math.pi);
    final radius = (core + span * 0.95) * (0.92 + 0.10 * breath) * ignite;

    final paint = Paint()
      ..shader = RadialGradient(
        colors: [
          _amber.withValues(alpha: 0.34 * ignite),
          _orange.withValues(alpha: 0.16 * ignite),
          _violet.withValues(alpha: 0.10 * ignite),
          Colors.transparent,
        ],
        stops: const [0.0, 0.42, 0.68, 1.0],
      ).createShader(Rect.fromCircle(center: center, radius: radius));

    canvas.drawCircle(center, radius, paint);
  }

  void _paintCorona(Canvas canvas, Offset center, double core, double span) {
    for (var i = 0; i < _tongues; i++) {
      final seed = i * 12.9898;
      final base = (i / _tongues) * 2 * math.pi;
      final drift = progress * 2 * math.pi * (i.isEven ? 0.18 : -0.13);
      final angle = base + drift;

      final flicker =
          0.5 + 0.5 * math.sin(progress * 2 * math.pi * (2 + i % 3) + seed);
      final stretch = 0.55 + _noise(seed + 2) * 0.45;
      final reach = span * stretch * (0.62 + 0.38 * flicker) * ignite;
      final width = core * (0.16 + 0.12 * _noise(seed + 6));

      final inner = center + Offset(math.cos(angle), math.sin(angle)) * core;
      final tip = center +
          Offset(math.cos(angle), math.sin(angle)) * (core + reach);
      final normal = Offset(-math.sin(angle), math.cos(angle));

      final path = Path()
        ..moveTo(inner.dx + normal.dx * width, inner.dy + normal.dy * width)
        ..quadraticBezierTo(
          center.dx + math.cos(angle + 0.16) * (core + reach * 0.55),
          center.dy + math.sin(angle + 0.16) * (core + reach * 0.55),
          tip.dx,
          tip.dy,
        )
        ..quadraticBezierTo(
          center.dx + math.cos(angle - 0.16) * (core + reach * 0.55),
          center.dy + math.sin(angle - 0.16) * (core + reach * 0.55),
          inner.dx - normal.dx * width,
          inner.dy - normal.dy * width,
        )
        ..close();

      final root = core / (core + reach);
      final mid = root + (1 - root) * 0.5;

      final paint = Paint()
        ..shader = RadialGradient(
          colors: [
            _amber.withValues(alpha: 0.95 * ignite),
            _amber.withValues(alpha: 0.95 * ignite),
            _orange.withValues(alpha: 0.60 * ignite),
            Colors.transparent,
          ],
          stops: [0.0, root, mid, 1.0],
        ).createShader(
          Rect.fromCircle(center: center, radius: core + reach),
        )
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 7);

      canvas.drawPath(path, paint);
    }
  }

  void _paintBolts(Canvas canvas, Offset center, double core, double span) {
    for (var i = 0; i < _bolts; i++) {
      final seed = i * 31.7;
      final phase = (progress * (1.7 + i % 3 * 0.4) + i / _bolts) % 1.0;
      final alive = math.sin(phase * math.pi);
      if (alive <= 0.18) continue;

      final angle = (i / _bolts) * 2 * math.pi +
          _noise(seed) * 0.9 +
          progress * 2 * math.pi * 0.12;
      final reach = span * (0.55 + 0.45 * _noise(seed + 5)) * alive;

      final trunk = _strike(center, core * 0.92, angle, reach, seed, 1.0);
      _drawBolt(canvas, trunk, alive, 2.2);

      final forkAt = 0.45 + 0.25 * _noise(seed + 9);
      final forkOrigin = center +
          Offset(math.cos(angle), math.sin(angle)) *
              (core * 0.92 + reach * forkAt);
      final forkAngle = angle + (_noise(seed + 13) - 0.5) * 1.1;
      final fork = _strike(
        forkOrigin,
        0,
        forkAngle,
        reach * (0.35 + 0.3 * _noise(seed + 17)),
        seed + 41,
        0.6,
      );
      _drawBolt(canvas, fork, alive * 0.8, 1.4);

      final tip = center +
          Offset(math.cos(angle), math.sin(angle)) * (core * 0.92 + reach);
      canvas.drawCircle(
        tip,
        2.8 * alive,
        Paint()
          ..color = Colors.white.withValues(alpha: 0.9 * alive * ignite)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 4),
      );
    }
  }

  Path _strike(
    Offset origin,
    double inset,
    double angle,
    double reach,
    double seed,
    double wobble,
  ) {
    final direction = Offset(math.cos(angle), math.sin(angle));
    final normal = Offset(-math.sin(angle), math.cos(angle));
    final path = Path();

    for (var s = 0; s <= _boltSegments; s++) {
      final t = s / _boltSegments;
      final along = origin + direction * (inset + reach * t);
      final sway = (_noise(seed + s * 5.1 + progress * 37) - 0.5) *
          reach *
          0.30 *
          wobble *
          math.sin(t * math.pi);
      final point = along + normal * sway;
      if (s == 0) {
        path.moveTo(point.dx, point.dy);
      } else {
        path.lineTo(point.dx, point.dy);
      }
    }
    return path;
  }

  void _drawBolt(Canvas canvas, Path path, double alive, double width) {
    canvas.drawPath(
      path,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = width * 3.4
        ..strokeCap = StrokeCap.round
        ..color = _violet.withValues(alpha: 0.34 * alive * ignite)
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8),
    );
    canvas.drawPath(
      path,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = width
        ..strokeCap = StrokeCap.round
        ..color = Color.lerp(_cyan, Colors.white, alive)!
            .withValues(alpha: (0.6 + 0.4 * alive) * ignite),
    );
  }

  double _noise(double seed) {
    final value = math.sin(seed * 91.7) * 43758.5453;
    return value - value.floorToDouble();
  }

  @override
  bool shouldRepaint(covariant _StarPainter oldDelegate) =>
      oldDelegate.progress != progress || oldDelegate.ignite != ignite;
}
