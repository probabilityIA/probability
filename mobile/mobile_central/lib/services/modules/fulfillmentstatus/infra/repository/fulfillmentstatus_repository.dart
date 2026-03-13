import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class FulfillmentStatusApiRepository implements IFulfillmentStatusRepository {
  final ApiClient _client;
  FulfillmentStatusApiRepository(this._client);

  @override
  Future<List<FulfillmentStatusInfo>> getFulfillmentStatuses() async {
    final response = await _client.get('/fulfillment-statuses');
    return (response.data['data'] as List<dynamic>?)?.map((e) => FulfillmentStatusInfo.fromJson(e)).toList() ?? [];
  }
}
