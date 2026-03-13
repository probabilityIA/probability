import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class ResourceApiRepository implements IResourceRepository {
  final ApiClient _client;

  ResourceApiRepository(this._client);

  @override
  Future<PaginatedResponse<Resource>> getResources(
      GetResourcesParams? params) async {
    final response = await _client.get(
      '/resources',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final resources = (data['data'] as List<dynamic>?)
            ?.map((r) => Resource.fromJson(r))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: resources, pagination: pagination);
  }

  @override
  Future<Resource> getResourceById(int id) async {
    final response = await _client.get('/resources/$id');
    return Resource.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Resource> createResource(CreateResourceDTO data) async {
    final response = await _client.post('/resources', data: data.toJson());
    return Resource.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Resource> updateResource(int id, UpdateResourceDTO data) async {
    final response =
        await _client.put('/resources/$id', data: data.toJson());
    return Resource.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteResource(int id) async {
    await _client.delete('/resources/$id');
  }
}
