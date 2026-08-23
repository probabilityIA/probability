import 'package:flutter/foundation.dart';

import '../../../../../core/network/api_client.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../shared/filters/filter_models.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/saved_comparison_use_cases.dart';
import '../../app/sync_providers.dart';
import '../../domain/entities.dart';
import '../../domain/saved_comparison_entities.dart';
import '../../infra/repository/saved_comparison_repository.dart';

enum ChannelFilterState { off, present, missing }

Pagination _pagination({
  required int page,
  required int perPage,
  required int total,
  required int totalPages,
}) {
  final lastPage = totalPages < 1 ? 1 : totalPages;
  return Pagination(
    currentPage: page,
    perPage: perPage,
    total: total,
    lastPage: lastPage,
    hasNext: page < lastPage,
    hasPrev: page > 1,
  );
}

class ProductMatrixProvider extends ChangeNotifier {
  ProductMatrixProvider({
    required ApiClient apiClient,
    SavedComparisonUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<MatrixRow>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final SavedComparisonUseCases? _injectedUseCases;
  late final PagedListController<MatrixRow> list;

  int? _businessId;
  String _search = '';
  String _searchBy = 'all';
  List<MatrixColumn> _columns = const <MatrixColumn>[];
  final Set<int> _presentIn = <int>{};
  final Set<int> _missingIn = <int>{};

  List<MatrixColumn> get columns => _columns;
  String get search => _search;
  String get searchBy => _searchBy;
  Set<int> get presentIn => Set.unmodifiable(_presentIn);
  Set<int> get missingIn => Set.unmodifiable(_missingIn);
  bool get hasFilters => _presentIn.isNotEmpty || _missingIn.isNotEmpty;

  ChannelFilterState stateFor(int integrationId) {
    if (_presentIn.contains(integrationId)) return ChannelFilterState.present;
    if (_missingIn.contains(integrationId)) return ChannelFilterState.missing;
    return ChannelFilterState.off;
  }

  FilterSelection get selection {
    final values = <String, String>{};
    for (final id in _presentIn) {
      values['ch_$id'] = 'present';
    }
    for (final id in _missingIn) {
      values['ch_$id'] = 'missing';
    }
    return FilterSelection(values);
  }

  Future<void> applySelection(FilterSelection next) {
    _presentIn.clear();
    _missingIn.clear();
    next.values.forEach((key, value) {
      final id = int.tryParse(key.replaceFirst('ch_', ''));
      if (id == null) return;
      if (value == 'present') {
        _presentIn.add(id);
      } else if (value == 'missing') {
        _missingIn.add(id);
      }
    });
    return list.refresh();
  }

  Future<void> clearFilters() {
    if (!hasFilters) return Future<void>.value();
    _presentIn.clear();
    _missingIn.clear();
    return list.refresh();
  }

  SavedComparisonUseCases get _useCases =>
      _injectedUseCases ??
      SavedComparisonUseCases(SavedComparisonApiRepository(_apiClient));

  Future<PaginatedResponse<MatrixRow>> _fetchPage(int page, int pageSize) async {
    final result = await _useCases.getMatchMatrix(
      businessId: _businessId,
      page: page,
      pageSize: pageSize,
      search: _search.isEmpty ? null : _search,
      searchBy: _searchBy,
      presentIn: _presentIn.toList(),
      missingIn: _missingIn.toList(),
    );
    if (result.columns.isNotEmpty) _columns = result.columns;

    return PaginatedResponse(
      data: result.rows,
      pagination: _pagination(
        page: result.page,
        perPage: pageSize,
        total: result.total,
        totalPages: result.totalPages,
      ),
    );
  }

  Future<void> load({int? businessId}) {
    _businessId = businessId;
    return list.refresh();
  }

  Future<void> setSearch(String value) {
    final next = value.trim();
    if (next == _search) return Future<void>.value();
    _search = next;
    return list.refresh();
  }

  Future<void> setSearchBy(String value) {
    if (value == _searchBy) return Future<void>.value();
    _searchBy = value;
    if (_search.isEmpty) {
      notifyListeners();
      return Future<void>.value();
    }
    return list.refresh();
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}

class FindingItemsProvider extends ChangeNotifier {
  FindingItemsProvider({
    required ApiClient apiClient,
    SavedComparisonUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<FindingItem>(fetcher: _fetchPage, pageSize: 50);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final SavedComparisonUseCases? _injectedUseCases;
  late final PagedListController<FindingItem> list;

  int? _businessId;
  String _code = '';
  String _search = '';

  String get code => _code;
  String get search => _search;

  SavedComparisonUseCases get _useCases =>
      _injectedUseCases ??
      SavedComparisonUseCases(SavedComparisonApiRepository(_apiClient));

  Future<PaginatedResponse<FindingItem>> _fetchPage(
    int page,
    int pageSize,
  ) async {
    if (_code.isEmpty) {
      return PaginatedResponse(
        data: const <FindingItem>[],
        pagination: _pagination(page: 1, perPage: pageSize, total: 0, totalPages: 1),
      );
    }

    final result = await _useCases.getFindingItems(
      code: _code,
      businessId: _businessId,
      page: page,
      pageSize: pageSize,
      search: _search.isEmpty ? null : _search,
    );

    return PaginatedResponse(
      data: result.items,
      pagination: _pagination(
        page: result.page,
        perPage: pageSize,
        total: result.total,
        totalPages: result.totalPages,
      ),
    );
  }

  Future<void> load({required String code, int? businessId}) {
    _code = code;
    _businessId = businessId;
    _search = '';
    return list.refresh();
  }

  Future<void> setSearch(String value) {
    final next = value.trim();
    if (next == _search) return Future<void>.value();
    _search = next;
    return list.refresh();
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}

class InventoryRowsProvider extends ChangeNotifier {
  InventoryRowsProvider({
    required ApiClient apiClient,
    SavedComparisonUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<InventoryCompareRow>(
      fetcher: _fetchPage,
      pageSize: 50,
    );
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final SavedComparisonUseCases? _injectedUseCases;
  late final PagedListController<InventoryCompareRow> list;

  int? _businessId;
  MyIntegration? _channel;
  bool _onlyDiff = false;
  bool _live = false;
  String _search = '';
  InventoryCompareTotals _totals = const InventoryCompareTotals();
  DateTime? _checkedAt;
  bool _fromCache = true;

  MyIntegration? get channel => _channel;
  bool get onlyDiff => _onlyDiff;
  bool get live => _live;
  String get search => _search;
  InventoryCompareTotals get totals => _totals;
  DateTime? get checkedAt => _checkedAt;
  bool get fromCache => _fromCache;

  SavedComparisonUseCases get _useCases =>
      _injectedUseCases ??
      SavedComparisonUseCases(SavedComparisonApiRepository(_apiClient));

  Future<PaginatedResponse<InventoryCompareRow>> _fetchPage(
    int page,
    int pageSize,
  ) async {
    final channel = _channel;
    final spec = channel == null ? null : syncProviderFor(channel.integrationTypeId);
    if (channel == null || spec == null) {
      return PaginatedResponse(
        data: const <InventoryCompareRow>[],
        pagination: _pagination(page: 1, perPage: pageSize, total: 0, totalPages: 1),
      );
    }

    final result = await _useCases.compareInventory(
      spec,
      InventoryCompareQuery(
        integrationId: channel.id,
        businessId: _businessId,
        page: page,
        pageSize: pageSize,
        snapshot: !_live,
        onlyDiff: _onlyDiff,
        search: _search.isEmpty ? null : _search,
      ),
    );

    _totals = result.totals;
    _checkedAt = result.checkedAt;
    _fromCache = result.fromCache;

    return PaginatedResponse(
      data: result.rows,
      pagination: _pagination(
        page: result.page,
        perPage: pageSize,
        total: result.total,
        totalPages: result.totalPages,
      ),
    );
  }

  Future<void> load({
    required MyIntegration channel,
    int? businessId,
    bool live = false,
  }) {
    _channel = channel;
    _businessId = businessId;
    _live = live;
    return list.refresh();
  }

  Future<void> toggleOnlyDiff() {
    _onlyDiff = !_onlyDiff;
    return list.refresh();
  }

  Future<void> setSearch(String value) {
    final next = value.trim();
    if (next == _search) return Future<void>.value();
    _search = next;
    return list.refresh();
  }

  Future<void> askChannel() {
    _live = true;
    return list.refresh();
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}

class InventoryMatrixProvider extends ChangeNotifier {
  InventoryMatrixProvider({
    required ApiClient apiClient,
    SavedComparisonUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<MatrixRow>(fetcher: _fetchPage, pageSize: 20);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final SavedComparisonUseCases? _injectedUseCases;
  late final PagedListController<MatrixRow> list;

  int? _businessId;
  String _search = '';
  String _searchBy = 'all';
  List<MyIntegration> _channels = const <MyIntegration>[];
  List<MatrixColumn> _columns = const <MatrixColumn>[];
  final Set<int> _presentIn = <int>{};
  final Set<int> _missingIn = <int>{};

  final Map<int, Map<String, InventoryCompareRow>> _stock =
      <int, Map<String, InventoryCompareRow>>{};
  final Map<int, DateTime?> _checkedAt = <int, DateTime?>{};
  final Set<int> _live = <int>{};

  bool _asking = false;
  String? _error;

  List<MatrixColumn> get columns => _columns;
  List<MyIntegration> get channels => _channels;
  String get search => _search;
  String get searchBy => _searchBy;
  bool get asking => _asking;
  String? get error => _error ?? list.error;
  bool get hasFilters => _presentIn.isNotEmpty || _missingIn.isNotEmpty;

  DateTime? get lastCheck {
    DateTime? latest;
    for (final when in _checkedAt.values) {
      if (when == null) continue;
      if (latest == null || when.isAfter(latest)) latest = when;
    }
    return latest;
  }

  InventoryCompareRow? stockFor(int integrationId, String sku) =>
      _stock[integrationId]?[sku];

  bool isLive(int integrationId) => _live.contains(integrationId);

  ChannelFilterState stateFor(int integrationId) {
    if (_presentIn.contains(integrationId)) return ChannelFilterState.present;
    if (_missingIn.contains(integrationId)) return ChannelFilterState.missing;
    return ChannelFilterState.off;
  }

  FilterSelection get selection {
    final values = <String, String>{};
    for (final id in _presentIn) {
      values['ch_$id'] = 'present';
    }
    for (final id in _missingIn) {
      values['ch_$id'] = 'missing';
    }
    return FilterSelection(values);
  }

  SavedComparisonUseCases get _useCases =>
      _injectedUseCases ??
      SavedComparisonUseCases(SavedComparisonApiRepository(_apiClient));

  Future<PaginatedResponse<MatrixRow>> _fetchPage(int page, int pageSize) async {
    final result = await _useCases.getMatchMatrix(
      businessId: _businessId,
      page: page,
      pageSize: pageSize,
      search: _search.isEmpty ? null : _search,
      searchBy: _searchBy,
      presentIn: _presentIn.toList(),
      missingIn: _missingIn.toList(),
    );
    if (result.columns.isNotEmpty) _columns = result.columns;

    await _loadStock(result.rows);

    final lastPage = result.totalPages < 1 ? 1 : result.totalPages;
    return PaginatedResponse(
      data: result.rows,
      pagination: _pagination(
        page: result.page,
        perPage: pageSize,
        total: result.total,
        totalPages: lastPage,
      ),
    );
  }

  Future<void> _loadStock(List<MatrixRow> rows, {Set<int>? live}) async {
    final skus = rows.map((row) => row.sku).where((sku) => sku.isNotEmpty).toList();
    if (skus.isEmpty || _channels.isEmpty) return;

    await Future.wait(_channels.map((channel) async {
      final spec = syncProviderFor(channel.integrationTypeId);
      if (spec == null || !spec.supportsCompareInventory) return;

      final askChannel = live?.contains(channel.id) == true;
      try {
        final page = await _useCases.compareInventory(
          spec,
          InventoryCompareQuery(
            integrationId: channel.id,
            businessId: _businessId,
            page: 1,
            pageSize: skus.length,
            snapshot: !askChannel,
            skus: skus,
          ),
        );
        final map = _stock.putIfAbsent(
          channel.id,
          () => <String, InventoryCompareRow>{},
        );
        for (final row in page.rows) {
          if (row.sku.isEmpty) continue;
          map[row.sku] = row;
        }
        _checkedAt[channel.id] = page.checkedAt;
        if (page.fromCache) {
          _live.remove(channel.id);
        } else {
          _live.add(channel.id);
        }
      } catch (_) {
        return;
      }
    }));
  }

  Future<void> load({
    required List<MyIntegration> integrations,
    int? businessId,
  }) {
    _channels = integrations
        .where((i) =>
            syncProviderFor(i.integrationTypeId)?.supportsCompareInventory == true)
        .toList();
    _businessId = businessId;
    return list.refresh();
  }

  Future<void> setSearch(String value) {
    final next = value.trim();
    if (next == _search) return Future<void>.value();
    _search = next;
    return list.refresh();
  }

  Future<void> setSearchBy(String value) {
    if (value == _searchBy) return Future<void>.value();
    _searchBy = value;
    if (_search.isEmpty) {
      notifyListeners();
      return Future<void>.value();
    }
    return list.refresh();
  }

  Future<void> applySelection(FilterSelection next) {
    _presentIn.clear();
    _missingIn.clear();
    next.values.forEach((key, value) {
      final id = int.tryParse(key.replaceFirst('ch_', ''));
      if (id == null) return;
      if (value == 'present') {
        _presentIn.add(id);
      } else if (value == 'missing') {
        _missingIn.add(id);
      }
    });
    return list.refresh();
  }

  Future<void> askChannels() async {
    if (_asking || _channels.isEmpty) return;
    _asking = true;
    _error = null;
    notifyListeners();

    try {
      await _loadStock(
        list.loadedItems,
        live: _channels.map((channel) => channel.id).toSet(),
      );
    } catch (e) {
      _error = parseError(e);
    } finally {
      _asking = false;
      notifyListeners();
    }
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
