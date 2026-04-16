import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class CustomerApiRepository implements ICustomerRepository {
  final ApiClient _client;

  CustomerApiRepository(this._client);

  @override
  Future<PaginatedResponse<CustomerInfo>> getCustomers(GetCustomersParams? params) async {
    final response = await _client.get('/customers', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => CustomerInfo.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<CustomerDetail> getCustomerById(int id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/customers/$id', queryParameters: qp);
    return CustomerDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<CustomerInfo> createCustomer(CreateCustomerDTO data, {int? businessId}) async {
    final body = data.toJson();
    if (businessId != null) body['business_id'] = businessId;
    final response = await _client.post('/customers', data: body);
    return CustomerInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<CustomerInfo> updateCustomer(int id, UpdateCustomerDTO data, {int? businessId}) async {
    final response = await _client.put('/customers/$id', data: data.toJson());
    return CustomerInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteCustomer(int id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    await _client.delete('/customers/$id', queryParameters: qp);
  }
}
