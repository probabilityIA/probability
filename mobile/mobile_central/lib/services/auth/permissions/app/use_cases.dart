import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class PermissionUseCases {
  final IPermissionRepository _repository;

  PermissionUseCases(this._repository);

  Future<PaginatedResponse<Permission>> getPermissions(
      GetPermissionsParams? params) {
    return _repository.getPermissions(params);
  }

  Future<Permission> getPermissionById(int id) {
    return _repository.getPermissionById(id);
  }

  Future<List<Permission>> getPermissionsByScope(int scopeId) {
    return _repository.getPermissionsByScope(scopeId);
  }

  Future<List<Permission>> getPermissionsByResource(String resource) {
    return _repository.getPermissionsByResource(resource);
  }

  Future<Permission> createPermission(CreatePermissionDTO data) {
    return _repository.createPermission(data);
  }

  Future<Permission> updatePermission(int id, UpdatePermissionDTO data) {
    return _repository.updatePermission(id, data);
  }

  Future<void> deletePermission(int id) {
    return _repository.deletePermission(id);
  }

  Future<void> createPermissionsBulk(BulkCreatePermissionsDTO data) {
    return _repository.createPermissionsBulk(data);
  }
}
