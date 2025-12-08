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
        console.error('Get Products Action Error:', error.message);
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

export const getProductByIdAction = async (id: string) => {
    try {
        return await (await getUseCases()).getProductById(id);
    } catch (error: any) {
        console.error('Get Product By Id Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const createProductAction = async (data: CreateProductDTO) => {
    try {
        return await (await getUseCases()).createProduct(data);
    } catch (error: any) {
        console.error('Create Product Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const updateProductAction = async (id: string, data: UpdateProductDTO) => {
    try {
        return await (await getUseCases()).updateProduct(id, data);
    } catch (error: any) {
        console.error('Update Product Action Error:', error.message);
        throw new Error(error.message);
    }
};

export const deleteProductAction = async (id: string) => {
    try {
        return await (await getUseCases()).deleteProduct(id);
    } catch (error: any) {
        console.error('Delete Product Action Error:', error.message);
        throw new Error(error.message);
    }
};

// ═══════════════════════════════════════════
// Product-Integration Management Actions
// ═══════════════════════════════════════════

export const addProductIntegrationAction = async (
    productId: string,
    data: AddProductIntegrationDTO
) => {
    try {
        return await (await getUseCases()).addProductIntegration(productId, data);
    } catch (error: any) {
        console.error('Add Product Integration Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al asociar producto con integración',
            data: null
        };
    }
};

export const removeProductIntegrationAction = async (
    productId: string,
    integrationId: number
) => {
    try {
        return await (await getUseCases()).removeProductIntegration(productId, integrationId);
    } catch (error: any) {
        console.error('Remove Product Integration Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al remover integración',
            error: error.message
        };
    }
};

export const getProductIntegrationsAction = async (productId: string) => {
    try {
        return await (await getUseCases()).getProductIntegrations(productId);
    } catch (error: any) {
        console.error('Get Product Integrations Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al obtener integraciones',
            data: [],
            total: 0
        };
    }
};
