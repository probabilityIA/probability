import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IIntegrationRepository {
  // Integrations
  Future<PaginatedResponse<Integration>> getIntegrations(
      GetIntegrationsParams? params);
  Future<Integration> getIntegrationById(int id);
  Future<Integration> getIntegrationByType(String type, {int? businessId});
  Future<Integration> createIntegration(CreateIntegrationDTO data);
  Future<Integration> updateIntegration(int id, UpdateIntegrationDTO data);
  Future<ActionResponse> deleteIntegration(int id);
  Future<ActionResponse> testConnection(int id);
  Future<ActionResponse> activateIntegration(int id);
  Future<ActionResponse> deactivateIntegration(int id);
  Future<Integration> setAsDefault(int id);
  Future<ActionResponse> syncOrders(int id, {SyncOrdersParams? params});
  Future<Map<String, dynamic>> getSyncStatus(int id, {int? businessId});
  Future<ActionResponse> testIntegration(int id);
  Future<ActionResponse> testConnectionRaw(
      String typeCode, Map<String, dynamic> config,
      Map<String, dynamic> credentials);
  Future<WebhookInfo> getWebhookUrl(int id);

  // Integration Types
  Future<List<IntegrationType>> getIntegrationTypes({int? categoryId});
  Future<List<IntegrationType>> getActiveIntegrationTypes();
  Future<IntegrationType> getIntegrationTypeById(int id);
  Future<IntegrationType> getIntegrationTypeByCode(String code);
  Future<IntegrationType> createIntegrationType(CreateIntegrationTypeDTO data);
  Future<IntegrationType> updateIntegrationType(
      int id, UpdateIntegrationTypeDTO data);
  Future<ActionResponse> deleteIntegrationType(int id);
  Future<Map<String, String>> getIntegrationTypePlatformCredentials(int id);

  // Integration Categories
  Future<List<IntegrationCategory>> getIntegrationCategories();
}
