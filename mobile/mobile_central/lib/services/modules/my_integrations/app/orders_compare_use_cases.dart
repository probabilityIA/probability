import '../domain/orders_compare_entities.dart';
import '../domain/ports.dart';

class OrdersCompareUseCases {
  OrdersCompareUseCases(this._repository);

  final IOrdersCompareRepository _repository;

  Future<OrdersComparePage> compare(OrdersCompareQuery query) =>
      _repository.compare(query);

  Future<OrdersApplyResult> apply(
    int integrationId,
    List<String> externalIds, {
    int? businessId,
  }) =>
      _repository.apply(integrationId, externalIds, businessId: businessId);
}
