import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class OrderUseCases {
  final IOrderRepository _repository;

  OrderUseCases(this._repository);

  Future<PaginatedResponse<Order>> getOrders(GetOrdersParams? params) {
    return _repository.getOrders(params);
  }

  Future<Order> getOrderById(String id) {
    return _repository.getOrderById(id);
  }

  Future<Order> createOrder(CreateOrderDTO data) {
    return _repository.createOrder(data);
  }

  Future<Order> updateOrder(String id, UpdateOrderDTO data) {
    return _repository.updateOrder(id, data);
  }

  Future<void> deleteOrder(String id) {
    return _repository.deleteOrder(id);
  }

  Future<Map<String, dynamic>> getOrderRaw(String id) {
    return _repository.getOrderRaw(id);
  }
}
