import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/navigation/app_modules.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../../auth/business/ui/providers/business_provider.dart';
import '../../../../integrations/core/ui/providers/integration_provider.dart';
import '../../../../auth/login/ui/providers/login_provider.dart';
import '../../domain/entities.dart';
import '../providers/dashboard_provider.dart';
import '../widgets/dashboard_sections.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _load();
      final integrations = context.read<IntegrationProvider>();
      if (integrations.integrationTypes.isEmpty) {
        integrations.fetchIntegrationTypes();
      }
    });
  }

  @override
  void didUpdateWidget(DashboardScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _load();
  }

  void _load() {
    context.read<DashboardProvider>().fetchStats(businessId: widget.businessId);
  }

  Map<String, String?> _channelLogos(BuildContext context) {
    final types = context.watch<IntegrationProvider>().integrationTypes;
    final map = <String, String?>{};
    for (final type in types) {
      map[type.name.toLowerCase()] = type.imageUrl;
      map[type.code.toLowerCase()] = type.imageUrl;
    }
    return map;
  }

  String? _activeBusinessName(BuildContext context, LoginProvider login) {
    if (!login.isSuperAdmin) return login.businessName;
    final id = widget.businessId;
    if (id == null) return null;
    final match = context
        .watch<BusinessProvider>()
        .businessesSimple
        .where((b) => b.id == id)
        .firstOrNull;
    return match?.name;
  }

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();

    return AppScaffold(
      showLogo: true,
      actions: [
        IconButton(
          icon: const Icon(Icons.notifications_none_rounded, size: 22),
          tooltip: 'Notificaciones',
          onPressed: () => context.go('/notifications'),
        ),
        IconButton(
          icon: const Icon(Icons.person_outline_rounded, size: 22),
          tooltip: 'Perfil',
          onPressed: () => context.go('/profile'),
        ),
      ],
      body: Consumer<DashboardProvider>(
        builder: (context, provider, _) {
          if (provider.isLoading && provider.stats == null) {
            return const AppListSkeleton();
          }
          if (provider.error != null && provider.stats == null) {
            return AppErrorState(message: provider.error!, onRetry: _load);
          }

          final stats = provider.stats;
          if (stats == null) {
            return AppEmptyState(
              icon: Icons.insights_outlined,
              title: 'Sin datos por ahora',
              message: 'Cuando entren ordenes vas a ver aqui el resumen de tu operacion.',
              actionLabel: 'Actualizar',
              onAction: _load,
            );
          }

          return RefreshIndicator(
            onRefresh: () async => _load(),
            color: Theme.of(context).colorScheme.primary,
            child: ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: AppSpacing.page,
              children: [
                _Greeting(
                  name: login.user?.name,
                  business: _activeBusinessName(context, login),
                ),
                const SizedBox(height: 12),
                _PeriodPicker(
                  selected: provider.period,
                  onSelected: (period) {
                    provider.setPeriod(period);
                    _load();
                  },
                ),
                const SizedBox(height: 12),
                _KpiGrid(
                  stats: stats,
                  channelLogos: _channelLogos(context),
                  effectiveness: provider.effectiveness,
                ),
                const SizedBox(height: 16),
                const _QuickActions(),
                if (stats.ordersByIntegrationType.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(
                    title: 'Ordenes por canal',
                    subtitle: 'De donde vienen tus ventas',
                  ),
                  DashboardChannelsSection(items: stats.ordersByIntegrationType),
                ],
                if (stats.shipmentsByStatus.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(title: 'Envios por estado'),
                  DashboardStatusSection(items: stats.shipmentsByStatus),
                ],
                if (stats.shipmentsByCarrier.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(
                    title: 'Envios por transportadora',
                    subtitle: 'Guias generadas por operador',
                  ),
                  DashboardCarriersSection(items: stats.shipmentsByCarrier),
                ],
                if (stats.topProducts.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(title: 'Productos mas vendidos'),
                  DashboardProductsSection(items: stats.topProducts),
                ],
                if (stats.topCustomers.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(title: 'Mejores clientes'),
                  DashboardCustomersSection(items: stats.topCustomers),
                ],
                if (stats.ordersByLocation.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  const AppSectionHeader(title: 'Ordenes por ciudad'),
                  DashboardLocationsSection(items: stats.ordersByLocation),
                ],
                const SizedBox(height: 24),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _Greeting extends StatelessWidget {
  const _Greeting({this.name, this.business});

  final String? name;
  final String? business;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final firstName = (name ?? '').split(' ').first;
    final hour = DateTime.now().hour;
    final salute = hour < 12
        ? 'Buenos dias'
        : hour < 19
            ? 'Buenas tardes'
            : 'Buenas noches';

    return Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Expanded(
          child: Text(
            firstName.isEmpty ? salute : '$salute, $firstName',
            style: theme.textTheme.titleLarge?.copyWith(fontSize: 18),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        if (business != null) ...[
          const SizedBox(width: 10),
          Flexible(
            child: Text(
              business!,
              style: theme.textTheme.labelSmall,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              textAlign: TextAlign.right,
            ),
          ),
        ],
      ],
    );
  }
}

class _KpiGrid extends StatelessWidget {
  const _KpiGrid({
    required this.stats,
    this.channelLogos = const {},
    this.effectiveness,
  });

  final DashboardStats stats;
  final Map<String, String?> channelLogos;
  final OrderEffectiveness? effectiveness;

  double get _revenue =>
      stats.topProducts.fold<double>(0, (acc, p) => acc + p.totalSold);

  int get _unitsSold =>
      stats.topProducts.fold<int>(0, (acc, p) => acc + p.orderCount);

  int get _shipments =>
      stats.shipmentsByStatus.fold<int>(0, (acc, s) => acc + s.count);

  int get _delivered => stats.shipmentsByStatus
      .where((s) => s.status.toLowerCase().contains('deliver'))
      .fold<int>(0, (acc, s) => acc + s.count);

  @override
  Widget build(BuildContext context) {
    final tiles = [
      AppKpiTile(
        label: 'Ordenes totales',
        value: AppFormat.number(stats.totalOrders),
        icon: Icons.receipt_long_outlined,
        onTap: () => context.go('/orders'),
        footer: effectiveness == null || !effectiveness!.hasData
            ? null
            : _EffectivenessRows(data: effectiveness!),
      ),
      _ChannelsKpiTile(
        channels: stats.ordersByIntegrationType,
        logos: channelLogos,
      ),
      AppKpiTile(
        label: 'Vendido',
        value: AppFormat.money(_revenue),
        trend: _unitsSold == 0 ? null : '${AppFormat.number(_unitsSold)} unidades',
        icon: Icons.payments_outlined,
        accent: AppColors.success,
        onTap: () => context.go('/inventory'),
      ),
      AppKpiTile(
        label: 'Guias generadas',
        value: AppFormat.number(_shipments),
        icon: Icons.local_shipping_outlined,
        accent: const Color(0xFFF97316),
        trend: _shipments == 0 ? null : '$_delivered entregadas',
        onTap: () => context.go('/orders/shipments'),
      ),
    ];

    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      padding: EdgeInsets.zero,
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        mainAxisSpacing: 9,
        crossAxisSpacing: 9,
        mainAxisExtent: 124,
      ),
      itemCount: tiles.length,
      itemBuilder: (context, index) => tiles[index],
    );
  }
}

class _QuickActions extends StatelessWidget {
  const _QuickActions();

  static const List<AppModule> _modules = [
    AppModule(label: 'Ordenes', route: '/orders', icon: Icons.receipt_long_outlined),
    AppModule(label: 'Envios', route: '/orders/shipments', icon: Icons.local_shipping_outlined),
    AppModule(label: 'Clientes', route: '/customers', icon: Icons.people_alt_outlined),
    AppModule(label: 'Billetera', route: '/wallet', icon: Icons.account_balance_wallet_outlined),
  ];

  @override
  Widget build(BuildContext context) {
    return Row(
      children: _modules
          .map((module) => Expanded(
                child: Padding(
                  padding: EdgeInsets.only(right: module == _modules.last ? 0 : 10),
                  child: _QuickAction(module: module),
                ),
              ))
          .toList(),
    );
  }
}

class _QuickAction extends StatelessWidget {
  const _QuickAction({required this.module});

  final AppModule module;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: AppColors.surface,
      borderRadius: AppRadius.lgAll,
      child: InkWell(
        borderRadius: AppRadius.lgAll,
        onTap: () => context.go(module.route),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 6),
          decoration: BoxDecoration(
            borderRadius: AppRadius.lgAll,
            border: Border.all(color: AppColors.border),
          ),
          child: Column(
            children: [
              Icon(module.icon, size: 19,
                  color: Theme.of(context).colorScheme.primary),
              const SizedBox(height: 6),
              Text(
                module.label,
                style: Theme.of(context).textTheme.labelSmall?.copyWith(
                      color: AppColors.textSecondary,
                      fontWeight: FontWeight.w600,
                    ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PeriodPicker extends StatelessWidget {
  const _PeriodPicker({required this.selected, required this.onSelected});

  final DashboardPeriod selected;
  final ValueChanged<DashboardPeriod> onSelected;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;

    return SizedBox(
      height: 32,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: EdgeInsets.zero,
        itemCount: DashboardPeriod.values.length,
        separatorBuilder: (context, index) => const SizedBox(width: 7),
        itemBuilder: (context, index) {
          final period = DashboardPeriod.values[index];
          final on = period == selected;
          return Material(
            color: on ? scheme.primaryContainer : AppColors.surface,
            borderRadius: AppRadius.pillAll,
            child: InkWell(
              borderRadius: AppRadius.pillAll,
              onTap: on ? null : () => onSelected(period),
              child: Container(
                alignment: Alignment.center,
                padding: const EdgeInsets.symmetric(horizontal: 14),
                decoration: BoxDecoration(
                  borderRadius: AppRadius.pillAll,
                  border: Border.all(
                    color: on ? scheme.primary : AppColors.border,
                  ),
                ),
                child: Text(
                  period.label,
                  style: TextStyle(
                    fontFamily: 'Inter',
                    fontSize: 12.5,
                    fontWeight: FontWeight.w600,
                    color: on
                        ? scheme.onPrimaryContainer
                        : AppColors.textSecondary,
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}

class _ChannelsKpiTile extends StatelessWidget {
  const _ChannelsKpiTile({required this.channels, required this.logos});

  final List<OrderCountByIntegrationType> channels;
  final Map<String, String?> logos;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    const accent = Color(0xFF0EA5E9);
    final sorted = [...channels]..sort((a, b) => b.count.compareTo(a.count));
    final shown = sorted.take(4).toList();
    final extra = sorted.length - shown.length;

    return Material(
      color: AppColors.surface,
      borderRadius: AppRadius.lgAll,
      child: InkWell(
        borderRadius: AppRadius.lgAll,
        onTap: () => context.go('/integrations'),
        child: Container(
          padding: const EdgeInsets.all(11),
          decoration: BoxDecoration(
            borderRadius: AppRadius.lgAll,
            border: Border.all(color: AppColors.border),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                children: [
                  Container(
                    width: 26,
                    height: 26,
                    decoration: BoxDecoration(
                      color: accent.withValues(alpha: 0.10),
                      borderRadius: AppRadius.smAll,
                    ),
                    child: const Icon(Icons.hub_outlined, size: 15, color: accent),
                  ),
                  const SizedBox(width: 9),
                  Expanded(
                    child: Text(
                      'Canales activos',
                      style: theme.textTheme.labelMedium,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  Text(
                    AppFormat.number(sorted.length),
                    style: theme.textTheme.headlineSmall,
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Wrap(
                      spacing: 4,
                      runSpacing: 4,
                      crossAxisAlignment: WrapCrossAlignment.center,
                      children: [
                        for (final channel in shown)
                          BrandLogo(
                            name: channel.integrationType,
                            imageUrl:
                                logos[channel.integrationType.toLowerCase()],
                            size: 22,
                            radius: 6,
                            padding: 3,
                          ),
                        if (extra > 0)
                          Text('+$extra', style: theme.textTheme.labelSmall),
                      ],
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _EffectivenessRows extends StatelessWidget {
  const _EffectivenessRows({required this.data});

  final OrderEffectiveness data;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        _MetricPill(
          label: 'EFECTIVIDAD',
          rate: data.effectivenessRate,
          count: data.delivered,
          caption: 'entregadas',
          color: AppColors.success,
        ),
        const SizedBox(height: 4),
        _MetricPill(
          label: 'DEVOLUCIONES',
          rate: data.returnRate,
          count: data.returned,
          caption: 'devueltas',
          color: AppColors.warning,
        ),
      ],
    );
  }
}

class _MetricPill extends StatelessWidget {
  const _MetricPill({
    required this.label,
    required this.rate,
    required this.count,
    required this.caption,
    required this.color,
  });

  final String label;
  final double rate;
  final int count;
  final String caption;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      children: [
        Container(
          width: 6,
          height: 6,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(width: 5),
        Text(
          '${(rate * 100).toStringAsFixed(1)}%',
          style: theme.textTheme.labelMedium?.copyWith(
            fontWeight: FontWeight.w700,
            color: color,
          ),
        ),
        const SizedBox(width: 5),
        Expanded(
          child: Text(
            '${AppFormat.number(count)} $caption',
            style: theme.textTheme.labelSmall?.copyWith(fontSize: 10),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}
