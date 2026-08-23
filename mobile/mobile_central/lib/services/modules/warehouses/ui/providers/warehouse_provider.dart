import 'package:flutter/foundation.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/warehouse_repository.dart';

class WarehouseProvider extends ChangeNotifier {
  WarehouseProvider({required ApiClient apiClient}) : _apiClient = apiClient {
    list = PagedListController<Warehouse>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  late final PagedListController<Warehouse> list;

  int? _businessId;
  String? _error;

  List<WarehouseLocation> _locations = [];
  bool _isLoadingLocations = false;

  List<Warehouse> get warehouses => list.loadedItems;
  bool get isLoading => list.isLoading;
  String? get error => _error ?? list.error;
  List<WarehouseLocation> get locations => _locations;
  bool get isLoadingLocations => _isLoadingLocations;

  WarehouseUseCases get _useCases =>
      WarehouseUseCases(WarehouseApiRepository(_apiClient));

  Future<PaginatedResponse<Warehouse>> _fetchPage(int page, int pageSize) {
    return _useCases.getWarehouses(GetWarehousesParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
    ));
  }

  Future<void> fetchWarehouses({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Warehouse? warehouseById(int id) {
    for (final warehouse in list.loadedItems) {
      if (warehouse.id == id) return warehouse;
    }
    return null;
  }

  Future<void> fetchLocations(int warehouseId, {int? businessId}) async {
    _isLoadingLocations = true;
    _locations = [];
    notifyListeners();
    try {
      _locations = await _useCases.getLocations(warehouseId, businessId: businessId);
    } catch (_) {
      _locations = [];
    }
    _isLoadingLocations = false;
    notifyListeners();
  }

  Future<Warehouse?> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    try {
      return await _useCases.createWarehouse(data, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Warehouse?> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    try {
      return await _useCases.updateWarehouse(id, data, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> deleteWarehouse(int id, {int? businessId}) async {
    try {
      await _useCases.deleteWarehouse(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
