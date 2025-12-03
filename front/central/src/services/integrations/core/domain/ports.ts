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
    activateIntegration(id: number): Promise<ActionResponse>;
    deactivateIntegration(id: number): Promise<ActionResponse>;
    setAsDefault(id: number): Promise<ActionResponse>;

    // Integration Types
    getIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>>;
    getActiveIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>>;
    getIntegrationTypeById(id: number): Promise<SingleResponse<IntegrationType>>;
    getIntegrationTypeByCode(code: string): Promise<SingleResponse<IntegrationType>>;
    createIntegrationType(data: CreateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    updateIntegrationType(id: number, data: UpdateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>>;
    deleteIntegrationType(id: number): Promise<ActionResponse>;
}
