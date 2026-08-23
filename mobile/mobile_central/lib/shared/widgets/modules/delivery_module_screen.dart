import 'package:flutter/material.dart';
import '../../../services/modules/drivers/ui/screens/driver_list_screen.dart';
import '../../../services/modules/routes/ui/screens/route_list_screen.dart';
import '../../../services/modules/vehicles/ui/screens/vehicle_list_screen.dart';
import 'module_tabs_scaffold.dart';

class DeliveryModuleScreen extends StatelessWidget {
  const DeliveryModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Ultima milla',
      subtitle: 'Rutas, conductores y vehiculos',
      initialTab: initialTab,
      tabs: const ['Rutas', 'Conductores', 'Vehiculos'],
      builder: (context, businessId) => [
        RouteListScreen(key: ValueKey('routes_$businessId'), businessId: businessId),
        DriverListScreen(key: ValueKey('drivers_$businessId'), businessId: businessId),
        VehicleListScreen(key: ValueKey('vehicles_$businessId'), businessId: businessId),
      ],
    );
  }
}
