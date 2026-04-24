import { IProductRepository } from '../domain/ports';
import {
    GetProductsParams,
    GetFamiliesParams,
    CreateProductDTO,
    UpdateProductDTO,
    CreateProductFamilyDTO,
    UpdateProductFamilyDTO,
    AddProductIntegrationDTO,
    UpdateProductIntegrationDTO
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

    async uploadProductImage(productId: string, formData: FormData, businessId?: number) {
        return this.repository.uploadProductImage(productId, formData, businessId);
    }

    async getProductFamilies(params?: GetFamiliesParams) {
        return this.repository.getProductFamilies(params);
    }

    async getProductFamilyById(familyId: number, businessId?: number) {
        return this.repository.getProductFamilyById(familyId, businessId);
    }

    async getFamilyVariants(familyId: number, businessId?: number) {
        return this.repository.getFamilyVariants(familyId, businessId);
    }

    async createProductFamily(data: CreateProductFamilyDTO, businessId?: number) {
        return this.repository.createProductFamily(data, businessId);
    }

    async updateProductFamily(familyId: number, data: UpdateProductFamilyDTO, businessId?: number) {
        return this.repository.updateProductFamily(familyId, data, businessId);
    }

    async deleteProductFamily(familyId: number, businessId?: number) {
        return this.repository.deleteProductFamily(familyId, businessId);
    }

    async addProductIntegration(productId: string, data: AddProductIntegrationDTO, businessId?: number) {
        return this.repository.addProductIntegration(productId, data, businessId);
    }

    async updateProductIntegration(productId: string, integrationId: number, data: UpdateProductIntegrationDTO, businessId?: number) {
        return this.repository.updateProductIntegration(productId, integrationId, data, businessId);
    }

    async removeProductIntegration(productId: string, integrationId: number, businessId?: number) {
        return this.repository.removeProductIntegration(productId, integrationId, businessId);
    }

    async getProductIntegrations(productId: string, businessId?: number) {
        return this.repository.getProductIntegrations(productId, businessId);
    }

    async lookupProductByExternalRef(
        integrationId: number,
        params: { external_variant_id?: string; external_sku?: string; external_product_id?: string; external_barcode?: string },
        businessId?: number
    ) {
        return this.repository.lookupProductByExternalRef(integrationId, params, businessId);
    }
}
