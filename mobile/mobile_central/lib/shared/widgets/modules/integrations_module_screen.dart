import 'package:flutter/material.dart';
import '../../../services/integrations/core/ui/screens/integration_list_screen.dart';
import '../../../services/modules/my_integrations/ui/screens/my_integrations_screen.dart';
import 'module_tabs_scaffold.dart';

class IntegrationsModuleScreen extends StatelessWidget {
  const IntegrationsModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Integraciones',
      subtitle: 'Canales, transporte y facturacion',
      initialTab: initialTab,
      tabs: const ['Mis integraciones', 'Catalogo'],
      builder: (context, businessId) => [
        MyIntegrationsScreen(key: ValueKey('my_integrations_$businessId'), businessId: businessId),
        const IntegrationListScreen(),
      ],
    );
  }
}
