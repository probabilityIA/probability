import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/integration_visibility.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/my_integrations_provider.dart';

class CoreScreen extends StatefulWidget {
  const CoreScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<CoreScreen> createState() => _CoreScreenState();
}

class _CoreScreenState extends State<CoreScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void didUpdateWidget(CoreScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _load();
  }

  void _load() {
    final provider = context.read<MyIntegrationsProvider>();
    provider.fetchIntegrations(businessId: widget.businessId);
    provider.fetchStats(businessId: widget.businessId);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Tus integraciones',
      subtitle: 'El nucleo de tu operacion',
      onBack: () => Navigator.of(context).maybePop(),
      body: Consumer<MyIntegrationsProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.integrations.isEmpty) {
            return const AppListSkeleton();
          }
          if (provider.error != null && provider.integrations.isEmpty) {
            return AppErrorState(message: provider.error!, onRetry: _load);
          }

          final visible = provider.integrations
              .where((i) => IntegrationVisibility.isVisible(
                    category: i.categoryCode,
                    type: i.integrationTypeCode,
                    name: i.integrationTypeName,
                  ))
              .toList();

          if (visible.isEmpty) {
            return AppEmptyState(
              icon: Icons.hub_outlined,
              title: 'Sin integraciones',
              message:
                  'Cuando conectes una tienda o un facturador lo vas a ver aqui.',
              actionLabel: 'Actualizar',
              onAction: _load,
            );
          }

          final sales = visible
              .where((i) => (i.categoryCode ?? '') == 'ecommerce')
              .toList();
          final others =
              visible.where((i) => !sales.contains(i)).toList();

          var totals = const IntegrationStats(integrationId: 0);
          for (final item in visible) {
            totals = totals + provider.statsFor(item.id);
          }

          return RefreshIndicator(
            onRefresh: () async => _load(),
            color: Theme.of(context).colorScheme.primary,
            child: ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.fromLTRB(16, 4, 16, 96),
              children: [
                _CoreCard(totals: totals),
                const SizedBox(height: 20),
                if (sales.isNotEmpty) ...[
                  const AppSectionHeader(title: 'Canales de venta'),
                  for (final item in sales) ...[
                    _ChannelCard(
                      integration: item,
                      stats: provider.statsFor(item.id),
                    ),
                    const SizedBox(height: 10),
                  ],
                  const SizedBox(height: 10),
                ],
                if (others.isNotEmpty) ...[
                  const AppSectionHeader(title: 'Otras integraciones'),
                  for (final item in others) ...[
                    _ChannelCard(
                      integration: item,
                      stats: provider.statsFor(item.id),
                      compact: true,
                    ),
                    const SizedBox(height: 10),
                  ],
                ],
              ],
            ),
          );
        },
      ),
    );
  }
}

class _CoreCard extends StatelessWidget {
  const _CoreCard({required this.totals});

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

class _ChannelCard extends StatelessWidget {
  const _ChannelCard({
    required this.integration,
    required this.stats,
    this.compact = false,
  });

  final MyIntegration integration;
  final IntegrationStats stats;
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
