import 'package:flutter/material.dart';
import '../../theme/app_colors.dart';
import '../../theme/app_tokens.dart';

class AppKpiTile extends StatelessWidget {
  const AppKpiTile({
    super.key,
    required this.label,
    required this.value,
    this.icon,
    this.trend,
    this.trendPositive = true,
    this.accent = AppColors.primary,
    this.onTap,
    this.footer,
  });

  final String label;
  final String value;
  final IconData? icon;
  final String? trend;
  final bool trendPositive;
  final Color accent;
  final VoidCallback? onTap;
  final Widget? footer;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final content = Container(
      padding: const EdgeInsets.all(11),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: AppRadius.lgAll,
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            children: [
              if (icon != null) ...[
                Container(
                  width: 26,
                  height: 26,
                  decoration: BoxDecoration(
                    color: accent.withValues(alpha: 0.10),
                    borderRadius: AppRadius.smAll,
                  ),
                  child: Icon(icon, size: 15, color: accent),
                ),
                const SizedBox(width: 9),
              ],
              Expanded(
                child: Text(
                  label,
                  style: theme.textTheme.labelMedium,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            value,
            style: theme.textTheme.headlineSmall?.copyWith(letterSpacing: -0.5),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          if (footer != null) ...[
            const SizedBox(height: 6),
            footer!,
          ],
          if (trend != null) ...[
            const SizedBox(height: 4),
            Row(
              children: [
                Icon(
                  trendPositive ? Icons.trending_up : Icons.trending_down,
                  size: 14,
                  color: trendPositive ? AppColors.success : AppColors.error,
                ),
                const SizedBox(width: 4),
                Text(
                  trend!,
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: trendPositive ? AppColors.success : AppColors.error,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ],
        ],
      ),
    );

    if (onTap == null) return content;
    return Material(
      color: Colors.transparent,
      child: InkWell(onTap: onTap, borderRadius: AppRadius.lgAll, child: content),
    );
  }
}
