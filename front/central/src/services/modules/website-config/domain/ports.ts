import { WebsiteConfigData, UpdateWebsiteConfigDTO } from './types';

export interface IWebsiteConfigRepository {
    getConfig(businessId?: number): Promise<WebsiteConfigData>;
    updateConfig(data: UpdateWebsiteConfigDTO, businessId?: number): Promise<WebsiteConfigData>;
}
