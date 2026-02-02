'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { OrderApiRepository } from '../repository/api-repository';
import { OrderUseCases } from '../../app/use-cases';
import {
    GetOrdersParams,
    CreateOrderDTO,
    UpdateOrderDTO
} from '../../domain/types';

/**
 * Helper para obtener use cases con token desde cookie o parámetro explícito
 * @param explicitToken Token opcional para iframes donde cookies están bloqueadas
 */
async function getUseCases(explicitToken?: string | null) {
    const token = await getAuthToken(explicitToken);
    const repository = new OrderApiRepository(token);
    return new OrderUseCases(repository);
}

/**
 * Get orders action
 * @param params Query parameters
 * @param token Token opcional para iframes (donde cookies están bloqueadas)
 */
export const getOrdersAction = async (params?: GetOrdersParams, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getOrders(params);
    } catch (error: any) {
        console.error('Get Orders Action Error:', error.message);
        throw new Error(error.message);
    }
};

/**
 * Get order by ID action
 * @param id Order ID
 * @param token Token opcional para iframes (donde cookies están bloqueadas)
 */
export const getOrderByIdAction = async (id: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getOrderById(id);
    } catch (error: any) {
        console.error('Get Order By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createOrderAction = async (data: CreateOrderDTO, token?: string | null) => {
    try {
        return await (await getUseCases(token)).createOrder(data);
    } catch (error: any) {
        console.error('Create Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateOrderAction = async (id: string, data: UpdateOrderDTO, token?: string | null) => {
    try {
        return await (await getUseCases(token)).updateOrder(id, data);
    } catch (error: any) {
        console.error('Update Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteOrderAction = async (id: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).deleteOrder(id);
    } catch (error: any) {
        console.error('Delete Order Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderRawAction = async (id: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getOrderRaw(id);
    } catch (error: any) {
        console.error('Get Order Raw Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getAIRecommendationAction = async (origin: string, destination: string, token?: string | null) => {
    try {
        return await (await getUseCases(token)).getAIRecommendation(origin, destination);
    } catch (error: any) {
        // No lanzar error, retornar null para que el componente maneje silenciosamente
        console.warn('AI Recommendation no disponible:', error.message);
        return null;
    }
};
