import { IBusinessRepository } from '../domain/ports';
import {
    Business,
    PaginatedResponse,
    GetBusinessesParams,
    SingleResponse,
    CreateBusinessDTO,
    UpdateBusinessDTO,
    ActionResponse,
    BusinessConfiguredResources,
    GetConfiguredResourcesParams,
    BusinessType,
    CreateBusinessTypeDTO,
    UpdateBusinessTypeDTO
} from '../domain/types';

export class BusinessUseCases {
    constructor(private repository: IBusinessRepository) { }

    // Business
    async getBusinesses(params?: GetBusinessesParams): Promise<PaginatedResponse<Business>> {
        return this.repository.getBusinesses(params);
    }

    async getBusinessById(id: number): Promise<SingleResponse<Business>> {
        return this.repository.getBusinessById(id);
    }

    async createBusiness(data: CreateBusinessDTO): Promise<SingleResponse<Business>> {
        return this.repository.createBusiness(data);
    }

    async updateBusiness(id: number, data: UpdateBusinessDTO): Promise<SingleResponse<Business>> {
        return this.repository.updateBusiness(id, data);
    }

    async deleteBusiness(id: number): Promise<ActionResponse> {
        return this.repository.deleteBusiness(id);
    }

    // Configured Resources
    async getConfiguredResources(params?: GetConfiguredResourcesParams): Promise<PaginatedResponse<BusinessConfiguredResources>> {
        return this.repository.getConfiguredResources(params);
    }

    async getBusinessConfiguredResources(id: number): Promise<SingleResponse<BusinessConfiguredResources>> {
        return this.repository.getBusinessConfiguredResources(id);
    }

    async activateResource(resourceId: number, businessId?: number): Promise<ActionResponse> {
        return this.repository.activateResource(resourceId, businessId);
    }

    async deactivateResource(resourceId: number, businessId?: number): Promise<ActionResponse> {
        return this.repository.deactivateResource(resourceId, businessId);
    }

    // Business Types
    async getBusinessTypes(): Promise<PaginatedResponse<BusinessType>> {
        return this.repository.getBusinessTypes();
    }

    async getBusinessTypeById(id: number): Promise<SingleResponse<BusinessType>> {
        return this.repository.getBusinessTypeById(id);
    }

    async createBusinessType(data: CreateBusinessTypeDTO): Promise<SingleResponse<BusinessType>> {
        return this.repository.createBusinessType(data);
    }

    async updateBusinessType(id: number, data: UpdateBusinessTypeDTO): Promise<SingleResponse<BusinessType>> {
        return this.repository.updateBusinessType(id, data);
    }

    async deleteBusinessType(id: number): Promise<ActionResponse> {
        return this.repository.deleteBusinessType(id);
    }
}
