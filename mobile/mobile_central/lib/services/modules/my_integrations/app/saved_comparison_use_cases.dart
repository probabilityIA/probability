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

  Future<FindingItemsPage> getFindingItems({
    required String code,
    int? businessId,
    int page = 1,
    int pageSize = 50,
    String? search,
  }) =>
      _repository.getFindingItems(
        code: code,
        businessId: businessId,
        page: page,
        pageSize: pageSize,
        search: search,
      );

  Future<MatrixPage> getMatchMatrix({
    int? businessId,
    int page = 1,
    int pageSize = 20,
    String? search,
    String searchBy = 'all',
    List<int> presentIn = const <int>[],
    List<int> missingIn = const <int>[],
  }) =>
      _repository.getMatchMatrix(
        businessId: businessId,
        page: page,
        pageSize: pageSize,
        search: search,
        searchBy: searchBy,
        presentIn: presentIn,
        missingIn: missingIn,
      );

  Future<InventoryComparePage> compareInventory(
    SyncProviderSpec spec,
    InventoryCompareQuery query,
  ) =>
      _repository.compareInventory(spec, query);
}
