import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/drivers_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class DriverProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<DriverInfo> _drivers = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  String _statusFilter = '';

  DriverProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<DriverInfo> get drivers => _drivers;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;

  DriverUseCases get _useCases =>
      DriverUseCases(DriverApiRepository(_apiClient));

  Future<void> fetchDrivers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetDriversParams(
        page: _page,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        status: _statusFilter.isNotEmpty ? _statusFilter : null,
        businessId: businessId,
      );
      final response = await _useCases.getDrivers(params);
      _drivers = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<DriverInfo?> getDriverById(int id, {int? businessId}) async {
    try {
      return await _useCases.getDriverById(id, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<DriverInfo?> createDriver(CreateDriverDTO data, {int? businessId}) async {
    try {
      final driver = await _useCases.createDriver(data, businessId: businessId);
      return driver;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) async {
    try {
      await _useCases.updateDriver(id, data, businessId: businessId);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteDriver(int id, {int? businessId}) async {
    try {
      await _useCases.deleteDriver(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({
    String? search,
    String? status,
  }) {
    _searchFilter = search ?? _searchFilter;
    _statusFilter = status ?? _statusFilter;
    _page = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _statusFilter = '';
    _page = 1;
  }
}
