import {
    Permission,
    PaginatedResponse,
    GetPermissionsParams,
    SingleResponse,
    CreatePermissionDTO,
    UpdatePermissionDTO,
    ActionResponse,
    BulkCreatePermissionsResponse
} from './types';

export interface IPermissionRepository {
    getPermissions(params?: GetPermissionsParams): Promise<PaginatedResponse<Permission>>;
    getPermissionById(id: number): Promise<SingleResponse<Permission>>;
    getPermissionsByScope(scopeId: number): Promise<PaginatedResponse<Permission>>;
    getPermissionsByResource(resource: string): Promise<PaginatedResponse<Permission>>;
    createPermission(data: CreatePermissionDTO): Promise<SingleResponse<Permission>>;
    updatePermission(id: number, data: UpdatePermissionDTO): Promise<SingleResponse<Permission>>;
    deletePermission(id: number): Promise<ActionResponse>;
    createPermissionsBulk(permissions: CreatePermissionDTO[]): Promise<BulkCreatePermissionsResponse>;
}
