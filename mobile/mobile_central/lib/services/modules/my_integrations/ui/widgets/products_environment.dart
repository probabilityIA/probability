import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/filters/filter_models.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/theme/channel_brand.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/image_memory.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/comparison_lists_provider.dart';
import '../providers/saved_comparison_provider.dart';
import '../providers/sync_activity_provider.dart';
import '../screens/finding_items_screen.dart';

class ProductsEnvironment extends StatefulWidget {
  const ProductsEnvironment({super.key, this.businessId});

  final int? businessId;

  @override
  State<ProductsEnvironment> createState() => _ProductsEnvironmentState();
}

class _ProductsEnvironmentState extends State<ProductsEnvironment> {
  final TextEditingController _search = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void didUpdateWidget(ProductsEnvironment oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _load();
  }

  void _load() {
    if (!mounted) return;
    final saved = context.read<SavedComparisonProvider>();
    saved.configure(businessId: widget.businessId);
    saved.loadFindings();
    context.read<ProductMatrixProvider>().load(businessId: widget.businessId);
  }

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  Future<void> _compare() async {
    final sync = context.read<SyncActivityProvider>();
    final saved = context.read<SavedComparisonProvider>();
    final matrix = context.read<ProductMatrixProvider>();
    final messenger = ScaffoldMessenger.of(context);

    if (sync.running) return;
    if (sync.eligible.isEmpty) {
      messenger.showSnackBar(const SnackBar(
        content: Text('Ningun canal conectado permite comparar productos'),
      ));
      return;
    }

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Comparar de nuevo'),
        content: Text(
          'Se le va a preguntar su catalogo a ${sync.eligible.length} canal'
          '${sync.eligible.length == 1 ? '' : 'es'}. Puede tardar varios '
          'minutos y reemplaza la comparacion guardada.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Comparar'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    await sync.runProducts();
    if (!mounted) return;

    await saved.loadFindings(force: true);
    await matrix.load(businessId: widget.businessId);
    if (!mounted) return;

    messenger.showSnackBar(
      const SnackBar(content: Text('Comparacion actualizada')),
    );
  }

  void _openFinding(Finding finding) {
    Navigator.of(context).push(MaterialPageRoute(
      builder: (_) => FindingItemsScreen(
        finding: finding,
        businessId: widget.businessId,
      ),
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Consumer2<ProductMatrixProvider, SavedComparisonProvider>(
      builder: (context, matrix, saved, _) {
        return Column(
          children: [
            _Filters(
              matrix: matrix,
              saved: saved,
              controller: _search,
              onOpenFinding: _openFinding,
              onCompare: _compare,
            ),
            Expanded(
              child: PaginatedListView<MatrixRow>(
                controller: matrix.list,
                unitLabel: 'productos',
                emptyIcon: Icons.grid_view_rounded,
                emptyTitle: matrix.hasFilters ? 'Sin coincidencias' : 'Sin productos',
                emptyMessage: matrix.hasFilters
                    ? 'Ningun producto cumple con los filtros de canal.'
                    : 'No hay una comparacion de productos guardada. Corre la '
                        'comparacion para llenarla.',
                placeholderHeight: 104,
                onRefresh: () async {
                  await saved.loadFindings(force: true);
                  await matrix.load(businessId: widget.businessId);
                },
                itemBuilder: (context, row, index) =>
                    _ProductCard(row: row, columns: matrix.columns),
              ),
            ),
          ],
        );
      },
    );
  }
}

class _Filters extends StatelessWidget {
  const _Filters({
    required this.matrix,
    required this.saved,
    required this.controller,
    required this.onOpenFinding,
    required this.onCompare,
  });

  final ProductMatrixProvider matrix;
  final SavedComparisonProvider saved;
  final TextEditingController controller;
  final void Function(Finding) onOpenFinding;
  final Future<void> Function() onCompare;

  List<FilterDimension> get _dimensions => matrix.columns
      .map((column) => FilterDimension(
            key: 'ch_${column.integrationId}',
            label: column.name,
            icon: Icons.storefront_outlined,
            imageUrl: column.imageUrl,
            accent: channelBrand(column.name.isEmpty ? column.code : column.name)
                .color,
            options: [
              FilterOption(value: 'present', label: 'Se vende en ${column.name}'),
              FilterOption(
                value: 'missing',
                label: 'No se vende en ${column.name}',
              ),
            ],
          ))
      .toList();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final findings = saved.findings.findings;
    final comparedAt = saved.findings.lastComparedAt;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        AppFilterBar(
          controller: controller,
          searchFields: const [
            SearchField(
              key: 'all',
              label: 'Todo',
              hint: 'SKU, producto o codigo de barras',
            ),
            SearchField(key: 'sku', label: 'SKU', hint: 'Buscar por SKU'),
            SearchField(
              key: 'name',
              label: 'Producto',
              hint: 'Buscar por nombre del producto',
            ),
            SearchField(
              key: 'barcode',
              label: 'Codigo de barras',
              hint: 'Buscar por codigo de barras',
            ),
          ],
          selectedField: matrix.searchBy,
          onFieldChanged: matrix.setSearchBy,
          onSearchChanged: matrix.setSearch,
          dimensions: _dimensions,
          selection: matrix.selection,
          onSelectionChanged: matrix.applySelection,
          trailing: _CompareButton(onCompare: onCompare),
        ),
        if (findings.isNotEmpty)
          SizedBox(
            height: 30,
            child: ListView.separated(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: findings.length,
              separatorBuilder: (context, index) => const SizedBox(width: 7),
              itemBuilder: (context, index) => _FindingChip(
                finding: findings[index],
                onTap: () => onOpenFinding(findings[index]),
              ),
            ),
          ),
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 9, 16, 7),
          child: Row(
            children: [
              Expanded(
                child: Text(
                  _summary(matrix),
                  style: theme.textTheme.labelSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
              if (comparedAt != null)
                Text(
                  'comparado ${AppFormat.relative(comparedAt)}',
                  style: theme.textTheme.labelSmall,
                ),
            ],
          ),
        ),
      ],
    );
  }

  String _summary(ProductMatrixProvider matrix) {
    final total = matrix.list.total;
    if (matrix.list.isLoading && total == 0) return 'Cargando productos';
    if (!matrix.hasFilters) return '${AppFormat.number(total)} productos';
    return '${AppFormat.number(total)} productos con el filtro';
  }
}

class _CompareButton extends StatelessWidget {
  const _CompareButton({required this.onCompare});

  final Future<void> Function() onCompare;

  @override
  Widget build(BuildContext context) {
    return Consumer<SyncActivityProvider>(
      builder: (context, sync, _) {
        final scheme = Theme.of(context).colorScheme;
        final running = sync.running;

        return Tooltip(
          message: running
              ? 'Comparando con los canales'
              : 'Volver a comparar con los canales',
          child: IconButton(
            onPressed: running ? null : onCompare,
            visualDensity: VisualDensity.compact,
            icon: running
                ? SizedBox(
                    width: 17,
                    height: 17,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      valueColor: AlwaysStoppedAnimation(scheme.primary),
                    ),
                  )
                : const Icon(Icons.sync_rounded, size: 21),
            color: AppColors.textMuted,
          ),
        );
      },
    );
  }
}

class _FindingChip extends StatelessWidget {
  const _FindingChip({required this.finding, required this.onTap});

  final Finding finding;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final tone = switch (finding.severity) {
      FindingSeverity.error => AppColors.error,
      FindingSeverity.warn => AppColors.warning,
      FindingSeverity.info => AppColors.info,
    };

    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 10),
        decoration: BoxDecoration(
          color: tone.withValues(alpha: 0.10),
          borderRadius: AppRadius.pillAll,
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              AppFormat.number(finding.count),
              style: TextStyle(
                fontFamily: 'Inter',
                fontSize: 12,
                fontWeight: FontWeight.w700,
                color: tone,
              ),
            ),
            const SizedBox(width: 6),
            Text(
              finding.title,
              style: TextStyle(
                fontFamily: 'Inter',
                fontSize: 11.5,
                fontWeight: FontWeight.w600,
                color: tone,
              ),
            ),
            const SizedBox(width: 3),
            Icon(Icons.chevron_right_rounded, size: 15, color: tone),
          ],
        ),
      ),
    );
  }
}

class _ProductCard extends StatelessWidget {
  const _ProductCard({required this.row, required this.columns});

  final MatrixRow row;
  final List<MatrixColumn> columns;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final pixels = ImageMemory.decodePixels(context, 40);

    return AppCard(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              ClipRRect(
                borderRadius: AppRadius.smAll,
                child: SizedBox(
                  width: 40,
                  height: 40,
                  child: row.imageUrl == null || row.imageUrl!.isEmpty
                      ? Container(
                          color: AppColors.surfaceMuted,
                          alignment: Alignment.center,
                          child: const Icon(Icons.inventory_2_outlined,
                              size: 18, color: AppColors.textDisabled),
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
                                size: 18, color: AppColors.textDisabled),
                          ),
                        ),
                ),
              ),
              const SizedBox(width: 11),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      row.name == null || row.name!.isEmpty ? row.sku : row.name!,
                      style: theme.textTheme.titleSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'SKU ${row.sku.isEmpty ? "-" : row.sku}',
                      style: theme.textTheme.labelSmall,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              Text(
                '${row.presentCount}/${columns.length}',
                style: theme.textTheme.labelMedium?.copyWith(
                  fontWeight: FontWeight.w700,
                  color: row.presentCount == columns.length
                      ? AppColors.success
                      : AppColors.textMuted,
                ),
              ),
            ],
          ),
          if (columns.isNotEmpty) ...[
            const SizedBox(height: 10),
            Wrap(
              spacing: 6,
              runSpacing: 6,
              children: columns
                  .map((column) => _ChannelTag(
                        column: column,
                        cell: row.cellFor(column.integrationId),
                      ))
                  .toList(),
            ),
          ],
        ],
      ),
    );
  }
}

class _ChannelTag extends StatelessWidget {
  const _ChannelTag({required this.column, required this.cell});

  final MatrixColumn column;
  final MatrixCell? cell;

  @override
  Widget build(BuildContext context) {
    final brand = channelBrand(column.name.isEmpty ? column.code : column.name);
    final present = cell?.present == true;
    final mismatched = cell?.mismatched == true;

    final color = !present
        ? AppColors.textDisabled
        : mismatched
            ? AppColors.warning
            : brand.color;

    return Tooltip(
      message: !present
          ? 'No esta publicado en ${column.name}'
          : mismatched
              ? 'Publicado en ${column.name} con el SKU ${cell?.sku ?? "distinto"}'
              : 'Publicado en ${column.name}',
      child: Opacity(
        opacity: present ? 1 : 0.42,
        child: Container(
          padding: const EdgeInsets.fromLTRB(4, 4, 8, 4),
          decoration: BoxDecoration(
            color: present ? color.withValues(alpha: 0.10) : Colors.transparent,
            borderRadius: AppRadius.pillAll,
            border: Border.all(
              color: present ? color.withValues(alpha: 0.45) : AppColors.border,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              BrandLogo(
                name: column.name,
                imageUrl: column.imageUrl,
                size: 22,
                radius: 999,
                padding: 3,
              ),
              const SizedBox(width: 6),
              Icon(
                !present
                    ? Icons.remove_rounded
                    : mismatched
                        ? Icons.error_outline_rounded
                        : Icons.check_rounded,
                size: 14,
                color: color,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
