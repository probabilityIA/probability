import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IOrderRepository {
  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params);
  Future<Order> getOrderById(String id);
  Future<Order> createOrder(CreateOrderDTO data);
  Future<Order> updateOrder(String id, UpdateOrderDTO data);
  Future<void> deleteOrder(String id);
  Future<Map<String, dynamic>> getOrderRaw(String id);
}
