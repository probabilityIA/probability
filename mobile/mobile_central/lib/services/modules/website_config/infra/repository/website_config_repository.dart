import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class WebsiteConfigApiRepository implements IWebsiteConfigRepository {
  final ApiClient _client;

  WebsiteConfigApiRepository(this._client);

  String _withBusinessId(String path, int? businessId) {
    if (businessId == null) return path;
    final sep = path.contains('?') ? '&' : '?';
    return '$path${sep}business_id=$businessId';
  }

  @override
  Future<WebsiteConfigData> getConfig({int? businessId}) async {
    final response = await _client.get(_withBusinessId('/website-config', businessId));
    return WebsiteConfigData.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<WebsiteConfigData> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) async {
    final response = await _client.put(
      _withBusinessId('/website-config', businessId),
      data: data.toJson(),
    );
    return WebsiteConfigData.fromJson(response.data['data'] ?? response.data);
  }
}
