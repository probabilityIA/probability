import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class RoleApiRepository implements IRoleRepository {
  final ApiClient _client;

  RoleApiRepository(this._client);

  @override
  Future<PaginatedResponse<Role>> getRoles(GetRolesParams? params) async {
    final response = await _client.get(
      '/roles',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final roles = (data['data'] as List<dynamic>?)
            ?.map((r) => Role.fromJson(r))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: roles, pagination: pagination);
  }

  @override
  Future<Role> getRoleById(int id) async {
    final response = await _client.get('/roles/$id');
    return Role.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<List<Role>> getRolesByScope(int scopeId) async {
    final response = await _client.get('/roles/scope/$scopeId');
    return (response.data['data'] as List<dynamic>?)
            ?.map((r) => Role.fromJson(r))
            .toList() ??
        [];
  }

  @override
  Future<List<Role>> getRolesByLevel(int level) async {
    final response = await _client.get('/roles/level/$level');
    return (response.data['data'] as List<dynamic>?)
            ?.map((r) => Role.fromJson(r))
            .toList() ??
        [];
  }

  @override
  Future<List<Role>> getSystemRoles() async {
    final response = await _client.get('/roles/system');
    return (response.data['data'] as List<dynamic>?)
            ?.map((r) => Role.fromJson(r))
            .toList() ??
        [];
  }

  @override
  Future<Role> createRole(CreateRoleDTO data) async {
    final response = await _client.post('/roles', data: data.toJson());
    return Role.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Role> updateRole(int id, UpdateRoleDTO data) async {
    final response = await _client.put('/roles/$id', data: data.toJson());
    return Role.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteRole(int id) async {
    await _client.delete('/roles/$id');
  }

  @override
  Future<void> assignPermissions(int roleId, AssignPermissionsDTO data) async {
    await _client.post('/roles/$roleId/permissions', data: data.toJson());
  }

  @override
  Future<RolePermissionsResponse> getRolePermissions(int roleId) async {
    final response = await _client.get('/roles/$roleId/permissions');
    return RolePermissionsResponse.fromJson(
        response.data['data'] ?? response.data);
  }

  @override
  Future<void> removePermissionFromRole(
      int roleId, int permissionId) async {
    await _client.delete('/roles/$roleId/permissions/$permissionId');
  }
}
