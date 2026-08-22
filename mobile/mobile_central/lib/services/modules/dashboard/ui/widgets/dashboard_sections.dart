import 'package:flutter/material.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class DashboardRankRow extends StatelessWidget {
  const DashboardRankRow({
    super.key,
    required this.leading,
    required this.title,
    required this.subtitle,
    required this.value,
    this.valueCaption,
    this.progress,
  });

  final Widget leading;
  final String title;
  final String subtitle;
  final String value;
  final String? valueCaption;
  final double? progress;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 9),
      child: Row(
        children: [
          leading,
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  subtitle,
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                if (progress != null) ...[
                  const SizedBox(height: 7),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(3),
                    child: LinearProgressIndicator(
                      value: progress!.clamp(0.02, 1),
                      minHeight: 4,
                      backgroundColor: AppColors.surfaceMuted,
                      valueColor: AlwaysStoppedAnimation(
                          Theme.of(context).colorScheme.primary),
                    ),
                  ),
                ],
              ],
            ),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                value,
                style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
              ),
              if (valueCaption != null) ...[
                const SizedBox(height: 2),
                Text(valueCaption!, style: theme.textTheme.labelSmall),
              ],
            ],
          ),
        ],
      ),
    );
  }
}

class DashboardRankBadge extends StatelessWidget {
  const DashboardRankBadge({super.key, required this.position});

  final int position;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 30,
      height: 30,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: position == 1
            ? Theme.of(context).colorScheme.primary
            : AppColors.surfaceMuted,
        borderRadius: AppRadius.smAll,
      ),
      child: Text(
        '$position',
        style: TextStyle(
          fontFamily: 'Inter',
          fontSize: 12,
          fontWeight: FontWeight.w700,
          color: position == 1 ? Colors.white : AppColors.textMuted,
        ),
      ),
    );
  }
}

class DashboardChannelsSection extends StatelessWidget {
  const DashboardChannelsSection({super.key, required this.items});

  final List<OrderCountByIntegrationType> items;

  @override
  Widget build(BuildContext context) {
    final sorted = [...items]..sort((a, b) => b.count.compareTo(a.count));
    final max = sorted.isEmpty ? 1 : sorted.first.count;

    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < sorted.length; i++) ...[
            if (i > 0) const Divider(height: 1),
            DashboardRankRow(
              leading: BrandLogo(name: sorted[i].integrationType, size: 34, radius: 9, padding: 5),
              title: sorted[i].integrationType,
              subtitle: 'Canal de venta',
              value: AppFormat.number(sorted[i].count),
              valueCaption: 'ordenes',
              progress: max == 0 ? 0 : sorted[i].count / max,
            ),
          ],
        ],
      ),
    );
  }
}

class DashboardCarriersSection extends StatelessWidget {
  const DashboardCarriersSection({super.key, required this.items});

  final List<ShipmentsByCarrier> items;

  @override
  Widget build(BuildContext context) {
    final sorted = [...items]..sort((a, b) => b.count.compareTo(a.count));
    final max = sorted.isEmpty ? 1 : sorted.first.count;

    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < sorted.length; i++) ...[
            if (i > 0) const Divider(height: 1),
            DashboardRankRow(
              leading: BrandLogo(name: sorted[i].carrier, size: 34, radius: 9, padding: 5),
              title: sorted[i].carrier,
              subtitle: 'Transportadora',
              value: AppFormat.number(sorted[i].count),
              valueCaption: 'guias',
              progress: max == 0 ? 0 : sorted[i].count / max,
            ),
          ],
        ],
      ),
    );
  }
}

class DashboardCustomersSection extends StatelessWidget {
  const DashboardCustomersSection({super.key, required this.items});

  final List<TopCustomer> items;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < items.length; i++) ...[
            if (i > 0) const Divider(height: 1),
            DashboardRankRow(
              leading: _InitialsAvatar(name: items[i].customerName),
              title: items[i].customerName,
              subtitle: items[i].customerEmail,
              value: AppFormat.number(items[i].orderCount),
              valueCaption: 'ordenes',
            ),
          ],
        ],
      ),
    );
  }
}

class DashboardProductsSection extends StatelessWidget {
  const DashboardProductsSection({super.key, required this.items});

  final List<TopProduct> items;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < items.length; i++) ...[
            if (i > 0) const Divider(height: 1),
            DashboardRankRow(
              leading: DashboardRankBadge(position: i + 1),
              title: items[i].productName,
              subtitle: 'SKU ${items[i].sku}',
              value: AppFormat.number(items[i].totalSold),
              valueCaption: 'unidades',
            ),
          ],
        ],
      ),
    );
  }
}

class DashboardStatusSection extends StatelessWidget {
  const DashboardStatusSection({super.key, required this.items});

  final List<ShipmentsByStatus> items;

  static const Map<String, String> _labels = {
    'created': 'Creada',
    'in_transit': 'En transito',
    'delivered': 'Entregada',
    'returned': 'Devuelta',
    'cancelled': 'Cancelada',
    'pending': 'Pendiente',
  };

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: items
          .map((item) => Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 9),
                decoration: BoxDecoration(
                  color: AppColors.surface,
                  borderRadius: AppRadius.pillAll,
                  border: Border.all(color: AppColors.border),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Container(
                      width: 7,
                      height: 7,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: _dotColor(item.status),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      _labels[item.status] ?? item.status,
                      style: Theme.of(context).textTheme.labelMedium?.copyWith(
                            color: AppColors.textSecondary,
                            fontWeight: FontWeight.w600,
                          ),
                    ),
                    const SizedBox(width: 7),
                    Text(
                      '${item.count}',
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(fontSize: 13),
                    ),
                  ],
                ),
              ))
          .toList(),
    );
  }

  Color _dotColor(String status) {
    switch (AppStatusChip.toneFromCode(status)) {
      case AppStatusTone.success:
        return AppColors.success;
      case AppStatusTone.warning:
        return AppColors.warning;
      case AppStatusTone.error:
        return AppColors.error;
      case AppStatusTone.info:
        return AppColors.info;
      case AppStatusTone.brand:
        return AppColors.primary;
      case AppStatusTone.neutral:
        return AppColors.textDisabled;
    }
  }
}

class DashboardLocationsSection extends StatelessWidget {
  const DashboardLocationsSection({super.key, required this.items});

  final List<OrderCountByLocation> items;

  @override
  Widget build(BuildContext context) {
    final sorted = [...items]..sort((a, b) => b.orderCount.compareTo(a.orderCount));
    final top = sorted.take(6).toList();
    final max = top.isEmpty ? 1 : top.first.orderCount;

    return AppCard(
      child: Column(
        children: [
          for (var i = 0; i < top.length; i++) ...[
            if (i > 0) const Divider(height: 1),
            DashboardRankRow(
              leading: const Icon(Icons.place_outlined, size: 20, color: AppColors.textMuted),
              title: top[i].city,
              subtitle: top[i].state,
              value: AppFormat.number(top[i].orderCount),
              valueCaption: 'ordenes',
              progress: max == 0 ? 0 : top[i].orderCount / max,
            ),
          ],
        ],
      ),
    );
  }
}

class _InitialsAvatar extends StatelessWidget {
  const _InitialsAvatar({required this.name});

  final String name;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    return Container(
      width: 34,
      height: 34,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        color: scheme.primaryContainer,
        shape: BoxShape.circle,
      ),
      child: Text(
        AppFormat.initials(name),
        style: TextStyle(
          fontFamily: 'Inter',
          fontSize: 12,
          fontWeight: FontWeight.w700,
          color: scheme.primary,
        ),
      ),
    );
  }
}
