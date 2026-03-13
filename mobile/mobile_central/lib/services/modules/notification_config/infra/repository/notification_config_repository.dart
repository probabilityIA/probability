import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class NotificationConfigApiRepository implements INotificationConfigRepository {
  final ApiClient _client;

  NotificationConfigApiRepository(this._client);

  String _withBusinessId(String path, int? businessId) {
    if (businessId == null) return path;
    final sep = path.contains('?') ? '&' : '?';
    return '$path${sep}business_id=$businessId';
  }

  @override
  Future<NotificationConfig> create(CreateConfigDTO dto, {int? businessId}) async {
    final path = _withBusinessId('/notification-configs', businessId ?? dto.businessId);
    final response = await _client.post(path, data: dto.toJson());
    return NotificationConfig.fromJson(response.data);
  }

  @override
  Future<NotificationConfig> getById(int id, {int? businessId}) async {
    final path = _withBusinessId('/notification-configs/$id', businessId);
    final response = await _client.get(path);
    return NotificationConfig.fromJson(response.data);
  }

  @override
  Future<NotificationConfig> update(int id, UpdateConfigDTO dto, {int? businessId}) async {
    final path = _withBusinessId('/notification-configs/$id', businessId);
    final response = await _client.put(path, data: dto.toJson());
    return NotificationConfig.fromJson(response.data);
  }

  @override
  Future<void> delete(int id, {int? businessId}) async {
    final path = _withBusinessId('/notification-configs/$id', businessId);
    await _client.delete(path);
  }

  @override
  Future<List<NotificationConfig>> list({ConfigFilter? filter}) async {
    final queryParams = filter?.toQueryParams();
    final response = await _client.get(
      '/notification-configs',
      queryParameters: queryParams,
    );
    final data = response.data;
    if (data is List) {
      return data.map((e) => NotificationConfig.fromJson(e)).toList();
    }
    return [];
  }

  @override
  Future<SyncConfigsResponse> syncByIntegration(SyncConfigsDTO dto, {int? businessId}) async {
    final path = _withBusinessId('/notification-configs/sync', businessId);
    final response = await _client.put(path, data: dto.toJson());
    return SyncConfigsResponse.fromJson(response.data);
  }
}
