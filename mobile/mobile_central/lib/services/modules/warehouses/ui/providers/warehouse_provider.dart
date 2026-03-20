import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/warehouse_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class WarehouseProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<Warehouse> _warehouses = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  WarehouseProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Warehouse> get warehouses => _warehouses;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  WarehouseUseCases get _useCases => WarehouseUseCases(WarehouseApiRepository(_apiClient));

  Future<void> fetchWarehouses({int? businessId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final params = GetWarehousesParams(page: _page, pageSize: _pageSize, businessId: businessId);
      final response = await _useCases.getWarehouses(params);
      _warehouses = response.data; _pagination = response.pagination;
    } catch (e) { _error = parseError(e); }
    _isLoading = false; notifyListeners();
  }

  Future<Warehouse?> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    try { return await _useCases.createWarehouse(data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<Warehouse?> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    try { return await _useCases.updateWarehouse(id, data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<bool> deleteWarehouse(int id, {int? businessId}) async {
    try { await _useCases.deleteWarehouse(id, businessId: businessId); return true; } catch (e) { _error = parseError(e); notifyListeners(); return false; }
  }

  void setPage(int page) { _page = page; }
}
