import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class PublicSiteApiRepository implements IPublicSiteRepository {
  final ApiClient _client;

  PublicSiteApiRepository(this._client);

  @override
  Future<PublicBusiness> getBusinessPage(String slug) async {
    final response = await _client.get('/public/tienda/$slug');
    return PublicBusiness.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<PaginatedResponse<PublicProduct>> getCatalog(
    String slug,
    GetPublicCatalogParams? params,
  ) async {
    final response = await _client.get(
      '/public/tienda/$slug/catalog',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => PublicProduct.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<PublicProduct> getProduct(String slug, String productId) async {
    final response = await _client.get('/public/tienda/$slug/product/$productId');
    return PublicProduct.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> submitContact(String slug, ContactFormDTO data) async {
    final response = await _client.post(
      '/public/tienda/$slug/contact',
      data: data.toJson(),
    );
    return Map<String, dynamic>.from(response.data);
  }
}
