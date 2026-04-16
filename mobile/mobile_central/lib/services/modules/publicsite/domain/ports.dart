import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IPublicSiteRepository {
  Future<PublicBusiness> getBusinessPage(String slug);
  Future<PaginatedResponse<PublicProduct>> getCatalog(String slug, GetPublicCatalogParams? params);
  Future<PublicProduct> getProduct(String slug, String productId);
  Future<Map<String, dynamic>> submitContact(String slug, ContactFormDTO data);
}
