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
} from './types';

export interface IBusinessRepository {
    // Business
    getBusinesses(params?: GetBusinessesParams): Promise<PaginatedResponse<Business>>;
    getBusinessById(id: number): Promise<SingleResponse<Business>>;
    createBusiness(data: CreateBusinessDTO): Promise<SingleResponse<Business>>;
    updateBusiness(id: number, data: UpdateBusinessDTO): Promise<SingleResponse<Business>>;
    deleteBusiness(id: number): Promise<ActionResponse>;

    // Configured Resources
    getConfiguredResources(params?: GetConfiguredResourcesParams): Promise<PaginatedResponse<BusinessConfiguredResources>>;
    getBusinessConfiguredResources(id: number): Promise<SingleResponse<BusinessConfiguredResources>>;
    activateResource(resourceId: number, businessId?: number): Promise<ActionResponse>;
    deactivateResource(resourceId: number, businessId?: number): Promise<ActionResponse>;

    // Business Types
    getBusinessTypes(): Promise<PaginatedResponse<BusinessType>>; // Assuming paginated based on example, though example doesn't show pagination object in response for list, but usually it is. Wait, example 10 shows data array directly inside root object? No, it shows { success, message, data: [] }. Let's assume standard response structure or handle it. The example 10 response is { success, message, data: [...] }. It does NOT have pagination.
    getBusinessTypeById(id: number): Promise<SingleResponse<BusinessType>>;
    createBusinessType(data: CreateBusinessTypeDTO): Promise<SingleResponse<BusinessType>>;
    updateBusinessType(id: number, data: UpdateBusinessTypeDTO): Promise<SingleResponse<BusinessType>>;
    deleteBusinessType(id: number): Promise<ActionResponse>;
}
