import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class ICustomerRepository {
  Future<PaginatedResponse<CustomerInfo>> getCustomers(GetCustomersParams? params);
  Future<CustomerDetail> getCustomerById(int id, {int? businessId});
  Future<CustomerInfo> createCustomer(CreateCustomerDTO data, {int? businessId});
  Future<CustomerInfo> updateCustomer(int id, UpdateCustomerDTO data, {int? businessId});
  Future<void> deleteCustomer(int id, {int? businessId});
}
