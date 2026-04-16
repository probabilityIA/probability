import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class ProductUseCases {
  final IProductRepository _repository;

  ProductUseCases(this._repository);

  Future<PaginatedResponse<Product>> getProducts(GetProductsParams? params) {
    return _repository.getProducts(params);
  }

  Future<Product> getProductById(String id, {int? businessId}) {
    return _repository.getProductById(id, businessId: businessId);
  }

  Future<Product> createProduct(CreateProductDTO data, {int? businessId}) {
    return _repository.createProduct(data, businessId: businessId);
  }

  Future<Product> updateProduct(String id, UpdateProductDTO data, {int? businessId}) {
    return _repository.updateProduct(id, data, businessId: businessId);
  }

  Future<void> deleteProduct(String id, {int? businessId}) {
    return _repository.deleteProduct(id, businessId: businessId);
  }

  Future<List<ProductIntegration>> getProductIntegrations(String productId, {int? businessId}) {
    return _repository.getProductIntegrations(productId, businessId: businessId);
  }

  Future<void> addProductIntegration(String productId, AddProductIntegrationDTO data, {int? businessId}) {
    return _repository.addProductIntegration(productId, data, businessId: businessId);
  }

  Future<void> removeProductIntegration(String productId, int integrationId, {int? businessId}) {
    return _repository.removeProductIntegration(productId, integrationId, businessId: businessId);
  }
}
