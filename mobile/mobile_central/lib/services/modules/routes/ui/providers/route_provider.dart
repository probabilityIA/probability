import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/route_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class RouteProvider extends ChangeNotifier {
  RouteProvider({required ApiClient apiClient}) : _apiClient = apiClient {
    list = PagedListController<RouteInfo>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  late final PagedListController<RouteInfo> list;

  RouteDetail? _selectedRoute;
  List<DriverOption> _availableDrivers = [];
  List<VehicleOption> _availableVehicles = [];
  List<AssignableOrder> _assignableOrders = [];
  bool _isLoadingDetail = false;
  String? _error;
  int? _businessId;
  String? _statusFilter;
  int? _driverFilter;

  List<RouteInfo> get routes => list.loadedItems;
  RouteDetail? get selectedRoute => _selectedRoute;
  List<DriverOption> get availableDrivers => _availableDrivers;
  List<VehicleOption> get availableVehicles => _availableVehicles;
  List<AssignableOrder> get assignableOrders => _assignableOrders;
  bool get isLoading => list.isLoading || _isLoadingDetail;
  String? get error => _error ?? list.error;

  RouteUseCases get _useCases => RouteUseCases(RouteApiRepository(_apiClient));

  Future<PaginatedResponse<RouteInfo>> _fetchPage(int page, int pageSize) {
    return _useCases.getRoutes(GetRoutesParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      status: _statusFilter,
      driverId: _driverFilter,
    ));
  }

  Future<void> fetchRoutes({int? businessId, String? status, int? driverId}) {
    _businessId = businessId;
    _statusFilter = status;
    _driverFilter = driverId;
    _error = null;
    return list.refresh();
  }

  Future<void> fetchRouteDetail(int id, {int? businessId}) async {
    _isLoadingDetail = true; _error = null; notifyListeners();
    try {
      _selectedRoute = await _useCases.getRouteById(id, businessId: businessId);
    } catch (e) { _error = parseError(e); }
    _isLoadingDetail = false; notifyListeners();
  }

  Future<RouteInfo?> createRoute(CreateRouteDTO data, {int? businessId}) async {
    try { return await _useCases.createRoute(data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<RouteInfo?> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    try { return await _useCases.updateRoute(id, data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<bool> deleteRoute(int id, {int? businessId}) async {
    try { await _useCases.deleteRoute(id, businessId: businessId); return true; } catch (e) { _error = parseError(e); notifyListeners(); return false; }
  }

  Future<RouteDetail?> startRoute(int id, {int? businessId}) async {
    try { return await _useCases.startRoute(id, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<RouteDetail?> completeRoute(int id, {int? businessId}) async {
    try { return await _useCases.completeRoute(id, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<void> fetchFormOptions({int? businessId}) async {
    try {
      final results = await Future.wait([
        _useCases.getAvailableDrivers(businessId: businessId),
        _useCases.getAvailableVehicles(businessId: businessId),
        _useCases.getAssignableOrders(businessId: businessId),
      ]);
      _availableDrivers = results[0] as List<DriverOption>;
      _availableVehicles = results[1] as List<VehicleOption>;
      _assignableOrders = results[2] as List<AssignableOrder>;
      notifyListeners();
    } catch (e) { _error = parseError(e); notifyListeners(); }
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
