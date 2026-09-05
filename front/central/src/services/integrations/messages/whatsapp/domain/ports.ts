import { WhatsAppProvisionResponse, WhatsAppTemplatesResponse } from './types';

export interface IWhatsAppRepository {
    getTemplatesStatus(businessId?: number, refresh?: boolean): Promise<WhatsAppTemplatesResponse>;
    provisionTemplates(businessId?: number): Promise<WhatsAppProvisionResponse>;
}
