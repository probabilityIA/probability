import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/integrations/core/app/use_cases.dart';
import 'package:mobile_central/services/integrations/core/domain/entities.dart';
import 'package:mobile_central/services/integrations/core/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

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
  int? capturedSyncOrdersId;
  SyncOrdersParams? capturedSyncOrdersParams;
  int? capturedSyncStatusId;
  int? capturedSyncStatusBusinessId;
  String? capturedTestConnectionRawTypeCode;
  Map<String, dynamic>? capturedTestConnectionRawConfig;
  Map<String, dynamic>? capturedTestConnectionRawCredentials;
  int? capturedGetTypeCategoryId;

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
    if (errorToThrow != null) throw errorToThrow!;
    return createIntegrationResult!;
  }

  @override
  Future<Integration> updateIntegration(
      int id, UpdateIntegrationDTO data) async {
    calls.add('updateIntegration');
    if (errorToThrow != null) throw errorToThrow!;
    return updateIntegrationResult!;
  }

  @override
  Future<ActionResponse> deleteIntegration(int id) async {
    calls.add('deleteIntegration');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteIntegrationResult!;
  }

  @override
  Future<ActionResponse> testConnection(int id) async {
    calls.add('testConnection');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return testConnectionResult!;
  }

  @override
  Future<ActionResponse> activateIntegration(int id) async {
    calls.add('activateIntegration');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return activateIntegrationResult!;
  }

  @override
  Future<ActionResponse> deactivateIntegration(int id) async {
    calls.add('deactivateIntegration');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deactivateIntegrationResult!;
  }

  @override
  Future<Integration> setAsDefault(int id) async {
    calls.add('setAsDefault');
    capturedId = id;
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
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypeByIdResult!;
  }

  @override
  Future<IntegrationType> getIntegrationTypeByCode(String code) async {
    calls.add('getIntegrationTypeByCode');
    if (errorToThrow != null) throw errorToThrow!;
    return getIntegrationTypeByCodeResult!;
  }

  @override
  Future<IntegrationType> createIntegrationType(
      CreateIntegrationTypeDTO data) async {
    calls.add('createIntegrationType');
    if (errorToThrow != null) throw errorToThrow!;
    return createIntegrationTypeResult!;
  }

  @override
  Future<IntegrationType> updateIntegrationType(
      int id, UpdateIntegrationTypeDTO data) async {
    calls.add('updateIntegrationType');
    if (errorToThrow != null) throw errorToThrow!;
    return updateIntegrationTypeResult!;
  }

  @override
  Future<ActionResponse> deleteIntegrationType(int id) async {
    calls.add('deleteIntegrationType');
    if (errorToThrow != null) throw errorToThrow!;
    return deleteIntegrationTypeResult!;
  }

  @override
  Future<Map<String, String>> getIntegrationTypePlatformCredentials(
      int id) async {
    calls.add('getIntegrationTypePlatformCredentials');
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

// --- Testable Provider ---
// The real IntegrationProvider creates use cases internally via a getter that
// depends on ApiClient. We create a testable version that accepts injected
// use cases and replicates the exact same state management logic.

class TestableIntegrationProvider {
  final IntegrationUseCases _useCases;

  List<Integration> _integrations = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String? _categoryFilter;
  String? _typeFilter;
  int? _businessIdFilter;
  bool? _isActiveFilter;
  String? _searchFilter;

  List<IntegrationType> _integrationTypes = [];
  bool _isLoadingTypes = false;

  List<IntegrationCategory> _integrationCategories = [];
  bool _isLoadingCategories = false;

  final List<String> notifications = [];

  TestableIntegrationProvider(this._useCases);

  // Getters
  List<Integration> get integrations => _integrations;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;
  List<IntegrationType> get integrationTypes => _integrationTypes;
  bool get isLoadingTypes => _isLoadingTypes;
  List<IntegrationCategory> get integrationCategories =>
      _integrationCategories;
  bool get isLoadingCategories => _isLoadingCategories;

  // Filter getters for test assertions
  String? get categoryFilter => _categoryFilter;
  String? get typeFilter => _typeFilter;
  int? get businessIdFilter => _businessIdFilter;
  bool? get isActiveFilter => _isActiveFilter;
  String? get searchFilter => _searchFilter;

  void _notifyListeners() {
    notifications.add('notified');
  }

  // ============================================
  // Integrations
  // ============================================

  Future<void> fetchIntegrations({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final params = GetIntegrationsParams(
        page: _page,
        pageSize: _pageSize,
        category: _categoryFilter,
        type: _typeFilter,
        businessId: businessId ?? _businessIdFilter,
        isActive: _isActiveFilter,
        search: _searchFilter,
      );
      final response = await _useCases.getIntegrations(params);
      _integrations = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Integration?> getIntegrationById(int id) async {
    try {
      return await _useCases.getIntegrationById(id);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Integration?> getIntegrationByType(String type,
      {int? businessId}) async {
    try {
      return await _useCases.getIntegrationByType(type,
          businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Integration?> createIntegration(CreateIntegrationDTO data) async {
    try {
      final integration = await _useCases.createIntegration(data);
      return integration;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Integration?> updateIntegration(
      int id, UpdateIntegrationDTO data) async {
    try {
      final integration = await _useCases.updateIntegration(id, data);
      return integration;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> deleteIntegration(int id) async {
    try {
      await _useCases.deleteIntegration(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<ActionResponse?> testConnection(int id) async {
    try {
      return await _useCases.testConnection(id);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> activateIntegration(int id) async {
    try {
      await _useCases.activateIntegration(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deactivateIntegration(int id) async {
    try {
      await _useCases.deactivateIntegration(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<Integration?> setAsDefault(int id) async {
    try {
      return await _useCases.setAsDefault(id);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<ActionResponse?> syncOrders(int id,
      {SyncOrdersParams? params}) async {
    try {
      return await _useCases.syncOrders(id, params: params);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> getSyncStatus(int id,
      {int? businessId}) async {
    try {
      return await _useCases.getSyncStatus(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<ActionResponse?> testConnectionRaw(String typeCode,
      Map<String, dynamic> config, Map<String, dynamic> credentials) async {
    try {
      return await _useCases.testConnectionRaw(typeCode, config, credentials);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  // ============================================
  // Integration Types
  // ============================================

  Future<void> fetchIntegrationTypes({int? categoryId}) async {
    _isLoadingTypes = true;
    _notifyListeners();

    try {
      _integrationTypes =
          await _useCases.getIntegrationTypes(categoryId: categoryId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoadingTypes = false;
    _notifyListeners();
  }

  Future<void> fetchActiveIntegrationTypes() async {
    _isLoadingTypes = true;
    _notifyListeners();

    try {
      _integrationTypes = await _useCases.getActiveIntegrationTypes();
    } catch (e) {
      _error = e.toString();
    }

    _isLoadingTypes = false;
    _notifyListeners();
  }

  Future<IntegrationType?> getIntegrationTypeById(int id) async {
    try {
      return await _useCases.getIntegrationTypeById(id);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<IntegrationType?> getIntegrationTypeByCode(String code) async {
    try {
      return await _useCases.getIntegrationTypeByCode(code);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  // ============================================
  // Integration Categories
  // ============================================

  Future<void> fetchIntegrationCategories() async {
    _isLoadingCategories = true;
    _notifyListeners();

    try {
      _integrationCategories = await _useCases.getIntegrationCategories();
    } catch (e) {
      _error = e.toString();
    }

    _isLoadingCategories = false;
    _notifyListeners();
  }

  // ============================================
  // Pagination & Filters
  // ============================================

  void setPage(int page) {
    _page = page;
  }

  void setFilters({
    String? category,
    String? type,
    int? businessId,
    bool? isActive,
    String? search,
  }) {
    _categoryFilter = category ?? _categoryFilter;
    _typeFilter = type ?? _typeFilter;
    _businessIdFilter = businessId ?? _businessIdFilter;
    _isActiveFilter = isActive ?? _isActiveFilter;
    _searchFilter = search ?? _searchFilter;
    _page = 1;
  }

  void resetFilters() {
    _categoryFilter = null;
    _typeFilter = null;
    _businessIdFilter = null;
    _isActiveFilter = null;
    _searchFilter = null;
    _page = 1;
  }

  void clearError() {
    _error = null;
    _notifyListeners();
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

Pagination _makePagination({
  int currentPage = 1,
  int total = 5,
  int lastPage = 1,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 20,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
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
  late TestableIntegrationProvider provider;

  setUp(() {
    mockRepo = MockIntegrationRepository();
    useCases = IntegrationUseCases(mockRepo);
    provider = TestableIntegrationProvider(useCases);
  });

  // ============================================
  // Initial State
  // ============================================

  group('initial state', () {
    test('starts with empty integrations list', () {
      expect(provider.integrations, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });

    test('starts with page 1', () {
      expect(provider.page, 1);
    });

    test('starts with pageSize 20', () {
      expect(provider.pageSize, 20);
    });

    test('starts with empty integration types list', () {
      expect(provider.integrationTypes, isEmpty);
    });

    test('starts with isLoadingTypes false', () {
      expect(provider.isLoadingTypes, false);
    });

    test('starts with empty integration categories list', () {
      expect(provider.integrationCategories, isEmpty);
    });

    test('starts with isLoadingCategories false', () {
      expect(provider.isLoadingCategories, false);
    });

    test('starts with null filters', () {
      expect(provider.categoryFilter, isNull);
      expect(provider.typeFilter, isNull);
      expect(provider.businessIdFilter, isNull);
      expect(provider.isActiveFilter, isNull);
      expect(provider.searchFilter, isNull);
    });
  });

  // ============================================
  // fetchIntegrations
  // ============================================

  group('fetchIntegrations', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [_makeIntegration()],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      expect(provider.notifications.length, 2);
    });

    test('populates integrations and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [
          _makeIntegration(id: 1),
          _makeIntegration(id: 2, name: 'ShopifyStore'),
        ],
        pagination: pagination,
      );

      await provider.fetchIntegrations();

      expect(provider.integrations.length, 2);
      expect(provider.integrations[0].id, 1);
      expect(provider.integrations[1].name, 'ShopifyStore');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchIntegrations();

      expect(provider.error, contains('Server error'));
      expect(provider.integrations, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      // First call fails
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchIntegrations();
      expect(provider.error, isNotNull);

      // Second call succeeds
      mockRepo.errorToThrow = null;
      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchIntegrations();

      expect(provider.error, isNull);
    });

    test('passes filter params to use cases', () async {
      provider.setFilters(
        category: 'ecommerce',
        type: 'sales_channel',
        businessId: 5,
        isActive: true,
        search: 'shop',
      );

      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      final captured = mockRepo.capturedGetIntegrationsParams;
      expect(captured, isNotNull);
      expect(captured!.category, 'ecommerce');
      expect(captured.type, 'sales_channel');
      expect(captured.businessId, 5);
      expect(captured.isActive, true);
      expect(captured.search, 'shop');
      expect(captured.page, 1);
      expect(captured.pageSize, 20);
    });

    test('businessId parameter overrides businessIdFilter', () async {
      provider.setFilters(businessId: 5);

      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations(businessId: 10);

      final captured = mockRepo.capturedGetIntegrationsParams;
      expect(captured!.businessId, 10);
    });

    test('uses businessIdFilter when no businessId parameter', () async {
      provider.setFilters(businessId: 5);

      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      final captured = mockRepo.capturedGetIntegrationsParams;
      expect(captured!.businessId, 5);
    });

    test('uses current page value', () async {
      provider.setPage(3);

      mockRepo.getIntegrationsResult = PaginatedResponse<Integration>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchIntegrations();

      final captured = mockRepo.capturedGetIntegrationsParams;
      expect(captured!.page, 3);
    });
  });

  // ============================================
  // getIntegrationById
  // ============================================

  group('getIntegrationById', () {
    test('returns Integration on success', () async {
      mockRepo.getIntegrationByIdResult =
          _makeIntegration(id: 42, name: 'FoundIntegration');

      final result = await provider.getIntegrationById(42);

      expect(result, isNotNull);
      expect(result!.id, 42);
      expect(result.name, 'FoundIntegration');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Not found');

      final result = await provider.getIntegrationById(99);

      expect(result, isNull);
      expect(provider.error, contains('Not found'));
    });
  });

  // ============================================
  // getIntegrationByType
  // ============================================

  group('getIntegrationByType', () {
    test('returns Integration on success', () async {
      mockRepo.getIntegrationByTypeResult =
          _makeIntegration(name: 'Shopify Store');

      final result = await provider.getIntegrationByType('shopify');

      expect(result, isNotNull);
      expect(result!.name, 'Shopify Store');
    });

    test('passes businessId to use case', () async {
      mockRepo.getIntegrationByTypeResult = _makeIntegration();

      await provider.getIntegrationByType('shopify', businessId: 7);

      expect(mockRepo.capturedType, 'shopify');
      expect(mockRepo.capturedBusinessId, 7);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Type not found');

      final result = await provider.getIntegrationByType('invalid');

      expect(result, isNull);
      expect(provider.error, contains('Type not found'));
    });
  });

  // ============================================
  // createIntegration
  // ============================================

  group('createIntegration', () {
    test('returns created Integration on success', () async {
      final dto = CreateIntegrationDTO(
        name: 'New Shop',
        code: 'new-shop',
        integrationTypeId: 1,
        category: 'ecommerce',
      );
      mockRepo.createIntegrationResult =
          _makeIntegration(id: 50, name: 'New Shop');

      final result = await provider.createIntegration(dto);

      expect(result, isNotNull);
      expect(result!.id, 50);
      expect(result.name, 'New Shop');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateIntegrationDTO(
        name: 'Fail',
        code: 'fail',
        integrationTypeId: 1,
        category: 'ecommerce',
      );

      final result = await provider.createIntegration(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  // ============================================
  // updateIntegration
  // ============================================

  group('updateIntegration', () {
    test('returns updated Integration on success', () async {
      final dto = UpdateIntegrationDTO(name: 'Updated');
      mockRepo.updateIntegrationResult =
          _makeIntegration(id: 5, name: 'Updated');

      final result = await provider.updateIntegration(5, dto);

      expect(result, isNotNull);
      expect(result!.id, 5);
      expect(result.name, 'Updated');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateIntegrationDTO(name: 'Fail');

      final result = await provider.updateIntegration(5, dto);

      expect(result, isNull);
      expect(provider.error, contains('Update failed'));
    });
  });

  // ============================================
  // deleteIntegration
  // ============================================

  group('deleteIntegration', () {
    test('returns true on success', () async {
      mockRepo.deleteIntegrationResult = _makeActionResponse();

      final result = await provider.deleteIntegration(7);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteIntegration(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  // ============================================
  // testConnection
  // ============================================

  group('testConnection', () {
    test('returns ActionResponse on success', () async {
      mockRepo.testConnectionResult =
          _makeActionResponse(message: 'Connection OK');

      final result = await provider.testConnection(10);

      expect(result, isNotNull);
      expect(result!.success, true);
      expect(result.message, 'Connection OK');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Connection refused');

      final result = await provider.testConnection(10);

      expect(result, isNull);
      expect(provider.error, contains('Connection refused'));
    });
  });

  // ============================================
  // activateIntegration
  // ============================================

  group('activateIntegration', () {
    test('returns true on success', () async {
      mockRepo.activateIntegrationResult = _makeActionResponse();

      final result = await provider.activateIntegration(3);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Activation failed');

      final result = await provider.activateIntegration(3);

      expect(result, false);
      expect(provider.error, contains('Activation failed'));
    });
  });

  // ============================================
  // deactivateIntegration
  // ============================================

  group('deactivateIntegration', () {
    test('returns true on success', () async {
      mockRepo.deactivateIntegrationResult = _makeActionResponse();

      final result = await provider.deactivateIntegration(4);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Deactivation failed');

      final result = await provider.deactivateIntegration(4);

      expect(result, false);
      expect(provider.error, contains('Deactivation failed'));
    });
  });

  // ============================================
  // setAsDefault
  // ============================================

  group('setAsDefault', () {
    test('returns Integration on success', () async {
      mockRepo.setAsDefaultResult =
          _makeIntegration(id: 8, name: 'DefaultStore');

      final result = await provider.setAsDefault(8);

      expect(result, isNotNull);
      expect(result!.id, 8);
      expect(result.name, 'DefaultStore');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Cannot set default');

      final result = await provider.setAsDefault(8);

      expect(result, isNull);
      expect(provider.error, contains('Cannot set default'));
    });
  });

  // ============================================
  // syncOrders
  // ============================================

  group('syncOrders', () {
    test('returns ActionResponse on success', () async {
      mockRepo.syncOrdersResult =
          _makeActionResponse(message: 'Sync started');

      final result = await provider.syncOrders(12);

      expect(result, isNotNull);
      expect(result!.success, true);
      expect(result.message, 'Sync started');
    });

    test('passes params to use case', () async {
      mockRepo.syncOrdersResult = _makeActionResponse();
      final params = SyncOrdersParams(
        createdAtMin: '2026-01-01',
        status: 'open',
      );

      await provider.syncOrders(12, params: params);

      expect(mockRepo.capturedSyncOrdersId, 12);
      expect(mockRepo.capturedSyncOrdersParams, params);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Sync failed');

      final result = await provider.syncOrders(12);

      expect(result, isNull);
      expect(provider.error, contains('Sync failed'));
    });
  });

  // ============================================
  // getSyncStatus
  // ============================================

  group('getSyncStatus', () {
    test('returns status map on success', () async {
      mockRepo.getSyncStatusResult = {
        'status': 'completed',
        'progress': 100,
      };

      final result = await provider.getSyncStatus(15);

      expect(result, isNotNull);
      expect(result!['status'], 'completed');
      expect(result['progress'], 100);
    });

    test('passes businessId to use case', () async {
      mockRepo.getSyncStatusResult = {'status': 'running'};

      await provider.getSyncStatus(15, businessId: 3);

      expect(mockRepo.capturedSyncStatusId, 15);
      expect(mockRepo.capturedSyncStatusBusinessId, 3);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Status fetch failed');

      final result = await provider.getSyncStatus(15);

      expect(result, isNull);
      expect(provider.error, contains('Status fetch failed'));
    });
  });

  // ============================================
  // testConnectionRaw
  // ============================================

  group('testConnectionRaw', () {
    test('returns ActionResponse on success', () async {
      mockRepo.testConnectionRawResult =
          _makeActionResponse(message: 'Raw test OK');

      final result = await provider.testConnectionRaw(
        'shopify',
        {'api_key': 'abc'},
        {'secret': 'xyz'},
      );

      expect(result, isNotNull);
      expect(result!.success, true);
      expect(result.message, 'Raw test OK');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Invalid config');

      final result = await provider.testConnectionRaw('shopify', {}, {});

      expect(result, isNull);
      expect(provider.error, contains('Invalid config'));
    });
  });

  // ============================================
  // fetchIntegrationTypes
  // ============================================

  group('fetchIntegrationTypes', () {
    test('populates integrationTypes on success', () async {
      mockRepo.getIntegrationTypesResult = [
        _makeIntegrationType(id: 1, name: 'Shopify'),
        _makeIntegrationType(id: 2, name: 'Amazon', code: 'amazon'),
      ];

      await provider.fetchIntegrationTypes();

      expect(provider.integrationTypes.length, 2);
      expect(provider.integrationTypes[0].name, 'Shopify');
      expect(provider.integrationTypes[1].code, 'amazon');
      expect(provider.isLoadingTypes, false);
    });

    test('sets loading state and notifies twice', () async {
      mockRepo.getIntegrationTypesResult = [];

      await provider.fetchIntegrationTypes();

      expect(provider.notifications.length, 2);
    });

    test('passes categoryId to use case', () async {
      mockRepo.getIntegrationTypesResult = [];

      await provider.fetchIntegrationTypes(categoryId: 3);

      expect(mockRepo.capturedGetTypeCategoryId, 3);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Types fetch failed');

      await provider.fetchIntegrationTypes();

      expect(provider.error, contains('Types fetch failed'));
      expect(provider.isLoadingTypes, false);
    });
  });

  // ============================================
  // fetchActiveIntegrationTypes
  // ============================================

  group('fetchActiveIntegrationTypes', () {
    test('populates integrationTypes with active types', () async {
      mockRepo.getActiveIntegrationTypesResult = [
        _makeIntegrationType(id: 1, name: 'ActiveShopify'),
      ];

      await provider.fetchActiveIntegrationTypes();

      expect(provider.integrationTypes.length, 1);
      expect(provider.integrationTypes[0].name, 'ActiveShopify');
      expect(provider.isLoadingTypes, false);
    });

    test('sets loading state and notifies twice', () async {
      mockRepo.getActiveIntegrationTypesResult = [];

      await provider.fetchActiveIntegrationTypes();

      expect(provider.notifications.length, 2);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Active types fetch failed');

      await provider.fetchActiveIntegrationTypes();

      expect(provider.error, contains('Active types fetch failed'));
      expect(provider.isLoadingTypes, false);
    });
  });

  // ============================================
  // getIntegrationTypeById
  // ============================================

  group('getIntegrationTypeById', () {
    test('returns IntegrationType on success', () async {
      mockRepo.getIntegrationTypeByIdResult =
          _makeIntegrationType(id: 10, name: 'MercadoLibre');

      final result = await provider.getIntegrationTypeById(10);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'MercadoLibre');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Type not found');

      final result = await provider.getIntegrationTypeById(999);

      expect(result, isNull);
      expect(provider.error, contains('Type not found'));
    });
  });

  // ============================================
  // getIntegrationTypeByCode
  // ============================================

  group('getIntegrationTypeByCode', () {
    test('returns IntegrationType on success', () async {
      mockRepo.getIntegrationTypeByCodeResult =
          _makeIntegrationType(name: 'Shopify', code: 'shopify');

      final result = await provider.getIntegrationTypeByCode('shopify');

      expect(result, isNotNull);
      expect(result!.code, 'shopify');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Code not found');

      final result = await provider.getIntegrationTypeByCode('invalid');

      expect(result, isNull);
      expect(provider.error, contains('Code not found'));
    });
  });

  // ============================================
  // fetchIntegrationCategories
  // ============================================

  group('fetchIntegrationCategories', () {
    test('populates integrationCategories on success', () async {
      mockRepo.getIntegrationCategoriesResult = [
        _makeCategory(id: 1, name: 'E-commerce'),
        _makeCategory(id: 2, name: 'Payments', code: 'payments'),
      ];

      await provider.fetchIntegrationCategories();

      expect(provider.integrationCategories.length, 2);
      expect(provider.integrationCategories[0].name, 'E-commerce');
      expect(provider.integrationCategories[1].code, 'payments');
      expect(provider.isLoadingCategories, false);
    });

    test('sets loading state and notifies twice', () async {
      mockRepo.getIntegrationCategoriesResult = [];

      await provider.fetchIntegrationCategories();

      expect(provider.notifications.length, 2);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Categories fetch failed');

      await provider.fetchIntegrationCategories();

      expect(provider.error, contains('Categories fetch failed'));
      expect(provider.isLoadingCategories, false);
    });
  });

  // ============================================
  // Pagination & Filters
  // ============================================

  group('setPage', () {
    test('updates page value', () {
      provider.setPage(5);

      expect(provider.page, 5);
    });

    test('can be set to 1', () {
      provider.setPage(5);
      provider.setPage(1);

      expect(provider.page, 1);
    });
  });

  group('setFilters', () {
    test('updates all filters when provided', () {
      provider.setFilters(
        category: 'ecommerce',
        type: 'sales_channel',
        businessId: 5,
        isActive: true,
        search: 'shopify',
      );

      expect(provider.categoryFilter, 'ecommerce');
      expect(provider.typeFilter, 'sales_channel');
      expect(provider.businessIdFilter, 5);
      expect(provider.isActiveFilter, true);
      expect(provider.searchFilter, 'shopify');
    });

    test('resets page to 1', () {
      provider.setPage(5);

      provider.setFilters(category: 'ecommerce');

      expect(provider.page, 1);
    });

    test('preserves existing filters when null passed', () {
      provider.setFilters(category: 'ecommerce', type: 'sales_channel');
      provider.setFilters(search: 'shop');

      expect(provider.categoryFilter, 'ecommerce');
      expect(provider.typeFilter, 'sales_channel');
      expect(provider.searchFilter, 'shop');
    });

    test('single filter update does not clear others', () {
      provider.setFilters(
        category: 'ecommerce',
        type: 'sales_channel',
        businessId: 5,
        isActive: true,
        search: 'shop',
      );

      provider.setFilters(search: 'new search');

      expect(provider.categoryFilter, 'ecommerce');
      expect(provider.typeFilter, 'sales_channel');
      expect(provider.businessIdFilter, 5);
      expect(provider.isActiveFilter, true);
      expect(provider.searchFilter, 'new search');
    });
  });

  group('resetFilters', () {
    test('clears all filters', () {
      provider.setFilters(
        category: 'ecommerce',
        type: 'sales_channel',
        businessId: 5,
        isActive: true,
        search: 'shopify',
      );

      provider.resetFilters();

      expect(provider.categoryFilter, isNull);
      expect(provider.typeFilter, isNull);
      expect(provider.businessIdFilter, isNull);
      expect(provider.isActiveFilter, isNull);
      expect(provider.searchFilter, isNull);
    });

    test('resets page to 1', () {
      provider.setPage(10);

      provider.resetFilters();

      expect(provider.page, 1);
    });
  });

  group('clearError', () {
    test('clears error and notifies', () async {
      // Trigger an error first
      mockRepo.errorToThrow = Exception('Some error');
      await provider.getIntegrationById(1);
      expect(provider.error, isNotNull);
      final notificationsBefore = provider.notifications.length;

      provider.clearError();

      expect(provider.error, isNull);
      expect(provider.notifications.length, notificationsBefore + 1);
    });

    test('can be called when no error exists', () {
      provider.clearError();

      expect(provider.error, isNull);
    });
  });
}
