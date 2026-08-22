import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/pagination/paged_collection.dart';

List<String> pageOf(int page, int size, int total) {
  final start = (page - 1) * size;
  final end = (start + size) > total ? total : (start + size);
  return [for (var i = start; i < end; i++) 'item-$i'];
}

void main() {
  group('PagedCollection', () {
    test('acumula paginas y respeta el total', () {
      final c = PagedCollection<String>(pageSize: 20, maxPagesInMemory: 8);
      c.setPage(1, pageOf(1, 20, 55), 55);
      c.setPage(2, pageOf(2, 20, 55), 55);

      expect(c.length, 40);
      expect(c.total, 55);
      expect(c.hasMore, isTrue);
      expect(c.peek(0), 'item-0');
      expect(c.peek(39), 'item-39');
    });

    test('la ultima pagina cierra hasMore y no infla la longitud', () {
      final c = PagedCollection<String>(pageSize: 20, maxPagesInMemory: 8);
      c.setPage(1, pageOf(1, 20, 25), 25);
      c.setPage(2, pageOf(2, 20, 25), 25);

      expect(c.length, 25);
      expect(c.hasMore, isFalse);
    });

    test('expulsa por LRU sin correr los indices', () {
      final c = PagedCollection<String>(pageSize: 10, maxPagesInMemory: 3);
      for (var page = 1; page <= 5; page++) {
        c.setPage(page, pageOf(page, 10, 100), 100);
      }

      c.at(45);
      final evicted = c.evictColdPages(protectedPages: {5});

      expect(evicted, isNotEmpty);
      expect(c.pagesInMemory, lessThanOrEqualTo(3));
      expect(c.length, 50);
      expect(c.peek(49), 'item-49');
    });

    test('el hueco queda nulo y la posicion se conserva', () {
      final c = PagedCollection<String>(pageSize: 10, maxPagesInMemory: 2);
      for (var page = 1; page <= 4; page++) {
        c.setPage(page, pageOf(page, 10, 100), 100);
      }
      c.at(35);
      c.evictColdPages(protectedPages: {4});

      expect(c.length, 40);
      expect(c.isHole(0), isTrue);
      expect(c.peek(39), 'item-39');
      expect(c.evictedPages, contains(1));
    });

    test('recargar una pagina expulsada la vuelve a llenar en su lugar', () {
      final c = PagedCollection<String>(pageSize: 10, maxPagesInMemory: 2);
      for (var page = 1; page <= 4; page++) {
        c.setPage(page, pageOf(page, 10, 100), 100);
      }
      c.evictColdPages();
      expect(c.isHole(0), isTrue);

      c.setPage(1, pageOf(1, 10, 100), 100);

      expect(c.peek(0), 'item-0');
      expect(c.isHole(0), isFalse);
      expect(c.length, 40);
    });

    test('nunca supera el tope de items vivos', () {
      final c = PagedCollection<String>(pageSize: 20, maxPagesInMemory: 8);
      for (var page = 1; page <= 500; page++) {
        c.setPage(page, pageOf(page, 20, 100000), 100000);
        c.at(((page - 1) * 20) + 5);
        c.evictColdPages(protectedPages: {page});
      }

      expect(c.length, 10000);
      expect(c.liveItemCount, lessThanOrEqualTo(8 * 20));
      expect(c.pagesInMemory, lessThanOrEqualTo(8));
    });

    test('reset limpia todo', () {
      final c = PagedCollection<String>(pageSize: 10, maxPagesInMemory: 2);
      c.setPage(1, pageOf(1, 10, 100), 100);
      c.reset();

      expect(c.length, 0);
      expect(c.total, 0);
      expect(c.isEmpty, isTrue);
      expect(c.hasMore, isFalse);
    });
  });
}
