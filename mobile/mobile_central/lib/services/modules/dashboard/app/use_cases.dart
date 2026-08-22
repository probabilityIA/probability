import '../domain/entities.dart';
import '../domain/ports.dart';

class DashboardUseCases {
  final IDashboardRepository _repository;

  DashboardUseCases(this._repository);

  Future<DashboardStatsResponse> getStats({
    int? businessId,
    int? integrationId,
    String? startDate,
    String? endDate,
  }) {
    return _repository.getStats(
      businessId: businessId,
      integrationId: integrationId,
      startDate: startDate,
      endDate: endDate,
    );
  }
}
