import {
    Product,
    PaginatedResponse,
    GetProductsParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    ActionResponse
} from './types';

export interface IProductRepository {
    getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>>;
    getProductById(id: string): Promise<SingleResponse<Product>>;
    createProduct(data: CreateProductDTO): Promise<SingleResponse<Product>>;
    updateProduct(id: string, data: UpdateProductDTO): Promise<SingleResponse<Product>>;
    deleteProduct(id: string): Promise<ActionResponse>;
}
