import 'entities.dart';

abstract class INotificationConfigRepository {
  Future<NotificationConfig> create(CreateConfigDTO dto, {int? businessId});
  Future<NotificationConfig> getById(int id, {int? businessId});
  Future<NotificationConfig> update(int id, UpdateConfigDTO dto, {int? businessId});
  Future<void> delete(int id, {int? businessId});
  Future<List<NotificationConfig>> list({ConfigFilter? filter});
  Future<SyncConfigsResponse> syncByIntegration(SyncConfigsDTO dto, {int? businessId});
}
