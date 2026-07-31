import { IWebsiteConfigRepository } from '../domain/ports';
import { UpdateWebsiteConfigDTO } from '../domain/types';

export class WebsiteConfigUseCases {
    constructor(private repository: IWebsiteConfigRepository) {}

    async getConfig(businessId?: number) {
        return this.repository.getConfig(businessId);
    }

    async updateConfig(data: UpdateWebsiteConfigDTO, businessId?: number) {
        return this.repository.updateConfig(data, businessId);
    }

    async uploadImage(formData: FormData, businessId?: number) {
        return this.repository.uploadImage(formData, businessId);
    }

    async deleteImage(imageUrl: string, businessId?: number) {
        return this.repository.deleteImage(imageUrl, businessId);
    }
}
