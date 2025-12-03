'use server';

import { cookies } from 'next/headers';
import { PermissionApiRepository } from '../repository/api-repository';
import { PermissionUseCases } from '../../app/use-cases';
import {
    GetPermissionsParams,
    CreatePermissionDTO,
    UpdatePermissionDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new PermissionApiRepository(token);
    return new PermissionUseCases(repository);
}

export const getPermissionsAction = async (params?: GetPermissionsParams) => {

    try {
        return await (await getUseCases()).getPermissions(params);
    } catch (error: any) {
        console.error('Get Permissions Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getPermissionByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getPermissionById(id);
    } catch (error: any) {
        console.error('Get Permission By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getPermissionsByScopeAction = async (scopeId: number) => {
    try {
        return await (await getUseCases()).getPermissionsByScope(scopeId);
    } catch (error: any) {
        console.error('Get Permissions By Scope Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getPermissionsByResourceAction = async (resource: string) => {
    try {
        return await (await getUseCases()).getPermissionsByResource(resource);
    } catch (error: any) {
        console.error('Get Permissions By Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createPermissionAction = async (data: CreatePermissionDTO) => {
    try {
        return await (await getUseCases()).createPermission(data);
    } catch (error: any) {
        console.error('Create Permission Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updatePermissionAction = async (id: number, data: UpdatePermissionDTO) => {
    try {
        return await (await getUseCases()).updatePermission(id, data);
    } catch (error: any) {
        console.error('Update Permission Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deletePermissionAction = async (id: number) => {
    try {
        return await (await getUseCases()).deletePermission(id);
    } catch (error: any) {
        console.error('Delete Permission Action Error:', error.message);
        throw new Error(error.message);
    }
};
