import { IOrderStatusMappingRepository } from '../domain/ports';
import {
    GetOrderStatusMappingsParams,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO
} from '../domain/types';

export class OrderStatusMappingUseCases {
    constructor(private repository: IOrderStatusMappingRepository) { }

    async getOrderStatusMappings(params?: GetOrderStatusMappingsParams) {
        return this.repository.getOrderStatusMappings(params);
    }

    async getOrderStatusMappingById(id: number) {
        return this.repository.getOrderStatusMappingById(id);
    }

    async createOrderStatusMapping(data: CreateOrderStatusMappingDTO) {
        return this.repository.createOrderStatusMapping(data);
    }

    async updateOrderStatusMapping(id: number, data: UpdateOrderStatusMappingDTO) {
        return this.repository.updateOrderStatusMapping(id, data);
    }

    async deleteOrderStatusMapping(id: number) {
        return this.repository.deleteOrderStatusMapping(id);
    }

    async toggleOrderStatusMappingActive(id: number) {
        return this.repository.toggleOrderStatusMappingActive(id);
    }
}
