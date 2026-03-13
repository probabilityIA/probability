import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class StorefrontApiRepository implements IStorefrontRepository {
  final ApiClient _client;

  StorefrontApiRepository(this._client);

  @override
  Future<PaginatedResponse<StorefrontProduct>> getCatalog(GetCatalogParams? params) async {
    final response = await _client.get(
      '/storefront/catalog',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => StorefrontProduct.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<StorefrontProduct> getProduct(String id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/storefront/catalog/$id', queryParameters: qp);
    return StorefrontProduct.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> createOrder(CreateStorefrontOrderDTO data, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.post(
      '/storefront/orders',
      data: data.toJson(),
      queryParameters: qp,
    );
    return Map<String, dynamic>.from(response.data);
  }

  @override
  Future<PaginatedResponse<StorefrontOrder>> getOrders(GetOrdersParams? params) async {
    final response = await _client.get(
      '/storefront/orders',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => StorefrontOrder.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<StorefrontOrder> getOrder(String id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/storefront/orders/$id', queryParameters: qp);
    return StorefrontOrder.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> register(RegisterDTO data) async {
    final response = await _client.post('/storefront/register', data: data.toJson());
    return Map<String, dynamic>.from(response.data);
  }
}
