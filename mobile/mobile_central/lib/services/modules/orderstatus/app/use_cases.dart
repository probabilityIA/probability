import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class OrderStatusUseCases {
  final IOrderStatusRepository _repository;
  OrderStatusUseCases(this._repository);

  Future<PaginatedResponse<OrderStatusMapping>> getMappings(GetOrderStatusMappingsParams? params) => _repository.getMappings(params);
  Future<OrderStatusMapping> getMappingById(int id) => _repository.getMappingById(id);
  Future<OrderStatusMapping> createMapping(CreateOrderStatusMappingDTO data) => _repository.createMapping(data);
  Future<OrderStatusMapping> updateMapping(int id, UpdateOrderStatusMappingDTO data) => _repository.updateMapping(id, data);
  Future<void> deleteMapping(int id) => _repository.deleteMapping(id);
  Future<OrderStatusMapping> toggleMappingActive(int id) => _repository.toggleMappingActive(id);
  Future<List<OrderStatusInfo>> getStatuses({bool? isActive}) => _repository.getStatuses(isActive: isActive);
  Future<OrderStatusInfo> createStatus(CreateOrderStatusDTO data) => _repository.createStatus(data);
  Future<OrderStatusInfo> getStatusById(int id) => _repository.getStatusById(id);
  Future<OrderStatusInfo> updateStatus(int id, CreateOrderStatusDTO data) => _repository.updateStatus(id, data);
  Future<void> deleteStatus(int id) => _repository.deleteStatus(id);
  Future<List<IntegrationTypeInfo>> getEcommerceIntegrationTypes() => _repository.getEcommerceIntegrationTypes();
  Future<List<ChannelStatusInfo>> getChannelStatuses(int integrationTypeId, {bool? isActive}) => _repository.getChannelStatuses(integrationTypeId, isActive: isActive);
  Future<ChannelStatusInfo> createChannelStatus(CreateChannelStatusDTO data) => _repository.createChannelStatus(data);
  Future<ChannelStatusInfo> updateChannelStatus(int id, UpdateChannelStatusDTO data) => _repository.updateChannelStatus(id, data);
  Future<void> deleteChannelStatus(int id) => _repository.deleteChannelStatus(id);
}
