import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/app/sync_providers.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/sync_entities.dart';

void main() {
  group('SyncRunRecord', () {
    test('parses an inventory run and maps it to a result', () {
      final record = SyncRunRecord.fromJson({
        'integration_id': 7,
        'kind': 'inventory',
        'status': 'completed',
        'finished_at': '2026-08-22T10:00:00Z',
        'total': 120,
        'updated': 80,
        'unchanged': 30,
        'skipped': 5,
        'failed': 5,
      });

      expect(record.integrationId, 7);
      expect(record.kind, SyncRunKind.inventory);
      expect(record.finishedOn, isNotNull);

      final result = record.toResult();
      expect(result, isA<InventorySyncResult>());
      expect((result as InventorySyncResult).updated, 80);
      expect(result.failed, 5);
    });

    test('parses a products run and exposes the pending count', () {
      final record = SyncRunRecord.fromJson({
        'integration_id': 9,
        'kind': 'products',
        'matched': 40,
        'not_associated': 3,
        'only_in_channel': 2,
        'channel_no_sku': 1,
        'sku_changed': 4,
        'sku_typo': 0,
      });

      final result = record.toResult() as ProductsSyncResult;
      expect(result.matched, 40);
      expect(result.pending, 10);
    });

    test('tolerates numbers arriving as strings or decimals', () {
      final record = SyncRunRecord.fromJson({
        'integration_id': '12',
        'kind': 'inventory',
        'total': 10.0,
        'updated': '4',
      });

      expect(record.integrationId, 12);
      expect(record.total, 10);
      expect(record.updated, 4);
    });

    test('an unknown kind falls back to inventory', () {
      final record = SyncRunRecord.fromJson({'integration_id': 1, 'kind': 'raro'});
      expect(record.kind, SyncRunKind.inventory);
    });
  });

  group('SyncStartResult', () {
    test('reads the correlation id from the envelope root', () {
      final result = SyncStartResult.fromJson({
        'success': true,
        'correlation_id': 'abc-123',
      });
      expect(result.tracked, isTrue);
      expect(result.correlationId, 'abc-123');
    });

    test('reads the correlation id nested under data', () {
      final result = SyncStartResult.fromJson({
        'success': true,
        'data': {'correlation_id': 'nested-1'},
      });
      expect(result.correlationId, 'nested-1');
      expect(result.tracked, isTrue);
    });

    test('a success without correlation id is not trackable', () {
      final result = SyncStartResult.fromJson({'success': true});
      expect(result.success, isTrue);
      expect(result.tracked, isFalse);
    });

    test('a null body is a failure', () {
      expect(SyncStartResult.fromJson(null).success, isFalse);
    });
  });

  group('buildApplyMessage', () {
    test('joins only the buckets that have items', () {
      expect(buildApplyMessage(10, 4, 6, 0, null), '4 creados, 6 actualizados de 10');
    });

    test('reports nothing to apply when the total is zero', () {
      expect(buildApplyMessage(0, 0, 0, 0, null), 'No habia productos por aplicar');
    });

    test('reports no changes when there were items but no movement', () {
      expect(buildApplyMessage(5, 0, 0, 0, null), 'Sin cambios');
    });

    test('appends the first error and how many more failed', () {
      final message = buildApplyMessage(3, 1, 0, 2, [
        {'sku': 'A-1', 'error': 'sin stock'},
        {'sku': 'A-2', 'error': 'sin precio'},
      ]);
      expect(message, contains('2 con error'));
      expect(message, contains('A-1: sin stock'));
      expect(message, contains('y 1 mas'));
    });
  });

  group('SyncProgress', () {
    test('a total of zero leaves the ratio undefined so the bar spins', () {
      expect(const SyncProgress().ratio, isNull);
      expect(const SyncProgress(processed: 5, total: 10).ratio, 0.5);
    });
  });

  group('SyncRunItemsQuery', () {
    test('omits the optional filters when they are empty', () {
      final params = const SyncRunItemsQuery(
        integrationId: 3,
        kind: SyncRunKind.products,
      ).toQueryParams();

      expect(params['kind'], 'products');
      expect(params.containsKey('group'), isFalse);
      expect(params.containsKey('q'), isFalse);
      expect(params.containsKey('business_id'), isFalse);
    });

    test('includes the filters when they carry a value', () {
      final params = const SyncRunItemsQuery(
        integrationId: 3,
        kind: SyncRunKind.inventory,
        group: 'failed',
        search: 'abc',
        businessId: 11,
      ).toQueryParams();

      expect(params['group'], 'failed');
      expect(params['q'], 'abc');
      expect(params['business_id'], 11);
    });
  });

  group('syncProviders', () {
    test('resolves a provider from an int or a string type id', () {
      expect(syncProviderFor(4)?.key, 'woocommerce');
      expect(syncProviderFor('4')?.key, 'woocommerce');
      expect(syncProviderFor(null), isNull);
      expect(syncProviderFor(999), isNull);
    });

    test('woocommerce keeps its distinct product event prefix', () {
      final woo = syncProviders[4]!;
      expect(woo.inventoryEventPrefix, 'woo');
      expect(woo.productPrefix, 'woocommerce');
      expect(woo.syncInventoryPath, '/woocommerce/inventory/sync');
    });

    test('siigo reconciles through its own start path', () {
      expect(syncProviders[8]!.reconcileProductsPath, '/siigo/products/reconcile/start');
    });

    test('only the channels that support it expose compare inventory', () {
      final withCompare = syncProviders.values
          .where((spec) => spec.supportsCompareInventory)
          .map((spec) => spec.key)
          .toSet();
      expect(withCompare, {'meli', 'woocommerce', 'siigo', 'tiendanube'});
    });

    test('apply bodies carry the direction and mode each backend expects', () {
      expect(syncProviders[4]!.applyBodyFor('createInChannel').toJson(), {'direction': 'to_woo'});
      expect(syncProviders[16]!.applyBodyFor('updateInProbability').toJson(),
          {'direction': 'to_probability', 'mode': 'update'});
      expect(syncProviders[8]!.applyBodyFor('createInProbability').toJson(), isEmpty);
    });

    test('siigo cannot associate and only applies into Probability', () {
      final siigo = syncProviders[8]!;
      expect(siigo.supportsAssociate, isFalse);
      expect(siigo.supportsApply('createInProbability'), isTrue);
      expect(siigo.supportsApply('createInChannel'), isFalse);
    });

    test('the subscribed event list covers every provider prefix', () {
      final types = globalSyncEventTypes;
      expect(types, contains('woo.inventory.sync.completed'));
      expect(types, contains('woocommerce.product.sync.completed'));
      expect(types, contains('meli.product.reconcile.completed'));
      expect(types.toSet().length, greaterThan(40));
    });

    test('orders compare covers the channels the backend accepts', () {
      expect(ordersCompareTypeIds, containsAll(<int>[1, 3, 4, 16, 17, 33]));
      expect(ordersCompareTypeIds, isNot(contains(8)));
    });
  });
}
