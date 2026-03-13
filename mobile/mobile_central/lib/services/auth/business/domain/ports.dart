import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IBusinessRepository {
  Future<PaginatedResponse<Business>> getBusinesses(
      GetBusinessesParams? params);
  Future<Business> getBusinessById(int id);
  Future<Business> createBusiness(CreateBusinessDTO data);
  Future<Business> updateBusiness(int id, UpdateBusinessDTO data);
  Future<void> deleteBusiness(int id);
  Future<void> activateBusiness(int id);
  Future<void> deactivateBusiness(int id);

  Future<List<BusinessSimple>> getBusinessesSimple();

  Future<List<ConfiguredResource>> getConfiguredResources(int businessId);
  Future<void> activateConfiguredResource(int resourceId);
  Future<void> deactivateConfiguredResource(int resourceId);

  Future<List<BusinessType>> getBusinessTypes();
  Future<BusinessType> createBusinessType(Map<String, dynamic> data);
  Future<BusinessType> updateBusinessType(int id, Map<String, dynamic> data);
  Future<void> deleteBusinessType(int id);
}
