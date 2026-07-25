import {
    Order,
    OrderHistory,
    PaginatedResponse,
    GetOrdersParams,
    SingleResponse,
    CreateOrderDTO,
    UpdateOrderDTO,
    ChangeOrderStatusDTO,
    ActionResponse,
    OrderNotifications
} from './types';

export interface IOrderRepository {
    getOrders(params?: GetOrdersParams): Promise<PaginatedResponse<Order>>;
    getOrderById(id: string): Promise<SingleResponse<Order>>;
    getOrderHistory(orderId: string): Promise<{ success: boolean; data: OrderHistory[] }>;
    createOrder(data: CreateOrderDTO): Promise<SingleResponse<Order>>;
    updateOrder(id: string, data: UpdateOrderDTO): Promise<SingleResponse<Order>>;
    changeOrderStatus(id: string, data: ChangeOrderStatusDTO): Promise<SingleResponse<Order>>;
    deleteOrder(id: string): Promise<ActionResponse>;
    getOrderNotifications(orderId: string, businessId?: number): Promise<OrderNotifications>;
    resendOrderNotification(orderId: string, action: string, businessId?: number): Promise<ActionResponse>;
    getOrderRaw(id: string): Promise<SingleResponse<any>>;
    getAIRecommendation(origin: string, destination: string): Promise<any>;
    requestConfirmation(orderId: string): Promise<ActionResponse>;
}
