import 'package:flutter/material.dart';
import '../../../services/modules/orders/ui/screens/order_list_screen.dart';
import '../../../services/modules/orderstatus/ui/screens/orderstatus_list_screen.dart';
import '../../../services/modules/shipments/ui/screens/shipment_list_screen.dart';
import 'module_tabs_scaffold.dart';

class OrdersModuleScreen extends StatelessWidget {
  const OrdersModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Ordenes',
      subtitle: 'Pedidos de todos los canales',
      initialTab: initialTab,
      tabs: const ['Ordenes', 'Envios', 'Estados'],
      builder: (context, businessId) => [
        OrderListScreen(key: ValueKey('orders_$businessId'), businessId: businessId),
        ShipmentListScreen(key: ValueKey('shipments_$businessId'), businessId: businessId),
        const OrderStatusListScreen(),
      ],
    );
  }
}
