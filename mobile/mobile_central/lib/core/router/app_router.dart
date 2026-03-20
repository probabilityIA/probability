import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

// Auth screens
import '../../services/auth/login/ui/providers/login_provider.dart';
import '../../services/auth/login/ui/screens/login_screen.dart';
import '../../services/auth/business/ui/screens/business_list_screen.dart';

// Module screens (standalone)
import '../../services/modules/customers/ui/screens/customer_list_screen.dart';
import '../../services/modules/pay/ui/screens/pay_screen.dart';
import '../../services/modules/wallet/ui/screens/wallet_screen.dart';
import '../../services/modules/dashboard/ui/screens/dashboard_screen.dart';
import '../../services/modules/publicsite/ui/screens/publicsite_screen.dart';

// Module wrapper screens (tabbed)
import '../../shared/widgets/modules/orders_module_screen.dart';
import '../../shared/widgets/modules/inventory_module_screen.dart';
import '../../shared/widgets/modules/delivery_module_screen.dart';
import '../../shared/widgets/modules/iam_module_screen.dart';
import '../../shared/widgets/modules/integrations_module_screen.dart';
import '../../shared/widgets/modules/storefront_module_screen.dart';
import '../../shared/widgets/modules/notifications_module_screen.dart';
import '../../shared/widgets/modules/invoicing_module_screen.dart';

// Shared
import '../../shared/widgets/app_shell.dart';
import '../../shared/widgets/business_selector_wrapper.dart';

class AppRouter {
  final LoginProvider loginProvider;

  AppRouter({required this.loginProvider});

  late final GoRouter router = GoRouter(
    refreshListenable: loginProvider,
    initialLocation: '/login',
    redirect: (context, state) {
      final isLoggedIn = loginProvider.isLoggedIn;
      final isLoginRoute = state.matchedLocation == '/login';

      if (!isLoggedIn && !isLoginRoute) return '/login';
      if (isLoggedIn && isLoginRoute) return '/dashboard';
      return null;
    },
    routes: [
      // Login (sin shell)
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),

      // Todas las rutas autenticadas dentro del ShellRoute con Drawer
      ShellRoute(
        builder: (context, state, child) => AppShell(child: child),
        routes: [
          // Dashboard (página inicial)
          GoRoute(
            path: '/dashboard',
            builder: (context, state) => _withBusiness(
              (id) => DashboardScreen(businessId: id),
            ),
          ),

          // ── Ventas ──
          GoRoute(
            path: '/orders',
            builder: (context, state) =>
                const OrdersModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/orders/shipments',
            builder: (context, state) =>
                const OrdersModuleScreen(initialTab: 1),
          ),
          GoRoute(
            path: '/orders/statuses',
            builder: (context, state) =>
                const OrdersModuleScreen(initialTab: 2),
          ),
          GoRoute(
            path: '/customers',
            builder: (context, state) => _withBusiness(
              (id) => CustomerListScreen(businessId: id),
            ),
          ),
          GoRoute(
            path: '/invoicing',
            builder: (context, state) => const InvoicingModuleScreen(),
          ),

          // ── Inventario ──
          GoRoute(
            path: '/inventory',
            builder: (context, state) =>
                const InventoryModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/inventory/warehouses',
            builder: (context, state) =>
                const InventoryModuleScreen(initialTab: 1),
          ),
          GoRoute(
            path: '/inventory/stock',
            builder: (context, state) =>
                const InventoryModuleScreen(initialTab: 2),
          ),

          // ── Logística / Última Milla ──
          GoRoute(
            path: '/delivery',
            builder: (context, state) =>
                const DeliveryModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/delivery/drivers',
            builder: (context, state) =>
                const DeliveryModuleScreen(initialTab: 1),
          ),
          GoRoute(
            path: '/delivery/vehicles',
            builder: (context, state) =>
                const DeliveryModuleScreen(initialTab: 2),
          ),

          // ── Integraciones ──
          GoRoute(
            path: '/integrations',
            builder: (context, state) =>
                const IntegrationsModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/integrations/catalog',
            builder: (context, state) =>
                const IntegrationsModuleScreen(initialTab: 1),
          ),

          // ── Configuración ──
          GoRoute(
            path: '/notifications',
            builder: (context, state) => const NotificationsModuleScreen(),
          ),
          GoRoute(
            path: '/storefront',
            builder: (context, state) =>
                const StorefrontModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/storefront/config',
            builder: (context, state) =>
                const StorefrontModuleScreen(initialTab: 1),
          ),

          // ── Finanzas ──
          GoRoute(
            path: '/wallet',
            builder: (context, state) => _withBusiness(
              (id) => WalletScreen(businessId: id),
            ),
          ),
          GoRoute(
            path: '/pay',
            builder: (context, state) => const PayScreen(),
          ),

          // ── IAM / Administración ──
          GoRoute(
            path: '/iam',
            builder: (context, state) =>
                const IamModuleScreen(initialTab: 0),
          ),
          GoRoute(
            path: '/iam/roles',
            builder: (context, state) =>
                const IamModuleScreen(initialTab: 1),
          ),
          GoRoute(
            path: '/iam/permissions',
            builder: (context, state) =>
                const IamModuleScreen(initialTab: 2),
          ),
          GoRoute(
            path: '/iam/resources',
            builder: (context, state) =>
                const IamModuleScreen(initialTab: 3),
          ),
          GoRoute(
            path: '/iam/actions',
            builder: (context, state) =>
                const IamModuleScreen(initialTab: 4),
          ),
          GoRoute(
            path: '/businesses',
            builder: (context, state) => const BusinessListScreen(),
          ),

          // ── Otros ──
          GoRoute(
            path: '/publicsite',
            builder: (context, state) => const PublicSiteScreen(),
          ),
        ],
      ),
    ],
  );

  static Widget _withBusiness(Widget Function(int?) screenBuilder) {
    return BusinessSelectorWrapper(
      builder: (context, businessId) {
        return screenBuilder(businessId == 0 ? null : businessId);
      },
    );
  }
}
