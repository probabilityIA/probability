import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../../../shared/filters/filter_models.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/theme/app_tokens.dart';
import '../../../../../shared/theme/channel_brand.dart';
import '../../../../../shared/utils/formatters.dart';
import '../../../../../shared/utils/image_memory.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../domain/entities.dart';
import '../../domain/saved_comparison_entities.dart';
import '../providers/comparison_lists_provider.dart';

class InventoryEnvironment extends StatefulWidget {
  const InventoryEnvironment({
    super.key,
    required this.integrations,
    this.businessId,
  });

  final List<MyIntegration> integrations;
  final int? businessId;

  @override
  State<InventoryEnvironment> createState() => _InventoryEnvironmentState();
}

class _InventoryEnvironmentState extends State<InventoryEnvironment> {
  final TextEditingController _search = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void didUpdateWidget(InventoryEnvironment oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) _load();
  }

  void _load() {
    if (!mounted) return;
    context.read<InventoryMatrixProvider>().load(
          integrations: widget.integrations,
          businessId: widget.businessId,
        );
  }

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  Future<void> _ask() async {
    final provider = context.read<InventoryMatrixProvider>();
    final messenger = ScaffoldMessenger.of(context);

    if (provider.channels.isEmpty) {
      messenger.showSnackBar(const SnackBar(
        content: Text('Ningun canal conectado permite comparar inventario'),
      ));
      return;
    }

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Preguntar el stock'),
        content: Text(
          'Se le va a preguntar su stock a ${provider.channels.length} canal'
          '${provider.channels.length == 1 ? '' : 'es'} para los productos de '
          'esta pagina. Solo lee, no cambia nada en el canal.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Preguntar'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    await provider.askChannels();
    if (!mounted) return;
    messenger.showSnackBar(
      const SnackBar(content: Text('Stock actualizado desde los canales')),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<InventoryMatrixProvider>(
      builder: (context, provider, _) {
        if (provider.channels.isEmpty) {
          return const AppEmptyState(
            icon: Icons.link_off_rounded,
            title: 'Sin canales para comparar',
            message:
                'Ninguno de tus canales conectados permite comparar inventario '
                'todavia.',
          );
        }

        return Column(
          children: [
            _Filters(
              provider: provider,
              controller: _search,
              onAsk: _ask,
            ),
            Expanded(
              child: PaginatedListView<MatrixRow>(
                controller: provider.list,
                unitLabel: 'productos',
                emptyIcon: Icons.inventory_2_outlined,
                emptyTitle: provider.hasFilters
                    ? 'Sin coincidencias'
                    : 'Sin comparacion',
                emptyMessage: provider.hasFilters
                    ? 'Ningun producto cumple con los filtros de canal.'
                    : 'No hay una comparacion guardada. Toca el boton para '
                        'preguntarle el stock a tus canales.',
                placeholderHeight: 112,
                onRefresh: () async => _load(),
                itemBuilder: (context, row, index) => _StockCard(
                  row: row,
                  columns: provider.columns,
                  provider: provider,
                ),
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
    required this.provider,
    required this.controller,
    required this.onAsk,
  });

  final InventoryMatrixProvider provider;
  final TextEditingController controller;
  final Future<void> Function() onAsk;

  List<FilterDimension> get _dimensions => provider.columns
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
    final checked = provider.lastCheck;

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
          selectedField: provider.searchBy,
          onFieldChanged: provider.setSearchBy,
          onSearchChanged: provider.setSearch,
          dimensions: _dimensions,
          selection: provider.selection,
          onSelectionChanged: provider.applySelection,
          trailing: _AskButton(asking: provider.asking, onAsk: onAsk),
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 2, 16, 7),
          child: Row(
            children: [
              Expanded(
                child: Text(
                  _summary(provider),
                  style: theme.textTheme.labelSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
              if (checked != null)
                Text(
                  'stock leido ${AppFormat.relative(checked)}',
                  style: theme.textTheme.labelSmall,
                ),
            ],
          ),
        ),
      ],
    );
  }

  String _summary(InventoryMatrixProvider provider) {
    final total = provider.list.total;
    if (provider.list.isLoading && total == 0) return 'Cargando productos';
    if (!provider.hasFilters) return '${AppFormat.number(total)} productos';
    return '${AppFormat.number(total)} productos con el filtro';
  }
}

class _AskButton extends StatelessWidget {
  const _AskButton({required this.asking, required this.onAsk});

  final bool asking;
  final Future<void> Function() onAsk;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;

    return Tooltip(
      message: asking
          ? 'Preguntando el stock a los canales'
          : 'Preguntarle el stock a los canales',
      child: IconButton(
        onPressed: asking ? null : onAsk,
        visualDensity: VisualDensity.compact,
        style: IconButton.styleFrom(
          backgroundColor: scheme.primaryContainer,
          foregroundColor: scheme.primary,
          shape: const CircleBorder(),
        ),
        icon: asking
            ? SizedBox(
                width: 17,
                height: 17,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation(scheme.primary),
                ),
              )
            : const Icon(Icons.sync_rounded, size: 20),
      ),
    );
  }
}

class _StockCard extends StatelessWidget {
  const _StockCard({
    required this.row,
    required this.columns,
    required this.provider,
  });

  final MatrixRow row;
  final List<MatrixColumn> columns;
  final InventoryMatrixProvider provider;

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
            ],
          ),
          if (columns.isNotEmpty) ...[
            const SizedBox(height: 10),
            Wrap(
              spacing: 6,
              runSpacing: 6,
              children: columns
                  .map((column) => _StockTag(
                        column: column,
                        cell: row.cellFor(column.integrationId),
                        stock: provider.stockFor(column.integrationId, row.sku),
                      ))
                  .toList(),
            ),
          ],
        ],
      ),
    );
  }
}

class _StockTag extends StatelessWidget {
  const _StockTag({
    required this.column,
    required this.cell,
    required this.stock,
  });

  final MatrixColumn column;
  final MatrixCell? cell;
  final InventoryCompareRow? stock;

  @override
  Widget build(BuildContext context) {
    final brand = channelBrand(column.name.isEmpty ? column.code : column.name);
    final present = cell?.present == true;

    if (!present) {
      return _Tag(
        column: column,
        color: AppColors.textDisabled,
        faded: true,
        child: const Text(
          'no publicado',
          style: TextStyle(
            fontFamily: 'Inter',
            fontSize: 11,
            fontStyle: FontStyle.italic,
            color: AppColors.textDisabled,
          ),
        ),
      );
    }

    if (stock == null) {
      return _Tag(
        column: column,
        color: brand.color,
        faded: true,
        child: const Text(
          'sin leer',
          style: TextStyle(
            fontFamily: 'Inter',
            fontSize: 11,
            fontStyle: FontStyle.italic,
            color: AppColors.textMuted,
          ),
        ),
      );
    }

    final delta = stock!.delta;
    final differs = stock!.needsUpdate && delta != null && delta != 0;
    final color = differs ? AppColors.warning : brand.color;

    return _Tag(
      column: column,
      color: color,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            stock!.channelQty == null
                ? '-'
                : AppFormat.number(stock!.channelQty),
            style: TextStyle(
              fontFamily: 'Inter',
              fontSize: 13,
              fontWeight: FontWeight.w700,
              color: color,
            ),
          ),
          if (differs) ...[
            const SizedBox(width: 5),
            Text(
              delta > 0 ? '+$delta' : '$delta',
              style: const TextStyle(
                fontFamily: 'Inter',
                fontSize: 11,
                fontWeight: FontWeight.w600,
                color: AppColors.warning,
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _Tag extends StatelessWidget {
  const _Tag({
    required this.column,
    required this.color,
    required this.child,
    this.faded = false,
  });

  final MatrixColumn column;
  final Color color;
  final Widget child;
  final bool faded;

  @override
  Widget build(BuildContext context) {
    return Opacity(
      opacity: faded ? 0.55 : 1,
      child: Container(
        padding: const EdgeInsets.fromLTRB(4, 4, 9, 4),
        decoration: BoxDecoration(
          color: faded ? Colors.transparent : color.withValues(alpha: 0.10),
          borderRadius: AppRadius.pillAll,
          border: Border.all(
            color: faded ? AppColors.border : color.withValues(alpha: 0.45),
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
            child,
          ],
        ),
      ),
    );
  }
}
