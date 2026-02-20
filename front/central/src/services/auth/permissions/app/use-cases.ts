import { IPermissionRepository } from '../domain/ports';
import {
    GetPermissionsParams,
    CreatePermissionDTO,
    UpdatePermissionDTO
} from '../domain/types';


export class PermissionUseCases {
    constructor(private repository: IPermissionRepository) { }

    async getPermissions(params?: GetPermissionsParams) {
        return this.repository.getPermissions(params);
    }

    async getPermissionById(id: number) {
        return this.repository.getPermissionById(id);
    }

    async getPermissionsByScope(scopeId: number) {
        return this.repository.getPermissionsByScope(scopeId);
    }

    async getPermissionsByResource(resource: string) {
        return this.repository.getPermissionsByResource(resource);
    }

    async createPermission(data: CreatePermissionDTO) {
        return this.repository.createPermission(data);
    }

    async updatePermission(id: number, data: UpdatePermissionDTO) {
        return this.repository.updatePermission(id, data);
    }

    async deletePermission(id: number) {
        return this.repository.deletePermission(id);
    }

    async createPermissionsBulk(permissions: CreatePermissionDTO[]) {
        return this.repository.createPermissionsBulk(permissions);
    }
}
