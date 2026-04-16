import 'entities.dart';

abstract class IDashboardRepository {
  Future<DashboardStatsResponse> getStats({int? businessId, int? integrationId});
}
