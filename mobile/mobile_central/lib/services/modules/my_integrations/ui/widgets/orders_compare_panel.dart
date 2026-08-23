import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../app/sync_providers.dart';
import '../../domain/entities.dart';
import '../../domain/orders_compare_entities.dart';
import '../providers/orders_compare_provider.dart';

class OrdersComparePanel extends StatefulWidget {
  const OrdersComparePanel({
    super.key,
    required this.integrations,
    this.businessId,
  });

  final List<MyIntegration> integrations;
  final int? businessId;

  @override
  State<OrdersComparePanel> createState() => _OrdersComparePanelState();
}

class _OrdersComparePanelState extends State<OrdersComparePanel> {
  final TextEditingController _searchController = TextEditingController();

  List<MyIntegration> get _channels => widget.integrations
      .where((i) => ordersCompareTypeIds.contains(i.integrationTypeId))
      .toList();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _prepare());
  }

  @override
  void didUpdateWidget(OrdersComparePanel oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _prepare();
  }

  void _prepare() {
    if (!mounted) return;
    final provider = context.read<OrdersCompareProvider>();
    provider.configure(businessId: widget.businessId);

    final channels = _channels;
    if (channels.isEmpty) return;
    if (provider.integrationId == null) {
      provider.selectChannel(channels.first.id);
    }
    if (provider.from.isEmpty) {
      final now = DateTime.now();
      provider.setRange(
        from: _iso(now.subtract(const Duration(days: 30))),
        to: _iso(now),
      );
    }
  }

  String _iso(DateTime value) =>
      '${value.year.toString().padLeft(4, '0')}-'
      '${value.month.toString().padLeft(2, '0')}-'
      '${value.day.toString().padLeft(2, '0')}';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _pickDate({required bool isFrom}) async {
    final provider = context.read<OrdersCompareProvider>();
    final current = DateTime.tryParse(isFrom ? provider.from : provider.to);
    final picked = await showDatePicker(
      context: context,
      initialDate: current ?? DateTime.now(),
      firstDate: DateTime(2018),
      lastDate: DateTime.now().add(const Duration(days: 1)),
    );
    if (picked == null) return;
    if (isFrom) {
      provider.setRange(from: _iso(picked));
    } else {
      provider.setRange(to: _iso(picked));
    }
  }

  Future<void> _apply() async {
    final provider = context.read<OrdersCompareProvider>();
    final messenger = ScaffoldMessenger.of(context);
    final count = provider.selection.length;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Crear en Probability'),
        content: Text(
          'Se van a crear $count orden${count == 1 ? '' : 'es'} tomando los '
          'datos del canal. Es una accion que escribe: no se puede deshacer '
          'desde aqui.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Crear'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    final ok = await provider.apply();
    if (!mounted) return;

    messenger.showSnackBar(SnackBar(
      content: Text(ok
          ? provider.lastApply?.summary ?? 'Ordenes enviadas a crear'
          : provider.error ?? 'No se pudieron crear las ordenes'),
    ));

    if (ok) {
      await Future<void>.delayed(const Duration(milliseconds: 2500));
      if (mounted) await provider.compare();
    }
  }

  @override
  Widget build(BuildContext context) {
    final channels = _channels;
    if (channels.isEmpty) {
      return const AppEmptyState(
        icon: Icons.link_off_rounded,
        title: 'Sin canales para comparar',
        message:
            'Ninguno de tus canales conectados permite comparar ordenes todavia.',
      );
    }

    return Consumer<OrdersCompareProvider>(
      builder: (context, provider, _) {
        return Column(
          children: [
            _Controls(
              channels: channels,
              provider: provider,
              searchController: _searchController,
              onPickDate: _pickDate,
            ),
            Expanded(
              child: !provider.started
                  ? const AppEmptyState(
                      icon: Icons.fact_check_outlined,
                      title: 'Compara con el canal',
                      message:
                          'Elige el canal y el rango de fechas, y toca Comparar '
                          'para preguntarle sus ordenes.',
                    )
                  : PaginatedListView<OrderCompareRow>(
                      controller: provider.rows,
                      unitLabel: 'ordenes',
                      emptyIcon: Icons.check_circle_outline_rounded,
                      emptyTitle: 'Sin diferencias',
                      emptyMessage:
                          'No hay diferencias en ese rango de fechas.',
                      placeholderHeight: 116,
                      onRefresh: provider.compare,
                      header: _Summary(provider: provider),
                      itemBuilder: (context, row, index) => _OrderCard(
                        row: row,
                        selected: provider.isSelected(row.externalId),
                        onToggle: () => provider.toggle(row.externalId),
                      ),
                    ),
            ),
            if (provider.hasSelection)
              _ApplyBar(provider: provider, onApply: _apply),
          ],
        );
      },
    );
  }
}

class _Controls extends StatelessWidget {
  const _Controls({
    required this.channels,
    required this.provider,
    required this.searchController,
    required this.onPickDate,
  });

  final List<MyIntegration> channels;
  final OrdersCompareProvider provider;
  final TextEditingController searchController;
  final Future<void> Function({required bool isFrom}) onPickDate;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 10, 16, 8),
      child: Column(
        children: [
          Row(
            children: [
              Expanded(
                child: DropdownButtonFormField<int>(
                  initialValue: provider.integrationId,
                  isDense: true,
                  decoration: const InputDecoration(
                    labelText: 'Canal',
                    isDense: true,
                  ),
                  items: channels
                      .map((channel) => DropdownMenuItem(
                            value: channel.id,
                            child: Text(
                              channel.integrationTypeName ?? channel.name,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ))
                      .toList(),
                  onChanged: (value) {
                    if (value != null) provider.selectChannel(value);
                  },
                ),
              ),
              const SizedBox(width: 8),
              FilledButton.icon(
                onPressed: provider.rows.isLoading ? null : provider.compare,
                icon: provider.rows.isLoading
                    ? const SizedBox(
                        width: 14,
                        height: 14,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation(Colors.white),
                        ),
                      )
                    : const Icon(Icons.refresh_rounded, size: 17),
                label: const Text('Comparar'),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(
                child: _DateField(
                  label: 'Desde',
                  value: provider.from,
                  onTap: () => onPickDate(isFrom: true),
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: _DateField(
                  label: 'Hasta',
                  value: provider.to,
                  onTap: () => onPickDate(isFrom: false),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: searchController,
                  textInputAction: TextInputAction.search,
                  decoration: const InputDecoration(
                    isDense: true,
                    hintText: 'Numero de orden o cliente',
                    prefixIcon: Icon(Icons.search_rounded, size: 19),
                  ),
                  onChanged: provider.setSearch,
                  onSubmitted: (_) => provider.compare(),
                ),
              ),
              const SizedBox(width: 8),
              _DiffToggle(provider: provider),
            ],
          ),
          if (provider.checkedAt != null) ...[
            const SizedBox(height: 7),
            Align(
              alignment: Alignment.centerLeft,
              child: Text(
                'Consultado ${AppFormat.relative(AppFormat.parseDate(provider.checkedAt!))}',
                style: theme.textTheme.labelSmall,
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _DateField extends StatelessWidget {
  const _DateField({
    required this.label,
    required this.value,
    required this.onTap,
  });

  final String label;
  final String value;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: AppRadius.smAll,
      child: InputDecorator(
        decoration: InputDecoration(labelText: label, isDense: true),
        child: Row(
          children: [
            Expanded(
              child: Text(
                value.isEmpty ? 'Sin fecha' : value,
                style: Theme.of(context).textTheme.bodyMedium,
              ),
            ),
            const Icon(Icons.event_rounded, size: 16, color: AppColors.textMuted),
          ],
        ),
      ),
    );
  }
}

class _DiffToggle extends StatelessWidget {
  const _DiffToggle({required this.provider});

  final OrdersCompareProvider provider;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final active = provider.onlyDiff;

    return Tooltip(
      message: 'Ocultar las ordenes que ya estan iguales en los dos lados',
      child: GestureDetector(
        onTap: provider.toggleOnlyDiff,
        behavior: HitTestBehavior.opaque,
        child: Container(
          height: 44,
          padding: const EdgeInsets.symmetric(horizontal: 12),
          alignment: Alignment.center,
          decoration: BoxDecoration(
            color: active ? scheme.primaryContainer : AppColors.surface,
            borderRadius: AppRadius.smAll,
            border: Border.all(
              color: active ? scheme.primary : AppColors.border,
            ),
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
      ),
    );
  }
}

class _Summary extends StatelessWidget {
  const _Summary({required this.provider});

  final OrdersCompareProvider provider;

  @override
  Widget build(BuildContext context) {
    final totals = provider.totals;
    final withoutInventory = totals.withoutInventory;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          height: 74,
          child: ListView(
            scrollDirection: Axis.horizontal,
            children: [
              _Kpi(label: 'En el canal', value: totals.inChannel),
              _Kpi(
                label: 'Faltan aqui',
                value: totals.toCreate,
                color: AppColors.info,
              ),
              _Kpi(
                label: 'En las dos',
                value: totals.inSync,
                color: AppColors.success,
              ),
              _Kpi(
                label: 'Estado distinto',
                value: totals.withStatusMismatch,
                color: AppColors.warning,
              ),
              _Kpi(label: 'Sin mover stock', value: withoutInventory),
            ],
          ),
        ),
        if (withoutInventory > 0) ...[
          const SizedBox(height: 8),
          _InventoryNotice(count: withoutInventory),
        ],
        if (provider.rows.loadedItems.any((row) => row.canCreate)) ...[
          const SizedBox(height: 6),
          Align(
            alignment: Alignment.centerLeft,
            child: TextButton.icon(
              onPressed: provider.toggleAllLoaded,
              icon: const Icon(Icons.checklist_rounded, size: 17),
              label: const Text('Marcar las que faltan'),
            ),
          ),
        ],
        const SizedBox(height: 4),
      ],
    );
  }
}

class _Kpi extends StatelessWidget {
  const _Kpi({required this.label, required this.value, this.color});

  final String label;
  final int value;
  final Color? color;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      width: 116,
      margin: const EdgeInsets.only(right: 8),
      padding: const EdgeInsets.symmetric(horizontal: 11, vertical: 9),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: AppRadius.smAll,
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            label.toUpperCase(),
            style: theme.textTheme.labelSmall?.copyWith(
              fontSize: 9.5,
              letterSpacing: 0.8,
              fontWeight: FontWeight.w700,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 5),
          Text(
            AppFormat.number(value),
            style: theme.textTheme.titleMedium?.copyWith(
              color: color,
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}

class _InventoryNotice extends StatelessWidget {
  const _InventoryNotice({required this.count});

  final int count;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 11, vertical: 9),
      decoration: BoxDecoration(
        color: AppColors.warning.withValues(alpha: 0.10),
        borderRadius: AppRadius.smAll,
        border: Border.all(color: AppColors.warning.withValues(alpha: 0.35)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.inventory_2_outlined, size: 16, color: AppColors.warning),
          const SizedBox(width: 9),
          Expanded(
            child: Text(
              '$count de las ordenes por crear entran como historicas: el canal '
              'ya las dio por entregadas, despachadas, canceladas o devueltas, '
              'asi que se crean completas pero no mueven inventario.',
              style: Theme.of(context).textTheme.labelSmall,
            ),
          ),
        ],
      ),
    );
  }
}

class _OrderCard extends StatelessWidget {
  const _OrderCard({
    required this.row,
    required this.selected,
    required this.onToggle,
  });

  final OrderCompareRow row;
  final bool selected;
  final VoidCallback onToggle;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return AppCard(
      padding: const EdgeInsets.fromLTRB(12, 11, 12, 11),
      onTap: row.canCreate ? onToggle : null,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (row.canCreate)
                SizedBox(
                  width: 26,
                  height: 26,
                  child: Checkbox(
                    value: selected,
                    onChanged: (_) => onToggle(),
                    visualDensity: VisualDensity.compact,
                  ),
                )
              else
                const SizedBox(width: 26),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      row.label,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      row.customerName.isEmpty ? 'Sin cliente' : row.customerName,
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
                  Text(
                    AppFormat.money(row.amount),
                    style: theme.textTheme.titleSmall?.copyWith(fontSize: 14),
                  ),
                  if (row.createdOn != null) ...[
                    const SizedBox(height: 2),
                    Text(
                      AppFormat.date(row.createdOn),
                      style: theme.textTheme.labelSmall,
                    ),
                  ],
                ],
              ),
            ],
          ),
          const SizedBox(height: 9),
          Row(
            children: [
              Expanded(
                child: _Side(
                  label: 'En el canal',
                  value: row.channelStatus.isEmpty ? '-' : row.channelStatus,
                ),
              ),
              const Icon(Icons.arrow_forward_rounded,
                  size: 14, color: AppColors.textDisabled),
              Expanded(
                child: _Side(
                  label: 'Aqui',
                  value: row.localStatus ?? 'no existe',
                  muted: row.localStatus == null,
                ),
              ),
            ],
          ),
          const SizedBox(height: 9),
          Wrap(
            spacing: 6,
            runSpacing: 6,
            children: [
              _situation(row),
              if (row.canCreate && !row.movesInventory)
                const AppStatusChip(
                  dense: true,
                  label: 'no mueve stock',
                  tone: AppStatusTone.neutral,
                ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _situation(OrderCompareRow row) {
    if (row.action == OrderCompareAction.create) {
      return const AppStatusChip(
        dense: true,
        label: 'falta en Probability',
        tone: AppStatusTone.info,
      );
    }
    if (row.action == OrderCompareAction.onlyInProbability) {
      return const AppStatusChip(
        dense: true,
        label: 'solo en Probability',
        tone: AppStatusTone.neutral,
      );
    }
    if (row.hasMismatch) {
      return AppStatusChip(
        dense: true,
        label: row.statusMismatch ? 'estado distinto' : 'total distinto',
        tone: AppStatusTone.warning,
      );
    }
    return const AppStatusChip(
      dense: true,
      label: 'en las dos',
      tone: AppStatusTone.success,
    );
  }
}

class _Side extends StatelessWidget {
  const _Side({required this.label, required this.value, this.muted = false});

  final String label;
  final String value;
  final bool muted;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label.toUpperCase(), style: theme.textTheme.labelSmall?.copyWith(fontSize: 9)),
        const SizedBox(height: 2),
        Text(
          value,
          style: theme.textTheme.bodySmall?.copyWith(
            color: muted ? AppColors.textDisabled : AppColors.textSecondary,
            fontStyle: muted ? FontStyle.italic : FontStyle.normal,
          ),
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
      ],
    );
  }
}

class _ApplyBar extends StatelessWidget {
  const _ApplyBar({required this.provider, required this.onApply});

  final OrdersCompareProvider provider;
  final VoidCallback onApply;

  @override
  Widget build(BuildContext context) {
    final count = provider.selection.length;

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 10, 16, 14),
      decoration: BoxDecoration(
        color: AppColors.surface,
        border: Border(top: BorderSide(color: AppColors.border)),
      ),
      child: SafeArea(
        top: false,
        child: Row(
          children: [
            Expanded(
              child: Text(
                '$count seleccionada${count == 1 ? '' : 's'} para crear aqui',
                style: Theme.of(context).textTheme.labelMedium,
              ),
            ),
            const SizedBox(width: 10),
            FilledButton.icon(
              onPressed: provider.applying ? null : onApply,
              icon: provider.applying
                  ? const SizedBox(
                      width: 14,
                      height: 14,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        valueColor: AlwaysStoppedAnimation(Colors.white),
                      ),
                    )
                  : const Icon(Icons.add_task_rounded, size: 17),
              label: const Text('Crear'),
            ),
          ],
        ),
      ),
    );
  }
}
