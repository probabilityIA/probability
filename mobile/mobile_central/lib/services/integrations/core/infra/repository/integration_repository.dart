import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class IntegrationApiRepository implements IIntegrationRepository {
  final ApiClient _client;

  IntegrationApiRepository(this._client);

  // ============================================
  // Integrations
  // ============================================

  @override
  Future<PaginatedResponse<Integration>> getIntegrations(
      GetIntegrationsParams? params) async {
    final response = await _client.get(
      '/integrations',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => Integration.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<Integration> getIntegrationById(int id) async {
    final response = await _client.get('/integrations/$id');
    return Integration.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Integration> getIntegrationByType(String type,
      {int? businessId}) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['business_id'] = businessId;
    final response = await _client.get(
      '/integrations/type/$type',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    return Integration.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Integration> createIntegration(CreateIntegrationDTO data) async {
    final response = await _client.post(
      '/integrations',
      data: data.toJson(),
    );
    return Integration.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Integration> updateIntegration(
      int id, UpdateIntegrationDTO data) async {
    final response = await _client.put(
      '/integrations/$id',
      data: data.toJson(),
    );
    return Integration.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ActionResponse> deleteIntegration(int id) async {
    final response = await _client.delete('/integrations/$id');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<ActionResponse> testConnection(int id) async {
    final response = await _client.post('/integrations/$id/test');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<ActionResponse> activateIntegration(int id) async {
    final response = await _client.put('/integrations/$id/activate');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<ActionResponse> deactivateIntegration(int id) async {
    final response = await _client.put('/integrations/$id/deactivate');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<Integration> setAsDefault(int id) async {
    final response = await _client.put('/integrations/$id/set-default');
    return Integration.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ActionResponse> syncOrders(int id,
      {SyncOrdersParams? params}) async {
    final response = await _client.post(
      '/integrations/$id/sync',
      data: params?.toJson(),
    );
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<Map<String, dynamic>> getSyncStatus(int id,
      {int? businessId}) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['business_id'] = businessId;
    final response = await _client.get(
      '/integrations/events/sync-status/$id',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    return response.data;
  }

  @override
  Future<ActionResponse> testIntegration(int id) async {
    final response = await _client.post('/integrations/$id/test');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<ActionResponse> testConnectionRaw(
      String typeCode, Map<String, dynamic> config,
      Map<String, dynamic> credentials) async {
    final response = await _client.post(
      '/integrations/test',
      data: {
        'type_code': typeCode,
        'config': config,
        'credentials': credentials,
      },
    );
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<WebhookInfo> getWebhookUrl(int id) async {
    final response = await _client.get('/integrations/$id/webhook');
    return WebhookInfo.fromJson(response.data['data'] ?? response.data);
  }

  // ============================================
  // Integration Types
  // ============================================

  @override
  Future<List<IntegrationType>> getIntegrationTypes(
      {int? categoryId}) async {
    final queryParams = <String, dynamic>{};
    if (categoryId != null) queryParams['category_id'] = categoryId;
    final response = await _client.get(
      '/integration-types',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => IntegrationType.fromJson(e)).toList();
  }

  @override
  Future<List<IntegrationType>> getActiveIntegrationTypes() async {
    final response = await _client.get('/integration-types/active');
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => IntegrationType.fromJson(e)).toList();
  }

  @override
  Future<IntegrationType> getIntegrationTypeById(int id) async {
    final response = await _client.get('/integration-types/$id');
    return IntegrationType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<IntegrationType> getIntegrationTypeByCode(String code) async {
    final response = await _client.get('/integration-types/code/$code');
    return IntegrationType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<IntegrationType> createIntegrationType(
      CreateIntegrationTypeDTO data) async {
    final response = await _client.post(
      '/integration-types',
      data: data.toJson(),
    );
    return IntegrationType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<IntegrationType> updateIntegrationType(
      int id, UpdateIntegrationTypeDTO data) async {
    final response = await _client.put(
      '/integration-types/$id',
      data: data.toJson(),
    );
    return IntegrationType.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ActionResponse> deleteIntegrationType(int id) async {
    final response = await _client.delete('/integration-types/$id');
    return ActionResponse.fromJson(response.data);
  }

  @override
  Future<Map<String, String>> getIntegrationTypePlatformCredentials(
      int id) async {
    final response =
        await _client.get('/integration-types/$id/platform-credentials');
    final data = response.data['data'] as Map<String, dynamic>? ?? {};
    return data.map((key, value) => MapEntry(key, value.toString()));
  }

  // ============================================
  // Integration Categories
  // ============================================

  @override
  Future<List<IntegrationCategory>> getIntegrationCategories() async {
    final response = await _client.get('/integration-categories');
    final data = response.data['data'] as List<dynamic>? ?? [];
    return data.map((e) => IntegrationCategory.fromJson(e)).toList();
  }
}
