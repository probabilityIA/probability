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
  Future<FindingItemsPage> getFindingItems({
    required String code,
    int? businessId,
    int page = 1,
    int pageSize = 50,
    String? search,
  }) async {
    final response = await _client.get(
      '/integrations/sync-runs/findings/items',
      queryParameters: <String, dynamic>{
        'code': code,
        'page': page,
        'page_size': pageSize,
        'business_id': ?businessId,
        if (search != null && search.isNotEmpty) 'q': search,
      },
    );
    final data = response.data;
    return FindingItemsPage.fromJson(
      data is Map<String, dynamic> ? data : <String, dynamic>{},
    );
  }

  @override
  Future<MatrixPage> getMatchMatrix({
    int? businessId,
    int page = 1,
    int pageSize = 20,
    String? search,
    List<int> presentIn = const <int>[],
    List<int> missingIn = const <int>[],
  }) async {
    final response = await _client.get(
      '/integrations/sync-runs/matrix',
      queryParameters: <String, dynamic>{
        'page': page,
        'page_size': pageSize,
        'business_id': ?businessId,
        if (search != null && search.isNotEmpty) 'q': search,
        if (presentIn.isNotEmpty) 'present_in': presentIn.join(','),
        if (missingIn.isNotEmpty) 'missing_in': missingIn.join(','),
      },
    );
    final data = response.data;
    return MatrixPage.fromJson(
      data is Map<String, dynamic> ? data : <String, dynamic>{},
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
