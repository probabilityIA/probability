import 'entities.dart';

abstract class IPaymentStatusRepository {
  Future<List<PaymentStatusInfo>> getPaymentStatuses({bool? isActive});
}
