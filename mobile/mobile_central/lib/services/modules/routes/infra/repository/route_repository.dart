import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class RouteApiRepository implements IRouteRepository {
  final ApiClient _client;
  RouteApiRepository(this._client);

  Map<String, dynamic>? _biz(int? businessId) =>
      businessId != null ? {'business_id': businessId} : null;

  @override
  Future<PaginatedResponse<RouteInfo>> getRoutes(GetRoutesParams? params) async {
    final response = await _client.get('/routes', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => RouteInfo.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<RouteDetail> getRouteById(int id, {int? businessId}) async {
    final response = await _client.get('/routes/$id', queryParameters: _biz(businessId));
    return RouteDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteInfo> createRoute(CreateRouteDTO data, {int? businessId}) async {
    final body = data.toJson();
    if (businessId != null) body['business_id'] = businessId;
    final response = await _client.post('/routes', data: body);
    return RouteInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteInfo> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    final response = await _client.put('/routes/$id', data: data.toJson());
    return RouteInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> deleteRoute(int id, {int? businessId}) async {
    final response = await _client.delete('/routes/$id', queryParameters: _biz(businessId));
    return response.data is Map<String, dynamic> ? response.data : {'message': 'Route deleted'};
  }

  @override
  Future<RouteDetail> startRoute(int id, {int? businessId}) async {
    final response = await _client.post('/routes/$id/start', data: businessId != null ? {'business_id': businessId} : null);
    return RouteDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteDetail> completeRoute(int id, {int? businessId}) async {
    final response = await _client.post('/routes/$id/complete', data: businessId != null ? {'business_id': businessId} : null);
    return RouteDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteStopInfo> addStop(int routeId, AddStopDTO data, {int? businessId}) async {
    final response = await _client.post('/routes/$routeId/stops', data: data.toJson());
    return RouteStopInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteStopInfo> updateStop(int routeId, int stopId, UpdateStopDTO data, {int? businessId}) async {
    final response = await _client.put('/routes/$routeId/stops/$stopId', data: data.toJson());
    return RouteStopInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Map<String, dynamic>> deleteStop(int routeId, int stopId, {int? businessId}) async {
    final response = await _client.delete('/routes/$routeId/stops/$stopId', queryParameters: _biz(businessId));
    return response.data is Map<String, dynamic> ? response.data : {'message': 'Stop deleted'};
  }

  @override
  Future<RouteStopInfo> updateStopStatus(int routeId, int stopId, UpdateStopStatusDTO data, {int? businessId}) async {
    final response = await _client.put('/routes/$routeId/stops/$stopId/status', data: data.toJson());
    return RouteStopInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<RouteDetail> reorderStops(int routeId, ReorderStopsDTO data, {int? businessId}) async {
    final response = await _client.put('/routes/$routeId/stops/reorder', data: data.toJson());
    return RouteDetail.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<List<DriverOption>> getAvailableDrivers({int? businessId}) async {
    final response = await _client.get('/routes/available-drivers', queryParameters: _biz(businessId));
    return (response.data['data'] as List<dynamic>?)?.map((e) => DriverOption.fromJson(e)).toList() ?? [];
  }

  @override
  Future<List<VehicleOption>> getAvailableVehicles({int? businessId}) async {
    final response = await _client.get('/routes/available-vehicles', queryParameters: _biz(businessId));
    return (response.data['data'] as List<dynamic>?)?.map((e) => VehicleOption.fromJson(e)).toList() ?? [];
  }

  @override
  Future<List<AssignableOrder>> getAssignableOrders({int? businessId}) async {
    final response = await _client.get('/routes/assignable-orders', queryParameters: _biz(businessId));
    return (response.data['data'] as List<dynamic>?)?.map((e) => AssignableOrder.fromJson(e)).toList() ?? [];
  }
}
