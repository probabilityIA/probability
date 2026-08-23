import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/utils/integration_visibility.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../../../integrations/core/ui/providers/integration_provider.dart';
import '../providers/my_integrations_provider.dart';
import '../providers/sync_activity_provider.dart';
import '../../domain/sync_entities.dart';
import '../widgets/environment_panel.dart';
import '../widgets/channel_cards.dart';
import '../widgets/module_toolbar.dart';

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

  Future<void> _load() async {
    final provider = context.read<MyIntegrationsProvider>();
    final sync = context.read<SyncActivityProvider>();
    await provider.fetchIntegrations(businessId: widget.businessId);
    if (!mounted) return;
    provider.fetchStats(businessId: widget.businessId);
    await sync.configure(
      integrations: provider.integrations,
      businessId: widget.businessId,
    );
  }

  Future<void> _runAction(MyIntegration item, ChannelAction action) async {
    final integrations = context.read<IntegrationProvider>();
    final messenger = ScaffoldMessenger.of(context);

    if (action == ChannelAction.toggle) {
      final ok = await showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: Text(item.isActive ? 'Desactivar integracion' : 'Activar integracion'),
          content: Text(
            item.isActive
                ? 'Dejara de sincronizar ordenes y productos hasta que la vuelvas a activar.'
                : 'Volvera a sincronizar ordenes y productos.',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: const Text('Cancelar'),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: Text(item.isActive ? 'Desactivar' : 'Activar'),
            ),
          ],
        ),
      );
      if (ok != true) return;
    }

    messenger.showSnackBar(
      SnackBar(content: Text('${action.runningLabel}...')),
    );

    var done = false;
    switch (action) {
      case ChannelAction.toggle:
        done = item.isActive
            ? await integrations.deactivateIntegration(item.id)
            : await integrations.activateIntegration(item.id);
      case ChannelAction.test:
        done = await integrations.testConnection(item.id) != null;
      case ChannelAction.sync:
        done = await integrations.syncOrders(item.id) != null;
    }

    if (!mounted) return;
    messenger.hideCurrentSnackBar();
    messenger.showSnackBar(
      SnackBar(
        content: Text(done ? action.doneLabel : 'No se pudo completar la accion'),
      ),
    );
    if (done) _load();
  }

  void _selectEnvironment(SyncEnvironment environment) {
    context.read<SyncActivityProvider>().setEnvironment(environment);
  }

  @override
  Widget build(BuildContext context) {
    final sync = context.watch<SyncActivityProvider>();

    return AppScaffold(
      bottom: ModuleToolbar(
        environment: sync.environment,
        running: sync.running,
        onEnvironment: _selectEnvironment,
      ),
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

          return EnvironmentPanel(
            environment: sync.environment,
            integrations: visible,
            statsFor: provider.statsFor,
            businessId: widget.businessId,
            onAction: _runAction,
            onRefresh: _load,
          );
        },
      ),
    );
  }
}

