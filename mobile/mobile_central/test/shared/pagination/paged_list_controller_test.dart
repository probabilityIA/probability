import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/pagination/paged_list_controller.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

class FakeSource {
  FakeSource(this.total);

  final int total;
  final List<int> requested = [];
  int failuresLeft = 0;

  Future<PaginatedResponse<String>> fetch(int page, int pageSize) async {
    requested.add(page);
    if (failuresLeft > 0) {
      failuresLeft--;
      throw StateError('boom');
    }
    final start = (page - 1) * pageSize;
    final end = (start + pageSize) > total ? total : (start + pageSize);
    final data = [for (var i = start; i < end; i++) 'item-$i'];
    return PaginatedResponse<String>(
      data: data,
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

void main() {
  group('PagedListController', () {
    test('refresh carga la primera pagina', () async {
      final source = FakeSource(120);
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);

      await c.refresh();

      expect(c.itemCount, 20);
      expect(c.total, 120);
      expect(c.hasMore, isTrue);
      expect(c.itemAt(0), 'item-0');
    });

    test('loadMore concatena y avanza de pagina', () async {
      final source = FakeSource(120);
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);

      await c.refresh();
      await c.loadMore();
      await c.loadMore();

      expect(c.itemCount, 60);
      expect(source.requested, [1, 2, 3]);
      expect(c.itemAt(59), 'item-59');
    });

    test('loadMore se detiene al llegar al final', () async {
      final source = FakeSource(30);
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);

      await c.refresh();
      await c.loadMore();
      await c.loadMore();

      expect(c.itemCount, 30);
      expect(c.hasMore, isFalse);
      expect(source.requested, [1, 2]);
    });

    test('no lanza dos veces la misma pagina', () async {
      final source = FakeSource(200);
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);
      await c.refresh();

      await Future.wait([c.loadMore(), c.loadMore(), c.loadMore()]);

      expect(source.requested, [1, 2]);
    });

    test('mantiene acotada la memoria recorriendo muchas paginas', () async {
      final source = FakeSource(20000);
      final c = PagedListController<String>(
        fetcher: source.fetch,
        pageSize: 20,
        maxPagesInMemory: 5,
      );

      await c.refresh();
      for (var i = 0; i < 60; i++) {
        c.itemAt(c.itemCount - 1);
        await Future<void>.delayed(Duration.zero);
        await c.loadMore();
      }
      await Future<void>.delayed(Duration.zero);

      expect(c.itemCount, 20 * 61);
      expect(c.liveItemCount, lessThanOrEqualTo(5 * 20));
      expect(c.pagesInMemory, lessThanOrEqualTo(5));
    });

    test('un hueco se vuelve a pedir al mirarlo', () async {
      final source = FakeSource(20000);
      final c = PagedListController<String>(
        fetcher: source.fetch,
        pageSize: 10,
        maxPagesInMemory: 2,
      );

      await c.refresh();
      for (var i = 0; i < 5; i++) {
        c.itemAt(c.itemCount - 1);
        await Future<void>.delayed(Duration.zero);
        await c.loadMore();
      }
      expect(c.isHole(0), isTrue);

      c.itemAt(0);
      await Future<void>.delayed(Duration.zero);
      await Future<void>.delayed(Duration.zero);

      expect(c.itemAt(0), 'item-0');
    });

    test('un error deja mensaje y no rompe la coleccion', () async {
      final source = FakeSource(120)..failuresLeft = 1;
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);

      await c.refresh();
      expect(c.error, isNotNull);
      expect(c.isEmpty, isTrue);

      await c.refresh();
      expect(c.error, isNull);
      expect(c.itemCount, 20);
    });

    test('refresh reinicia el recorrido', () async {
      final source = FakeSource(120);
      final c = PagedListController<String>(fetcher: source.fetch, pageSize: 20);

      await c.refresh();
      await c.loadMore();
      expect(c.itemCount, 40);

      await c.refresh();
      expect(c.itemCount, 20);
      expect(c.itemAt(0), 'item-0');
    });
  });
}
