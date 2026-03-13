import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/inventory_repository.dart';

class InventoryProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<InventoryLevel> _inventoryLevels = [];
  List<StockMovement> _movements = [];
  List<MovementType> _movementTypes = [];
  Pagination? _pagination;
  Pagination? _movementsPagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  bool? _lowStockFilter;

  InventoryProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<InventoryLevel> get inventoryLevels => _inventoryLevels;
  List<StockMovement> get movements => _movements;
  List<MovementType> get movementTypes => _movementTypes;
  Pagination? get pagination => _pagination;
  Pagination? get movementsPagination => _movementsPagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;

  InventoryUseCases get _useCases =>
      InventoryUseCases(InventoryApiRepository(_apiClient));

  Future<void> fetchWarehouseInventory(int warehouseId, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetInventoryParams(
        page: _page,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        lowStock: _lowStockFilter,
        businessId: businessId,
      );
      final response = await _useCases.getWarehouseInventory(warehouseId, params);
      _inventoryLevels = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<List<InventoryLevel>> getProductInventory(String productId, {int? businessId}) async {
    try {
      return await _useCases.getProductInventory(productId, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return [];
    }
  }

  Future<StockMovement?> adjustStock(AdjustStockDTO data, {int? businessId}) async {
    try {
      final movement = await _useCases.adjustStock(data, businessId: businessId);
      return movement;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> transferStock(TransferStockDTO data, {int? businessId}) async {
    try {
      await _useCases.transferStock(data, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<void> fetchMovements({GetMovementsParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.getMovements(params);
      _movements = response.data;
      _movementsPagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchMovementTypes({GetMovementTypesParams? params}) async {
    try {
      final response = await _useCases.getMovementTypes(params);
      _movementTypes = response.data;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
    }
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({
    String? search,
    bool? lowStock,
  }) {
    _searchFilter = search ?? _searchFilter;
    _lowStockFilter = lowStock ?? _lowStockFilter;
    _page = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _lowStockFilter = null;
    _page = 1;
  }
}
