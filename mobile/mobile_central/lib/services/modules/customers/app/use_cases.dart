import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class CustomerUseCases {
  final ICustomerRepository _repository;

  CustomerUseCases(this._repository);

  Future<PaginatedResponse<CustomerInfo>> getCustomers(GetCustomersParams? params) {
    return _repository.getCustomers(params);
  }

  Future<CustomerDetail> getCustomerById(int id, {int? businessId}) {
    return _repository.getCustomerById(id, businessId: businessId);
  }

  Future<CustomerInfo> createCustomer(CreateCustomerDTO data, {int? businessId}) {
    return _repository.createCustomer(data, businessId: businessId);
  }

  Future<CustomerInfo> updateCustomer(int id, UpdateCustomerDTO data, {int? businessId}) {
    return _repository.updateCustomer(id, data, businessId: businessId);
  }

  Future<void> deleteCustomer(int id, {int? businessId}) {
    return _repository.deleteCustomer(id, businessId: businessId);
  }
}
