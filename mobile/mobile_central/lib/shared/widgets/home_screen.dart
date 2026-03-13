import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Probability Central'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              context.read<LoginProvider>().logout();
            },
          ),
        ],
      ),
      body: Consumer<LoginProvider>(
        builder: (context, loginProvider, child) {
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Bienvenido, ${loginProvider.user?.name ?? ""}',
                        style: Theme.of(context).textTheme.titleLarge,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        loginProvider.user?.email ?? '',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                              color: Colors.grey[600],
                            ),
                      ),
                      if (loginProvider.isSuperAdmin)
                        const Padding(
                          padding: EdgeInsets.only(top: 8),
                          child: Chip(
                            label: Text('Super Admin'),
                            avatar: Icon(Icons.shield, size: 16),
                          ),
                        ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),
              _MenuSection(
                title: 'Autenticación',
                items: [
                  _MenuItem(icon: Icons.people, label: 'Usuarios', route: '/users'),
                  _MenuItem(icon: Icons.admin_panel_settings, label: 'Roles', route: '/roles'),
                  _MenuItem(icon: Icons.security, label: 'Permisos', route: '/permissions'),
                  _MenuItem(icon: Icons.business, label: 'Negocios', route: '/businesses'),
                  _MenuItem(icon: Icons.widgets, label: 'Recursos', route: '/resources'),
                  _MenuItem(icon: Icons.touch_app, label: 'Acciones', route: '/actions'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Ventas',
                items: [
                  _MenuItem(icon: Icons.shopping_cart, label: 'Órdenes', route: '/orders'),
                  _MenuItem(icon: Icons.inventory_2, label: 'Productos', route: '/products'),
                  _MenuItem(icon: Icons.person_outline, label: 'Clientes', route: '/customers'),
                  _MenuItem(icon: Icons.receipt_long, label: 'Facturación', route: '/invoicing'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Estados',
                items: [
                  _MenuItem(icon: Icons.assignment, label: 'Estado Orden', route: '/orderstatus'),
                  _MenuItem(icon: Icons.payment, label: 'Estado Pago', route: '/paymentstatus'),
                  _MenuItem(icon: Icons.local_shipping, label: 'Estado Envío', route: '/fulfillmentstatus'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Logística',
                items: [
                  _MenuItem(icon: Icons.local_shipping, label: 'Envíos', route: '/shipments'),
                  _MenuItem(icon: Icons.warehouse, label: 'Bodegas', route: '/warehouses'),
                  _MenuItem(icon: Icons.inventory, label: 'Inventario', route: '/inventory'),
                  _MenuItem(icon: Icons.person_pin, label: 'Conductores', route: '/drivers'),
                  _MenuItem(icon: Icons.directions_car, label: 'Vehículos', route: '/vehicles'),
                  _MenuItem(icon: Icons.route, label: 'Rutas', route: '/routes'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Finanzas',
                items: [
                  _MenuItem(icon: Icons.account_balance_wallet, label: 'Pagos', route: '/pay'),
                  _MenuItem(icon: Icons.wallet, label: 'Wallet', route: '/wallet'),
                  _MenuItem(icon: Icons.dashboard, label: 'Dashboard', route: '/dashboard'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Configuración',
                items: [
                  _MenuItem(icon: Icons.notifications, label: 'Notificaciones', route: '/notification-config'),
                  _MenuItem(icon: Icons.storefront, label: 'Storefront', route: '/storefront'),
                  _MenuItem(icon: Icons.public, label: 'Sitio Público', route: '/publicsite'),
                  _MenuItem(icon: Icons.web, label: 'Config Web', route: '/website-config'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Integraciones',
                items: [
                  _MenuItem(icon: Icons.extension, label: 'Mis Integraciones', route: '/my-integrations'),
                  _MenuItem(icon: Icons.hub, label: 'Catálogo', route: '/integrations'),
                ],
              ),
            ],
          );
        },
      ),
    );
  }
}

class _MenuSection extends StatelessWidget {
  final String title;
  final List<_MenuItem> items;

  const _MenuSection({required this.title, required this.items});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
        ),
        const SizedBox(height: 12),
        GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 3,
            mainAxisSpacing: 12,
            crossAxisSpacing: 12,
            childAspectRatio: 0.85,
          ),
          itemCount: items.length,
          itemBuilder: (context, index) {
            final item = items[index];
            return Card(
              child: InkWell(
                borderRadius: BorderRadius.circular(12),
                onTap: () => context.push(item.route),
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(item.icon, size: 28, color: Colors.deepPurple),
                      const SizedBox(height: 6),
                      Flexible(
                        child: Text(
                          item.label,
                          textAlign: TextAlign.center,
                          overflow: TextOverflow.ellipsis,
                          maxLines: 2,
                          style: const TextStyle(fontSize: 11),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        ),
      ],
    );
  }
}

class _MenuItem {
  final IconData icon;
  final String label;
  final String route;

  const _MenuItem({
    required this.icon,
    required this.label,
    required this.route,
  });
}
