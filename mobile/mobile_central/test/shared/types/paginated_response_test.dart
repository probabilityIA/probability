import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

void main() {
  group('Pagination.fromEnvelope', () {
    test('lee la forma plana que devuelve produccion', () {
      final p = Pagination.fromEnvelope(const {
        'success': true,
        'page': 1,
        'page_size': 20,
        'total': 47073,
        'total_pages': 2354,
        'data': <dynamic>[],
      });

      expect(p.total, 47073);
      expect(p.currentPage, 1);
      expect(p.perPage, 20);
      expect(p.lastPage, 2354);
      expect(p.hasNext, isTrue);
      expect(p.hasPrev, isFalse);
    });

    test('la ultima pagina de la forma plana cierra hasNext', () {
      final p = Pagination.fromEnvelope(const {
        'page': 2354,
        'page_size': 20,
        'total': 47073,
        'total_pages': 2354,
      });

      expect(p.hasNext, isFalse);
      expect(p.hasPrev, isTrue);
    });

    test('sigue leyendo la forma anidada', () {
      final p = Pagination.fromEnvelope(const {
        'pagination': {
          'current_page': 3,
          'per_page': 20,
          'total': 100,
          'last_page': 5,
          'has_next': true,
          'has_prev': true,
        },
      });

      expect(p.total, 100);
      expect(p.currentPage, 3);
      expect(p.lastPage, 5);
      expect(p.hasNext, isTrue);
    });

    test('deduce last_page cuando no viene', () {
      final p = Pagination.fromEnvelope(const {
        'page': 1,
        'page_size': 20,
        'total': 55,
      });

      expect(p.lastPage, 3);
      expect(p.hasNext, isTrue);
    });

    test('un sobre sin datos de paginacion no explota', () {
      final p = Pagination.fromEnvelope(const {'data': <dynamic>[]});

      expect(p.total, 0);
      expect(p.currentPage, 1);
      expect(p.hasNext, isFalse);
    });
  });
}
