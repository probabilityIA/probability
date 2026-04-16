import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class WarehouseApiRepository implements IWarehouseRepository {
  final ApiClient _client;
  WarehouseApiRepository(this._client);

  Map<String, dynamic>? _biz(int? businessId) =>
      businessId != null ? {'business_id': businessId} : null;

  @override
  Future<PaginatedResponse<Warehouse>> getWarehouses(GetWarehousesParams? params) async {
    final response = await _client.get('/warehouses', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => Warehouse.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<WarehouseDetail> getWarehouseById(int id, {int? businessId}) async {
    final response = await _client.get('/warehouses/$id', queryParameters: _biz(businessId));
    return WarehouseDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Warehouse> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    final body = data.toJson();
    if (businessId != null) body['business_id'] = businessId;
    final response = await _client.post('/warehouses', data: body);
    return Warehouse.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Warehouse> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    final response = await _client.put('/warehouses/$id', data: data.toJson());
    return Warehouse.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteWarehouse(int id, {int? businessId}) async {
    await _client.delete('/warehouses/$id', queryParameters: _biz(businessId));
  }

  @override
  Future<List<WarehouseLocation>> getLocations(int warehouseId, {int? businessId}) async {
    final response = await _client.get('/warehouses/$warehouseId/locations', queryParameters: _biz(businessId));
    return (response.data['data'] as List<dynamic>?)?.map((e) => WarehouseLocation.fromJson(e)).toList() ?? [];
  }

  @override
  Future<WarehouseLocation> createLocation(int warehouseId, CreateLocationDTO data, {int? businessId}) async {
    final response = await _client.post('/warehouses/$warehouseId/locations', data: data.toJson());
    return WarehouseLocation.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<WarehouseLocation> updateLocation(int warehouseId, int locationId, UpdateLocationDTO data, {int? businessId}) async {
    final response = await _client.put('/warehouses/$warehouseId/locations/$locationId', data: data.toJson());
    return WarehouseLocation.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteLocation(int warehouseId, int locationId, {int? businessId}) async {
    await _client.delete('/warehouses/$warehouseId/locations/$locationId', queryParameters: _biz(businessId));
  }
}
