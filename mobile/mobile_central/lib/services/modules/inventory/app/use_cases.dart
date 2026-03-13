import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class InventoryUseCases {
  final IInventoryRepository _repository;

  InventoryUseCases(this._repository);

  Future<List<InventoryLevel>> getProductInventory(String productId, {int? businessId}) {
    return _repository.getProductInventory(productId, businessId: businessId);
  }

  Future<PaginatedResponse<InventoryLevel>> getWarehouseInventory(int warehouseId, GetInventoryParams? params) {
    return _repository.getWarehouseInventory(warehouseId, params);
  }

  Future<StockMovement> adjustStock(AdjustStockDTO data, {int? businessId}) {
    return _repository.adjustStock(data, businessId: businessId);
  }

  Future<Map<String, dynamic>> transferStock(TransferStockDTO data, {int? businessId}) {
    return _repository.transferStock(data, businessId: businessId);
  }

  Future<PaginatedResponse<StockMovement>> getMovements(GetMovementsParams? params) {
    return _repository.getMovements(params);
  }

  Future<PaginatedResponse<MovementType>> getMovementTypes(GetMovementTypesParams? params) {
    return _repository.getMovementTypes(params);
  }
}
