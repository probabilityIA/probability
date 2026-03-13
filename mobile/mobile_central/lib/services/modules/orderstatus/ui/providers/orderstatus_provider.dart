import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/orderstatus_repository.dart';

class OrderStatusProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<OrderStatusMapping> _mappings = [];
  List<OrderStatusInfo> _statuses = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  OrderStatusProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<OrderStatusMapping> get mappings => _mappings;
  List<OrderStatusInfo> get statuses => _statuses;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  OrderStatusUseCases get _useCases => OrderStatusUseCases(OrderStatusApiRepository(_apiClient));

  Future<void> fetchMappings({int? integrationTypeId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final params = GetOrderStatusMappingsParams(page: _page, pageSize: _pageSize, integrationTypeId: integrationTypeId);
      final response = await _useCases.getMappings(params);
      _mappings = response.data;
      _pagination = response.pagination;
    } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }

  Future<void> fetchStatuses({bool? isActive}) async {
    try {
      _statuses = await _useCases.getStatuses(isActive: isActive);
      notifyListeners();
    } catch (e) { _error = e.toString(); notifyListeners(); }
  }

  void setPage(int page) { _page = page; }
}
