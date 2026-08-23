import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../app/sync_providers.dart';
import '../../domain/sync_entities.dart';
import '../providers/sync_activity_provider.dart';

class ChannelSyncStrip extends StatelessWidget {
  const ChannelSyncStrip({
    super.key,
    required this.integrationId,
    required this.integrationTypeId,
  });

  final int integrationId;
  final int integrationTypeId;

  @override
  Widget build(BuildContext context) {
    if (syncProviderFor(integrationTypeId) == null) return const SizedBox.shrink();

    return Consumer<SyncActivityProvider>(
      builder: (context, sync, _) {
        final state = sync.stateFor(integrationId);
        final result = sync.resultFor(integrationId);
        final busy = state == SyncNodeState.active ||
            state == SyncNodeState.scan ||
            state == SyncNodeState.queued;

        if (busy) {
          return _Running(
            state: state,
            progress: sync.progressFor(integrationId),
          );
        }

        if (result != null) return _Outcome(result: result);

        final inventory = sync.lastRunFor(integrationId, SyncRunKind.inventory);
        final products = sync.lastRunFor(integrationId, SyncRunKind.products);
        final latest = _latest(inventory, products);
        if (latest == null) return const SizedBox.shrink();

        return _LastRun(record: latest);
      },
    );
  }

  SyncRunRecord? _latest(SyncRunRecord? a, SyncRunRecord? b) {
    final left = a?.finishedOn;
    final right = b?.finishedOn;
    if (left == null) return b;
    if (right == null) return a;
    return left.isAfter(right) ? a : b;
  }
}

class _Running extends StatelessWidget {
  const _Running({required this.state, required this.progress});

  final SyncNodeState state;
  final SyncProgress progress;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final label = switch (state) {
      SyncNodeState.queued => 'En cola',
      SyncNodeState.scan => 'Comparando con el canal',
      _ => 'Sincronizando',
    };
    final ratio = progress.ratio;

    return Padding(
      padding: const EdgeInsets.only(top: 10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              SizedBox(
                width: 12,
                height: 12,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation(theme.colorScheme.primary),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  label,
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.primary,
                    fontWeight: FontWeight.w600,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              if (progress.total > 0)
                Text(
                  '${AppFormat.number(progress.processed)} / ${AppFormat.number(progress.total)}',
                  style: theme.textTheme.labelSmall,
                ),
            ],
          ),
          const SizedBox(height: 7),
          ClipRRect(
            borderRadius: BorderRadius.circular(3),
            child: LinearProgressIndicator(
              value: ratio,
              minHeight: 4,
              backgroundColor: AppColors.surfaceMuted,
              valueColor: AlwaysStoppedAnimation(theme.colorScheme.primary),
            ),
          ),
        ],
      ),
    );
  }
}

class _Outcome extends StatelessWidget {
  const _Outcome({required this.result});

  final SyncResult result;

  @override
  Widget build(BuildContext context) {
    final (icon, color, text) = switch (result) {
      SyncErrorResult(message: final message) => (
          Icons.error_outline_rounded,
          AppColors.error,
          message,
        ),
      InventorySyncResult(
        updated: final updated,
        unchanged: final unchanged,
        failed: final failed,
      ) =>
        (
          failed > 0 ? Icons.warning_amber_rounded : Icons.check_circle_outline_rounded,
          failed > 0 ? AppColors.warning : AppColors.success,
          _inventoryText(updated, unchanged, failed),
        ),
      ProductsSyncResult(matched: final matched, pending: final pending) => (
          pending > 0 ? Icons.warning_amber_rounded : Icons.check_circle_outline_rounded,
          pending > 0 ? AppColors.warning : AppColors.success,
          pending > 0
              ? '${AppFormat.number(matched)} asociados, ${AppFormat.number(pending)} por revisar'
              : '${AppFormat.number(matched)} asociados, sin pendientes',
        ),
      _ => (Icons.info_outline_rounded, AppColors.textMuted, ''),
    };

    if (text.isEmpty) return const SizedBox.shrink();
    return _Line(icon: icon, color: color, text: text);
  }

  String _inventoryText(int updated, int unchanged, int failed) {
    final parts = <String>[];
    if (updated > 0) parts.add('${AppFormat.number(updated)} actualizados');
    if (unchanged > 0) parts.add('${AppFormat.number(unchanged)} sin cambio');
    if (failed > 0) parts.add('${AppFormat.number(failed)} con error');
    return parts.isEmpty ? 'Sin cambios' : parts.join(', ');
  }
}

class _LastRun extends StatelessWidget {
  const _LastRun({required this.record});

  final SyncRunRecord record;

  @override
  Widget build(BuildContext context) {
    final when = record.finishedOn;
    final kind = record.kind == SyncRunKind.inventory ? 'inventario' : 'productos';
    final ago = when == null ? '' : ' ${AppFormat.relative(when)}';

    return _Line(
      icon: Icons.history_rounded,
      color: AppColors.textMuted,
      text: 'Ultima comparacion de $kind$ago',
    );
  }
}

class _Line extends StatelessWidget {
  const _Line({required this.icon, required this.color, required this.text});

  final IconData icon;
  final Color color;
  final String text;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(top: 10),
      child: Row(
        children: [
          Icon(icon, size: 14, color: color),
          const SizedBox(width: 7),
          Expanded(
            child: Text(
              text,
              style: Theme.of(context).textTheme.labelSmall?.copyWith(color: color),
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      ),
    );
  }
}
