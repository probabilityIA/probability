import { IStorefrontRepository } from '../domain/ports';
import { CreateStorefrontOrderDTO, RegisterDTO } from '../domain/types';

export class StorefrontUseCases {
    constructor(private repository: IStorefrontRepository) {}

    async getCatalog(params?: { page?: number; page_size?: number; search?: string; category?: string; business_id?: number }) {
        return this.repository.getCatalog(params);
    }

    async getProduct(id: string, businessId?: number) {
        return this.repository.getProduct(id, businessId);
    }

    async createOrder(data: CreateStorefrontOrderDTO, businessId?: number) {
        return this.repository.createOrder(data, businessId);
    }

    async getOrders(params?: { page?: number; page_size?: number; business_id?: number }) {
        return this.repository.getOrders(params);
    }

    async getOrder(id: string, businessId?: number) {
        return this.repository.getOrder(id, businessId);
    }

    async register(data: RegisterDTO) {
        return this.repository.register(data);
    }
}
