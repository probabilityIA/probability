import '../domain/entities.dart';
import '../domain/ports.dart';

class MobileUseCases {
  final IMobileRepository _repository;

  MobileUseCases(this._repository);

  Future<MobileOrderFull> getOrderFull(String orderId, {int? businessId}) {
    return _repository.getOrderFull(orderId, businessId: businessId);
  }
}
