import '../domain/entities.dart';
import '../domain/ports.dart';

class FulfillmentStatusUseCases {
  final IFulfillmentStatusRepository _repository;
  FulfillmentStatusUseCases(this._repository);
  Future<List<FulfillmentStatusInfo>> getFulfillmentStatuses() => _repository.getFulfillmentStatuses();
}
