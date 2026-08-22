import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class DashboardApiRepository implements IDashboardRepository {
  final ApiClient _client;

  DashboardApiRepository(this._client);

  @override
  Future<DashboardStatsResponse> getStats({
    int? businessId,
    int? integrationId,
    String? startDate,
    String? endDate,
  }) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['business_id'] = businessId;
    if (integrationId != null) queryParams['integration_id'] = integrationId;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _client.get(
      '/dashboard/stats',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    return DashboardStatsResponse.fromJson(response.data);
  }
}
