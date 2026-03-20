import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../../services/auth/login/ui/providers/login_provider.dart';
import 'network_avatar.dart';

/// Main navigation shell for the app.
/// Provides a Scaffold with a persistent Drawer accessible via hamburger icon.
/// The child screens provide their own AppBar content; this shell only adds the Drawer.
class AppShell extends StatelessWidget {
  final Widget child;

  const AppShell({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return _AppShellScaffold(child: child);
  }
}

class _AppShellScaffold extends StatefulWidget {
  final Widget child;
  const _AppShellScaffold({required this.child});

  @override
  State<_AppShellScaffold> createState() => _AppShellScaffoldState();
}

class _AppShellScaffoldState extends State<_AppShellScaffold> {
  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: _scaffoldKey,
      drawer: const _AppDrawer(),
      body: Stack(
        children: [
          // Child content fills the entire screen
          widget.child,
          // Hamburger button floating top-left
          Positioned(
            top: MediaQuery.of(context).padding.top + 8,
            left: 8,
            child: Material(
              color: Theme.of(context).colorScheme.surface.withAlpha(220),
              borderRadius: BorderRadius.circular(12),
              elevation: 2,
              child: InkWell(
                borderRadius: BorderRadius.circular(12),
                onTap: () => _scaffoldKey.currentState?.openDrawer(),
                child: Padding(
                  padding: const EdgeInsets.all(10),
                  child: Icon(
                    Icons.menu,
                    color: Theme.of(context).colorScheme.onSurface,
                    size: 24,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Drawer
// ---------------------------------------------------------------------------

class _AppDrawer extends StatelessWidget {
  const _AppDrawer();

  @override
  Widget build(BuildContext context) {
    final loginProvider = context.watch<LoginProvider>();
    final currentLocation = GoRouterState.of(context).matchedLocation;
    final colorScheme = Theme.of(context).colorScheme;

    return Drawer(
      child: SafeArea(
        child: Column(
          children: [
            // ── Header ──────────────────────────────────────────────
            _DrawerHeader(loginProvider: loginProvider),
            const Divider(height: 1),

            // ── Navigation items ────────────────────────────────────
            Expanded(
              child: ListView(
                padding: EdgeInsets.zero,
                children: [
                  _NavTile(
                    icon: Icons.dashboard_outlined,
                    label: 'Dashboard',
                    route: '/dashboard',
                    currentRoute: currentLocation,
                  ),
                  _NavTile(
                    icon: Icons.shopping_cart_outlined,
                    label: 'Ordenes',
                    route: '/orders',
                    currentRoute: currentLocation,
                    activeOnPrefix: true,
                  ),
                  _NavTile(
                    icon: Icons.people_outline,
                    label: 'Clientes',
                    route: '/customers',
                    currentRoute: currentLocation,
                  ),
                  _NavTile(
                    icon: Icons.receipt_long_outlined,
                    label: 'Facturacion',
                    route: '/invoicing',
                    currentRoute: currentLocation,
                  ),
                  _NavTile(
                    icon: Icons.inventory_2_outlined,
                    label: 'Inventario',
                    route: '/inventory',
                    currentRoute: currentLocation,
                    activeOnPrefix: true,
                  ),
                  _NavTile(
                    icon: Icons.route_outlined,
                    label: 'Ultima Milla',
                    route: '/delivery',
                    currentRoute: currentLocation,
                    activeOnPrefix: true,
                  ),
                  _NavTile(
                    icon: Icons.hub_outlined,
                    label: 'Integraciones',
                    route: '/integrations',
                    currentRoute: currentLocation,
                    activeOnPrefix: true,
                  ),
                  _NavTile(
                    icon: Icons.notifications_outlined,
                    label: 'Notificaciones',
                    route: '/notifications',
                    currentRoute: currentLocation,
                  ),
                  _NavTile(
                    icon: Icons.storefront_outlined,
                    label: 'Tienda',
                    route: '/storefront',
                    currentRoute: currentLocation,
                    activeOnPrefix: true,
                  ),
                  _NavTile(
                    icon: Icons.wallet_outlined,
                    label: 'Billetera',
                    route: '/wallet',
                    currentRoute: currentLocation,
                  ),
                  _NavTile(
                    icon: Icons.payment_outlined,
                    label: 'Pagos',
                    route: '/pay',
                    currentRoute: currentLocation,
                  ),
                  if (_canAccessIAM(loginProvider)) ...[
                    _NavTile(
                      icon: Icons.admin_panel_settings_outlined,
                      label: 'Administracion',
                      route: '/iam',
                      currentRoute: currentLocation,
                      activeOnPrefix: true,
                    ),
                    _NavTile(
                      icon: Icons.business_outlined,
                      label: 'Negocios',
                      route: '/businesses',
                      currentRoute: currentLocation,
                    ),
                  ],
                ],
              ),
            ),

            // ── Footer: Logout ──────────────────────────────────────
            const Divider(height: 1),
            ListTile(
              leading: Icon(Icons.logout, color: colorScheme.error),
              title: Text(
                'Cerrar Sesion',
                style: TextStyle(color: colorScheme.error),
              ),
              onTap: () {
                Navigator.pop(context); // close drawer
                loginProvider.logout();
                context.go('/login');
              },
            ),
          ],
        ),
      ),
    );
  }

  /// Returns true when the IAM section should be visible.
  bool _canAccessIAM(LoginProvider provider) {
    if (provider.isSuperAdmin) return true;
    // Show IAM section if the user has any IAM-related permission.
    const iamResources = [
      'Usuarios', 'Users', 'Empleados',
      'Roles', 'Roles y Permisos',
      'Permisos', 'Permissions',
      'Empresas',
    ];
    for (final res in iamResources) {
      if (provider.hasPermission(res, 'Read')) return true;
    }
    return false;
  }
}

// ---------------------------------------------------------------------------
// Drawer header — user info + super admin badge
// ---------------------------------------------------------------------------

class _DrawerHeader extends StatelessWidget {
  final LoginProvider loginProvider;

  const _DrawerHeader({required this.loginProvider});

  @override
  Widget build(BuildContext context) {
    final user = loginProvider.user;
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.fromLTRB(16, 24, 16, 16),
      decoration: BoxDecoration(
        color: colorScheme.primaryContainer.withAlpha(80),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Avatar
          NetworkAvatar(
            imageUrl: user?.avatarUrl,
            fallbackText: user?.name ?? '?',
            radius: 30,
            backgroundColor: colorScheme.primary,
            foregroundColor: colorScheme.onPrimary,
          ),
          const SizedBox(height: 12),

          // Name
          Text(
            user?.name ?? '',
            style: textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 2),

          // Email
          Text(
            user?.email ?? '',
            style: textTheme.bodySmall?.copyWith(
              color: colorScheme.onSurfaceVariant,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),

          // Super Admin badge
          if (loginProvider.isSuperAdmin) ...[
            const SizedBox(height: 8),
            Chip(
              avatar: Icon(Icons.shield_outlined,
                  size: 16, color: colorScheme.primary),
              label: const Text('Super Admin'),
              labelStyle: textTheme.labelSmall,
              visualDensity: VisualDensity.compact,
              padding: EdgeInsets.zero,
            ),
          ],
        ],
      ),
    );
  }
}



// ---------------------------------------------------------------------------
// Navigation tile (single drawer item)
// ---------------------------------------------------------------------------

class _NavTile extends StatelessWidget {
  final IconData icon;
  final String label;
  final String route;
  final String currentRoute;
  final bool activeOnPrefix;

  const _NavTile({
    required this.icon,
    required this.label,
    required this.route,
    required this.currentRoute,
    this.activeOnPrefix = false,
  });

  bool get _isActive =>
      activeOnPrefix ? currentRoute.startsWith(route) : currentRoute == route;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 1),
      child: ListTile(
        dense: true,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        selected: _isActive,
        selectedTileColor: colorScheme.primaryContainer.withAlpha(120),
        selectedColor: colorScheme.primary,
        leading: Icon(icon, size: 22),
        title: Text(label),
        onTap: () {
          Navigator.pop(context); // close drawer first
          if (!_isActive) {
            context.go(route);
          }
        },
      ),
    );
  }
}
