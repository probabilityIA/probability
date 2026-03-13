import 'package:flutter/material.dart';
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
                  _MenuItem(
                    icon: Icons.people,
                    label: 'Usuarios',
                    route: '/users',
                  ),
                  _MenuItem(
                    icon: Icons.admin_panel_settings,
                    label: 'Roles',
                    route: '/roles',
                  ),
                  _MenuItem(
                    icon: Icons.security,
                    label: 'Permisos',
                    route: '/permissions',
                  ),
                  _MenuItem(
                    icon: Icons.business,
                    label: 'Negocios',
                    route: '/businesses',
                  ),
                  _MenuItem(
                    icon: Icons.widgets,
                    label: 'Recursos',
                    route: '/resources',
                  ),
                  _MenuItem(
                    icon: Icons.touch_app,
                    label: 'Acciones',
                    route: '/actions',
                  ),
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
            childAspectRatio: 1,
          ),
          itemCount: items.length,
          itemBuilder: (context, index) {
            final item = items[index];
            return Card(
              child: InkWell(
                borderRadius: BorderRadius.circular(12),
                onTap: () {
                  // Navigation handled by GoRouter
                },
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(item.icon, size: 32, color: Colors.deepPurple),
                    const SizedBox(height: 8),
                    Text(
                      item.label,
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 12),
                    ),
                  ],
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
