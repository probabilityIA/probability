import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/core/app/use_cases.dart';
import 'package:mobile_central/services/integrations/core/domain/entities.dart';
import 'package:mobile_central/services/integrations/core/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockIntegrationRepository implements IIntegrationRepository {
  final List<String> calls = [];

  // Configurable return values
  PaginatedResponse<Integration>? getIntegrationsResult;
  Integration? getIntegrationByIdResult;
  Integration? getIntegrationByTypeResult;
  Integration? createIntegrationResult;
  Integration? updateIntegrationResult;
  ActionResponse? deleteIntegrationResult;
  ActionResponse? testConnectionResult;
  ActionResponse? activateIntegrationResult;
  ActionResponse? deactivateIntegrationResult;
  Integration? setAsDefaultResult;
  ActionResponse? syncOrdersResult;
  Map<String, dynamic>? getSyncStatusResult;
  ActionResponse? testIntegrationResult;
  ActionResponse? testConnectionRawResult;
  WebhookInfo? getWebhookUrlResult;

  List<IntegrationType>? getIntegrationTypesResult;
  List<IntegrationType>? getActiveIntegrationTypesResult;
  IntegrationType? getIntegrationTypeByIdResult;
  IntegrationType? getIntegrationTypeByCodeResult;
  IntegrationType? createIntegrationTypeResult;
  IntegrationType? updateIntegrationTypeResult;
  ActionResponse? deleteIntegrationTypeResult;
  Map<String, String>? getIntegrationTypePlatformCredentialsResult;

  List<IntegrationCategory>? getIntegrationCategoriesResult;

  // Configurable error
  Exception? errorToThrow;

  // Captured arguments
  GetIntegrationsParams? capturedGetIntegrationsParams;
  int? capturedId;
  String? capturedType;
  int? capturedBusinessId;
  CreateIntegrationDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateIntegrationDTO? capturedUpdateData;
  int? capturedDeleteId;
  int? capturedTestConnectionId;
  int? capturedActivateId;
  int? capturedDeactivateId;
  int? capturedSetAsDefaultId;
  int? capturedSyncOrdersId;
  SyncOrdersParams? capturedSyncOrdersParams;
  int? capturedSyncStatusId;
  int? capturedSyncStatusBusinessId;
  int? capturedTestIntegrationId;
  String? capturedTestConnectionRawTypeCode;
  Map<String, dynamic>? capturedTestConnectionRawConfig;
  Map<String, dynamic>? capturedTestConnectionRawCredentials;
  int? capturedWebhookUrlId;

  int? capturedGetTypeCategoryId;
  int? capturedGetTypeByIdId;
  String? capturedGetTypeByCodeCode;
  CreateIntegrationTypeDTO? capturedCreateTypeData;
  int? capturedUpdateTypeId;
  UpdateIntegrationTypeDTO? capturedUpdateTypeData;
  int? capturedDeleteTypeId;
  int? capturedPlatformCredentialsId;

  @override
  Future<PaginatedResponse<Integration>> getIntegrations(
      GetIntegrationsParams? params) async {
    calls.add('getIntegrations');
    capturedGetIntegrationsParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationsResult!;
  }

  @override
  Future<Integration> getIntegrationById(int id) async {
    calls.add('getIntegrationById');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationByIdResult!;
  }

  @override
  Future<Integration> getIntegrationByType(String type,
      {int? businessId}) async {
    calls.add('getIntegrationByType');
    capturedType = type;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationByTypeResult!;
  }

  @override
  Future<Integration> createIntegration(CreateIntegrationDTO data) async {
    calls.add('createIntegration');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createIntegrationResult!;
  }

  @override
  Future<Integration> updateIntegration(
      int id, UpdateIntegrationDTO data) async {
    calls.add('updateIntegration');
    capturedUpdateId = id;
    capturedUpdateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateIntegrationResult!;
  }

  @override
  Future<ActionResponse> deleteIntegration(int id) async {
    calls.add('deleteIntegration');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteIntegrationResult!;
  }

  @override
  Future<ActionResponse> testConnection(int id) async {
    calls.add('testConnection');
    capturedTestConnectionId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return testConnectionResult!;
  }

  @override
  Future<ActionResponse> activateIntegration(int id) async {
    calls.add('activateIntegration');
    capturedActivateId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return activateIntegrationResult!;
  }

  @override
  Future<ActionResponse> deactivateIntegration(int id) async {
    calls.add('deactivateIntegration');
    capturedDeactivateId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deactivateIntegrationResult!;
  }

  @override
  Future<Integration> setAsDefault(int id) async {
    calls.add('setAsDefault');
    capturedSetAsDefaultId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return setAsDefaultResult!;
  }

  @override
  Future<ActionResponse> syncOrders(int id,
      {SyncOrdersParams? params}) async {
    calls.add('syncOrders');
    capturedSyncOrdersId = id;
    capturedSyncOrdersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return syncOrdersResult!;
  }

  @override
  Future<Map<String, dynamic>> getSyncStatus(int id,
      {int? businessId}) async {
    calls.add('getSyncStatus');
    capturedSyncStatusId = id;
    capturedSyncStatusBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getSyncStatusResult!;
  }

  @override
  Future<ActionResponse> testIntegration(int id) async {
    calls.add('testIntegration');
    capturedTestIntegrationId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return testIntegrationResult!;
  }

  @override
  Future<ActionResponse> testConnectionRaw(String typeCode,
      Map<String, dynamic> config, Map<String, dynamic> credentials) async {
    calls.add('testConnectionRaw');
    capturedTestConnectionRawTypeCode = typeCode;
    capturedTestConnectionRawConfig = config;
    capturedTestConnectionRawCredentials = credentials;
    if (errorToThrow != null) throw errorToThrow!;
    return testConnectionRawResult!;
  }

  @override
  Future<WebhookInfo> getWebhookUrl(int id) async {
    calls.add('getWebhookUrl');
    capturedWebhookUrlId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getWebhookUrlResult!;
  }

  @override
  Future<List<IntegrationType>> getIntegrationTypes(
      {int? categoryId}) async {
    calls.add('getIntegrationTypes');
    capturedGetTypeCategoryId = categoryId;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypesResult!;
  }

  @override
  Future<List<IntegrationType>> getActiveIntegrationTypes() async {
    calls.add('getActiveIntegrationTypes');
    if (errorToThrow != null) throw errorToThrow!;
    return getActiveIntegrationTypesResult!;
  }

  @override
  Future<IntegrationType> getIntegrationTypeById(int id) async {
    calls.add('getIntegrationTypeById');
    capturedGetTypeByIdId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypeByIdResult!;
  }

  @override
  Future<IntegrationType> getIntegrationTypeByCode(String code) async {
    calls.add('getIntegrationTypeByCode');
    capturedGetTypeByCodeCode = code;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypeByCodeResult!;
  }

  @override
  Future<IntegrationType> createIntegrationType(
      CreateIntegrationTypeDTO data) async {
    calls.add('createIntegrationType');
    capturedCreateTypeData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createIntegrationTypeResult!;
  }

  @override
  Future<IntegrationType> updateIntegrationType(
      int id, UpdateIntegrationTypeDTO data) async {
    calls.add('updateIntegrationType');
    capturedUpdateTypeId = id;
    capturedUpdateTypeData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateIntegrationTypeResult!;
  }

  @override
  Future<ActionResponse> deleteIntegrationType(int id) async {
    calls.add('deleteIntegrationType');
    capturedDeleteTypeId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteIntegrationTypeResult!;
  }

  @override
  Future<Map<String, String>> getIntegrationTypePlatformCredentials(
      int id) async {
    calls.add('getIntegrationTypePlatformCredentials');
    capturedPlatformCredentialsId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypePlatformCredentialsResult!;
  }

  @override
  Future<List<IntegrationCategory>> getIntegrationCategories() async {
    calls.add('getIntegrationCategories');
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationCategoriesResult!;
  }
}

// --- Helpers ---

Integration _makeIntegration({int id = 1, String name = 'TestIntegration'}) {
  return Integration(
    id: id,
    name: name,
    code: 'test-code',
    integrationTypeId: 1,
    type: 'sales_channel',
    category: 'ecommerce',
    isActive: true,
    isDefault: false,
    isTesting: false,
    config: IntegrationConfig(),
    createdById: 1,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  );
}

IntegrationType _makeIntegrationType(
    {int id = 1, String name = 'Shopify', String code = 'shopify'}) {
  return IntegrationType(
    id: id,
    name: name,
    code: code,
    isActive: true,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  );
}

IntegrationCategory _makeCategory(
    {int id = 1, String name = 'E-commerce', String code = 'ecommerce'}) {
  return IntegrationCategory(
    id: id,
    code: code,
    name: name,
    displayOrder: 1,
    isActive: true,
    isVisible: true,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  );
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 20,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

ActionResponse _makeActionResponse(
    {bool success = true, String message = 'OK'}) {
  return ActionResponse(success: success, message: message);
}

// --- Tests ---

void main() {
  late MockIntegrationRepository mockRepo;
  late IntegrationUseCases useCases;

  setUp(() {
    mockRepo = MockIntegrationRepository();
    useCases = IntegrationUseCases(mockRepo);
  });

  // ============================================
  // Integrations
  // ============================================

  group('Integrations', () {
    group('getIntegrations', () {
      test('delegates to repository and returns result', () async {
        final expected = PaginatedResponse<Integration>(
          data: [_makeIntegration()],
          pagination: _makePagination(),
        );
        mockRepo.getIntegrationsResult = expected;
        final params = GetIntegrationsParams(page: 1, pageSize: 20);

        final result = await useCases.getIntegrations(params);

        expect(result.data.length, 1);
        expect(result.data[0].name, 'TestIntegration');
        expect(mockRepo.calls, ['getIntegrations']);
        expect(mockRepo.capturedGetIntegrationsParams, params);
      });

      test('passes null params to repository', () async {
        mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
          data: [],
          pagination: _makePagination(),
        );

        await useCases.getIntegrations(null);

        expect(mockRepo.capturedGetIntegrationsParams, isNull);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Network error');

        expect(() => useCases.getIntegrations(null), throwsException);
      });
    });

    group('getIntegrationById', () {
      test('delegates to repository with correct id', () async {
        mockRepo.getIntegrationByIdResult =
            _makeIntegration(id: 42, name: 'Found');

        final result = await useCases.getIntegrationById(42);

        expect(result.id, 42);
        expect(result.name, 'Found');
        expect(mockRepo.capturedId, 42);
        expect(mockRepo.calls, ['getIntegrationById']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Not found');

        expect(() => useCases.getIntegrationById(99), throwsException);
      });
    });

    group('getIntegrationByType', () {
      test('delegates to repository with correct type', () async {
        mockRepo.getIntegrationByTypeResult =
            _makeIntegration(name: 'Shopify Integration');

        final result = await useCases.getIntegrationByType('shopify');

        expect(result.name, 'Shopify Integration');
        expect(mockRepo.capturedType, 'shopify');
        expect(mockRepo.capturedBusinessId, isNull);
        expect(mockRepo.calls, ['getIntegrationByType']);
      });

      test('passes businessId when provided', () async {
        mockRepo.getIntegrationByTypeResult = _makeIntegration();

        await useCases.getIntegrationByType('shopify', businessId: 5);

        expect(mockRepo.capturedType, 'shopify');
        expect(mockRepo.capturedBusinessId, 5);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Not found');

        expect(
            () => useCases.getIntegrationByType('invalid'), throwsException);
      });
    });

    group('createIntegration', () {
      test('delegates to repository with correct data', () async {
        final dto = CreateIntegrationDTO(
          name: 'New Shop',
          code: 'new-shop',
          integrationTypeId: 1,
          category: 'ecommerce',
        );
        mockRepo.createIntegrationResult =
            _makeIntegration(id: 99, name: 'New Shop');

        final result = await useCases.createIntegration(dto);

        expect(result.id, 99);
        expect(result.name, 'New Shop');
        expect(mockRepo.capturedCreateData, dto);
        expect(mockRepo.calls, ['createIntegration']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Validation error');
        final dto = CreateIntegrationDTO(
          name: 'Fail',
          code: 'fail',
          integrationTypeId: 1,
          category: 'ecommerce',
        );

        expect(() => useCases.createIntegration(dto), throwsException);
      });
    });

    group('updateIntegration', () {
      test('delegates to repository with correct id and data', () async {
        final dto = UpdateIntegrationDTO(name: 'Updated');
        mockRepo.updateIntegrationResult =
            _makeIntegration(id: 5, name: 'Updated');

        final result = await useCases.updateIntegration(5, dto);

        expect(result.id, 5);
        expect(result.name, 'Updated');
        expect(mockRepo.capturedUpdateId, 5);
        expect(mockRepo.capturedUpdateData, dto);
        expect(mockRepo.calls, ['updateIntegration']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Update failed');
        final dto = UpdateIntegrationDTO(name: 'Fail');

        expect(() => useCases.updateIntegration(5, dto), throwsException);
      });
    });

    group('deleteIntegration', () {
      test('delegates to repository with correct id', () async {
        mockRepo.deleteIntegrationResult = _makeActionResponse();

        final result = await useCases.deleteIntegration(7);

        expect(result.success, true);
        expect(mockRepo.capturedDeleteId, 7);
        expect(mockRepo.calls, ['deleteIntegration']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Delete failed');

        expect(() => useCases.deleteIntegration(7), throwsException);
      });
    });

    group('testConnection', () {
      test('delegates to repository with correct id', () async {
        mockRepo.testConnectionResult =
            _makeActionResponse(message: 'Connection OK');

        final result = await useCases.testConnection(10);

        expect(result.success, true);
        expect(result.message, 'Connection OK');
        expect(mockRepo.capturedTestConnectionId, 10);
        expect(mockRepo.calls, ['testConnection']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Connection refused');

        expect(() => useCases.testConnection(10), throwsException);
      });
    });

    group('activateIntegration', () {
      test('delegates to repository with correct id', () async {
        mockRepo.activateIntegrationResult = _makeActionResponse();

        final result = await useCases.activateIntegration(3);

        expect(result.success, true);
        expect(mockRepo.capturedActivateId, 3);
        expect(mockRepo.calls, ['activateIntegration']);
      });
    });

    group('deactivateIntegration', () {
      test('delegates to repository with correct id', () async {
        mockRepo.deactivateIntegrationResult = _makeActionResponse();

        final result = await useCases.deactivateIntegration(4);

        expect(result.success, true);
        expect(mockRepo.capturedDeactivateId, 4);
        expect(mockRepo.calls, ['deactivateIntegration']);
      });
    });

    group('setAsDefault', () {
      test('delegates to repository with correct id', () async {
        mockRepo.setAsDefaultResult =
            _makeIntegration(id: 8, name: 'Default');

        final result = await useCases.setAsDefault(8);

        expect(result.id, 8);
        expect(mockRepo.capturedSetAsDefaultId, 8);
        expect(mockRepo.calls, ['setAsDefault']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Cannot set default');

        expect(() => useCases.setAsDefault(8), throwsException);
      });
    });

    group('syncOrders', () {
      test('delegates to repository with correct id and no params', () async {
        mockRepo.syncOrdersResult =
            _makeActionResponse(message: 'Sync started');

        final result = await useCases.syncOrders(12);

        expect(result.success, true);
        expect(result.message, 'Sync started');
        expect(mockRepo.capturedSyncOrdersId, 12);
        expect(mockRepo.capturedSyncOrdersParams, isNull);
        expect(mockRepo.calls, ['syncOrders']);
      });

      test('delegates to repository with params', () async {
        mockRepo.syncOrdersResult = _makeActionResponse();
        final params = SyncOrdersParams(
          createdAtMin: '2026-01-01',
          status: 'open',
        );

        await useCases.syncOrders(12, params: params);

        expect(mockRepo.capturedSyncOrdersId, 12);
        expect(mockRepo.capturedSyncOrdersParams, params);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Sync failed');

        expect(() => useCases.syncOrders(12), throwsException);
      });
    });

    group('getSyncStatus', () {
      test('delegates to repository with correct id', () async {
        mockRepo.getSyncStatusResult = {'status': 'completed', 'progress': 100};

        final result = await useCases.getSyncStatus(15);

        expect(result['status'], 'completed');
        expect(result['progress'], 100);
        expect(mockRepo.capturedSyncStatusId, 15);
        expect(mockRepo.capturedSyncStatusBusinessId, isNull);
        expect(mockRepo.calls, ['getSyncStatus']);
      });

      test('passes businessId when provided', () async {
        mockRepo.getSyncStatusResult = {'status': 'running'};

        await useCases.getSyncStatus(15, businessId: 3);

        expect(mockRepo.capturedSyncStatusId, 15);
        expect(mockRepo.capturedSyncStatusBusinessId, 3);
      });
    });

    group('testIntegration', () {
      test('delegates to repository with correct id', () async {
        mockRepo.testIntegrationResult =
            _makeActionResponse(message: 'Test passed');

        final result = await useCases.testIntegration(20);

        expect(result.success, true);
        expect(result.message, 'Test passed');
        expect(mockRepo.capturedTestIntegrationId, 20);
        expect(mockRepo.calls, ['testIntegration']);
      });
    });

    group('testConnectionRaw', () {
      test('delegates to repository with correct arguments', () async {
        mockRepo.testConnectionRawResult =
            _makeActionResponse(message: 'Raw test OK');
        final config = {'api_key': 'abc123'};
        final credentials = {'secret': 'xyz'};

        final result =
            await useCases.testConnectionRaw('shopify', config, credentials);

        expect(result.success, true);
        expect(result.message, 'Raw test OK');
        expect(mockRepo.capturedTestConnectionRawTypeCode, 'shopify');
        expect(mockRepo.capturedTestConnectionRawConfig, config);
        expect(mockRepo.capturedTestConnectionRawCredentials, credentials);
        expect(mockRepo.calls, ['testConnectionRaw']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Invalid credentials');

        expect(
          () => useCases.testConnectionRaw('shopify', {}, {}),
          throwsException,
        );
      });
    });

    group('getWebhookUrl', () {
      test('delegates to repository with correct id', () async {
        mockRepo.getWebhookUrlResult = WebhookInfo(
          url: 'https://example.com/webhook',
          method: 'POST',
          description: 'Order webhook',
          events: ['order.created'],
        );

        final result = await useCases.getWebhookUrl(25);

        expect(result.url, 'https://example.com/webhook');
        expect(result.method, 'POST');
        expect(result.events, ['order.created']);
        expect(mockRepo.capturedWebhookUrlId, 25);
        expect(mockRepo.calls, ['getWebhookUrl']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Webhook not found');

        expect(() => useCases.getWebhookUrl(25), throwsException);
      });
    });
  });

  // ============================================
  // Integration Types
  // ============================================

  group('Integration Types', () {
    group('getIntegrationTypes', () {
      test('delegates to repository without categoryId', () async {
        mockRepo.getIntegrationTypesResult = [
          _makeIntegrationType(),
          _makeIntegrationType(id: 2, name: 'Amazon', code: 'amazon'),
        ];

        final result = await useCases.getIntegrationTypes();

        expect(result.length, 2);
        expect(result[0].name, 'Shopify');
        expect(result[1].code, 'amazon');
        expect(mockRepo.capturedGetTypeCategoryId, isNull);
        expect(mockRepo.calls, ['getIntegrationTypes']);
      });

      test('passes categoryId when provided', () async {
        mockRepo.getIntegrationTypesResult = [_makeIntegrationType()];

        await useCases.getIntegrationTypes(categoryId: 3);

        expect(mockRepo.capturedGetTypeCategoryId, 3);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Fetch failed');

        expect(() => useCases.getIntegrationTypes(), throwsException);
      });
    });

    group('getActiveIntegrationTypes', () {
      test('delegates to repository', () async {
        mockRepo.getActiveIntegrationTypesResult = [_makeIntegrationType()];

        final result = await useCases.getActiveIntegrationTypes();

        expect(result.length, 1);
        expect(result[0].isActive, true);
        expect(mockRepo.calls, ['getActiveIntegrationTypes']);
      });
    });

    group('getIntegrationTypeById', () {
      test('delegates to repository with correct id', () async {
        mockRepo.getIntegrationTypeByIdResult =
            _makeIntegrationType(id: 10, name: 'MercadoLibre');

        final result = await useCases.getIntegrationTypeById(10);

        expect(result.id, 10);
        expect(result.name, 'MercadoLibre');
        expect(mockRepo.capturedGetTypeByIdId, 10);
        expect(mockRepo.calls, ['getIntegrationTypeById']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Not found');

        expect(() => useCases.getIntegrationTypeById(999), throwsException);
      });
    });

    group('getIntegrationTypeByCode', () {
      test('delegates to repository with correct code', () async {
        mockRepo.getIntegrationTypeByCodeResult =
            _makeIntegrationType(name: 'Shopify', code: 'shopify');

        final result = await useCases.getIntegrationTypeByCode('shopify');

        expect(result.code, 'shopify');
        expect(mockRepo.capturedGetTypeByCodeCode, 'shopify');
        expect(mockRepo.calls, ['getIntegrationTypeByCode']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Not found');

        expect(
          () => useCases.getIntegrationTypeByCode('invalid'),
          throwsException,
        );
      });
    });

    group('createIntegrationType', () {
      test('delegates to repository with correct data', () async {
        final dto =
            CreateIntegrationTypeDTO(name: 'WhatsApp', categoryId: 2);
        mockRepo.createIntegrationTypeResult =
            _makeIntegrationType(id: 50, name: 'WhatsApp', code: 'whatsapp');

        final result = await useCases.createIntegrationType(dto);

        expect(result.id, 50);
        expect(result.name, 'WhatsApp');
        expect(mockRepo.capturedCreateTypeData, dto);
        expect(mockRepo.calls, ['createIntegrationType']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Duplicate code');
        final dto =
            CreateIntegrationTypeDTO(name: 'Fail', categoryId: 1);

        expect(() => useCases.createIntegrationType(dto), throwsException);
      });
    });

    group('updateIntegrationType', () {
      test('delegates to repository with correct id and data', () async {
        final dto = UpdateIntegrationTypeDTO(name: 'Updated Name');
        mockRepo.updateIntegrationTypeResult =
            _makeIntegrationType(id: 5, name: 'Updated Name');

        final result = await useCases.updateIntegrationType(5, dto);

        expect(result.id, 5);
        expect(result.name, 'Updated Name');
        expect(mockRepo.capturedUpdateTypeId, 5);
        expect(mockRepo.capturedUpdateTypeData, dto);
        expect(mockRepo.calls, ['updateIntegrationType']);
      });
    });

    group('deleteIntegrationType', () {
      test('delegates to repository with correct id', () async {
        mockRepo.deleteIntegrationTypeResult = _makeActionResponse();

        final result = await useCases.deleteIntegrationType(11);

        expect(result.success, true);
        expect(mockRepo.capturedDeleteTypeId, 11);
        expect(mockRepo.calls, ['deleteIntegrationType']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('In use');

        expect(() => useCases.deleteIntegrationType(11), throwsException);
      });
    });

    group('getIntegrationTypePlatformCredentials', () {
      test('delegates to repository with correct id', () async {
        mockRepo.getIntegrationTypePlatformCredentialsResult = {
          'api_key': 'key123',
          'api_secret': 'secret456',
        };

        final result =
            await useCases.getIntegrationTypePlatformCredentials(7);

        expect(result['api_key'], 'key123');
        expect(result['api_secret'], 'secret456');
        expect(mockRepo.capturedPlatformCredentialsId, 7);
        expect(mockRepo.calls, ['getIntegrationTypePlatformCredentials']);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Access denied');

        expect(
          () => useCases.getIntegrationTypePlatformCredentials(7),
          throwsException,
        );
      });
    });
  });

  // ============================================
  // Integration Categories
  // ============================================

  group('Integration Categories', () {
    group('getIntegrationCategories', () {
      test('delegates to repository and returns result', () async {
        mockRepo.getIntegrationCategoriesResult = [
          _makeCategory(),
          _makeCategory(id: 2, name: 'Payments', code: 'payments'),
        ];

        final result = await useCases.getIntegrationCategories();

        expect(result.length, 2);
        expect(result[0].name, 'E-commerce');
        expect(result[1].code, 'payments');
        expect(mockRepo.calls, ['getIntegrationCategories']);
      });

      test('returns empty list when no categories exist', () async {
        mockRepo.getIntegrationCategoriesResult = [];

        final result = await useCases.getIntegrationCategories();

        expect(result, isEmpty);
      });

      test('propagates repository errors', () async {
        mockRepo.errorToThrow = Exception('Fetch failed');

        expect(
            () => useCases.getIntegrationCategories(), throwsException);
      });
    });
  });
}
