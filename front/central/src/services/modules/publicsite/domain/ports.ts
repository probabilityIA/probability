import { PublicBusiness, PublicProduct, PaginatedResponse, ContactFormDTO } from './types';

export interface IPublicSiteRepository {
    getBusinessPage(slug: string): Promise<PublicBusiness>;
    getCatalog(slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }): Promise<PaginatedResponse<PublicProduct>>;
    getProduct(slug: string, productId: string): Promise<PublicProduct>;
    submitContact(slug: string, data: ContactFormDTO): Promise<{ message: string }>;
}
