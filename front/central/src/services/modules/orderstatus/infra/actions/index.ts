'use server';

import { cookies } from 'next/headers';
import { OrderStatusMappingApiRepository } from '../repository/api-repository';
import { OrderStatusMappingUseCases } from '../../app/use-cases';
import {
    GetOrderStatusMappingsParams,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new OrderStatusMappingApiRepository(token);
    return new OrderStatusMappingUseCases(repository);
}

export const getOrderStatusMappingsAction = async (params?: GetOrderStatusMappingsParams) => {
    try {
        return await (await getUseCases()).getOrderStatusMappings(params);
    } catch (error: any) {
        console.error('Get Order Status Mappings Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderStatusMappingByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getOrderStatusMappingById(id);
    } catch (error: any) {
        console.error('Get Order Status Mapping By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createOrderStatusMappingAction = async (data: CreateOrderStatusMappingDTO) => {
    try {
        return await (await getUseCases()).createOrderStatusMapping(data);
    } catch (error: any) {
        console.error('Create Order Status Mapping Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateOrderStatusMappingAction = async (id: number, data: UpdateOrderStatusMappingDTO) => {
    try {
        return await (await getUseCases()).updateOrderStatusMapping(id, data);
    } catch (error: any) {
        console.error('Update Order Status Mapping Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteOrderStatusMappingAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteOrderStatusMapping(id);
    } catch (error: any) {
        console.error('Delete Order Status Mapping Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const toggleOrderStatusMappingActiveAction = async (id: number) => {
    try {
        return await (await getUseCases()).toggleOrderStatusMappingActive(id);
    } catch (error: any) {
        console.error('Toggle Order Status Mapping Active Action Error:', error.message);
        throw new Error(error.message);
    }
};
