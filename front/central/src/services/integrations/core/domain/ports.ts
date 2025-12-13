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
    UpdateIntegrationTypeDTO
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
    activateIntegration(id: number): Promise<SingleResponse<Integration>>;
    deactivateIntegration(id: number): Promise<SingleResponse<Integration>>;
    setAsDefault(id: number): Promise<SingleResponse<Integration>>;
    syncOrders(id: number): Promise<ActionResponse>;
    testIntegration(id: number): Promise<ActionResponse>;
    testConnectionRaw(typeCode: string, config: any, credentials: any): Promise<ActionResponse>;

    // Integration Types
    getIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>>;
    getActiveIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>>;
    getIntegrationTypeById(id: number): Promise<SingleResponse<IntegrationType>>;
    getIntegrationTypeByCode(code: string): Promise<SingleResponse<IntegrationType>>;
    createIntegrationType(data: CreateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    updateIntegrationType(id: number, data: UpdateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    deleteIntegrationType(id: number): Promise<ActionResponse>;
}
