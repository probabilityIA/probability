'use server';

import { cookies } from 'next/headers';
import { RoleApiRepository } from '../repository/api-repository';
import { RoleUseCases } from '../../app/use-cases';
import {
    GetRolesParams,
    CreateRoleDTO,
    UpdateRoleDTO,
    AssignPermissionsDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new RoleApiRepository(token);
    return new RoleUseCases(repository);
}

export const getRolesAction = async (params?: GetRolesParams) => {
    try {
        return await (await getUseCases()).getRoles(params);
    } catch (error: any) {
        console.error('Get Roles Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getRoleByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getRoleById(id);
    } catch (error: any) {
        console.error('Get Role By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getRolesByScopeAction = async (scopeId: number) => {
    try {
        return await (await getUseCases()).getRolesByScope(scopeId);
    } catch (error: any) {
        console.error('Get Roles By Scope Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getRolesByLevelAction = async (level: number) => {
    try {
        return await (await getUseCases()).getRolesByLevel(level);
    } catch (error: any) {
        console.error('Get Roles By Level Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getSystemRolesAction = async () => {
    try {
        return await (await getUseCases()).getSystemRoles();
    } catch (error: any) {
        console.error('Get System Roles Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createRoleAction = async (data: CreateRoleDTO) => {
    try {
        return await (await getUseCases()).createRole(data);
    } catch (error: any) {
        console.error('Create Role Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateRoleAction = async (id: number, data: UpdateRoleDTO) => {
    try {
        return await (await getUseCases()).updateRole(id, data);
    } catch (error: any) {
        console.error('Update Role Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteRoleAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteRole(id);
    } catch (error: any) {
        console.error('Delete Role Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const assignPermissionsAction = async (id: number, data: AssignPermissionsDTO) => {
    try {
        return await (await getUseCases()).assignPermissions(id, data);
    } catch (error: any) {
        console.error('Assign Permissions Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getRolePermissionsAction = async (id: number) => {
    try {
        return await (await getUseCases()).getRolePermissions(id);
    } catch (error: any) {
        console.error('Get Role Permissions Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const removePermissionFromRoleAction = async (roleId: number, permissionId: number) => {
    try {
        return await (await getUseCases()).removePermissionFromRole(roleId, permissionId);
    } catch (error: any) {
        console.error('Remove Permission From Role Action Error:', error.message);
        throw new Error(error.message);
    }
};
