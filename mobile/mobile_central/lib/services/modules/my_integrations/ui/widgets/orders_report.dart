import 'package:flutter/material.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';

class OrdersReport extends StatelessWidget {
  const OrdersReport({
    super.key,
    required this.integrations,
    required this.statsFor,
  });

  final List<MyIntegration> integrations;
  final IntegrationStats Function(int id) statsFor;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    var totals = const IntegrationStats(integrationId: 0);
    for (final item in integrations) {
      totals = totals + statsFor(item.id);
    }

    final rows = [...integrations]..sort((a, b) =>
        statsFor(b.id).ordersCount.compareTo(statsFor(a.id).ordersCount));
    final maxOrders = rows.isEmpty ? 0 : statsFor(rows.first.id).ordersCount;

    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(16, 4, 16, 96),
      children: [
        Text(
          'Cuantas ordenes entraron por cada origen y en que estado van.',
          style: theme.textTheme.labelSmall,
        ),
        const SizedBox(height: 12),
        _TotalsGrid(totals: totals, origins: integrations.length),
        const SizedBox(height: 14),
        AppCard(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'COMO VAN TODAS LAS ORDENES',
                style: theme.textTheme.labelSmall?.copyWith(
                  letterSpacing: 0.8,
                  fontWeight: FontWeight.w700,
                  fontSize: 10,
                ),
              ),
              const SizedBox(height: 10),
              StatusBar(stats: totals, height: 8),
              const SizedBox(height: 10),
              StatusLegend(stats: totals),
            ],
          ),
        ),
        const SizedBox(height: 18),
        const AppSectionHeader(title: 'Por origen'),
        AppCard(
          padding: const EdgeInsets.symmetric(vertical: 4),
          child: Column(
            children: [
              for (var i = 0; i < rows.length; i++) ...[
                if (i > 0) const Divider(height: 1),
                _OriginRow(
                  integration: rows[i],
                  stats: statsFor(rows[i].id),
                  maxOrders: maxOrders,
                ),
              ],
            ],
          ),
        ),
        const SizedBox(height: 12),
        Text(
          'Cada barra usa la misma escala: la mas larga es el origen con mas '
          'ordenes.',
          style: theme.textTheme.labelSmall,
        ),
      ],
    );
  }
}

class _TotalsGrid extends StatelessWidget {
  const _TotalsGrid({required this.totals, required this.origins});

  final IntegrationStats totals;
  final int origins;

  @override
  Widget build(BuildContext context) {
    final total = totals.ordersCount;
    String pct(int value) {
      if (total == 0) return '0% del total';
      return '${(value * 100 / total).round()}% del total';
    }

    return Column(
      children: [
        _ReportTile(
          label: 'ORDENES',
          value: AppFormat.number(total),
          caption: 'en $origins origen${origins == 1 ? '' : 'es'}',
        ),
        const SizedBox(height: 9),
        Row(
          children: [
            Expanded(
              child: _ReportTile(
                label: 'EN CURSO',
                value: AppFormat.number(totals.ordersInProgress),
                caption: pct(totals.ordersInProgress),
                color: AppColors.primary,
              ),
            ),
            const SizedBox(width: 9),
            Expanded(
              child: _ReportTile(
                label: 'ENTREGADAS',
                value: AppFormat.number(totals.ordersDelivered),
                caption: pct(totals.ordersDelivered),
                color: AppColors.success,
              ),
            ),
          ],
        ),
        const SizedBox(height: 9),
        Row(
          children: [
            Expanded(
              child: _ReportTile(
                label: 'CANCELADAS',
                value: AppFormat.number(totals.ordersCancelled),
                caption: pct(totals.ordersCancelled),
                color: AppColors.error,
              ),
            ),
            const SizedBox(width: 9),
            Expanded(
              child: _ReportTile(
                label: 'DEVUELTAS',
                value: AppFormat.number(totals.ordersReturned),
                caption: pct(totals.ordersReturned),
                color: AppColors.warning,
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _ReportTile extends StatelessWidget {
  const _ReportTile({
    required this.label,
    required this.value,
    required this.caption,
    this.color,
  });

  final String label;
  final String value;
  final String caption;
  final Color? color;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: AppRadius.lgAll,
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              if (color != null) ...[
                Container(
                  width: 6,
                  height: 6,
                  decoration: BoxDecoration(color: color, shape: BoxShape.circle),
                ),
                const SizedBox(width: 6),
              ],
              Text(
                label,
                style: theme.textTheme.labelSmall?.copyWith(
                  fontSize: 9.5,
                  letterSpacing: 0.8,
                  fontWeight: FontWeight.w700,
                  color: color ?? AppColors.textMuted,
                ),
              ),
            ],
          ),
          const SizedBox(height: 6),
          Text(value, style: theme.textTheme.headlineSmall),
          const SizedBox(height: 2),
          Text(caption, style: theme.textTheme.labelSmall?.copyWith(fontSize: 10)),
        ],
      ),
    );
  }
}

class _OriginRow extends StatelessWidget {
  const _OriginRow({
    required this.integration,
    required this.stats,
    required this.maxOrders,
  });

  final MyIntegration integration;
  final IntegrationStats stats;
  final int maxOrders;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = integration.integrationTypeName ?? integration.name;
    final ratio = maxOrders == 0 ? 0.0 : stats.ordersCount / maxOrders;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 11),
      child: Column(
        children: [
          Row(
            children: [
              BrandLogo(
                name: name,
                imageUrl: integration.imageUrl,
                size: 28,
                radius: 7,
                padding: 4,
              ),
              const SizedBox(width: 9),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      name,
                      style: theme.textTheme.titleSmall?.copyWith(fontSize: 13),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    if (stats.lastOrderAt != null)
                      Text(
                        'Ult. orden ${AppFormat.relative(AppFormat.parseDate(stats.lastOrderAt!))}',
                        style: theme.textTheme.labelSmall?.copyWith(fontSize: 10),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    AppFormat.number(stats.ordersCount),
                    style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
                  ),
                  Text(
                    '${AppFormat.number(stats.productsCount)} prod.',
                    style: theme.textTheme.labelSmall?.copyWith(fontSize: 10),
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 8),
          if (stats.ordersCount == 0)
            Align(
              alignment: Alignment.centerLeft,
              child: Text(
                'Sin ordenes registradas',
                style: theme.textTheme.labelSmall
                    ?.copyWith(fontSize: 10, fontStyle: FontStyle.italic),
              ),
            )
          else
            Row(
              children: [
                Expanded(
                  flex: (ratio * 1000).round().clamp(1, 1000),
                  child: StatusBar(stats: stats, height: 6),
                ),
                if (ratio < 0.999)
                  Expanded(
                    flex: ((1 - ratio) * 1000).round().clamp(1, 1000),
                    child: const SizedBox(height: 6),
                  ),
              ],
            ),
        ],
      ),
    );
  }
}

class StatusBar extends StatelessWidget {
  const StatusBar({super.key, required this.stats, this.height = 6});

  final IntegrationStats stats;
  final double height;

  @override
  Widget build(BuildContext context) {
    final parts = <({int count, Color color})>[
      (count: stats.ordersInProgress, color: AppColors.primary),
      (count: stats.ordersDelivered, color: AppColors.success),
      (count: stats.ordersCancelled, color: AppColors.error),
      (count: stats.ordersReturned, color: AppColors.warning),
    ].where((p) => p.count > 0).toList();

    return ClipRRect(
      borderRadius: BorderRadius.circular(999),
      child: SizedBox(
        height: height,
        child: parts.isEmpty
            ? Container(color: AppColors.surfaceMuted)
            : Row(
                children: [
                  for (final part in parts)
                    Expanded(
                      flex: part.count.clamp(1, 1 << 30),
                      child: Container(color: part.color),
                    ),
                ],
              ),
      ),
    );
  }
}

class StatusLegend extends StatelessWidget {
  const StatusLegend({super.key, required this.stats});

  final IntegrationStats stats;

  @override
  Widget build(BuildContext context) {
    final items = <({int count, String label, Color color})>[
      (count: stats.ordersInProgress, label: 'en curso', color: AppColors.primary),
      (count: stats.ordersDelivered, label: 'entregadas', color: AppColors.success),
      (count: stats.ordersCancelled, label: 'canceladas', color: AppColors.error),
      (count: stats.ordersReturned, label: 'devueltas', color: AppColors.warning),
    ];

    return Wrap(
      spacing: 12,
      runSpacing: 6,
      children: [
        for (final item in items)
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 6,
                height: 6,
                decoration:
                    BoxDecoration(color: item.color, shape: BoxShape.circle),
              ),
              const SizedBox(width: 5),
              Text(
                '${AppFormat.number(item.count)} ${item.label}',
                style: Theme.of(context)
                    .textTheme
                    .labelSmall
                    ?.copyWith(fontSize: 11),
              ),
            ],
          ),
      ],
    );
  }
}
