import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class DriverApiRepository implements IDriverRepository {
  final ApiClient _client;

  DriverApiRepository(this._client);

  Map<String, dynamic>? _withBusinessId(int? businessId) {
    if (businessId == null) return null;
    return {'business_id': businessId};
  }

  @override
  Future<PaginatedResponse<DriverInfo>> getDrivers(GetDriversParams? params) async {
    final response = await _client.get(
      '/drivers',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => DriverInfo.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<DriverInfo> getDriverById(int id, {int? businessId}) async {
    final response = await _client.get(
      '/drivers/$id',
      queryParameters: _withBusinessId(businessId),
    );
    return DriverInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<DriverInfo> createDriver(CreateDriverDTO data, {int? businessId}) async {
    final response = await _client.post(
      '/drivers',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return DriverInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<DriverInfo> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) async {
    final response = await _client.put(
      '/drivers/$id',
      data: data.toJson(),
      queryParameters: _withBusinessId(businessId),
    );
    return DriverInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> deleteDriver(int id, {int? businessId}) async {
    final response = await _client.delete(
      '/drivers/$id',
      queryParameters: _withBusinessId(businessId),
    );
    return response.data is Map<String, dynamic>
        ? response.data
        : {'message': 'Driver deleted'};
  }
}
