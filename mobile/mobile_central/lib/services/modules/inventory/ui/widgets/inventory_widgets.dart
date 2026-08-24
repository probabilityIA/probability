import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class InventoryLevelCard extends StatelessWidget {
  const InventoryLevelCard({super.key, required this.level, this.onTap});

  final InventoryLevel level;
  final VoidCallback? onTap;

  ({String label, AppStatusTone tone}) get _badge {
    final reorder = level.reorderPoint ?? 0;
    if (level.availableQty <= 0) return (label: 'Agotado', tone: AppStatusTone.error);
    if (reorder > 0 && level.availableQty <= reorder) {
      return (label: 'Reponer', tone: AppStatusTone.warning);
    }
    return (label: 'Disponible', tone: AppStatusTone.success);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final badge = _badge;
    final max = level.maxStock ?? 0;
    final ratio = max > 0 ? (level.quantity / max).clamp(0.0, 1.0) : null;

    return AppCard(
      padding: const EdgeInsets.all(13),
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      level.productName ?? 'Producto',
                      style: theme.textTheme.titleSmall,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 3),
                    Text('SKU ${level.productSku ?? '-'}', style: theme.textTheme.labelSmall),
                  ],
                ),
              ),
              const SizedBox(width: 10),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    AppFormat.number(level.availableQty),
                    style: theme.textTheme.titleMedium,
                  ),
                  Text('disponibles', style: theme.textTheme.labelSmall),
                ],
              ),
            ],
          ),
          const SizedBox(height: 10),
          if (ratio != null) ...[
            ClipRRect(
              borderRadius: BorderRadius.circular(3),
              child: LinearProgressIndicator(
                value: ratio,
                minHeight: 4,
                backgroundColor: AppColors.surfaceMuted,
                valueColor: AlwaysStoppedAnimation(
                  badge.tone == AppStatusTone.error
                      ? AppColors.error
                      : badge.tone == AppStatusTone.warning
                          ? AppColors.warning
                          : AppColors.success,
                ),
              ),
            ),
            const SizedBox(height: 10),
          ],
          Row(
            children: [
              AppStatusChip(dense: true, label: badge.label, tone: badge.tone),
              const SizedBox(width: 8),
              Text(
                'Total ${AppFormat.number(level.quantity)}  ·  Reservado ${AppFormat.number(level.reservedQty)}',
                style: theme.textTheme.labelSmall,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class StockMovementTile extends StatelessWidget {
  const StockMovementTile({super.key, required this.movement});

  final StockMovement movement;

  bool get _isIn => movement.quantity >= 0;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _isIn ? AppColors.success : AppColors.error;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 34,
            height: 34,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.11),
              borderRadius: AppRadius.smAll,
            ),
            child: Icon(
              _isIn ? Icons.south_west_rounded : Icons.north_east_rounded,
              size: 17,
              color: color,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  movement.productName ?? 'Producto',
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  '${movement.movementTypeName}  ·  ${movement.warehouseName ?? ''}',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  AppFormat.relative(AppFormat.parseDate(movement.createdAt)),
                  style: theme.textTheme.labelSmall,
                ),
              ],
            ),
          ),
          const SizedBox(width: 10),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '${_isIn ? '+' : ''}${AppFormat.number(movement.quantity)}',
                style: theme.textTheme.titleSmall?.copyWith(color: color),
              ),
              const SizedBox(height: 2),
              Text(
                '${AppFormat.number(movement.previousQty)} a ${AppFormat.number(movement.newQty)}',
                style: theme.textTheme.labelSmall,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
