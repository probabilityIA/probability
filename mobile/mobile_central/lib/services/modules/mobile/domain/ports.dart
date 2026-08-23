import 'entities.dart';

abstract class IMobileRepository {
  Future<MobileOrderFull> getOrderFull(String orderId, {int? businessId});
}
