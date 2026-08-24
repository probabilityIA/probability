import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/pagination/paged_list_controller.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';
import 'package:mobile_central/shared/widgets/ui/paginated_list_view.dart';

class Source {
  Source({required this.total, this.delay = Duration.zero});

  final int total;
  final Duration delay;
  final List<int> requested = [];

  Future<PaginatedResponse<String>> fetch(int page, int pageSize) async {
    requested.add(page);
    if (delay > Duration.zero) await Future<void>.delayed(delay);
    final start = (page - 1) * pageSize;
    final end = (start + pageSize) > total ? total : (start + pageSize);
    return PaginatedResponse<String>(
      data: [for (var i = start; i < end; i++) 'fila-$i'],
      pagination: Pagination(
        currentPage: page,
        perPage: pageSize,
        total: total,
        lastPage: (total / pageSize).ceil(),
        hasNext: end < total,
        hasPrev: page > 1,
      ),
    );
  }
}

Widget _app(PagedListController<String> controller) {
  return MaterialApp(
    home: Scaffold(
      body: PaginatedListView<String>(
        controller: controller,
        unitLabel: 'filas',
        placeholderHeight: 60,
        emptyIcon: Icons.inbox_outlined,
        emptyTitle: 'Vacio',
        emptyMessage: 'No hay nada',
        itemBuilder: (context, item, index) =>
            SizedBox(height: 60, child: Text(item)),
      ),
    ),
  );
}

void main() {
  group('PaginatedListView', () {
    testWidgets('muestra el estado vacio cuando no hay datos', (tester) async {
      final source = Source(total: 0);
      final controller =
          PagedListController<String>(fetcher: source.fetch, pageSize: 20);
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      expect(find.text('Vacio'), findsOneWidget);
      expect(find.text('No hay nada'), findsOneWidget);
    });

    testWidgets('el pie dice cuantos se ven de cuantos hay', (tester) async {
      final source = Source(total: 800);
      final controller =
          PagedListController<String>(fetcher: source.fetch, pageSize: 3);
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      expect(find.text('fila-0'), findsOneWidget);
      expect(find.text('3 de 800 filas'), findsOneWidget);
    });

    testWidgets('el pie omite el total cuando ya no hay mas', (tester) async {
      final source = Source(total: 4);
      final controller =
          PagedListController<String>(fetcher: source.fetch, pageSize: 20);
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      expect(find.text('4 filas'), findsOneWidget);
    });

    testWidgets('scrollear hasta el final pide la pagina siguiente',
        (tester) async {
      final source = Source(total: 500);
      final controller =
          PagedListController<String>(fetcher: source.fetch, pageSize: 20);
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();
      expect(source.requested, [1]);

      await tester.drag(find.byType(ListView), const Offset(0, -2000));
      await tester.pumpAndSettle();

      expect(source.requested, contains(2));
      expect(controller.itemCount, greaterThan(20));
    });

    testWidgets('solo construye los items visibles, no los 500',
        (tester) async {
      final source = Source(total: 500);
      final controller =
          PagedListController<String>(fetcher: source.fetch, pageSize: 20);
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      final construidos = find.textContaining('fila-').evaluate().length;
      expect(construidos, lessThan(20));
    });

    testWidgets('una posicion expulsada se pinta como esqueleto y se repone',
        (tester) async {
      final source = Source(total: 5000);
      final controller = PagedListController<String>(
        fetcher: source.fetch,
        pageSize: 20,
        maxPagesInMemory: 3,
      );
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      for (var i = 0; i < 12; i++) {
        await tester.drag(find.byType(ListView), const Offset(0, -1200));
        await tester.pumpAndSettle();
      }

      expect(controller.pagesInMemory, lessThanOrEqualTo(3));
      expect(controller.itemCount, greaterThan(60));
      expect(controller.isHole(0), isTrue);

      await tester.drag(find.byType(ListView), const Offset(0, 20000));
      await tester.pumpAndSettle();

      expect(controller.itemAt(0), 'fila-0');
      expect(find.text('fila-0'), findsOneWidget);
    });

    testWidgets('el scroll no salta cuando se expulsan paginas',
        (tester) async {
      final source = Source(total: 5000);
      final controller = PagedListController<String>(
        fetcher: source.fetch,
        pageSize: 20,
        maxPagesInMemory: 3,
      );
      await controller.refresh();

      await tester.pumpWidget(_app(controller));
      await tester.pumpAndSettle();

      for (var i = 0; i < 8; i++) {
        await tester.drag(find.byType(ListView), const Offset(0, -1200));
        await tester.pumpAndSettle();
      }

      final scrollable = tester.widget<ListView>(find.byType(ListView));
      final position = scrollable.controller!.position;
      final antes = position.pixels;

      await tester.pump();
      await tester.pumpAndSettle();

      expect(position.pixels, antes);
      expect(controller.pagesInMemory, lessThanOrEqualTo(3));
    });
  });
}
