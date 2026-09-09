import {
    StorefrontProduct,
    StorefrontOrder,
    StorefrontClient,
    CatalogFiltersResult,
    CatalogLayout,
    CatalogBanner,
    CreateStorefrontOrderDTO,
    CreateClientDTO,
    CreateClientResult,
    PaginatedResponse,
} from './types';

export interface IStorefrontRepository {
    getCatalog(params?: { page?: number; page_size?: number; search?: string; category?: string; family_id?: number; business_id?: number }): Promise<PaginatedResponse<StorefrontProduct>>;
    getCatalogFilters(businessId?: number): Promise<CatalogFiltersResult>;
    updateCatalogLayout(layout: CatalogLayout, businessId?: number): Promise<CatalogLayout>;
    updateCatalogBanner(enabled: boolean, businessId?: number): Promise<CatalogBanner>;
    uploadCatalogBannerImage(formData: FormData, businessId?: number): Promise<CatalogBanner>;
    getProduct(id: string, businessId?: number): Promise<StorefrontProduct>;
    createOrder(data: CreateStorefrontOrderDTO, businessId?: number): Promise<{ message: string }>;
    getOrders(params?: { page?: number; page_size?: number; business_id?: number }): Promise<PaginatedResponse<StorefrontOrder>>;
    getOrder(id: string, businessId?: number): Promise<StorefrontOrder>;
    createClient(data: CreateClientDTO, businessId?: number): Promise<CreateClientResult>;
    getClients(params?: { page?: number; page_size?: number; business_id?: number }): Promise<PaginatedResponse<StorefrontClient>>;
}
