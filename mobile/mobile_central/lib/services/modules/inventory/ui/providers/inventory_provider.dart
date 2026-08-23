import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/inventory_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class InventoryProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final InventoryUseCases? _injectedUseCases;

  List<MovementType> _movementTypes = [];
  String? _error;
  String _searchFilter = '';
  bool? _lowStockFilter;
  int? _warehouseId;
  int? _businessId;
  GetMovementsParams? _movementsParams;

  late final PagedListController<InventoryLevel> levels;
  late final PagedListController<StockMovement> movementsList;

  InventoryProvider({required ApiClient apiClient, InventoryUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    levels = PagedListController<InventoryLevel>(fetcher: _fetchLevels);
    movementsList = PagedListController<StockMovement>(fetcher: _fetchMovements);
    levels.addListener(notifyListeners);
    movementsList.addListener(notifyListeners);
  }

  List<InventoryLevel> get inventoryLevels => levels.loadedItems;
  List<StockMovement> get movements => movementsList.loadedItems;
  List<MovementType> get movementTypes => _movementTypes;
  bool get isLoading => levels.isLoading || movementsList.isLoading;
  String? get error => _error ?? levels.error ?? movementsList.error;

  InventoryUseCases get _useCases =>
      _injectedUseCases ?? InventoryUseCases(InventoryApiRepository(_apiClient));

  Future<PaginatedResponse<InventoryLevel>> _fetchLevels(int page, int pageSize) {
    return _useCases.getWarehouseInventory(
      _warehouseId ?? 0,
      GetInventoryParams(
        page: page,
        pageSize: pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        lowStock: _lowStockFilter,
        businessId: _businessId,
      ),
    );
  }

  Future<void> fetchWarehouseInventory(int warehouseId, {int? businessId}) {
    _warehouseId = warehouseId;
    _businessId = businessId;
    _error = null;
    return levels.refresh();
  }

  Future<List<InventoryLevel>> getProductInventory(String productId, {int? businessId}) async {
    try {
      return await _useCases.getProductInventory(productId, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return [];
    }
  }

  Future<StockMovement?> adjustStock(AdjustStockDTO data, {int? businessId}) async {
    try {
      final movement = await _useCases.adjustStock(data, businessId: businessId);
      return movement;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> transferStock(TransferStockDTO data, {int? businessId}) async {
    try {
      await _useCases.transferStock(data, businessId: businessId);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<PaginatedResponse<StockMovement>> _fetchMovements(int page, int pageSize) {
    final base = _movementsParams;
    return _useCases.getMovements(GetMovementsParams(
      page: page,
      pageSize: pageSize,
      warehouseId: base?.warehouseId,
      productId: base?.productId,
      type: base?.type,
      businessId: base?.businessId,
    ));
  }

  Future<void> fetchMovements({GetMovementsParams? params}) {
    _movementsParams = params;
    _error = null;
    return movementsList.refresh();
  }

  Future<void> fetchMovementTypes({GetMovementTypesParams? params}) async {
    try {
      final response = await _useCases.getMovementTypes(params);
      _movementTypes = response.data;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
    }
  }

  void setFilters({
    String? search,
    bool? lowStock,
  }) {
    _searchFilter = search ?? _searchFilter;
    _lowStockFilter = lowStock ?? _lowStockFilter;
  }

  void resetFilters() {
    _searchFilter = '';
    _lowStockFilter = null;
  }

  @override
  void dispose() {
    levels.removeListener(notifyListeners);
    movementsList.removeListener(notifyListeners);
    levels.dispose();
    movementsList.dispose();
    super.dispose();
  }
}
