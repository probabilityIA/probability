import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../../services/auth/business/ui/providers/business_provider.dart';
import '../../../services/auth/login/ui/providers/login_provider.dart';
import '../../../services/modules/orders/ui/screens/order_list_screen.dart';
import '../../../services/modules/orderstatus/ui/screens/orderstatus_list_screen.dart';
import '../../../services/modules/shipments/ui/screens/shipment_list_screen.dart';

/// Module wrapper that groups Orders, Shipments and Order Statuses
/// behind a TabBar, replicating the Next.js "subnavbar" pattern.
class OrdersModuleScreen extends StatefulWidget {
  final int initialTab;

  const OrdersModuleScreen({super.key, this.initialTab = 0});

  @override
  State<OrdersModuleScreen> createState() => _OrdersModuleScreenState();
}

class _OrdersModuleScreenState extends State<OrdersModuleScreen> {
  int? _selectedBusinessId;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final login = context.read<LoginProvider>();
      if (login.isSuperAdmin) {
        final biz = context.read<BusinessProvider>();
        if (biz.businessesSimple.isEmpty) {
          biz.fetchBusinessesSimple();
        }
        if (biz.selectedBusinessId != null) {
          setState(() => _selectedBusinessId = biz.selectedBusinessId);
        }
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final login = context.watch<LoginProvider>();
    final isSuperAdmin = login.isSuperAdmin;
    final effectiveBusinessId = isSuperAdmin ? _selectedBusinessId : null;

    if (isSuperAdmin && _selectedBusinessId == null) {
      return Scaffold(
        appBar: AppBar(
          title: const Text('Ordenes'),
        ),
        body: _buildBusinessGate(),
      );
    }

    return DefaultTabController(
      length: 3,
      initialIndex: widget.initialTab,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Ordenes'),
          bottom: _buildTabBar(),
          actions: isSuperAdmin ? [_buildBusinessChip()] : null,
        ),
        body: TabBarView(
          children: [
            OrderListScreen(
              key: ValueKey('orders_$effectiveBusinessId'),
              businessId: effectiveBusinessId,
            ),
            ShipmentListScreen(
              key: ValueKey('shipments_$effectiveBusinessId'),
              businessId: effectiveBusinessId,
            ),
            const OrderStatusListScreen(),
          ],
        ),
      ),
    );
  }

  PreferredSizeWidget _buildTabBar() {
    return const TabBar(
      isScrollable: true,
      tabAlignment: TabAlignment.start,
      tabs: [
        Tab(icon: Icon(Icons.shopping_cart), text: 'Ordenes'),
        Tab(icon: Icon(Icons.local_shipping), text: 'Envios'),
        Tab(icon: Icon(Icons.assignment), text: 'Estados'),
      ],
    );
  }

  Widget _buildBusinessGate() {
    return Consumer<BusinessProvider>(
      builder: (context, bizProvider, _) {
        if (bizProvider.isLoading) {
          return const Center(child: CircularProgressIndicator());
        }

        if (bizProvider.businessesSimple.isEmpty) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.store, size: 64, color: Colors.grey.shade400),
                const SizedBox(height: 16),
                Text(
                  bizProvider.error ?? 'Selecciona un negocio',
                  style: TextStyle(fontSize: 16, color: Colors.grey.shade600),
                ),
                const SizedBox(height: 12),
                FilledButton.icon(
                  onPressed: () => bizProvider.fetchBusinessesSimple(),
                  icon: const Icon(Icons.refresh),
                  label: const Text('Cargar negocios'),
                ),
              ],
            ),
          );
        }

        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: bizProvider.businessesSimple.length,
          itemBuilder: (context, index) {
            final biz = bizProvider.businessesSimple[index];
            return Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                leading: CircleAvatar(
                  backgroundColor:
                      Theme.of(context).colorScheme.primaryContainer,
                  child: Text(
                      biz.name.isNotEmpty ? biz.name[0].toUpperCase() : '?'),
                ),
                title: Text(biz.name),
                subtitle: Text('ID: ${biz.id}'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {
                  context.read<BusinessProvider>().setSelectedBusinessId(biz.id);
                  setState(() => _selectedBusinessId = biz.id);
                },
              ),
            );
          },
        );
      },
    );
  }

  Widget _buildBusinessChip() {
    return Consumer<BusinessProvider>(
      builder: (context, bizProvider, _) {
        final biz = bizProvider.businessesSimple
            .where((b) => b.id == _selectedBusinessId)
            .firstOrNull;
        return Padding(
          padding: const EdgeInsets.only(right: 8),
          child: ActionChip(
            avatar: const Icon(Icons.business, size: 16),
            label: Text(biz?.name ?? 'Negocio $_selectedBusinessId',
                style: const TextStyle(fontSize: 12)),
            onPressed: () {
              context.read<BusinessProvider>().setSelectedBusinessId(null);
              setState(() => _selectedBusinessId = null);
            },
          ),
        );
      },
    );
  }
}
