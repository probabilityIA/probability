'use server';

import { cookies } from 'next/headers';
import { CustomerApiRepository } from '../repository/api-repository';
import { CustomerUseCases } from '../../app/use-cases';
import { GetCustomersParams, PaginationParams, CreateCustomerDTO, UpdateCustomerDTO } from '../../domain/types';

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
        throw new Error(error.message);
    }
};

export const getCustomerByIdAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getCustomerById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createCustomerAction = async (data: CreateCustomerDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createCustomer(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateCustomerAction = async (id: number, data: UpdateCustomerDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateCustomer(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteCustomerAction = async (id: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteCustomer(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerSummaryAction = async (customerId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getCustomerSummary(customerId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerAddressesAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerAddresses(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerProductsAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerProducts(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getCustomerOrderItemsAction = async (customerId: number, params?: PaginationParams) => {
    try {
        return await (await getUseCases()).getCustomerOrderItems(customerId, params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
