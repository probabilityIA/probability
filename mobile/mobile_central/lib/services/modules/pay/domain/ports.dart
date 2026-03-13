import 'entities.dart';

abstract class IPayGatewayRepository {
  Future<PaymentGatewayTypesResponse> listPaymentGatewayTypes();
}
