import 'dart:async';

import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/core/network/sse_client.dart';
import 'package:mobile_central/services/modules/my_integrations/app/sync_providers.dart';
import 'package:mobile_central/services/modules/my_integrations/app/sync_use_cases.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/entities.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/ports.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/sync_entities.dart';
import 'package:mobile_central/services/modules/my_integrations/ui/providers/sync_activity_provider.dart';

class FakeSse implements SseSource {
  final StreamController<SseEvent> _controller = StreamController<SseEvent>.broadcast();

  int connections = 0;

  @override
  Stream<SseEvent> get events => _controller.stream;

  @override
  Future<void> connect({required int businessId, required List<String> eventTypes}) async {
    connections++;
  }

  void emit(String type, Map<String, dynamic> data) {
    _controller.add(SseEvent(type: type, data: <String, dynamic>{'type': type, 'data': data}));
  }

  @override
  void disconnect() {}

  @override
  void dispose() {
    _controller.close();
  }
}

class FakeSyncRuns implements ISyncRunsRepository {
  FakeSyncRuns({this.correlationId = 'corr-1', this.startSucceeds = true});

  String correlationId;
  bool startSucceeds;
  List<SyncRunRecord> runs = <SyncRunRecord>[];

  final List<String> calls = <String>[];
  List<String>? lastSkus;
  String? lastApplyAction;

  @override
  Future<List<SyncRunRecord>> listLastRuns({int? businessId}) async {
    calls.add('listLastRuns');
    return runs;
  }

  @override
  Future<List<SyncRunDetail>> listRunItems(SyncRunItemsQuery query) async => const [];

  @override
  Future<SyncStartResult> syncInventory(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) async {
    calls.add('syncInventory');
    lastSkus = skus;
    return _start();
  }

  @override
  Future<SyncStartResult> reconcileProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
  }) async {
    calls.add('reconcileProducts');
    return _start();
  }

  @override
  Future<SyncStartResult> associateProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) async {
    calls.add('associateProducts');
    return _start();
  }

  @override
  Future<SyncStartResult> applyProducts(
    SyncProviderSpec spec,
    int integrationId,
    String action, {
    int? businessId,
    List<String>? skus,
  }) async {
    calls.add('applyProducts');
    lastApplyAction = action;
    return _start();
  }

  SyncStartResult _start() => startSucceeds
      ? SyncStartResult(success: true, correlationId: correlationId)
      : const SyncStartResult(success: false, message: 'No se pudo iniciar');
}

MyIntegration _woo({int id = 1, bool active = true}) => MyIntegration(
      id: id,
      createdAt: '',
      updatedAt: '',
      businessId: 2,
      integrationTypeId: 4,
      name: 'Mi tienda',
      isActive: active,
    );

MyIntegration _unsupported() => MyIntegration(
      id: 99,
      createdAt: '',
      updatedAt: '',
      businessId: 2,
      integrationTypeId: 404,
      name: 'Canal raro',
      isActive: true,
    );

SyncActivityProvider _build(FakeSyncRuns repo, FakeSse sse) => SyncActivityProvider(
      apiClient: ApiClient(),
      useCases: SyncRunsUseCases(repo),
      sseClient: sse,
    );

void main() {
  group('eligible', () {
    test('leaves out the inactive and the channels with no provider', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);

      await provider.configure(
        integrations: [_woo(id: 1), _woo(id: 2, active: false), _unsupported()],
        businessId: 2,
      );

      expect(provider.eligible.map((i) => i.id), [1]);
      expect(sse.connections, 1);
      provider.dispose();
    });
  });

  group('runInventoryOne', () {
    test('finishes on the completed event and keeps the counters', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runInventoryOne(1);
      await Future<void>.delayed(const Duration(milliseconds: 10));

      sse.emit('woo.inventory.sync.completed', {
        'correlation_id': 'corr-1',
        'total': 50,
        'updated': 30,
        'unchanged': 15,
        'skipped': 3,
        'failed': 2,
      });
      await run;

      final result = provider.resultFor(1) as InventorySyncResult;
      expect(result.updated, 30);
      expect(result.failed, 2);
      expect(provider.stateFor(1), SyncNodeState.done);
      expect(provider.running, isFalse);
      provider.dispose();
    });

    test('progress events feed the bar', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runInventoryOne(1);
      await Future<void>.delayed(const Duration(milliseconds: 10));

      sse.emit('woo.inventory.sync.started', {'correlation_id': 'corr-1', 'total': 40});
      sse.emit('woo.inventory.sync.progress', {'correlation_id': 'corr-1', 'processed': 10});
      await Future<void>.delayed(const Duration(milliseconds: 10));

      expect(provider.progressFor(1).processed, 10);
      expect(provider.progressFor(1).total, 40);
      expect(provider.progressFor(1).ratio, 0.25);

      sse.emit('woo.inventory.sync.completed', {'correlation_id': 'corr-1', 'total': 40});
      await run;
      provider.dispose();
    });

    test('an item event becomes a detail line with its tone', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runInventoryOne(1);
      await Future<void>.delayed(const Duration(milliseconds: 10));

      sse.emit('woo.inventory.sync.item', {
        'correlation_id': 'corr-1',
        'sku': 'A-1',
        'action': 'failed',
        'error': 'sin stock',
      });
      sse.emit('woo.inventory.sync.item', {
        'correlation_id': 'corr-1',
        'sku': 'A-2',
        'action': 'unchanged',
        'quantity': 4,
      });
      await Future<void>.delayed(const Duration(milliseconds: 10));

      final details = provider.detailsFor(1);
      expect(details.length, 2);
      expect(details.first.tone, 'error');
      expect(details.first.label, 'sin stock');
      expect(details.last.tone, 'warn');
      expect(details.last.group, 'skipped');

      sse.emit('woo.inventory.sync.completed', {'correlation_id': 'corr-1', 'total': 2});
      await run;
      provider.dispose();
    });

    test('a start that fails leaves the error without waiting for the timeout', () async {
      final repo = FakeSyncRuns(startSucceeds: false);
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      await provider.runInventoryOne(1);

      expect(provider.stateFor(1), SyncNodeState.error);
      expect((provider.resultFor(1) as SyncErrorResult).message, 'No se pudo iniciar');
      provider.dispose();
    });

    test('forwards the selected skus to the channel', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runInventoryOne(1, skus: ['A-1', 'A-2']);
      await Future<void>.delayed(const Duration(milliseconds: 10));
      sse.emit('woo.inventory.sync.completed', {'correlation_id': 'corr-1', 'total': 2});
      await run;

      expect(repo.lastSkus, ['A-1', 'A-2']);
      provider.dispose();
    });
  });

  group('events that arrive before the correlation id is known', () {
    test('are buffered and replayed, not lost', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      sse.emit('woo.inventory.sync.completed', {
        'correlation_id': 'corr-1',
        'total': 12,
        'updated': 12,
      });
      await Future<void>.delayed(const Duration(milliseconds: 10));

      await provider.runInventoryOne(1);

      final result = provider.resultFor(1) as InventorySyncResult;
      expect(result.updated, 12);
      provider.dispose();
    });
  });

  group('runProducts', () {
    test('a reconcile completed fills the matched and pending counters', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runProducts();
      await Future<void>.delayed(const Duration(milliseconds: 10));
      expect(provider.stateFor(1), SyncNodeState.scan);

      sse.emit('woo.product.reconcile.completed', {
        'correlation_id': 'corr-1',
        'matched': 25,
        'not_associated': 2,
        'only_in_channel': 1,
      });
      await run;

      final result = provider.resultFor(1) as ProductsSyncResult;
      expect(result.matched, 25);
      expect(result.pending, 3);
      expect(provider.stateFor(1), SyncNodeState.done);
      provider.dispose();
    });

    test('a reconcile that reports an error marks the node as failed', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runProducts();
      await Future<void>.delayed(const Duration(milliseconds: 10));
      sse.emit('woo.product.reconcile.completed', {
        'correlation_id': 'corr-1',
        'error': 'credenciales vencidas',
      });
      await run;

      expect(provider.stateFor(1), SyncNodeState.error);
      expect((provider.resultFor(1) as SyncErrorResult).message, 'credenciales vencidas');
      provider.dispose();
    });
  });

  group('runProductAction', () {
    test('skips an action the channel does not support', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);
      repo.calls.clear();

      await provider.runProductAction(1, ProductActionKey.updateInProbability);

      expect(repo.calls, isEmpty);
      provider.dispose();
    });

    test('createBothSides runs the two directions in order', () async {
      final repo = FakeSyncRuns(correlationId: '');
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);
      repo.calls.clear();

      await provider.runProductAction(1, ProductActionKey.createBothSides);

      expect(repo.calls.where((c) => c == 'applyProducts').length, 2);
      expect(repo.lastApplyAction, 'createInProbability');
      provider.dispose();
    });
  });

  group('loadLastRuns', () {
    test('indexes the runs by integration and kind', () async {
      final repo = FakeSyncRuns();
      repo.runs = [
        SyncRunRecord.fromJson({'integration_id': 1, 'kind': 'inventory', 'updated': 4}),
        SyncRunRecord.fromJson({'integration_id': 1, 'kind': 'products', 'matched': 9}),
      ];
      final sse = FakeSse();
      final provider = _build(repo, sse);

      await provider.configure(integrations: [_woo()], businessId: 2);

      expect(provider.lastRunFor(1, SyncRunKind.inventory)?.updated, 4);
      expect(provider.lastRunFor(1, SyncRunKind.products)?.matched, 9);
      expect(provider.lastRunFor(2, SyncRunKind.inventory), isNull);
      provider.dispose();
    });
  });

  group('reset', () {
    test('clears the live state but keeps the last runs', () async {
      final repo = FakeSyncRuns();
      repo.runs = [
        SyncRunRecord.fromJson({'integration_id': 1, 'kind': 'inventory', 'updated': 4}),
      ];
      final sse = FakeSse();
      final provider = _build(repo, sse);
      await provider.configure(integrations: [_woo()], businessId: 2);

      final run = provider.runInventoryOne(1);
      await Future<void>.delayed(const Duration(milliseconds: 10));
      sse.emit('woo.inventory.sync.completed', {'correlation_id': 'corr-1', 'total': 3});
      await run;

      expect(provider.finished, isTrue);
      provider.reset();

      expect(provider.resultFor(1), isNull);
      expect(provider.stateFor(1), SyncNodeState.idle);
      expect(provider.finished, isFalse);
      expect(provider.lastRunFor(1, SyncRunKind.inventory)?.updated, 4);
      provider.dispose();
    });
  });

  group('environment', () {
    test('only inventory and products can be run from the toolbar', () async {
      final repo = FakeSyncRuns();
      final sse = FakeSse();
      final provider = _build(repo, sse);

      expect(provider.environment, SyncEnvironment.overview);
      expect(provider.canRun, isFalse);

      provider.setEnvironment(SyncEnvironment.inventory);
      expect(provider.canRun, isTrue);

      provider.setEnvironment(SyncEnvironment.ordersCompare);
      expect(provider.canRun, isFalse);
      provider.dispose();
    });
  });
}
