import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../app/sync_providers.dart';
import '../../domain/sync_entities.dart';
import '../providers/sync_activity_provider.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/saved_comparison_provider.dart';
import '../screens/finding_items_screen.dart';
import '../screens/inventory_compare_screen.dart';
import 'channel_cards.dart';
import 'channel_data_sheet.dart';
import 'inventory_environment.dart';
import 'module_toolbar.dart';
import 'products_environment.dart';
import 'orders_compare_panel.dart';
import 'orders_report.dart';
import 'saved_comparison_views.dart';

class EnvironmentPanel extends StatelessWidget {
  const EnvironmentPanel({
    super.key,
    required this.environment,
    required this.integrations,
    required this.statsFor,
    required this.onAction,
    required this.onRefresh,
    this.businessId,
  });

  final SyncEnvironment environment;
  final List<MyIntegration> integrations;
  final IntegrationStats Function(int) statsFor;
  final ChannelActionHandler onAction;
  final Future<void> Function() onRefresh;
  final int? businessId;

  @override
  Widget build(BuildContext context) {
    if (environment == SyncEnvironment.overview) {
      return _Overview(
        integrations: integrations,
        statsFor: statsFor,
        onAction: onAction,
        onRefresh: onRefresh,
      );
    }

    if (environment == SyncEnvironment.products) {
      return ProductsEnvironment(businessId: businessId);
    }

    if (environment == SyncEnvironment.inventory) {
      return InventoryEnvironment(
        integrations: integrations,
        businessId: businessId,
      );
    }

    if (environment == SyncEnvironment.ordersCompare) {
      return OrdersComparePanel(
        integrations: integrations,
        businessId: businessId,
      );
    }

    final spec = environmentSpec(environment);

    if (environment == SyncEnvironment.invoicing) {
      return _Notice(
        icon: spec.icon,
        title: 'Facturar',
        message: 'Facturacion desde el hub: proximamente.',
      );
    }

    return _SavedEnvironment(
      spec: spec,
      integrations: integrations,
      businessId: businessId,
    );
  }
}

class _Overview extends StatelessWidget {
  const _Overview({
    required this.integrations,
    required this.statsFor,
    required this.onAction,
    required this.onRefresh,
  });

  final List<MyIntegration> integrations;
  final IntegrationStats Function(int) statsFor;
  final ChannelActionHandler onAction;
  final Future<void> Function() onRefresh;

  @override
  Widget build(BuildContext context) {
    final sales = integrations
        .where((i) => (i.categoryCode ?? '') == 'ecommerce')
        .toList();
    final others = integrations.where((i) => !sales.contains(i)).toList();

    var totals = const IntegrationStats(integrationId: 0);
    for (final item in integrations) {
      totals = totals + statsFor(item.id);
    }

    return RefreshIndicator(
      onRefresh: onRefresh,
      color: Theme.of(context).colorScheme.primary,
      child: ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.fromLTRB(16, 12, 16, 96),
        children: [
          CoreCard(totals: totals),
          const SizedBox(height: 16),
          OrdersReport(
            integrations: integrations,
            statsFor: statsFor,
            embedded: true,
          ),
          const SizedBox(height: 16),
          if (sales.isNotEmpty) ...[
            const AppSectionHeader(title: 'Canales de venta'),
            for (final item in sales) ...[
              ChannelCard(
                integration: item,
                stats: statsFor(item.id),
                onAction: onAction,
              ),
              const SizedBox(height: 10),
            ],
            const SizedBox(height: 6),
          ],
          if (others.isNotEmpty) ...[
            const AppSectionHeader(title: 'Otras integraciones'),
            for (final item in others) ...[
              ChannelCard(
                integration: item,
                stats: statsFor(item.id),
                compact: true,
                onAction: onAction,
              ),
              const SizedBox(height: 10),
            ],
          ],
        ],
      ),
    );
  }
}

class _SavedEnvironment extends StatefulWidget {
  const _SavedEnvironment({
    required this.spec,
    required this.integrations,
    this.businessId,
  });

  final EnvironmentSpec spec;
  final List<MyIntegration> integrations;
  final int? businessId;

  @override
  State<_SavedEnvironment> createState() => _SavedEnvironmentState();
}

class _SavedEnvironmentState extends State<_SavedEnvironment> {
  DataApplyResult? _lastApply;
  bool _undoing = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void didUpdateWidget(_SavedEnvironment oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId ||
        oldWidget.spec.environment != widget.spec.environment) {
      _load();
    }
  }

  Future<void> _load({bool force = false}) async {
    if (!mounted) return;
    final saved = context.read<SavedComparisonProvider>();
    saved.configure(businessId: widget.businessId);

    switch (widget.spec.environment) {
      case SyncEnvironment.products:
        await saved.loadFindings(force: force);
      case SyncEnvironment.data:
        await saved.loadDataSummary(force: force);
      case SyncEnvironment.inventory:
        await saved.loadInventorySnapshots(widget.integrations, force: force);
      default:
        return;
    }
  }

  void _openFinding(Finding finding) {
    Navigator.of(context).push(MaterialPageRoute(
      builder: (_) => FindingItemsScreen(
        finding: finding,
        businessId: widget.businessId,
      ),
    ));
  }

  String? _logoFor(int integrationId) {
    for (final item in widget.integrations) {
      if (item.id == integrationId) return item.imageUrl;
    }
    return null;
  }

  Future<void> _pickData(
    DataSummaryRow row,
    DataSummaryCell cell,
    DataMode mode,
  ) async {
    final messenger = ScaffoldMessenger.of(context);
    final result = await showChannelDataSheet(
      context,
      target: ChannelDataTarget(
        field: row.field,
        label: row.label.isEmpty ? row.field : row.label,
        cell: cell,
        mode: mode,
        logoUrl: _logoFor(cell.integrationId),
      ),
      businessId: widget.businessId,
    );
    if (result == null || !mounted) return;

    setState(() => _lastApply = result);
    messenger.showSnackBar(SnackBar(
      content: Text('${result.applied} productos actualizados'),
    ));
    await _load(force: true);
  }

  Future<void> _undo() async {
    final batch = _lastApply?.batchId;
    if (batch == null || batch.isEmpty) return;

    final saved = context.read<SavedComparisonProvider>();
    final messenger = ScaffoldMessenger.of(context);

    setState(() => _undoing = true);
    try {
      final reverted = await saved.undoChannelData(batch);
      if (!mounted) return;
      setState(() => _lastApply = null);
      messenger.showSnackBar(SnackBar(
        content: Text('$reverted productos devueltos a como estaban'),
      ));
      await _load(force: true);
    } catch (e) {
      if (!mounted) return;
      messenger.showSnackBar(
        const SnackBar(content: Text('No se pudo deshacer')),
      );
    } finally {
      if (mounted) setState(() => _undoing = false);
    }
  }

  void _openChannel(MyIntegration channel) {
    Navigator.of(context).push(MaterialPageRoute(
      builder: (_) => InventoryCompareScreen(
        channel: channel,
        businessId: widget.businessId,
      ),
    ));
  }

  List<MyIntegration> get _comparableChannels => widget.integrations
      .where((i) =>
          syncProviderFor(i.integrationTypeId)?.supportsCompareInventory == true)
      .toList();

  @override
  Widget build(BuildContext context) {
    final spec = widget.spec;

    return Consumer<SavedComparisonProvider>(
      builder: (context, saved, _) {
        return RefreshIndicator(
          onRefresh: () => _load(force: true),
          color: Theme.of(context).colorScheme.primary,
          child: ListView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 96),
            children: [
              _Intro(spec: spec),
              const SizedBox(height: 12),
              _body(saved),
              const SizedBox(height: 12),
              if (spec.runLabel != null)
                _RunCard(spec: spec, onFinished: () => _load(force: true)),
            ],
          ),
        );
      },
    );
  }

  Widget _body(SavedComparisonProvider saved) {
    switch (widget.spec.environment) {
      case SyncEnvironment.products:
        if (saved.loadingFindings && saved.findings.isEmpty) {
          return const AppLoading(label: 'Leyendo la ultima comparacion');
        }
        if (saved.findingsError != null && saved.findings.isEmpty) {
          return _Notice(
            icon: Icons.error_outline_rounded,
            title: 'No se pudo leer',
            message: saved.findingsError!,
          );
        }
        return FindingsSummaryView(
          report: saved.findings,
          onOpenFinding: _openFinding,
        );

      case SyncEnvironment.data:
        if (saved.loadingData && saved.dataSummary.isEmpty) {
          return const AppLoading(label: 'Leyendo los datos guardados');
        }
        if (saved.dataError != null && saved.dataSummary.isEmpty) {
          return _Notice(
            icon: Icons.error_outline_rounded,
            title: 'No se pudo leer',
            message: saved.dataError!,
          );
        }
        return DataSummaryView(
          summary: saved.dataSummary,
          onPick: _pickData,
          logoFor: _logoFor,
          lastApply: _lastApply,
          undoing: _undoing,
          onUndo: _undo,
        );

      case SyncEnvironment.inventory:
        return InventorySnapshotView(
          channels: _comparableChannels,
          provider: saved,
          onRefreshChannel: (channel) =>
              saved.loadInventorySnapshot(channel, live: true),
          onOpenChannel: _openChannel,
        );

      default:
        return const SizedBox.shrink();
    }
  }
}

class _Intro extends StatelessWidget {
  const _Intro({required this.spec});

  final EnvironmentSpec spec;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      child: Row(
        children: [
          Container(
            width: 34,
            height: 34,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer,
              borderRadius: AppRadius.smAll,
            ),
            child: Icon(spec.icon, size: 18, color: theme.colorScheme.primary),
          ),
          const SizedBox(width: 11),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(spec.label, style: theme.textTheme.titleSmall),
                const SizedBox(height: 3),
                Text(spec.hint, style: theme.textTheme.labelSmall),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _RunCard extends StatelessWidget {
  const _RunCard({required this.spec, required this.onFinished});

  final EnvironmentSpec spec;
  final Future<void> Function() onFinished;

  @override
  Widget build(BuildContext context) {
    return Consumer<SyncActivityProvider>(
      builder: (context, sync, _) {
        final theme = Theme.of(context);
        final channels = sync.eligible;

        if (channels.isEmpty) {
          return const _Notice(
            icon: Icons.link_off_rounded,
            title: 'Sin canales que lo permitan',
            message:
                'Ninguno de tus canales conectados soporta esta accion todavia.',
          );
        }

        return AppCard(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Se ejecuta sobre ${channels.length} canal${channels.length == 1 ? '' : 'es'}',
                style: theme.textTheme.labelSmall,
              ),
              const SizedBox(height: 11),
              SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  onPressed: sync.running
                      ? null
                      : () async {
                          await sync.runCurrent();
                          await onFinished();
                        },
                  icon: sync.running
                      ? const SizedBox(
                          width: 14,
                          height: 14,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation(Colors.white),
                          ),
                        )
                      : Icon(spec.icon, size: 17),
                  label: Text(sync.running ? 'En curso' : spec.runLabel!),
                ),
              ),
              if (sync.finished) ...[
                const SizedBox(height: 8),
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton.icon(
                    onPressed: sync.reset,
                    icon: const Icon(Icons.restart_alt_rounded, size: 17),
                    label: const Text('Reiniciar'),
                  ),
                ),
              ],
              const SizedBox(height: 11),
              Text(
                'Correr vuelve a preguntarle a cada canal y reemplaza la '
                'comparacion guardada de arriba.',
                style: theme.textTheme.labelSmall,
              ),
            ],
          ),
        );
      },
    );
  }
}

class _Notice extends StatelessWidget {
  const _Notice({
    required this.icon,
    required this.title,
    required this.message,
  });

  final IconData icon;
  final String title;
  final String message;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 19, color: AppColors.textMuted),
          const SizedBox(width: 11),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: theme.textTheme.titleSmall),
                const SizedBox(height: 4),
                Text(message, style: theme.textTheme.bodySmall),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
