import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IVehicleRepository {
  Future<PaginatedResponse<VehicleInfo>> getVehicles(GetVehiclesParams? params);
  Future<VehicleInfo> getVehicleById(int id, {int? businessId});
  Future<VehicleInfo> createVehicle(CreateVehicleDTO data, {int? businessId});
  Future<VehicleInfo> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId});
  Future<Map<String, dynamic>> deleteVehicle(int id, {int? businessId});
}
