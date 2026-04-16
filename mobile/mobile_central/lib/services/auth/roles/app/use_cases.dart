import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class RoleUseCases {
  final IRoleRepository _repository;

  RoleUseCases(this._repository);

  Future<PaginatedResponse<Role>> getRoles(GetRolesParams? params) {
    return _repository.getRoles(params);
  }

  Future<Role> getRoleById(int id) {
    return _repository.getRoleById(id);
  }

  Future<List<Role>> getRolesByScope(int scopeId) {
    return _repository.getRolesByScope(scopeId);
  }

  Future<List<Role>> getRolesByLevel(int level) {
    return _repository.getRolesByLevel(level);
  }

  Future<List<Role>> getSystemRoles() {
    return _repository.getSystemRoles();
  }

  Future<Role> createRole(CreateRoleDTO data) {
    return _repository.createRole(data);
  }

  Future<Role> updateRole(int id, UpdateRoleDTO data) {
    return _repository.updateRole(id, data);
  }

  Future<void> deleteRole(int id) {
    return _repository.deleteRole(id);
  }

  Future<void> assignPermissions(int roleId, AssignPermissionsDTO data) {
    return _repository.assignPermissions(roleId, data);
  }

  Future<RolePermissionsResponse> getRolePermissions(int roleId) {
    return _repository.getRolePermissions(roleId);
  }

  Future<void> removePermissionFromRole(int roleId, int permissionId) {
    return _repository.removePermissionFromRole(roleId, permissionId);
  }
}
