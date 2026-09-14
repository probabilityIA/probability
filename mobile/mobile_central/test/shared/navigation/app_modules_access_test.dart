import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/navigation/app_modules.dart';

void main() {
  bool Function(String) navOf(Set<String> keys) => keys.contains;

  group('AppModules con navegacion del backend', () {
    test('un rol sin clientes no puede entrar a /customers', () {
      final allowed = AppModules.isRouteAllowed(
        '/customers',
        isSuperAdmin: false,
        canNav: navOf({'orders', 'shipments'}),
      );
      expect(allowed, isFalse);
    });

    test('envios usa su propia clave aunque cuelgue de /orders', () {
      final canNav = navOf({'shipments'});
      expect(AppModules.isRouteAllowed('/orders/shipments', isSuperAdmin: false, canNav: canNav), isTrue);
      expect(AppModules.isRouteAllowed('/orders', isSuperAdmin: false, canNav: canNav), isFalse);
    });

    test('el super admin entra a todo lo visible', () {
      expect(AppModules.isRouteAllowed('/customers', isSuperAdmin: true), isTrue);
    });

    test('el listado de modulos solo muestra lo que el backend permite', () {
      final groups = AppModules.visibleGroupsFor(
        isSuperAdmin: false,
        canNav: navOf({'orders'}),
      );
      final routes = groups.expand((g) => g.modules).map((m) => m.route).toSet();
      expect(routes.contains('/orders'), isTrue);
      expect(routes.contains('/customers'), isFalse);
      expect(routes.contains('/core'), isTrue);
    });
  });
}
