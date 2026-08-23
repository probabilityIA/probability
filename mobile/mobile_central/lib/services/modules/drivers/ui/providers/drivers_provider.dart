import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/drivers_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class DriverProvider extends ChangeNotifier {
  DriverProvider({required ApiClient apiClient}) : _apiClient = apiClient {
    list = PagedListController<DriverInfo>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  late final PagedListController<DriverInfo> list;

  int? _businessId;
  String? _error;
  String _searchFilter = '';
  String _statusFilter = '';

  List<DriverInfo> get drivers => list.loadedItems;
  bool get isLoading => list.isLoading;
  String? get error => _error ?? list.error;

  DriverUseCases get _useCases =>
      DriverUseCases(DriverApiRepository(_apiClient));

  Future<PaginatedResponse<DriverInfo>> _fetchPage(int page, int pageSize) {
    return _useCases.getDrivers(GetDriversParams(
      page: page,
      pageSize: pageSize,
      search: _searchFilter.isNotEmpty ? _searchFilter : null,
      status: _statusFilter.isNotEmpty ? _statusFilter : null,
      businessId: _businessId,
    ));
  }

  Future<void> fetchDrivers({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
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

  void setFilters({
    String? search,
    String? status,
  }) {
    _searchFilter = search ?? _searchFilter;
    _statusFilter = status ?? _statusFilter;
  }

  void resetFilters() {
    _searchFilter = '';
    _statusFilter = '';
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
