import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class BusinessApiRepository implements IBusinessRepository {
  final ApiClient _client;

  BusinessApiRepository(this._client);

  @override
  Future<PaginatedResponse<Business>> getBusinesses(
      GetBusinessesParams? params) async {
    final response = await _client.get(
      '/businesses',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final businesses = (data['data'] as List<dynamic>?)
            ?.map((b) => Business.fromJson(b))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: businesses, pagination: pagination);
  }

  @override
  Future<Business> getBusinessById(int id) async {
    final response = await _client.get('/businesses/$id');
    return Business.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Business> createBusiness(CreateBusinessDTO data) async {
    final response = await _client.post('/businesses', data: data.toJson());
    return Business.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Business> updateBusiness(int id, UpdateBusinessDTO data) async {
    final response =
        await _client.put('/businesses/$id', data: data.toJson());
    return Business.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteBusiness(int id) async {
    await _client.delete('/businesses/$id');
  }

  @override
  Future<void> activateBusiness(int id) async {
    await _client.post('/businesses/$id/activate');
  }

  @override
  Future<void> deactivateBusiness(int id) async {
    await _client.post('/businesses/$id/deactivate');
  }

  @override
  Future<List<BusinessSimple>> getBusinessesSimple() async {
    final response = await _client.get('/businesses/simple');
    final data = response.data;
    final list = data is List ? data : (data['data'] as List<dynamic>?) ?? [];
    return list.map((b) => BusinessSimple.fromJson(b)).toList();
  }

  @override
  Future<List<ConfiguredResource>> getConfiguredResources(
      int businessId) async {
    final response =
        await _client.get('/businesses/$businessId/configured-resources');
    final data = response.data;
    final list = data is List ? data : (data['data'] as List<dynamic>?) ?? [];
    return list.map((r) => ConfiguredResource.fromJson(r)).toList();
  }

  @override
  Future<void> activateConfiguredResource(int resourceId) async {
    await _client.post(
        '/businesses/configured-resources/$resourceId/activate');
  }

  @override
  Future<void> deactivateConfiguredResource(int resourceId) async {
    await _client.post(
        '/businesses/configured-resources/$resourceId/deactivate');
  }

  @override
  Future<List<BusinessType>> getBusinessTypes() async {
    final response = await _client.get('/business-types');
    final data = response.data;
    final list = data is List ? data : (data['data'] as List<dynamic>?) ?? [];
    return list.map((t) => BusinessType.fromJson(t)).toList();
  }

  @override
  Future<BusinessType> createBusinessType(Map<String, dynamic> data) async {
    final response = await _client.post('/business-types', data: data);
    return BusinessType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<BusinessType> updateBusinessType(
      int id, Map<String, dynamic> data) async {
    final response = await _client.put('/business-types/$id', data: data);
    return BusinessType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteBusinessType(int id) async {
    await _client.delete('/business-types/$id');
  }
}
