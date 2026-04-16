import '../domain/entities.dart';
import '../domain/ports.dart';

class DashboardUseCases {
  final IDashboardRepository _repository;

  DashboardUseCases(this._repository);

  Future<DashboardStatsResponse> getStats({int? businessId, int? integrationId}) {
    return _repository.getStats(businessId: businessId, integrationId: integrationId);
  }
}
