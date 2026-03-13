import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class VehicleUseCases {
  final IVehicleRepository _repository;

  VehicleUseCases(this._repository);

  Future<PaginatedResponse<VehicleInfo>> getVehicles(GetVehiclesParams? params) {
    return _repository.getVehicles(params);
  }

  Future<VehicleInfo> getVehicleById(int id, {int? businessId}) {
    return _repository.getVehicleById(id, businessId: businessId);
  }

  Future<VehicleInfo> createVehicle(CreateVehicleDTO data, {int? businessId}) {
    return _repository.createVehicle(data, businessId: businessId);
  }

  Future<VehicleInfo> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) {
    return _repository.updateVehicle(id, data, businessId: businessId);
  }

  Future<Map<String, dynamic>> deleteVehicle(int id, {int? businessId}) {
    return _repository.deleteVehicle(id, businessId: businessId);
  }
}
