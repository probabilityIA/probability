import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/image_memory.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/comparison_lists_provider.dart';
import '../providers/sync_activity_provider.dart';
import '../widgets/saved_comparison_views.dart';

class InventoryCompareScreen extends StatefulWidget {
  const InventoryCompareScreen({
    super.key,
    required this.channel,
    this.businessId,
  });

  final MyIntegration channel;
  final int? businessId;

  @override
  State<InventoryCompareScreen> createState() => _InventoryCompareScreenState();
}

class _InventoryCompareScreenState extends State<InventoryCompareScreen> {
  final TextEditingController _search = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      context.read<InventoryRowsProvider>().load(
            channel: widget.channel,
            businessId: widget.businessId,
          );
    });
  }

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  Future<void> _push() async {
    final rows = context.read<InventoryRowsProvider>();
    final sync = context.read<SyncActivityProvider>();
    final messenger = ScaffoldMessenger.of(context);
    final pending = rows.totals.toUpdate;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Enviar stock al canal'),
        content: Text(
          'Se le va a mandar a ${widget.channel.integrationTypeName ?? widget.channel.name} '
          'el stock de Probability para los $pending productos que quedaron '
          'distintos. Es una accion que escribe en el canal.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Enviar'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    await sync.runInventoryOne(widget.channel.id);
    if (!mounted) return;

    messenger.showSnackBar(
      const SnackBar(content: Text('Sincronizacion lanzada')),
    );
    await rows.load(
      channel: widget.channel,
      businessId: widget.businessId,
      live: true,
    );
  }

  @override
  Widget build(BuildContext context) {
    final name = widget.channel.integrationTypeName ?? widget.channel.name;

    return Consumer<InventoryRowsProvider>(
      builder: (context, provider, _) {
        final totals = provider.totals;

        return AppScaffold(
          title: name,
          subtitle: 'Stock del canal frente al de Probability',
          onBack: () => Navigator.of(context).maybePop(),
          body: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 10, 16, 8),
                child: Column(
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _search,
                            textInputAction: TextInputAction.search,
                            decoration: const InputDecoration(
                              isDense: true,
                              hintText: 'Buscar SKU o producto',
                              prefixIcon: Icon(Icons.search_rounded, size: 19),
                            ),
                            onSubmitted: provider.setSearch,
                          ),
                        ),
                        const SizedBox(width: 8),
                        _DiffToggle(provider: provider),
                      ],
                    ),
                    const SizedBox(height: 9),
                    Row(
                      children: [
                        Expanded(
                          child: SavedStamp(
                            when: provider.checkedAt,
                            live: !provider.fromCache,
                          ),
                        ),
                        TextButton.icon(
                          onPressed: provider.list.isLoading
                              ? null
                              : provider.askChannel,
                          icon: const Icon(Icons.cloud_sync_outlined, size: 17),
                          label: const Text('Preguntar al canal'),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              Expanded(
                child: PaginatedListView<InventoryCompareRow>(
                  controller: provider.list,
                  unitLabel: 'productos',
                  emptyIcon: Icons.inventory_2_outlined,
                  emptyTitle: 'Sin comparacion',
                  emptyMessage:
                      'No hay una comparacion de inventario guardada para este '
                      'canal. Toca "Preguntar al canal" para hacer una.',
                  placeholderHeight: 86,
                  onRefresh: () => provider.load(
                    channel: widget.channel,
                    businessId: widget.businessId,
                  ),
                  header: _Totals(totals: totals),
                  itemBuilder: (context, row, index) => _StockCard(row: row),
                ),
              ),
              if (totals.toUpdate > 0)
                _PushBar(count: totals.toUpdate, onPush: _push),
            ],
          ),
        );
      },
    );
  }
}

class _DiffToggle extends StatelessWidget {
  const _DiffToggle({required this.provider});

  final InventoryRowsProvider provider;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final active = provider.onlyDiff;

    return GestureDetector(
      onTap: provider.toggleOnlyDiff,
      behavior: HitTestBehavior.opaque,
      child: Container(
        height: 44,
        padding: const EdgeInsets.symmetric(horizontal: 12),
        alignment: Alignment.center,
        decoration: BoxDecoration(
          color: active ? scheme.primaryContainer : AppColors.surface,
          borderRadius: AppRadius.smAll,
          border: Border.all(color: active ? scheme.primary : AppColors.border),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.filter_alt_rounded,
              size: 15,
              color: active ? scheme.primary : AppColors.textMuted,
            ),
            const SizedBox(width: 6),
            Text(
              'Solo dif.',
              style: TextStyle(
                fontFamily: 'Inter',
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: active ? scheme.primary : AppColors.textMuted,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Totals extends StatelessWidget {
  const _Totals({required this.totals});

  final InventoryCompareTotals totals;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Expanded(
            child: _Chip(
              label: 'por actualizar',
              value: totals.toUpdate,
              color: totals.toUpdate > 0 ? AppColors.warning : AppColors.textMuted,
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: _Chip(
              label: 'iguales',
              value: totals.unchanged,
              color: AppColors.success,
            ),
          ),
          if (totals.skipped > 0) ...[
            const SizedBox(width: 8),
            Expanded(
              child: _Chip(
                label: 'omitidos',
                value: totals.skipped,
                color: AppColors.textMuted,
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip({required this.label, required this.value, required this.color});

  final String label;
  final int value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
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
              fontSize: 16,
              fontWeight: FontWeight.w700,
              color: color,
            ),
          ),
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 10.5,
              fontWeight: FontWeight.w500,
              color: color,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }
}

class _StockCard extends StatelessWidget {
  const _StockCard({required this.row});

  final InventoryCompareRow row;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final pixels = ImageMemory.decodePixels(context, 36);
    final delta = row.delta;

    return AppCard(
      padding: const EdgeInsets.all(11),
      child: Row(
        children: [
          ClipRRect(
            borderRadius: AppRadius.smAll,
            child: SizedBox(
              width: 36,
              height: 36,
              child: row.imageUrl == null || row.imageUrl!.isEmpty
                  ? Container(
                      color: AppColors.surfaceMuted,
                      alignment: Alignment.center,
                      child: const Icon(Icons.inventory_2_outlined,
                          size: 16, color: AppColors.textDisabled),
                    )
                  : Image.network(
                      row.imageUrl!,
                      fit: BoxFit.cover,
                      cacheWidth: pixels,
                      cacheHeight: pixels,
                      errorBuilder: (context, error, stack) => Container(
                        color: AppColors.surfaceMuted,
                        alignment: Alignment.center,
                        child: const Icon(Icons.broken_image_outlined,
                            size: 16, color: AppColors.textDisabled),
                      ),
                    ),
            ),
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  row.name.isEmpty ? row.sku : row.name,
                  style: theme.textTheme.titleSmall?.copyWith(fontSize: 13),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  row.reason != null && row.reason!.isNotEmpty
                      ? row.reason!
                      : 'SKU ${row.sku}',
                  style: theme.textTheme.labelSmall,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Row(
                children: [
                  _Qty(label: 'aqui', value: row.probabilityQty),
                  const SizedBox(width: 10),
                  _Qty(
                    label: 'canal',
                    value: row.channelQty,
                    highlight: row.needsUpdate,
                  ),
                ],
              ),
              if (delta != null && delta != 0) ...[
                const SizedBox(height: 3),
                Text(
                  delta > 0 ? '+$delta' : '$delta',
                  style: theme.textTheme.labelSmall?.copyWith(
                    fontWeight: FontWeight.w700,
                    color: delta > 0 ? AppColors.success : AppColors.error,
                  ),
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }
}

class _Qty extends StatelessWidget {
  const _Qty({required this.label, required this.value, this.highlight = false});

  final String label;
  final int? value;
  final bool highlight;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      children: [
        Text(
          value == null ? '-' : AppFormat.number(value),
          style: theme.textTheme.titleSmall?.copyWith(
            fontSize: 14,
            color: highlight ? AppColors.warning : null,
          ),
        ),
        Text(label, style: theme.textTheme.labelSmall?.copyWith(fontSize: 9.5)),
      ],
    );
  }
}

class _PushBar extends StatelessWidget {
  const _PushBar({required this.count, required this.onPush});

  final int count;
  final VoidCallback onPush;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.fromLTRB(16, 10, 16, 14),
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(top: BorderSide(color: AppColors.border)),
      ),
      child: SafeArea(
        top: false,
        child: Row(
          children: [
            Expanded(
              child: Text(
                '$count producto${count == 1 ? '' : 's'} quedaria distinto',
                style: Theme.of(context).textTheme.labelMedium,
              ),
            ),
            const SizedBox(width: 10),
            FilledButton.icon(
              onPressed: onPush,
              icon: const Icon(Icons.upload_rounded, size: 17),
              label: const Text('Enviar stock'),
            ),
          ],
        ),
      ),
    );
  }
}
