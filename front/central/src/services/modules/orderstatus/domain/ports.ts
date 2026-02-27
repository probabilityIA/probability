import {
    OrderStatusMapping,
    PaginatedResponse,
    GetOrderStatusMappingsParams,
    SingleResponse,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    ActionResponse,
    OrderStatusInfo,
    CreateOrderStatusDTO,
    UpdateOrderStatusDTO,
    EcommerceIntegrationType,
    ChannelStatusInfo,
    CreateChannelStatusDTO,
    UpdateChannelStatusDTO
} from './types';

export interface IOrderStatusMappingRepository {
    getOrderStatusMappings(params?: GetOrderStatusMappingsParams): Promise<PaginatedResponse<OrderStatusMapping>>;
    getOrderStatusMappingById(id: number): Promise<SingleResponse<OrderStatusMapping>>;
    createOrderStatusMapping(data: CreateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>>;
    updateOrderStatusMapping(id: number, data: UpdateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>>;
    deleteOrderStatusMapping(id: number): Promise<ActionResponse>;
    toggleOrderStatusMappingActive(id: number): Promise<SingleResponse<OrderStatusMapping>>;
    getOrderStatuses(isActive?: boolean): Promise<{ success: boolean; data: OrderStatusInfo[]; message?: string }>;

    // CRUD para estados de Probability
    createOrderStatus(data: CreateOrderStatusDTO): Promise<SingleResponse<OrderStatusInfo>>;
    getOrderStatusById(id: number): Promise<SingleResponse<OrderStatusInfo>>;
    updateOrderStatus(id: number, data: UpdateOrderStatusDTO): Promise<SingleResponse<OrderStatusInfo>>;
    deleteOrderStatus(id: number): Promise<ActionResponse>;

    // Estados por canal de integración (ecommerce)
    getEcommerceIntegrationTypes(): Promise<{ success: boolean; data: EcommerceIntegrationType[]; message?: string }>;
    getChannelStatuses(integrationTypeId: number, isActive?: boolean): Promise<{ success: boolean; data: ChannelStatusInfo[]; message?: string }>;
    createChannelStatus(data: CreateChannelStatusDTO): Promise<SingleResponse<ChannelStatusInfo>>;
    updateChannelStatus(id: number, data: UpdateChannelStatusDTO): Promise<SingleResponse<ChannelStatusInfo>>;
    deleteChannelStatus(id: number): Promise<ActionResponse>;
}
