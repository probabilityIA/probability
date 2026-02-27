'use server';

import { cookies } from 'next/headers';
import { OrderStatusMappingApiRepository } from '../repository/api-repository';
import { OrderStatusMappingUseCases } from '../../app/use-cases';
import {
    GetOrderStatusMappingsParams,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    OrderStatusInfo,
    CreateOrderStatusDTO,
    UpdateOrderStatusDTO,
    CreateChannelStatusDTO,
    UpdateChannelStatusDTO
} from '../../domain/types';
import { env } from '@/shared/config/env';

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

export const getOrderStatusesAction = async (isActive?: boolean): Promise<{ success: boolean; data: OrderStatusInfo[]; message?: string }> => {
    try {
        return await (await getUseCases()).getOrderStatuses(isActive);
    } catch (error: any) {
        console.error('Get Order Statuses Action Error:', error.message);
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

// ============================================
// CRUD para estados de Probability
// ============================================

export const createOrderStatusAction = async (data: CreateOrderStatusDTO) => {
    try {
        return await (await getUseCases()).createOrderStatus(data);
    } catch (error: any) {
        console.error('Create Order Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderStatusByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getOrderStatusById(id);
    } catch (error: any) {
        console.error('Get Order Status By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateOrderStatusAction = async (id: number, data: UpdateOrderStatusDTO) => {
    try {
        return await (await getUseCases()).updateOrderStatus(id, data);
    } catch (error: any) {
        console.error('Update Order Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteOrderStatusAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteOrderStatus(id);
    } catch (error: any) {
        console.error('Delete Order Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

// ============================================
// Simple Actions - Para Dropdowns/Selectores
// ============================================

// ============================================
// Estados por canal de integración (ecommerce)
// ============================================

export const getEcommerceIntegrationTypesAction = async () => {
    try {
        return await (await getUseCases()).getEcommerceIntegrationTypes();
    } catch (error: any) {
        console.error('Get Ecommerce Integration Types Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getChannelStatusesAction = async (integrationTypeId: number, isActive?: boolean) => {
    try {
        return await (await getUseCases()).getChannelStatuses(integrationTypeId, isActive);
    } catch (error: any) {
        console.error('Get Channel Statuses Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createChannelStatusAction = async (data: CreateChannelStatusDTO) => {
    try {
        return await (await getUseCases()).createChannelStatus(data);
    } catch (error: any) {
        console.error('Create Channel Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateChannelStatusAction = async (id: number, data: UpdateChannelStatusDTO) => {
    try {
        return await (await getUseCases()).updateChannelStatus(id, data);
    } catch (error: any) {
        console.error('Update Channel Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteChannelStatusAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteChannelStatus(id);
    } catch (error: any) {
        console.error('Delete Channel Status Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const getOrderStatusesSimpleAction = async (isActive: boolean = true): Promise<import('../../domain/types').OrderStatusesSimpleResponse> => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || '';

        const url = `${env.API_BASE_URL}/order-statuses/simple?is_active=${isActive}`;

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error('Error al obtener estados de orden');
        }

        return await response.json();
    } catch (error: any) {
        console.error('Get Order Statuses Simple Action Error:', error.message);
        return {
            success: false,
            message: error.message,
            data: [],
        };
    }
};
