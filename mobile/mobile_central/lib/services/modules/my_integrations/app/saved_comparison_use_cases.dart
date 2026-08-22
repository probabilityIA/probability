import '../domain/ports.dart';
import '../domain/saved_comparison_entities.dart';
import 'sync_providers.dart';

class SavedComparisonUseCases {
  SavedComparisonUseCases(this._repository);

  final ISavedComparisonRepository _repository;

  Future<FindingsReport> getFindings({int? businessId}) =>
      _repository.getFindings(businessId: businessId);

  Future<DataSummary> getDataSummary({int? businessId}) =>
      _repository.getDataSummary(businessId: businessId);

  Future<InventoryComparePage> compareInventory(
    SyncProviderSpec spec,
    InventoryCompareQuery query,
  ) =>
      _repository.compareInventory(spec, query);
}
