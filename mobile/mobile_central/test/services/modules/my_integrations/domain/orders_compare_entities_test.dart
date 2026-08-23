import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/orders_compare_entities.dart';

void main() {
  group('OrderCompareRow', () {
    test('parses a row that only exists in the channel', () {
      final row = OrderCompareRow.fromJson({
        'external_id': 'MLC-1',
        'number': '1001',
        'customer_name': 'Ana',
        'channel_status': 'paid',
        'raw_status': 'paid',
        'total': 125000,
        'items': 2,
        'created_at': '2026-08-01T10:00:00Z',
        'action': 'create',
        'moves_inventory': true,
      });

      expect(row.canCreate, isTrue);
      expect(row.label, '1001');
      expect(row.localStatus, isNull);
      expect(row.amount, 125000);
      expect(row.hasMismatch, isFalse);
    });

    test('falls back to the external id when there is no number', () {
      final row = OrderCompareRow.fromJson({'external_id': 'X-9', 'action': 'create'});
      expect(row.label, 'X-9');
    });

    test('an unknown action is treated as in sync, never as creatable', () {
      final row = OrderCompareRow.fromJson({'external_id': 'X', 'action': 'raro'});
      expect(row.action, OrderCompareAction.inSync);
      expect(row.canCreate, isFalse);
    });

    test('uses the local total when the channel total is zero', () {
      final row = OrderCompareRow.fromJson({
        'external_id': 'X',
        'total': 0,
        'local_total': 44000,
        'action': 'only_in_probability',
      });
      expect(row.amount, 44000);
    });

    test('a placeholder date is discarded instead of shown as year 1', () {
      final row = OrderCompareRow.fromJson({
        'external_id': 'X',
        'created_at': '0001-01-01T00:00:00Z',
      });
      expect(row.createdOn, isNull);
    });

    test('reports a mismatch when the status or the total differ', () {
      final status = OrderCompareRow.fromJson({
        'external_id': 'X',
        'action': 'in_sync',
        'status_mismatch': true,
      });
      final total = OrderCompareRow.fromJson({
        'external_id': 'Y',
        'action': 'in_sync',
        'total_mismatch': true,
      });
      expect(status.hasMismatch, isTrue);
      expect(total.hasMismatch, isTrue);
    });
  });

  group('OrderCompareTotals', () {
    test('what the channel has is what is in sync plus what is missing', () {
      final totals = OrderCompareTotals.fromJson({'in_sync': 8, 'to_create': 3});
      expect(totals.inChannel, 11);
    });

    test('a missing totals block is all zeros, not a crash', () {
      expect(OrderCompareTotals.fromJson(null).inChannel, 0);
    });
  });

  group('OrdersComparePage', () {
    test('parses rows and totals, and defaults the page when absent', () {
      final page = OrdersComparePage.fromJson({
        'rows': [
          {'external_id': 'A', 'action': 'create'},
          {'external_id': 'B', 'action': 'in_sync'},
        ],
        'totals': {'to_create': 1, 'in_sync': 1},
        'total': 2,
        'total_pages': 1,
        'checked_at': '2026-08-22T12:00:00Z',
      });

      expect(page.rows.length, 2);
      expect(page.page, 1);
      expect(page.pageSize, 20);
      expect(page.totals.toCreate, 1);
      expect(page.checkedAt, isNotNull);
    });

    test('an empty payload yields an empty page', () {
      final page = OrdersComparePage.fromJson(<String, dynamic>{});
      expect(page.rows, isEmpty);
      expect(page.total, 0);
    });
  });

  group('OrdersCompareQuery', () {
    test('only sends only_diff when it is on', () {
      final off = const OrdersCompareQuery(integrationId: 1).toQueryParams();
      final on = const OrdersCompareQuery(integrationId: 1, onlyDiff: true).toQueryParams();
      expect(off.containsKey('only_diff'), isFalse);
      expect(on['only_diff'], 'true');
    });

    test('drops the empty dates and the empty search', () {
      final params = const OrdersCompareQuery(
        integrationId: 1,
        from: '',
        to: '',
        search: '',
      ).toQueryParams();
      expect(params.containsKey('from'), isFalse);
      expect(params.containsKey('to'), isFalse);
      expect(params.containsKey('q'), isFalse);
    });
  });

  group('OrdersApplyResult', () {
    test('summarizes what was queued, skipped and failed', () {
      final result = OrdersApplyResult.fromJson({
        'queued': ['A', 'B'],
        'skipped': ['C'],
        'failed': {'D': 'sin stock'},
      });
      expect(result.summary, contains('2 ordenes enviadas a crear'));
      expect(result.summary, contains('1 ya existian'));
      expect(result.summary, contains('1 fallaron'));
      expect(result.failed['D'], 'sin stock');
    });

    test('says plainly when nothing was created', () {
      expect(const OrdersApplyResult().summary, 'No se creo ninguna orden');
    });

    test('uses the singular for a single order', () {
      final result = OrdersApplyResult.fromJson({'queued': ['A']});
      expect(result.summary, contains('1 orden enviada a crear'));
    });
  });
}
