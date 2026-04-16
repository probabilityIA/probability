import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class BusinessUseCases {
  final IBusinessRepository _repository;

  BusinessUseCases(this._repository);

  Future<PaginatedResponse<Business>> getBusinesses(
      GetBusinessesParams? params) {
    return _repository.getBusinesses(params);
  }

  Future<Business> getBusinessById(int id) {
    return _repository.getBusinessById(id);
  }

  Future<Business> createBusiness(CreateBusinessDTO data) {
    return _repository.createBusiness(data);
  }

  Future<Business> updateBusiness(int id, UpdateBusinessDTO data) {
    return _repository.updateBusiness(id, data);
  }

  Future<void> deleteBusiness(int id) {
    return _repository.deleteBusiness(id);
  }

  Future<void> activateBusiness(int id) {
    return _repository.activateBusiness(id);
  }

  Future<void> deactivateBusiness(int id) {
    return _repository.deactivateBusiness(id);
  }

  Future<List<BusinessSimple>> getBusinessesSimple() {
    return _repository.getBusinessesSimple();
  }

  Future<List<ConfiguredResource>> getConfiguredResources(int businessId) {
    return _repository.getConfiguredResources(businessId);
  }

  Future<void> activateConfiguredResource(int resourceId) {
    return _repository.activateConfiguredResource(resourceId);
  }

  Future<void> deactivateConfiguredResource(int resourceId) {
    return _repository.deactivateConfiguredResource(resourceId);
  }

  Future<List<BusinessType>> getBusinessTypes() {
    return _repository.getBusinessTypes();
  }

  Future<BusinessType> createBusinessType(Map<String, dynamic> data) {
    return _repository.createBusinessType(data);
  }

  Future<BusinessType> updateBusinessType(int id, Map<String, dynamic> data) {
    return _repository.updateBusinessType(id, data);
  }

  Future<void> deleteBusinessType(int id) {
    return _repository.deleteBusinessType(id);
  }
}
