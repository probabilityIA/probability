import {
    Product,
    PaginatedResponse,
    GetProductsParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    ActionResponse,
    AddProductIntegrationDTO,
    ProductIntegrationsResponse
} from './types';

export interface IProductRepository {
    getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>>;
    getProductById(id: string, businessId?: number): Promise<SingleResponse<Product>>;
    createProduct(data: CreateProductDTO, businessId?: number): Promise<SingleResponse<Product>>;
    updateProduct(id: string, data: UpdateProductDTO, businessId?: number): Promise<SingleResponse<Product>>;
    deleteProduct(id: string, businessId?: number): Promise<ActionResponse>;

    // Product-Integration Management
    addProductIntegration(productId: string, data: AddProductIntegrationDTO, businessId?: number): Promise<SingleResponse<any>>;
    removeProductIntegration(productId: string, integrationId: number, businessId?: number): Promise<ActionResponse>;
    getProductIntegrations(productId: string, businessId?: number): Promise<ProductIntegrationsResponse>;
}
