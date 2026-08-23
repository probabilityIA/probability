import '../../../../../core/network/api_client.dart';
import '../../app/sync_providers.dart';
import '../../domain/ports.dart';
import '../../domain/sync_entities.dart';

class SyncRunsApiRepository implements ISyncRunsRepository {
  SyncRunsApiRepository(this._client);

  final ApiClient _client;

  @override
  Future<List<SyncRunRecord>> listLastRuns({int? businessId}) async {
    final response = await _client.get(
      '/integrations/sync-runs',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
    );
    final data = response.data;
    final rows = (data is Map ? data['data'] : data) as List<dynamic>?;
    return (rows ?? const <dynamic>[])
        .whereType<Map<String, dynamic>>()
        .map(SyncRunRecord.fromJson)
        .toList();
  }

  @override
  Future<List<SyncRunDetail>> listRunItems(SyncRunItemsQuery query) async {
    final response = await _client.get(
      '/integrations/sync-runs/items',
      queryParameters: query.toQueryParams(),
    );
    final data = response.data;
    final rows = (data is Map ? data['data'] : data) as List<dynamic>?;
    return (rows ?? const <dynamic>[])
        .whereType<Map<String, dynamic>>()
        .map(SyncRunDetail.fromJson)
        .toList();
  }

  @override
  Future<SyncStartResult> syncInventory(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) async {
    return _start(
      spec.syncInventoryPath,
      <String, dynamic>{
        'integration_id': integrationId,
        'business_id': ?businessId,
        if (skus != null && skus.isNotEmpty) 'skus': skus,
      },
    );
  }

  @override
  Future<SyncStartResult> reconcileProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
  }) async {
    return _start(
      spec.reconcileProductsPath,
      <String, dynamic>{
        'integration_id': integrationId,
        'business_id': ?businessId,
      },
    );
  }

  @override
  Future<SyncStartResult> associateProducts(
    SyncProviderSpec spec,
    int integrationId, {
    int? businessId,
    List<String>? skus,
  }) async {
    return _start(
      spec.associateProductsPath,
      <String, dynamic>{
        'integration_id': integrationId,
        'business_id': ?businessId,
        if (skus != null && skus.isNotEmpty) 'skus': skus,
      },
    );
  }

  @override
  Future<SyncStartResult> applyProducts(
    SyncProviderSpec spec,
    int integrationId,
    String action, {
    int? businessId,
    List<String>? skus,
  }) async {
    return _start(
      spec.applyProductsPath,
      <String, dynamic>{
        'integration_id': integrationId,
        'business_id': ?businessId,
        ...spec.applyBodyFor(action).toJson(),
        if (skus != null && skus.isNotEmpty) 'skus': skus,
      },
    );
  }

  Future<SyncStartResult> _start(String path, Map<String, dynamic> body) async {
    try {
      final response = await _client.post(path, data: body);
      final data = response.data;
      if (data is Map<String, dynamic>) return SyncStartResult.fromJson(data);
      return const SyncStartResult(success: true);
    } catch (e) {
      return SyncStartResult(success: false, message: _reason(e));
    }
  }

  String _reason(Object error) {
    final text = error.toString();
    return text.length > 160 ? text.substring(0, 160) : text;
  }
}
