import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class InventoryApiRepository implements IInventoryRepository {
  final ApiClient _client;

  InventoryApiRepository(this._client);

  Map<String, dynamic>? _withBusinessId(int? businessId) {
    if (businessId == null) return null;
    return {'business_id': businessId};
  }

  @override
  Future<List<InventoryLevel>> getProductInventory(String productId, {int? businessId}) async {
    final response = await _client.get(
      '/inventory/product/$productId',
      queryParameters: _withBusinessId(businessId),
    );
    final data = response.data['data'] ?? response.data;
    if (data is List) {
      return data.map((e) => InventoryLevel.fromJson(e)).toList();
    }
    return [];
  }

  @override
  Future<PaginatedResponse<InventoryLevel>> getWarehouseInventory(int warehouseId, GetInventoryParams? params) async {
    final response = await _client.get(
      '/inventory/warehouse/$warehouseId',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => InventoryLevel.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<StockMovement> adjustStock(AdjustStockDTO data, {int? businessId}) async {
    final response = await _client.post(
      '/inventory/adjust',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return StockMovement.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> transferStock(TransferStockDTO data, {int? businessId}) async {
    final response = await _client.post(
      '/inventory/transfer',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return response.data is Map<String, dynamic>
        ? response.data
        : {'message': 'Transfer completed'};
  }

  @override
  Future<PaginatedResponse<StockMovement>> getMovements(GetMovementsParams? params) async {
    final response = await _client.get(
      '/inventory/movements',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => StockMovement.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<PaginatedResponse<MovementType>> getMovementTypes(GetMovementTypesParams? params) async {
    final response = await _client.get(
      '/inventory/movement-types',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => MovementType.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }
}
