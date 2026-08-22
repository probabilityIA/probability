import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/saved_comparison_entities.dart';

void main() {
  group('FindingChannelSummary', () {
    test('pending adds up every bucket that is not a clean match', () {
      final channel = FindingChannelSummary.fromJson({
        'integration_id': 4,
        'integration_name': 'WooCommerce',
        'matched': 100,
        'not_associated': 2,
        'only_in_channel': 3,
        'channel_no_sku': 1,
        'sku_changed': 4,
        'sku_typo': 5,
        'compared_at': '2026-08-20T10:00:00Z',
      });

      expect(channel.pending, 15);
      expect(channel.isClean, isFalse);
      expect(channel.comparedAt, isNotNull);
    });

    test('a channel with only matches is clean', () {
      final channel = FindingChannelSummary.fromJson({
        'integration_id': 1,
        'matched': 40,
      });
      expect(channel.isClean, isTrue);
      expect(channel.comparedAt, isNull);
    });

    test('a placeholder compared_at is treated as never compared', () {
      final channel = FindingChannelSummary.fromJson({
        'integration_id': 1,
        'compared_at': '0001-01-01T00:00:00Z',
      });
      expect(channel.comparedAt, isNull);
    });
  });

  group('FindingsReport', () {
    test('reports the most recent comparison across channels', () {
      final report = FindingsReport.fromJson({
        'total': 2,
        'channels': [
          {'integration_id': 1, 'compared_at': '2026-08-10T10:00:00Z'},
          {'integration_id': 2, 'compared_at': '2026-08-21T10:00:00Z'},
          {'integration_id': 3},
        ],
      });

      expect(report.lastComparedAt, isNotNull);
      expect(report.lastComparedAt!.toUtc().day, 21);
    });

    test('a null payload is an empty report, not a crash', () {
      final report = FindingsReport.fromJson(null);
      expect(report.isEmpty, isTrue);
      expect(report.lastComparedAt, isNull);
    });

    test('parses findings with their severity', () {
      final report = FindingsReport.fromJson({
        'findings': [
          {'code': 'a', 'severity': 'error', 'title': 'X', 'count': 3},
          {'code': 'b', 'severity': 'raro', 'title': 'Y', 'count': 1},
        ],
      });
      expect(report.findings.first.severity, FindingSeverity.error);
      expect(report.findings.last.severity, FindingSeverity.info);
    });
  });

  group('DataSummary', () {
    test('only the fields with something to bring are actionable', () {
      final summary = DataSummary.fromJson({
        'snapshot_at': '2026-08-19T08:00:00Z',
        'data': [
          {
            'field': 'name',
            'label': 'Nombre',
            'cells': [
              {'integration_id': 1, 'can_fill': 5, 'can_overwrite': 0},
            ],
          },
          {
            'field': 'image',
            'label': 'Imagen',
            'cells': [
              {'integration_id': 1, 'can_fill': 0, 'can_overwrite': 0},
            ],
          },
        ],
      });

      expect(summary.rows.length, 2);
      expect(summary.actionable.length, 1);
      expect(summary.actionable.first.field, 'name');
      expect(summary.snapshotAt, isNotNull);
    });

    test('totals add across the channels of a field', () {
      final row = DataSummaryRow.fromJson({
        'field': 'name',
        'cells': [
          {'integration_id': 1, 'can_fill': 5, 'can_overwrite': 2},
          {'integration_id': 2, 'can_fill': 3, 'can_overwrite': 1},
        ],
      });
      expect(row.totalFill, 8);
      expect(row.totalOverwrite, 3);
    });

    test('a null payload is empty', () {
      expect(DataSummary.fromJson(null).isEmpty, isTrue);
    });
  });

  group('InventoryComparePage', () {
    test('marks a saved read apart from a live one', () {
      final saved = InventoryComparePage.fromJson({
        'rows': [],
        'from_cache': true,
        'checked_at': '2026-08-18T09:00:00Z',
      });
      final live = InventoryComparePage.fromJson({'rows': []});

      expect(saved.fromCache, isTrue);
      expect(saved.checkedAt, isNotNull);
      expect(live.fromCache, isFalse);
    });

    test('parses rows and keeps the null quantities as null', () {
      final page = InventoryComparePage.fromJson({
        'rows': [
          {
            'product_id': 'p1',
            'sku': 'A-1',
            'name': 'Camiseta',
            'action': 'update',
            'probability_qty': 10,
            'channel_qty': 4,
            'delta': -6,
          },
          {'product_id': 'p2', 'sku': 'A-2', 'action': 'skip'},
        ],
        'totals': {'total': 2, 'to_update': 1, 'skipped': 1},
      });

      expect(page.rows.first.needsUpdate, isTrue);
      expect(page.rows.first.delta, -6);
      expect(page.rows.last.probabilityQty, isNull);
      expect(page.rows.last.action, InventoryCompareAction.skip);
      expect(page.totals.toUpdate, 1);
    });
  });

  group('InventoryCompareQuery', () {
    test('asks for the snapshot by default', () {
      final body = const InventoryCompareQuery(integrationId: 4).toBody();
      expect(body['source'], 'snapshot');
    });

    test('a live read omits the source so the backend asks the channel', () {
      final body = const InventoryCompareQuery(
        integrationId: 4,
        snapshot: false,
      ).toBody();
      expect(body.containsKey('source'), isFalse);
    });

    test('drops the empty filters', () {
      final body = const InventoryCompareQuery(
        integrationId: 4,
        search: '',
        skus: <String>[],
      ).toBody();
      expect(body.containsKey('q'), isFalse);
      expect(body.containsKey('skus'), isFalse);
      expect(body.containsKey('only_diff'), isFalse);
    });
  });
}
