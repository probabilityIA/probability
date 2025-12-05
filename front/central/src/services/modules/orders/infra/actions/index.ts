'use server';

import { cookies } from 'next/headers';
import { OrderApiRepository } from '../repository/api-repository';
import { OrderUseCases } from '../../app/use-cases';
import {
    GetOrdersParams,
    CreateOrderDTO,
    UpdateOrderDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new OrderApiRepository(token);
    return new OrderUseCases(repository);
}

export const getOrdersAction = async (params?: GetOrdersParams) => {
    try {
        return await (await getUseCases()).getOrders(params);
    } catch (error: any) {
        console.error('Get Orders Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderByIdAction = async (id: string) => {
    try {
        return await (await getUseCases()).getOrderById(id);
    } catch (error: any) {
        console.error('Get Order By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createOrderAction = async (data: CreateOrderDTO) => {
    try {
        return await (await getUseCases()).createOrder(data);
    } catch (error: any) {
        console.error('Create Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateOrderAction = async (id: string, data: UpdateOrderDTO) => {
    try {
        return await (await getUseCases()).updateOrder(id, data);
    } catch (error: any) {
        console.error('Update Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteOrderAction = async (id: string) => {
    try {
        return await (await getUseCases()).deleteOrder(id);
    } catch (error: any) {
        console.error('Delete Order Action Error:', error.message);
        throw new Error(error.message);
    }
};
