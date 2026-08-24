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

  static const List<String> _deliveredStatuses = ['delivered'];
  static const List<String> _returnedStatuses = [
    'refunded',
    'restocked',
    'returned',
  ];

  Future<int> _countOrders({
    required String status,
    int? businessId,
    String? startDate,
    String? endDate,
  }) async {
    try {
      final response = await _client.get(
        '/orders',
        queryParameters: <String, dynamic>{
          'page': 1,
          'page_size': 1,
          'status': status,
          'business_id': ?businessId,
          'start_date': ?startDate,
          'end_date': ?endDate,
        },
      );
      final data = response.data;
      final total = data is Map ? data['total'] : null;
      if (total is int) return total;
      if (total is num) return total.toInt();
      return 0;
    } catch (_) {
      return 0;
    }
  }

  @override
  Future<OrderEffectiveness> getEffectiveness({
    int? businessId,
    String? startDate,
    String? endDate,
  }) async {
    Future<int> count(String status) => _countOrders(
          status: status,
          businessId: businessId,
          startDate: startDate,
          endDate: endDate,
        );

    final counts = await Future.wait([
      ..._deliveredStatuses.map(count),
      ..._returnedStatuses.map(count),
    ]);

    final delivered = counts
        .take(_deliveredStatuses.length)
        .fold<int>(0, (acc, value) => acc + value);
    final returned = counts
        .skip(_deliveredStatuses.length)
        .fold<int>(0, (acc, value) => acc + value);

    return OrderEffectiveness(delivered: delivered, returned: returned);
  }
}
