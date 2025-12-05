'use server';

import { cookies } from 'next/headers';
import { ResourceApiRepository } from '../repository/api-repository';
import { ResourceUseCases } from '../../app/use-cases';
import {
    GetResourcesParams,
    CreateResourceDTO,
    UpdateResourceDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new ResourceApiRepository(token);
    return new ResourceUseCases(repository);
}

export const getResourcesAction = async (params?: GetResourcesParams) => {
    try {
        return await (await getUseCases()).getResources(params);
    } catch (error: any) {
        console.error('Get Resources Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getResourceByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getResourceById(id);
    } catch (error: any) {
        console.error('Get Resource By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createResourceAction = async (data: CreateResourceDTO) => {
    try {
        return await (await getUseCases()).createResource(data);
    } catch (error: any) {
        console.error('Create Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateResourceAction = async (id: number, data: UpdateResourceDTO) => {
    try {
        return await (await getUseCases()).updateResource(id, data);
    } catch (error: any) {
        console.error('Update Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteResourceAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteResource(id);
    } catch (error: any) {
        console.error('Delete Resource Action Error:', error.message);
        throw new Error(error.message);
    }
};
