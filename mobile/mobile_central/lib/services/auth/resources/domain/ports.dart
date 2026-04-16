import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IResourceRepository {
  Future<PaginatedResponse<Resource>> getResources(
      GetResourcesParams? params);
  Future<Resource> getResourceById(int id);
  Future<Resource> createResource(CreateResourceDTO data);
  Future<Resource> updateResource(int id, UpdateResourceDTO data);
  Future<void> deleteResource(int id);
}
