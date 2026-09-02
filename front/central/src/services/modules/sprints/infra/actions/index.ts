'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { SprintApiRepository } from '../repository/api-repository';
import { SprintUseCases } from '../../app/use-cases';
import {
    CreateSprintDTO,
    UpdateSprintDTO,
    ListSprintsParams,
    SprintStatus,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new SprintApiRepository(token);
    return new SprintUseCases(repository);
}

export const listSprintsAction = async (params?: ListSprintsParams) => {
    return (await getUseCases()).list(params);
};

export const getSprintAction = async (id: number) => {
    return (await getUseCases()).get(id);
};

export const createSprintAction = async (data: CreateSprintDTO) => {
    const r = await (await getUseCases()).create(data);
    revalidatePath('/tickets');
    return r;
};

export const updateSprintAction = async (id: number, data: UpdateSprintDTO) => {
    const r = await (await getUseCases()).update(id, data);
    revalidatePath('/tickets');
    return r;
};

export const deleteSprintAction = async (id: number) => {
    await (await getUseCases()).remove(id);
    revalidatePath('/tickets');
};

export const changeSprintStatusAction = async (id: number, status: SprintStatus) => {
    const r = await (await getUseCases()).changeStatus(id, status);
    revalidatePath('/tickets');
    return r;
};
