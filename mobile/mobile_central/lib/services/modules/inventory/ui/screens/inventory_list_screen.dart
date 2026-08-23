import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../../shared/theme/app_colors.dart';
import '../../../../../shared/widgets/ui/ui.dart';
import '../../../warehouses/ui/providers/warehouse_provider.dart';
import '../../domain/entities.dart';
import '../providers/inventory_provider.dart';
import '../widgets/inventory_widgets.dart';
import 'stock_adjust_sheet.dart';

class InventoryListScreen extends StatefulWidget {
  const InventoryListScreen({super.key, this.businessId});

  final int? businessId;

  @override
  State<InventoryListScreen> createState() => _InventoryListScreenState();
}

class _InventoryListScreenState extends State<InventoryListScreen> {
  final _searchController = TextEditingController();
  Timer? _debounce;
  int? _warehouseId;
  String _view = 'stock';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _bootstrap());
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _searchController.dispose();
    super.dispose();
  }

  void _onSearch(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 400), () {
      if (!mounted) return;
      context.read<InventoryProvider>().setFilters(search: value.trim());
      _refresh();
    });
  }

  Future<void> _bootstrap() async {
    final warehouses = context.read<WarehouseProvider>();
    if (warehouses.warehouses.isEmpty) {
      await warehouses.fetchWarehouses(businessId: widget.businessId);
    }
    if (!mounted) return;
    final first = warehouses.warehouses.isNotEmpty ? warehouses.warehouses.first.id : null;
    setState(() => _warehouseId = _warehouseId ?? first);
    _refresh();
  }

  void _refresh() {
    final provider = context.read<InventoryProvider>();
    if (_view == 'stock') {
      if (_warehouseId != null) {
        provider.fetchWarehouseInventory(_warehouseId!, businessId: widget.businessId);
      }
    } else {
      provider.fetchMovements(
        params: GetMovementsParams(
          businessId: widget.businessId,
          warehouseId: _warehouseId,
        ),
      );
    }
  }

  Future<void> _openAdjust(InventoryLevel level) async {
    final changed = await showStockAdjustSheet(
      context,
      level: level,
      businessId: widget.businessId,
    );
    if (changed == true) _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 12),
        _WarehousePicker(
          selected: _warehouseId,
          onSelected: (id) {
            setState(() => _warehouseId = id);
            _refresh();
          },
        ),
        const SizedBox(height: 10),
        AppFilterChips(
          options: const [
            (value: 'stock', label: 'Existencias'),
            (value: 'movements', label: 'Movimientos'),
          ],
          selected: _view,
          onSelected: (value) {
            setState(() => _view = value);
            _refresh();
          },
        ),
        const SizedBox(height: 4),
        Expanded(
          child: _view == 'stock' ? _buildStock() : _buildMovements(),
        ),
      ],
    );
  }

  Widget _buildStock() {
    return Consumer<InventoryProvider>(
      builder: (context, provider, _) {
        return Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 6, 16, 10),
              child: AppSearchField(
                controller: _searchController,
                hintText: 'Producto o SKU',
                onChanged: _onSearch,
              ),
            ),
            Expanded(
              child: PaginatedListView<InventoryLevel>(
                controller: provider.levels,
                unitLabel: 'productos',
                placeholderHeight: 104,
                emptyIcon: Icons.inventory_2_outlined,
                emptyTitle: 'Sin existencias',
                emptyMessage: _searchController.text.trim().isEmpty
                    ? 'Esta bodega no tiene productos con stock registrado.'
                    : 'Ningun producto coincide con la busqueda.',
                itemBuilder: (context, level, index) => InventoryLevelCard(
                  level: level,
                  onTap: () => _openAdjust(level),
                ),
              ),
            ),
          ],
        );
      },
    );
  }

  Widget _buildMovements() {
    return Consumer<InventoryProvider>(
      builder: (context, provider, _) {
        return PaginatedListView<StockMovement>(
          controller: provider.movementsList,
          unitLabel: 'movimientos',
          placeholderHeight: 72,
          emptyIcon: Icons.swap_vert_rounded,
          emptyTitle: 'Sin movimientos',
          emptyMessage:
              'Aun no hay entradas ni salidas registradas en esta bodega.',
          itemBuilder: (context, movement, index) => AppCard(
            padding: const EdgeInsets.symmetric(horizontal: 4),
            child: StockMovementTile(movement: movement),
          ),
        );
      },
    );
  }
}

class _WarehousePicker extends StatelessWidget {
  const _WarehousePicker({required this.selected, required this.onSelected});

  final int? selected;
  final ValueChanged<int> onSelected;

  @override
  Widget build(BuildContext context) {
    return Consumer<WarehouseProvider>(
      builder: (context, provider, _) {
        if (provider.warehouses.isEmpty) return const SizedBox.shrink();

        return SizedBox(
          height: 38,
          child: ListView.separated(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 16),
            itemCount: provider.warehouses.length,
            separatorBuilder: (context, index) => const SizedBox(width: 8),
            itemBuilder: (context, index) {
              final warehouse = provider.warehouses[index];
              final active = warehouse.id == selected;
              return ChoiceChip(
                label: Text(warehouse.name),
                selected: active,
                onSelected: (_) => onSelected(warehouse.id),
                avatar: Icon(
                  Icons.warehouse_outlined,
                  size: 15,
                  color: active ? AppColors.primaryDark : AppColors.textMuted,
                ),
                labelStyle: TextStyle(
                  fontFamily: 'Inter',
                  fontSize: 12.5,
                  fontWeight: FontWeight.w600,
                  color: active ? AppColors.primaryDark : AppColors.textSecondary,
                ),
              );
            },
          ),
        );
      },
    );
  }
}
