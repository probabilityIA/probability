import 'package:flutter/foundation.dart';

import '../../../../../core/network/api_client.dart';
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
  List<MatrixColumn> _columns = const <MatrixColumn>[];
  final Set<int> _presentIn = <int>{};
  final Set<int> _missingIn = <int>{};

  List<MatrixColumn> get columns => _columns;
  String get search => _search;
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
