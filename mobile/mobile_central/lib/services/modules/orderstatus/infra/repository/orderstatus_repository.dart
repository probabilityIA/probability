import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class OrderStatusApiRepository implements IOrderStatusRepository {
  final ApiClient _client;
  OrderStatusApiRepository(this._client);

  @override
  Future<PaginatedResponse<OrderStatusMapping>> getMappings(GetOrderStatusMappingsParams? params) async {
    final response = await _client.get('/order-status-mappings', queryParameters: params?.toQueryParams());
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)?.map((e) => OrderStatusMapping.fromJson(e)).toList() ?? [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<OrderStatusMapping> getMappingById(int id) async {
    final response = await _client.get('/order-status-mappings/$id');
    return OrderStatusMapping.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<OrderStatusMapping> createMapping(CreateOrderStatusMappingDTO data) async {
    final response = await _client.post('/order-status-mappings', data: data.toJson());
    return OrderStatusMapping.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<OrderStatusMapping> updateMapping(int id, UpdateOrderStatusMappingDTO data) async {
    final response = await _client.put('/order-status-mappings/$id', data: data.toJson());
    return OrderStatusMapping.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteMapping(int id) async => await _client.delete('/order-status-mappings/$id');

  @override
  Future<OrderStatusMapping> toggleMappingActive(int id) async {
    final response = await _client.put('/order-status-mappings/$id/toggle-active');
    return OrderStatusMapping.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<List<OrderStatusInfo>> getStatuses({bool? isActive}) async {
    final qp = isActive != null ? {'is_active': isActive} : null;
    final response = await _client.get('/order-statuses', queryParameters: qp);
    return (response.data['data'] as List<dynamic>?)?.map((e) => OrderStatusInfo.fromJson(e)).toList() ?? [];
  }

  @override
  Future<OrderStatusInfo> createStatus(CreateOrderStatusDTO data) async {
    final response = await _client.post('/order-statuses', data: data.toJson());
    return OrderStatusInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<OrderStatusInfo> getStatusById(int id) async {
    final response = await _client.get('/order-statuses/$id');
    return OrderStatusInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<OrderStatusInfo> updateStatus(int id, CreateOrderStatusDTO data) async {
    final response = await _client.put('/order-statuses/$id', data: data.toJson());
    return OrderStatusInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteStatus(int id) async => await _client.delete('/order-statuses/$id');

  @override
  Future<List<IntegrationTypeInfo>> getEcommerceIntegrationTypes() async {
    final response = await _client.get('/order-statuses/ecommerce-types');
    return (response.data['data'] as List<dynamic>?)?.map((e) => IntegrationTypeInfo.fromJson(e)).toList() ?? [];
  }

  @override
  Future<List<ChannelStatusInfo>> getChannelStatuses(int integrationTypeId, {bool? isActive}) async {
    final qp = <String, dynamic>{'integration_type_id': integrationTypeId};
    if (isActive != null) qp['is_active'] = isActive;
    final response = await _client.get('/channel-statuses', queryParameters: qp);
    return (response.data['data'] as List<dynamic>?)?.map((e) => ChannelStatusInfo.fromJson(e)).toList() ?? [];
  }

  @override
  Future<ChannelStatusInfo> createChannelStatus(CreateChannelStatusDTO data) async {
    final response = await _client.post('/channel-statuses', data: data.toJson());
    return ChannelStatusInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ChannelStatusInfo> updateChannelStatus(int id, UpdateChannelStatusDTO data) async {
    final response = await _client.put('/channel-statuses/$id', data: data.toJson());
    return ChannelStatusInfo.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteChannelStatus(int id) async => await _client.delete('/channel-statuses/$id');
}
