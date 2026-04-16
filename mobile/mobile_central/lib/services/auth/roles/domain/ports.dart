import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IRoleRepository {
  Future<PaginatedResponse<Role>> getRoles(GetRolesParams? params);
  Future<Role> getRoleById(int id);
  Future<List<Role>> getRolesByScope(int scopeId);
  Future<List<Role>> getRolesByLevel(int level);
  Future<List<Role>> getSystemRoles();
  Future<Role> createRole(CreateRoleDTO data);
  Future<Role> updateRole(int id, UpdateRoleDTO data);
  Future<void> deleteRole(int id);
  Future<void> assignPermissions(int roleId, AssignPermissionsDTO data);
  Future<RolePermissionsResponse> getRolePermissions(int roleId);
  Future<void> removePermissionFromRole(int roleId, int permissionId);
}
