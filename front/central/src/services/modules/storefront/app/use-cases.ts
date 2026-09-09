import { IStorefrontRepository } from '../domain/ports';
import { CreateStorefrontOrderDTO, CreateClientDTO, CatalogLayout } from '../domain/types';

export class StorefrontUseCases {
    constructor(private repository: IStorefrontRepository) {}

    async getCatalog(params?: { page?: number; page_size?: number; search?: string; category?: string; family_id?: number; business_id?: number }) {
        return this.repository.getCatalog(params);
    }

    async getCatalogFilters(businessId?: number) {
        return this.repository.getCatalogFilters(businessId);
    }

    async updateCatalogLayout(layout: CatalogLayout, businessId?: number) {
        return this.repository.updateCatalogLayout(layout, businessId);
    }

    async updateCatalogBanner(enabled: boolean, businessId?: number) {
        return this.repository.updateCatalogBanner(enabled, businessId);
    }

    async uploadCatalogBannerImage(formData: FormData, businessId?: number) {
        return this.repository.uploadCatalogBannerImage(formData, businessId);
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

    async createClient(data: CreateClientDTO, businessId?: number) {
        return this.repository.createClient(data, businessId);
    }

    async getClients(params?: { page?: number; page_size?: number; business_id?: number }) {
        return this.repository.getClients(params);
    }
}
