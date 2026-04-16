import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class RouteUseCases {
  final IRouteRepository _repository;
  RouteUseCases(this._repository);

  Future<PaginatedResponse<RouteInfo>> getRoutes(GetRoutesParams? params) => _repository.getRoutes(params);
  Future<RouteDetail> getRouteById(int id, {int? businessId}) => _repository.getRouteById(id, businessId: businessId);
  Future<RouteInfo> createRoute(CreateRouteDTO data, {int? businessId}) => _repository.createRoute(data, businessId: businessId);
  Future<RouteInfo> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) => _repository.updateRoute(id, data, businessId: businessId);
  Future<Map<String, dynamic>> deleteRoute(int id, {int? businessId}) => _repository.deleteRoute(id, businessId: businessId);

  Future<RouteDetail> startRoute(int id, {int? businessId}) => _repository.startRoute(id, businessId: businessId);
  Future<RouteDetail> completeRoute(int id, {int? businessId}) => _repository.completeRoute(id, businessId: businessId);

  Future<RouteStopInfo> addStop(int routeId, AddStopDTO data, {int? businessId}) => _repository.addStop(routeId, data, businessId: businessId);
  Future<RouteStopInfo> updateStop(int routeId, int stopId, UpdateStopDTO data, {int? businessId}) => _repository.updateStop(routeId, stopId, data, businessId: businessId);
  Future<Map<String, dynamic>> deleteStop(int routeId, int stopId, {int? businessId}) => _repository.deleteStop(routeId, stopId, businessId: businessId);
  Future<RouteStopInfo> updateStopStatus(int routeId, int stopId, UpdateStopStatusDTO data, {int? businessId}) => _repository.updateStopStatus(routeId, stopId, data, businessId: businessId);
  Future<RouteDetail> reorderStops(int routeId, ReorderStopsDTO data, {int? businessId}) => _repository.reorderStops(routeId, data, businessId: businessId);

  Future<List<DriverOption>> getAvailableDrivers({int? businessId}) => _repository.getAvailableDrivers(businessId: businessId);
  Future<List<VehicleOption>> getAvailableVehicles({int? businessId}) => _repository.getAvailableVehicles(businessId: businessId);
  Future<List<AssignableOrder>> getAssignableOrders({int? businessId}) => _repository.getAssignableOrders(businessId: businessId);
}
