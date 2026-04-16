import { IPublicSiteRepository } from '../domain/ports';
import { ContactFormDTO } from '../domain/types';

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

    async submitContact(slug: string, data: ContactFormDTO) {
        return this.repository.submitContact(slug, data);
    }
}
