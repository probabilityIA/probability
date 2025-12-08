import { IOrderRepository } from '../domain/ports';
import {
    GetOrdersParams,
    CreateOrderDTO,
    UpdateOrderDTO
} from '../domain/types';

export class OrderUseCases {
    constructor(private repository: IOrderRepository) { }

    async getOrders(params?: GetOrdersParams) {
        return this.repository.getOrders(params);
    }

    async getOrderById(id: string) {
        return this.repository.getOrderById(id);
    }

    async createOrder(data: CreateOrderDTO) {
        return this.repository.createOrder(data);
    }

    async updateOrder(id: string, data: UpdateOrderDTO) {
        return this.repository.updateOrder(id, data);
    }

    async deleteOrder(id: string) {
        return this.repository.deleteOrder(id);
    }

    async getOrderRaw(id: string) {
        return this.repository.getOrderRaw(id);
    }
}
