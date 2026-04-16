import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IWarehouseRepository {
  Future<PaginatedResponse<Warehouse>> getWarehouses(GetWarehousesParams? params);
  Future<WarehouseDetail> getWarehouseById(int id, {int? businessId});
  Future<Warehouse> createWarehouse(CreateWarehouseDTO data, {int? businessId});
  Future<Warehouse> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId});
  Future<void> deleteWarehouse(int id, {int? businessId});
  // Locations
  Future<List<WarehouseLocation>> getLocations(int warehouseId, {int? businessId});
  Future<WarehouseLocation> createLocation(int warehouseId, CreateLocationDTO data, {int? businessId});
  Future<WarehouseLocation> updateLocation(int warehouseId, int locationId, UpdateLocationDTO data, {int? businessId});
  Future<void> deleteLocation(int warehouseId, int locationId, {int? businessId});
}
