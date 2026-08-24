import 'package:flutter/foundation.dart';

import '../../../../../core/errors/error_parser.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/orders_compare_use_cases.dart';
import '../../domain/orders_compare_entities.dart';
import '../../infra/repository/orders_compare_repository.dart';

class OrdersCompareProvider extends ChangeNotifier {
  OrdersCompareProvider({
    required ApiClient apiClient,
    OrdersCompareUseCases? useCases,
  })  : _apiClient = apiClient,
        _injectedUseCases = useCases {
    rows = PagedListController<OrderCompareRow>(fetcher: _fetchPage);
    rows.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final OrdersCompareUseCases? _injectedUseCases;
  late final PagedListController<OrderCompareRow> rows;

  int? _businessId;
  int? _integrationId;
  String _from = '';
  String _to = '';
  bool _onlyDiff = true;
  String _search = '';

  bool _started = false;
  bool _applying = false;
  String? _error;
  OrderCompareTotals _totals = const OrderCompareTotals();
  String? _checkedAt;
  OrdersApplyResult? _lastApply;
  final Set<String> _selection = <String>{};

  int? get integrationId => _integrationId;
  String get from => _from;
  String get to => _to;
  bool get onlyDiff => _onlyDiff;
  String get search => _search;
  bool get started => _started;
  bool get applying => _applying;
  String? get error => _error ?? rows.error;
  OrderCompareTotals get totals => _totals;
  String? get checkedAt => _checkedAt;
  OrdersApplyResult? get lastApply => _lastApply;
  Set<String> get selection => Set.unmodifiable(_selection);
  bool get hasSelection => _selection.isNotEmpty;

  bool isSelected(String externalId) => _selection.contains(externalId);

  OrdersCompareUseCases get _useCases =>
      _injectedUseCases ??
      OrdersCompareUseCases(OrdersCompareApiRepository(_apiClient));

  Future<PaginatedResponse<OrderCompareRow>> _fetchPage(
    int page,
    int pageSize,
  ) async {
    final id = _integrationId;
    if (id == null) {
      return PaginatedResponse(
        data: const <OrderCompareRow>[],
        pagination: Pagination(
          currentPage: 1,
          perPage: 20,
          total: 0,
          lastPage: 1,
          hasNext: false,
          hasPrev: false,
        ),
      );
    }

    final result = await _useCases.compare(OrdersCompareQuery(
      integrationId: id,
      businessId: _businessId,
      from: _from.isEmpty ? null : _from,
      to: _to.isEmpty ? null : _to,
      page: page,
      pageSize: pageSize,
      onlyDiff: _onlyDiff,
      search: _search.isEmpty ? null : _search,
    ));

    _totals = result.totals;
    _checkedAt = result.checkedAt;

    final lastPage = result.totalPages < 1 ? 1 : result.totalPages;

    return PaginatedResponse(
      data: result.rows,
      pagination: Pagination(
        currentPage: result.page,
        perPage: result.pageSize,
        total: result.total,
        lastPage: lastPage,
        hasNext: result.page < lastPage,
        hasPrev: result.page > 1,
      ),
    );
  }

  void configure({int? businessId}) {
    if (_businessId == businessId) return;
    _businessId = businessId;
    _integrationId = null;
    _started = false;
    _selection.clear();
    _lastApply = null;
    _totals = const OrderCompareTotals();
    notifyListeners();
  }

  void selectChannel(int integrationId) {
    if (_integrationId == integrationId) return;
    _integrationId = integrationId;
    _started = false;
    _selection.clear();
    _lastApply = null;
    _error = null;
    _totals = const OrderCompareTotals();
    notifyListeners();
  }

  void setRange({String? from, String? to}) {
    _from = from ?? _from;
    _to = to ?? _to;
    notifyListeners();
  }

  void toggleOnlyDiff() {
    _onlyDiff = !_onlyDiff;
    notifyListeners();
  }

  void setSearch(String value) {
    _search = value.trim();
    notifyListeners();
  }

  void toggle(String externalId) {
    if (!_selection.remove(externalId)) _selection.add(externalId);
    notifyListeners();
  }

  void toggleAllLoaded() {
    final creatable = rows.loadedItems
        .where((row) => row.canCreate)
        .map((row) => row.externalId)
        .toList();
    if (creatable.isEmpty) return;

    final allSelected = creatable.every(_selection.contains);
    if (allSelected) {
      _selection.removeAll(creatable);
    } else {
      _selection.addAll(creatable);
    }
    notifyListeners();
  }

  Future<void> compare() async {
    if (_integrationId == null) return;
    _started = true;
    _error = null;
    _lastApply = null;
    _selection.clear();
    notifyListeners();
    await rows.refresh();
  }

  Future<bool> apply() async {
    final id = _integrationId;
    if (id == null || _selection.isEmpty || _applying) return false;

    _applying = true;
    _error = null;
    notifyListeners();

    try {
      _lastApply = await _useCases.apply(
        id,
        _selection.toList(),
        businessId: _businessId,
      );
      _selection.clear();
      return true;
    } catch (e) {
      _error = parseError(e);
      return false;
    } finally {
      _applying = false;
      notifyListeners();
    }
  }

  @override
  void dispose() {
    rows.removeListener(notifyListeners);
    rows.dispose();
    super.dispose();
  }
}
