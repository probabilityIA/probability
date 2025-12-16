'use server';

import { cookies } from 'next/headers';
import { ActionApiRepository } from '../repository/api-repository';
import { ActionUseCases } from '../../app/use-cases';
import {
    GetActionsParams,
    CreateActionDTO,
    UpdateActionDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new ActionApiRepository(token);
    return new ActionUseCases(repository);
}

export const getActionsAction = async (params?: GetActionsParams) => {
    try {
        return await (await getUseCases()).getActions(params);
    } catch (error: any) {
        console.error('Get Actions Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getActionByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getActionById(id);
    } catch (error: any) {
        console.error('Get Action By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createActionAction = async (data: CreateActionDTO) => {
    try {
        return await (await getUseCases()).createAction(data);
    } catch (error: any) {
        console.error('Create Action Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateActionAction = async (id: number, data: UpdateActionDTO) => {
    try {
        return await (await getUseCases()).updateAction(id, data);
    } catch (error: any) {
        console.error('Update Action Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteActionAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteAction(id);
    } catch (error: any) {
        console.error('Delete Action Action Error:', error.message);
        throw new Error(error.message);
    }
};
