import {
    Order,
    PaginatedResponse,
    GetOrdersParams,
    SingleResponse,
    CreateOrderDTO,
    UpdateOrderDTO,
    ActionResponse
} from './types';

export interface IOrderRepository {
    getOrders(params?: GetOrdersParams): Promise<PaginatedResponse<Order>>;
    getOrderById(id: string): Promise<SingleResponse<Order>>;
    createOrder(data: CreateOrderDTO): Promise<SingleResponse<Order>>;
    updateOrder(id: string, data: UpdateOrderDTO): Promise<SingleResponse<Order>>;
    deleteOrder(id: string): Promise<ActionResponse>;
    getOrderRaw(id: string): Promise<SingleResponse<any>>;
}
