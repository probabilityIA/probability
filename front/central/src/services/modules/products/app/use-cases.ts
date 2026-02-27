import { IProductRepository } from '../domain/ports';
import {
    GetProductsParams,
    CreateProductDTO,
    UpdateProductDTO,
    AddProductIntegrationDTO
} from '../domain/types';

export class ProductUseCases {
    constructor(private repository: IProductRepository) { }

    async getProducts(params?: GetProductsParams) {
        return this.repository.getProducts(params);
    }

    async getProductById(id: string, businessId?: number) {
        return this.repository.getProductById(id, businessId);
    }

    async createProduct(data: CreateProductDTO, businessId?: number) {
        return this.repository.createProduct(data, businessId);
    }

    async updateProduct(id: string, data: UpdateProductDTO, businessId?: number) {
        return this.repository.updateProduct(id, data, businessId);
    }

    async deleteProduct(id: string, businessId?: number) {
        return this.repository.deleteProduct(id, businessId);
    }

    // ═══════════════════════════════════════════
    // Product-Integration Management
    // ═══════════════════════════════════════════

    async addProductIntegration(productId: string, data: AddProductIntegrationDTO, businessId?: number) {
        return this.repository.addProductIntegration(productId, data, businessId);
    }

    async removeProductIntegration(productId: string, integrationId: number, businessId?: number) {
        return this.repository.removeProductIntegration(productId, integrationId, businessId);
    }

    async getProductIntegrations(productId: string, businessId?: number) {
        return this.repository.getProductIntegrations(productId, businessId);
    }
}
