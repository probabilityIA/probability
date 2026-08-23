import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../../../../core/network/api_client.dart';
import '../../../../../core/network/sse_client.dart';
import '../../app/sync_providers.dart';
import '../../app/sync_use_cases.dart';
import '../../domain/entities.dart';
import '../../domain/sync_entities.dart';
import '../../infra/repository/sync_runs_repository.dart';

const Duration syncTimeout = Duration(seconds: 90);
const Duration applyTimeout = Duration(minutes: 10);
const Duration applySettle = Duration(seconds: 4);
const Duration runClockSkew = Duration(seconds: 60);
const int maxPendingRuns = 20;
const int maxPendingEventsPerRun = 100;
const int maxDetailsPerIntegration = 400;

class _BufferedEvent {
  const _BufferedEvent(this.type, this.data);

  final String type;
  final Map<String, dynamic> data;
}

class SyncActivityProvider extends ChangeNotifier {
  SyncActivityProvider({
    required ApiClient apiClient,
    SyncRunsUseCases? useCases,
    SseSource? sseClient,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases,
        _sse = sseClient ?? SseClient();

  final ApiClient _apiClient;
  final SyncRunsUseCases? _injectedUseCases;
  final SseSource _sse;

  StreamSubscription<SseEvent>? _subscription;

  List<MyIntegration> _integrations = const <MyIntegration>[];
  int? _businessId;

  SyncMode _mode = SyncMode.idle;
  bool _running = false;
  SyncEnvironment _environment = SyncEnvironment.overview;

  final Map<int, SyncNodeState> _nodes = <int, SyncNodeState>{};
  final Map<int, SyncProgress> _progress = <int, SyncProgress>{};
  final Map<int, SyncResult> _results = <int, SyncResult>{};
  final Map<int, List<SyncRunDetail>> _details = <int, List<SyncRunDetail>>{};
  final Map<int, Map<SyncRunKind, SyncRunRecord>> _lastRuns =
      <int, Map<SyncRunKind, SyncRunRecord>>{};
  final Map<int, ProductActionKey?> _actionBusy = <int, ProductActionKey?>{};
  final Map<int, ProductActionResult?> _actionResult = <int, ProductActionResult?>{};

  final Map<String, int> _corrToIntegration = <String, int>{};
  final Map<String, List<_BufferedEvent>> _pending = <String, List<_BufferedEvent>>{};
  final Map<int, Completer<void>> _completion = <int, Completer<void>>{};
  final Map<int, Completer<void>> _applyCompletion = <int, Completer<void>>{};
  final Map<int, Timer> _timers = <int, Timer>{};

  SyncMode get mode => _mode;
  bool get running => _running;
  SyncEnvironment get environment => _environment;
  int? get businessId => _businessId;
  Map<int, SyncNodeState> get nodes => Map.unmodifiable(_nodes);
  Map<int, SyncProgress> get progress => Map.unmodifiable(_progress);
  Map<int, SyncResult> get results => Map.unmodifiable(_results);

  SyncNodeState stateFor(int id) => _nodes[id] ?? SyncNodeState.idle;
  SyncProgress progressFor(int id) => _progress[id] ?? const SyncProgress();
  SyncResult? resultFor(int id) => _results[id];
  List<SyncRunDetail> detailsFor(int id) => _details[id] ?? const <SyncRunDetail>[];
  SyncRunRecord? lastRunFor(int id, SyncRunKind kind) => _lastRuns[id]?[kind];
  ProductActionKey? actionBusyFor(int id) => _actionBusy[id];
  ProductActionResult? actionResultFor(int id) => _actionResult[id];

  bool get canRun =>
      _environment == SyncEnvironment.inventory || _environment == SyncEnvironment.products;

  bool get finished =>
      !_running &&
      _nodes.isNotEmpty &&
      _nodes.values.every((s) => s == SyncNodeState.done || s == SyncNodeState.error);

  List<MyIntegration> get eligible => _integrations
      .where((i) => i.isActive && syncProviderFor(i.integrationTypeId) != null)
      .toList();

  SyncRunsUseCases get _useCases =>
      _injectedUseCases ?? SyncRunsUseCases(SyncRunsApiRepository(_apiClient));

  void setEnvironment(SyncEnvironment value) {
    if (_environment == value) return;
    _environment = value;
    notifyListeners();
  }

  Future<void> configure({
    required List<MyIntegration> integrations,
    int? businessId,
  }) async {
    final changedBusiness = _businessId != businessId;
    _integrations = integrations;
    _businessId = businessId;

    if (changedBusiness) {
      reset();
      _lastRuns.clear();
      await _listen();
    } else if (_subscription == null) {
      await _listen();
    }

    await loadLastRuns();
  }

  Future<void> _listen() async {
    await _subscription?.cancel();
    _subscription = _sse.events.listen(_onEvent);
    await _sse.connect(
      businessId: _businessId ?? 0,
      eventTypes: globalSyncEventTypes,
    );
  }

  Future<void> loadLastRuns() async {
    try {
      final rows = await _useCases.listLastRuns(businessId: _businessId);
      _lastRuns.clear();
      for (final row in rows) {
        _lastRuns.putIfAbsent(row.integrationId, () => <SyncRunKind, SyncRunRecord>{});
        _lastRuns[row.integrationId]![row.kind] = row;
      }
      notifyListeners();
    } catch (_) {
      return;
    }
  }

  void _onEvent(SseEvent event) {
    final data = event.data['data'];
    if (data is! Map<String, dynamic>) return;
    final correlationId = data['correlation_id']?.toString() ?? '';
    if (correlationId.isEmpty) return;

    final id = _corrToIntegration[correlationId];
    if (id == null) {
      _buffer(correlationId, event.type, data);
      return;
    }
    _apply(id, event.type, data);
  }

  void _buffer(String correlationId, String type, Map<String, dynamic> data) {
    final list = _pending.putIfAbsent(correlationId, () => <_BufferedEvent>[]);
    if (list.length >= maxPendingEventsPerRun) return;
    list.add(_BufferedEvent(type, data));
    while (_pending.length > maxPendingRuns) {
      _pending.remove(_pending.keys.first);
    }
  }

  void _drain(String correlationId) {
    final buffered = _pending.remove(correlationId);
    if (buffered == null) return;
    final id = _corrToIntegration[correlationId];
    if (id == null) return;
    for (final item in buffered) {
      _apply(id, item.type, item.data);
    }
  }

  void _apply(int id, String type, Map<String, dynamic> data) {
    if (type.endsWith('.inventory.sync.started')) {
      _setNode(id, SyncNodeState.active);
      _progress[id] = SyncProgress(total: _int(data['total']));
    } else if (type.endsWith('.inventory.sync.item')) {
      _pushDetail(id, _detailFromItem(data));
    } else if (type.endsWith('.inventory.sync.progress')) {
      _progress[id] = SyncProgress(
        processed: _int(data['processed']),
        total: _int(data['total']) > 0 ? _int(data['total']) : progressFor(id).total,
      );
    } else if (type.endsWith('.product.sync.started')) {
      _setNode(id, SyncNodeState.active);
      _progress[id] = SyncProgress(total: _int(data['total']));
    } else if (type.endsWith('.product.sync.progress')) {
      _progress[id] = SyncProgress(
        processed: _int(data['processed']),
        total: _int(data['total']) > 0 ? _int(data['total']) : progressFor(id).total,
      );
    } else if (type.endsWith('.product.sync.completed')) {
      _onApplyCompleted(id, data);
      return;
    } else if (type.endsWith('.product.reconcile.started')) {
      _setNode(id, SyncNodeState.scan);
    } else if (type.endsWith('.product.reconcile.completed')) {
      _onReconcileCompleted(id, data);
      return;
    } else if (type.endsWith('.inventory.sync.completed')) {
      _onInventoryCompleted(id, data);
      return;
    } else {
      return;
    }
    notifyListeners();
  }

  void _onApplyCompleted(int id, Map<String, dynamic> data) {
    final created = _int(data['created']);
    final updated = _int(data['updated']);
    final failed = _int(data['failed']);
    final total = _int(data['total']) > 0 ? _int(data['total']) : created + updated + failed;
    final failedItems = data['failed_items'];

    _actionResult[id] = ProductActionResult(
      ok: failed == 0 && total > 0,
      message: buildApplyMessage(
        total,
        created,
        updated,
        failed,
        failedItems is List ? failedItems : null,
      ),
    );
    _nodes[id] = failed > 0 ? SyncNodeState.error : SyncNodeState.done;
    _finish(_applyCompletion, id);
    notifyListeners();
  }

  void _onReconcileCompleted(int id, Map<String, dynamic> data) {
    final error = data['error'];
    if (error != null && error.toString().isNotEmpty) {
      _nodes[id] = SyncNodeState.error;
      _results[id] = SyncErrorResult(error.toString());
    } else {
      _results[id] = ProductsSyncResult(
        matched: _int(data['matched']),
        notAssociated: _int(data['not_associated']),
        onlyInProbability: _int(data['only_in_probability']),
        onlyInChannel: _int(data['only_in_channel']),
        channelNoSku: _int(data['channel_no_sku']),
        skuChanged: _int(data['sku_changed']),
        skuTypo: _int(data['sku_typo']),
      );
      _nodes[id] = SyncNodeState.done;
    }
    _finish(_completion, id);
    notifyListeners();
  }

  void _onInventoryCompleted(int id, Map<String, dynamic> data) {
    final total = _int(data['total']);
    _progress[id] = SyncProgress(processed: total, total: total);
    _results[id] = InventorySyncResult(
      total: total,
      updated: _int(data['updated']),
      unchanged: _int(data['unchanged']),
      skipped: _int(data['skipped']),
      failed: _int(data['failed']),
    );

    final failedSkus = data['failed_skus'];
    if (failedSkus is List) {
      for (final raw in failedSkus) {
        if (raw is String) {
          _pushDetail(
            id,
            SyncRunDetail(
              sku: raw.isEmpty ? '(sin sku)' : raw,
              label: 'fallo al actualizar',
              tone: 'error',
              group: 'failed',
            ),
          );
        } else if (raw is Map) {
          final sku = raw['sku']?.toString() ?? '';
          _pushDetail(
            id,
            SyncRunDetail(
              sku: sku.isEmpty ? '(sin sku)' : sku,
              label: raw['error']?.toString() ?? 'fallo al actualizar',
              tone: 'error',
              group: 'failed',
            ),
          );
        }
      }
    }

    _nodes[id] = SyncNodeState.done;
    _finish(_completion, id);
    notifyListeners();
  }

  SyncRunDetail _detailFromItem(Map<String, dynamic> data) {
    final action = data['action']?.toString() ?? '';
    final failed = RegExp('fail|error', caseSensitive: false).hasMatch(action);
    final skipped = RegExp('skip|omit|unchanged', caseSensitive: false).hasMatch(action);
    final sku = data['sku']?.toString() ?? '';
    final quantity = data['quantity']?.toString() ?? '-';

    return SyncRunDetail(
      sku: sku.isEmpty ? '(sin sku)' : sku,
      label: failed
          ? (data['error'] ?? data['message'] ?? 'fallo al actualizar').toString()
          : '${action.isEmpty ? 'actualizado' : action} - $quantity u.',
      tone: failed ? 'error' : (skipped ? 'warn' : 'ok'),
      group: failed ? 'failed' : (skipped ? 'skipped' : 'updated'),
    );
  }

  void _pushDetail(int id, SyncRunDetail detail) {
    final list = _details.putIfAbsent(id, () => <SyncRunDetail>[]);
    if (list.length > maxDetailsPerIntegration) return;
    list.add(detail);
  }

  void _setNode(int id, SyncNodeState state) {
    if (_nodes[id] == state) return;
    _nodes[id] = state;
  }

  void _finish(Map<int, Completer<void>> registry, int id) {
    _timers.remove(id)?.cancel();
    final completer = registry.remove(id);
    if (completer != null && !completer.isCompleted) completer.complete();
  }

  Future<void> _await(
    Map<int, Completer<void>> registry,
    int id,
    String correlationId,
    Duration timeout,
    Future<void> Function() onTimeout,
  ) async {
    final completer = Completer<void>();
    registry[id] = completer;
    _timers[id] = Timer(timeout, () async {
      registry.remove(id);
      await onTimeout();
      if (!completer.isCompleted) completer.complete();
    });
    _drain(correlationId);
    await completer.future;
    _timers.remove(id)?.cancel();
  }

  Future<bool> _hydrateFromLastRun(int id, SyncRunKind kind, DateTime launchedAt) async {
    try {
      final rows = await _useCases.listLastRuns(businessId: _businessId);
      SyncRunRecord? row;
      for (final candidate in rows) {
        if (candidate.integrationId == id && candidate.kind == kind) {
          row = candidate;
          break;
        }
      }
      final finishedAt = row?.finishedOn;
      if (row == null || finishedAt == null) return false;
      if (finishedAt.isBefore(launchedAt.subtract(runClockSkew))) return false;

      _results[id] = row.toResult();
      if (kind == SyncRunKind.inventory) {
        _progress[id] = SyncProgress(processed: row.total, total: row.total);
      }
      _nodes[id] = row.failed > 0 ? SyncNodeState.error : SyncNodeState.done;
      notifyListeners();
      return true;
    } catch (_) {
      return false;
    }
  }

  Future<void> _syncInventoryOne(MyIntegration integration, {List<String>? skus}) async {
    final spec = syncProviderFor(integration.integrationTypeId);
    if (spec == null) return;

    _setNode(integration.id, SyncNodeState.active);
    notifyListeners();

    final launchedAt = DateTime.now();
    final start = await _useCases.syncInventory(
      spec,
      integration.id,
      businessId: _businessId,
      skus: skus,
    );

    if (!start.tracked) {
      _nodes[integration.id] = SyncNodeState.error;
      _results[integration.id] = SyncErrorResult(start.message ?? 'No se pudo iniciar');
      notifyListeners();
      return;
    }

    final correlationId = start.correlationId!;
    _corrToIntegration[correlationId] = integration.id;

    await _await(_completion, integration.id, correlationId, syncTimeout, () async {
      final recovered =
          await _hydrateFromLastRun(integration.id, SyncRunKind.inventory, launchedAt);
      if (!recovered) {
        _nodes[integration.id] = SyncNodeState.error;
        _results[integration.id] = const SyncErrorResult('Continua en segundo plano');
        notifyListeners();
      }
    });
  }

  Future<void> _reconcileOne(MyIntegration integration) async {
    final spec = syncProviderFor(integration.integrationTypeId);
    if (spec == null) return;

    _setNode(integration.id, SyncNodeState.scan);
    _details[integration.id] = <SyncRunDetail>[];
    notifyListeners();

    final launchedAt = DateTime.now();
    final start = await _useCases.reconcileProducts(
      spec,
      integration.id,
      businessId: _businessId,
    );

    if (!start.tracked) {
      _nodes[integration.id] = SyncNodeState.error;
      _results[integration.id] = SyncErrorResult(start.message ?? 'No se pudo comparar');
      notifyListeners();
      return;
    }

    final correlationId = start.correlationId!;
    _corrToIntegration[correlationId] = integration.id;

    await _await(_completion, integration.id, correlationId, syncTimeout, () async {
      final recovered =
          await _hydrateFromLastRun(integration.id, SyncRunKind.products, launchedAt);
      if (!recovered) {
        _nodes[integration.id] = SyncNodeState.error;
        _results[integration.id] = const SyncErrorResult('Continua en segundo plano');
        notifyListeners();
      }
    });
  }

  Future<void> runInventoryOne(int integrationId, {List<String>? skus}) async {
    if (_running) return;
    MyIntegration? integration;
    for (final item in eligible) {
      if (item.id == integrationId) integration = item;
    }
    if (integration == null) return;

    _running = true;
    _mode = SyncMode.inventory;
    _results.remove(integrationId);
    _details[integrationId] = <SyncRunDetail>[];
    _progress[integrationId] = const SyncProgress();
    notifyListeners();

    await _syncInventoryOne(integration, skus: skus);

    _running = false;
    _mode = SyncMode.idle;
    notifyListeners();
    await loadLastRuns();
  }

  Future<void> runInventory() async {
    final targets = eligible;
    if (_running || targets.isEmpty) return;

    _running = true;
    _mode = SyncMode.inventory;
    _results.clear();
    _progress.clear();
    _details.clear();
    _corrToIntegration.clear();
    for (final integration in targets) {
      _nodes[integration.id] = SyncNodeState.queued;
    }
    notifyListeners();

    for (final integration in targets) {
      await _syncInventoryOne(integration);
    }

    _running = false;
    _mode = SyncMode.idle;
    notifyListeners();
    await loadLastRuns();
  }

  Future<void> runProducts() async {
    final targets = eligible;
    if (_running || targets.isEmpty) return;

    _running = true;
    _mode = SyncMode.products;
    _results.clear();
    _progress.clear();
    _details.clear();
    _corrToIntegration.clear();
    for (final integration in targets) {
      _nodes[integration.id] = SyncNodeState.scan;
    }
    notifyListeners();

    await Future.wait(targets.map(_reconcileOne));

    _running = false;
    _mode = SyncMode.idle;
    notifyListeners();
    await loadLastRuns();
  }

  Future<void> runCurrent() async {
    if (_environment == SyncEnvironment.products) return runProducts();
    if (_environment == SyncEnvironment.inventory) return runInventory();
  }

  Future<void> runProductAction(
    int integrationId,
    ProductActionKey action, {
    List<String>? skus,
  }) async {
    if (_actionBusy[integrationId] != null) return;

    MyIntegration? integration;
    for (final item in eligible) {
      if (item.id == integrationId) integration = item;
    }
    final spec = integration == null ? null : syncProviderFor(integration.integrationTypeId);
    if (integration == null || spec == null) return;

    final steps = action == ProductActionKey.createBothSides
        ? <ProductActionKey>[
            ProductActionKey.createInChannel,
            ProductActionKey.createInProbability,
          ]
        : <ProductActionKey>[action];

    final runnable = steps.where((step) {
      if (step == ProductActionKey.associate) return spec.supportsAssociate;
      return spec.supportsApply(step.name);
    }).toList();
    if (runnable.isEmpty) return;

    _actionBusy[integrationId] = action;
    _actionResult[integrationId] = null;
    notifyListeners();

    var ok = true;
    var message = '';

    try {
      for (final step in runnable) {
        final result = step == ProductActionKey.associate
            ? await _useCases.associateProducts(
                spec,
                integrationId,
                businessId: _businessId,
                skus: skus,
              )
            : await _useCases.applyProducts(
                spec,
                integrationId,
                step.name,
                businessId: _businessId,
                skus: skus,
              );

        if (!result.success) {
          ok = false;
          message = result.message ?? 'No se pudo aplicar';
          break;
        }
        message = result.message ?? 'Aplicado';

        final correlationId = result.correlationId ?? '';
        if (correlationId.isEmpty) continue;

        _corrToIntegration[correlationId] = integrationId;
        _actionResult[integrationId] = const ProductActionResult(
          ok: true,
          pending: true,
          message: 'Aplicando...',
        );
        notifyListeners();

        await _await(_applyCompletion, integrationId, correlationId, applyTimeout, () async {
          _actionResult[integrationId] = const ProductActionResult(
            ok: true,
            pending: true,
            message: 'Sigue aplicandose en segundo plano',
          );
          notifyListeners();
        });
      }

      final current = _actionResult[integrationId];
      if (current == null || current.message == 'Aplicando...') {
        _actionResult[integrationId] = ProductActionResult(ok: ok, message: message);
      }

      if (ok) {
        await Future<void>.delayed(applySettle);
        await _reconcileOne(integration);
      }
    } catch (_) {
      _actionResult[integrationId] =
          const ProductActionResult(ok: false, message: 'No se pudo aplicar');
    } finally {
      _actionBusy[integrationId] = null;
      notifyListeners();
      await loadLastRuns();
    }
  }

  void reset() {
    for (final timer in _timers.values) {
      timer.cancel();
    }
    _timers.clear();
    _pending.clear();
    _corrToIntegration.clear();
    _completion.clear();
    _applyCompletion.clear();
    _mode = SyncMode.idle;
    _running = false;
    _nodes.clear();
    _progress.clear();
    _results.clear();
    _details.clear();
    _actionResult.clear();
    notifyListeners();
  }

  static int _int(dynamic value) {
    if (value is int) return value;
    if (value is num) return value.round();
    if (value is String) return int.tryParse(value) ?? 0;
    return 0;
  }

  @override
  void dispose() {
    for (final timer in _timers.values) {
      timer.cancel();
    }
    _timers.clear();
    _subscription?.cancel();
    _sse.dispose();
    super.dispose();
  }
}
