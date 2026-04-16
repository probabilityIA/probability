import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IProductRepository {
  Future<PaginatedResponse<Product>> getProducts(GetProductsParams? params);
  Future<Product> getProductById(String id, {int? businessId});
  Future<Product> createProduct(CreateProductDTO data, {int? businessId});
  Future<Product> updateProduct(String id, UpdateProductDTO data, {int? businessId});
  Future<void> deleteProduct(String id, {int? businessId});
  Future<List<ProductIntegration>> getProductIntegrations(String productId, {int? businessId});
  Future<void> addProductIntegration(String productId, AddProductIntegrationDTO data, {int? businessId});
  Future<void> removeProductIntegration(String productId, int integrationId, {int? businessId});
}
