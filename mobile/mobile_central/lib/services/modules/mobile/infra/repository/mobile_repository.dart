import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class MobileApiRepository implements IMobileRepository {
  final ApiClient _client;

  MobileApiRepository(this._client);

  @override
  Future<MobileOrderFull> getOrderFull(String orderId, {int? businessId}) async {
    final response = await _client.get(
      '/mobile/orders/$orderId/full',
      queryParameters: businessId != null ? {'business_id': businessId} : null,
    );
    final data = response.data['data'] ?? response.data;
    return MobileOrderFull.fromJson(Map<String, dynamic>.from(data));
  }
}
