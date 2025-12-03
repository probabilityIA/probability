import {
    Integration,
    PaginatedResponse,
    GetIntegrationsParams,
    SingleResponse,
    CreateIntegrationDTO,
    UpdateIntegrationDTO,
    ActionResponse
} from './types';

export interface IIntegrationRepository {
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
}
