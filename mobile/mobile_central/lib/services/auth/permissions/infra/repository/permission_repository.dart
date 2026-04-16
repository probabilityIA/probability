import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class PermissionApiRepository implements IPermissionRepository {
  final ApiClient _client;

  PermissionApiRepository(this._client);

  @override
  Future<PaginatedResponse<Permission>> getPermissions(
      GetPermissionsParams? params) async {
    final response = await _client.get(
      '/permissions',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final permissions = (data['data'] as List<dynamic>?)
            ?.map((p) => Permission.fromJson(p))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: permissions, pagination: pagination);
  }

  @override
  Future<Permission> getPermissionById(int id) async {
    final response = await _client.get('/permissions/$id');
    return Permission.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<List<Permission>> getPermissionsByScope(int scopeId) async {
    final response = await _client.get('/permissions/scope/$scopeId');
    return (response.data['data'] as List<dynamic>?)
            ?.map((p) => Permission.fromJson(p))
            .toList() ??
        [];
  }

  @override
  Future<List<Permission>> getPermissionsByResource(String resource) async {
    final response = await _client.get('/permissions/resource/$resource');
    return (response.data['data'] as List<dynamic>?)
            ?.map((p) => Permission.fromJson(p))
            .toList() ??
        [];
  }

  @override
  Future<Permission> createPermission(CreatePermissionDTO data) async {
    final response = await _client.post('/permissions', data: data.toJson());
    return Permission.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Permission> updatePermission(
      int id, UpdatePermissionDTO data) async {
    final response =
        await _client.put('/permissions/$id', data: data.toJson());
    return Permission.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deletePermission(int id) async {
    await _client.delete('/permissions/$id');
  }

  @override
  Future<void> createPermissionsBulk(BulkCreatePermissionsDTO data) async {
    await _client.post('/permissions/bulk', data: data.toJson());
  }
}
