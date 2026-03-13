import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IPermissionRepository {
  Future<PaginatedResponse<Permission>> getPermissions(
      GetPermissionsParams? params);
  Future<Permission> getPermissionById(int id);
  Future<List<Permission>> getPermissionsByScope(int scopeId);
  Future<List<Permission>> getPermissionsByResource(String resource);
  Future<Permission> createPermission(CreatePermissionDTO data);
  Future<Permission> updatePermission(int id, UpdatePermissionDTO data);
  Future<void> deletePermission(int id);
  Future<void> createPermissionsBulk(BulkCreatePermissionsDTO data);
}
