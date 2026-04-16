import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class VehicleApiRepository implements IVehicleRepository {
  final ApiClient _client;

  VehicleApiRepository(this._client);

  Map<String, dynamic>? _withBusinessId(int? businessId) {
    if (businessId == null) return null;
    return {'business_id': businessId};
  }

  @override
  Future<PaginatedResponse<VehicleInfo>> getVehicles(GetVehiclesParams? params) async {
    final response = await _client.get(
      '/vehicles',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => VehicleInfo.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<VehicleInfo> getVehicleById(int id, {int? businessId}) async {
    final response = await _client.get(
      '/vehicles/$id',
      queryParameters: _withBusinessId(businessId),
    );
    return VehicleInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<VehicleInfo> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    final response = await _client.post(
      '/vehicles',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return VehicleInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<VehicleInfo> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    final response = await _client.put(
      '/vehicles/$id',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return VehicleInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> deleteVehicle(int id, {int? businessId}) async {
    final response = await _client.delete(
      '/vehicles/$id',
      queryParameters: _withBusinessId(businessId),
    );
    return response.data is Map<String, dynamic>
        ? response.data
        : {'message': 'Vehicle deleted'};
  }
}
