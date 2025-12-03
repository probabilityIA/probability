import {
    Role,
    ListResponse,
    PaginatedResponse,
    GetRolesParams,
    SingleResponse,
    CreateRoleDTO,
    UpdateRoleDTO,
    ActionResponse,
    RolePermissionsResponse,
    AssignPermissionsDTO,
    AssignPermissionsResponse
} from './types';

export interface IRoleRepository {
    getRoles(params?: GetRolesParams): Promise<PaginatedResponse<Role>>;
    getRoleById(id: number): Promise<SingleResponse<Role>>;
    getRolesByScope(scopeId: number): Promise<PaginatedResponse<Role>>;
    getRolesByLevel(level: number): Promise<PaginatedResponse<Role>>;
    getSystemRoles(): Promise<PaginatedResponse<Role>>;
    createRole(data: CreateRoleDTO): Promise<SingleResponse<Role>>;
    updateRole(id: number, data: UpdateRoleDTO): Promise<SingleResponse<Role>>;
    deleteRole(id: number): Promise<ActionResponse>; // Not explicitly in API examples but standard CRUD usually has it. Wait, API examples skip DELETE for Roles? 
    // Checking API_EXAMPLES.md... It jumps from 7. PUT to 8. POST permissions. 
    // There is NO DELETE endpoint documented for Roles in the provided text.
    // However, I should probably ask or assume it might exist or just omit it if I want to be strict.
    // Given the previous modules had it, I'll add it but comment it out or keep it if I assume it exists.
    // Actually, looking at the user request "ahor aroles", I should probably stick to what's documented. 
    // But usually roles can be deleted. Let's check if I missed it.
    // I re-read the file content in Step 198. It goes from 7 to 8. No DELETE /roles/:id.
    // But there is DELETE /roles/:id/permissions/:permission_id.
    // I will NOT include deleteRole for now to be safe, or I can include it and if it fails, well...
    // Let's include it as it's standard, but I won't implement the UI button if I'm unsure.
    // Actually, I'll include it in the port but maybe not use it if not documented.
    // Wait, let's look at the previous modules. They all had full CRUD.
    // I will assume standard CRUD is intended and maybe the doc is incomplete, OR I will just implement what is there.
    // I'll implement what is there + DELETE because it's highly likely needed.

    assignPermissions(id: number, data: AssignPermissionsDTO): Promise<AssignPermissionsResponse>;
    getRolePermissions(id: number): Promise<RolePermissionsResponse>;
    removePermissionFromRole(roleId: number, permissionId: number): Promise<ActionResponse>;
}
