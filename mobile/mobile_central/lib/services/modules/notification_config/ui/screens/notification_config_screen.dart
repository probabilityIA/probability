import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../providers/notification_config_provider.dart';

const Map<String, String> notificationCategoryLabels = {
  'order': 'Ordenes',
  'shipment': 'Envios',
  'wallet': 'Billetera',
  'invoice': 'Facturacion',
};

class NotificationConfigScreen extends StatefulWidget {
  const NotificationConfigScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<NotificationConfigScreen> createState() => _NotificationConfigScreenState();
}

class _NotificationConfigScreenState extends State<NotificationConfigScreen> {
  int? _saving;
  String _channel = '';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh());
  }

  @override
  void didUpdateWidget(NotificationConfigScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _refresh();
  }

  void _refresh() {
    context.read<NotificationConfigProvider>().fetchConfigs(
          filter: ConfigFilter(businessId: widget.businessId),
        );
  }

  Future<void> _toggle(NotificationConfig config, bool value) async {
    setState(() => _saving = config.id);

    final provider = context.read<NotificationConfigProvider>();
    final result = await provider.updateConfig(
      config.id,
      UpdateConfigDTO(enabled: value),
      businessId: widget.businessId,
    );

    if (!mounted) return;
    setState(() => _saving = null);

    if (result == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(provider.error ?? 'No se pudo guardar')),
      );
    }
    _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<NotificationConfigProvider>(
      builder: (context, provider, _) {
        if (provider.isLoading && provider.configs.isEmpty) {
          return const AppListSkeleton();
        }
        if (provider.error != null && provider.configs.isEmpty) {
          return AppErrorState(message: provider.error!, onRetry: _refresh);
        }
        if (provider.configs.isEmpty) {
          return AppEmptyState(
            icon: Icons.notifications_none_rounded,
            title: 'Sin notificaciones configuradas',
            message: 'Conecta un canal como WhatsApp para avisarle a tus clientes.',
            actionLabel: 'Actualizar',
            onAction: _refresh,
          );
        }

        final channels = <String>{
          ...provider.configs
              .map((c) => c.notificationTypeName ?? '')
              .where((c) => c.isNotEmpty),
        }.toList()
          ..sort();

        final rows = _channel.isEmpty
            ? provider.configs
            : provider.configs
                .where((c) => (c.notificationTypeName ?? '') == _channel)
                .toList();

        final active = provider.configs.where((c) => c.enabled).length;

        return RefreshIndicator(
          onRefresh: () async => _refresh(),
          color: AppColors.primary,
          child: ListView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: AppSpacing.page,
            children: [
              AppCard(
                child: Row(
                  children: [
                    Container(
                      width: 44,
                      height: 44,
                      alignment: Alignment.center,
                      decoration: BoxDecoration(
                        color: AppColors.primarySoft,
                        borderRadius: AppRadius.mdAll,
                      ),
                      child: const Icon(
                        Icons.notifications_active_outlined,
                        size: 21,
                        color: AppColors.primary,
                      ),
                    ),
                    const SizedBox(width: 13),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            '$active de ${provider.configs.length} eventos activos',
                            style: Theme.of(context).textTheme.titleSmall,
                          ),
                          const SizedBox(height: 2),
                          Text(
                            'Cada evento avisa al cliente por su canal',
                            style: Theme.of(context).textTheme.labelSmall,
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
              if (channels.length > 1) ...[
                AppFilterChips(
                  options: [
                    (value: '', label: 'Todos'),
                    ...channels.map((c) => (value: c, label: c)),
                  ],
                  selected: _channel,
                  onSelected: (value) => setState(() => _channel = value),
                ),
                const SizedBox(height: 14),
              ],
              AppCard(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Column(
                  children: [
                    for (var i = 0; i < rows.length; i++) ...[
                      if (i > 0) const Divider(height: 1),
                      _ConfigTile(
                        config: rows[i],
                        saving: _saving == rows[i].id,
                        onChanged: (value) => _toggle(rows[i], value),
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 26),
            ],
          ),
        );
      },
    );
  }
}

class _ConfigTile extends StatelessWidget {
  const _ConfigTile({
    required this.config,
    required this.saving,
    required this.onChanged,
  });

  final NotificationConfig config;
  final bool saving;
  final ValueChanged<bool> onChanged;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final channel = config.notificationTypeName ?? 'Canal';
    final category = config.eventType ?? '';

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 11),
      child: Row(
        children: [
          BrandLogo(name: channel, size: 36, radius: 10, padding: 5),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  config.notificationEventName ?? 'Evento',
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 3),
                Row(
                  children: [
                    Text(channel, style: theme.textTheme.labelSmall),
                    if (category.isNotEmpty) ...[
                      const SizedBox(width: 6),
                      AppStatusChip(
                        dense: true,
                        label: notificationCategoryLabels[category] ?? category,
                        tone: AppStatusTone.neutral,
                      ),
                    ],
                  ],
                ),
                const SizedBox(height: 3),
                Text(
                  'Actualizado ${AppFormat.relative(AppFormat.parseDate(config.updatedAt))}',
                  style: theme.textTheme.labelSmall,
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          if (saving)
            const SizedBox(
              height: 20,
              width: 20,
              child: CircularProgressIndicator(strokeWidth: 2),
            )
          else
            Switch(value: config.enabled, onChanged: onChanged),
        ],
      ),
    );
  }
}
