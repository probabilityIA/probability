import { WebsiteConfigData, UpdateWebsiteConfigDTO } from './types';

export interface IWebsiteConfigRepository {
    getConfig(businessId?: number): Promise<WebsiteConfigData>;
    updateConfig(data: UpdateWebsiteConfigDTO, businessId?: number): Promise<WebsiteConfigData>;
    uploadImage(formData: FormData, businessId?: number): Promise<{ success: boolean; image_url: string }>;
    deleteImage(imageUrl: string, businessId?: number): Promise<{ success: boolean }>;
}
