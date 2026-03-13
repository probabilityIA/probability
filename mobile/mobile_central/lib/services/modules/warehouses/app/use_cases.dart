import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class WarehouseUseCases {
  final IWarehouseRepository _repository;

  WarehouseUseCases(this._repository);

  Future<PaginatedResponse<Warehouse>> getWarehouses(GetWarehousesParams? params) {
    return _repository.getWarehouses(params);
  }

  Future<WarehouseDetail> getWarehouseById(int id, {int? businessId}) {
    return _repository.getWarehouseById(id, businessId: businessId);
  }

  Future<Warehouse> createWarehouse(CreateWarehouseDTO data, {int? businessId}) {
    return _repository.createWarehouse(data, businessId: businessId);
  }

  Future<Warehouse> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) {
    return _repository.updateWarehouse(id, data, businessId: businessId);
  }

  Future<void> deleteWarehouse(int id, {int? businessId}) {
    return _repository.deleteWarehouse(id, businessId: businessId);
  }

  Future<List<WarehouseLocation>> getLocations(int warehouseId, {int? businessId}) {
    return _repository.getLocations(warehouseId, businessId: businessId);
  }

  Future<WarehouseLocation> createLocation(int warehouseId, CreateLocationDTO data, {int? businessId}) {
    return _repository.createLocation(warehouseId, data, businessId: businessId);
  }

  Future<WarehouseLocation> updateLocation(int warehouseId, int locationId, UpdateLocationDTO data, {int? businessId}) {
    return _repository.updateLocation(warehouseId, locationId, data, businessId: businessId);
  }

  Future<void> deleteLocation(int warehouseId, int locationId, {int? businessId}) {
    return _repository.deleteLocation(warehouseId, locationId, businessId: businessId);
  }
}
