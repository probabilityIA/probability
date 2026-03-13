import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class PaymentStatusApiRepository implements IPaymentStatusRepository {
  final ApiClient _client;
  PaymentStatusApiRepository(this._client);

  @override
  Future<List<PaymentStatusInfo>> getPaymentStatuses({bool? isActive}) async {
    final qp = isActive != null ? {'is_active': isActive} : null;
    final response = await _client.get('/payment-statuses', queryParameters: qp);
    return (response.data['data'] as List<dynamic>?)?.map((e) => PaymentStatusInfo.fromJson(e)).toList() ?? [];
  }
}
