import {
    OrderStatusMapping,
    PaginatedResponse,
    GetOrderStatusMappingsParams,
    SingleResponse,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    ActionResponse
} from './types';

export interface IOrderStatusMappingRepository {
    getOrderStatusMappings(params?: GetOrderStatusMappingsParams): Promise<PaginatedResponse<OrderStatusMapping>>;
    getOrderStatusMappingById(id: number): Promise<SingleResponse<OrderStatusMapping>>;
    createOrderStatusMapping(data: CreateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>>;
    updateOrderStatusMapping(id: number, data: UpdateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>>;
    deleteOrderStatusMapping(id: number): Promise<ActionResponse>;
    toggleOrderStatusMappingActive(id: number): Promise<SingleResponse<OrderStatusMapping>>;
}
