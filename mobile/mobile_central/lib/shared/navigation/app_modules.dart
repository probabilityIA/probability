import 'package:flutter/material.dart';

enum ModuleStage { prod, beta, development }

extension ModuleStageX on ModuleStage {
  bool get isVisible => this != ModuleStage.development;
  bool get showsBadge => this == ModuleStage.beta;

  String get label {
    switch (this) {
      case ModuleStage.prod:
        return 'Produccion';
      case ModuleStage.beta:
        return 'Beta';
      case ModuleStage.development:
        return 'En desarrollo';
    }
  }
}

class AppModule {
  const AppModule({
    required this.label,
    required this.route,
    required this.icon,
    this.description,
    this.matchPrefix = false,
    this.resources = const [],
    this.stage = ModuleStage.prod,
  });

  final String label;
  final String route;
  final IconData icon;
  final String? description;
  final bool matchPrefix;
  final List<String> resources;
  final ModuleStage stage;

  bool get isVisible => stage.isVisible;

  bool isActive(String location) =>
      matchPrefix ? location.startsWith(route) : location == route;

  bool owns(String location) =>
      location == route || location.startsWith('$route/') || isActive(location);
}

class AppModuleGroup {
  const AppModuleGroup({required this.title, required this.modules});

  final String title;
  final List<AppModule> modules;

  List<AppModule> get visibleModules =>
      modules.where((module) => module.isVisible).toList();
}

class AppModules {
  const AppModules._();

  static const AppModule dashboard = AppModule(
    label: 'Inicio',
    route: '/dashboard',
    icon: Icons.space_dashboard_outlined,
    description: 'Resumen de la operacion',
  );

  static const List<AppModuleGroup> groups = [
    AppModuleGroup(
      title: 'Ventas',
      modules: [
        AppModule(
          label: 'Ordenes',
          route: '/orders',
          icon: Icons.receipt_long_outlined,
          description: 'Pedidos de todos los canales',
          matchPrefix: true,
          resources: ['Ordenes', 'Orders'],
        ),
        AppModule(
          label: 'Clientes',
          route: '/customers',
          icon: Icons.people_alt_outlined,
          description: 'Directorio y compras',
          resources: ['Clientes', 'Customers'],
        ),
        AppModule(
          label: 'Facturacion',
          route: '/invoicing',
          icon: Icons.description_outlined,
          description: 'Facturas y notas credito',
          resources: ['Facturacion', 'Invoicing'],
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Logistica',
      modules: [
        AppModule(
          label: 'Envios',
          route: '/orders/shipments',
          icon: Icons.local_shipping_outlined,
          description: 'Guias y seguimiento',
          resources: ['Envios', 'Shipments'],
        ),
        AppModule(
          label: 'Ultima milla',
          route: '/delivery',
          icon: Icons.alt_route_outlined,
          description: 'Rutas, conductores y vehiculos',
          matchPrefix: true,
          resources: ['Rutas', 'Routes'],
          stage: ModuleStage.development,
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Inventario',
      modules: [
        AppModule(
          label: 'Productos',
          route: '/inventory',
          icon: Icons.sell_outlined,
          description: 'Catalogo y precios',
          resources: ['Productos', 'Products'],
        ),
        AppModule(
          label: 'Bodegas',
          route: '/inventory/warehouses',
          icon: Icons.warehouse_outlined,
          description: 'Ubicaciones y ocupacion',
          resources: ['Bodegas', 'Warehouses'],
        ),
        AppModule(
          label: 'Stock',
          route: '/inventory/stock',
          icon: Icons.inventory_2_outlined,
          description: 'Existencias y movimientos',
          resources: ['Inventario', 'Inventory'],
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Finanzas',
      modules: [
        AppModule(
          label: 'Billetera',
          route: '/wallet',
          icon: Icons.account_balance_wallet_outlined,
          description: 'Saldo y movimientos',
          resources: ['Billetera', 'Wallet'],
        ),
        AppModule(
          label: 'Pagos',
          route: '/pay',
          icon: Icons.credit_card_outlined,
          description: 'Pasarelas y recaudo',
          resources: ['Pagos', 'Pay'],
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Canales',
      modules: [
        AppModule(
          label: 'Integraciones',
          route: '/integrations',
          icon: Icons.hub_outlined,
          description: 'Tiendas, transporte y facturacion',
          matchPrefix: true,
          resources: ['Integraciones', 'Integrations'],
        ),
        AppModule(
          label: 'Tienda online',
          route: '/storefront',
          icon: Icons.storefront_outlined,
          description: 'Catalogo publico y sitio',
          matchPrefix: true,
          resources: ['Tienda', 'Storefront'],
        ),
        AppModule(
          label: 'Notificaciones',
          route: '/notifications',
          icon: Icons.notifications_none_outlined,
          description: 'Eventos y plantillas',
          resources: ['Notificaciones', 'Notifications'],
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Administracion',
      modules: [
        AppModule(
          label: 'Usuarios y roles',
          route: '/iam',
          icon: Icons.admin_panel_settings_outlined,
          description: 'Accesos y permisos',
          matchPrefix: true,
          resources: ['Usuarios', 'Users', 'Roles', 'Permisos', 'Permissions'],
        ),
        AppModule(
          label: 'Negocios',
          route: '/businesses',
          icon: Icons.apartment_outlined,
          description: 'Empresas de la plataforma',
          resources: ['Empresas', 'Businesses'],
        ),
      ],
    ),
  ];

  static List<AppModule> get all =>
      groups.expand((group) => group.modules).toList();

  static List<AppModuleGroup> get visibleGroups => groups
      .map((group) => AppModuleGroup(
            title: group.title,
            modules: group.visibleModules,
          ))
      .where((group) => group.modules.isNotEmpty)
      .toList();

  static bool isRouteAvailable(String location) {
    for (final module in all) {
      if (module.owns(location)) return module.isVisible;
    }
    return true;
  }

  static AppModule? byRoute(String location) {
    if (dashboard.isActive(location)) return dashboard;
    AppModule? best;
    for (final module in all) {
      if (location == module.route) return module;
      if (location.startsWith('${module.route}/') || module.isActive(location)) {
        if (best == null || module.route.length > best.route.length) best = module;
      }
    }
    return best;
  }
}

class AppBottomTab {
  const AppBottomTab({
    required this.label,
    required this.route,
    required this.icon,
    required this.activeIcon,
    this.prefixes = const [],
  });

  final String label;
  final String route;
  final IconData icon;
  final IconData activeIcon;
  final List<String> prefixes;

  bool matches(String location) => matchScore(location) > 0;

  int matchScore(String location) {
    if (location == route) return route.length + 1;
    var best = 0;
    for (final prefix in prefixes) {
      if (location == prefix || location.startsWith('$prefix/')) {
        if (prefix.length > best) best = prefix.length;
      }
    }
    return best;
  }
}

const List<AppBottomTab> appBottomTabs = [
  AppBottomTab(
    label: 'Inicio',
    route: '/dashboard',
    icon: Icons.space_dashboard_outlined,
    activeIcon: Icons.space_dashboard,
  ),
  AppBottomTab(
    label: 'Ordenes',
    route: '/orders',
    icon: Icons.receipt_long_outlined,
    activeIcon: Icons.receipt_long,
    prefixes: ['/orders'],
  ),
  AppBottomTab(
    label: 'Envios',
    route: '/orders/shipments',
    icon: Icons.local_shipping_outlined,
    activeIcon: Icons.local_shipping,
    prefixes: ['/orders/shipments'],
  ),
  AppBottomTab(
    label: 'Inventario',
    route: '/inventory',
    icon: Icons.inventory_2_outlined,
    activeIcon: Icons.inventory_2,
    prefixes: ['/inventory'],
  ),
  AppBottomTab(
    label: 'Mas',
    route: '/more',
    icon: Icons.grid_view_outlined,
    activeIcon: Icons.grid_view,
    prefixes: ['/more'],
  ),
];
