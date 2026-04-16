import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IInventoryRepository {
  Future<List<InventoryLevel>> getProductInventory(String productId, {int? businessId});
  Future<PaginatedResponse<InventoryLevel>> getWarehouseInventory(int warehouseId, GetInventoryParams? params);
  Future<StockMovement> adjustStock(AdjustStockDTO data, {int? businessId});
  Future<Map<String, dynamic>> transferStock(TransferStockDTO data, {int? businessId});
  Future<PaginatedResponse<StockMovement>> getMovements(GetMovementsParams? params);
  Future<PaginatedResponse<MovementType>> getMovementTypes(GetMovementTypesParams? params);
}
