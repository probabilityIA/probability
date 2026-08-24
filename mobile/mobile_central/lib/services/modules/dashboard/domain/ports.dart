import 'entities.dart';

abstract class IDashboardRepository {
  Future<DashboardStatsResponse> getStats({
    int? businessId,
    int? integrationId,
    String? startDate,
    String? endDate,
  });

  Future<OrderEffectiveness> getEffectiveness({
    int? businessId,
    String? startDate,
    String? endDate,
  });
}
