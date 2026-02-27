'use server';

import { cookies } from 'next/headers';
import { CustomerApiRepository } from '../repository/api-repository';
import { CustomerUseCases } from '../../app/use-cases';
import { GetCustomersParams, CreateCustomerDTO, UpdateCustomerDTO } from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new CustomerApiRepository(token);
    return new CustomerUseCases(repository);
}

export const getCustomersAction = async (params?: GetCustomersParams) => {
    try {
        return await (await getUseCases()).getCustomers(params);
    } catch (error: any) {
        console.error('Get Customers Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getCustomerByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getCustomerById(id);
    } catch (error: any) {
        console.error('Get Customer By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createCustomerAction = async (data: CreateCustomerDTO) => {
    try {
        return await (await getUseCases()).createCustomer(data);
    } catch (error: any) {
        console.error('Create Customer Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateCustomerAction = async (id: number, data: UpdateCustomerDTO) => {
    try {
        return await (await getUseCases()).updateCustomer(id, data);
    } catch (error: any) {
        console.error('Update Customer Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteCustomerAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteCustomer(id);
    } catch (error: any) {
        console.error('Delete Customer Action Error:', error.message);
        throw new Error(error.message);
    }
};
