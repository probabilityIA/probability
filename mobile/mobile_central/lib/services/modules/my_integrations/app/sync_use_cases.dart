import '../domain/ports.dart';
import '../domain/sync_entities.dart';
import 'sync_providers.dart';

class SyncRunsUseCases {
  SyncRunsUseCases(this._repository);

  final ISyncRunsRepository _repository;

  Future<List<SyncRunRecord>> listLastRuns({int? businessId}) =>
      _repository.listLastRuns(businessId: businessId);

  Future<List<SyncRunDetail>> listRunItems(SyncRunItemsQuery query) =>
      _repository.listRunItems(query);

  Future<SyncStartResult> syncInventory(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) =>
      _repository.syncInventory(spec, integrationId, businessId: businessId, skus: skus);

  Future<SyncStartResult> reconcileProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
  }) =>
      _repository.reconcileProducts(spec, integrationId, businessId: businessId);

  Future<SyncStartResult> associateProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) =>
      _repository.associateProducts(spec, integrationId, businessId: businessId, skus: skus);

  Future<SyncStartResult> applyProducts(
    SyncProviderSpec spec,
    int integrationId,
    String action, {
    int? businessId,
    List<String>? skus,
  }) =>
      _repository.applyProducts(spec, integrationId, action, businessId: businessId, skus: skus);
}
