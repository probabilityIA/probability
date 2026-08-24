import 'package:flutter/material.dart';
import '../../../services/modules/inventory/ui/screens/inventory_list_screen.dart';
import '../../../services/modules/products/ui/screens/product_list_screen.dart';
import '../../../services/modules/warehouses/ui/screens/warehouse_list_screen.dart';
import 'module_tabs_scaffold.dart';

class InventoryModuleScreen extends StatelessWidget {
  const InventoryModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Inventario',
      subtitle: 'Catalogo, bodegas y existencias',
      initialTab: initialTab,
      tabs: const ['Productos', 'Bodegas', 'Stock'],
      builder: (context, businessId) => [
        ProductListScreen(key: ValueKey('products_$businessId'), businessId: businessId),
        WarehouseListScreen(key: ValueKey('warehouses_$businessId'), businessId: businessId),
        InventoryListScreen(key: ValueKey('stock_$businessId'), businessId: businessId),
      ],
    );
  }
}
