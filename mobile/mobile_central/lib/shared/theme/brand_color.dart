import 'package:flutter/material.dart';

import 'app_colors.dart';

class BrandColor {
  const BrandColor._();

  static Color? parse(String? hex) {
    if (hex == null) return null;
    var value = hex.trim().replaceAll('#', '');
    if (value.length == 3) {
      value = value.split('').map((c) => '$c$c').join();
    }
    if (value.length == 6) value = 'FF$value';
    if (value.length != 8) return null;
    final parsed = int.tryParse(value, radix: 16);
    if (parsed == null) return null;
    return Color(parsed);
  }

  static Color resolve(String? hex) => parse(hex) ?? AppColors.primary;

  static Color soft(Color base) =>
      Color.alphaBlend(base.withValues(alpha: 0.12), Colors.white);

  static Color dark(Color base) {
    final hsl = HSLColor.fromColor(base);
    return hsl.withLightness((hsl.lightness - 0.16).clamp(0.0, 1.0)).toColor();
  }

  static Color onBrand(Color base) =>
      base.computeLuminance() > 0.6 ? AppColors.textPrimary : Colors.white;
}
