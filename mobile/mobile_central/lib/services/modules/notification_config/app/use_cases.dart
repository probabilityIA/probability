import '../domain/entities.dart';
import '../domain/ports.dart';

class NotificationConfigUseCases {
  final INotificationConfigRepository _repository;

  NotificationConfigUseCases(this._repository);

  Future<NotificationConfig> create(CreateConfigDTO dto, {int? businessId}) {
    return _repository.create(dto, businessId: businessId);
  }

  Future<NotificationConfig> getById(int id, {int? businessId}) {
    return _repository.getById(id, businessId: businessId);
  }

  Future<NotificationConfig> update(int id, UpdateConfigDTO dto, {int? businessId}) {
    return _repository.update(id, dto, businessId: businessId);
  }

  Future<void> delete(int id, {int? businessId}) {
    return _repository.delete(id, businessId: businessId);
  }

  Future<List<NotificationConfig>> list({ConfigFilter? filter}) {
    return _repository.list(filter: filter);
  }

  Future<SyncConfigsResponse> syncByIntegration(SyncConfigsDTO dto, {int? businessId}) {
    return _repository.syncByIntegration(dto, businessId: businessId);
  }
}
