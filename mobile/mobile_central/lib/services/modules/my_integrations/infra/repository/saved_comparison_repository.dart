import '../../../../../core/network/api_client.dart';
import '../../app/sync_providers.dart';
import '../../domain/ports.dart';
import '../../domain/saved_comparison_entities.dart';

class SavedComparisonApiRepository implements ISavedComparisonRepository {
  SavedComparisonApiRepository(this._client);

  final ApiClient _client;

  @override
  Future<FindingsReport> getFindings({int? businessId}) async {
    final response = await _client.get(
      '/integrations/sync-runs/findings',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
    );
    final data = response.data;
    return FindingsReport.fromJson(
      data is Map && data['data'] is Map<String, dynamic>
          ? data['data'] as Map<String, dynamic>
          : null,
    );
  }

  @override
  Future<DataSummary> getDataSummary({int? businessId}) async {
    final response = await _client.get(
      '/integrations/sync-runs/data-summary',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
    );
    final data = response.data;
    return DataSummary.fromJson(
      data is Map<String, dynamic> ? data : null,
    );
  }

  @override
  Future<InventoryComparePage> compareInventory(
    SyncProviderSpec spec,
    InventoryCompareQuery query,
  ) async {
    final response = await _client.post(
      spec.compareInventoryPath,
      data: query.toBody(),
    );
    final data = response.data;
    if (data is Map && data['success'] == false) {
      throw Exception(data['message'] ?? data['error'] ?? 'No se pudo leer el stock');
    }
    return InventoryComparePage.fromJson(
      data is Map && data['data'] is Map<String, dynamic>
          ? data['data'] as Map<String, dynamic>
          : <String, dynamic>{},
    );
  }
}
