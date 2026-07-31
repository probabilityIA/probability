import { PublicBusiness, PublicProduct, PaginatedResponse, ContactFormDTO, CreateCheckoutInput, CheckoutSession, TiendaSession, TiendaRegisterDTO, TiendaOrderSummary, TiendaAddress, TiendaAddressInput, TiendaProfileUpdate } from './types';

export interface IPublicSiteRepository {
    getBusinessPage(slug: string): Promise<PublicBusiness>;
    getCatalog(slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }): Promise<PaginatedResponse<PublicProduct>>;
    getProduct(slug: string, productId: string): Promise<PublicProduct>;
    getCategories(slug: string): Promise<string[]>;
    submitContact(slug: string, data: ContactFormDTO): Promise<{ message: string }>;
    createCheckoutSession(slug: string, data: CreateCheckoutInput, token?: string): Promise<CheckoutSession>;
    getCheckoutStatus(slug: string, reference: string): Promise<{ status: string }>;
    registerCustomer(slug: string, data: TiendaRegisterDTO): Promise<{ success: boolean; message?: string }>;
    getSession(slug: string, token: string): Promise<TiendaSession>;
    getMyPrices(slug: string, token: string): Promise<Record<string, number>>;
    getMyOrders(slug: string, token: string, page?: number): Promise<{ data: TiendaOrderSummary[]; total: number }>;
    updateProfile(slug: string, token: string, data: TiendaProfileUpdate): Promise<{ success: boolean }>;
    uploadAvatar(slug: string, token: string, formData: FormData): Promise<{ success: boolean; avatar_url: string }>;
    listAddresses(slug: string, token: string): Promise<TiendaAddress[]>;
    addAddress(slug: string, token: string, data: TiendaAddressInput): Promise<{ success: boolean }>;
    setPrimaryAddress(slug: string, token: string, id: number): Promise<{ success: boolean }>;
    deleteAddress(slug: string, token: string, id: number): Promise<{ success: boolean }>;
}
