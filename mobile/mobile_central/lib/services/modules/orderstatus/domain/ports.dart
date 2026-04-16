import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IOrderStatusRepository {
  Future<PaginatedResponse<OrderStatusMapping>> getMappings(GetOrderStatusMappingsParams? params);
  Future<OrderStatusMapping> getMappingById(int id);
  Future<OrderStatusMapping> createMapping(CreateOrderStatusMappingDTO data);
  Future<OrderStatusMapping> updateMapping(int id, UpdateOrderStatusMappingDTO data);
  Future<void> deleteMapping(int id);
  Future<OrderStatusMapping> toggleMappingActive(int id);
  Future<List<OrderStatusInfo>> getStatuses({bool? isActive});
  Future<OrderStatusInfo> createStatus(CreateOrderStatusDTO data);
  Future<OrderStatusInfo> getStatusById(int id);
  Future<OrderStatusInfo> updateStatus(int id, CreateOrderStatusDTO data);
  Future<void> deleteStatus(int id);
  Future<List<IntegrationTypeInfo>> getEcommerceIntegrationTypes();
  Future<List<ChannelStatusInfo>> getChannelStatuses(int integrationTypeId, {bool? isActive});
  Future<ChannelStatusInfo> createChannelStatus(CreateChannelStatusDTO data);
  Future<ChannelStatusInfo> updateChannelStatus(int id, UpdateChannelStatusDTO data);
  Future<void> deleteChannelStatus(int id);
}
