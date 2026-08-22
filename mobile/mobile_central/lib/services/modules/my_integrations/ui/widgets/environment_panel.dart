import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../app/sync_providers.dart';
import '../../domain/sync_entities.dart';
import '../providers/sync_activity_provider.dart';
import '../providers/saved_comparison_provider.dart';
import 'module_toolbar.dart';
import 'orders_compare_panel.dart';
import 'orders_report.dart';
import 'saved_comparison_views.dart';

class EnvironmentPanel extends StatelessWidget {
  const EnvironmentPanel({
    super.key,
    required this.environment,
    required this.integrations,
    required this.statsFor,
    this.businessId,
  });

  final SyncEnvironment environment;
  final List<MyIntegration> integrations;
  final IntegrationStats Function(int) statsFor;
  final int? businessId;

  @override
  Widget build(BuildContext context) {
    if (environment == SyncEnvironment.overview) {
      return OrdersReport(integrations: integrations, statsFor: statsFor);
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
        return FindingsSummaryView(report: saved.findings);

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
        return DataSummaryView(summary: saved.dataSummary);

      case SyncEnvironment.inventory:
        return InventorySnapshotView(
          channels: _comparableChannels,
          provider: saved,
          onRefreshChannel: (channel) =>
              saved.loadInventorySnapshot(channel, live: true),
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
