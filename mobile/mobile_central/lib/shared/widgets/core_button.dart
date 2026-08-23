import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

class CoreButton extends StatelessWidget {
  const CoreButton({super.key, required this.onTap, this.size = 62});

  final VoidCallback onTap;
  final double size;

  @override
  Widget build(BuildContext context) {
    final brand = Theme.of(context).colorScheme.primary;

    return SizedBox(
      width: size + 26,
      height: size + 26,
      child: Center(
        child: Material(
          color: Colors.transparent,
          shape: const CircleBorder(),
          clipBehavior: Clip.antiAlias,
          child: InkWell(
            onTap: onTap,
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
    );
  }
}
