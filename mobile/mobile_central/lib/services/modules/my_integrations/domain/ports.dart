import '../../../../shared/types/paginated_response.dart';
import '../app/sync_providers.dart';
import 'entities.dart';
import 'orders_compare_entities.dart';
import 'saved_comparison_entities.dart';
import 'sync_entities.dart';

abstract class IMyIntegrationsRepository {
  Future<PaginatedResponse<MyIntegration>> getIntegrations(GetMyIntegrationsParams? params);
  Future<MyIntegration> getIntegrationById(int id, {int? businessId});
  Future<List<IntegrationStats>> getStats({int? businessId});
}

abstract class ISyncRunsRepository {
  Future<List<SyncRunRecord>> listLastRuns({int? businessId});

  Future<List<SyncRunDetail>> listRunItems(SyncRunItemsQuery query);

  Future<SyncStartResult> syncInventory(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  });

  Future<SyncStartResult> reconcileProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
  });

  Future<SyncStartResult> associateProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  });

  Future<SyncStartResult> applyProducts(
    SyncProviderSpec spec,
    int integrationId,
    String action, {
    int? businessId,
    List<String>? skus,
  });
}

abstract class IOrdersCompareRepository {
  Future<OrdersComparePage> compare(OrdersCompareQuery query);

  Future<OrdersApplyResult> apply(
    int integrationId,
    List<String> externalIds, {
    int? businessId,
  });
}

abstract class ISavedComparisonRepository {
  Future<FindingsReport> getFindings({int? businessId});

  Future<DataSummary> getDataSummary({int? businessId});

  Future<FindingItemsPage> getFindingItems({
    required String code,
    int? businessId,
    int page,
    int pageSize,
    String? search,
  });

  Future<MatrixPage> getMatchMatrix({
    int? businessId,
    int page,
    int pageSize,
    String? search,
    String searchBy,
    List<int> presentIn,
    List<int> missingIn,
  });

  Future<InventoryComparePage> compareInventory(
    SyncProviderSpec spec,
    InventoryCompareQuery query,
  );
}
