import 'package:flutter/material.dart';
import '../../../services/modules/invoicing/ui/screens/invoice_list_screen.dart';
import 'module_tabs_scaffold.dart';

class InvoicingModuleScreen extends StatelessWidget {
  const InvoicingModuleScreen({super.key, this.initialTab = 0});

  final int initialTab;

  @override
  Widget build(BuildContext context) {
    return ModuleTabsScaffold(
      title: 'Facturacion',
      subtitle: 'Facturas y notas credito',
      initialTab: initialTab,
      tabs: const ['Facturas'],
      builder: (context, businessId) => [
        InvoiceListScreen(key: ValueKey('invoices_$businessId'), businessId: businessId),
      ],
    );
  }
}
