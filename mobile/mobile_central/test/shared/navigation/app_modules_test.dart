import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/navigation/app_modules.dart';

void main() {
  group('estado de modulos', () {
    test('ultima milla esta en desarrollo y no se muestra', () {
      final delivery =
          AppModules.all.firstWhere((m) => m.route == '/delivery');

      expect(delivery.stage, ModuleStage.development);
      expect(delivery.isVisible, isFalse);
    });

    test('ningun grupo visible expone un modulo en desarrollo', () {
      final visibles = AppModules.visibleGroups
          .expand((g) => g.modules)
          .map((m) => m.route)
          .toList();

      expect(visibles, isNot(contains('/delivery')));
      expect(visibles, isNotEmpty);
    });

    test('un grupo que se queda sin modulos desaparece', () {
      for (final group in AppModules.visibleGroups) {
        expect(group.modules, isNotEmpty, reason: group.title);
      }
    });

    test('las rutas de un modulo en desarrollo no estan disponibles', () {
      expect(AppModules.isRouteAvailable('/delivery'), isFalse);
      expect(AppModules.isRouteAvailable('/delivery/drivers'), isFalse);
      expect(AppModules.isRouteAvailable('/delivery/vehicles'), isFalse);
    });

    test('las rutas de produccion siguen disponibles', () {
      expect(AppModules.isRouteAvailable('/orders'), isTrue);
      expect(AppModules.isRouteAvailable('/orders/shipments'), isTrue);
      expect(AppModules.isRouteAvailable('/inventory/stock'), isTrue);
      expect(AppModules.isRouteAvailable('/dashboard'), isTrue);
    });

    test('solo beta lleva insignia', () {
      expect(ModuleStage.prod.showsBadge, isFalse);
      expect(ModuleStage.beta.showsBadge, isTrue);
      expect(ModuleStage.development.showsBadge, isFalse);
      expect(ModuleStage.beta.label, 'Beta');
    });
  });

  group('modulos por rol', () {
    test('el catalogo de integraciones es solo del super admin', () {
      final integraciones =
          AppModules.all.firstWhere((m) => m.route == '/integrations');

      expect(integraciones.superAdminOnly, isTrue);
      expect(
        AppModules.isRouteAllowed('/integrations', isSuperAdmin: false),
        isFalse,
      );
      expect(
        AppModules.isRouteAllowed('/integrations', isSuperAdmin: true),
        isTrue,
      );
    });

    test('el usuario normal no ve el catalogo en el menu', () {
      final rutas = AppModules.visibleGroupsFor(isSuperAdmin: false)
          .expand((g) => g.modules)
          .map((m) => m.route)
          .toList();

      expect(rutas, isNot(contains('/integrations')));
      expect(rutas, contains('/core'));
    });

    test('el super admin si lo ve', () {
      final rutas = AppModules.visibleGroupsFor(isSuperAdmin: true)
          .expand((g) => g.modules)
          .map((m) => m.route)
          .toList();

      expect(rutas, contains('/integrations'));
    });

    test('tus integraciones esta disponible para todos', () {
      expect(AppModules.isRouteAllowed('/core', isSuperAdmin: false), isTrue);
      expect(AppModules.isRouteAllowed('/core', isSuperAdmin: true), isTrue);
    });

    test('ultima milla sigue bloqueada para ambos', () {
      expect(AppModules.isRouteAllowed('/delivery', isSuperAdmin: true), isFalse);
      expect(AppModules.isRouteAllowed('/delivery', isSuperAdmin: false), isFalse);
    });
  });

  group('barra inferior', () {
    test('son cuatro pestanias y ninguna es de envios', () {
      final labels = appBottomTabs.map((t) => t.label).toList();

      expect(labels, ['Inicio', 'Ordenes', 'Inventario', 'Mas']);
    });

    test('las guias quedan bajo la pestania de ordenes', () {
      final ordenes = appBottomTabs.firstWhere((t) => t.label == 'Ordenes');

      expect(ordenes.matches('/orders/shipments'), isTrue);
      expect(ordenes.matches('/orders/shipments/9'), isTrue);
    });

    test('ninguna pestania apunta a un modulo oculto', () {
      for (final tab in appBottomTabs) {
        expect(AppModules.isRouteAvailable(tab.route), isTrue,
            reason: tab.label);
      }
    });

    test('gana la pestania mas especifica, no la primera que coincide', () {
      int indexFor(String location) {
        var best = -1;
        var bestScore = 0;
        for (var i = 0; i < appBottomTabs.length; i++) {
          final score = appBottomTabs[i].matchScore(location);
          if (score > bestScore) {
            bestScore = score;
            best = i;
          }
        }
        return best;
      }

      final ordenes = appBottomTabs.indexWhere((t) => t.label == 'Ordenes');
      final inventario = appBottomTabs.indexWhere((t) => t.label == 'Inventario');

      expect(indexFor('/orders'), ordenes);
      expect(indexFor('/orders/shipments'), ordenes);
      expect(indexFor('/orders/shipments/9'), ordenes);
      expect(indexFor('/inventory/stock'), inventario);
    });
  });
}
