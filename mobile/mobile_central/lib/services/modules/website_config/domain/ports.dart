import 'entities.dart';

abstract class IWebsiteConfigRepository {
  Future<WebsiteConfigData> getConfig({int? businessId});
  Future<WebsiteConfigData> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId});
}
