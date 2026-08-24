import 'package:flutter/material.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import 'channel_sync_strip.dart';

class CoreCard extends StatelessWidget {
  const CoreCard({super.key, required this.totals});

  final IntegrationStats totals;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final brand = theme.colorScheme.primary;

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 22, horizontal: 18),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppRadius.xl),
        border: Border.all(color: AppColors.border),
        boxShadow: [
          BoxShadow(
            color: brand.withValues(alpha: 0.10),
            blurRadius: 26,
            offset: const Offset(0, 8),
          ),
        ],
      ),
      child: Column(
        children: [
          Text(
            'N U C L E O',
            style: theme.textTheme.labelSmall?.copyWith(
              letterSpacing: 3,
              fontSize: 9.5,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 6),
          const AppLogo(height: 26),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _CoreStat(
                value: AppFormat.number(totals.ordersCount),
                label: 'ordenes',
              ),
              Container(
                width: 1,
                height: 30,
                margin: const EdgeInsets.symmetric(horizontal: 22),
                color: AppColors.border,
              ),
              _CoreStat(
                value: AppFormat.number(totals.productsCount),
                label: 'productos',
              ),
            ],
          ),
          const SizedBox(height: 16),
          Wrap(
            spacing: 7,
            runSpacing: 7,
            alignment: WrapAlignment.center,
            children: [
              _StatDot(
                count: totals.ordersInProgress,
                label: 'en curso',
                color: AppColors.primary,
              ),
              _StatDot(
                count: totals.ordersDelivered,
                label: 'entregadas',
                color: AppColors.success,
              ),
              _StatDot(
                count: totals.ordersCancelled,
                label: 'canceladas',
                color: AppColors.error,
              ),
              _StatDot(
                count: totals.ordersReturned,
                label: 'devueltas',
                color: AppColors.warning,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _CoreStat extends StatelessWidget {
  const _CoreStat({required this.value, required this.label});

  final String value;
  final String label;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      children: [
        Text(value, style: theme.textTheme.headlineSmall),
        const SizedBox(height: 1),
        Text(label, style: theme.textTheme.labelSmall),
      ],
    );
  }
}

class _StatDot extends StatelessWidget {
  const _StatDot({
    required this.count,
    required this.label,
    required this.color,
  });

  final int count;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 5),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: AppRadius.pillAll,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 6,
            height: 6,
            decoration: BoxDecoration(color: color, shape: BoxShape.circle),
          ),
          const SizedBox(width: 6),
          Text(
            '${AppFormat.number(count)} $label',
            style: Theme.of(context)
                .textTheme
                .labelSmall
                ?.copyWith(fontSize: 11),
          ),
        ],
      ),
    );
  }
}

enum ChannelAction {
  toggle,
  test,
  sync;

  String get runningLabel {
    switch (this) {
      case ChannelAction.toggle:
        return 'Cambiando estado';
      case ChannelAction.test:
        return 'Probando conexion';
      case ChannelAction.sync:
        return 'Sincronizando ordenes';
    }
  }

  String get doneLabel {
    switch (this) {
      case ChannelAction.toggle:
        return 'Estado actualizado';
      case ChannelAction.test:
        return 'La conexion funciona';
      case ChannelAction.sync:
        return 'Sincronizacion lanzada';
    }
  }
}

typedef ChannelActionHandler = Future<void> Function(
  MyIntegration integration,
  ChannelAction action,
);

class ChannelCard extends StatelessWidget {
  const ChannelCard({
    super.key,
    required this.integration,
    required this.stats,
    required this.onAction,
    this.compact = false,
  });

  final MyIntegration integration;
  final IntegrationStats stats;
  final ChannelActionHandler onAction;
  final bool compact;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = integration.integrationTypeName ?? integration.name;

    return AppCard(
      padding: const EdgeInsets.all(13),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(
                name: name,
                imageUrl: integration.imageUrl,
                size: 38,
                radius: 9,
                padding: 5,
              ),
              const SizedBox(width: 11),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      name,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      integration.name,
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              AppStatusChip(
                dense: true,
                label: integration.isActive ? 'Activa' : 'Inactiva',
                tone: integration.isActive
                    ? AppStatusTone.success
                    : AppStatusTone.neutral,
              ),
              const SizedBox(width: 2),
              PopupMenuButton<ChannelAction>(
                icon: const Icon(Icons.tune_rounded, size: 19),
                tooltip: 'Acciones',
                onSelected: (action) => onAction(integration, action),
                itemBuilder: (context) => [
                  PopupMenuItem(
                    value: ChannelAction.toggle,
                    child: Row(
                      children: [
                        Icon(
                          integration.isActive
                              ? Icons.toggle_on_rounded
                              : Icons.toggle_off_rounded,
                          size: 19,
                        ),
                        const SizedBox(width: 10),
                        Text(integration.isActive ? 'Desactivar' : 'Activar'),
                      ],
                    ),
                  ),
                  const PopupMenuItem(
                    value: ChannelAction.test,
                    child: Row(
                      children: [
                        Icon(Icons.wifi_tethering_rounded, size: 19),
                        SizedBox(width: 10),
                        Text('Probar conexion'),
                      ],
                    ),
                  ),
                  const PopupMenuItem(
                    value: ChannelAction.sync,
                    child: Row(
                      children: [
                        Icon(Icons.sync_rounded, size: 19),
                        SizedBox(width: 10),
                        Text('Sincronizar ordenes'),
                      ],
                    ),
                  ),
                ],
              ),
            ],
          ),
          if (!compact) ...[
            const SizedBox(height: 12),
            Row(
              children: [
                _MiniStat(
                  value: AppFormat.number(stats.ordersCount),
                  label: 'ordenes',
                ),
                const SizedBox(width: 20),
                _MiniStat(
                  value: AppFormat.number(stats.productsCount),
                  label: 'productos',
                ),
                const Spacer(),
                if (stats.lastOrderAt != null)
                  Text(
                    AppFormat.relative(
                      AppFormat.parseDate(stats.lastOrderAt!),
                    ),
                    style: theme.textTheme.labelSmall,
                  ),
              ],
            ),
            ChannelSyncStrip(
              integrationId: integration.id,
              integrationTypeId: integration.integrationTypeId,
            ),
            if (stats.hasOrders) ...[
              const SizedBox(height: 10),
              _ProgressBar(stats: stats),
              const SizedBox(height: 9),
              Wrap(
                spacing: 6,
                runSpacing: 6,
                children: [
                  if (stats.ordersInProgress > 0)
                    _StatDot(
                      count: stats.ordersInProgress,
                      label: 'en curso',
                      color: AppColors.primary,
                    ),
                  if (stats.ordersDelivered > 0)
                    _StatDot(
                      count: stats.ordersDelivered,
                      label: 'entregadas',
                      color: AppColors.success,
                    ),
                  if (stats.ordersCancelled > 0)
                    _StatDot(
                      count: stats.ordersCancelled,
                      label: 'canceladas',
                      color: AppColors.error,
                    ),
                  if (stats.ordersReturned > 0)
                    _StatDot(
                      count: stats.ordersReturned,
                      label: 'devueltas',
                      color: AppColors.warning,
                    ),
                ],
              ),
            ],
          ],
        ],
      ),
    );
  }
}

class _MiniStat extends StatelessWidget {
  const _MiniStat({required this.value, required this.label});

  final String value;
  final String label;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Row(
      crossAxisAlignment: CrossAxisAlignment.baseline,
      textBaseline: TextBaseline.alphabetic,
      children: [
        Text(
          value,
          style: theme.textTheme.titleSmall?.copyWith(fontSize: 15),
        ),
        const SizedBox(width: 4),
        Text(label, style: theme.textTheme.labelSmall),
      ],
    );
  }
}

class _ProgressBar extends StatelessWidget {
  const _ProgressBar({required this.stats});

  final IntegrationStats stats;

  @override
  Widget build(BuildContext context) {
    final total = stats.ordersCount == 0 ? 1 : stats.ordersCount;
    final parts = <({int count, Color color})>[
      (count: stats.ordersInProgress, color: AppColors.primary),
      (count: stats.ordersDelivered, color: AppColors.success),
      (count: stats.ordersCancelled, color: AppColors.error),
      (count: stats.ordersReturned, color: AppColors.warning),
    ].where((p) => p.count > 0).toList();

    return ClipRRect(
      borderRadius: BorderRadius.circular(999),
      child: SizedBox(
        height: 5,
        child: Row(
          children: [
            for (final part in parts)
              Expanded(
                flex: (part.count * 1000 ~/ total).clamp(1, 1000),
                child: Container(color: part.color),
              ),
            if (parts.isEmpty)
              Expanded(child: Container(color: AppColors.surfaceMuted)),
          ],
        ),
      ),
    );
  }
}
