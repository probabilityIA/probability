import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class PublicSiteUseCases {
  final IPublicSiteRepository _repository;

  PublicSiteUseCases(this._repository);

  Future<PublicBusiness> getBusinessPage(String slug) {
    return _repository.getBusinessPage(slug);
  }

  Future<PaginatedResponse<PublicProduct>> getCatalog(String slug, GetPublicCatalogParams? params) {
    return _repository.getCatalog(slug, params);
  }

  Future<PublicProduct> getProduct(String slug, String productId) {
    return _repository.getProduct(slug, productId);
  }

  Future<Map<String, dynamic>> submitContact(String slug, ContactFormDTO data) {
    return _repository.submitContact(slug, data);
  }
}
