import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class StorefrontUseCases {
  final IStorefrontRepository _repository;

  StorefrontUseCases(this._repository);

  Future<PaginatedResponse<StorefrontProduct>> getCatalog(GetCatalogParams? params) {
    return _repository.getCatalog(params);
  }

  Future<StorefrontProduct> getProduct(String id, {int? businessId}) {
    return _repository.getProduct(id, businessId: businessId);
  }

  Future<Map<String, dynamic>> createOrder(CreateStorefrontOrderDTO data, {int? businessId}) {
    return _repository.createOrder(data, businessId: businessId);
  }

  Future<PaginatedResponse<StorefrontOrder>> getOrders(GetOrdersParams? params) {
    return _repository.getOrders(params);
  }

  Future<StorefrontOrder> getOrder(String id, {int? businessId}) {
    return _repository.getOrder(id, businessId: businessId);
  }

  Future<Map<String, dynamic>> register(RegisterDTO data) {
    return _repository.register(data);
  }
}
