import 'package:flutter/material.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';

enum AppStatusTone { neutral, success, warning, error, info, brand }

class AppStatusChip extends StatelessWidget {
  const AppStatusChip({
    super.key,
    required this.label,
    this.tone = AppStatusTone.neutral,
    this.icon,
    this.dense = false,
  });

  final String label;
  final AppStatusTone tone;
  final IconData? icon;
  final bool dense;

  static AppStatusTone toneFromCode(String? code) {
    final value = (code ?? '').toLowerCase();
    if (value.contains('cancel') || value.contains('fail') ||
        value.contains('error') || value.contains('reject') ||
        value.contains('rechaz') || value.contains('anul')) {
      return AppStatusTone.error;
    }
    if (value.contains('deliver') || value.contains('entreg') ||
        value.contains('paid') || value.contains('pagad') ||
        value.contains('complet') || value.contains('success') ||
        value.contains('active') || value.contains('activo') ||
        value.contains('emitid') || value.contains('confirm')) {
      return AppStatusTone.success;
    }
    if (value.contains('pend') || value.contains('wait') ||
        value.contains('espera') || value.contains('draft') ||
        value.contains('borrador') || value.contains('retry')) {
      return AppStatusTone.warning;
    }
    if (value.contains('transit') || value.contains('ship') ||
        value.contains('envi') || value.contains('proces') ||
        value.contains('progress') || value.contains('sync')) {
      return AppStatusTone.info;
    }
    return AppStatusTone.neutral;
  }

  ({Color bg, Color fg}) get _palette {
    switch (tone) {
      case AppStatusTone.success:
        return (bg: AppColors.successSoft, fg: const Color(0xFF047857));
      case AppStatusTone.warning:
        return (bg: AppColors.warningSoft, fg: const Color(0xFFB45309));
      case AppStatusTone.error:
        return (bg: AppColors.errorSoft, fg: const Color(0xFFB91C1C));
      case AppStatusTone.info:
        return (bg: AppColors.infoSoft, fg: const Color(0xFF1D4ED8));
      case AppStatusTone.brand:
        return (bg: AppColors.primarySoft, fg: AppColors.primaryDark);
      case AppStatusTone.neutral:
        return (bg: AppColors.surfaceMuted, fg: AppColors.textSecondary);
    }
  }

  @override
  Widget build(BuildContext context) {
    final palette = _palette;
    return Container(
      padding: EdgeInsets.symmetric(
        horizontal: dense ? 8 : 10,
        vertical: dense ? 3 : 5,
      ),
      decoration: BoxDecoration(
        color: palette.bg,
        borderRadius: AppRadius.pillAll,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[
            Icon(icon, size: dense ? 12 : 14, color: palette.fg),
            const SizedBox(width: 5),
          ],
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: dense ? 11 : 12,
              fontWeight: FontWeight.w600,
              color: palette.fg,
              height: 1.2,
            ),
          ),
        ],
      ),
    );
  }
}
