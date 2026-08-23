import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/sync_entities.dart';
import 'package:mobile_central/services/modules/my_integrations/ui/widgets/module_toolbar.dart';

Widget _wrap(Widget child) => MaterialApp(home: Scaffold(body: child));

Future<void> _wide(WidgetTester tester) async {
  tester.view.physicalSize = const Size(2400, 1200);
  tester.view.devicePixelRatio = 1;
  addTearDown(tester.view.reset);
}

void main() {
  group('coreEnvironments', () {
    test('covers every environment the front offers', () {
      expect(
        coreEnvironments.map((spec) => spec.environment).toSet(),
        SyncEnvironment.values.toSet(),
      );
    });

    test('facturar stays disabled, like in the front', () {
      expect(environmentSpec(SyncEnvironment.invoicing).enabled, isFalse);
    });

    test('only products and inventory can be launched from the panel', () {
      final runnable = coreEnvironments
          .where((spec) => spec.runLabel != null)
          .map((spec) => spec.environment)
          .toSet();
      expect(runnable, {SyncEnvironment.products, SyncEnvironment.inventory});
    });

    test('every environment carries the hint that explains it', () {
      for (final spec in coreEnvironments) {
        expect(spec.hint, isNotEmpty, reason: spec.label);
        expect(spec.label, isNotEmpty);
      }
    });
  });

  group('ModuleToolbar', () {
    testWidgets('renders both views and every environment chip', (tester) async {
      await _wide(tester);
      await tester.pumpWidget(_wrap(ModuleToolbar(
        view: CoreView.diagrama,
        environment: SyncEnvironment.overview,
        onView: (_) {},
        onEnvironment: (_) {},
      )));

      expect(find.text('Diagrama'), findsOneWidget);
      expect(find.text('Informe'), findsOneWidget);
      expect(find.text('Vista general'), findsOneWidget);

      for (final spec in coreEnvironments) {
        expect(find.text(spec.label), findsOneWidget, reason: spec.label);
      }
    });

    testWidgets('tapping a view reports the change', (tester) async {
      CoreView? picked;
      await tester.pumpWidget(_wrap(ModuleToolbar(
        view: CoreView.diagrama,
        environment: SyncEnvironment.overview,
        onView: (value) => picked = value,
        onEnvironment: (_) {},
      )));

      await tester.tap(find.text('Informe'));
      expect(picked, CoreView.informe);
    });

    testWidgets('tapping an environment reports it', (tester) async {
      await _wide(tester);
      SyncEnvironment? picked;
      await tester.pumpWidget(_wrap(ModuleToolbar(
        view: CoreView.informe,
        environment: SyncEnvironment.overview,
        onView: (_) {},
        onEnvironment: (value) => picked = value,
      )));

      await tester.tap(find.text('Comparar ordenes'));
      expect(picked, SyncEnvironment.ordersCompare);
    });

    testWidgets('a disabled environment does not report a change', (tester) async {
      await _wide(tester);
      SyncEnvironment? picked;
      await tester.pumpWidget(_wrap(ModuleToolbar(
        view: CoreView.informe,
        environment: SyncEnvironment.overview,
        onView: (_) {},
        onEnvironment: (value) => picked = value,
      )));

      await tester.tap(find.text('Facturar'));
      expect(picked, isNull);
    });

    testWidgets('while a run is in flight no environment can be switched', (tester) async {
      SyncEnvironment? picked;
      await tester.pumpWidget(_wrap(ModuleToolbar(
        view: CoreView.diagrama,
        environment: SyncEnvironment.inventory,
        running: true,
        onView: (_) {},
        onEnvironment: (value) => picked = value,
      )));

      await tester.tap(find.text('Vista general'));
      expect(picked, isNull);
    });
  });
}
