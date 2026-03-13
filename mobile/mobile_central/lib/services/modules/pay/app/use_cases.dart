import '../domain/entities.dart';
import '../domain/ports.dart';

class PayUseCases {
  final IPayGatewayRepository _repository;

  PayUseCases(this._repository);

  Future<PaymentGatewayTypesResponse> listPaymentGatewayTypes() {
    return _repository.listPaymentGatewayTypes();
  }
}
