'use server';

import { cookies } from 'next/headers';
import { UserApiRepository } from '../repository/api-repository';
import { UserUseCases } from '../../app/use-cases';
import {
    GetUsersParams,
    CreateUserDTO,
    UpdateUserDTO,
    AssignRolesDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new UserApiRepository(token);
    return new UserUseCases(repository);
}

export const getUsersAction = async (params?: GetUsersParams) => {
    try {
        return await (await getUseCases()).getUsers(params);
    } catch (error: any) {
        console.error('Get Users Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getUserByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getUserById(id);
    } catch (error: any) {
        console.error('Get User By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createUserAction = async (data: CreateUserDTO) => {
    try {
        return await (await getUseCases()).createUser(data);
    } catch (error: any) {
        console.error('Create User Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateUserAction = async (id: number, data: UpdateUserDTO) => {
    try {
        return await (await getUseCases()).updateUser(id, data);
    } catch (error: any) {
        console.error('Update User Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteUserAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteUser(id);
    } catch (error: any) {
        console.error('Delete User Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const assignRolesAction = async (id: number, data: AssignRolesDTO) => {
    try {
        return await (await getUseCases()).assignRoles(id, data);
    } catch (error: any) {
        console.error('Assign Roles Action Error:', error.message);
        throw new Error(error.message);
    }
};
