import '../domain/entities.dart';
import '../domain/ports.dart';

class PaymentStatusUseCases {
  final IPaymentStatusRepository _repository;
  PaymentStatusUseCases(this._repository);
  Future<List<PaymentStatusInfo>> getPaymentStatuses({bool? isActive}) => _repository.getPaymentStatuses(isActive: isActive);
}
