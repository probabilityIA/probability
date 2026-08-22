import '../../../../shared/types/paginated_response.dart';
import '../app/sync_providers.dart';
import 'entities.dart';
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
