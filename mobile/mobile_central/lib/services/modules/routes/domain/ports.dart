import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IRouteRepository {
  // Route CRUD
  Future<PaginatedResponse<RouteInfo>> getRoutes(GetRoutesParams? params);
  Future<RouteDetail> getRouteById(int id, {int? businessId});
  Future<RouteInfo> createRoute(CreateRouteDTO data, {int? businessId});
  Future<RouteInfo> updateRoute(int id, UpdateRouteDTO data, {int? businessId});
  Future<Map<String, dynamic>> deleteRoute(int id, {int? businessId});

  // Route lifecycle
  Future<RouteDetail> startRoute(int id, {int? businessId});
  Future<RouteDetail> completeRoute(int id, {int? businessId});

  // Stop management
  Future<RouteStopInfo> addStop(int routeId, AddStopDTO data, {int? businessId});
  Future<RouteStopInfo> updateStop(int routeId, int stopId, UpdateStopDTO data, {int? businessId});
  Future<Map<String, dynamic>> deleteStop(int routeId, int stopId, {int? businessId});
  Future<RouteStopInfo> updateStopStatus(int routeId, int stopId, UpdateStopStatusDTO data, {int? businessId});
  Future<RouteDetail> reorderStops(int routeId, ReorderStopsDTO data, {int? businessId});

  // Form options
  Future<List<DriverOption>> getAvailableDrivers({int? businessId});
  Future<List<VehicleOption>> getAvailableVehicles({int? businessId});
  Future<List<AssignableOrder>> getAssignableOrders({int? businessId});
}
