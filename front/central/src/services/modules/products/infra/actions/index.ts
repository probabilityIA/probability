'use server';

import { cookies } from 'next/headers';
import { ProductApiRepository } from '../repository/api-repository';
import { ProductUseCases } from '../../app/use-cases';
import {
    GetProductsParams,
    CreateProductDTO,
    UpdateProductDTO,
    AddProductIntegrationDTO
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new ProductApiRepository(token);
    return new ProductUseCases(repository);
}

export const getProductsAction = async (params?: GetProductsParams) => {
    try {
        return await (await getUseCases()).getProducts(params);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener productos',
            data: [],
            total: 0,
            page: params?.page || 1,
            page_size: params?.page_size || 20,
            total_pages: 0
        };
    }
};

export const getProductByIdAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductById(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createProductAction = async (data: CreateProductDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createProduct(data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateProductAction = async (id: string, data: UpdateProductDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateProduct(id, data, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteProductAction = async (id: string, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteProduct(id, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

// ═══════════════════════════════════════════
// Product-Integration Management Actions
// ═══════════════════════════════════════════

export const addProductIntegrationAction = async (
    productId: string,
    data: AddProductIntegrationDTO,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).addProductIntegration(productId, data, businessId);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al asociar producto con integración',
            data: null
        };
    }
};

export const removeProductIntegrationAction = async (
    productId: string,
    integrationId: number,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).removeProductIntegration(productId, integrationId, businessId);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al remover integración',
            error: error.message
        };
    }
};

export const getProductIntegrationsAction = async (productId: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductIntegrations(productId, businessId);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener integraciones',
            data: [],
            total: 0
        };
    }
};
