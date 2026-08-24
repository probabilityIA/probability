import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/theme/app_colors.dart';
import 'package:mobile_central/shared/theme/app_theme.dart';
import 'package:mobile_central/shared/theme/brand_color.dart';

void main() {
  group('BrandColor.parse', () {
    test('lee el formato que manda el API', () {
      expect(BrandColor.parse('#DC2626'), const Color(0xFFDC2626));
      expect(BrandColor.parse('DC2626'), const Color(0xFFDC2626));
      expect(BrandColor.parse('#a363e0'), const Color(0xFFA363E0));
    });

    test('acepta la forma corta de tres digitos', () {
      expect(BrandColor.parse('#f00'), const Color(0xFFFF0000));
    });

    test('devuelve nulo si el color no sirve', () {
      expect(BrandColor.parse(null), isNull);
      expect(BrandColor.parse(''), isNull);
      expect(BrandColor.parse('rojo'), isNull);
      expect(BrandColor.parse('#12345'), isNull);
    });

    test('resolve cae a la marca Probability', () {
      expect(BrandColor.resolve(null), AppColors.primary);
      expect(BrandColor.resolve('no-es-color'), AppColors.primary);
      expect(BrandColor.resolve('#DC2626'), const Color(0xFFDC2626));
    });
  });

  group('contraste', () {
    test('un color oscuro lleva texto blanco', () {
      expect(BrandColor.onBrand(const Color(0xFF0F172A)), Colors.white);
    });

    test('un color muy claro lleva texto oscuro', () {
      expect(BrandColor.onBrand(const Color(0xFFFDE68A)), AppColors.textPrimary);
    });
  });

  group('tema por negocio', () {
    test('el tema toma el color del negocio', () {
      final theme = AppTheme.lightFor(const Color(0xFFDC2626));

      expect(theme.colorScheme.primary, const Color(0xFFDC2626));
      expect(theme.tabBarTheme.labelColor, const Color(0xFFDC2626));
    });

    test('sin negocio queda la marca Probability', () {
      expect(AppTheme.light.colorScheme.primary, AppColors.primary);
    });

    test('los fondos NO cambian con el negocio', () {
      final rojo = AppTheme.lightFor(const Color(0xFFDC2626));
      final base = AppTheme.light;

      expect(rojo.scaffoldBackgroundColor, base.scaffoldBackgroundColor);
      expect(rojo.colorScheme.surface, base.colorScheme.surface);
      expect(rojo.colorScheme.onSurface, base.colorScheme.onSurface);
    });
  });
}
