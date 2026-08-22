import '../../../../../core/network/api_client.dart';
import '../../domain/orders_compare_entities.dart';
import '../../domain/ports.dart';

class OrdersCompareApiRepository implements IOrdersCompareRepository {
  OrdersCompareApiRepository(this._client);

  final ApiClient _client;

  @override
  Future<OrdersComparePage> compare(OrdersCompareQuery query) async {
    final response = await _client.get(
      '/orders-compare',
      queryParameters: query.toQueryParams(),
    );
    final data = response.data;
    final payload = data is Map && data['data'] is Map<String, dynamic>
        ? data['data'] as Map<String, dynamic>
        : <String, dynamic>{};
    return OrdersComparePage.fromJson(payload);
  }

  @override
  Future<OrdersApplyResult> apply(
    int integrationId,
    List<String> externalIds, {
    int? businessId,
  }) async {
    final response = await _client.post(
      '/orders-compare/apply',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
      data: <String, dynamic>{
        'integration_id': integrationId,
        'external_ids': externalIds,
      },
    );
    final data = response.data;
    return OrdersApplyResult.fromJson(
      data is Map && data['data'] is Map<String, dynamic>
          ? data['data'] as Map<String, dynamic>
          : null,
    );
  }
}
