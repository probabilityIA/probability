import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/my_integrations_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class MyIntegrationsProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final MyIntegrationsUseCases? _injectedUseCases;

  List<MyIntegration> _integrations = [];
  MyIntegration? _selectedIntegration;
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _categoryFilter = '';
  bool? _activeFilter;

  MyIntegrationsProvider({required ApiClient apiClient, MyIntegrationsUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases;

  List<MyIntegration> get integrations => _integrations;
  MyIntegration? get selectedIntegration => _selectedIntegration;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  MyIntegrationsUseCases get _useCases =>
      _injectedUseCases ?? MyIntegrationsUseCases(MyIntegrationsApiRepository(_apiClient));

  Future<void> fetchIntegrations({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetMyIntegrationsParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        categoryCode: _categoryFilter.isNotEmpty ? _categoryFilter : null,
        isActive: _activeFilter,
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

  Future<void> fetchIntegrationById(int id, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _selectedIntegration = await _useCases.getIntegrationById(id, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({String? categoryCode, bool? isActive}) {
    _categoryFilter = categoryCode ?? _categoryFilter;
    _activeFilter = isActive ?? _activeFilter;
    _page = 1;
  }

  void resetFilters() {
    _categoryFilter = '';
    _activeFilter = null;
    _page = 1;
  }
}
