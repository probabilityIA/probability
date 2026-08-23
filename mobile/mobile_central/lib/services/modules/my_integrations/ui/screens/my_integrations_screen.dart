import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../../../../shared/utils/integration_visibility.dart';
import '../providers/my_integrations_provider.dart';
import '../widgets/integration_actions_sheet.dart';

class MyIntegrationsScreen extends StatefulWidget {
  const MyIntegrationsScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<MyIntegrationsScreen> createState() => _MyIntegrationsScreenState();
}

class _MyIntegrationsScreenState extends State<MyIntegrationsScreen> {
  String _category = 'all';

  static const Map<String, String> _categoryLabels = {
    'all': 'Todas',
    'ecommerce': 'Tiendas',
    'shipping': 'Transporte',
    'invoicing': 'Facturacion',
    'messaging': 'Mensajeria',
    'payment': 'Pagos',
    'platform': 'Plataforma',
  };

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  void _load() {
    context.read<MyIntegrationsProvider>().fetchIntegrations(businessId: widget.businessId);
  }

  List<MyIntegration> _visible(List<MyIntegration> rows) => rows
      .where((r) => IntegrationVisibility.isVisible(
            category: r.categoryCode,
            type: r.integrationTypeCode,
            name: r.integrationTypeName,
          ))
      .toList();

  List<MyIntegration> _filter(List<MyIntegration> rows) {
    final visible = _visible(rows);
    if (_category == 'all') return visible;
    return visible.where((r) => (r.categoryCode ?? '') == _category).toList();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<MyIntegrationsProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.integrations.isEmpty) {
          return const AppListSkeleton();
        }

        if (provider.error != null && provider.integrations.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _load);
        }

        if (_visible(provider.integrations).isEmpty) {
          return AppEmptyState(
            icon: Icons.hub_outlined,
            title: 'Sin integraciones conectadas',
            message: 'Conecta tu primera tienda, transportadora o facturador desde el catalogo.',
            actionLabel: 'Recargar',
            onAction: _load,
          );
        }

        final rows = _filter(provider.integrations);
        final categories = <String>{
          'all',
          ..._visible(provider.integrations)
              .map((e) => e.categoryCode ?? '')
              .where((e) => e.isNotEmpty),
        }.toList();

        return RefreshIndicator(
          onRefresh: () async => _load(),
          color: AppColors.primary,
          child: Column(
            children: [
              const SizedBox(height: 12),
              AppFilterChips(
                selected: _category,
                onSelected: (value) => setState(() => _category = value),
                options: categories
                    .map((code) => (value: code, label: _categoryLabels[code] ?? code))
                    .toList(),
              ),
              const SizedBox(height: 4),
              Expanded(
                child: ListView.separated(
                  padding: AppSpacing.page,
                  itemCount: rows.length,
                  separatorBuilder: (context, index) => const SizedBox(height: 10),
                  itemBuilder: (context, index) => _IntegrationCard(
                    integration: rows[index],
                    onTap: () async {
                      await showIntegrationActionsSheet(context, integration: rows[index]);
                      if (context.mounted) _load();
                    },
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _IntegrationCard extends StatelessWidget {
  const _IntegrationCard({required this.integration, this.onTap});

  final MyIntegration integration;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final active = integration.isActive;

    return AppCard(
      padding: const EdgeInsets.all(14),
      onTap: onTap,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          BrandLogo(
            name: integration.integrationTypeName ?? integration.name,
            imageUrl: integration.imageUrl,
            size: 46,
            radius: 12,
            padding: 7,
          ),
          const SizedBox(width: 13),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        integration.name,
                        style: theme.textTheme.titleSmall,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    const SizedBox(width: 8),
                    AppStatusChip(
                      dense: true,
                      label: active ? 'Conectada' : 'Inactiva',
                      tone: active ? AppStatusTone.success : AppStatusTone.neutral,
                    ),
                  ],
                ),
                const SizedBox(height: 3),
                Text(
                  integration.integrationTypeName ?? 'Integracion',
                  style: theme.textTheme.bodySmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 9),
                Row(
                  children: [
                    const Icon(Icons.schedule, size: 13, color: AppColors.textDisabled),
                    const SizedBox(width: 5),
                    Text(
                      'Actualizada ${AppFormat.relative(AppFormat.parseDate(integration.updatedAt))}',
                      style: theme.textTheme.labelSmall,
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
