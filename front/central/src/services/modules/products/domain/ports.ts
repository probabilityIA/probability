import {
    Product,
    ProductFamily,
    PaginatedResponse,
    GetProductsParams,
    GetFamiliesParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    CreateProductFamilyDTO,
    UpdateProductFamilyDTO,
    ActionResponse,
    AddProductIntegrationDTO,
    UpdateProductIntegrationDTO,
    ProductIntegrationsResponse,
    UploadImageResponse
} from './types';

export interface IProductRepository {
    getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>>;
    getProductById(id: string, businessId?: number): Promise<SingleResponse<Product>>;
    createProduct(data: CreateProductDTO, businessId?: number): Promise<SingleResponse<Product>>;
    updateProduct(id: string, data: UpdateProductDTO, businessId?: number): Promise<SingleResponse<Product>>;
    deleteProduct(id: string, businessId?: number): Promise<ActionResponse>;
    uploadProductImage(productId: string, formData: FormData, businessId?: number): Promise<UploadImageResponse>;

    getProductFamilies(params?: GetFamiliesParams): Promise<PaginatedResponse<ProductFamily>>;
    getProductFamilyById(familyId: number, businessId?: number): Promise<SingleResponse<ProductFamily>>;
    getFamilyVariants(familyId: number, businessId?: number): Promise<{ success: boolean; data: Product[] }>;
    createProductFamily(data: CreateProductFamilyDTO, businessId?: number): Promise<SingleResponse<ProductFamily>>;
    updateProductFamily(familyId: number, data: UpdateProductFamilyDTO, businessId?: number): Promise<SingleResponse<ProductFamily>>;
    deleteProductFamily(familyId: number, businessId?: number): Promise<ActionResponse>;

    addProductIntegration(productId: string, data: AddProductIntegrationDTO, businessId?: number): Promise<SingleResponse<any>>;
    updateProductIntegration(productId: string, integrationId: number, data: UpdateProductIntegrationDTO, businessId?: number): Promise<SingleResponse<any>>;
    removeProductIntegration(productId: string, integrationId: number, businessId?: number): Promise<ActionResponse>;
    getProductIntegrations(productId: string, businessId?: number): Promise<ProductIntegrationsResponse>;
    lookupProductByExternalRef(integrationId: number, params: { external_variant_id?: string; external_sku?: string; external_product_id?: string; external_barcode?: string }, businessId?: number): Promise<SingleResponse<Product>>;
}
