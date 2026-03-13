import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/route_repository.dart';

class RouteProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<RouteInfo> _routes = [];
  RouteDetail? _selectedRoute;
  List<DriverOption> _availableDrivers = [];
  List<VehicleOption> _availableVehicles = [];
  List<AssignableOrder> _assignableOrders = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  RouteProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<RouteInfo> get routes => _routes;
  RouteDetail? get selectedRoute => _selectedRoute;
  List<DriverOption> get availableDrivers => _availableDrivers;
  List<VehicleOption> get availableVehicles => _availableVehicles;
  List<AssignableOrder> get assignableOrders => _assignableOrders;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  RouteUseCases get _useCases => RouteUseCases(RouteApiRepository(_apiClient));

  Future<void> fetchRoutes({int? businessId, String? status, int? driverId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final params = GetRoutesParams(page: _page, pageSize: _pageSize, businessId: businessId, status: status, driverId: driverId);
      final response = await _useCases.getRoutes(params);
      _routes = response.data; _pagination = response.pagination;
    } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }

  Future<void> fetchRouteDetail(int id, {int? businessId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      _selectedRoute = await _useCases.getRouteById(id, businessId: businessId);
    } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }

  Future<RouteInfo?> createRoute(CreateRouteDTO data, {int? businessId}) async {
    try { return await _useCases.createRoute(data, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
  }

  Future<RouteInfo?> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    try { return await _useCases.updateRoute(id, data, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
  }

  Future<bool> deleteRoute(int id, {int? businessId}) async {
    try { await _useCases.deleteRoute(id, businessId: businessId); return true; } catch (e) { _error = e.toString(); notifyListeners(); return false; }
  }

  Future<RouteDetail?> startRoute(int id, {int? businessId}) async {
    try { return await _useCases.startRoute(id, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
  }

  Future<RouteDetail?> completeRoute(int id, {int? businessId}) async {
    try { return await _useCases.completeRoute(id, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
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
    } catch (e) { _error = e.toString(); notifyListeners(); }
  }

  void setPage(int page) { _page = page; }
}
