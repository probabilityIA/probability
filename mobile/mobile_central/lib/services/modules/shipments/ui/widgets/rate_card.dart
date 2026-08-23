import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/rate_pricing.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class RateCard extends StatelessWidget {
  const RateCard({
    super.key,
    required this.rate,
    required this.options,
    required this.selected,
    required this.onTap,
  });

  final EnvioClickRate rate;
  final RatePricingOptions options;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final breakdown = RatePricing.breakdown(rate, options);
    final appliesCod = RatePricing.appliesCod(rate, options);

    return AppCard(
      padding: const EdgeInsets.all(14),
      borderColor: selected ? AppColors.primary : AppColors.border,
      color: selected ? AppColors.primarySoft.withValues(alpha: 0.35) : AppColors.surface,
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(name: rate.carrier, size: 40, radius: 11, padding: 6),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      rate.carrier,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${rate.product}  ·  ${rate.deliveryDays} dias',
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    AppFormat.money(breakdown.total),
                    style: theme.textTheme.titleMedium?.copyWith(
                      color: selected ? AppColors.primaryDark : AppColors.textPrimary,
                    ),
                  ),
                  Text('total', style: theme.textTheme.labelSmall),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
            decoration: BoxDecoration(
              color: AppColors.surfaceMuted,
              borderRadius: AppRadius.mdAll,
            ),
            child: Column(
              children: [
                _Line(label: 'Flete', value: breakdown.flete),
                if (breakdown.minimumInsurance > 0)
                  _Line(label: 'Seguro minimo', value: breakdown.minimumInsurance),
                if (breakdown.extraInsurance > 0)
                  _Line(label: 'Seguro extra', value: breakdown.extraInsurance),
                const Divider(height: 14),
                _Line(label: 'Costo de la guia', value: breakdown.guideCost, strong: true),
                if (appliesCod)
                  _Line(label: 'Comision contra entrega', value: breakdown.carrierFee),
              ],
            ),
          ),
          if (options.cod && rate.cod != true) ...[
            const SizedBox(height: 10),
            Row(
              children: [
                const Icon(Icons.info_outline, size: 14, color: AppColors.warning),
                const SizedBox(width: 6),
                Expanded(
                  child: Text(
                    'Esta tarifa no soporta contra entrega',
                    style: theme.textTheme.labelSmall?.copyWith(color: const Color(0xFFB45309)),
                  ),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }
}

class _Line extends StatelessWidget {
  const _Line({required this.label, required this.value, this.strong = false});

  final String label;
  final double value;
  final bool strong;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        children: [
          Expanded(
            child: Text(
              label,
              style: strong
                  ? theme.textTheme.titleSmall?.copyWith(fontSize: 13)
                  : theme.textTheme.labelSmall,
            ),
          ),
          Text(
            AppFormat.money(value),
            style: strong
                ? theme.textTheme.titleSmall?.copyWith(fontSize: 13)
                : theme.textTheme.labelMedium?.copyWith(color: AppColors.textSecondary),
          ),
        ],
      ),
    );
  }
}
