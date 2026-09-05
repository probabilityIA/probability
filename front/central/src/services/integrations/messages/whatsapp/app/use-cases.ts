import { IWhatsAppRepository } from '../domain/ports';
import { WhatsAppProvisionResponse, WhatsAppTemplatesResponse } from '../domain/types';

export class WhatsAppUseCases {
    constructor(private readonly repository: IWhatsAppRepository) {}

    async getTemplatesStatus(businessId?: number, refresh = false): Promise<WhatsAppTemplatesResponse> {
        return this.repository.getTemplatesStatus(businessId, refresh);
    }

    async provisionTemplates(businessId?: number): Promise<WhatsAppProvisionResponse> {
        return this.repository.provisionTemplates(businessId);
    }
}
