'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { StorefrontApiRepository } from '../repository/api-repository';
import { StorefrontUseCases } from '../../app/use-cases';
import { CreateStorefrontOrderDTO, CreateClientDTO, CatalogLayout } from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new StorefrontApiRepository(token);
    return new StorefrontUseCases(repository);
}

export const getCatalogAction = async (params?: { page?: number; page_size?: number; search?: string; category?: string; family_id?: number; business_id?: number }) => {
    try {
        return await (await getUseCases()).getCatalog(params);
    } catch (error: any) {
        console.error('Get Catalog Action Error:', error.message);
        return { data: [], total: 0, page: 1, page_size: 12, total_pages: 0 };
    }
};

export const getCatalogFiltersAction = async (businessId?: number) => {
    try {
        return await (await getUseCases()).getCatalogFilters(businessId);
    } catch (error: any) {
        console.error('Get Catalog Filters Action Error:', error.message);
        return { categories: [], families: [], layout: { columns: 4, rows: 3 }, banner: { enabled: false, image_url: '' } };
    }
};

export const updateCatalogLayoutAction = async (layout: CatalogLayout, businessId?: number) => {
    try {
        const result = await (await getUseCases()).updateCatalogLayout(layout, businessId);
        return { success: true as const, layout: result };
    } catch (error: any) {
        console.error('Update Catalog Layout Action Error:', error.message);
        return { success: false as const, message: error.message || 'Error al guardar el diseno del catalogo' };
    }
};

export const updateCatalogBannerAction = async (enabled: boolean, businessId?: number) => {
    try {
        const result = await (await getUseCases()).updateCatalogBanner(enabled, businessId);
        return { success: true as const, banner: result };
    } catch (error: any) {
        console.error('Update Catalog Banner Action Error:', error.message);
        return { success: false as const, message: error.message || 'Error al guardar el banner' };
    }
};

export const uploadCatalogBannerImageAction = async (formData: FormData, businessId?: number) => {
    try {
        const result = await (await getUseCases()).uploadCatalogBannerImage(formData, businessId);
        return { success: true as const, banner: result };
    } catch (error: any) {
        console.error('Upload Catalog Banner Action Error:', error.message);
        return { success: false as const, message: error.message || 'Error al subir la imagen del banner' };
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

export const createClientAction = async (data: CreateClientDTO, businessId?: number) => {
    try {
        const result = await (await getUseCases()).createClient(data, businessId);
        return { success: true as const, ...result };
    } catch (error: any) {
        console.error('Create Client Action Error:', error.message);
        return { success: false as const, message: error.message || 'Error al crear el cliente' };
    }
};

export const getClientsAction = async (params?: { page?: number; page_size?: number; business_id?: number }) => {
    try {
        return await (await getUseCases()).getClients(params);
    } catch (error: any) {
        console.error('Get Clients Action Error:', error.message);
        return { data: [], total: 0, page: 1, page_size: 20, total_pages: 0 };
    }
};
