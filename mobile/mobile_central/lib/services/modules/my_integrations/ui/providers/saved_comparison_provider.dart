import 'package:flutter/foundation.dart';

import '../../../../../core/errors/error_parser.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/saved_comparison_use_cases.dart';
import '../../app/sync_providers.dart';
import '../../domain/entities.dart';
import '../../domain/saved_comparison_entities.dart';
import '../../infra/repository/saved_comparison_repository.dart';

class ChannelSnapshot {
  const ChannelSnapshot({
    required this.integrationId,
    this.page,
    this.error,
    this.loading = false,
  });

  final int integrationId;
  final InventoryComparePage? page;
  final String? error;
  final bool loading;

  bool get hasData => page != null;

  bool get isSaved => page?.fromCache == true;

  int get toUpdate => page?.totals.toUpdate ?? 0;

  int get total => page?.total ?? 0;

  DateTime? get checkedAt => page?.checkedAt;
}

class SavedComparisonProvider extends ChangeNotifier {
  SavedComparisonProvider({
    required ApiClient apiClient,
    SavedComparisonUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases;

  final ApiClient _apiClient;
  final SavedComparisonUseCases? _injectedUseCases;

  int? _businessId;

  FindingsReport _findings = const FindingsReport();
  bool _loadingFindings = false;
  String? _findingsError;

  DataSummary _dataSummary = const DataSummary();
  bool _loadingData = false;
  String? _dataError;

  final Map<int, ChannelSnapshot> _snapshots = <int, ChannelSnapshot>{};

  FindingsReport get findings => _findings;
  bool get loadingFindings => _loadingFindings;
  String? get findingsError => _findingsError;

  DataSummary get dataSummary => _dataSummary;
  bool get loadingData => _loadingData;
  String? get dataError => _dataError;

  ChannelSnapshot snapshotFor(int integrationId) =>
      _snapshots[integrationId] ?? ChannelSnapshot(integrationId: integrationId);

  bool get anySnapshotLoading =>
      _snapshots.values.any((snapshot) => snapshot.loading);

  DateTime? get lastInventoryCheck {
    DateTime? latest;
    for (final snapshot in _snapshots.values) {
      final when = snapshot.checkedAt;
      if (when == null) continue;
      if (latest == null || when.isAfter(latest)) latest = when;
    }
    return latest;
  }

  int get inventoryToUpdate =>
      _snapshots.values.fold(0, (sum, snapshot) => sum + snapshot.toUpdate);

  SavedComparisonUseCases get _useCases =>
      _injectedUseCases ??
      SavedComparisonUseCases(SavedComparisonApiRepository(_apiClient));

  void configure({int? businessId}) {
    if (_businessId == businessId) return;
    _businessId = businessId;
    _findings = const FindingsReport();
    _dataSummary = const DataSummary();
    _snapshots.clear();
    _findingsError = null;
    _dataError = null;
    notifyListeners();
  }

  Future<void> loadFindings({bool force = false}) async {
    if (_loadingFindings) return;
    if (!force && !_findings.isEmpty) return;

    _loadingFindings = true;
    _findingsError = null;
    notifyListeners();

    try {
      _findings = await _useCases.getFindings(businessId: _businessId);
    } catch (e) {
      _findingsError = parseError(e);
    } finally {
      _loadingFindings = false;
      notifyListeners();
    }
  }

  Future<void> loadDataSummary({bool force = false}) async {
    if (_loadingData) return;
    if (!force && !_dataSummary.isEmpty) return;

    _loadingData = true;
    _dataError = null;
    notifyListeners();

    try {
      _dataSummary = await _useCases.getDataSummary(businessId: _businessId);
    } catch (e) {
      _dataError = parseError(e);
    } finally {
      _loadingData = false;
      notifyListeners();
    }
  }

  Future<void> loadInventorySnapshots(
    List<MyIntegration> integrations, {
    bool force = false,
  }) async {
    final targets = integrations
        .where((i) => syncProviderFor(i.integrationTypeId)?.supportsCompareInventory == true)
        .toList();

    await Future.wait(targets.map((integration) =>
        loadInventorySnapshot(integration, force: force)));
  }

  Future<void> loadInventorySnapshot(
    MyIntegration integration, {
    bool force = false,
    bool live = false,
  }) async {
    final spec = syncProviderFor(integration.integrationTypeId);
    if (spec == null || !spec.supportsCompareInventory) return;

    final current = _snapshots[integration.id];
    if (current != null && current.loading) return;
    if (!force && !live && current?.hasData == true) return;

    _snapshots[integration.id] = ChannelSnapshot(
      integrationId: integration.id,
      page: current?.page,
      loading: true,
    );
    notifyListeners();

    try {
      final page = await _useCases.compareInventory(
        spec,
        InventoryCompareQuery(
          integrationId: integration.id,
          businessId: _businessId,
          snapshot: !live,
        ),
      );
      _snapshots[integration.id] = ChannelSnapshot(
        integrationId: integration.id,
        page: page,
      );
    } catch (e) {
      _snapshots[integration.id] = ChannelSnapshot(
        integrationId: integration.id,
        page: current?.page,
        error: parseError(e),
      );
    } finally {
      notifyListeners();
    }
  }
}
