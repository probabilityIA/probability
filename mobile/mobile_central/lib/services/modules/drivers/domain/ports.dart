import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IDriverRepository {
  Future<PaginatedResponse<DriverInfo>> getDrivers(GetDriversParams? params);
  Future<DriverInfo> getDriverById(int id, {int? businessId});
  Future<DriverInfo> createDriver(CreateDriverDTO data, {int? businessId});
  Future<DriverInfo> updateDriver(int id, UpdateDriverDTO data, {int? businessId});
  Future<Map<String, dynamic>> deleteDriver(int id, {int? businessId});
}
