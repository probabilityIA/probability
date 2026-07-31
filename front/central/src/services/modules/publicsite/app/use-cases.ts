import { IPublicSiteRepository } from '../domain/ports';
import { ContactFormDTO, CreateCheckoutInput, TiendaRegisterDTO, TiendaAddressInput, TiendaProfileUpdate } from '../domain/types';

export class PublicSiteUseCases {
    constructor(private repository: IPublicSiteRepository) {}

    async getBusinessPage(slug: string) {
        return this.repository.getBusinessPage(slug);
    }

    async getCatalog(slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }) {
        return this.repository.getCatalog(slug, params);
    }

    async getProduct(slug: string, productId: string) {
        return this.repository.getProduct(slug, productId);
    }

    async getCategories(slug: string) {
        return this.repository.getCategories(slug);
    }

    async submitContact(slug: string, data: ContactFormDTO) {
        return this.repository.submitContact(slug, data);
    }

    async createCheckoutSession(slug: string, data: CreateCheckoutInput, token?: string) {
        return this.repository.createCheckoutSession(slug, data, token);
    }

    async getCheckoutStatus(slug: string, reference: string) {
        return this.repository.getCheckoutStatus(slug, reference);
    }

    async registerCustomer(slug: string, data: TiendaRegisterDTO) {
        return this.repository.registerCustomer(slug, data);
    }

    async getSession(slug: string, token: string) {
        return this.repository.getSession(slug, token);
    }

    async getMyPrices(slug: string, token: string) {
        return this.repository.getMyPrices(slug, token);
    }

    async getMyOrders(slug: string, token: string, page?: number) {
        return this.repository.getMyOrders(slug, token, page);
    }

    async updateProfile(slug: string, token: string, data: TiendaProfileUpdate) {
        return this.repository.updateProfile(slug, token, data);
    }

    async uploadAvatar(slug: string, token: string, formData: FormData) {
        return this.repository.uploadAvatar(slug, token, formData);
    }

    async listAddresses(slug: string, token: string) {
        return this.repository.listAddresses(slug, token);
    }

    async addAddress(slug: string, token: string, data: TiendaAddressInput) {
        return this.repository.addAddress(slug, token, data);
    }

    async setPrimaryAddress(slug: string, token: string, id: number) {
        return this.repository.setPrimaryAddress(slug, token, id);
    }

    async deleteAddress(slug: string, token: string, id: number) {
        return this.repository.deleteAddress(slug, token, id);
    }
}
