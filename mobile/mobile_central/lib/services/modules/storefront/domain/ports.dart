import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IStorefrontRepository {
  Future<PaginatedResponse<StorefrontProduct>> getCatalog(GetCatalogParams? params);
  Future<StorefrontProduct> getProduct(String id, {int? businessId});
  Future<Map<String, dynamic>> createOrder(CreateStorefrontOrderDTO data, {int? businessId});
  Future<PaginatedResponse<StorefrontOrder>> getOrders(GetOrdersParams? params);
  Future<StorefrontOrder> getOrder(String id, {int? businessId});
  Future<Map<String, dynamic>> register(RegisterDTO data);
}
