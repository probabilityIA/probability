import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../services/auth/business/ui/providers/business_provider.dart';
import 'business_switcher.dart';

class DraggableBusinessButton extends StatefulWidget {
  const DraggableBusinessButton({super.key});

  @override
  State<DraggableBusinessButton> createState() =>
      _DraggableBusinessButtonState();
}

class _DraggableBusinessButtonState extends State<DraggableBusinessButton> {
  static const double _height = 46;
  static const double _peek = 24;
  static const double _margin = 8;
  static const Duration _hideAfter = Duration(seconds: 3);

  Offset? _position;
  bool _dragging = false;
  bool _hidden = true;
  bool _onRight = true;
  Timer? _autoHide;

  @override
  void dispose() {
    _autoHide?.cancel();
    super.dispose();
  }

  void _scheduleHide() {
    _autoHide?.cancel();
    _autoHide = Timer(_hideAfter, () {
      if (mounted) setState(() => _hidden = true);
    });
  }

  void _onTap() {
    if (_hidden) {
      setState(() => _hidden = false);
      _scheduleHide();
      return;
    }
    _autoHide?.cancel();
    showBusinessSwitcher(context).then((_) {
      if (mounted) setState(() => _hidden = true);
    });
  }

  double _width(String label) => 74 + label.length * 7.2;

  @override
  Widget build(BuildContext context) {
    return Consumer<BusinessProvider>(
      builder: (context, provider, _) {
        final selected = provider.selectedBusinessId;
        final business = selected == null
            ? null
            : provider.businessesSimple
                .where((b) => b.id == selected)
                .firstOrNull;
        final raw = business?.name ?? 'Elegir negocio';
        final label = raw.length > 16 ? '${raw.substring(0, 15)}…' : raw;

        return LayoutBuilder(
          builder: (context, constraints) {
            final width = _width(label).clamp(96.0, constraints.maxWidth - 24);
            final maxTop = constraints.maxHeight - _height - _margin;
            final position = _position ??
                Offset(constraints.maxWidth - width - _margin, maxTop * 0.72);

            var left = position.dx;
            if (!_dragging) {
              left = _onRight
                  ? constraints.maxWidth - (_hidden ? _peek : width + _margin)
                  : (_hidden ? _peek - width : _margin);
            }
            final top = position.dy.clamp(_margin, maxTop <= 0 ? 0.0 : maxTop);

            return Stack(
              children: [
                AnimatedPositioned(
                  duration: Duration(milliseconds: _dragging ? 0 : 220),
                  curve: Curves.easeOutCubic,
                  left: left,
                  top: top,
                  child: GestureDetector(
                    onTap: _onTap,
                    onPanStart: (_) {
                      _autoHide?.cancel();
                      setState(() {
                        _dragging = true;
                        _hidden = false;
                        _position = Offset(
                          _onRight
                              ? constraints.maxWidth - width - _margin
                              : _margin,
                          top,
                        );
                      });
                    },
                    onPanUpdate: (details) {
                      setState(() {
                        _position = Offset(
                          (_position!.dx + details.delta.dx)
                              .clamp(-width, constraints.maxWidth),
                          (_position!.dy + details.delta.dy)
                              .clamp(_margin, maxTop <= 0 ? 0.0 : maxTop),
                        );
                      });
                    },
                    onPanEnd: (_) {
                      final center = _position!.dx + width / 2;
                      setState(() {
                        _onRight = center > constraints.maxWidth / 2;
                        _dragging = false;
                        _hidden = true;
                      });
                    },
                    child: _Pill(
                      label: label,
                      width: width,
                      height: _height,
                      hidden: _hidden && !_dragging,
                      onRight: _onRight,
                    ),
                  ),
                ),
              ],
            );
          },
        );
      },
    );
  }
}

class _Pill extends StatelessWidget {
  const _Pill({
    required this.label,
    required this.width,
    required this.height,
    required this.hidden,
    required this.onRight,
  });

  final String label;
  final double width;
  final double height;
  final bool hidden;
  final bool onRight;

  @override
  Widget build(BuildContext context) {
    final radius = Radius.circular(height / 2);
    final brand = Theme.of(context).colorScheme.primary;

    if (hidden) {
      return Container(
        width: width,
        height: height,
        decoration: BoxDecoration(
          color: brand,
          borderRadius: BorderRadius.horizontal(
            left: onRight ? radius : Radius.zero,
            right: onRight ? Radius.zero : radius,
          ),
          boxShadow: [
            BoxShadow(
              color: brand.withValues(alpha: 0.28),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: Align(
          alignment: onRight ? Alignment.centerLeft : Alignment.centerRight,
          child: const Padding(
            padding: EdgeInsets.symmetric(horizontal: 3),
            child: Icon(Icons.swap_horiz_rounded, size: 17, color: Colors.white),
          ),
        ),
      );
    }

    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        color: brand,
        borderRadius: BorderRadius.all(radius),
        boxShadow: [
          BoxShadow(
            color: brand.withValues(alpha: 0.28),
            blurRadius: 12,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.swap_horiz_rounded, size: 18, color: Colors.white),
          const SizedBox(width: 8),
          Flexible(
            child: Text(
              label,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              softWrap: false,
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: Colors.white,
              ),
            ),
          ),
          const SizedBox(width: 4),
        ],
      ),
    );
  }
}
