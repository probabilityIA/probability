import 'package:flutter/material.dart';
import '../../theme/app_colors.dart';

enum AppLogoVariant { mark, wordmark }

class AppLogo extends StatelessWidget {
  const AppLogo({
    super.key,
    this.variant = AppLogoVariant.wordmark,
    this.height = 28,
  });

  final AppLogoVariant variant;
  final double height;

  @override
  Widget build(BuildContext context) {
    final asset = variant == AppLogoVariant.mark
        ? 'assets/images/logo_probability.png'
        : 'assets/images/logo_probability_wide.png';
    return Image.asset(
      asset,
      height: height,
      fit: BoxFit.contain,
      filterQuality: FilterQuality.high,
      errorBuilder: (context, error, stack) => Text(
        'ProbabilityIA',
        style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: AppColors.primary,
              fontWeight: FontWeight.w700,
            ),
      ),
    );
  }
}

class AppLogoBadge extends StatelessWidget {
  const AppLogoBadge({super.key, this.size = 40});

  final double size;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: AppColors.primarySoft,
        borderRadius: BorderRadius.circular(size * 0.28),
      ),
      padding: EdgeInsets.all(size * 0.14),
      child: const AppLogo(variant: AppLogoVariant.mark),
    );
  }
}
