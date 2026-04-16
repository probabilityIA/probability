import { IRoleRepository } from '../domain/ports';
import {
    GetRolesParams,
    CreateRoleDTO,
    UpdateRoleDTO,
    AssignPermissionsDTO
} from '../domain/types';

export class RoleUseCases {
    constructor(private repository: IRoleRepository) { }

    async getRoles(params?: GetRolesParams) {
        return this.repository.getRoles(params);
    }

    async getRoleById(id: number) {
        return this.repository.getRoleById(id);
    }

    async getRolesByScope(scopeId: number) {
        return this.repository.getRolesByScope(scopeId);
    }

    async getRolesByLevel(level: number) {
        return this.repository.getRolesByLevel(level);
    }

    async getSystemRoles() {
        return this.repository.getSystemRoles();
    }

    async createRole(data: CreateRoleDTO) {
        return this.repository.createRole(data);
    }

    async updateRole(id: number, data: UpdateRoleDTO) {
        return this.repository.updateRole(id, data);
    }

    async deleteRole(id: number) {
        return this.repository.deleteRole(id);
    }

    async assignPermissions(id: number, data: AssignPermissionsDTO) {
        return this.repository.assignPermissions(id, data);
    }

    async getRolePermissions(id: number) {
        return this.repository.getRolePermissions(id);
    }

    async removePermissionFromRole(roleId: number, permissionId: number) {
        return this.repository.removePermissionFromRole(roleId, permissionId);
    }
}
