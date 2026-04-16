import '../domain/entities.dart';
import '../domain/ports.dart';

class WebsiteConfigUseCases {
  final IWebsiteConfigRepository _repository;

  WebsiteConfigUseCases(this._repository);

  Future<WebsiteConfigData> getConfig({int? businessId}) {
    return _repository.getConfig(businessId: businessId);
  }

  Future<WebsiteConfigData> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) {
    return _repository.updateConfig(data, businessId: businessId);
  }
}
