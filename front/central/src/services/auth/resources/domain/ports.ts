import {
    Resource,
    PaginatedResponse,
    GetResourcesParams,
    SingleResponse,
    CreateResourceDTO,
    UpdateResourceDTO,
    ActionResponse
} from './types';

export interface IResourceRepository {
    getResources(params?: GetResourcesParams): Promise<PaginatedResponse<Resource>>;
    getResourceById(id: number): Promise<SingleResponse<Resource>>;
    createResource(data: CreateResourceDTO): Promise<SingleResponse<Resource>>;
    updateResource(id: number, data: UpdateResourceDTO): Promise<SingleResponse<Resource>>;
    deleteResource(id: number): Promise<ActionResponse>;
}
