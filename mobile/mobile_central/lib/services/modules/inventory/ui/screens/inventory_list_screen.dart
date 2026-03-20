import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/inventory_provider.dart';
import '../../domain/entities.dart';
import '../../../warehouses/ui/providers/warehouse_provider.dart';

class InventoryListScreen extends StatefulWidget {
  final int? businessId;

  const InventoryListScreen({super.key, this.businessId});

  @override
  State<InventoryListScreen> createState() => _InventoryListScreenState();
}

class _InventoryListScreenState extends State<InventoryListScreen> {
  final _searchController = TextEditingController();
  int? _selectedWarehouseId;
  bool _lowStockOnly = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadWarehouses();
    });
  }

  @override
  void didUpdateWidget(InventoryListScreen oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.businessId != widget.businessId) {
      _searchController.clear();
      setState(() {
        _selectedWarehouseId = null;
        _lowStockOnly = false;
      });
      _loadWarehouses();
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadWarehouses() {
    context
        .read<WarehouseProvider>()
        .fetchWarehouses(businessId: widget.businessId);
  }

  void _fetchInventory() {
    if (_selectedWarehouseId == null) return;
    final provider = context.read<InventoryProvider>();
    provider.setFilters(
      search: _searchController.text.isNotEmpty ? _searchController.text : null,
      lowStock: _lowStockOnly ? true : null,
    );
    provider.fetchWarehouseInventory(_selectedWarehouseId!,
        businessId: widget.businessId);
  }

  void _onClearSearch() {
    _searchController.clear();
    final provider = context.read<InventoryProvider>();
    provider.setFilters(search: null);
    provider.setPage(1);
    if (_selectedWarehouseId != null) {
      provider.fetchWarehouseInventory(_selectedWarehouseId!,
          businessId: widget.businessId);
    }
  }

  void _onWarehouseChanged(int? warehouseId) {
    setState(() {
      _selectedWarehouseId = warehouseId;
      _lowStockOnly = false;
    });
    _searchController.clear();
    final provider = context.read<InventoryProvider>();
    provider.resetFilters();
    if (warehouseId != null) {
      provider.fetchWarehouseInventory(warehouseId,
          businessId: widget.businessId);
    }
  }

  bool _isLowStock(InventoryLevel level) {
    if (level.reorderPoint != null) return level.availableQty <= level.reorderPoint!;
    if (level.minStock != null) return level.availableQty <= level.minStock!;
    return false;
  }

  void _showAdjustStockForm({String? productId}) {
    if (_selectedWarehouseId == null) return;

    final productIdCtrl = TextEditingController(text: productId ?? '');
    final quantityCtrl = TextEditingController();
    final reasonCtrl = TextEditingController();
    final notesCtrl = TextEditingController();
    final formKey = GlobalKey<FormState>();
    bool isSaving = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) {
        return StatefulBuilder(
          builder: (ctx, setModalState) {
            return Padding(
              padding: EdgeInsets.only(
                left: 16,
                right: 16,
                top: 16,
                bottom: MediaQuery.of(ctx).viewInsets.bottom + 16,
              ),
              child: Form(
                key: formKey,
                child: SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: Text(
                              'Ajustar stock',
                              style: Theme.of(ctx).textTheme.titleLarge,
                            ),
                          ),
                          IconButton(
                            icon: const Icon(Icons.close),
                            onPressed: () => Navigator.pop(ctx),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: productIdCtrl,
                        decoration: const InputDecoration(
                          labelText: 'ID del producto *',
                          border: OutlineInputBorder(),
                          hintText: 'UUID del producto',
                        ),
                        readOnly: productId != null,
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: quantityCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Cantidad *',
                          border: OutlineInputBorder(),
                          helperText:
                              'Positivo para agregar, negativo para quitar',
                        ),
                        keyboardType: const TextInputType.numberWithOptions(
                            signed: true),
                        validator: (v) {
                          if (v == null || v.trim().isEmpty) return 'Requerido';
                          if (int.tryParse(v.trim()) == null) {
                            return 'Ingrese un numero valido';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: reasonCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Razon *',
                          border: OutlineInputBorder(),
                          hintText: 'Conteo fisico, correccion, etc.',
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty)
                            ? 'Requerido'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: notesCtrl,
                        decoration: const InputDecoration(
                          labelText: 'Notas',
                          border: OutlineInputBorder(),
                        ),
                        maxLines: 3,
                      ),
                      const SizedBox(height: 20),
                      FilledButton(
                        onPressed: isSaving
                            ? null
                            : () async {
                                if (!formKey.currentState!.validate()) return;
                                setModalState(() => isSaving = true);
                                final provider =
                                    context.read<InventoryProvider>();
                                final dto = AdjustStockDTO(
                                  productId: productIdCtrl.text.trim(),
                                  warehouseId: _selectedWarehouseId!,
                                  quantity:
                                      int.parse(quantityCtrl.text.trim()),
                                  reason: reasonCtrl.text.trim(),
                                  notes: notesCtrl.text.trim().isNotEmpty
                                      ? notesCtrl.text.trim()
                                      : null,
                                );
                                final result = await provider.adjustStock(
                                    dto,
                                    businessId: widget.businessId);
                                setModalState(() => isSaving = false);
                                if (result != null && ctx.mounted) {
                                  Navigator.pop(ctx);
                                  provider.fetchWarehouseInventory(
                                      _selectedWarehouseId!,
                                      businessId: widget.businessId);
                                  if (mounted) {
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      const SnackBar(
                                          content:
                                              Text('Stock ajustado')),
                                    );
                                  }
                                } else if (mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    SnackBar(
                                      content: Text(provider.error ??
                                          'Error al ajustar stock'),
                                      backgroundColor:
                                          Theme.of(context).colorScheme.error,
                                    ),
                                  );
                                }
                              },
                        child: isSaving
                            ? const SizedBox(
                                height: 20,
                                width: 20,
                                child: CircularProgressIndicator(
                                    strokeWidth: 2, color: Colors.white),
                              )
                            : const Text('Ajustar stock'),
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Inventario')),
      floatingActionButton: _selectedWarehouseId != null
          ? FloatingActionButton(
              onPressed: () => _showAdjustStockForm(),
              child: const Icon(Icons.tune),
            )
          : null,
      body: Column(
        children: [
          // Warehouse selector
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 4),
            child: Consumer<WarehouseProvider>(
              builder: (context, whProvider, _) {
                if (whProvider.isLoading && whProvider.warehouses.isEmpty) {
                  return const LinearProgressIndicator();
                }
                return DropdownButtonFormField<int>(
                  initialValue: _selectedWarehouseId,
                  decoration: const InputDecoration(
                    labelText: 'Seleccionar bodega',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.warehouse),
                  ),
                  hint: const Text('Selecciona una bodega'),
                  items: whProvider.warehouses
                      .map((wh) => DropdownMenuItem(
                            value: wh.id,
                            child: Text('${wh.name} (${wh.code})'),
                          ))
                      .toList(),
                  onChanged: (v) => _onWarehouseChanged(v),
                );
              },
            ),
          ),

          if (_selectedWarehouseId != null) ...[
            // Search bar + low stock filter
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 4),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _searchController,
                      decoration: InputDecoration(
                        hintText: 'Buscar por producto o SKU...',
                        prefixIcon: const Icon(Icons.search),
                        suffixIcon: _searchController.text.isNotEmpty
                            ? IconButton(
                                icon: const Icon(Icons.clear),
                                onPressed: _onClearSearch,
                              )
                            : null,
                        border: const OutlineInputBorder(),
                        contentPadding:
                            const EdgeInsets.symmetric(horizontal: 12),
                      ),
                      onSubmitted: (_) => _fetchInventory(),
                    ),
                  ),
                  const SizedBox(width: 8),
                  IconButton.filled(
                    onPressed: _fetchInventory,
                    icon: const Icon(Icons.search),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  FilterChip(
                    label: const Text('Solo stock bajo'),
                    selected: _lowStockOnly,
                    onSelected: (v) {
                      setState(() => _lowStockOnly = v);
                      final provider = context.read<InventoryProvider>();
                      provider.setFilters(lowStock: v ? true : null);
                      provider.setPage(1);
                      provider.fetchWarehouseInventory(_selectedWarehouseId!,
                          businessId: widget.businessId);
                    },
                  ),
                ],
              ),
            ),
          ],

          // Content
          Expanded(
            child: _selectedWarehouseId == null
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.inventory_2_outlined,
                            size: 48, color: theme.disabledColor),
                        const SizedBox(height: 16),
                        const Text('Selecciona una bodega para ver inventario'),
                      ],
                    ),
                  )
                : Consumer<InventoryProvider>(
                    builder: (context, provider, _) {
                      return _buildContent(provider, theme);
                    },
                  ),
          ),

          // Pagination
          if (_selectedWarehouseId != null)
            Consumer<InventoryProvider>(
              builder: (context, provider, _) {
                if (provider.pagination != null && !provider.isLoading) {
                  return _buildPagination(provider);
                }
                return const SizedBox.shrink();
              },
            ),
        ],
      ),
    );
  }

  Widget _buildContent(InventoryProvider provider, ThemeData theme) {
    if (provider.isLoading && provider.inventoryLevels.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (provider.error != null && provider.inventoryLevels.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 48, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text(provider.error!, textAlign: TextAlign.center),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () => provider.fetchWarehouseInventory(
                  _selectedWarehouseId!,
                  businessId: widget.businessId),
              icon: const Icon(Icons.refresh),
              label: const Text('Reintentar'),
            ),
          ],
        ),
      );
    }

    if (provider.inventoryLevels.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.inventory_outlined,
                size: 48, color: theme.disabledColor),
            const SizedBox(height: 16),
            const Text('No hay inventario en esta bodega'),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () => provider.fetchWarehouseInventory(
          _selectedWarehouseId!,
          businessId: widget.businessId),
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        itemCount: provider.inventoryLevels.length,
        itemBuilder: (context, index) {
          final level = provider.inventoryLevels[index];
          final lowStock = _isLowStock(level);

          return Card(
            margin: const EdgeInsets.only(bottom: 8),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Product info + adjust button
                  Row(
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              level.productName ?? level.productId,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w600, fontSize: 14),
                            ),
                            if (level.productSku != null &&
                                level.productSku!.isNotEmpty)
                              Text(
                                'SKU: ${level.productSku}',
                                style: theme.textTheme.bodySmall
                                    ?.copyWith(fontFamily: 'monospace'),
                              ),
                          ],
                        ),
                      ),
                      if (lowStock)
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: Colors.red.withValues(alpha: 0.15),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Text(
                            'Stock bajo',
                            style: TextStyle(
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              color: Colors.red,
                            ),
                          ),
                        )
                      else
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(
                            color: Colors.green.withValues(alpha: 0.15),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Text(
                            'OK',
                            style: TextStyle(
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              color: Colors.green,
                            ),
                          ),
                        ),
                      const SizedBox(width: 4),
                      IconButton(
                        icon: const Icon(Icons.tune, size: 20),
                        tooltip: 'Ajustar stock',
                        onPressed: () => _showAdjustStockForm(
                            productId: level.productId),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  // Quantities row
                  Row(
                    children: [
                      _buildQtyChip('Cantidad', level.quantity, theme),
                      const SizedBox(width: 8),
                      _buildQtyChip(
                        'Reservado',
                        level.reservedQty,
                        theme,
                        color: level.reservedQty > 0 ? Colors.orange : null,
                      ),
                      const SizedBox(width: 8),
                      _buildQtyChip('Disponible', level.availableQty, theme,
                          bold: true),
                    ],
                  ),
                  if (level.minStock != null || level.maxStock != null) ...[
                    const SizedBox(height: 4),
                    Text(
                      'Min: ${level.minStock ?? "-"} / Max: ${level.maxStock ?? "-"}',
                      style: theme.textTheme.bodySmall
                          ?.copyWith(color: theme.disabledColor),
                    ),
                  ],
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildQtyChip(String label, int value, ThemeData theme,
      {Color? color, bool bold = false}) {
    return Expanded(
      child: Column(
        children: [
          Text(label, style: theme.textTheme.labelSmall),
          const SizedBox(height: 2),
          Text(
            value.toString(),
            style: TextStyle(
              fontSize: 16,
              fontWeight: bold ? FontWeight.bold : FontWeight.w500,
              color: color,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPagination(InventoryProvider provider) {
    final pagination = provider.pagination!;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            'Pag. ${pagination.currentPage} de ${pagination.lastPage}  (${pagination.total} total)',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          Row(
            children: [
              IconButton(
                icon: const Icon(Icons.chevron_left),
                onPressed: pagination.hasPrev
                    ? () {
                        provider.setPage(pagination.currentPage - 1);
                        provider.fetchWarehouseInventory(
                            _selectedWarehouseId!,
                            businessId: widget.businessId);
                      }
                    : null,
              ),
              IconButton(
                icon: const Icon(Icons.chevron_right),
                onPressed: pagination.hasNext
                    ? () {
                        provider.setPage(pagination.currentPage + 1);
                        provider.fetchWarehouseInventory(
                            _selectedWarehouseId!,
                            businessId: widget.businessId);
                      }
                    : null,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
