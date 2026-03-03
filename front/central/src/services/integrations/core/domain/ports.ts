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
} from './types';

export interface IIntegrationRepository {
    // Integrations
    getIntegrations(params?: GetIntegrationsParams): Promise<PaginatedResponse<Integration>>;
    getIntegrationById(id: number): Promise<SingleResponse<Integration>>;
    getIntegrationByType(type: string, businessId?: number): Promise<SingleResponse<Integration>>;
    createIntegration(data: CreateIntegrationDTO): Promise<SingleResponse<Integration>>;
    updateIntegration(id: number, data: UpdateIntegrationDTO): Promise<SingleResponse<Integration>>;
    deleteIntegration(id: number): Promise<ActionResponse>;
    testConnection(id: number): Promise<ActionResponse>;
    activateIntegration(id: number): Promise<ActionResponse>;
    deactivateIntegration(id: number): Promise<ActionResponse>;
    setAsDefault(id: number): Promise<SingleResponse<Integration>>;
    syncOrders(id: number, params?: SyncOrdersParams): Promise<ActionResponse>;
    getSyncStatus(id: number, businessId?: number): Promise<{ success: boolean; in_progress: boolean; sync_state?: any }>;
    testIntegration(id: number): Promise<ActionResponse>;
    testConnectionRaw(typeCode: string, config: any, credentials: any): Promise<ActionResponse>;
    getWebhookUrl(id: number): Promise<WebhookResponse>;
    listWebhooks(id: number): Promise<ListWebhooksResponse>;
    deleteWebhook(id: number, webhookId: string): Promise<DeleteWebhookResponse>;
    verifyWebhooks(id: number): Promise<VerifyWebhooksResponse>;
    createWebhook(id: number): Promise<CreateWebhookResponse>;

    // Integration Types
    getIntegrationTypes(categoryId?: number): Promise<SingleResponse<IntegrationType[]>>;
    getActiveIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>>;
    getIntegrationTypeById(id: number): Promise<SingleResponse<IntegrationType>>;
    getIntegrationTypeByCode(code: string): Promise<SingleResponse<IntegrationType>>;
    createIntegrationType(data: CreateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    updateIntegrationType(id: number, data: UpdateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    deleteIntegrationType(id: number): Promise<ActionResponse>;
    getIntegrationTypePlatformCredentials(id: number): Promise<{ success: boolean; message: string; data: Record<string, string> }>;

    // Integration Categories
    getIntegrationCategories(): Promise<IntegrationCategoriesResponse>;
}
