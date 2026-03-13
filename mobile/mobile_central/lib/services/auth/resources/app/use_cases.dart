import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class ResourceUseCases {
  final IResourceRepository _repository;

  ResourceUseCases(this._repository);

  Future<PaginatedResponse<Resource>> getResources(
      GetResourcesParams? params) {
    return _repository.getResources(params);
  }

  Future<Resource> getResourceById(int id) {
    return _repository.getResourceById(id);
  }

  Future<Resource> createResource(CreateResourceDTO data) {
    return _repository.createResource(data);
  }

  Future<Resource> updateResource(int id, UpdateResourceDTO data) {
    return _repository.updateResource(id, data);
  }

  Future<void> deleteResource(int id) {
    return _repository.deleteResource(id);
  }
}
