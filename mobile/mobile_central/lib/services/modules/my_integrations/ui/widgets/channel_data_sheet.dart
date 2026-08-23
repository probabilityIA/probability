import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../core/network/sse_client.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/saved_comparison_provider.dart';

const List<String> channelDataEventTypes = <String>[
  'products.channel_data.progress',
  'products.channel_data.completed',
  'products.channel_data.failed',
];

class ChannelDataTarget {
  const ChannelDataTarget({
    required this.field,
    required this.label,
    required this.cell,
    required this.mode,
  });

  final String field;
  final String label;
  final DataSummaryCell cell;
  final DataMode mode;
}

Future<DataApplyResult?> showChannelDataSheet(
  BuildContext context, {
  required ChannelDataTarget target,
  int? businessId,
}) {
  return showModalBottomSheet<DataApplyResult>(
    context: context,
    isScrollControlled: true,
    backgroundColor: AppColors.surface,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(AppRadius.xl)),
    ),
    builder: (_) => ChannelDataSheet(target: target, businessId: businessId),
  );
}

class ChannelDataSheet extends StatefulWidget {
  const ChannelDataSheet({super.key, required this.target, this.businessId});

  final ChannelDataTarget target;
  final int? businessId;

  @override
  State<ChannelDataSheet> createState() => _ChannelDataSheetState();
}

class _ChannelDataSheetState extends State<ChannelDataSheet> {
  final SseClient _sse = SseClient();
  StreamSubscription<SseEvent>? _subscription;

  DataPreview? _preview;
  bool _loading = true;
  bool _applying = false;
  String? _error;
  String? _batchId;
  int _processed = 0;
  int _total = 0;

  @override
  void initState() {
    super.initState();
    _loadPreview();
  }

  @override
  void dispose() {
    _subscription?.cancel();
    _sse.dispose();
    super.dispose();
  }

  Future<void> _loadPreview() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final preview =
          await context.read<SavedComparisonProvider>().previewChannelData(
                integrationId: widget.target.cell.integrationId,
                field: widget.target.field,
                mode: widget.target.mode,
              );
      if (!mounted) return;
      setState(() => _preview = preview);
    } catch (e) {
      if (!mounted) return;
      setState(() => _error = parseError(e));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _listen() async {
    await _subscription?.cancel();
    _subscription = _sse.events.listen(_onEvent);
    await _sse.connect(
      businessId: widget.businessId ?? 0,
      eventTypes: channelDataEventTypes,
    );
  }

  void _onEvent(SseEvent event) {
    final data = event.data['data'];
    if (data is! Map<String, dynamic>) return;
    if (_batchId == null || data['batch_id']?.toString() != _batchId) return;

    if (event.type.endsWith('.progress')) {
      setState(() {
        _processed = int.tryParse('${data['processed']}') ?? _processed;
        _total = int.tryParse('${data['total']}') ?? _total;
      });
      return;
    }
    if (event.type.endsWith('.failed')) {
      setState(() {
        _applying = false;
        _error = data['error']?.toString() ?? 'No se pudo completar';
      });
      return;
    }
    if (event.type.endsWith('.completed')) {
      final applied = int.tryParse('${data['applied']}') ?? 0;
      Navigator.of(context).pop(DataApplyResult(
        batchId: _batchId!,
        field: widget.target.field,
        applied: applied,
      ));
    }
  }

  Future<void> _apply() async {
    final saved = context.read<SavedComparisonProvider>();
    setState(() {
      _applying = true;
      _error = null;
      _processed = 0;
      _total = _preview?.total ?? 0;
    });

    await _listen();

    try {
      final result = await saved.applyChannelData(
        integrationId: widget.target.cell.integrationId,
        field: widget.target.field,
        mode: widget.target.mode,
      );
      if (!mounted) return;
      if (result.batchId.isEmpty) {
        setState(() {
          _applying = false;
          _error = 'No se pudo iniciar la actualizacion';
        });
        return;
      }
      setState(() => _batchId = result.batchId);
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _applying = false;
        _error = parseError(e);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final target = widget.target;
    final preview = _preview;
    final maxHeight = MediaQuery.sizeOf(context).height * 0.85;

    return SafeArea(
      top: false,
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: maxHeight),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 14, 20, 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Actualizar ${target.label.toLowerCase()} en Probability',
                    style: theme.textTheme.titleMedium,
                  ),
                  const SizedBox(height: 5),
                  Text(
                    'Se escribe en tus productos con lo que dice '
                    '${target.cell.integrationName}. El canal no se modifica. '
                    '${target.mode.help}',
                    style: theme.textTheme.labelSmall,
                  ),
                ],
              ),
            ),
            Flexible(
              child: ListView(
                shrinkWrap: true,
                padding: const EdgeInsets.fromLTRB(20, 4, 20, 8),
                children: [
                  if (_loading)
                    const Padding(
                      padding: EdgeInsets.symmetric(vertical: 28),
                      child: AppLoading(label: 'Revisando que se veria afectado'),
                    )
                  else if (preview != null) ...[
                    Row(
                      children: [
                        Expanded(
                          child: _Count(
                            label: 'se llenan',
                            value: preview.wouldFill,
                            color: AppColors.success,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: _Count(
                            label: 'se reemplazan',
                            value: preview.wouldReplace,
                            color: preview.wouldReplace > 0
                                ? AppColors.warning
                                : AppColors.textDisabled,
                          ),
                        ),
                      ],
                    ),
                    if (preview.conflicts.isNotEmpty) ...[
                      const SizedBox(height: 12),
                      _Conflicts(conflicts: preview.conflicts),
                    ],
                    if (preview.samples.isNotEmpty) ...[
                      const SizedBox(height: 12),
                      Text('Ejemplos', style: theme.textTheme.titleSmall),
                      const SizedBox(height: 8),
                      for (final sample in preview.samples) ...[
                        _Sample(sample: sample),
                        const SizedBox(height: 8),
                      ],
                    ],
                    const SizedBox(height: 4),
                    Text(
                      'Queda registrado quien escribio cada dato y se puede '
                      'deshacer completo despues de aplicar.',
                      style: theme.textTheme.labelSmall,
                    ),
                  ],
                ],
              ),
            ),
            if (_error != null)
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 4, 20, 4),
                child: Text(
                  _error!,
                  style: theme.textTheme.labelSmall
                      ?.copyWith(color: AppColors.error),
                ),
              ),
            if (_applying && _error == null)
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 4, 20, 4),
                child: Column(
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            _total > 0
                                ? '${AppFormat.number(_processed)} de ${AppFormat.number(_total)}'
                                : 'Preparando...',
                            style: theme.textTheme.labelSmall,
                          ),
                        ),
                        Text(
                          '${_total > 0 ? ((_processed / _total) * 100).round() : 0}%',
                          style: theme.textTheme.labelSmall?.copyWith(
                            fontWeight: FontWeight.w700,
                            color: theme.colorScheme.primary,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 6),
                    ClipRRect(
                      borderRadius: BorderRadius.circular(3),
                      child: LinearProgressIndicator(
                        value: _total > 0 ? _processed / _total : null,
                        minHeight: 4,
                        backgroundColor: AppColors.surfaceMuted,
                      ),
                    ),
                  ],
                ),
              ),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 10, 20, 14),
              child: Row(
                children: [
                  TextButton(
                    onPressed: _applying ? null : () => Navigator.pop(context),
                    child: const Text('Cancelar'),
                  ),
                  const Spacer(),
                  FilledButton.icon(
                    onPressed: _loading ||
                            _applying ||
                            preview == null ||
                            !preview.hasSomething
                        ? null
                        : _apply,
                    icon: _applying
                        ? const SizedBox(
                            width: 14,
                            height: 14,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation(Colors.white),
                            ),
                          )
                        : const Icon(Icons.download_rounded, size: 17),
                    label: Text(
                      preview == null || !preview.hasSomething
                          ? 'Nada por actualizar'
                          : 'Actualizar ${AppFormat.number(preview.total)}',
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Count extends StatelessWidget {
  const _Count({required this.label, required this.value, required this.color});

  final String label;
  final int value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.10),
        borderRadius: AppRadius.smAll,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            AppFormat.number(value),
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 20,
              fontWeight: FontWeight.w700,
              color: color,
            ),
          ),
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

class _Conflicts extends StatelessWidget {
  const _Conflicts({required this.conflicts});

  final List<DataConflict> conflicts;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.warning.withValues(alpha: 0.10),
        borderRadius: AppRadius.smAll,
        border: Border.all(color: AppColors.warning.withValues(alpha: 0.35)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.warning_amber_rounded,
                  size: 15, color: AppColors.warning),
              const SizedBox(width: 7),
              Expanded(
                child: Text(
                  'Ojo, este dato ya lo escribio alguien mas',
                  style: theme.textTheme.labelMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                    color: AppColors.warning,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 7),
          for (final conflict in conflicts)
            Padding(
              padding: const EdgeInsets.only(bottom: 3),
              child: Text(
                '${AppFormat.number(conflict.count)} '
                '${conflict.count == 1 ? "producto viene" : "productos vienen"} '
                'de ${conflict.who}'
                '${conflict.lastChangeAt != null ? " - ${AppFormat.relative(conflict.lastChangeAt)}" : ""}',
                style: theme.textTheme.labelSmall,
              ),
            ),
        ],
      ),
    );
  }
}

class _Sample extends StatelessWidget {
  const _Sample({required this.sample});

  final DataSample sample;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.all(11),
      decoration: BoxDecoration(
        color: AppColors.surfaceMuted,
        borderRadius: AppRadius.smAll,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            sample.sku,
            style: theme.textTheme.labelMedium?.copyWith(
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 6),
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('HOY',
                        style: theme.textTheme.labelSmall?.copyWith(fontSize: 9)),
                    const SizedBox(height: 2),
                    Text(
                      sample.isEmptyToday
                          ? 'vacio'
                          : DataSample.clean(sample.current),
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: AppColors.textDisabled,
                        fontStyle: sample.isEmptyToday
                            ? FontStyle.italic
                            : FontStyle.normal,
                      ),
                    ),
                  ],
                ),
              ),
              const Padding(
                padding: EdgeInsets.symmetric(horizontal: 8),
                child: Icon(Icons.arrow_forward_rounded,
                    size: 13, color: AppColors.textDisabled),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('QUEDARIA',
                        style: theme.textTheme.labelSmall?.copyWith(fontSize: 9)),
                    const SizedBox(height: 2),
                    Text(
                      DataSample.clean(sample.incoming),
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: AppColors.success,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
