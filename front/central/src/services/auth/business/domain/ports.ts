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
    UpdateBusinessTypeDTO,
    BusinessFiscalProfile,
    UpdateFiscalProfileDTO
} from './types';

export interface IBusinessRepository {
    getBusinesses(params?: GetBusinessesParams): Promise<PaginatedResponse<Business>>;
    getBusinessById(id: number): Promise<SingleResponse<Business>>;
    createBusiness(data: CreateBusinessDTO): Promise<SingleResponse<Business>>;
    updateBusiness(id: number, data: UpdateBusinessDTO): Promise<SingleResponse<Business>>;
    deleteBusiness(id: number): Promise<ActionResponse>;
    activateBusiness(id: number): Promise<ActionResponse>;
    deactivateBusiness(id: number): Promise<ActionResponse>;

    getConfiguredResources(params?: GetConfiguredResourcesParams): Promise<PaginatedResponse<BusinessConfiguredResources>>;
    getBusinessConfiguredResources(id: number): Promise<SingleResponse<BusinessConfiguredResources>>;
    activateResource(resourceId: number, businessId?: number): Promise<ActionResponse>;
    deactivateResource(resourceId: number, businessId?: number): Promise<ActionResponse>;

    getFiscalProfile(businessId: number): Promise<SingleResponse<BusinessFiscalProfile>>;
    updateFiscalProfile(businessId: number, data: UpdateFiscalProfileDTO): Promise<SingleResponse<BusinessFiscalProfile>>;

    getBusinessTypes(): Promise<PaginatedResponse<BusinessType>>;
    getBusinessTypeById(id: number): Promise<SingleResponse<BusinessType>>;
    createBusinessType(data: CreateBusinessTypeDTO): Promise<SingleResponse<BusinessType>>;
    updateBusinessType(id: number, data: UpdateBusinessTypeDTO): Promise<SingleResponse<BusinessType>>;
    deleteBusinessType(id: number): Promise<ActionResponse>;
}
