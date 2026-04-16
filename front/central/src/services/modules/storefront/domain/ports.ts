import {
    StorefrontProduct,
    StorefrontOrder,
    CreateStorefrontOrderDTO,
    RegisterDTO,
    PaginatedResponse,
} from './types';

export interface IStorefrontRepository {
    getCatalog(params?: { page?: number; page_size?: number; search?: string; category?: string; business_id?: number }): Promise<PaginatedResponse<StorefrontProduct>>;
    getProduct(id: string, businessId?: number): Promise<StorefrontProduct>;
    createOrder(data: CreateStorefrontOrderDTO, businessId?: number): Promise<{ message: string }>;
    getOrders(params?: { page?: number; page_size?: number; business_id?: number }): Promise<PaginatedResponse<StorefrontOrder>>;
    getOrder(id: string, businessId?: number): Promise<StorefrontOrder>;
    register(data: RegisterDTO): Promise<{ message: string }>;
}
