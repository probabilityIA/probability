'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { OrderApiRepository } from '../repository/api-repository';
import { OrderUseCases } from '../../app/use-cases';
import {
    GetOrdersParams,
    CreateOrderDTO,
    UpdateOrderDTO
} from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
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

export const getOrderRawAction = async (id: string) => {
    try {
        return await (await getUseCases()).getOrderRaw(id);
    } catch (error: any) {
        console.error('Get Order Raw Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getAIRecommendationAction = async (origin: string, destination: string) => {
    try {
        return await (await getUseCases()).getAIRecommendation(origin, destination);
    } catch (error: any) {
        console.warn('AI Recommendation no disponible:', error.message);
        return null;
    }
};
