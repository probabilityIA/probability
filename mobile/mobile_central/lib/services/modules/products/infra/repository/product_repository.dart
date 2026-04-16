import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class ProductApiRepository implements IProductRepository {
  final ApiClient _client;

  ProductApiRepository(this._client);

  @override
  Future<PaginatedResponse<Product>> getProducts(GetProductsParams? params) async {
    final response = await _client.get('/products', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => Product.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<Product> getProductById(String id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/products/$id', queryParameters: qp);
    return Product.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Product> createProduct(CreateProductDTO data, {int? businessId}) async {
    final response = await _client.post('/products', data: data.toJson());
    return Product.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Product> updateProduct(String id, UpdateProductDTO data, {int? businessId}) async {
    final response = await _client.put('/products/$id', data: data.toJson());
    return Product.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteProduct(String id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    await _client.delete('/products/$id', queryParameters: qp);
  }

  @override
  Future<List<ProductIntegration>> getProductIntegrations(String productId, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/products/$productId/integrations', queryParameters: qp);
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => ProductIntegration.fromJson(e)).toList();
  }

  @override
  Future<void> addProductIntegration(String productId, AddProductIntegrationDTO data, {int? businessId}) async {
    await _client.post('/products/$productId/integrations', data: data.toJson());
  }

  @override
  Future<void> removeProductIntegration(String productId, int integrationId, {int? businessId}) async {
    await _client.delete('/products/$productId/integrations/$integrationId');
  }
}
