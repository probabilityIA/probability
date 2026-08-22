import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../domain/sync_entities.dart';
import '../providers/sync_activity_provider.dart';
import 'module_toolbar.dart';
import 'orders_compare_panel.dart';
import 'orders_report.dart';

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

    return _Pending(spec: spec);
  }
}

class _Pending extends StatelessWidget {
  const _Pending({required this.spec});

  final EnvironmentSpec spec;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final runnable = spec.runLabel != null;

    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 96),
      children: [
        AppCard(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
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
                    child: Text(spec.label, style: theme.textTheme.titleSmall),
                  ),
                ],
              ),
              const SizedBox(height: 11),
              Text(spec.hint, style: theme.textTheme.bodySmall),
            ],
          ),
        ),
        const SizedBox(height: 12),
        if (runnable) _RunCard(spec: spec),
      ],
    );
  }
}

class _RunCard extends StatelessWidget {
  const _RunCard({required this.spec});

  final EnvironmentSpec spec;

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
                  onPressed: sync.running ? null : sync.runCurrent,
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
                'El detalle por producto llega en una fase siguiente. Por ahora el '
                'avance y el resultado de cada canal se ven en el diagrama.',
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
