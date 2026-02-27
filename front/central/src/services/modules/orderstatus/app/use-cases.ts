import { IOrderStatusMappingRepository } from '../domain/ports';
import {
    GetOrderStatusMappingsParams,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    CreateOrderStatusDTO,
    UpdateOrderStatusDTO,
    CreateChannelStatusDTO,
    UpdateChannelStatusDTO
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

    async getOrderStatuses(isActive?: boolean) {
        return this.repository.getOrderStatuses(isActive);
    }

    async createOrderStatus(data: CreateOrderStatusDTO) {
        return this.repository.createOrderStatus(data);
    }

    async getOrderStatusById(id: number) {
        return this.repository.getOrderStatusById(id);
    }

    async updateOrderStatus(id: number, data: UpdateOrderStatusDTO) {
        return this.repository.updateOrderStatus(id, data);
    }

    async deleteOrderStatus(id: number) {
        return this.repository.deleteOrderStatus(id);
    }

    async getEcommerceIntegrationTypes() {
        return this.repository.getEcommerceIntegrationTypes();
    }

    async getChannelStatuses(integrationTypeId: number, isActive?: boolean) {
        return this.repository.getChannelStatuses(integrationTypeId, isActive);
    }

    async createChannelStatus(data: CreateChannelStatusDTO) {
        return this.repository.createChannelStatus(data);
    }

    async updateChannelStatus(id: number, data: UpdateChannelStatusDTO) {
        return this.repository.updateChannelStatus(id, data);
    }

    async deleteChannelStatus(id: number) {
        return this.repository.deleteChannelStatus(id);
    }
}
