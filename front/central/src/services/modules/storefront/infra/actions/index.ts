'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { StorefrontApiRepository } from '../repository/api-repository';
import { StorefrontUseCases } from '../../app/use-cases';
import { CreateStorefrontOrderDTO, RegisterDTO } from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new StorefrontApiRepository(token);
    return new StorefrontUseCases(repository);
}

export const getCatalogAction = async (params?: { page?: number; page_size?: number; search?: string; category?: string; business_id?: number }) => {
    try {
        return await (await getUseCases()).getCatalog(params);
    } catch (error: any) {
        console.error('Get Catalog Action Error:', error.message);
        return { data: [], total: 0, page: 1, page_size: 12, total_pages: 0 };
    }
};

export const getProductAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProduct(id, businessId);
    } catch (error: any) {
        console.error('Get Product Action Error:', error.message);
        return null;
    }
};

export const createOrderAction = async (data: CreateStorefrontOrderDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createOrder(data, businessId);
    } catch (error: any) {
        console.error('Create Order Action Error:', error.message);
        return { success: false, message: error.message || 'Error al crear pedido' };
    }
};

export const getOrdersAction = async (params?: { page?: number; page_size?: number; business_id?: number }) => {
    try {
        return await (await getUseCases()).getOrders(params);
    } catch (error: any) {
        console.error('Get Orders Action Error:', error.message);
        return { data: [], total: 0, page: 1, page_size: 10, total_pages: 0 };
    }
};

export const getOrderAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getOrder(id, businessId);
    } catch (error: any) {
        console.error('Get Order Action Error:', error.message);
        return null;
    }
};

export const registerAction = async (data: RegisterDTO) => {
    // Registration doesn't need auth token
    const repository = new StorefrontApiRepository();
    const useCases = new StorefrontUseCases(repository);
    try {
        return await useCases.register(data);
    } catch (error: any) {
        console.error('Register Action Error:', error.message);
        return { success: false, message: error.message || 'Error al registrarse' };
    }
};
