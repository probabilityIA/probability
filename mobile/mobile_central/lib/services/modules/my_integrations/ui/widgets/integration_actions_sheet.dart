import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../../integrations/core/ui/providers/integration_provider.dart';
import '../../domain/entities.dart';

Future<bool?> showIntegrationActionsSheet(
  BuildContext context, {
  required MyIntegration integration,
}) {
  return showModalBottomSheet<bool>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    builder: (sheetContext) => _ActionsSheet(integration: integration),
  );
}

class _ActionsSheet extends StatefulWidget {
  const _ActionsSheet({required this.integration});

  final MyIntegration integration;

  @override
  State<_ActionsSheet> createState() => _ActionsSheetState();
}

class _ActionsSheetState extends State<_ActionsSheet> {
  String? _running;
  String? _result;
  bool _resultOk = true;

  Future<void> _run(String key, Future<bool> Function() action) async {
    setState(() {
      _running = key;
      _result = null;
    });

    final ok = await action();

    if (!mounted) return;
    setState(() {
      _running = null;
      _resultOk = ok;
      _result = ok ? 'Listo' : (context.read<IntegrationProvider>().error ?? 'No se pudo completar');
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final integration = widget.integration;
    final provider = context.read<IntegrationProvider>();

    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 4, 20, 24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
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
                    Text(integration.name, style: theme.textTheme.titleMedium),
                    const SizedBox(height: 2),
                    Text(
                      integration.integrationTypeName ?? '',
                      style: theme.textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
              AppStatusChip(
                dense: true,
                label: integration.isActive ? 'Activa' : 'Inactiva',
                tone: integration.isActive ? AppStatusTone.success : AppStatusTone.neutral,
              ),
            ],
          ),
          const SizedBox(height: 18),
          _ActionRow(
            icon: Icons.wifi_tethering_rounded,
            label: 'Probar conexion',
            busy: _running == 'test',
            onTap: () => _run('test', () async {
              final response = await provider.testConnection(integration.id);
              return response != null;
            }),
          ),
          const Divider(height: 1),
          _ActionRow(
            icon: Icons.sync_rounded,
            label: 'Sincronizar ordenes',
            busy: _running == 'sync',
            onTap: () => _run('sync', () async {
              final response = await provider.syncOrders(integration.id);
              return response != null;
            }),
          ),
          const Divider(height: 1),
          _ActionRow(
            icon: integration.isActive
                ? Icons.pause_circle_outline_rounded
                : Icons.play_circle_outline_rounded,
            label: integration.isActive ? 'Desactivar' : 'Activar',
            busy: _running == 'toggle',
            tone: integration.isActive ? AppColors.error : AppColors.success,
            onTap: () => _run('toggle', () async {
              return integration.isActive
                  ? provider.deactivateIntegration(integration.id)
                  : provider.activateIntegration(integration.id);
            }),
          ),
          if (_result != null) ...[
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 11),
              decoration: BoxDecoration(
                color: _resultOk ? AppColors.successSoft : AppColors.errorSoft,
                borderRadius: AppRadius.mdAll,
              ),
              child: Row(
                children: [
                  Icon(
                    _resultOk ? Icons.check_circle_outline : Icons.error_outline,
                    size: 17,
                    color: _resultOk ? const Color(0xFF047857) : AppColors.error,
                  ),
                  const SizedBox(width: 9),
                  Expanded(
                    child: Text(
                      _result!,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: _resultOk ? const Color(0xFF047857) : const Color(0xFFB91C1C),
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
          const SizedBox(height: 18),
          Text(
            'Actualizada ${AppFormat.relative(AppFormat.parseDate(integration.updatedAt))}',
            style: theme.textTheme.labelSmall,
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}

class _ActionRow extends StatelessWidget {
  const _ActionRow({
    required this.icon,
    required this.label,
    required this.onTap,
    this.busy = false,
    this.tone = AppColors.textSecondary,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool busy;
  final Color tone;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: busy ? null : onTap,
      borderRadius: AppRadius.mdAll,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15),
        child: Row(
          children: [
            Icon(icon, size: 20, color: tone),
            const SizedBox(width: 13),
            Expanded(
              child: Text(
                label,
                style: Theme.of(context).textTheme.titleSmall?.copyWith(color: tone),
              ),
            ),
            if (busy)
              const SizedBox(
                height: 17,
                width: 17,
                child: CircularProgressIndicator(strokeWidth: 2),
              )
            else
              const Icon(Icons.chevron_right_rounded, size: 20, color: AppColors.textDisabled),
          ],
        ),
      ),
    );
  }
}
