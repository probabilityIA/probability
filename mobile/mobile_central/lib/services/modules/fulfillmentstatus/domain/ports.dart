import 'entities.dart';

abstract class IFulfillmentStatusRepository {
  Future<List<FulfillmentStatusInfo>> getFulfillmentStatuses();
}
