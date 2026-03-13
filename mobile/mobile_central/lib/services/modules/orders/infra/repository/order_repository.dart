import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class OrderApiRepository implements IOrderRepository {
  final ApiClient _client;

  OrderApiRepository(this._client);

  @override
  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params) async {
    final response = await _client.get(
      '/orders',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => Order.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<Order> getOrderById(String id) async {
    final response = await _client.get('/orders/$id');
    return Order.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Order> createOrder(CreateOrderDTO data) async {
    final response = await _client.post('/orders', data: data.toJson());
    return Order.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<Order> updateOrder(String id, UpdateOrderDTO data) async {
    final response = await _client.put('/orders/$id', data: data.toJson());
    return Order.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteOrder(String id) async {
    await _client.delete('/orders/$id');
  }

  @override
  Future<Map<String, dynamic>> getOrderRaw(String id) async {
    final response = await _client.get('/orders/$id/raw');
    return response.data['data'] ?? response.data;
  }
}
