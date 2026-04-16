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
                title: 'Administracion',
                items: [
                  _MenuItem(icon: Icons.people, label: 'Usuarios', route: '/iam'),
                  _MenuItem(icon: Icons.admin_panel_settings, label: 'Roles', route: '/iam/roles'),
                  _MenuItem(icon: Icons.security, label: 'Permisos', route: '/iam/permissions'),
                  _MenuItem(icon: Icons.business, label: 'Negocios', route: '/businesses'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Ventas',
                items: [
                  _MenuItem(icon: Icons.shopping_cart, label: 'Ordenes', route: '/orders'),
                  _MenuItem(icon: Icons.local_shipping, label: 'Envios', route: '/orders/shipments'),
                  _MenuItem(icon: Icons.person_outline, label: 'Clientes', route: '/customers'),
                  _MenuItem(icon: Icons.receipt_long, label: 'Facturacion', route: '/invoicing'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Inventario',
                items: [
                  _MenuItem(icon: Icons.inventory_2, label: 'Productos', route: '/inventory'),
                  _MenuItem(icon: Icons.warehouse, label: 'Bodegas', route: '/inventory/warehouses'),
                  _MenuItem(icon: Icons.assessment, label: 'Stock', route: '/inventory/stock'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Ultima Milla',
                items: [
                  _MenuItem(icon: Icons.route, label: 'Rutas', route: '/delivery'),
                  _MenuItem(icon: Icons.badge, label: 'Conductores', route: '/delivery/drivers'),
                  _MenuItem(icon: Icons.directions_car, label: 'Vehiculos', route: '/delivery/vehicles'),
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
                title: 'Tienda',
                items: [
                  _MenuItem(icon: Icons.storefront, label: 'Catalogo', route: '/storefront'),
                  _MenuItem(icon: Icons.web, label: 'Config Web', route: '/storefront/config'),
                  _MenuItem(icon: Icons.public, label: 'Sitio Publico', route: '/publicsite'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Configuracion',
                items: [
                  _MenuItem(icon: Icons.notifications, label: 'Notificaciones', route: '/notifications'),
                ],
              ),
              const SizedBox(height: 16),
              _MenuSection(
                title: 'Integraciones',
                items: [
                  _MenuItem(icon: Icons.extension, label: 'Mis Integraciones', route: '/integrations'),
                  _MenuItem(icon: Icons.hub, label: 'Catalogo', route: '/integrations/catalog'),
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
