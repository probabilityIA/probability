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
    this.navKey,
    this.stage = ModuleStage.prod,
    this.superAdminOnly = false,
  });

  final String label;
  final String route;
  final IconData icon;
  final String? description;
  final bool matchPrefix;
  final String? navKey;
  final ModuleStage stage;
  final bool superAdminOnly;

  bool get isVisible => stage.isVisible;

  bool isActive(String location) =>
      matchPrefix ? location.startsWith(route) : location == route;

  bool allowedBy({required bool isSuperAdmin, bool Function(String key)? canNav}) {
    if (isSuperAdmin || navKey == null) return true;
    if (canNav == null) return false;
    return canNav(navKey!);
  }

  bool owns(String location) =>
      location == route || location.startsWith('$route/') || isActive(location);
}

class AppModuleGroup {
  const AppModuleGroup({required this.title, required this.modules});

  final String title;
  final List<AppModule> modules;

  List<AppModule> get visibleModules =>
      modules.where((module) => module.isVisible).toList();

  List<AppModule> visibleFor({
    required bool isSuperAdmin,
    bool Function(String key)? canNav,
  }) =>
      modules
          .where((module) => module.isVisible)
          .where((module) => isSuperAdmin || !module.superAdminOnly)
          .where((module) => module.allowedBy(isSuperAdmin: isSuperAdmin, canNav: canNav))
          .toList();
}

class AppModules {
  const AppModules._();

  static const AppModule dashboard = AppModule(
    label: 'Inicio',
    route: '/dashboard',
    icon: Icons.space_dashboard_outlined,
    description: 'Resumen de la operaci\u00f3n',
  );

  static const List<AppModuleGroup> groups = [
    AppModuleGroup(
      title: 'Ventas',
      modules: [
        AppModule(
          label: '\u00d3rdenes',
          route: '/orders',
          navKey: 'orders',
          icon: Icons.receipt_long_outlined,
          description: 'Pedidos de todos los canales',
          matchPrefix: true,
        ),
        AppModule(
          label: 'Clientes',
          route: '/customers',
          navKey: 'customers',
          icon: Icons.people_alt_outlined,
          description: 'Directorio y compras',
        ),
        AppModule(
          label: 'Facturaci\u00f3n',
          route: '/invoicing',
          navKey: 'invoicing',
          icon: Icons.description_outlined,
          description: 'Facturas y notas cr\u00e9dito',
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Log\u00edstica',
      modules: [
        AppModule(
          label: 'Env\u00edos',
          route: '/orders/shipments',
          navKey: 'shipments',
          icon: Icons.local_shipping_outlined,
          description: 'Gu\u00edas y seguimiento',
        ),
        AppModule(
          label: 'Ultima milla',
          route: '/delivery',
          navKey: 'delivery',
          icon: Icons.alt_route_outlined,
          description: 'Rutas, conductores y vehiculos',
          matchPrefix: true,
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
          navKey: 'products',
          icon: Icons.sell_outlined,
          description: 'Catalogo y precios',
        ),
        AppModule(
          label: 'Bodegas',
          route: '/inventory/warehouses',
          navKey: 'warehouses',
          icon: Icons.warehouse_outlined,
          description: 'Ubicaciones y ocupacion',
        ),
        AppModule(
          label: 'Stock',
          route: '/inventory/stock',
          navKey: 'inventory',
          icon: Icons.inventory_2_outlined,
          description: 'Existencias y movimientos',
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Finanzas',
      modules: [
        AppModule(
          label: 'Billetera',
          route: '/wallet',
          navKey: 'wallet',
          icon: Icons.account_balance_wallet_outlined,
          description: 'Saldo y movimientos',
        ),
        AppModule(
          label: 'Pagos',
          route: '/pay',
          navKey: 'integrations',
          icon: Icons.credit_card_outlined,
          description: 'Pasarelas y recaudo',
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Canales',
      modules: [
        AppModule(
          label: 'Integraciones',
          route: '/integrations',
          navKey: 'integrations',
          icon: Icons.hub_outlined,
          description: 'Catalogo de conectores',
          matchPrefix: true,
          superAdminOnly: true,
        ),
        AppModule(
          label: 'Tus integraciones',
          route: '/core',
          icon: Icons.auto_awesome_outlined,
          description: 'Lo que tienes conectado',
        ),
        AppModule(
          label: 'Tienda online',
          route: '/storefront',
          navKey: 'storefront',
          icon: Icons.storefront_outlined,
          description: 'Catalogo p\u00fablico y sitio',
          matchPrefix: true,
        ),
        AppModule(
          label: 'Notificaciones',
          route: '/notifications',
          navKey: 'notifications',
          icon: Icons.notifications_none_outlined,
          description: 'Eventos y plantillas',
        ),
      ],
    ),
    AppModuleGroup(
      title: 'Administracion',
      modules: [
        AppModule(
          label: 'Usuarios y roles',
          route: '/iam',
          navKey: 'users',
          icon: Icons.admin_panel_settings_outlined,
          description: 'Accesos y permisos',
          matchPrefix: true,
        ),
        AppModule(
          label: 'Negocios',
          route: '/businesses',
          navKey: 'businesses',
          icon: Icons.apartment_outlined,
          description: 'Empresas de la plataforma',
        ),
      ],
    ),
  ];

  static List<AppModule> get all =>
      groups.expand((group) => group.modules).toList();

  static List<AppModuleGroup> get visibleGroups =>
      visibleGroupsFor(isSuperAdmin: true);

  static List<AppModuleGroup> visibleGroupsFor({
    required bool isSuperAdmin,
    bool Function(String key)? canNav,
  }) =>
      groups
          .map((group) => AppModuleGroup(
                title: group.title,
                modules: group.visibleFor(isSuperAdmin: isSuperAdmin, canNav: canNav),
              ))
          .where((group) => group.modules.isNotEmpty)
          .toList();

  static bool isRouteAllowed(
    String location, {
    required bool isSuperAdmin,
    bool Function(String key)? canNav,
  }) {
    AppModule? owner;
    for (final module in all) {
      if (module.owns(location) && (owner == null || module.route.length > owner.route.length)) {
        owner = module;
      }
    }
    if (owner == null) return true;
    if (!owner.isVisible) return false;
    if (owner.superAdminOnly && !isSuperAdmin) return false;
    return owner.allowedBy(isSuperAdmin: isSuperAdmin, canNav: canNav);
  }

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
    label: '\u00d3rdenes',
    route: '/orders',
    icon: Icons.receipt_long_outlined,
    activeIcon: Icons.receipt_long,
    prefixes: ['/orders'],
  ),
  AppBottomTab(
    label: 'Inventario',
    route: '/inventory',
    icon: Icons.inventory_2_outlined,
    activeIcon: Icons.inventory_2,
    prefixes: ['/inventory'],
  ),
  AppBottomTab(
    label: 'M\u00e1s',
    route: '/more',
    icon: Icons.grid_view_outlined,
    activeIcon: Icons.grid_view,
    prefixes: ['/more'],
  ),
];
