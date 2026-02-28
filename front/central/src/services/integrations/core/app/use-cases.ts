import { IIntegrationRepository } from '../domain/ports';
import {
    Integration,
    PaginatedResponse,
    GetIntegrationsParams,
    SingleResponse,
    CreateIntegrationDTO,
    UpdateIntegrationDTO,
    ActionResponse,
    IntegrationType,
    CreateIntegrationTypeDTO,
    UpdateIntegrationTypeDTO,
    WebhookResponse,
    ListWebhooksResponse,
    DeleteWebhookResponse,
    VerifyWebhooksResponse,
    CreateWebhookResponse,
    SyncOrdersParams,
    IntegrationCategoriesResponse
} from '../domain/types';

export class IntegrationUseCases {
    constructor(private readonly repository: IIntegrationRepository) { }

    async getIntegrations(params?: GetIntegrationsParams): Promise<PaginatedResponse<Integration>> {
        return this.repository.getIntegrations(params);
    }

    async getIntegrationById(id: number): Promise<SingleResponse<Integration>> {
        return this.repository.getIntegrationById(id);
    }

    async getIntegrationByType(type: string, businessId?: number): Promise<SingleResponse<Integration>> {
        return this.repository.getIntegrationByType(type, businessId);
    }

    async createIntegration(data: CreateIntegrationDTO): Promise<SingleResponse<Integration>> {
        return this.repository.createIntegration(data);
    }

    async updateIntegration(id: number, data: UpdateIntegrationDTO): Promise<SingleResponse<Integration>> {
        return this.repository.updateIntegration(id, data);
    }

    async deleteIntegration(id: number): Promise<ActionResponse> {
        return this.repository.deleteIntegration(id);
    }

    async testConnection(id: number): Promise<ActionResponse> {
        return this.repository.testConnection(id);
    }

    async activateIntegration(id: number): Promise<ActionResponse> {
        return this.repository.activateIntegration(id);
    }

    async deactivateIntegration(id: number): Promise<ActionResponse> {
        return this.repository.deactivateIntegration(id);
    }

    async setAsDefault(id: number): Promise<ActionResponse> {
        return this.repository.setAsDefault(id);
    }

    async syncOrders(id: number, params?: SyncOrdersParams): Promise<ActionResponse> {
        return this.repository.syncOrders(id, params);
    }

    async getSyncStatus(id: number, businessId?: number): Promise<{ success: boolean; in_progress: boolean; sync_state?: any }> {
        return this.repository.getSyncStatus(id, businessId);
    }

    async testIntegration(id: number): Promise<ActionResponse> {
        return this.repository.testIntegration(id);
    }

    async testConnectionRaw(typeCode: string, config: any, credentials: any): Promise<ActionResponse> {
        return this.repository.testConnectionRaw(typeCode, config, credentials);
    }

    // Integration Types
    async getIntegrationTypes(categoryId?: number): Promise<SingleResponse<IntegrationType[]>> {
        return this.repository.getIntegrationTypes(categoryId);
    }

    async getActiveIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>> {
        return this.repository.getActiveIntegrationTypes();
    }

    async getIntegrationTypeById(id: number): Promise<SingleResponse<IntegrationType>> {
        return this.repository.getIntegrationTypeById(id);
    }

    async getIntegrationTypeByCode(code: string): Promise<SingleResponse<IntegrationType>> {
        return this.repository.getIntegrationTypeByCode(code);
    }

    async createIntegrationType(data: CreateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>> {
        return this.repository.createIntegrationType(data);
    }

    async updateIntegrationType(id: number, data: UpdateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>> {
        return this.repository.updateIntegrationType(id, data);
    }

    async deleteIntegrationType(id: number): Promise<ActionResponse> {
        return this.repository.deleteIntegrationType(id);
    }

    async getIntegrationTypePlatformCredentials(id: number): Promise<{ success: boolean; message: string; data: Record<string, string> }> {
        return this.repository.getIntegrationTypePlatformCredentials(id);
    }

    async getWebhookUrl(id: number): Promise<WebhookResponse> {
        return this.repository.getWebhookUrl(id);
    }

    async listWebhooks(id: number): Promise<ListWebhooksResponse> {
        return this.repository.listWebhooks(id);
    }

    async deleteWebhook(id: number, webhookId: string): Promise<DeleteWebhookResponse> {
        return this.repository.deleteWebhook(id, webhookId);
    }

    async verifyWebhooks(id: number): Promise<VerifyWebhooksResponse> {
        return this.repository.verifyWebhooks(id);
    }

    async createWebhook(id: number): Promise<CreateWebhookResponse> {
        return this.repository.createWebhook(id);
    }

    // Integration Categories
    async getIntegrationCategories(): Promise<IntegrationCategoriesResponse> {
        return this.repository.getIntegrationCategories();
    }
}
