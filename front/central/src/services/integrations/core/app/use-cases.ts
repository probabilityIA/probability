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
    UpdateIntegrationTypeDTO
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

    async syncOrders(id: number): Promise<ActionResponse> {
        return this.repository.syncOrders(id);
    }

    async testIntegration(id: number): Promise<ActionResponse> {
        return this.repository.testIntegration(id);
    }

    async testConnectionRaw(typeCode: string, config: any, credentials: any): Promise<ActionResponse> {
        return this.repository.testConnectionRaw(typeCode, config, credentials);
    }

    // Integration Types
    async getIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>> {
        return this.repository.getIntegrationTypes();
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
}
