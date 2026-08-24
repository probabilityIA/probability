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
    String searchBy = 'all',
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
        if (search != null && search.isNotEmpty && searchBy != 'all')
          'search_by': searchBy,
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
  Future<DataPreview> getDataPreview({
    required int integrationId,
    required String field,
    required DataMode mode,
    int? businessId,
  }) async {
    final response = await _client.post(
      '/integrations/sync-runs/data-preview',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
      data: <String, dynamic>{
        'integration_id': integrationId,
        'field': field,
        'mode': mode.code,
        'business_id': ?businessId,
      },
    );
    final data = response.data;
    return DataPreview.fromJson(
      data is Map<String, dynamic> ? data : <String, dynamic>{},
    );
  }

  @override
  Future<DataApplyResult> applyChannelData({
    required int integrationId,
    required String field,
    required DataMode mode,
    int? businessId,
  }) async {
    final response = await _client.post(
      '/products/channel-data/apply',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
      data: <String, dynamic>{
        'integration_id': integrationId,
        'field': field,
        'mode': mode.code,
      },
    );
    final data = response.data;
    if (data is Map && data['success'] == false) {
      throw Exception(data['message'] ?? data['error'] ?? 'No se pudo aplicar');
    }
    return DataApplyResult.fromJson(
      data is Map && data['data'] is Map<String, dynamic>
          ? data['data'] as Map<String, dynamic>
          : null,
    );
  }

  @override
  Future<int> undoChannelData({
    required String batchId,
    int? businessId,
  }) async {
    final response = await _client.post(
      '/products/channel-data/undo',
      queryParameters: <String, dynamic>{'business_id': ?businessId},
      data: <String, dynamic>{'batch_id': batchId},
    );
    final data = response.data;
    if (data is Map && data['success'] == false) {
      throw Exception(data['message'] ?? data['error'] ?? 'No se pudo deshacer');
    }
    final payload = data is Map ? data['data'] : null;
    if (payload is Map && payload['reverted'] != null) {
      return int.tryParse(payload['reverted'].toString()) ?? 0;
    }
    return 0;
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
