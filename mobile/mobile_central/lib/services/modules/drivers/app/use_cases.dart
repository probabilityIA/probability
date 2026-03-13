import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class DriverUseCases {
  final IDriverRepository _repository;

  DriverUseCases(this._repository);

  Future<PaginatedResponse<DriverInfo>> getDrivers(GetDriversParams? params) {
    return _repository.getDrivers(params);
  }

  Future<DriverInfo> getDriverById(int id, {int? businessId}) {
    return _repository.getDriverById(id, businessId: businessId);
  }

  Future<DriverInfo> createDriver(CreateDriverDTO data, {int? businessId}) {
    return _repository.createDriver(data, businessId: businessId);
  }

  Future<DriverInfo> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) {
    return _repository.updateDriver(id, data, businessId: businessId);
  }

  Future<Map<String, dynamic>> deleteDriver(int id, {int? businessId}) {
    return _repository.deleteDriver(id, businessId: businessId);
  }
}
