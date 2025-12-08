import { IProductRepository } from '../domain/ports';
import {
    GetProductsParams,
    CreateProductDTO,
    UpdateProductDTO
} from '../domain/types';

export class ProductUseCases {
    constructor(private repository: IProductRepository) { }

    async getProducts(params?: GetProductsParams) {
        return this.repository.getProducts(params);
    }

    async getProductById(id: string) {
        return this.repository.getProductById(id);
    }

    async createProduct(data: CreateProductDTO) {
        return this.repository.createProduct(data);
    }

    async updateProduct(id: string, data: UpdateProductDTO) {
        return this.repository.updateProduct(id, data);
    }

    async deleteProduct(id: string) {
        return this.repository.deleteProduct(id);
    }
}
