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

  group('barra inferior', () {
    test('la pestania de envios apunta a guias, no a ultima milla', () {
      final envios = appBottomTabs.firstWhere((t) => t.label == 'Envios');

      expect(envios.route, '/orders/shipments');
      expect(AppModules.isRouteAvailable(envios.route), isTrue);
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
      final envios = appBottomTabs.indexWhere((t) => t.label == 'Envios');

      expect(indexFor('/orders'), ordenes);
      expect(indexFor('/orders/shipments'), envios);
      expect(indexFor('/orders/shipments/9'), envios);
    });
  });
}
