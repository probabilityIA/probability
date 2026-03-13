import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class PayGatewayApiRepository implements IPayGatewayRepository {
  final ApiClient _client;

  PayGatewayApiRepository(this._client);

  @override
  Future<PaymentGatewayTypesResponse> listPaymentGatewayTypes() async {
    final response = await _client.get('/integration-types/active');
    final data = response.data;

    final allTypes = (data['data'] as List<dynamic>?) ?? [];

    final paymentTypes = allTypes
        .where((it) =>
            (it['category']?['name'] == 'Pagos') ||
            (it['integration_category']?['name'] == 'Pagos'))
        .map((it) => PaymentGatewayType(
              id: it['id'] ?? 0,
              name: it['name'] ?? '',
              code: it['code'] ?? '',
              imageUrl: it['image_url'],
              isActive: it['is_active'] ?? true,
              inDevelopment: it['in_development'] ?? false,
            ))
        .toList();

    return PaymentGatewayTypesResponse(
      success: true,
      data: paymentTypes,
    );
  }
}
