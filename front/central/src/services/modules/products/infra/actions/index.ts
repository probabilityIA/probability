'use server';

import { cookies } from 'next/headers';
import { ProductApiRepository } from '../repository/api-repository';
import { ProductUseCases } from '../../app/use-cases';
import {
    GetProductsParams,
    GetFamiliesParams,
    CreateProductDTO,
    UpdateProductDTO,
    CreateProductFamilyDTO,
    UpdateProductFamilyDTO,
    AddProductIntegrationDTO,
    UpdateProductIntegrationDTO
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

export const uploadProductImageAction = async (productId: string, formData: FormData, businessId?: number) => {
    try {
        return await (await getUseCases()).uploadProductImage(productId, formData, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al subir imagen', image_url: '' };
    }
};

export const getProductFamiliesAction = async (params?: GetFamiliesParams) => {
    try {
        return await (await getUseCases()).getProductFamilies(params);
    } catch (error: any) {
        return {
            success: false,
            message: error.message || 'Error al obtener familias',
            data: [],
            total: 0,
            page: params?.page || 1,
            page_size: params?.page_size || 20,
            total_pages: 0
        };
    }
};

export const getProductFamilyByIdAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductFamilyById(familyId, businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getFamilyVariantsAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).getFamilyVariants(familyId, businessId);
    } catch (error: any) {
        return { success: false, data: [] };
    }
};

export const createProductFamilyAction = async (data: CreateProductFamilyDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).createProductFamily(data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al crear familia', data: null };
    }
};

export const updateProductFamilyAction = async (familyId: number, data: UpdateProductFamilyDTO, businessId?: number) => {
    try {
        return await (await getUseCases()).updateProductFamily(familyId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al actualizar familia', data: null };
    }
};

export const deleteProductFamilyAction = async (familyId: number, businessId?: number) => {
    try {
        return await (await getUseCases()).deleteProductFamily(familyId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al eliminar familia', error: error.message };
    }
};

export const addProductIntegrationAction = async (
    productId: string,
    data: AddProductIntegrationDTO,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).addProductIntegration(productId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al asociar integración', data: null };
    }
};

export const updateProductIntegrationAction = async (
    productId: string,
    integrationId: number,
    data: UpdateProductIntegrationDTO,
    businessId?: number
) => {
    try {
        return await (await getUseCases()).updateProductIntegration(productId, integrationId, data, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al actualizar mapping', data: null };
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
        return { success: false, message: error.message || 'Error al remover integración', error: error.message };
    }
};

export const getProductIntegrationsAction = async (productId: string, businessId?: number) => {
    try {
        return await (await getUseCases()).getProductIntegrations(productId, businessId);
    } catch (error: any) {
        return { success: false, message: error.message || 'Error al obtener integraciones', data: [], total: 0 };
    }
};
