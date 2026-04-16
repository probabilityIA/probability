import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class IntegrationUseCases {
  final IIntegrationRepository _repository;

  IntegrationUseCases(this._repository);

  // Integrations
  Future<PaginatedResponse<Integration>> getIntegrations(
      GetIntegrationsParams? params) {
    return _repository.getIntegrations(params);
  }

  Future<Integration> getIntegrationById(int id) {
    return _repository.getIntegrationById(id);
  }

  Future<Integration> getIntegrationByType(String type, {int? businessId}) {
    return _repository.getIntegrationByType(type, businessId: businessId);
  }

  Future<Integration> createIntegration(CreateIntegrationDTO data) {
    return _repository.createIntegration(data);
  }

  Future<Integration> updateIntegration(int id, UpdateIntegrationDTO data) {
    return _repository.updateIntegration(id, data);
  }

  Future<ActionResponse> deleteIntegration(int id) {
    return _repository.deleteIntegration(id);
  }

  Future<ActionResponse> testConnection(int id) {
    return _repository.testConnection(id);
  }

  Future<ActionResponse> activateIntegration(int id) {
    return _repository.activateIntegration(id);
  }

  Future<ActionResponse> deactivateIntegration(int id) {
    return _repository.deactivateIntegration(id);
  }

  Future<Integration> setAsDefault(int id) {
    return _repository.setAsDefault(id);
  }

  Future<ActionResponse> syncOrders(int id, {SyncOrdersParams? params}) {
    return _repository.syncOrders(id, params: params);
  }

  Future<Map<String, dynamic>> getSyncStatus(int id, {int? businessId}) {
    return _repository.getSyncStatus(id, businessId: businessId);
  }

  Future<ActionResponse> testIntegration(int id) {
    return _repository.testIntegration(id);
  }

  Future<ActionResponse> testConnectionRaw(
      String typeCode, Map<String, dynamic> config,
      Map<String, dynamic> credentials) {
    return _repository.testConnectionRaw(typeCode, config, credentials);
  }

  Future<WebhookInfo> getWebhookUrl(int id) {
    return _repository.getWebhookUrl(id);
  }

  // Integration Types
  Future<List<IntegrationType>> getIntegrationTypes({int? categoryId}) {
    return _repository.getIntegrationTypes(categoryId: categoryId);
  }

  Future<List<IntegrationType>> getActiveIntegrationTypes() {
    return _repository.getActiveIntegrationTypes();
  }

  Future<IntegrationType> getIntegrationTypeById(int id) {
    return _repository.getIntegrationTypeById(id);
  }

  Future<IntegrationType> getIntegrationTypeByCode(String code) {
    return _repository.getIntegrationTypeByCode(code);
  }

  Future<IntegrationType> createIntegrationType(
      CreateIntegrationTypeDTO data) {
    return _repository.createIntegrationType(data);
  }

  Future<IntegrationType> updateIntegrationType(
      int id, UpdateIntegrationTypeDTO data) {
    return _repository.updateIntegrationType(id, data);
  }

  Future<ActionResponse> deleteIntegrationType(int id) {
    return _repository.deleteIntegrationType(id);
  }

  Future<Map<String, String>> getIntegrationTypePlatformCredentials(int id) {
    return _repository.getIntegrationTypePlatformCredentials(id);
  }

  // Integration Categories
  Future<List<IntegrationCategory>> getIntegrationCategories() {
    return _repository.getIntegrationCategories();
  }
}
