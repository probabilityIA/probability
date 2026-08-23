import 'package:flutter/material.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/saved_comparison_provider.dart';

class SavedStamp extends StatelessWidget {
  const SavedStamp({super.key, required this.when, this.live = false});

  final DateTime? when;
  final bool live;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final text = when == null
        ? 'Todavia no se ha comparado'
        : live
            ? 'Consultado al canal ${AppFormat.relative(when)}'
            : 'Ultima comparacion guardada ${AppFormat.relative(when)}';

    return Row(
      children: [
        Icon(
          when == null ? Icons.help_outline_rounded : Icons.history_rounded,
          size: 14,
          color: AppColors.textMuted,
        ),
        const SizedBox(width: 7),
        Expanded(
          child: Text(
            text,
            style: theme.textTheme.labelSmall,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}

class FindingsSummaryView extends StatelessWidget {
  const FindingsSummaryView({
    super.key,
    required this.report,
    required this.onOpenFinding,
  });

  final FindingsReport report;
  final void Function(Finding) onOpenFinding;

  @override
  Widget build(BuildContext context) {
    if (report.channels.isEmpty && report.findings.isEmpty) {
      return const AppCard(
        child: Text(
          'Todavia no hay una comparacion de productos guardada. Corre la '
          'comparacion para ver que producto esta en cada canal.',
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (report.channels.isNotEmpty) ...[
          const AppSectionHeader(title: 'Por canal'),
          for (final channel in report.channels) ...[
            _ChannelFindingCard(channel: channel),
            const SizedBox(height: 10),
          ],
        ],
        if (report.findings.isNotEmpty) ...[
          const SizedBox(height: 6),
          const AppSectionHeader(title: 'Que hay por revisar'),
          for (final finding in report.findings) ...[
            _FindingCard(
              finding: finding,
              onTap: () => onOpenFinding(finding),
            ),
            const SizedBox(height: 10),
          ],
        ],
      ],
    );
  }
}

class _ChannelFindingCard extends StatelessWidget {
  const _ChannelFindingCard({required this.channel});

  final FindingChannelSummary channel;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(name: channel.integrationName, size: 30, radius: 8, padding: 4),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  channel.integrationName,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              AppStatusChip(
                dense: true,
                label: channel.isClean ? 'sin pendientes' : '${channel.pending} por revisar',
                tone: channel.isClean ? AppStatusTone.success : AppStatusTone.warning,
              ),
            ],
          ),
          const SizedBox(height: 10),
          Wrap(
            spacing: 6,
            runSpacing: 6,
            children: [
              _Pill(label: 'asociados', value: channel.matched, tone: AppColors.success),
              if (channel.notAssociated > 0)
                _Pill(label: 'sin asociar', value: channel.notAssociated),
              if (channel.onlyInChannel > 0)
                _Pill(label: 'solo en el canal', value: channel.onlyInChannel),
              if (channel.channelNoSku > 0)
                _Pill(label: 'sin SKU', value: channel.channelNoSku),
              if (channel.skuChanged > 0)
                _Pill(label: 'SKU cambiado', value: channel.skuChanged),
              if (channel.skuTypo > 0)
                _Pill(label: 'SKU con error', value: channel.skuTypo),
            ],
          ),
          const SizedBox(height: 9),
          SavedStamp(when: channel.comparedAt),
        ],
      ),
    );
  }
}

class _FindingCard extends StatelessWidget {
  const _FindingCard({required this.finding, required this.onTap});

  final Finding finding;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final tone = switch (finding.severity) {
      FindingSeverity.error => AppColors.error,
      FindingSeverity.warn => AppColors.warning,
      FindingSeverity.info => AppColors.info,
    };

    return AppCard(
      onTap: onTap,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 42,
            height: 42,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: tone.withValues(alpha: 0.12),
              borderRadius: AppRadius.smAll,
            ),
            child: Text(
              AppFormat.compact(finding.count),
              style: theme.textTheme.titleSmall?.copyWith(color: tone, fontSize: 13),
            ),
          ),
          const SizedBox(width: 11),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(finding.title, style: theme.textTheme.titleSmall),
                const SizedBox(height: 3),
                Text(finding.detail, style: theme.textTheme.bodySmall),
                if (finding.channels.isNotEmpty) ...[
                  const SizedBox(height: 7),
                  Wrap(
                    spacing: 6,
                    runSpacing: 6,
                    children: finding.channels
                        .map((name) => AppStatusChip(
                              dense: true,
                              label: name,
                              tone: AppStatusTone.neutral,
                            ))
                        .toList(),
                  ),
                ],
              ],
            ),
          ),
          const Icon(Icons.chevron_right_rounded, color: AppColors.textDisabled),
        ],
      ),
    );
  }
}

class InventorySnapshotView extends StatelessWidget {
  const InventorySnapshotView({
    super.key,
    required this.channels,
    required this.provider,
    required this.onRefreshChannel,
    required this.onOpenChannel,
  });

  final List<MyIntegration> channels;
  final SavedComparisonProvider provider;
  final void Function(MyIntegration) onRefreshChannel;
  final void Function(MyIntegration) onOpenChannel;

  @override
  Widget build(BuildContext context) {
    if (channels.isEmpty) {
      return const AppCard(
        child: Text(
          'Ninguno de tus canales conectados permite comparar inventario. '
          'Los que si lo permiten muestran aqui la ultima comparacion guardada.',
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        for (final channel in channels) ...[
          _InventoryChannelCard(
            channel: channel,
            snapshot: provider.snapshotFor(channel.id),
            onRefresh: () => onRefreshChannel(channel),
            onTap: () => onOpenChannel(channel),
          ),
          const SizedBox(height: 10),
        ],
      ],
    );
  }
}

class _InventoryChannelCard extends StatelessWidget {
  const _InventoryChannelCard({
    required this.channel,
    required this.snapshot,
    required this.onRefresh,
    required this.onTap,
  });

  final MyIntegration channel;
  final ChannelSnapshot snapshot;
  final VoidCallback onRefresh;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = channel.integrationTypeName ?? channel.name;
    final totals = snapshot.page?.totals;

    return AppCard(
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              BrandLogo(
                name: name,
                imageUrl: channel.imageUrl,
                size: 30,
                radius: 8,
                padding: 4,
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  name,
                  style: theme.textTheme.titleSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              IconButton(
                onPressed: snapshot.loading ? null : onRefresh,
                icon: snapshot.loading
                    ? const SizedBox(
                        width: 15,
                        height: 15,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.refresh_rounded, size: 19),
                tooltip: 'Preguntarle el stock al canal ahora',
                visualDensity: VisualDensity.compact,
              ),
            ],
          ),
          if (snapshot.error != null) ...[
            const SizedBox(height: 8),
            Text(
              snapshot.error!,
              style: theme.textTheme.labelSmall?.copyWith(color: AppColors.error),
            ),
          ] else if (!snapshot.hasData) ...[
            const SizedBox(height: 8),
            Text(
              snapshot.loading
                  ? 'Leyendo la ultima comparacion'
                  : 'Sin comparacion guardada todavia',
              style: theme.textTheme.labelSmall,
            ),
          ] else ...[
            const SizedBox(height: 10),
            Wrap(
              spacing: 6,
              runSpacing: 6,
              children: [
                _Pill(
                  label: 'por actualizar',
                  value: totals?.toUpdate ?? 0,
                  tone: (totals?.toUpdate ?? 0) > 0 ? AppColors.warning : null,
                ),
                _Pill(label: 'iguales', value: totals?.unchanged ?? 0, tone: AppColors.success),
                if ((totals?.skipped ?? 0) > 0)
                  _Pill(label: 'omitidos', value: totals!.skipped),
              ],
            ),
            const SizedBox(height: 9),
            Row(
              children: [
                Expanded(
                  child: SavedStamp(
                    when: snapshot.checkedAt,
                    live: !snapshot.isSaved,
                  ),
                ),
                Text(
                  'Ver productos',
                  style: theme.textTheme.labelSmall?.copyWith(
                    color: theme.colorScheme.primary,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const Icon(Icons.chevron_right_rounded,
                    size: 17, color: AppColors.textDisabled),
              ],
            ),
          ],
        ],
      ),
    );
  }
}

class DataSummaryView extends StatelessWidget {
  const DataSummaryView({
    super.key,
    required this.summary,
    required this.onPick,
    this.lastApply,
    this.onUndo,
    this.undoing = false,
  });

  final DataSummary summary;
  final void Function(DataSummaryRow, DataSummaryCell, DataMode) onPick;
  final DataApplyResult? lastApply;
  final VoidCallback? onUndo;
  final bool undoing;

  @override
  Widget build(BuildContext context) {
    if (summary.isEmpty) {
      return const AppCard(
        child: Text(
          'Todavia no hay una lectura guardada de los datos de tus canales. '
          'Se calcula cuando se compara el catalogo.',
        ),
      );
    }

    final rows = summary.actionable;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SavedStamp(when: summary.snapshotAt),
        if (lastApply != null) ...[
          const SizedBox(height: 12),
          _UndoCard(result: lastApply!, onUndo: onUndo, undoing: undoing),
        ],
        const SizedBox(height: 12),
        if (rows.isEmpty)
          const AppCard(
            child: Text(
              'Ningun canal tiene datos que puedan entrar a Probability: '
              'todo lo que se puede traer ya esta aqui.',
            ),
          )
        else
          for (final row in rows) ...[
            _DataRowCard(row: row, onPick: onPick),
            const SizedBox(height: 10),
          ],
      ],
    );
  }
}

class _UndoCard extends StatelessWidget {
  const _UndoCard({
    required this.result,
    required this.onUndo,
    required this.undoing,
  });

  final DataApplyResult result;
  final VoidCallback? onUndo;
  final bool undoing;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AppCard(
      child: Row(
        children: [
          const Icon(Icons.check_circle_outline_rounded,
              size: 19, color: AppColors.success),
          const SizedBox(width: 11),
          Expanded(
            child: Text(
              '${AppFormat.number(result.applied)} productos actualizados',
              style: theme.textTheme.bodySmall,
            ),
          ),
          TextButton.icon(
            onPressed: undoing ? null : onUndo,
            icon: undoing
                ? const SizedBox(
                    width: 14,
                    height: 14,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.undo_rounded, size: 17),
            label: const Text('Deshacer'),
          ),
        ],
      ),
    );
  }
}

class _DataRowCard extends StatelessWidget {
  const _DataRowCard({required this.row, required this.onPick});

  final DataSummaryRow row;
  final void Function(DataSummaryRow, DataSummaryCell, DataMode) onPick;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(row.label.isEmpty ? row.field : row.label,
              style: theme.textTheme.titleSmall),
          if (row.note != null && row.note!.isNotEmpty) ...[
            const SizedBox(height: 3),
            Text(row.note!, style: theme.textTheme.labelSmall),
          ],
          const SizedBox(height: 10),
          for (final cell in row.cells.where((c) => c.hasSomething)) ...[
            Padding(
              padding: const EdgeInsets.only(bottom: 9),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      BrandLogo(
                        name: cell.integrationName,
                        size: 20,
                        radius: 5,
                        padding: 3,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          cell.integrationName,
                          style: theme.textTheme.bodySmall,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 7),
                  Wrap(
                    spacing: 7,
                    runSpacing: 7,
                    children: [
                      if (cell.canFill > 0)
                        _ActionPill(
                          label: 'Llenar vacios',
                          value: cell.canFill,
                          tone: AppColors.info,
                          onTap: () => onPick(row, cell, DataMode.fillEmpty),
                        ),
                      if (cell.canOverwrite > 0)
                        _ActionPill(
                          label: 'Reemplazar',
                          value: cell.canOverwrite,
                          tone: AppColors.warning,
                          onTap: () => onPick(row, cell, DataMode.overwrite),
                        ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _ActionPill extends StatelessWidget {
  const _ActionPill({
    required this.label,
    required this.value,
    required this.tone,
    required this.onTap,
  });

  final String label;
  final int value;
  final Color tone;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: tone.withValues(alpha: 0.12),
      borderRadius: AppRadius.pillAll,
      child: InkWell(
        borderRadius: AppRadius.pillAll,
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(11, 7, 9, 7),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                '$label ${AppFormat.number(value)}',
                style: TextStyle(
                  fontFamily: 'Inter',
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: tone,
                ),
              ),
              const SizedBox(width: 4),
              Icon(Icons.chevron_right_rounded, size: 15, color: tone),
            ],
          ),
        ),
      ),
    );
  }
}

class _Pill extends StatelessWidget {
  const _Pill({required this.label, required this.value, this.tone});

  final String label;
  final int value;
  final Color? tone;

  @override
  Widget build(BuildContext context) {
    final color = tone ?? AppColors.textMuted;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 5),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.10),
        borderRadius: AppRadius.pillAll,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            AppFormat.number(value),
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 12,
              fontWeight: FontWeight.w700,
              color: color,
            ),
          ),
          const SizedBox(width: 5),
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 11,
              fontWeight: FontWeight.w500,
              color: color,
            ),
          ),
        ],
      ),
    );
  }
}
