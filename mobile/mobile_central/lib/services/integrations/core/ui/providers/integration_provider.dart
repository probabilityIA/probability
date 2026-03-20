import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/integration_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class IntegrationProvider extends ChangeNotifier {
  final ApiClient _apiClient;

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

  // Integration Types
  List<IntegrationType> _integrationTypes = [];
  bool _isLoadingTypes = false;

  // Integration Categories
  List<IntegrationCategory> _integrationCategories = [];
  bool _isLoadingCategories = false;

  IntegrationProvider({required ApiClient apiClient})
      : _apiClient = apiClient;

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

  IntegrationUseCases get _useCases =>
      IntegrationUseCases(IntegrationApiRepository(_apiClient));

  // ============================================
  // Integrations
  // ============================================

  Future<void> fetchIntegrations({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

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
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Integration?> getIntegrationById(int id) async {
    try {
      return await _useCases.getIntegrationById(id);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Integration?> getIntegrationByType(String type,
      {int? businessId}) async {
    try {
      return await _useCases.getIntegrationByType(type,
          businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Integration?> createIntegration(CreateIntegrationDTO data) async {
    try {
      final integration = await _useCases.createIntegration(data);
      return integration;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Integration?> updateIntegration(
      int id, UpdateIntegrationDTO data) async {
    try {
      final integration = await _useCases.updateIntegration(id, data);
      return integration;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> deleteIntegration(int id) async {
    try {
      await _useCases.deleteIntegration(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<ActionResponse?> testConnection(int id) async {
    try {
      return await _useCases.testConnection(id);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> activateIntegration(int id) async {
    try {
      await _useCases.activateIntegration(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deactivateIntegration(int id) async {
    try {
      await _useCases.deactivateIntegration(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<Integration?> setAsDefault(int id) async {
    try {
      return await _useCases.setAsDefault(id);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<ActionResponse?> syncOrders(int id,
      {SyncOrdersParams? params}) async {
    try {
      return await _useCases.syncOrders(id, params: params);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> getSyncStatus(int id,
      {int? businessId}) async {
    try {
      return await _useCases.getSyncStatus(id, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<ActionResponse?> testConnectionRaw(
      String typeCode, Map<String, dynamic> config,
      Map<String, dynamic> credentials) async {
    try {
      return await _useCases.testConnectionRaw(
          typeCode, config, credentials);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  // ============================================
  // Integration Types
  // ============================================

  Future<void> fetchIntegrationTypes({int? categoryId}) async {
    _isLoadingTypes = true;
    notifyListeners();

    try {
      _integrationTypes =
          await _useCases.getIntegrationTypes(categoryId: categoryId);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoadingTypes = false;
    notifyListeners();
  }

  Future<void> fetchActiveIntegrationTypes() async {
    _isLoadingTypes = true;
    notifyListeners();

    try {
      _integrationTypes = await _useCases.getActiveIntegrationTypes();
    } catch (e) {
      _error = parseError(e);
    }

    _isLoadingTypes = false;
    notifyListeners();
  }

  Future<IntegrationType?> getIntegrationTypeById(int id) async {
    try {
      return await _useCases.getIntegrationTypeById(id);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<IntegrationType?> getIntegrationTypeByCode(String code) async {
    try {
      return await _useCases.getIntegrationTypeByCode(code);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  // ============================================
  // Integration Categories
  // ============================================

  Future<void> fetchIntegrationCategories() async {
    _isLoadingCategories = true;
    notifyListeners();

    try {
      _integrationCategories = await _useCases.getIntegrationCategories();
    } catch (e) {
      _error = parseError(e);
    }

    _isLoadingCategories = false;
    notifyListeners();
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
    notifyListeners();
  }
}
